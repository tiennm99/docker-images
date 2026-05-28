# Scribe-2.2 Cleanup Verification Review

**Date:** 2026-05-28
**Scope:** Uncommitted changes in `scribe-2.2/` vs HEAD
**Reviewer:** code-reviewer

## Files Reviewed

- `scribe-2.2/Dockerfile`
- `scribe-2.2/README.md`
- `scribe-2.2/config/scribe.conf`
- `scribe-2.2/example/README.md`
- `scribe-2.2/example/docker-compose.yml`
- `scribe-2.2/example/sender/main.go`
- `scribe-2.2/scripts/docker-entrypoint.sh`

## Overall Assessment

Pure cosmetic cleanup. No functional regressions. All build steps, runtime behavior, ports, volumes, users, configs unchanged. One concern: removed metadata LABELs.

## Critical Issues

None.

## High Priority

None.

## Medium Priority — Concerns (do not auto-revert)

### C1. Removed LABEL metadata in runtime stage

- HEAD had: `LABEL maintainer`, `LABEL description`, `LABEL version="2.2"`
- Working tree keeps only `LABEL maintainer`
- `description` and `version` labels removed
- Impact:
  - `docker inspect scribe:2.2` no longer reports version/description
  - Image registries (Harbor, GHCR) lose searchable description
  - Image scanners / SBOM tools lose declared version field
  - `org.opencontainers.image.version` / `.description` were not used (only plain `version`/`description`), so OCI compliance was already partial
- Status: **flag as intentional-or-regression question for user**. Image tag (`scribe:2.2`) still encodes the version, so loss is minor. But `description` had documentation value
- Recommended (if user wants to re-add): use OCI-standard labels:
  - `org.opencontainers.image.version="2.2"`
  - `org.opencontainers.image.description="Facebook Legacy Scribe - real-time log aggregation daemon"`
  - `org.opencontainers.image.source="https://github.com/facebookarchive/scribe"`

## Low Priority

None.

## Per-Question Verification

### 1. Dockerfile build steps & package list

- **Build steps:** no removals/alterations. Only banner-comment lines and inline `# ──` comments removed. All RUN blocks identical
- **Builder yum list:** same 14 packages, sorted alphabetically:
  - HEAD: `gcc, gcc-c++, make, autoconf, automake, libtool, git, wget, flex, bison, boost-devel, libevent-devel, openssl-devel, zlib-devel` (14)
  - Working: `autoconf, automake, bison, boost-devel, flex, gcc, gcc-c++, git, libevent-devel, libtool, make, openssl-devel, wget, zlib-devel` (14)
  - Set-equal, count-equal. One-package-per-line improves readability
- **Runtime yum list:** unchanged (boost-filesystem, boost-system, boost-thread, libevent, openssl)
- **EXPOSE 1463, USER scribe, VOLUME, ENTRYPOINT:** all unchanged

### 2. Removed LABELs

- `LABEL maintainer="Tien Nguyen Minh <tiennm99@outlook.com>"` still present
- `LABEL description=...` and `LABEL version="2.2"` removed
- See C1 above. Flagged, not auto-reverted

### 3. Go file `example/sender/main.go`

Four diff hunks, all non-logic:
- L5 (file header comment): em-dash to hyphen
- L125 comment: `messageBegin — detect` to `messageBegin: detect`
- L161 comment: `no return field present → treat as OK` to `no return field present; treat as OK`
- L211 string literal (user-visible): `Hello from Go sender — message %d/%d` to `Hello from Go sender - message %d/%d`

No logic touched. No control flow, no struct fields, no Thrift wire-encoding bytes, no retry logic, no timeouts changed.

**README sync check:** `example/README.md:41` already updated to `Hello from Go sender - message N/5`. Consistent with main.go output. Confirmed via grep — only the two expected occurrences exist, both hyphen.

### 4. `config/scribe.conf`

Only comment formatting:
- L1 header em-dash to hyphen
- 4 section-header `# ── Name ──────` lines simplified to `# Name`

Directives unchanged: `port=1463`, `max_msg_per_second=100000`, `max_msg_size=16777216`, `check_interval=3`, default store block (file, /var/log/scribe, hourly rotation, write_meta=true). No store directive added/removed.

### 5. Em-dash semantic breakage

Reviewed all em-dash to hyphen and em-dash to `:`/`;` replacements:

- `example/README.md:51` `Message category → log subdirectory name` to `Message category to log subdirectory name` — arrow replaced with "to". Reads naturally
- `example/README.md:62` `[link] — see the [other link]` to `[link] - see the [other link]` — fine
- `main.go:161` `→ treat as OK` to `; treat as OK` — semantically equivalent (consequence vs sequence). OK
- `main.go:125` `messageBegin — detect` to `messageBegin: detect` — annotation style. OK
- All other em-dash to hyphen replacements are parenthetical-aside style. No semantic loss

No broken Markdown links, no broken Mermaid/code-fence boundaries.

### 6. Residual Unicode

Grep for `[‐-―─-╿‘’“”…←-↓]` across all of `scribe-2.2/`: **0 matches**. Cleanup is complete.

Also grep for `→ — – ─ │ ┌ └ ├ ┤`: **0 matches** in any of the 7 modified files.

## Positive Observations

- Alphabetized yum package list improves diff readability for future deps changes
- One-package-per-line aids future bisect of build failures
- ASCII-only sources improve grep/sed/IDE tooling on Windows where Unicode rendering varies
- Go string literal change keeps user-visible message ASCII — safer for log-grepping pipelines downstream
- No commit yet, so easy to fold back C1 if user agrees

## Trust-Boundary / Security Review

- No auth changes, no input validation paths touched
- No new env vars, no new ports, no new file paths
- Non-root `scribe` user retained
- No new ADD/COPY from external sources
- No secrets in diff

## Metrics

- Files changed: 7
- Lines net: comment/whitespace-only delta
- Functional code changes: 0
- Logic changes: 0
- Build-graph changes: 0
- Removed metadata labels: 2

## Recommended Actions

1. Decide on C1 (LABEL description/version removal):
   - **Keep removed** — accept minor metadata loss
   - **Re-add** — preferably as OCI-standard `org.opencontainers.image.*` labels
2. If C1 accepted as-is, commit cleanup as `style(scribe): normalize ASCII formatting in scribe-2.2/`

## Unresolved Questions

1. Was removal of `LABEL description` and `LABEL version` intentional, or collateral damage during banner-comment cleanup? Image tag still carries version, so functional impact is nil, but discoverability via `docker inspect` is lost.
