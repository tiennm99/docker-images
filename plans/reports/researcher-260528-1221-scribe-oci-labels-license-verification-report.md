# Scribe OCI Image Labels — License & Metadata Verification

## Q1: License Identifier & Multi-License Analysis

**Finding: SPDX identifier is `Apache-2.0` (single license).**

### Evidence
- **Scribe upstream:** https://github.com/facebookarchive/scribe — Licensed under Apache License 2.0 (verified via LICENSE file).
- **SPDX value:** `Apache-2.0`
- **Quote:** "Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License."

### Multi-License Check (Thrift Dependency)
Dockerfile explicitly builds Apache Thrift 0.9.1 (line 28). Thrift is also Apache-2.0 licensed.

**Conclusion:** No additional license expression needed. Both Scribe and bundled Thrift are Apache-2.0. SPDX label value:
```
LABEL org.opencontainers.image.licenses="Apache-2.0"
```

Not required to use SPDX expression (`Apache-2.0 AND Apache-2.0`) — single `Apache-2.0` is correct.

---

## Q2: Additional OCI Image Annotations — Recommendation Table

| Label | Recommend? | Reason |
|-------|------------|--------|
| `image.title` | **YES** | Short human-readable name; trivial to add. E.g., `"Scribe Log Aggregation Daemon"` |
| `image.url` | NO | Project is archived (read-only); homepage is same as `image.source`. Redundant. |
| `image.documentation` | NO | No maintained docs site. README is in archived repo — already discoverable via `image.source`. Skip. |
| `image.vendor` | NO | Original author (Facebook) no longer maintains. Current maintainer is builder (you). Not load-bearing for legacy image. Skip. |
| `image.authors` | **MAYBE** | If you want attribution, add `"Tien Nguyen Minh"` (matches Dockerfile maintainer). Otherwise skip — not critical for archived image. |
| `image.created` | **YES** | Build timestamp is meaningful for reproducibility. Add via `--build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')` at build time. |
| `image.revision` | NO | Dockerfile uses `git clone --depth 1` (line 58), losing commit history. Revision would be stale/misleading. Skip unless build args fix the clone depth. |

---

## Recommended Additions (Final Label Set)

```dockerfile
# Add these two lines to the LABEL block (lines 71-73)
LABEL org.opencontainers.image.title="Scribe Log Aggregation Daemon"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.created="${BUILD_DATE}"
```

Optional (low priority):
```dockerfile
LABEL org.opencontainers.image.authors="Tien Nguyen Minh <tiennm99@outlook.com>"
```

Build command to inject BUILD_DATE:
```bash
docker build --build-arg BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')" -t scribe:2.2 .
```

---

## Unresolved Questions

None. License is clear (Apache-2.0), dependency check complete, and OCI spec is authoritative.
