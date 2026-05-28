---
status: superseded
supersededBy: 260528-1324-rename-couchbase-to-convention-b
---

# Image Naming Convention - Decision Parked

> **Superseded by [260528-1324-rename-couchbase-to-convention-b](../260528-1324-rename-couchbase-to-convention-b/plan.md) on 2026-05-28.** User chose Convention B; see superseding plan for rationale.

Status: **SUPERSEDED** - Convention B locked in via 260528-1324
Created: 2026-05-28 11:10
Owner: tiennm99

## Question

How should this repo name its Docker images?

- Convention A (current for 3/4): `<service>-<major-line>:<full-version>` e.g. `couchbase-2.5:2.5.2`
- Convention B (current for scribe): `<service>:<full-version>` e.g. `scribe:2.2` (Docker-Hub-official style)

## Context

- 4 images total at brainstorm time: `couchbase-2.5`, `gradle-8`, `oraclejdk-8`, `scribe`
- `gradle-8` and `oraclejdk-8` dropped 2026-05-28 (this commit) - decision now only spans `couchbase-2.5` + scribe
- scribe diverges from the rest
- Subdir already named `scribe-2.2/` regardless of which convention wins

## Recommendation (from brainstorm)

**Convention A**, rename scribe to `scribe-2.2`. Decisive arguments:

1. `latest` semantics break under B for archived software (newest-wins misleads users searching for legacy versions)
2. Attestation lineage cleaner per-line under A
3. Docker-official `python:3.x` analogy doesn't transfer - those are sliding-window builds; this repo is a museum of frozen artifacts
4. Migration A = ~5 edits; B = ~15 edits + orphan-package cleanup + new `latest` policy

See reports:
- `plans/reports/brainstorm-260528-1007-image-naming-convention-report.md` (initial)
- `plans/reports/brainstorm-260528-1015-image-naming-convention-reframe-report.md` (re-frame, stronger reasoning, same conclusion)

## Decisions pending

1. **Approve Convention A?** If yes: rename scribe to `scribe-2.2` (~5 edits).
2. **Codify rule** in README + CLAUDE.md as part of same change, or separate task?
3. **Orphan packages** post-rename: delete `ghcr.io/tiennm99/scribe` + Docker Hub counterpart, or tombstone in README?

## When to revisit

- Before adding a new image (so the new image follows the chosen convention)
- Or sooner if user wants to lock the rule down
