# Image Naming Convention -- Final Verdict (Round 3)

State verified 2026-05-28:

- 2 images only: couchbase-2.5/, scribe-2.2/ (commit ff5e5f6 dropped gradle/oraclejdk)
- scribe carries OCI labels including image.title, image.licenses=Apache-2.0, image.source (commits afaa631 + 4ef7dfb). Couchbase does NOT.
- Workflow labels-block pattern proven in publish-scribe-2.2.yml L48-53.

## Does new state alter the verdict?

Two real changes:

1. Sample collapsed to 2. Defend-the-convention frame is dead. 50/50 split = clean pick, no incumbent.
2. OCI image.title on scribe. Weakens ONE sub-argument (image-name-as-identifier) but does not touch the load-bearing one: latest semantics across major lines. docker inspect showing the title does not help a user typing docker pull -- they see the name, not labels, until after pull.

Neither change reaches the Round-2 core argument: archived-software repo = museum, latest per major line is the only stable semantic.

## Three options

### Option A -- rename scribe to scribe-2.2

- Edits: workflow L43,44,67 (3); scribe-2.2/README.md (4); scribe-2.2/example/docker-compose.yml (1); root README.md L10 (1). ~9 mechanical edits.
- latest: stable per line. scribe-2.2:latest = canonical 2.2 build, always.
- New couchbase-6 later: new subdir, new workflow, new package. Trivial. No collision.
- Hypothetical scribe-3: new subdir scribe-3.x/, new package. Old scribe-2.2 untouched, no latest drift.
- Discoverability: GHCR list shows couchbase-2.5, scribe-2.2 side-by-side, alphabetical, lineage explicit.
- OCI label: complementary, not redundant. Name says which line; label says what it is.

### Option B -- rename couchbase-2.5 to couchbase

- Edits: workflow L43,44,61 (3); root README.md L9 link + L15 pull (2); couchbase-2.5/docker-compose.example.yml (1). ~6 edits + 1 orphan package.
- latest: BREAKS the moment couchbase-6 lands. couchbase:latest = newest-published = 6.x. User searching for 2.5 gets 6 silently. Must drop latest or pin awkwardly.
- New couchbase-6 later: same package, new tag 6.x.x. Cheap, but forces latest-policy reckoning. Two workflows pushing to one package -- tag uniqueness saves data but latest flaps.
- Hypothetical scribe-3: same story. Repeated problem.
- Discoverability: GHCR shows couchbase + scribe. Major lines hidden behind tag list. Docker-official muscle memory yes, but this is not sliding-window software.
- OCI label: helps identification post-pull, does NOT solve latest ambiguity pre-pull.

### Option C -- codified asymmetry (lazy migration)

- Rule: 1 line per service = service:version; more-than-1 line = rename to service-major:full lazily.
- Under C, today state = couchbase + scribe (B-shaped). Lazy migration = the painful B-to-A rename at the worst moment (adding new line under deadline).
- latest: works today, breaks at migration.
- Forces the most expensive migration (rename + orphan + consumer breakage) at the moment of highest cognitive load. Anti-KISS.

## Verdict: Option A (unchanged from prior rounds)

1. Decisive argument (latest stability across major lines for archived software) unchanged by both new state-changes. Sample size and OCI labels do not touch it.
2. OCI image.title is COMPLEMENTARY to the line-suffix name, not a substitute. Pull-time discoverability (docker search, GHCR URL, badge) happens before any inspect.
3. Migration cost asymmetry stays: A = ~9 edits, no latest-policy rewrite. B = ~6 edits + new latest policy + landmine for next major line + orphan. C = pay later, more, at worse time.
4. With only 2 images, A locks in the rule cheaply NOW. Cost of B later (when couchbase-6 arrives) is provably higher than A now. Asymmetric optionality favors A.
5. couchbase-2.5-is-the-lone-convention-follower framing reversible: scribe rename has cheaper blast radius (7-week-old workflow, 1 example-compose consumer, fresh package). Sunk-cost cuts against B, not A.

## Codified rule (README + CLAUDE.md, copy-pasteable)

```
## Naming convention

This repo packages archived/legacy software. Each major version line is a
write-once artifact with its own base OS, Dockerfile, attestation history,
and latest pointer.

Naming rule (all four MUST match):

    Subdirectory      :  <service>-<major-line>/
    GHCR image        :  <service>-<major-line>
    Docker Hub image  :  <service>-<major-line>
    Workflow file     :  publish-<service>-<major-line>.yml

Examples: couchbase-2.5, scribe-2.2. Future: couchbase-6, redis-3.

Tags published per image:

    - <full-version>      e.g. 2.5.2, 2.2, 8u201
    - latest              alias for canonical build of THIS major line

Adding a new major line of an existing service = new subdir + new workflow
+ new package. NOT a new tag under the existing package. Keeps latest
semantics stable per line and attestation lineage clean.

OCI image-spec labels (org.opencontainers.image.*) REQUIRED on every image
for metadata-level identification independent of name. See
scribe-2.2/Dockerfile for the canonical set.
```

## Migration order if A approved

1. Rename literals in scribe (9 edits above).
2. Push, verify GHCR package scribe-2.2 populates with latest + 2.2.
3. Delete orphan ghcr.io/tiennm99/scribe + Docker Hub tiennm99/scribe via UI.
4. Add naming-convention section to root README.md + project-local CLAUDE.md.
5. Separate task: port OCI labels to couchbase using the labels: pattern from publish-scribe-2.2.yml L48-53.

## Unresolved questions

1. Orphan ghcr.io/tiennm99/scribe after rename -- delete or tombstone? Lean delete (7 weeks old).
2. Pre-check Docker Hub docker.io/tiennm99/scribe-2.2 namespace availability before rename push? Trivial check, do it.
3. CLAUDE.md scope: user-global vs project-local. Convention belongs in project-local D:/tiennm99/docker-images/CLAUDE.md (create if absent). Confirm.
4. OCI labels on couchbase -- same task as rename, or separate? Lean separate to keep rename diff minimal.
