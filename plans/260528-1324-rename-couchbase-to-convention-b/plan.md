---
title: "Apply Convention B naming: rename couchbase image and codify rule"
description: "Migrate this repo to Docker-Hub-style image naming (<service>:<version>). Rename ghcr.io/tiennm99/couchbase-2.5 -> couchbase. Port OCI labels to couchbase. Codify the rule in README. Close out prior naming-decision plans."
status: completed
priority: P2
branch: "main"
tags: [docker, naming, convention, couchbase, oci]
blockedBy: []
blocks: [260528-1110-image-naming-convention-decision]
created: "2026-05-28T06:30:23.203Z"
createdBy: "ck:plan"
source: skill
---

# Apply Convention B naming: rename couchbase image and codify rule

## Overview

Three brainstorm rounds had recommended Convention A (`<service>-<major>:<full-version>`). User overruled on YAGNI grounds — A's load-bearing argument (`latest` semantics breaking on multi-major-line coexistence) is a future hypothetical, not a current need. Repo has only 2 images, both single-major-line. Convention B (`<service>:<full-version>`) is the Docker-Hub standard, shorter, cheaper to migrate.

This plan implements Convention B: rename `ghcr.io/tiennm99/couchbase-2.5` to `couchbase` (and Docker Hub counterpart), keep `scribe` as-is (already conformant), port OCI labels to couchbase using the proven scribe pattern, codify the rule in `README.md`, and close out prior decision plans.

## Phases

| Phase | Name | Status |
|-------|------|--------|
| 1 | [Codify naming rule and supersede prior plans](./phase-01-codify-naming-rule-and-supersede-prior-plans.md) | Completed (`2f7bfde`) |
| 2 | [Rename couchbase image to <service>:<version>](./phase-02-rename-couchbase-image-to-service-version.md) | Completed (`3a0e31b`) |
| 3 | [Port OCI labels to couchbase image](./phase-03-port-oci-labels-to-couchbase-image.md) | Completed (`afa2930`) |

## Dependencies

- **Supersedes:** `plans/260528-1110-image-naming-convention-decision/` (parked "decide later" plan; this plan locks in Convention B and closes the loop).
- **Builds on:** `plans/260528-1221-scribe-oci-labels-expansion/` (completed via commits `afaa631` + `4ef74f3`). That plan's OCI label pattern is the template for Phase 03.

## Research input

- `plans/reports/researcher-260528-1324-couchbase-license-and-source-verification-report.md` — confirms Couchbase 2.5.2 license is `Apache-2.0`; recommends `image.source=https://github.com/couchbase/docker` (general repo, not pinned commit).

## Commit strategy

One commit per phase, three commits total. Pattern matches the just-shipped scribe work: `docs:`/`refactor:`/`feat:`.

## Post-execution manual actions (not part of cook)

1. Delete orphan packages via UI: `ghcr.io/tiennm99/couchbase-2.5` and `docker.io/${DOCKER_USERNAME}/couchbase-2.5`.
2. Validate: `docker pull ghcr.io/tiennm99/couchbase:2.5.2` and `docker inspect` shows OCI labels.

## Out of scope

- Subdirectory rename (`couchbase-2.5/` stays; subdir name is decoupled from image name under Convention B — codified in Phase 01's rule).
- Workflow filename rename (`publish-couchbase-2.5.yml` stays; path-trigger needs to match the subdir).
- New project-local `CLAUDE.md` — defaulted to README-only for KISS; flagged in Phase 01 as an optional add.
