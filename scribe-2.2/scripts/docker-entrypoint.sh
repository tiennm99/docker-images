#!/bin/bash
# docker-entrypoint.sh – Scribe daemon startup script
#
# Reads SCRIBE_CONFIG env var (default: /etc/scribe/scribe.conf)
# and launches scribed with that configuration file.
#
# Usage inside container:
#   docker run ... scribed           (uses default config)
#   SCRIBE_CONFIG=/my.conf docker run ...
set -euo pipefail

SCRIBE_CONFIG="${SCRIBE_CONFIG:-/etc/scribe/scribe.conf}"

if [ ! -f "$SCRIBE_CONFIG" ]; then
    echo "[scribe] ERROR: config file not found: $SCRIBE_CONFIG" >&2
    exit 1
fi

echo "[scribe] Starting scribed with config: $SCRIBE_CONFIG"
exec scribed "$SCRIBE_CONFIG"
