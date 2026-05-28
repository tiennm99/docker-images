# Code Review: Project Cleanup and Compaction

**Agent:** `code-reviewer` | **Date:** 2026-05-28 17:21 | **Scope:** D:\tiennm99\docker-images, read-only project scan after Convention B + OCI labels rollout completed.

**Constraint from user:** cleanup or compact logic only; NO new features.

## Overall Assessment

Small, tidy repo. The four publish workflows are 98% identical and the strongest cleanup target. Dockerfiles reasonable; a few `LABEL`/layer compaction wins. Several known small drifts (maintainer label, README "Corretto" inaccuracy, `template.yml` staleness) are all easy one-liners.

## 1. Dead code / unused files

| File | Finding | Fix | Severity |
|---|---|---|---|
| `.github/workflows/template.yml` | Stale: placeholder `my-docker-hub-namespace/my-docker-hub-repository`, lacks `paths:`, lacks per-image tags/labels, no `DOCKER_USERNAME` secret. Does not reflect real workflows. | Update OR delete. | MEDIUM |
| `couchbase-2.5/Dockerfile:21` | `dos2unix` installed only to fix CRLF on `scripts/couchbase-start` (line 31). `.gitattributes` forces LF for `run` and `*.sh`, so the dance is unnecessary. | Drop `dos2unix` from apt install + drop the `dos2unix` invocation. | MEDIUM |
| `couchbase-2.5/Dockerfile:30` | `ADD` used for a local file with no URL/tar semantics. | Use `COPY scripts/couchbase-start /usr/local/bin/`. | LOW |
| `gradle-8/Dockerfile:14-57` | Two large `echo "Testing …" && which …` blocks verify tools just installed. Pure noise. Copy-paste from upstream `gradle/docker-gradle`. | Drop both `echo "Testing …"` sections + their `which` chains. | MEDIUM |
| `gradle-8/Dockerfile:91-95` | `RUN gradle --version` then `USER root` switch back. Test layer + user toggle. | Drop the test `RUN` or move smoke testing out of image build. | LOW |
| `scribe-2.2/example/logs/.gitkeep` | Empty bind-mount placeholder. Acceptable. | None. | NIT |

## 2. DRY violations across the 4 workflows

| File | Finding | Fix | Severity |
|---|---|---|---|
| All 4 publish-*.yml | Identical permissions, login steps, checkout, metadata-action, build-push, attestation. Only image/context/tags/labels differ. | Suggestion: consolidate into reusable `workflow_call` or composite action. User to decide. | MEDIUM |
| All 4 publish-*.yml | All use identical pinned SHAs across the board. No copy-paste drift. | None. | — |
| All 4 publish-*.yml | All use `name: Publish Docker image` (same name 4x). Actions UI shows 4 identical names. | Rename each `name:` to `Publish <service> Docker image`. | LOW |
| publish-scribe-2.2.yml:43-44 vs publish-couchbase-2.5.yml:43-44 etc. | Order differs: scribe/couchbase list `ghcr.io/...` first; oraclejdk/gradle list `${{ secrets.DOCKER_USERNAME }}/...` first. | Pick one order, apply to all 4. | LOW |
| publish-scribe-2.2.yml:43 (and 3 others) | Uses `ghcr.io/${{ github.actor }}/scribe` for image + `subject-name`. `github.actor` is triggering user, not repo owner. If anyone else runs `workflow_dispatch`, image publishes under a different namespace. | Replace `${{ github.actor }}` with `${{ github.repository_owner }}` (or hardcode `tiennm99`) in all 4 workflows, both `images:` and `subject-name:`. | **HIGH** |

## 3. DRY violations across Dockerfiles

| File | Finding | Fix | Severity |
|---|---|---|---|
| `scribe-2.2/Dockerfile:6-7` and `:78-79` | Identical CentOS-vault `sed` runs in builder + runtime stages. Acceptable for multi-stage. | None. | NIT |
| `scribe-2.2` and `gradle-8` Dockerfiles | Both `sed` CentOS vault repos but different regex patterns. | Align if convenient. No functional bug. | NIT |
| `couchbase-2.5/Dockerfile:3` `Tien Nguyen Minh` vs `oraclejdk-8/Dockerfile:3` and `gradle-8/Dockerfile:3` `miti99`; `scribe-2.2/Dockerfile:70` `Tien Nguyen Minh` | Maintainer label drift. | Unify. | MEDIUM |
| All 4 Dockerfiles | Each repeats 6 `LABEL` lines. Best practice: combine. | Squash into single multi-line `LABEL` per image. | LOW |

## 4. README simplification

| File:line | Finding | Fix | Severity |
|---|---|---|---|
| `README.md:10` | "JDK 8 + Corretto base" — but `gradle-8/Dockerfile:1` is `FROM ghcr.io/tiennm99/oraclejdk:8u201`. Not Corretto. | Change to "JDK 8 (Oracle) base". | MEDIUM |
| `README.md:28-68` | Per-image sections under Couchbase/Gradle/Oracle JDK contain only Credits. Gradle and Oracle JDK have no usage. Inconsistent. | Either add `docker pull` snippets, or drop sections. | LOW |
| `README.md:14` | "Naming convention" verbose for a 4-image repo. | Trim to 2-3 lines. | NIT |
| `README.md:70-74` | "GitHub Actions template" section linked to template.yml. | Tie to template.yml decision. | LOW |
| `scribe-2.2/README.md:107` | Facebook archived URL works but content is frozen. | Add "(archived)" inline. | NIT |

