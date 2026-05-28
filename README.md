# Docker Images

Custom Docker images of legacy/discontinued software — built and published to GHCR.

## Image index

| Image | Tag | GHCR package | Notes |
|-------|-----|-------------|-------|
| Couchbase 2.5 | `2.5.2` | [couchbase-2.5](https://github.com/tiennm99/docker-images/pkgs/container/couchbase-2.5) | Legacy NoSQL; no official image exists |
| Scribe 2.2 | `2.2` | [scribe](https://github.com/tiennm99/docker-images/pkgs/container/scribe) | Facebook's legacy log aggregation daemon |

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

## Scribe 2.2

Facebook's legacy real-time log aggregation daemon.

See [scribe-2.2/README.md](scribe-2.2/README.md) for usage, configuration, and examples.

### Credits

- https://github.com/facebookarchive/scribe
- https://archive.apache.org/dist/thrift/0.9.1/

## GitHub Actions template

### Credits

- https://docs.github.com/en/actions/tutorials/publish-packages/publish-docker-images
