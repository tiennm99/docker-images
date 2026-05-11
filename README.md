# Docker Images

Custom Docker images of legacy/discontinued software — built and published to GHCR.

## Image index

| Image | Tag | GHCR package | Notes |
|-------|-----|-------------|-------|
| Couchbase 2.5 | `2.5.2` | [couchbase-2.5](https://github.com/tiennm99/docker-images/pkgs/container/couchbase-2.5) | Legacy NoSQL; no official image exists |
| Gradle 8 | latest | — | JDK 8 + Corretto base |
| Oracle JDK 8 | `8u201` | — | Archived Oracle JDK build |
| Scribe | — | — | Facebook's legacy log aggregation daemon |

## Couchbase-2.5

```bash
docker pull ghcr.io/tiennm99/couchbase-2.5:2.5.2
```

More options: https://github.com/tiennm99/docker-images/pkgs/container/couchbase-2.5

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

## Scribe

Facebook's legacy real-time log aggregation daemon.

See [scribe/README.md](scribe/README.md) for usage, configuration, and examples.

### Credits

- https://github.com/facebookarchive/scribe
- https://archive.apache.org/dist/thrift/0.9.1/

## GitHub Actions template

### Credits

- https://docs.github.com/en/actions/tutorials/publish-packages/publish-docker-images
