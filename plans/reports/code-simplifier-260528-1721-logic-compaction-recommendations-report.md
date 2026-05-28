# Logic Compaction Recommendations

**Agent:** `code-simplifier` | **Date:** 2026-05-28 17:21 | **Scope:** D:\tiennm99\docker-images, REPORT-ONLY pass (no edits).

**User constraint:** cleanup or compact logic only; NO new features.

## 1. Workflow DRY (4 files, 279 lines total -> potential ~120 lines)

- **Recommendation:** Extract a reusable `workflow_call` workflow at `.github/workflows/_publish-image.yml`. Each per-image file shrinks to ~25 lines (path trigger + `uses:` + inputs). Net save: ~155 lines (~55%).
- **Concrete sketch** (`_publish-image.yml`):

```yaml
name: _publish-image
on:
  workflow_call:
    inputs:
      image:       { required: true,  type: string }   # e.g. "couchbase"
      context:     { required: true,  type: string }   # e.g. "couchbase-2.5"
      tags:        { required: true,  type: string }   # multi-line "type=raw,value=..."
      labels:      { required: true,  type: string }   # multi-line OCI labels
jobs:
  push:
    runs-on: ubuntu-latest
    permissions: { packages: write, contents: read, attestations: write, id-token: write }
    steps:
      - uses: actions/checkout@v5
      - uses: docker/login-action@f4ef78c0...
        with: { username: ${{ secrets.DOCKER_USERNAME }}, password: ${{ secrets.DOCKER_PASSWORD }} }
      - uses: docker/login-action@65b78e6e...
        with: { registry: ghcr.io, username: ${{ github.actor }}, password: ${{ secrets.GITHUB_TOKEN }} }
      - id: meta
        uses: docker/metadata-action@9ec57ed1...
        with:
          images: |
            ghcr.io/${{ github.actor }}/${{ inputs.image }}
            ${{ secrets.DOCKER_USERNAME }}/${{ inputs.image }}
          tags:   ${{ inputs.tags }}
          labels: ${{ inputs.labels }}
      - id: push
        uses: docker/build-push-action@3b5e8027...
        with: { context: ${{ inputs.context }}, push: true, tags: ${{ steps.meta.outputs.tags }}, labels: ${{ steps.meta.outputs.labels }} }
      - uses: actions/attest-build-provenance@v3
        with: { subject-name: ghcr.io/${{ github.actor }}/${{ inputs.image }}, subject-digest: ${{ steps.push.outputs.digest }}, push-to-registry: true }
```

  Note: `secrets: inherit` must be passed by the caller, OR pass secrets explicitly as `secrets:` inputs. Recommend `secrets: inherit` (simplest). **Also see code-reviewer report HIGH finding: replace `${{ github.actor }}` with `${{ github.repository_owner }}` in image-path sites before lifting into the reusable workflow.**
- **Trade-off:** A one-off image tweak (e.g. unique build-arg) requires either a new input or breaking out of the reusable workflow.
- **Effort:** S. **Risk:** LOW (action versions/inputs preserved 1:1; published artifact unchanged). **Apply?** YES.

## 2. Dockerfile LABEL compaction (4 files)

- **Recommendation:** Merge 5-6 separate `LABEL` lines per Dockerfile into a single multi-line `LABEL`. Same metadata, fewer image layers, Docker best practice. Saves ~4 lines/file + 5 image layers/image.
- **Sketch:**
  ```dockerfile
  LABEL maintainer="..." \
        org.opencontainers.image.title="Couchbase Server 2.5" \
        org.opencontainers.image.description="..." \
        org.opencontainers.image.version="2.5.2" \
        org.opencontainers.image.source="https://github.com/couchbase/docker" \
        org.opencontainers.image.licenses="Apache-2.0"
  ```
- **Effort:** S. **Risk:** LOW. **Apply?** YES.

## 3. Dockerfile RUN/ENV consolidation

