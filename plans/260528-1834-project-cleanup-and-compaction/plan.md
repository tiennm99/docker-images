---
title: "Project cleanup and compaction"
description: "Apply read-only review findings from code-reviewer + code-simplifier agents (2026-05-28 17:21 session). Mix of safety, hygiene, polish, and one structural DRY refactor. Cook in a future session after open questions are answered."
status: pending
priority: P2
branch: "main"
tags: [cleanup, refactor, dockerfile, workflows, dry, kiss]
blockedBy: []
blocks: []
created: "2026-05-28T11:34:48.715Z"
createdBy: "ck:plan"
source: skill
---

# Project cleanup and compaction

## Overview

Two ClaudeKit agents (`code-reviewer` + `code-simplifier`) scanned the whole project after the Convention B + OCI labels rollout completed. Findings consolidate into 4 phases ordered by risk and ROI. No new features — pure cleanup, simplification, consistency, and one structural DRY refactor.

Underlying agent reports are committed under `plans/reports/` for self-contained context.

## Phases

| Phase | Name | Status |
|-------|------|--------|
| 1 | [High-severity correctness fixes](./phase-01-high-severity-correctness-fixes.md) | Pending |
| 2 | [Dockerfile and workflow hygiene](./phase-02-dockerfile-and-workflow-hygiene.md) | Pending |
| 3 | [Low-priority polish and consistency](./phase-03-low-priority-polish-and-consistency.md) | Pending |
| 4 | [Workflow_call DRY refactor](./phase-04-workflow-call-dry-refactor.md) | Superseded: one `publish.yml` reads each folder's `image.json` (2026-11-14) |

## Phase ordering rationale

- **Phase 1 first** — correctness bugs (namespace leak via `github.actor`; silently swallowed curl exit code) and the one factual README error. Tiny diff, immediate safety win.
- **Phase 2 next** — Dockerfile and workflow hygiene. Independent file edits, no structural changes. Decisions on maintainer canonical + template.yml fate must be locked before this phase starts.
- **Phase 3** — pure polish (LABEL compaction, workflow name uniqueness, ENV→ARG, .gitignore). Can ship anytime, no dependencies.
- **Phase 4 last** — the structural workflow_call refactor depends on Phase 1's `github.actor` fix already being baked in (otherwise the refactor would carry the bug forward). Also depends on the maintainer + reuse-style decisions.

## Open questions (MUST resolve before Phase 2 cook)

1. **Maintainer canonical:** `miti99 <tiennm99@outlook.com>` (matches GHCR namespace) or `Tien Nguyen Minh <tiennm99@outlook.com>` (human-readable, matches couchbase + scribe)?
2. **`.github/workflows/template.yml` fate:** update to a real-workflow scaffold OR delete + remove the "GitHub Actions template" section from README?
3. **Workflow reuse style (Phase 4):** repo-private `workflow_call` at `.github/workflows/_publish-image.yml`, OR composite action at `.github/actions/publish-image/`?

## Source reports

- `plans/reports/code-reviewer-260528-1721-project-cleanup-review-report.md` — full code-reviewer findings (8 categories, severity-tagged)
- `plans/reports/code-simplifier-260528-1721-logic-compaction-recommendations-report.md` — code-simplifier findings + a concrete `workflow_call` sketch

## Commit strategy

One commit per phase, with the option for Phase 1 to split into two commits if `fix(workflows):` + `fix(oraclejdk):` semantics matter more than commit-count economy. Cook decides.

## Out of scope (explicitly NOT cleanup)

- Adding renovate/dependabot
- Adding image scanning (trivy/grype)
- Adding multi-arch builds (arm64)
- Adding smoke tests for published images
- Adding cosign signing beyond existing attestation
- Any further image renames (Convention B is locked in)

## Post-execution manual actions

- After Phase 1 pushes: re-trigger any `workflow_dispatch` that used a non-`tiennm99` actor and watch that publishes now use `repository_owner`. (No orphan packages expected since this fix only matters when someone other than the owner manually triggers.)
- After Phase 4 pushes: verify each per-image workflow still triggers on the same `paths:` and publishes identical `tags:` + `labels:`. Compare `docker inspect` output before and after.
