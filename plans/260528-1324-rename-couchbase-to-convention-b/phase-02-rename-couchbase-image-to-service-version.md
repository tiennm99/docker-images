---
phase: 2
title: "Rename couchbase image to <service>:<version>"
status: code-applied
priority: P1
effort: "10m"
dependencies: [1]
---

# Phase 2: Rename couchbase image to `<service>:<version>`

## Overview

Rename the published image from `couchbase-2.5` to `couchbase` across GHCR + Docker Hub. Keep subdirectory and workflow filename unchanged (path-trigger anchor). Add `2.5` major-line alias tag while we're here.

## Requirements

- Functional: after push + workflow run, `docker pull ghcr.io/tiennm99/couchbase:2.5.2` resolves; old `ghcr.io/tiennm99/couchbase-2.5:2.5.2` becomes orphan.
- Non-functional: subdirectory `couchbase-2.5/` not renamed; workflow filename not renamed.

## Architecture

Three files touch the literal string `couchbase-2.5` in an image-reference context:
- `.github/workflows/publish-couchbase-2.5.yml` - `images:` block (L42-43) + `subject-name:` (L61)
- `README.md` - image-index table row link target (`pkgs/container/couchbase-2.5` -> `pkgs/container/couchbase`) and the `docker pull ghcr.io/tiennm99/couchbase-2.5:2.5.2` example block in the Couchbase section.
- `couchbase-2.5/docker-compose.example.yml` - L3 `image: ghcr.io/tiennm99/couchbase-2.5:2.5.2`

Subdir name and workflow filename are NOT image references; they stay.

## Related Code Files

- Modify: `D:\tiennm99\docker-images\.github\workflows\publish-couchbase-2.5.yml`
- Modify: `D:\tiennm99\docker-images\README.md`
- Modify: `D:\tiennm99\docker-images\couchbase-2.5\docker-compose.example.yml`

## Implementation Steps

1. Edit `.github/workflows/publish-couchbase-2.5.yml`:

   **Before (L42-47):**
   ```yaml
           images: |
             ghcr.io/${{ github.actor }}/couchbase-2.5
             ${{ secrets.DOCKER_USERNAME }}/couchbase-2.5
           tags: |
             type=raw,value=latest
             type=raw,value=2.5.2
   ```

   **After:**
   ```yaml
           images: |
             ghcr.io/${{ github.actor }}/couchbase
             ${{ secrets.DOCKER_USERNAME }}/couchbase
           tags: |
             type=raw,value=latest
             type=raw,value=2.5
             type=raw,value=2.5.2
   ```

   **Before (L61):**
   ```yaml
           subject-name: ghcr.io/${{ github.actor }}/couchbase-2.5
   ```

   **After:**
   ```yaml
           subject-name: ghcr.io/${{ github.actor }}/couchbase
   ```

2. Edit `README.md`:
   - Image-index table row: link `[couchbase-2.5](https://github.com/tiennm99/docker-images/pkgs/container/couchbase-2.5)` -> `[couchbase](https://github.com/tiennm99/docker-images/pkgs/container/couchbase)`. Display text in the link also flips to `couchbase`.
   - Couchbase section `docker pull ghcr.io/tiennm99/couchbase-2.5:2.5.2` example -> `docker pull ghcr.io/tiennm99/couchbase:2.5.2`
   - "More options" GHCR URL in the same section also flips.

3. Edit `couchbase-2.5/docker-compose.example.yml` L3:
   - `image: ghcr.io/tiennm99/couchbase-2.5:2.5.2` -> `image: ghcr.io/tiennm99/couchbase:2.5.2`

4. Commit: `refactor(couchbase): rename image package to couchbase (Convention B)`

## Success Criteria

- [ ] Workflow `images:` block lists `couchbase`, not `couchbase-2.5`
- [ ] Workflow `tags:` includes new `2.5` major-line alias
- [ ] Workflow `subject-name:` ends in `/couchbase`
- [ ] README image-index row link target is `pkgs/container/couchbase`
- [ ] `couchbase-2.5/docker-compose.example.yml` pulls from `ghcr.io/tiennm99/couchbase:2.5.2`
- [ ] `grep -r "tiennm99/couchbase-2.5" D:\tiennm99\docker-images\` (excluding plans/reports) returns 0 results
- [ ] One commit, message starts with `refactor(couchbase):`

## Risk Assessment

- **Risk:** stale `docker pull ghcr.io/tiennm99/couchbase-2.5:2.5.2` references in user environments break after orphan deletion. **Mitigation:** repo is one-maintainer; no known external consumers. Old package remains pullable until user deletes via UI (manual step, deferred).
- **Risk:** `2.5` tag alias confuses users into expecting it to track patches within 2.5.x. **Mitigation:** only one 2.5.x release exists (2.5.2); `2.5` is effectively a synonym. Codified rule in Phase 01 calls it an "optional alias".

## Validation (after push, in Phase 03 or later)

- Workflow run completes successfully on GitHub Actions
- `docker pull ghcr.io/tiennm99/couchbase:2.5.2` resolves
- `docker pull ghcr.io/tiennm99/couchbase:latest` resolves
- `docker pull ghcr.io/tiennm99/couchbase:2.5` resolves
- Root README image-index link resolves to a real GHCR page
