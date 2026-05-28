---
phase: 1
title: "Codify naming rule and supersede prior plans"
status: completed
priority: P2
effort: "15m"
commit: "2f7bfde"
dependencies: []
---

# Phase 1: Codify naming rule and supersede prior plans

## Overview

Document Convention B in `README.md` and close out the two prior decision plans (parked `260528-1110-…` and completed `260528-1221-…`). One docs commit.

## Requirements

- Functional: rule is discoverable in README; prior plans no longer show as "pending" in `ck plan status`.
- Non-functional: plain ASCII; under ~25 lines added to README; no breaking link changes.

## Architecture

Single concern: decision record. Touches 3 markdown files. No code.

## Related Code Files

- Modify: `D:\tiennm99\docker-images\README.md` (append new "## Naming convention" section after the image-index table at L7-10)
- Modify: `D:\tiennm99\docker-images\plans\260528-1110-image-naming-convention-decision\plan.md` (frontmatter `status: superseded` + supersede pointer)
- Modify: `D:\tiennm99\docker-images\plans\260528-1221-scribe-oci-labels-expansion\plan.md` (frontmatter `status: completed` + commit refs)

## Implementation Steps

1. Append the codified rule to `README.md` immediately after the image-index table (before any per-image sections). Use this exact block (plain ASCII, no smart quotes, no em-dashes):

   ```markdown
   ## Naming convention

   This repo uses Docker-Hub-style image names: `<service>:<full-version>`.

   - Subdirectory: `<service>-<major-line>/` (e.g. `couchbase-2.5/`, `scribe-2.2/`). Subdir holds the build context for one major version line.
   - GHCR + Docker Hub image: `<service>` (e.g. `couchbase`, `scribe`). No major-line suffix in the image name.
   - Workflow file: `publish-<service>-<major-line>.yml` (matches subdir for path-trigger clarity).
   - Tags published per workflow:
     - `<full-version>` (e.g. `2.5.2`, `2.2`)
     - `<major-line>` (e.g. `2.5`) - optional alias
     - `latest` - canonical build of the only major line this repo carries for the service.

   If a second major line of the same service is ever added to this repo, drop `latest` from both workflows and publish `<major-line>-latest` tags instead.
   ```

2. Update `plans/260528-1110-image-naming-convention-decision/plan.md`:
   - Set frontmatter `status: superseded` (or add the field if missing).
   - Add a top-of-body note: `> **Superseded by [260528-1324-rename-couchbase-to-convention-b](../260528-1324-rename-couchbase-to-convention-b/plan.md) on 2026-05-28.** User chose Convention B; see superseding plan for rationale.`

3. Update `plans/260528-1221-scribe-oci-labels-expansion/plan.md`:
   - Set frontmatter `status: completed`.
   - Add a top-of-body note: `> **Completed 2026-05-28 via commits `afaa631` (ASCII cleanup) and `4ef74f3` (OCI labels). Pattern reused as template for Phase 03 of [260528-1324-rename-couchbase-to-convention-b](../260528-1324-rename-couchbase-to-convention-b/plan.md).**`

4. Commit: `docs: codify Convention B naming rule and close prior decision plans`

## Success Criteria

- [ ] README contains "Naming convention" section with the 5 bullet structure above
- [ ] Parked 1110 plan marked superseded with link forward
- [ ] Completed 1221 plan marked completed with commit refs
- [ ] One commit, message starts with `docs:`
- [ ] `git diff --stat` for this phase touches exactly 3 files

## Risk Assessment

- **Risk:** future contributor adds a 2nd major-line service (e.g. `couchbase-6`) and forgets the `latest` rule, causing `couchbase:latest` to flip to the newer version mid-flight. **Mitigation:** the codified rule explicitly calls this out; future cross-plan scan would catch it.
- **Risk:** prior plans' `status: superseded` field may not be a recognized `ck plan` status. **Mitigation:** if `ck plan status` rejects, fall back to a body-level note plus removing the file from active scan via different means (e.g. archive flag). Verify after step 2.

## Open question (not blocking)

- Should the rule also live in a project-local `CLAUDE.md` so future Claude Code sessions auto-load it? Defaulted to README-only (KISS — repo has no CLAUDE.md today). Trivial to add later if it becomes painful.
