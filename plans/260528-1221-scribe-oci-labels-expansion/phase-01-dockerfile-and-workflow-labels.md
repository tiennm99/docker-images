# Phase 01 — Dockerfile + Workflow Labels

**Status:** pending
**Effort:** ~20m

## Context Links

- Research report: `plans/reports/researcher-260528-1221-scribe-oci-labels-license-verification-report.md`
- metadata-action docs: https://github.com/docker/metadata-action (auto-injects 8 OCI labels by default; `labels:` input overrides them)

## Key findings driving the design

1. `docker/metadata-action@9ec57ed` (workflow line 40) auto-injects: `title`, `description`, `url`, `source`, `version`, `created`, `revision`, `licenses` — values derived from THIS repo's GitHub metadata.
2. CLI `--label` (which `docker/build-push-action` uses via `labels: ${{ steps.meta.outputs.labels }}`) OVERRIDES Dockerfile `LABEL` instructions at build time. So workflow builds today shadow the Dockerfile's labels for `version`, `description`, `source`.
3. Auto-injected `source` and `title` are WRONG for this image (point to our build repo / `docker-images` name, not Scribe upstream). Must be overridden in the workflow's `metadata-action.labels:` block.
4. Auto-injected `created` is correct (build timestamp). Keep it.
5. Manual `docker build` (no CI) skips metadata-action entirely, so Dockerfile labels still matter for that path.

## Final label set on published image

| Key | Value | Source |
|-----|-------|--------|
| `org.opencontainers.image.title` | `Facebook Legacy Scribe` | Dockerfile + workflow override |
| `org.opencontainers.image.description` | `Facebook Legacy Scribe - real-time log aggregation daemon` | Dockerfile + workflow override |
| `org.opencontainers.image.version` | `2.2` | Dockerfile + workflow override |
| `org.opencontainers.image.source` | `https://github.com/facebookarchive/scribe` | Dockerfile + workflow override |
| `org.opencontainers.image.licenses` | `Apache-2.0` | Dockerfile + workflow override |
| `org.opencontainers.image.created` | (build timestamp) | metadata-action auto |
| `org.opencontainers.image.url` | (our repo URL) | metadata-action auto (acceptable — IS where image is built) |
| `org.opencontainers.image.revision` | (git SHA of our repo) | metadata-action auto (acceptable — identifies build recipe commit) |
| `maintainer` | `Tien Nguyen Minh <tiennm99@outlook.com>` | Dockerfile only (not OCI; metadata-action does not auto-inject) |

## Files to modify

### 1. `scribe-2.2/Dockerfile` — replace lines 70-73

**Before:**
```dockerfile
LABEL maintainer="Tien Nguyen Minh <tiennm99@outlook.com>"
LABEL org.opencontainers.image.version="2.2"
LABEL org.opencontainers.image.description="Facebook Legacy Scribe - real-time log aggregation daemon"
LABEL org.opencontainers.image.source="https://github.com/facebookarchive/scribe"
```

**After:**
```dockerfile
LABEL maintainer="Tien Nguyen Minh <tiennm99@outlook.com>"
LABEL org.opencontainers.image.title="Facebook Legacy Scribe"
LABEL org.opencontainers.image.description="Facebook Legacy Scribe - real-time log aggregation daemon"
LABEL org.opencontainers.image.version="2.2"
LABEL org.opencontainers.image.source="https://github.com/facebookarchive/scribe"
LABEL org.opencontainers.image.licenses="Apache-2.0"
```

Order: `maintainer` first, then OCI labels grouped (title, description, version, source, licenses). Plain ASCII only.

### 2. `.github/workflows/publish-scribe-2.2.yml` — extend metadata-action step (lines 38-47)

**Before (lines 38-47):**
```yaml
      - name: Extract metadata (tags, labels) for Docker
        id: meta
        uses: docker/metadata-action@9ec57ed1fcdbf14dcef7dfbe97b2010124a938b7
        with:
          images: |
            ghcr.io/${{ github.actor }}/scribe
            ${{ secrets.DOCKER_USERNAME }}/scribe
          tags: |
            type=raw,value=latest
            type=raw,value=2.2
```

