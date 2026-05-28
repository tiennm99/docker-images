# Image Naming Convention - Brainstorm

## Current state (verified)

| Subdir | GHCR image | Tags | Notes |
|---|---|---|---|
| couchbase-2.5/ | `couchbase-2.5` | `latest`, `2.5.2` | matches subdir |
| gradle-8/ | `gradle-8` | `latest`, `8.13` | matches subdir |
| oraclejdk-8/ | `oraclejdk-8` | `latest`, `8u201` | matches subdir |
| scribe-2.2/ | `scribe` | `latest`, `2.2` | **outlier** |

Scribe diverges in `publish-scribe-2.2.yml` L43-44 (`images:`), L61 (`subject-name:`), and `README.md` L12 (package link). Subdir already correctly named `scribe-2.2/`. Workflow first added ~7 weeks ago - package fresh, blast radius low.

## Convention

`<service>-<major-line>` = **image name** (preserved product line). `<full-version>` = **tag** (specific build). Differs from Docker Hub officials (`python`, `node`) because this repo specializes in archived/legacy software - each major line is its own write-once artifact with potentially divergent Dockerfile/base OS, not a sliding window. Scribe is same shape (archived 2.2, no 2.3 ever) - line name belongs in image name, not just tag.

## Options for scribe

### Option A - Rename to `scribe-2.2` (full alignment)

- Changes: workflow L43-44 + L61, README L12, possibly `scribe-2.2/README.md` and `scribe-2.2/example/*` if any pull-URL refs
- Pros: consistent with 3/3 siblings; reserves `scribe` for hypothetical fork; alphabetical sort lines up; one rule to teach; future-proof
- Cons: breaks any existing `docker pull ghcr.io/tiennm99/scribe:2.2`; old GHCR + Docker Hub packages orphaned (inert)
- Blast radius: low. 7 weeks old, no known external consumers

### Option B - Keep as `scribe` (document exception)

- Pros: zero migration, no orphans, shorter name
- Cons: convention dies on first stress test; future contributors must memorize special case; irregular sort; retroactively ugly if `scribe-3.x` appears
- Blast radius: zero now, debt later

### Option C - Dual-publish for transition

- Pros: smoothest external transition
- Cons: doubles artifacts; sunset date easy to forget; complexity unjustified for a personal legacy archive
- Blast radius: medium

### Option D - Drop suffix from all siblings (inverted)

- Reject: breaks multi-major-line model; requires renaming 3 existing packages; strictly worse

## Trade-off matrix

| Dim | A (rename) | B (keep) | C (dual) | D (drop suffix) |
|---|---|---|---|---|
| Consistency | full | one outlier | during transition | inverted |
| Migration cost | low (~5 edits) | zero | low + ongoing | high (3 pkgs) |
| Orphans | 2 | 0 | 2 eventually | 6+ |
| Future-proof | yes | no | yes | no |
| Discoverability | alphabetical | irregular | irregular | alphabetical |
| Teach cost | 1 rule | 1 rule + exception | 1 rule + sunset | 1 worse rule |

## Recommendation: Option A - rename to `scribe-2.2`

- Package 7 weeks old, no known external consumers - blast radius minimal
- ~5 line edits, one workflow run
- Establishes rule cleanly on first stress test
- Orphan packages inert; deletable in 30 seconds via UI

## `latest` tag policy - keep unchanged

For write-once archived images, `latest` is semantically equivalent to canonical version tag. Usual objection ("silently shifts") doesn't apply - nothing shifts. Ergonomic default for casual `docker pull`. Edge case: CVE rebuild would shift latest, but consumers can pin to version tag.

## Repo-wide rule to codify

Add to `README.md` (new "Naming convention" section) and `CLAUDE.md`:

```
Naming convention for new images
- Subdirectory:    <service>-<major-line>          e.g.  redis-3/, mongodb-3.6/
- GHCR + Docker Hub image name: <service>-<major-line>  (= subdir name)
- Workflow file:   .github/workflows/publish-<service>-<major-line>.yml
- Tags published:
    - latest                  (alias for canonical/only version)
    - <full-version>          e.g. 3.6.23, 7.4.5, 8u201
- All three (subdir, image name, workflow basename) MUST match.
```

Rule: **subdir = image name = workflow basename**, all three derived from `<service>-<major-line>`.

## Risks / considerations

- Attestation `subject-name:` (L61) must also update - easy to miss
- Orphan packages harmless but visible in GHCR UI - decide: delete vs tombstone in README
- Grep `scribe-2.2/` for any `ghcr.io/tiennm99/scribe` literal before committing
- Commit message: `refactor(scribe): align image name with subdir naming convention`

## Success metrics

- `gh api /users/tiennm99/packages/container/scribe-2.2` returns 200 after first workflow run
- README image-index links resolve
- Naming-convention section exists in README, referenced from CLAUDE.md
- Next image added follows rule without prompting

## Next steps

1. User confirms Option A (or picks B/C)
2. If A: edit workflow + README + grep/fix any pull-URL refs in scribe-2.2/
3. Codify naming rule in README + project CLAUDE.md
4. Decide on orphan deletion vs retention

## Unresolved questions

1. Delete orphaned `ghcr.io/tiennm99/scribe` + `docker.io/tiennm99/scribe` after rename, or tombstone in README? (Lean: delete.)
2. Does `scribe-2.2/example/docker-compose.yml` reference the published image, or only local build context? Quick grep needed before rename.
3. README image-index table shows GHCR only - add Docker Hub coords too? (Adjacent, out of scope.)
4. `gradle-8` / `oraclejdk-8` GHCR links missing from README (show `-`) - intentional pending first publish, or oversight?
