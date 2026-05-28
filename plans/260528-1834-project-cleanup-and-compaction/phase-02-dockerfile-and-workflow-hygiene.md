---
phase: 2
title: "Dockerfile and workflow hygiene"
status: pending
priority: P2
effort: "30m"
dependencies: [1]
---

# Phase 2: Dockerfile and workflow hygiene

## Overview

Independent file-level cleanups: unify maintainer label, drop dead `dos2unix` dance from couchbase, drop noisy "Testing" echo chains from gradle, decide template.yml fate. Blocked on Phase 1's correctness fixes (avoids interleaving safety fixes with hygiene churn) and on user decisions on open questions.

## Requirements

- Functional: all images still build and publish with identical contracts (image name, tags, OCI labels, layer-count may go down).
- Non-functional: maintainer label consistent across all 4 Dockerfiles; obsolete tooling/lines removed.

## Open questions (must answer before cook)

1. Maintainer canonical: `miti99 <tiennm99@outlook.com>` or `Tien Nguyen Minh <tiennm99@outlook.com>`? Defaulting in this phase to **`Tien Nguyen Minh <tiennm99@outlook.com>`** since 2 of 4 Dockerfiles + most git history use it. Cook can override.
2. `template.yml` fate: **delete** (since cleanup #7 finds it stale and unused) or update? Defaulting to **delete** + remove the README section that points at it.

## Architecture

Four independent edits:
- `couchbase-2.5/Dockerfile`: drop `dos2unix` from apt install list, drop `dos2unix` shell step, switch `ADD` → `COPY` for local script
- `gradle-8/Dockerfile`: drop 2 long `echo "Testing …" && which …` blocks (lines ~33-50 in current file)
- All 4 Dockerfiles: unify maintainer label
- `.github/workflows/template.yml`: delete
- `README.md`: remove "## GitHub Actions template" section (currently L70-74)

## Related Code Files

- Modify: `D:\tiennm99\docker-images\couchbase-2.5\Dockerfile`
- Modify: `D:\tiennm99\docker-images\gradle-8\Dockerfile`
- Modify: `D:\tiennm99\docker-images\oraclejdk-8\Dockerfile`
- Modify: `D:\tiennm99\docker-images\scribe-2.2\Dockerfile`
- Delete: `D:\tiennm99\docker-images\.github\workflows\template.yml`
- Modify: `D:\tiennm99\docker-images\README.md`

## Implementation Steps

### Step 1 — Unify maintainer label

Default canonical: `Tien Nguyen Minh <tiennm99@outlook.com>`.

- `couchbase-2.5/Dockerfile:3` — already `Tien Nguyen Minh <tiennm99@outlook.com>`, no change
- `scribe-2.2/Dockerfile:70` — already `Tien Nguyen Minh <tiennm99@outlook.com>`, no change
- `oraclejdk-8/Dockerfile:3` — change `miti99 <tiennm99@outlook.com>` → `Tien Nguyen Minh <tiennm99@outlook.com>`
- `gradle-8/Dockerfile:3` — change `miti99 <tiennm99@outlook.com>` → `Tien Nguyen Minh <tiennm99@outlook.com>` (added in earlier session commit `54e1481`)

### Step 2 — couchbase Dockerfile dos2unix removal

Verify first: `git show HEAD:couchbase-2.5/scripts/couchbase-start | file -` (LF or CRLF?). `.gitattributes` has `run text eol=lf` so it should be LF on checkout.

- `couchbase-2.5/Dockerfile:16` — drop `dos2unix` from apt install list
- `couchbase-2.5/Dockerfile:25-27` — change:

  **Before:**
  ```dockerfile
  ADD scripts/couchbase-start /usr/local/bin/
  RUN dos2unix /usr/local/bin/couchbase-start && \
      chmod +x /usr/local/bin/couchbase-start
  ```

  **After:**
  ```dockerfile
  COPY scripts/couchbase-start /usr/local/bin/
  RUN chmod +x /usr/local/bin/couchbase-start
  ```

### Step 3 — gradle Dockerfile noise removal

`gradle-8/Dockerfile` lines ~32-50 contain:

```dockerfile
    && echo "Testing common utilities" \
    && which awk \
    && which curl \
    && which cut \
    && which grep \
    && which gunzip \
    && which sha256sum \
    && which sed \
    && which tar \
    && which tr \
    && which unzip \
    && which wget \
    \
    && echo "Testing VCSes" \
    && which git \
    && which git-lfs \
    && which hg \
    && which svn
```

Drop both blocks. The preceding `apt-get install` already errors if any package fails — these `which` checks are belt-and-suspenders inherited from upstream `gradle/docker-gradle` and add zero value here.

Also consider dropping the trailing version-check `RUN` near the end (`RUN ... && gradle --version`) per agent finding §1. Cook decides — if dropped, also drop the `USER root` reset that follows it.

### Step 4 — template.yml deletion

- `git rm .github/workflows/template.yml`
- `README.md` — remove the "## GitHub Actions template" section (currently L70-74), including the credits link bullet underneath.

### Commit

`refactor: unify maintainer label, drop dead dos2unix and gradle noise, remove stale template`

## Success Criteria

- [ ] `grep "LABEL maintainer" {couchbase-2.5,gradle-8,oraclejdk-8,scribe-2.2}/Dockerfile` shows identical value on all 4
- [ ] `grep "dos2unix" couchbase-2.5/Dockerfile` returns 0 results
- [ ] `couchbase-2.5/Dockerfile` uses `COPY` for the script, not `ADD`
- [ ] `grep -c "which " gradle-8/Dockerfile` significantly lower (target: 0 — verify gradle still builds without those)
- [ ] `.github/workflows/template.yml` does not exist
- [ ] `README.md` has no "## GitHub Actions template" section
- [ ] One commit (or split per file if cook prefers) starting with `refactor:`
- [ ] All 4 images still build locally without errors

## Risk Assessment

- **Risk:** `couchbase-start` script actually has CRLF on disk despite `.gitattributes`. **Mitigation:** verify before commit; if CRLF, also fix the script with a one-shot dos2unix conversion + commit, then drop the runtime dependency.
- **Risk:** gradle build relies on one of those `which` checks failing-fast to catch a missing dependency. **Mitigation:** the `yum install -y` step above already fails if any package is missing. Drop is safe.
- **Risk:** deleting `template.yml` removes a discoverable reference for future contributors. **Mitigation:** Phase 4's reusable workflow (`_publish-image.yml`) IS the new template; per-image files become the example. Document this transition in Phase 4 README updates.

## Validation

- Local build all 4 images: `docker build {couchbase-2.5,gradle-8,oraclejdk-8,scribe-2.2}/ -t local-test:<svc>`; verify exit code 0 each. `docker rmi` after.
- `docker inspect local-test:couchbase --format '{{.Config.Labels.maintainer}}'` returns the unified value.
- Push and watch 4 CI runs complete green.
