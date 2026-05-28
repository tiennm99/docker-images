---
phase: 1
title: "High-severity correctness fixes"
status: pending
priority: P1
effort: "15m"
dependencies: []
---

# Phase 1: High-severity correctness fixes

## Overview

Three correctness bugs / factual errors found by the review agents. Tiny diffs, immediate safety win, blocks Phase 4 (workflow refactor would otherwise inherit the `github.actor` bug).

## Requirements

- Functional: workflows publish to the repo owner's namespace regardless of who triggers `workflow_dispatch`; oraclejdk Dockerfile fails fast if the JDK sha256 URL 404s; README accurately states the Gradle base.
- Non-functional: no change to image contracts, label values, tags, or subdirectory layout.

## Architecture

Three independent file edits:
- Workflows: 4 files, swap `${{ github.actor }}` → `${{ github.repository_owner }}` everywhere (8 sites: 2 per workflow — `images:` block + `subject-name:`).
- Dockerfile: 1 file, replace the swallowed-exit-code curl with a two-step fetch-then-read using `curl -fL`.
- README: 1 file, one-line text fix.

## Related Code Files

- Modify: `D:\tiennm99\docker-images\.github\workflows\publish-couchbase-2.5.yml`
- Modify: `D:\tiennm99\docker-images\.github\workflows\publish-gradle-8.yml`
- Modify: `D:\tiennm99\docker-images\.github\workflows\publish-oraclejdk-8.yml`
- Modify: `D:\tiennm99\docker-images\.github\workflows\publish-scribe-2.2.yml`
- Modify: `D:\tiennm99\docker-images\oraclejdk-8\Dockerfile`
- Modify: `D:\tiennm99\docker-images\README.md`

## Implementation Steps

### Step 1 — `github.actor` → `github.repository_owner` (4 workflow files)

For each of the 4 `publish-*.yml` files, replace **all** occurrences of `${{ github.actor }}` with `${{ github.repository_owner }}`.

Sites per workflow:
- `metadata-action.images:` block — the `ghcr.io/${{ github.actor }}/<service>` line
- `attest-build-provenance.subject-name:` — `ghcr.io/${{ github.actor }}/<service>`
- `docker/login-action` (ghcr.io step) `username:` — keep `${{ github.actor }}` here, since this IS the auth user, not the repo namespace. (Important distinction: login uses the triggering actor's token; the repo path uses the owner.)

Wait — re-verify: `docker/login-action` uses `${{ github.actor }}` paired with `${{ secrets.GITHUB_TOKEN }}`. The `GITHUB_TOKEN` is scoped to the workflow run, not the actor, so leaving `github.actor` for login is fine and conventional. The fix is ONLY for the two image-path sites.

### Step 2 — Fix oraclejdk Dockerfile curl exit-code swallow

`oraclejdk-8/Dockerfile` currently at L19-21:

```dockerfile
    JAVA_SHA256="$(curl -L "$JAVA_PKG".sha256)" ;  \
    curl -L --output /tmp/jdk.tar.gz "$JAVA_PKG" && \
    echo "$JAVA_SHA256" */tmp/jdk.tar.gz | sha256sum -c -; \
```

Replace with:

```dockerfile
    curl -fL --output /tmp/jdk.sha256 "$JAVA_PKG".sha256 && \
    JAVA_SHA256="$(cat /tmp/jdk.sha256)" && \
    curl -fL --output /tmp/jdk.tar.gz "$JAVA_PKG" && \
    echo "$JAVA_SHA256" */tmp/jdk.tar.gz | sha256sum -c - && \
    rm -f /tmp/jdk.sha256 && \
```

Changes: `-f` flag on both curls so HTTP errors propagate exit code; sha256 stored in tmp file then read, so the `curl` exit code isn't swallowed by command substitution; both curls now chain with `&&` for consistent fail-fast.

Open question for cook: does the Huawei `.sha256` mirror return the file content as `<hash> *filename` or just `<hash>`? If it's just `<hash>`, the existing `echo "$JAVA_SHA256" */tmp/jdk.tar.gz | sha256sum -c -` line is correct. If it includes `*filename` already, the `*/tmp/jdk.tar.gz` prefix double-pollutes. Verify with `curl -fL https://repo.huaweicloud.com/java/jdk/8u201-b09/jdk-8u201-linux-x64.tar.gz.sha256` before committing.

### Step 3 — README L10 "Corretto" fix

`README.md` L10:

**Before:**
```
| Gradle 8 | `8.13` | [gradle](https://github.com/tiennm99/docker-images/pkgs/container/gradle) | JDK 8 + Corretto base |
```

**After:**
```
| Gradle 8 | `8.13` | [gradle](https://github.com/tiennm99/docker-images/pkgs/container/gradle) | Oracle JDK 8 base |
```

### Commit

`fix: namespace via repository_owner, harden oraclejdk sha256 fetch, correct gradle base in README`

Or split into two:
- `fix(workflows): use repository_owner instead of actor for image namespace`
- `fix(oraclejdk): propagate curl exit code on sha256 fetch`
- `docs(readme): correct gradle base description`

Cook decides.

## Success Criteria

- [ ] All 4 `publish-*.yml` files have ZERO `${{ github.actor }}` in `images:` and `subject-name:` (only the ghcr.io login `username:` line retains `github.actor`)
- [ ] `oraclejdk-8/Dockerfile` uses `curl -fL` for both fetches; sha256 stored via intermediate file (no command-substitution swallow)
- [ ] README L10 Gradle row note reads "Oracle JDK 8 base"
- [ ] `grep "github.actor" .github/workflows/publish-*.yml` returns only the 4 ghcr.io login lines (one per workflow)
- [ ] One commit OR three semantic commits (cook chooses)
- [ ] Post-push: each workflow run still publishes to `ghcr.io/tiennm99/<service>:<tags>` (no namespace flip from `tiennm99` because owner == actor for `push` triggers from main)

## Risk Assessment

- **Risk:** Huawei `.sha256` URL format mismatch breaks the existing checksum line. **Mitigation:** verify URL format manually before committing (see Step 2 open question). If it has `*filename` already, drop the `*/tmp/jdk.tar.gz` prefix.
- **Risk:** `repository_owner` resolves differently than `actor` if the repo is forked AND a fork user triggers a workflow. **Mitigation:** that's exactly the point — fork users SHOULD NOT publish under the upstream namespace.
- **Risk:** the README change is visible to anyone reading the repo; if Convention B was wrong about Oracle base, this would be wrong too. **Mitigation:** `gradle-8/Dockerfile:1` is `FROM ghcr.io/tiennm99/oraclejdk:8u201` — that's authoritative; "Oracle JDK 8" is correct.

## Validation

- After push, manually trigger `workflow_dispatch` on one workflow (`publish-scribe-2.2.yml` is cheapest) and confirm publish lands at `ghcr.io/tiennm99/scribe`, not `ghcr.io/<actor>/scribe`.
- Confirm `oraclejdk-8` builds locally with `docker build oraclejdk-8/ -t oraclejdk:fix-verify`; rmi after.
