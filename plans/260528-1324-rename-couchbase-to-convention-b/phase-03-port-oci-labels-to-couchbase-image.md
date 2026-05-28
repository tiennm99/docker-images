---
phase: 3
title: "Port OCI labels to couchbase image"
status: code-applied
priority: P2
effort: "10m"
dependencies: [2]
---

# Phase 3: Port OCI labels to couchbase image

## Overview

Add OCI image-spec labels to couchbase using the proven pattern from `scribe-2.2/Dockerfile` + `.github/workflows/publish-scribe-2.2.yml`. Both Dockerfile (manual-build path) and workflow `metadata-action.labels:` block (CI-build path) set the same values, since CLI `--label` from metadata-action overrides Dockerfile `LABEL` instructions.

## Requirements

- Functional: published image at `ghcr.io/tiennm99/couchbase:2.5.2` carries 5 OCI labels with correct values; metadata-action's auto-defaults (derived from THIS repo) do not shadow upstream-Couchbase metadata.
- Non-functional: plain ASCII; pattern matches scribe verbatim for consistency.

## Architecture

Mirror scribe-2.2's structure:
- `Dockerfile` runtime stage carries 5 `LABEL org.opencontainers.image.*` lines after `LABEL maintainer`.
- Workflow `metadata-action` step has a `labels:` block in YAML setting the same 5 keys (overrides metadata-action defaults).

## Research input

- `plans/reports/researcher-260528-1324-couchbase-license-and-source-verification-report.md`
  - License: `Apache-2.0` (Couchbase Server 2.5.2, Dec 2012; BSL adoption was later, post-2021)
  - `image.source`: `https://github.com/couchbase/docker` (general repo, not pinned commit)

## Label values

| Key | Value |
|-----|-------|
| `org.opencontainers.image.title` | `Couchbase Server 2.5` |
| `org.opencontainers.image.description` | `Legacy Couchbase Server (no official upstream image)` |
| `org.opencontainers.image.version` | `2.5.2` |
| `org.opencontainers.image.source` | `https://github.com/couchbase/docker` |
| `org.opencontainers.image.licenses` | `Apache-2.0` |

`org.opencontainers.image.created`, `url`, `revision` left to metadata-action defaults (acceptable).

## Related Code Files

- Modify: `D:\tiennm99\docker-images\couchbase-2.5\Dockerfile` (insert 5 LABEL lines after L3 `LABEL maintainer`)
- Modify: `D:\tiennm99\docker-images\.github\workflows\publish-couchbase-2.5.yml` (add `labels:` block after `tags:` in metadata-action step)

## Implementation Steps

1. Edit `couchbase-2.5/Dockerfile`:

   **Before (L3):**
   ```dockerfile
   LABEL maintainer="Tien Nguyen Minh <tiennm99@outlook.com>"
   ```

   **After:**
   ```dockerfile
   LABEL maintainer="Tien Nguyen Minh <tiennm99@outlook.com>"
   LABEL org.opencontainers.image.title="Couchbase Server 2.5"
   LABEL org.opencontainers.image.description="Legacy Couchbase Server (no official upstream image)"
   LABEL org.opencontainers.image.version="2.5.2"
   LABEL org.opencontainers.image.source="https://github.com/couchbase/docker"
   LABEL org.opencontainers.image.licenses="Apache-2.0"
   ```

2. Edit `.github/workflows/publish-couchbase-2.5.yml` metadata-action step. Insert `labels:` block immediately after the `tags:` block (after Phase 02's edits):

   ```yaml
             labels: |
               org.opencontainers.image.title=Couchbase Server 2.5
               org.opencontainers.image.description=Legacy Couchbase Server (no official upstream image)
               org.opencontainers.image.version=2.5.2
               org.opencontainers.image.source=https://github.com/couchbase/docker
               org.opencontainers.image.licenses=Apache-2.0
   ```

   Do NOT quote values in the `labels:` block - metadata-action splits on `=` only, quotes become part of the value.

3. Commit: `feat(couchbase): expand OCI image-spec labels`

## Success Criteria

- [ ] Couchbase Dockerfile has 5 OCI LABEL lines after maintainer
- [ ] Workflow `metadata-action` step has `labels:` block with the same 5 keys, no quoted values
- [ ] One commit, message starts with `feat(couchbase):`
- [ ] After workflow runs: `docker inspect ghcr.io/tiennm99/couchbase:2.5.2 --format '{{json .Config.Labels}}'` returns all 5 OCI keys with upstream-Couchbase values (not auto-derived from this repo)

## Risk Assessment

- **Risk:** license verification relied on Wikipedia + Couchbase blog; if 2.5.2 binaries are actually licensed differently from Couchbase open-source code, the label is wrong. **Mitigation:** label is metadata, not legal claim; updateable. If Couchbase publishes a contradicting LICENSE.txt for 2.5.2, swap the label value in a one-line follow-up commit.
- **Risk:** workflow `labels:` block YAML indentation drift (it sits inside `metadata-action.with:`). **Mitigation:** match scribe's exact indentation in `publish-scribe-2.2.yml` L48-53 (proven working).

## Validation (after push)

- Workflow run on `main` after this phase's commit completes successfully.
- `docker pull ghcr.io/tiennm99/couchbase:2.5.2` then `docker inspect couchbase:2.5.2 --format '{{json .Config.Labels}}' | jq`. Expect:
  - All 5 OCI keys above with specified values
  - `org.opencontainers.image.created` set to recent ISO-8601 timestamp (metadata-action default)
  - `org.opencontainers.image.url` and `org.opencontainers.image.revision` set to this repo's metadata (metadata-action default, acceptable)
  - `maintainer` set to `Tien Nguyen Minh <tiennm99@outlook.com>` (Dockerfile only)
- Manual `docker build couchbase-2.5/ -t couchbase:test`; `docker inspect couchbase:test --format '{{json .Config.Labels}}'` shows 5 OCI keys + maintainer (no `created`/`url`/`revision` since metadata-action absent on manual build).
