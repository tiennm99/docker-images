# Code Review: scribe-2.2 OCI labels + cleanup verification

**Date:** 2026-05-28
**Scope:** uncommitted `scribe-2.2/` changes vs HEAD (`ff5e5f6`)
**Focus:** OCI label re-addition delta (cleanup already cleared in `code-review-260528-0925-scribe-cleanup-verification-report.md`)

## Scope
- Files changed: 7 (Dockerfile, README, config/scribe.conf, example/README, example/docker-compose.yml, example/sender/main.go, scripts/docker-entrypoint.sh)
- LOC: +38 / -57
- Delta under review: 3 `LABEL org.opencontainers.image.*` lines at Dockerfile:71-73

## Overall Assessment
Clean. OCI label addition is correct, well-formed, properly placed. Cleanup portion unchanged from prior approval. Ready to commit.

## OCI Label Verification

### Key names — all spec-compliant
Verified against https://github.com/opencontainers/image-spec/blob/main/annotations.md:
- `org.opencontainers.image.version` — OK, "version of the packaged software"
- `org.opencontainers.image.description` — OK, "human-readable description"
- `org.opencontainers.image.source` — OK, "URL to get source code"

### Values — well-formed
- `"2.2"` — plain ASCII string, OK
- `"Facebook Legacy Scribe - real-time log aggregation daemon"` — ASCII hyphen (was `–` U+2013 before cleanup), no smart quotes, no escaping issues
- `"https://github.com/facebookarchive/scribe"` — HTTP 200 verified, sane

### Placement — correct
Dockerfile:68-73 — runtime stage (`FROM centos:7`), after `LABEL maintainer`. Builder stage correctly left unlabeled (only runtime image gets metadata).

## Optional OCI labels worth considering (NOT auto-added)

Recommendations only — user decides:

| Label | Suggested value | Rationale |
|-------|----------------|-----------|
| `image.title` | `"Facebook Legacy Scribe"` | Short human-readable name, complements description. Most-commonly-used OCI label after version. Worth adding. |
| `image.licenses` | `"Apache-2.0"` | Scribe's license (SPDX expression). Useful for downstream consumers / compliance scans. Worth adding. |
| `image.url` | could point to this repo or README | Marginal value since `image.source` already covers provenance. Skip unless you want to distinguish "where to learn" from "where the source lives". |
| `image.revision` | scribe upstream git SHA at build time | Scribe is `--depth 1` cloned (Dockerfile:60) → revision drifts silently between rebuilds. Pinning would need a build ARG. Defer; not blocker. |
| `image.created` | RFC 3339 timestamp | BuildKit auto-injects this; manual labeling not needed. |
| `image.vendor` | `"tiennm99"` or org | Optional, marginal. Skip. |

**Strongest recommendation:** add `image.title` + `image.licenses`. The other four are skippable for an archive image.

## Diff Drift Check
Confirmed: full scribe-2.2 diff = (previously approved cleanup) + (3 OCI labels). No unintended edits.

Per-file:
- **Dockerfile** — cleanup (box-drawing comments, alpha-sorted builder deps, en-dash → hyphen) + OCI label swap. Clean.
- **README.md** — 2 en-dash → hyphen (cleanup). Clean.
- **config/scribe.conf** — 6 box-drawing/en-dash → ASCII (cleanup). Clean.
- **example/README.md** — 8 en-dash → hyphen (cleanup). Clean.
- **example/docker-compose.yml** — 1 en-dash → colon (cleanup). Clean.
- **example/sender/main.go** — 4 en-dash → hyphen/semicolon (cleanup, includes user-visible string `Hello from Go sender - message N/5`). Clean.
- **scripts/docker-entrypoint.sh** — 1 en-dash → hyphen (cleanup). Clean.

No new files. No deletions. No file moves.

## Critical Issues
None.

## High Priority
None.

## Medium Priority
None blocking. See "Optional OCI labels" above for incremental improvement.

## Low Priority
- CRLF warnings from git (LF will be replaced by CRLF on Windows checkout). Pre-existing repo behavior, not a delta concern.

## Positive Observations
- Migrated from non-standard `LABEL description` / `LABEL version` to OCI-standard keys — proper fix to the prior concern.
- All ASCII; no Unicode dashes carried over from cleanup pass into label values.
- Single-line `LABEL` per key (readable; multi-line continuation would also be fine, current style is clearer).
- `image.source` URL matches the actual upstream cloned at Dockerfile:60 (`https://github.com/facebookarchive/scribe.git`) — consistent.

## Recommended Actions
1. **Ship as-is** — no blockers.
2. *(Optional, recommend)* — add two more OCI labels before committing:
   ```dockerfile
   LABEL org.opencontainers.image.title="Facebook Legacy Scribe"
   LABEL org.opencontainers.image.licenses="Apache-2.0"
   ```
   Confirm Scribe upstream license is Apache-2.0 first (facebookarchive/scribe LICENSE file).

## Metrics
- OCI spec conformance: 3/3 keys valid
- Label value well-formedness: pass
- Source URL reachability: HTTP 200
- Drift beyond approved cleanup: zero

## Unresolved Questions
- Do you want to add `image.title` + `image.licenses` now, or ship the minimal 3-label set? (Recommend adding both; trivial extra metadata.)
- Confirm Scribe upstream license = Apache-2.0 before setting `image.licenses` (worth a 5s check at https://github.com/facebookarchive/scribe/blob/master/LICENSE).
