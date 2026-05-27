---
phase: 1
title: Go raw-socket sender
status: completed
priority: P2
effort: 2h
dependencies: []
---

# Phase 1: Go raw-socket sender

## Overview

A tiny Go program that connects to the Scribe daemon over TCP and sends `LogEntry`
messages by hand-encoding the Thrift framed binary `Log()` RPC — stdlib only, no Apache
Thrift dependency, no generated bindings. Compiles to a small static binary in a
multi-stage Docker build.

## Requirements

- Functional: connect to `SCRIBE_HOST:SCRIBE_PORT`, send `SCRIBE_COUNT` messages under
  category `SCRIBE_CATEGORY`, parse the `ResultCode` reply, log OK/TRY_LATER, exit 0 on success.
- Functional: retry the initial connect (~10s, ~1s interval) so a cold `docker compose up`
  does not fail before Scribe is listening.
- Non-functional: no external Go modules (stdlib `net`, `encoding/binary`, `os`, `time`, `log`).
- Non-functional: static binary (`CGO_ENABLED=0`); minimal runtime image.

## Architecture

### Scribe Thrift contract (from facebookarchive/scribe IDL)

```thrift
enum ResultCode { OK = 0, TRY_LATER = 1 }
struct LogEntry { 1: string category; 2: string message; }
service scribe { ResultCode Log(1: list<LogEntry> messages); }
```

### Wire format to produce

- **Transport:** `TFramedTransport` → 4-byte big-endian frame length prefixes the payload.
- **Protocol:** `TBinaryProtocol`, **non-strict** message header (matches README client which
  uses `strictWrite=False`; Scribe C++ server reads with `strict_read=false`, accepts it).

Thrift type IDs used: `STOP=0`, `I32=8`, `STRING=11`, `STRUCT=12`, `LIST=15`. All ints big-endian.

**Request payload (`Log` CALL):**
1. Message begin (non-strict): `string "Log"` (i32 len + bytes) · `i8 msgType=1 (CALL)` · `i32 seqid=0`
2. Args struct `Log_args`:
   - field begin: `i8 type=15 (LIST)` · `i16 id=1`
   - list begin: `i8 elemType=12 (STRUCT)` · `i32 size=N`
   - per `LogEntry`:
     - field: `i8 11 (STRING)` · `i16 1` · writeString(category)
     - field: `i8 11 (STRING)` · `i16 2` · writeString(message)
     - `i8 0 (STOP)`
   - args struct: `i8 0 (STOP)`
3. Prepend `i32` total length → frame.

`writeString(s)` = `i32 len(s)` + UTF-8 bytes.

**Response parse (minimal):**
1. Read 4-byte frame length, then that many bytes.
2. Message begin (non-strict): i32 name-len + name + `i8 msgType` + `i32 seqid`.
3. Field begin `i8 type` · `i16 id`; for return value id=0 type=I32 → read `i32` ResultCode.
4. Map 0→OK, 1→TRY_LATER; anything else → log raw + treat non-zero as failure.

### Connect-retry

Loop `net.DialTimeout("tcp", addr, 2s)` up to 10 attempts, 1s between failures; fatal after.

## Related Code Files

- Create: `scribe/example/sender/main.go` — sender program (encode/send/parse + retry).
- Create: `scribe/example/sender/go.mod` — `module scribe-example-sender` / `go 1.22`, stdlib only.
- Create: `scribe/example/sender/Dockerfile` — multi-stage build.

### Dockerfile (multi-stage)

```dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod .
COPY main.go .
RUN CGO_ENABLED=0 go build -o /sender .

FROM alpine:3.20
COPY --from=build /sender /sender
ENTRYPOINT ["/sender"]
```

### main.go env inputs (with defaults)

| Env | Default | Purpose |
|-----|---------|---------|
| `SCRIBE_HOST` | `scribe` | daemon hostname (compose service name) |
| `SCRIBE_PORT` | `1463` | Thrift port |
| `SCRIBE_CATEGORY` | `example` | message category → log subdir |
| `SCRIBE_COUNT` | `5` | number of LogEntry messages to send |

## Implementation Steps

1. Create `scribe/example/sender/go.mod` (module name + `go 1.22`).
2. Write `main.go`:
   - read env with defaults; build `bytes.Buffer` helpers for `i8/i16/i32/string`.
   - `buildLogFrame(entries)` → payload then 4-byte length prefix.
   - connect-retry loop → write frame → read+parse response frame → log ResultCode.
   - `os.Exit(0)` on OK; non-zero exit on connect failure or non-OK.
3. Write the multi-stage `Dockerfile`.
4. Local sanity build: `docker build -t scribe-sender scribe/example/sender` (verify compiles).

## Success Criteria

- [ ] `docker build scribe/example/sender` succeeds (binary compiles, no module fetch).
- [ ] Run against a live scribe container: sender logs `sent N messages, result=OK` and exits 0.
- [ ] Sender survives a cold start (scribe not yet listening) via the retry loop.

## Risk Assessment

- **Hand-encoded Thrift correctness** — wire format above is the stable 0.9.1 binary layout;
  keep helpers tiny and comment each field. Mitigation: validate end-to-end in Phase 2 by
  confirming the written log file contains the exact messages.
- **Strict vs non-strict header** — if the server unexpectedly rejects non-strict, switch
  message-begin to strict (`i32 version=0x80010001`, name, `i32 seqid`). Documented here as
  the fallback; default is non-strict to match the README client.