## 5. Inconsistencies

- Maintainer drift (see §3)
- README "Corretto" inaccuracy (see §4)
- Workflow string quoting: consistent `'main'`, `'scribe-2.2/**'`. No issue.
- `.gitattributes:3-4` covers entry scripts. Verified LF.
- `images:` ordering inconsistent (see §2).

## 6. Workflow triggers

All 4 `paths:` triggers correctly scope per-subdir. No stale paths. No issues.

## 7. Dockerfile compaction

| File:line | Finding | Fix | Severity |
|---|---|---|---|
| `couchbase-2.5/Dockerfile:13-18` | 4 `ENV` lines used only at build time become persistent runtime env vars. | Convert to `ARG`. Keep `PATH` as `ENV`. | LOW |
| `gradle-8/Dockerfile:12,74` | `GRADLE_HOME` and `GRADLE_VERSION` declared far from use. Stylistic. | Optional move. | NIT |
| `oraclejdk-8/Dockerfile:24` | `JAVA_SHA256="$(curl -L "$JAVA_PKG".sha256)"` — `curl` exit code swallowed by command substitution. If URL 404s, build silently proceeds with empty checksum. Pre-existing safety issue. | Split: `curl -fL -o /tmp/jdk.sha256 "$JAVA_PKG".sha256; JAVA_SHA256="$(cat /tmp/jdk.sha256)"`. `-f` propagates HTTP errors. | **HIGH** |
| `oraclejdk-8/Dockerfile:25,28` | `set -eux` then mixed `;` and `&&` separators in same `RUN`. Inconsistent. | Use `&&` throughout, or `set -e` + `;`. | LOW |
| `oraclejdk-8/Dockerfile:24` | `echo "$JAVA_SHA256" */tmp/jdk.tar.gz | sha256sum -c -` — suspicious `*` prefix. Depends on whether Huawei `.sha256` already includes `*filename`. | Inspect actual mirror content, adjust. | MEDIUM |
| `couchbase-2.5/Dockerfile:20-23` | `apt-get autoremove` without `-y`. | Add `-y`. | LOW |
| All Dockerfiles | LABEL line layer creation (modern BuildKit folds, but still recommended). | See §3. | LOW |

## 8. Repo-level hygiene

| File | Finding | Fix | Severity |
|---|---|---|---|
| `.gitignore` | No `.idea/`, `.vscode/`, `.DS_Store`, `Thumbs.db`. | Add common IDE/OS entries. | LOW |
| `LICENSE` | Apache-2.0 matches `org.opencontainers.image.licenses=Apache-2.0` on 3 images; oraclejdk correctly uses `LicenseRef-Oracle-BCL`. | None. | — |
| `.gitattributes:1` | `* text=auto` + `*.sh text eol=lf` is fine. | None. | — |
| Tracked files | No IDE/OS metadata or build artifacts visible. | None. | — |

## Edge cases found

- `${{ github.actor }}` vs `${{ github.repository_owner }}` namespace divergence (HIGH).
- `oraclejdk-8` curl-of-sha256 silent failure (HIGH).
- `oraclejdk-8` possible double-`*` in checksum input (MEDIUM).
- `couchbase-2.5` `dos2unix` redundant given `.gitattributes` (MEDIUM).

## Positive observations

- All 4 workflows use identical, current pinned SHAs — no drift.
- `paths:` triggers correctly scope per-subdir.
- `scribe-2.2/Dockerfile` is a clean two-stage build.
- Go sender self-contained, well-commented.
- `org.opencontainers.image.licenses` correctly distinguishes Apache-2.0 vs Oracle BCL.
- `.gitattributes` enforces LF for shell scripts.

## Top 5 cleanups by ROI

1. **Replace `${{ github.actor }}` with `${{ github.repository_owner }}`** in 4 workflows (HIGH, ~5 min). Prevents accidental namespace leak.
2. **Fix `oraclejdk-8/Dockerfile:24` swallowed curl error** (HIGH, ~5 min). Two-line change.
3. **Normalize maintainer label** across 4 Dockerfiles (MEDIUM, ~2 min).
4. **Fix `README.md:10` "Corretto" → "Oracle JDK 8"** (MEDIUM, ~30 sec).
5. **Decide `template.yml` fate** (MEDIUM, ~10 min): update or delete + remove README section.

## Out of scope (NEW features only, NOT in cleanup scope)

- Reusable `workflow_call` to deduplicate 4 publish workflows (suggestion-only; structural refactor — see Phase 4 of plan).
- SBOM generation / Trivy scan.
- Multi-arch image builds.
- Dependabot for action SHA updates.

## Unresolved questions

- `template.yml`: delete or kept-updated?
- Huawei `.sha256` mirror format: does it include `*filename` or just hash?
- Maintainer canonical: `miti99` (GHCR namespace match) or `Tien Nguyen Minh` (human-readable)?
