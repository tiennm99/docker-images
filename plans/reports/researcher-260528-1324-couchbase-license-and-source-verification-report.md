# Couchbase 2.5 License & Source Verification

## Q1: License (SPDX Identifier)

**Finding:** `Apache-2.0`

**Evidence:** Wikipedia: "Couchbase Server is a packaged version of Couchbase's open source software technology and is available in a community edition without recent bug fixes with an Apache 2.0 license" ([Couchbase Server - Wikipedia](https://en.wikipedia.org/wiki/Couchbase_Server)). Couchbase Server 2.0 released Dec 2012; v2.5.2 inherits same Apache 2.0 license applied since open-source inception (2011, post Membase/CouchOne merger).

**SPDX ID:** `Apache-2.0`

No mixing of licenses in Server 2.5.x; only one license applies to the binary.

---

## Q2: `image.source` URL

**Recommendation:** `https://github.com/couchbase/docker`

**One-line reason:** General upstream repo URL maximizes durability; pinned commit hashes are infrastructure detail, not user-facing metadata.

**Alternatives rejected:**
- Pinned branch URLs (`…/tree/14f9e5f37…`) — commit references rot when branches are deleted or rebased; violates OCI-label durability principle (user should find repo, not fight dead links).
- Scribe precedent uses `https://github.com/facebookarchive/scribe` (upstream source root), not a pinned archive commit, so Couchbase should match that pattern.

---

## Unresolved Questions

None. Both values verified via public sources.
