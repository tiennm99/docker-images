---
phase: 3
title: "Low-priority polish and consistency"
status: pending
priority: P3
effort: "20m"
dependencies: []
---

# Phase 3: Low-priority polish and consistency

## Overview

Polish-only changes: compact 5-line OCI LABEL blocks into one multi-line `LABEL`, disambiguate the 4 workflow `name:` fields, normalize `images:` ordering, convert build-only `ENV`s to `ARG`, expand `.gitignore`. No correctness risk, no behavior change.

Independent of other phases — can ship before, after, or between them.

## Requirements

- Functional: image contracts unchanged (same labels, same tags, same image namespace, same behavior).
- Non-functional: fewer image layers, less visual noise, more discoverable Actions UI, less env-pollution at runtime.

## Architecture

Five independent edits:
1. Multi-line `LABEL` compaction in all 4 Dockerfiles (5 single-line LABELs → 1 multi-line LABEL)
2. Workflow `name:` disambiguation (4 workflows currently all read `Publish Docker image`)
3. `images:` list ordering uniformity (some workflows list ghcr first, others list Docker Hub first)
4. `couchbase-2.5/Dockerfile` `ENV` → `ARG` for build-only constants
5. `.gitignore` additions for common IDE / OS files

## Related Code Files

- Modify: all 4 `<service>-<major>/Dockerfile`
- Modify: all 4 `.github/workflows/publish-<service>-<major>.yml`
- Modify: `D:\tiennm99\docker-images\.gitignore`

## Implementation Steps

### Step 1 — LABEL compaction (4 Dockerfiles)

Replace 5 separate `LABEL org.opencontainers.image.*` lines with 1 multi-line LABEL. Keep the existing `LABEL maintainer` either separate (1 layer) or fold into the multi-line block (0 extra layer). Recommend folding for max compaction.

Example for couchbase-2.5/Dockerfile:

**Before (L3-8):**
```dockerfile
LABEL maintainer="Tien Nguyen Minh <tiennm99@outlook.com>"
LABEL org.opencontainers.image.title="Couchbase Server 2.5"
LABEL org.opencontainers.image.description="Legacy Couchbase Server (no official upstream image)"
LABEL org.opencontainers.image.version="2.5.2"
LABEL org.opencontainers.image.source="https://github.com/couchbase/docker"
LABEL org.opencontainers.image.licenses="Apache-2.0"
```

**After:**
```dockerfile
LABEL maintainer="Tien Nguyen Minh <tiennm99@outlook.com>" \
      org.opencontainers.image.title="Couchbase Server 2.5" \
      org.opencontainers.image.description="Legacy Couchbase Server (no official upstream image)" \
      org.opencontainers.image.version="2.5.2" \
      org.opencontainers.image.source="https://github.com/couchbase/docker" \
      org.opencontainers.image.licenses="Apache-2.0"
```

Apply same pattern to `scribe-2.2`, `oraclejdk-8`, `gradle-8` Dockerfiles.

Net effect: -5 layers per image; -4 LOC per Dockerfile (continuation lines balance).

### Step 2 — Workflow `name:` disambiguation

All 4 workflow files currently read `name: Publish Docker image`. GitHub Actions UI shows 4 identical names — annoying when triaging runs.

Rename each:
- `publish-couchbase-2.5.yml` L1 → `name: Publish couchbase Docker image`
- `publish-scribe-2.2.yml` L1 → `name: Publish scribe Docker image`
- `publish-oraclejdk-8.yml` L1 → `name: Publish oraclejdk Docker image`
- `publish-gradle-8.yml` L1 → `name: Publish gradle Docker image`

### Step 3 — `images:` list ordering uniformity

Code-reviewer flagged that scribe + couchbase list `ghcr.io/...` first, while oraclejdk + gradle list `${{ secrets.DOCKER_USERNAME }}/...` first. Cosmetic but worth normalizing.

Recommend: GHCR first (matches commit-history convention from scribe + couchbase rounds). Flip oraclejdk and gradle.

### Step 4 — couchbase ENV → ARG

`couchbase-2.5/Dockerfile` L8-11 declare `CB_VERSION`, `CB_RELEASE_URL`, `CB_PACKAGE`, `CB_SHA256` as `ENV`. These are used only in the install `RUN` at L20-23. They persist into the running container as no-op env vars.

Change `ENV` → `ARG` for these 4 (leave `PATH` as `ENV` — it's runtime-relevant).

```dockerfile
ARG CB_VERSION=2.5.2
ARG CB_RELEASE_URL=https://packages.couchbase.com/releases
ARG CB_PACKAGE=couchbase-server-enterprise_2.5.2_x86_64.deb
ARG CB_SHA256=27a79a65758023c34ed900e8ef8c54bab4a65f4c84b7c94359cba910800a4b19

ENV PATH=$PATH:/opt/couchbase/bin:/opt/couchbase/bin/tools:/opt/couchbase/bin/install
```

### Step 5 — `.gitignore` expansion

Append:
```
# IDEs
.idea/
.vscode/

# OS metadata
.DS_Store
Thumbs.db
```

### Commit

`chore: compact LABEL blocks, disambiguate workflow names, normalize image ordering, ENV->ARG, expand .gitignore`

Cook can split if commit-message length becomes painful.

## Success Criteria

- [ ] Each of 4 Dockerfiles has exactly 1 `LABEL` directive (multi-line OK) covering maintainer + 5 OCI keys
- [ ] `grep -c "^LABEL " {couchbase-2.5,gradle-8,oraclejdk-8,scribe-2.2}/Dockerfile` returns `1` per file
- [ ] Each workflow `name:` field is unique
- [ ] All 4 workflow `images:` blocks list GHCR before Docker Hub
- [ ] `couchbase-2.5/Dockerfile` has 4 new `ARG` lines, 4 fewer `ENV` lines for CB build constants
- [ ] `.gitignore` contains `.idea/`, `.vscode/`, `.DS_Store`, `Thumbs.db`
- [ ] One or more commits starting with `chore:` or `style:`
- [ ] `docker inspect` shows identical Config.Labels output before/after for all 4 images (same keys, same values)

## Risk Assessment

- **Risk:** Multi-line `LABEL` syntax — wrong continuation could break build. **Mitigation:** verify with `docker build` per image; first build is the validator.
- **Risk:** `ARG` change to couchbase Dockerfile means the values vanish from the running container's `env`. **Mitigation:** intentional — they were never used at runtime; if anyone WAS depending on them, that's a behavior leak this phase fixes.
- **Risk:** Renaming workflow `name:` does not change the workflow file path, so existing run history stays linked. **Mitigation:** none needed; safe rename.

## Validation

- Local build each image; `docker inspect <img> --format '{{json .Config.Labels}}' | jq .` shows all 5 OCI keys + maintainer
- After push, the 4 workflows appear in Actions UI with distinct names
