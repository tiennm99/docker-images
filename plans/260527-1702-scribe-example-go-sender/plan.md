---
title: Scribe runnable example with Go test sender
description: >-
  Self-contained scribe/example/ that boots the Scribe daemon and a Go sender so
  docker compose up writes visible log files to host ./logs
status: completed
priority: P2
branch: main
tags:
  - scribe
  - example
  - docker-compose
  - golang
  - thrift
blockedBy: []
blocks: []
created: '2026-05-27T10:27:22.391Z'
createdBy: 'ck:plan'
source: skill
---

# Scribe runnable example with Go test sender

## Overview

Add `scribe/example/` — a one-command demo proving the Scribe config actually writes logs.
`docker compose up` starts the Scribe daemon (published image, mounting the repo's
`config/scribe.conf`) plus a tiny Go sender that connects over Thrift, pushes a few
`LogEntry` messages, and exits. Resulting log files land in host-visible `./logs`.

Scribe is receive-only: it writes nothing until a client sends Thrift messages. The Go
sender is what makes the example "write logs". The sender hand-encodes the Thrift framed
binary `Log()` call using only the Go stdlib (no codegen, no external modules) — chosen for
the smallest, most deterministic build (approach G1 from brainstorm).

## Phases

| Phase | Name | Status |
|-------|------|--------|
| 1 | [Go raw-socket sender](./phase-01-go-raw-socket-sender.md) | Completed |
| 2 | [Compose orchestration and docs](./phase-02-compose-orchestration-and-docs.md) | Completed |

## Key Decisions (from brainstorm)

- Sender language: **Go**, raw-socket stdlib only (no Apache Thrift lib, no generated bindings).
- Logs: **host bind mount** `scribe/example/logs` → container `/var/log/scribe`.
- Scribe service uses the **published image** `ghcr.io/tiennm99/scribe:2.2` (no local build needed).
- README's existing Python client snippet stays as-is (documented "real API" reference).

## Dependencies

None. New folder only; touches no existing image build. Reuses `scribe/config/scribe.conf`
(read-only mount) and the published `scribe:2.2` image.