**After:**
```yaml
      - name: Extract metadata (tags, labels) for Docker
        id: meta
        uses: docker/metadata-action@9ec57ed1fcdbf14dcef7dfbe97b2010124a938b7
        with:
          images: |
            ghcr.io/${{ github.actor }}/scribe
            ${{ secrets.DOCKER_USERNAME }}/scribe
          tags: |
            type=raw,value=latest
            type=raw,value=2.2
          labels: |
            org.opencontainers.image.title=Facebook Legacy Scribe
            org.opencontainers.image.description=Facebook Legacy Scribe - real-time log aggregation daemon
            org.opencontainers.image.version=2.2
            org.opencontainers.image.source=https://github.com/facebookarchive/scribe
            org.opencontainers.image.licenses=Apache-2.0
```

Note: do NOT quote values in the metadata-action `labels:` block — quotes become part of the label value. metadata-action splits on `=` only.

## Implementation steps

1. Edit `scribe-2.2/Dockerfile` lines 70-73 per spec above.
2. Edit `.github/workflows/publish-scribe-2.2.yml` adding the `labels:` block after `tags:`.
3. Lint/validate yaml manually (no schema check needed — small change).
4. Stage only these 2 files. Leave the pending ASCII cleanup unstaged.
5. Commit: `feat(scribe): expand OCI image-spec labels` (conventional commit, no plan refs in message per global rules).
6. Push. Workflow triggers because path filter `scribe-2.2/**` and `.github/workflows/publish-scribe-2.2.yml` both match.
7. Validate (see below).

## Commit plan

**One commit for this plan.** Bundles Dockerfile + workflow because they describe the same change (label set expansion) and must land together to avoid a window where workflow builds and Dockerfile disagree.

The existing pending ASCII cleanup is a separate logical commit — handle BEFORE or AFTER this plan, never in the same commit. Plan does not stage or modify those files.

Commit message:
```
feat(scribe): expand OCI image-spec labels

Add image.title and image.licenses; override metadata-action defaults
in publish workflow so published image carries Scribe upstream source
and correct title rather than this repo's auto-derived values.
```

## Validation

After workflow run completes and image is published:

```bash
docker buildx imagetools inspect ghcr.io/tiennm99/scribe:2.2 --format '{{json .Manifest}}'
# OR inspect labels specifically:
docker pull ghcr.io/tiennm99/scribe:2.2
docker inspect ghcr.io/tiennm99/scribe:2.2 --format '{{json .Config.Labels}}' | jq
```

Expected: all 5 OCI keys above present with the specified values; `created` set to a recent ISO-8601 timestamp; `url`/`revision` present (metadata-action defaults); `maintainer` present.

For manual-build validation (no CI):
```bash
docker build -t scribe:2.2-local scribe-2.2
docker inspect scribe:2.2-local --format '{{json .Config.Labels}}' | jq
```
Expected: 5 OCI keys + `maintainer` from Dockerfile; no `created`/`url`/`revision` (those only exist via metadata-action).

## Todo

- [ ] Edit `scribe-2.2/Dockerfile` lines 70-73 — final 6-label block
- [ ] Edit `.github/workflows/publish-scribe-2.2.yml` — add `labels:` block to metadata-action step
- [ ] `git add scribe-2.2/Dockerfile .github/workflows/publish-scribe-2.2.yml` (explicit paths; do NOT stage other modified files)
- [ ] Verify `git status` shows only those two files staged; ASCII-cleanup files still unstaged
- [ ] Commit with conventional message above
- [ ] Push to `main`
- [ ] Wait for `publish-scribe-2.2.yml` workflow to complete (GitHub Actions tab)
- [ ] Run `docker inspect` validation against `ghcr.io/tiennm99/scribe:2.2`
- [ ] Confirm all 5 OCI labels carry expected values; tick this phase complete

## Success criteria

- Published image at `ghcr.io/tiennm99/scribe:2.2` AND `${DOCKER_USERNAME}/scribe:2.2` carries:
  - `org.opencontainers.image.title=Facebook Legacy Scribe`
  - `org.opencontainers.image.source=https://github.com/facebookarchive/scribe` (NOT this repo's URL)
  - `org.opencontainers.image.licenses=Apache-2.0`
  - `org.opencontainers.image.version=2.2`
  - `org.opencontainers.image.description` with `Facebook Legacy Scribe - real-time log aggregation daemon`
- Workflow continues to publish `2.2` and `latest` tags (no regression).
- Manual `docker build scribe-2.2/` still works and produces labeled image.

## Rollback

Single commit revert: `git revert <sha>` and force-trigger workflow. No data, no schema, no dependents — rollback is purely cosmetic metadata.
