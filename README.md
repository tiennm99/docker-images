# Docker Images

Custom Docker images of legacy/discontinued software — built and published to GHCR.

## Image index

| Image | Tag | GHCR package | Notes |
|-------|-----|-------------|-------|
| Couchbase 2.5 | `2.5.2` | [couchbase](https://github.com/tiennm99/docker-images/pkgs/container/couchbase) | Legacy NoSQL; no official image exists |
| Gradle 8 | latest | — | JDK 8 + Corretto base |
| Oracle JDK 8 | `8u201` | [oraclejdk](https://github.com/tiennm99/docker-images/pkgs/container/oraclejdk) | Archived Oracle JDK build |
| Scribe 2.2 | `2.2` | [scribe](https://github.com/tiennm99/docker-images/pkgs/container/scribe) | Facebook's legacy log aggregation daemon |

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

## Couchbase-2.5

```bash
docker pull ghcr.io/tiennm99/couchbase:2.5.2
```

More options: https://github.com/tiennm99/docker-images/pkgs/container/couchbase

Docker Compose example: [couchbase-2.5/docker-compose.example.yml](couchbase-2.5/docker-compose.example.yml)

### Credits

- https://docs.couchbase.com/couchbase-manual-2.5/cb-install/
- https://github.com/couchbase/docker/tree/14f9e5f37723ef15affa8cf3227b15409805f668/enterprise/couchbase-server/2.5.2
- https://github.com/couchbase/docker/tree/a9e69006c4dae3d6b8aa4055436c83e9cd51d9e4/enterprise/couchbase-server/2.5.2

## Gradle 8

### Credits

- https://gradle.org/releases/#8.13
- https://github.com/gradle/docker-gradle/blob/8/jdk8-corretto/Dockerfile

## Oracle JDK 8

### Credits

- https://www.oracle.com/java/technologies/javase/javase8-archive-downloads.html
- https://github.com/oracle/docker-images/blob/main/OracleJava/8/jdk/Dockerfile.ol8
- https://repo.huaweicloud.com/java/jdk/8u201-b09/

## Scribe 2.2

Facebook's legacy real-time log aggregation daemon.

See [scribe-2.2/README.md](scribe-2.2/README.md) for usage, configuration, and examples.

### Credits

- https://github.com/facebookarchive/scribe
- https://archive.apache.org/dist/thrift/0.9.1/

## GitHub Actions template

### Credits

- https://docs.github.com/en/actions/tutorials/publish-packages/publish-docker-images
