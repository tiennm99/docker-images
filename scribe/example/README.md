# Scribe Runnable Example

A one-command demo that proves the Scribe config actually writes logs.

`docker compose up --build` starts the Scribe daemon (using the repo's default
[`config/scribe.conf`](../config/scribe.conf)) plus a tiny **Go sender** that connects over
Thrift, pushes a few messages, and exits. The resulting log files appear on your host under
[`./logs`](./logs).

> Scribe is receive-only — it writes **nothing** until a client sends it messages. The Go
> sender is what makes this example produce log files.

---

## Run

```bash
cd scribe/example
docker compose up --build
```

You'll see the sender connect and report:

```
sender-1  | sent 5 messages to category "example", result=OK
```

The sender then exits; the Scribe daemon keeps running. Press `Ctrl+C`, then `docker compose down` to stop it.

---

## Verify

Log files are written on the host under `./logs/<category>/`:

```bash
ls logs/example/
cat logs/example/example-*_00000
```

You should see the five `Hello from Go sender — message N/5` lines.

---

## Customize

Adjust the sender via environment variables in [`docker-compose.yml`](./docker-compose.yml):

| Env | Default | Purpose |
|-----|---------|---------|
| `SCRIBE_CATEGORY` | `example` | Message category → log subdirectory name |
| `SCRIBE_COUNT` | `5` | Number of messages to send |
| `SCRIBE_HOST` | `scribe` | Daemon hostname (compose service name) |
| `SCRIBE_PORT` | `1463` | Thrift port |

To change *where/how* Scribe stores logs (rotation, store type, per-category routing), edit
[`config/scribe.conf`](../config/scribe.conf) — see the
[Configuration Reference](../README.md#configuration-reference).

The sender ([`sender/main.go`](./sender/main.go)) hand-encodes the Thrift framed binary
`Log()` call using only the Go standard library — no Apache Thrift dependency, no generated
bindings. For the higher-level client API, see the Python snippet in the
[main README](../README.md#sending-log-messages).

---

## Troubleshoot

- **`logs/` stays empty** — no messages were written yet. Check the `sender` output in the
  compose logs; it exits non-zero and logs the cause if the send failed.
- **Permission denied writing logs (native Linux hosts)** — Scribe runs as a non-root user
  inside the container, which may not be able to write to the bind-mounted host directory.
  Fix by loosening the local dir before running:

  ```bash
  chmod 777 logs
  ```

  On Docker Desktop (Windows/macOS) this is normally not needed.
- **Sender can't connect** — it retries for ~10s on startup. If Scribe takes longer to boot,
  re-run `docker compose up` or increase the retry window in `sender/main.go`.

---

## Stop

```bash
docker compose down
```
