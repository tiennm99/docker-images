// Command scribe-example-sender pushes a few LogEntry messages to a Scribe
// daemon to demonstrate that the example writes real log files.
//
// It speaks Scribe's Thrift RPC by hand-encoding the framed binary `Log()` call
// using only the Go standard library - no Apache Thrift dependency and no
// generated bindings. The Scribe Thrift contract:
//
//	enum ResultCode { OK = 0, TRY_LATER = 1 }
//	struct LogEntry { 1: string category; 2: string message; }
//	service scribe { ResultCode Log(1: list<LogEntry> messages); }
package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

// Thrift binary protocol type IDs and message types.
const (
	tSTOP   = 0
	tI32    = 8
	tSTRING = 11
	tSTRUCT = 12
	tLIST   = 15

	msgCALL = 1
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

// --- TBinaryProtocol writers (all integers big-endian / network order) ---

func writeI16(b *bytes.Buffer, v int16) {
	var tmp [2]byte
	binary.BigEndian.PutUint16(tmp[:], uint16(v))
	b.Write(tmp[:])
}

func writeI32(b *bytes.Buffer, v int32) {
	var tmp [4]byte
	binary.BigEndian.PutUint32(tmp[:], uint32(v))
	b.Write(tmp[:])
}

func writeString(b *bytes.Buffer, s string) {
	writeI32(b, int32(len(s)))
	b.WriteString(s)
}

func writeFieldBegin(b *bytes.Buffer, typeID byte, id int16) {
	b.WriteByte(typeID)
	writeI16(b, id)
}

// buildLogFrame builds a framed Thrift `Log()` CALL carrying one LogEntry per
// message, all under the same category.
//
// Uses the non-strict TBinaryProtocol message header (string name, byte type,
// i32 seqid) to match the README Python client (strictWrite=False); the Scribe
// C++ server reads with strict_read=false and accepts it.
func buildLogFrame(category string, messages []string) []byte {
	var p bytes.Buffer

	// messageBegin (non-strict): name, type, seqid
	writeString(&p, "Log")
	p.WriteByte(msgCALL)
	writeI32(&p, 0)

	// Log_args struct: field 1 = list<LogEntry> messages
	writeFieldBegin(&p, tLIST, 1)
	p.WriteByte(tSTRUCT) // list element type
	writeI32(&p, int32(len(messages)))
	for _, m := range messages {
		writeFieldBegin(&p, tSTRING, 1) // LogEntry.category
		writeString(&p, category)
		writeFieldBegin(&p, tSTRING, 2) // LogEntry.message
		writeString(&p, m)
		p.WriteByte(tSTOP) // end LogEntry
	}
	p.WriteByte(tSTOP) // end Log_args

	// Frame the payload with a 4-byte big-endian length prefix (TFramedTransport).
	var frame bytes.Buffer
	writeI32(&frame, int32(p.Len()))
	frame.Write(p.Bytes())
	return frame.Bytes()
}

// readResultCode reads one framed Thrift reply and extracts the ResultCode.
// Handles both strict and non-strict message headers, since the server's write
// strictness is independent of ours.
func readResultCode(r io.Reader) (int32, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return -1, fmt.Errorf("read frame length: %w", err)
	}
	frame := make([]byte, binary.BigEndian.Uint32(lenBuf[:]))
	if _, err := io.ReadFull(r, frame); err != nil {
		return -1, fmt.Errorf("read frame body: %w", err)
	}
	br := bytes.NewReader(frame)

	// messageBegin: detect strict vs non-strict from the first i32.
	var i32buf [4]byte
	if _, err := io.ReadFull(br, i32buf[:]); err != nil {
		return -1, fmt.Errorf("read message header: %w", err)
	}
	if first := int32(binary.BigEndian.Uint32(i32buf[:])); first < 0 {
		// strict: first = version|type, then string name, then i32 seqid
		var nameLen [4]byte
		if _, err := io.ReadFull(br, nameLen[:]); err != nil {
			return -1, err
		}
		if _, err := io.CopyN(io.Discard, br, int64(binary.BigEndian.Uint32(nameLen[:]))); err != nil {
			return -1, err
		}
		if _, err := io.ReadFull(br, i32buf[:]); err != nil { // seqid
			return -1, err
		}
	} else {
		// non-strict: first = name length, then name, then i8 type, then i32 seqid
		if _, err := io.CopyN(io.Discard, br, int64(first)); err != nil {
			return -1, err
		}
		if _, err := br.ReadByte(); err != nil { // message type
			return -1, err
		}
		if _, err := io.ReadFull(br, i32buf[:]); err != nil { // seqid
			return -1, err
		}
	}

	// Log_result struct: first field is the return value (id 0, type I32).
	typeID, err := br.ReadByte()
	if err != nil {
		return -1, err
	}
	if typeID == tSTOP {
		return 0, nil // no return field present; treat as OK
	}
	var idBuf [2]byte
	if _, err := io.ReadFull(br, idBuf[:]); err != nil {
		return -1, err
	}
	if typeID != tI32 {
		return -1, fmt.Errorf("unexpected reply field type %d", typeID)
	}
	if _, err := io.ReadFull(br, i32buf[:]); err != nil {
		return -1, err
	}
	return int32(binary.BigEndian.Uint32(i32buf[:])), nil
}

// dialWithRetry retries the initial connection so a cold `docker compose up`
// does not fail before the Scribe daemon is listening.
func dialWithRetry(addr string, attempts int, interval time.Duration) (net.Conn, error) {
	var lastErr error
	for i := 0; i < attempts; i++ {
		conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
		if err == nil {
			return conn, nil
		}
		lastErr = err
		log.Printf("scribe not reachable at %s (attempt %d/%d): %v", addr, i+1, attempts, err)
		time.Sleep(interval)
	}
	return nil, lastErr
}

func main() {
	host := env("SCRIBE_HOST", "scribe")
	port := env("SCRIBE_PORT", "1463")
	category := env("SCRIBE_CATEGORY", "example")
	count := envInt("SCRIBE_COUNT", 5)

	addr := net.JoinHostPort(host, port)
	conn, err := dialWithRetry(addr, 10, 1*time.Second)
	if err != nil {
		log.Fatalf("could not connect to scribe at %s: %v", addr, err)
	}
	defer conn.Close()

	// Bound the whole exchange so a wedged server that accepts the connection
	// but never replies cannot hang this one-shot sender indefinitely.
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	messages := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		messages = append(messages, fmt.Sprintf("Hello from Go sender - message %d/%d", i, count))
	}

	if _, err := conn.Write(buildLogFrame(category, messages)); err != nil {
		log.Fatalf("send failed: %v", err)
	}

	rc, err := readResultCode(bufio.NewReader(conn))
	if err != nil {
		log.Fatalf("read reply failed: %v", err)
	}

	switch rc {
	case 0:
		log.Printf("sent %d messages to category %q, result=OK", count, category)
	case 1:
		log.Printf("sent %d messages to category %q, result=TRY_LATER (scribe busy)", count, category)
		os.Exit(2)
	default:
		log.Printf("sent %d messages to category %q, unexpected result code=%d", count, category, rc)
		os.Exit(3)
	}
}
