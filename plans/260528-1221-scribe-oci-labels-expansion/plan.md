---
title: "Expand OCI image-spec labels on scribe-2.2"
description: "Add title/licenses labels and override workflow defaults so published image has correct upstream metadata."
status: completed
priority: P3
effort: 20m
branch: main
tags: [docker, oci, metadata, scribe]
created: 2026-05-28
---

# Scribe-2.2 OCI Labels Expansion

> **Completed 2026-05-28 via commits `afaa631` (ASCII cleanup) and `4ef74f3` (OCI labels). Pattern reused as template for Phase 03 of [260528-1324-rename-couchbase-to-convention-b](../260528-1324-rename-couchbase-to-convention-b/plan.md).**

Single-pass metadata change. Adds `image.title` and `image.licenses` to the Dockerfile, and overrides the GitHub Actions `docker/metadata-action` defaults so the published image carries correct Scribe-upstream values rather than this repo's defaults.

## Why this is not a one-line change

`docker/metadata-action` auto-injects 8 OCI labels with values derived from THIS repo (`tiennm99/docker-images`). CLI `--label` overrides Dockerfile `LABEL` instructions for the same key, so the Dockerfile labels are effectively shadowed at workflow build time. Fix: override the relevant keys in the workflow's `metadata-action` `labels:` block. Keep Dockerfile labels as the manual-build fallback.

## Decisions locked

- Title: `Facebook Legacy Scribe` (consistent with existing description prefix).
- Licenses: `Apache-2.0` (verified, single license, no SPDX expression needed).
- `image.created`: let `metadata-action` auto-inject (build timestamp). No `--build-arg` plumbing.
- Skip: `url`, `documentation`, `vendor`, `revision`, `authors`.
- Override approach: BOTH Dockerfile (manual builds) AND workflow `labels:` block (CI builds). Same key/value in both places — single source of truth in plan, dual placement for two build paths.

## Phases

- [ ] Phase 01 — Dockerfile + workflow labels ([phase-01-dockerfile-and-workflow-labels.md](phase-01-dockerfile-and-workflow-labels.md))

## Files touched

- `scribe-2.2/Dockerfile` (runtime stage, lines 70-73)
- `.github/workflows/publish-scribe-2.2.yml` (metadata-action step, lines 38-47)

## Out of scope

- Existing uncommitted ASCII cleanup across `scribe-2.2/**` — separate logical commit, untouched by this plan.
- Other images in the repo (only `scribe-2.2/` touched).
