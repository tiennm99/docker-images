---
phase: 2
title: Compose orchestration and docs
status: completed
priority: P2
effort: 1h
dependencies:
  - 1
---

# Phase 2: Compose orchestration and docs

## Overview

Wire the Scribe daemon and the Go sender (Phase 1) into a single `docker compose up`
demo, bind-mount logs to a host-visible folder, and document run + verify steps. End
result: one command produces real log files under `scribe/example/logs/`.

## Requirements

- Functional: `docker compose up` (from `scribe/example/`) starts Scribe with the repo's
  `config/scribe.conf`, runs the sender once, and writes files to host `./logs`.
- Functional: logs are inspectable on the host without entering the container.
- Non-functional: scribe service uses published `ghcr.io/tiennm99/scribe:2.2` (no local build).
- Non-functional: example README is self-contained (run, verify, troubleshoot).

## Architecture

### docker-compose.yml

```yaml
services:
  scribe:
    image: ghcr.io/tiennm99/scribe:2.2
    ports:
      - "1463:1463"
    volumes:
      - ../config/scribe.conf:/etc/scribe/scribe.conf:ro
      - ./logs:/var/log/scribe
    restart: unless-stopped

  sender:
    build: ./sender
    depends_on:
      - scribe
    environment:
      - SCRIBE_HOST=scribe
      - SCRIBE_PORT=1463
      - SCRIBE_CATEGORY=example
      - SCRIBE_COUNT=5
```

Notes:
- Sender reaches Scribe over the compose network by service name `scribe:1463` (internal),
  independent of the host `1463:1463` mapping.
- `depends_on` only orders start, not readiness — the Phase 1 retry loop covers the gap.
- Sender is one-shot: it sends then exits; scribe keeps running (`restart: unless-stopped`).

### Expected on-disk result

```
scribe/example/logs/
└── example/
    └── example-<YYYY-MM-DD>_00000        # contains the 5 sent messages
```

(Category `example` → subdir under the default catch-all store's `file_path=/var/log/scribe`.)

## Related Code Files

- Create: `scribe/example/docker-compose.yml`
- Create: `scribe/example/logs/.gitkeep` — keep the bind-mount dir in git, empty.
- Create: `scribe/example/README.md` — run/verify/troubleshoot for the example.
- Modify: `scribe/README.md` — add a short "Runnable example" pointer linking to `example/`.

### scribe/example/README.md outline

1. **What it does** — one command boots Scribe + sends test messages → log files appear.
2. **Run** — `cd scribe/example && docker compose up --build`.
3. **Verify** — show `scribe/example/logs/example/example-*_00000` and `cat` a line.
4. **Customize** — env vars (`SCRIBE_CATEGORY`, `SCRIBE_COUNT`), pointer to `config/scribe.conf`.
5. **Troubleshoot** — empty logs dir = no messages sent yet; Linux bind-mount permissions caveat.
6. **Stop** — `docker compose down`.

## Implementation Steps

1. Create `scribe/example/docker-compose.yml` per spec above.
2. Create `scribe/example/logs/.gitkeep`.
3. Write `scribe/example/README.md` per outline.
4. Add a "Runnable example" subsection to `scribe/README.md` linking to `example/`.
5. End-to-end test: `cd scribe/example && docker compose up --build`; confirm a log file
   appears under `logs/example/` containing the sent messages; `docker compose down`.

## Success Criteria

- [ ] `docker compose up --build` starts both services without manual steps.
- [ ] After the sender exits, `scribe/example/logs/example/example-*_00000` exists on host.
- [ ] That file contains the test messages sent by the Go sender.
- [ ] `scribe/example/README.md` documents run, verify, and the Linux permission caveat.
- [ ] `scribe/README.md` links to the example.

## Risk Assessment

- **Bind-mount permissions (native Linux)** — Scribe runs as non-root user `scribe`
  (`scribe/Dockerfile:107-110`); host `./logs` may reject its writes. On Docker Desktop
  (Windows/macOS) this is normally fine. Mitigation: document `chmod 777 ./logs` (or matching
  uid) in the README troubleshoot section.
- **Stale relative paths** — `../config/scribe.conf` and `./sender` build context are relative
  to `scribe/example/`; compose must be run from that dir. Documented in README run step.
- **Empty logs confusion** — if the sender fails silently, logs dir stays empty. Mitigation:
  sender logs success/failure to stdout (`docker compose` output) and exits non-zero on failure.
