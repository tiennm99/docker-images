# Image Naming Convention - Re-brainstorm (A vs B, no scribe special-casing)

## TL;DR

**Recommend Convention A** (`<service>-<major-line>:<full-version>`). Same conclusion as previous brainstorm, narrower reasoning. Single strongest argument: `latest` semantics break under B for archived software.

## The two conventions

| Convention | Image name | Tag | Example | Used by |
|---|---|---|---|---|
| A | `<service>-<major-line>` | `<full-version>` | `couchbase-2.5:2.5.2` | 3 of 4 images in repo (current) |
| B | `<service>` | `<full-version>` | `couchbase:2.5.2` | scribe (current); Docker Hub officials (`python`, `redis`, `postgres`) |

## Steel-mans

### A — museum collection
Each major line is a write-once artifact: different upstream source, often different base OS, separate EOL/CVE story. Per-package `latest`, README, badge, attestation history. GHCR list IS the catalog.

### B — library catalog
One package per service, tag namespace handles divergence (`3.11-bookworm`, `3.11-alpine`). Adding a new major line is `new subdir + new tag` instead of `new package`. Matches `docker pull python:3.11` muscle memory.

## Where the choice actually matters

### Adding Couchbase 6 hypothetically

| Aspect | A | B |
|---|---|---|
| GHCR discoverability | 2 rows side-by-side | 1 row, tags hidden |
| README image-index table | grows linearly | stays narrow |
| GHCR package list | growing alphabetical | compact cluster |
| Workflow files | 1:1 with subdir | 1:1 subdir, but **multiple workflows push to same package** |
| Attestation `subject-name` | per-line lineage | mixed lineage |
| `latest` semantics | bounded per line | ambiguous - newest wins? meaningless? |
| Adding a build | heavyweight (new package) | lightweight (new tag) |

### Multiple builds within one line (alpha/stable)
Tag suffixes work in both. Wash.

### Base-OS divergence
Docker officials use tag suffixes (`python:3.11-alpine`). But that pattern is for SAME software on different runtimes. Here, `couchbase-2.5` (CentOS) vs `couchbase-6` (Ubuntu) is DIFFERENT software trees - more like `python` vs `pypy`. Analogy doesn't transfer.

## `latest` tag - the decisive argument

### Under A
`latest` per package = canonical build of that frozen line. Stable, bounded, doesn't drift across major lines.

### Under B
- `couchbase:latest` = newest tag wins (Docker Hub default). User landing here for legacy 2.5 gets 6.6.0. **Misleading.**
- Fix: don't publish `latest` under B. Every README example must spell out a version tag. No quick-pull demo.
- Alternative: pin `latest` to a specific version. Which? "Oldest" is weird; "most-popular" needs telemetry; "last added" is arbitrary.

A's `latest` is mildly defensible. B's `latest` is broken or must be dropped. Net point for A.

## Workflow / CI

- A: each workflow owns its package. Zero concurrency risk between major lines.
- B: multiple workflows push to same package on different tags. Tag uniqueness avoids data loss, but:
  - `latest` pointer flaps if both publish `latest` (must disable on all but one)
  - Attestation history mixes lines - harder to audit
  - Concurrency groups need per-tag rescoping

Net point for A.

## Migration cost if B were chosen

| File | Edit |
|---|---|
| 3x `publish-*.yml` (couchbase, gradle, oraclejdk) | strip `-<major>` from `images:` and `subject-name:` |
| `publish-scribe-2.2.yml` | no change (already B-style) |
| `README.md` | image-index table + pull examples + package links |
| 4x `<subdir>/README.md` | regrep for pull URLs |
| 4x `<subdir>/example/*` | regrep for pull URLs |
| Tag policy | drop `latest` or pin per-service - new decision needed |
| GHCR + Docker Hub cleanup | 6 orphan packages |

~12-15 edits + new tag-policy decision + 6 orphans + concurrency rescoping. **~3x cost of A-renames-scribe**, lands in worse place.

A-route: ~5 edits, no tag-policy question, 2 orphans.

## Why A wins (5 reasons)

1. **`latest` breaks under B for archives.** Strongest argument. A's `latest` per-line is stable; B's is misleading or must be dropped.
2. **Attestation lineage per-line under A, mixed under B.** Matters for supply-chain trust on a repo whose whole point is "I built these, trust them".
3. **Docker-official analogy doesn't transfer.** Officials are sliding-window for maintained software; this repo is a museum of frozen artifacts. Different domain.
4. **Migration asymmetry.** Don't pay 3x to land in worse place.
5. **Per-line discoverability matches the artifact model.** Each exhibit gets its own plinth. GHCR package list IS the catalog.

The one cost of A - divergence from Docker official conventions - is a feature here. Visitors are specifically looking for legacy software; seeing 4 distinct rows is more discoverable than one mystery `couchbase` with hidden tag list.

## Codified rule

```
## Naming convention

This repo packages archived/legacy software. Each major version line is a
write-once artifact with its own base OS, Dockerfile, attestation history,
and `latest` pointer.

Rule:
  Subdirectory     :  <service>-<major-line>/
  GHCR image       :  <service>-<major-line>
  Docker Hub image :  <service>-<major-line>
  Workflow file    :  publish-<service>-<major-line>.yml
  Tags published   :
                      - <full-version>      e.g. 2.5.2, 8.13, 8u201, 2.2
                      - latest              (alias for canonical build of THIS line)

Constraint:
  subdir == GHCR image == Docker Hub image == workflow basename. No exceptions.

Why not <service>:<full-version>?
  This repo is a museum, not a sliding window. `latest` per major line is
  meaningful and stable. `latest` across a service would be misleading
  (newest version wins, even when users came here specifically for the
  archived line). Attestation lineage is also cleaner per line.

New major line of an existing service?
  New subdir, new workflow, new GHCR + Docker Hub package. Not a new tag
  under the existing package.
```

## Unresolved questions

1. `latest` on CVE rebuild - shifts digest. Document "pin by version tag for reproducibility"?
2. Orphan cleanup for `scribe` + `scribe:latest` - delete vs tombstone? (Lean: delete.)
3. Docker Hub coords in README image-index table - add column? (Defer until rename lands.)
4. `gradle-8` / `oraclejdk-8` GHCR links show `-` - pending first publish or never triggered?
5. Codified rule in README only, or README + CLAUDE.md? (Lean: both.)
