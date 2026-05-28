---
phase: 4
title: "Workflow_call DRY refactor"
status: pending
priority: P3
effort: "60m"
dependencies: [1, 2]
---

# Phase 4: Workflow_call DRY refactor

## Overview

Extract the 98%-identical publish logic from the 4 `publish-*.yml` workflows into one reusable workflow at `.github/workflows/_publish-image.yml`. Each per-image file shrinks from ~70 lines to ~25 lines (path trigger + `uses:` + inputs). Net save: ~155 lines (~55% of workflow YAML).

Blocked on Phase 1 (would otherwise carry the `github.actor` bug forward into the reusable workflow) and Phase 2 (whose decisions affect what's left to consolidate, e.g. should `name:` per-image stay or move into the reusable).

This is a STRUCTURAL refactor — the published artifact MUST remain identical (same image namespace, tags, labels, attestation). Verify with `docker inspect` before and after.

## Open question (must answer before cook)

**Workflow reuse style:** repo-private `workflow_call` at `.github/workflows/_publish-image.yml` OR composite action at `.github/actions/publish-image/action.yml`?

| Aspect | `workflow_call` | composite action |
|--------|-----------------|------------------|
| Reusable across repos | No (without `uses: <org>/<repo>/.github/workflows/...@<ref>`) | Yes (cleaner cross-repo via `uses: <org>/<repo>/.github/actions/...@<ref>`) |
| Per-image file is still a workflow | Yes (keeps `on:` triggers in per-image file) | Yes |
| Secrets handling | `secrets: inherit` or explicit | Cannot access secrets directly; caller must pass as inputs |
| Concurrency / matrix support | Native | Native via caller |
| Existing pattern parity | Identical structure to per-image files today | Slightly different style |

Recommend `workflow_call` because (a) the use case is single-repo, (b) `secrets: inherit` keeps the caller side minimal, (c) per-image YAML files remain proper workflows that the GitHub UI shows under "Actions". Cook may override.

## Requirements

- Functional: every per-image workflow still triggers on its own `paths:` glob; publishes identical image namespace, tags, labels, attestation; runs in the same order/duration as today.
- Non-functional: net reduction in YAML LOC across `.github/workflows/`; identical SHA-pinned action versions in the reusable workflow; per-image files only carry image-specific inputs.

## Architecture

Two file types after refactor:

### Reusable workflow: `.github/workflows/_publish-image.yml`

```yaml
name: _publish-image
on:
  workflow_call:
    inputs:
      image:
        required: true
        type: string
        description: "GHCR + Docker Hub image name (e.g. couchbase, scribe)"
      context:
        required: true
        type: string
        description: "Docker build context directory (e.g. couchbase-2.5)"
      tags:
        required: true
        type: string
        description: "metadata-action tags block (newline-separated type=raw,value=... entries)"
      labels:
        required: true
        type: string
        description: "metadata-action labels block (newline-separated key=value entries)"
jobs:
  push:
    runs-on: ubuntu-latest
    permissions:
      packages: write
      contents: read
      attestations: write
      id-token: write
    steps:
      - uses: actions/checkout@v5
      - uses: docker/login-action@f4ef78c080cd8ba55a85445d5b36e214a81df20a
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}
      - uses: docker/login-action@65b78e6e13532edd9afa3aa52ac7964289d1a9c1
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - id: meta
        uses: docker/metadata-action@9ec57ed1fcdbf14dcef7dfbe97b2010124a938b7
        with:
          images: |
            ghcr.io/${{ github.repository_owner }}/${{ inputs.image }}
            ${{ secrets.DOCKER_USERNAME }}/${{ inputs.image }}
          tags: ${{ inputs.tags }}
          labels: ${{ inputs.labels }}
      - id: push
        uses: docker/build-push-action@3b5e8027fcad23fda98b2e3ac259d8d67585f671
        with:
          context: ${{ inputs.context }}
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
      - uses: actions/attest-build-provenance@v3
        with:
          subject-name: ghcr.io/${{ github.repository_owner }}/${{ inputs.image }}
          subject-digest: ${{ steps.push.outputs.digest }}
          push-to-registry: true
```

Note: `${{ github.actor }}` retained for ghcr.io LOGIN (correct — that's the triggering user's token). `${{ github.repository_owner }}` used for image NAMESPACE (correct per Phase 01 fix).

### Per-image caller: `.github/workflows/publish-couchbase-2.5.yml` (example)

```yaml
name: Publish couchbase Docker image
on:
  push:
    branches: ['main']
    paths:
      - 'couchbase-2.5/**'
      - '.github/workflows/publish-couchbase-2.5.yml'
      - '.github/workflows/_publish-image.yml'
  workflow_dispatch:
jobs:
  publish:
    uses: ./.github/workflows/_publish-image.yml
    secrets: inherit
    with:
      image: couchbase
      context: couchbase-2.5
      tags: |
        type=raw,value=latest
        type=raw,value=2.5
        type=raw,value=2.5.2
      labels: |
        org.opencontainers.image.title=Couchbase Server 2.5
        org.opencontainers.image.description=Legacy Couchbase Server (no official upstream image)
        org.opencontainers.image.version=2.5.2
        org.opencontainers.image.source=https://github.com/couchbase/docker
        org.opencontainers.image.licenses=Apache-2.0
```

Repeat the 4-input pattern for scribe, oraclejdk, gradle. Note added `paths:` entry for `_publish-image.yml` itself so changes to the reusable workflow re-trigger all 4 builds.

## Related Code Files

- Create: `D:\tiennm99\docker-images\.github\workflows\_publish-image.yml`
- Rewrite: all 4 `D:\tiennm99\docker-images\.github\workflows\publish-<service>-<major>.yml`

## Implementation Steps

1. Create `_publish-image.yml` (reusable) with the structure above. Use the same SHA-pinned actions as the current workflows (no version bumps in this phase).
2. Rewrite `publish-couchbase-2.5.yml` to be the thin caller. Keep `paths:` triggers identical, add `_publish-image.yml` to the trigger list.
3. Repeat for scribe, oraclejdk, gradle.
4. Verify locally: count `wc -l .github/workflows/*.yml` before vs after.
5. Commit: `refactor(workflows): extract reusable _publish-image workflow_call`.
6. Push to a feature branch first; confirm one workflow run succeeds before pushing to main. (Optional — cook decides if main-push is acceptable given the verification plan.)

## Success Criteria

- [ ] `.github/workflows/_publish-image.yml` exists with `workflow_call` trigger and 4 inputs
- [ ] Each per-image workflow file is < 35 lines (input-data only)
- [ ] Each per-image workflow's `paths:` includes both its subdir AND `.github/workflows/_publish-image.yml`
- [ ] Total LOC across `.github/workflows/` drops by ~150 lines vs Phase 3 end state
- [ ] After push: `docker inspect ghcr.io/tiennm99/<svc>:<tag> --format '{{json .Config.Labels}}'` returns the same set of OCI keys + values as before the refactor
- [ ] Each workflow's `subject-name` attestation resolves to `ghcr.io/tiennm99/<svc>` (not `<actor>/svc`)
- [ ] One commit starting with `refactor(workflows):`

## Risk Assessment

- **Risk:** `workflow_call` inputs of type `string` with multi-line content (tags, labels) may not round-trip cleanly through `metadata-action`. **Mitigation:** GitHub's documented `workflow_call` `inputs.<name>` `type: string` preserves newlines. Validate with one workflow before refactoring all 4.
- **Risk:** `secrets: inherit` exposes ALL repo secrets to the called workflow. **Mitigation:** acceptable since `_publish-image.yml` is repo-private, not a third-party reusable. If you'd prefer least-privilege, pass `DOCKER_USERNAME`/`DOCKER_PASSWORD` as explicit `secrets:` inputs.
- **Risk:** action SHAs in the reusable workflow drift from a per-image override someone added in the past. **Mitigation:** comparing the 4 current workflows shows all SHAs are identical. Phase 4 makes them shared by construction.
- **Risk:** GitHub Actions UI may render reusable-workflow runs slightly differently — could surprise muscle memory. **Mitigation:** that's a one-time adjustment.
- **Risk:** if the reusable workflow has a bug, all 4 image publishes break simultaneously instead of just one. **Mitigation:** that's the standard DRY trade-off; offset by easier auditing (one place to fix).

## Validation

- Before refactor: capture `docker inspect ghcr.io/tiennm99/<svc>:<latest-tag>` for all 4 services, save outputs locally for diff.
- After refactor: trigger one `workflow_dispatch` per image; wait for completion; re-inspect; diff before/after. All differences should be limited to the new `image.revision` (commit SHA) and `image.created` timestamp (build time).
- Manually verify the 4 `Actions` UI runs show distinct names (Phase 03 disambiguation should already cover this).
