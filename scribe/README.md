# Scribe

Facebook's legacy **Scribe** log aggregation daemon, packaged as a Docker image.

Scribe is a server for aggregating log data streamed in real time from a large number of servers. It accepts messages over a Thrift RPC interface and writes them to configurable back-end stores (files, network relays, etc.).

> **Note:** Scribe was open-sourced by Facebook and is now archived at
> [facebookarchive/scribe](https://github.com/facebookarchive/scribe).
> This image exists because no official Docker image is available.

---

## Quick Start

```bash
docker run -d \
  -p 1463:1463 \
  -v $(pwd)/logs:/var/log/scribe \
  ghcr.io/tiennm99/scribe:2.2
```

Scribe will listen on port **1463** (Thrift binary protocol) and write incoming log messages to `/var/log/scribe` inside the container.

---

## Custom Configuration

Mount your own `scribe.conf` over the default one:

```bash
docker run -d \
  -p 1463:1463 \
  -v $(pwd)/my-scribe.conf:/etc/scribe/scribe.conf:ro \
  -v $(pwd)/logs:/var/log/scribe \
  ghcr.io/tiennm99/scribe:2.2
```

Or set `SCRIBE_CONFIG` to point to a different path inside the container:

```bash
docker run -d \
  -p 1463:1463 \
  -e SCRIBE_CONFIG=/opt/scribe/custom.conf \
  -v $(pwd)/custom.conf:/opt/scribe/custom.conf:ro \
  -v $(pwd)/logs:/var/log/scribe \
  ghcr.io/tiennm99/scribe:2.2
```

---

## Docker Compose

```yaml
services:
  scribe:
    image: ghcr.io/tiennm99/scribe:2.2
    ports:
      - "1463:1463"
    volumes:
      - scribe_logs:/var/log/scribe
      - ./config/scribe.conf:/etc/scribe/scribe.conf:ro
    restart: unless-stopped

volumes:
  scribe_logs:
```

---

## Configuration Reference

The default config at [`config/scribe.conf`](config/scribe.conf) demonstrates the most common options:

| Key | Default | Description |
|-----|---------|-------------|
| `port` | `1463` | TCP port scribed listens on |
| `max_msg_per_second` | `100000` | Rate limit across all categories (0 = unlimited) |
| `max_msg_size` | `16777216` | Max message size in bytes |
| `check_interval` | `3` | Seconds between store health checks |

### Store Types

| Type | Description |
|------|-------------|
| `file` | Write messages to rotating files on disk |
| `network` | Relay messages to an upstream Scribe instance |
| `buffer` | Buffer messages in memory with a file fallback |
| `null` | Discard messages (useful for testing) |
| `multi` | Fan-out to multiple stores simultaneously |

Full configuration documentation:
https://github.com/facebookarchive/scribe/wiki/Scribe-Configuration

---

## Sending Log Messages

Use the Scribe Thrift client library for your language, or a simple Python example:

```python
from scribe import scribe
from thrift.transport import TSocket, TTransport
from thrift.protocol import TBinaryProtocol

socket    = TSocket.TSocket(host="localhost", port=1463)
transport = TTransport.TFramedTransport(socket)
protocol  = TBinaryProtocol.TBinaryProtocol(trans=transport, strictRead=False, strictWrite=False)
client    = scribe.Client(protocol)

transport.open()
log_entry = scribe.LogEntry(category="myapp", message="Hello, Scribe!")
result    = client.Log(messages=[log_entry])
transport.close()
```

---

## Volumes

| Path | Purpose |
|------|---------|
| `/var/log/scribe` | Log output directory (mount a host path or named volume) |
| `/etc/scribe` | Configuration directory |

---

## Ports

| Port | Protocol | Description |
|------|----------|-------------|
| `1463` | TCP | Scribe Thrift RPC interface |

---

## Building Locally

```bash
docker build -t scribe:2.2 .
```

The Dockerfile uses a two-stage build:
1. **builder** – compiles Apache Thrift 0.9.1 and Scribe from source on Ubuntu 14.04
2. **runtime** – copies only the compiled binary and required shared libraries

---

## Credits

- https://github.com/facebookarchive/scribe
- https://archive.apache.org/dist/thrift/0.9.1/