- **couchbase-2.5:** lines 13-15 are 3 separate `ENV`s for CB version constants; collapse to one `ENV` block. Lines 30-32: `ADD` followed by 2-cmd `RUN` (`dos2unix` + `chmod`) — leave (already minimal).
- **oraclejdk-8:** already a single `RUN` — nothing to merge.
- **gradle-8:** 4 separate `RUN` blocks each with `set -o errexit -o nounset`. Merging would help layers but harms readability. Leave as-is.
- **scribe-2.2:** builder stage has 3 sequential `RUN`s (thrift build, fb303, scribe clone). Could merge but hurts debuggability. Leave.
- **`set -eux` vs `set -o errexit -o nounset` drift:** oraclejdk + scribe use `set -eux`; gradle uses long form. Unify to `set -eux` (3 chars vs ~30, identical semantics minus `-x` tracing which helps CI). Cosmetic.
- **Effort:** S. **Risk:** LOW. **Apply?** YES for LABEL+ENV only; NO for further RUN merging.

## 4. README compaction (74 -> ~45 lines)

- **Recommendation:** Per-image `##` sections (lines 28-69) are mostly empty headings with a Credits sublist. Image-index table already covers pull commands implicitly via the GHCR link column. Collapse all per-image sections into a single `## Credits` section grouped by image; drop the redundant `docker pull` example for Couchbase (table link suffices) — OR keep one canonical "Usage" snippet for the whole repo.
- Couchbase-2.5 still uniquely has a docker-compose example link — move to the table's Notes column as `compose example` link.
- **Effort:** S. **Risk:** LOW (docs only). **Apply?** YES.

## 5. scribe-2.2 subtree

- `example/logs/` is **empty** — bind-mount target intended by `docker-compose.yml`. Either add a `.gitkeep` or drop the directory and document the bind mount creates it. **Apply?** YES (drop empty dir; one-line `mkdir` note in scribe README).
- `example/sender/` (Dockerfile + go.mod + main.go) appears to be a real demo; verify referenced in `scribe-2.2/README.md` before removing.

## 6. couchbase-2.5/scripts

- Only `couchbase-start` exists; referenced by Dockerfile `ADD`. Nothing to drop.

## 7. template.yml

- `.github/workflows/template.yml` (65 lines) is a reference scaffold not invoked by any trigger except `release.published`. If the team is comfortable with the new reusable workflow as the canonical template, **delete it**. Saves 65 lines. **Risk:** LOW (zero references). **Apply?** YES if the reusable workflow lands; otherwise keep.

## Summary table

| # | Recommendation | Lines saved | Risk | Effort | Apply? |
|---|---------------|-------------|------|--------|--------|
| 1 | Reusable `workflow_call` for 4 publish workflows | ~155 | LOW | S | YES |
| 2 | Merge 5 LABEL lines -> 1 multi-line LABEL (x4) | ~16 + 20 layers | LOW | S | YES |
| 3 | Consolidate 3 `ENV` lines in couchbase Dockerfile | ~2 | LOW | S | YES |
| 3b | Unify `set -eux` vs `set -o errexit -o nounset` | ~0 (cosmetic) | LOW | S | DEPENDS |
| 4 | Collapse per-image README sections -> single Credits | ~25 | LOW | S | YES |
| 5 | Drop empty `scribe-2.2/example/logs/` dir | -1 dir | LOW | S | YES |
| 7 | Delete `.github/workflows/template.yml` after #1 lands | ~65 | LOW | S | YES (post #1) |

**Total saved if all applied: ~260 lines across the repo (~35-40% of YAML+Dockerfile+README surface area), plus ~20 fewer image layers.**

## Out of scope (NEW features only, NOT in cleanup scope)

- Adding renovate/dependabot for action SHA pinning
- Adding image-scanning (trivy/grype) steps
- Adding build matrices for multi-arch (arm64)
- Adding tests/smoke-tests for the published images
- Adding signing (cosign) beyond the existing attestation
- Renaming workflow files (already done this session under Convention B)

## Unresolved questions

1. Should the reusable workflow live at `.github/workflows/_publish-image.yml` (private to repo, `workflow_call`) or as a **composite action** at `.github/actions/publish-image/`? Composite action is reusable across repos and slightly more concise; reusable workflow keeps the per-image file as a true 1:1 mapping of `on:` triggers. Lead's call.
2. Confirm `secrets: inherit` is acceptable, or prefer explicit `secrets:` block for least-privilege clarity.
3. Confirm scribe `example/sender/` is actively referenced (don't yank without checking the scribe README).
