#!/usr/bin/env bash
set -euo pipefail

mkdir -p /tmp/go-build-cache /tmp/go-mod-cache
export GOCACHE=/tmp/go-build-cache
export GOMODCACHE=/tmp/go-mod-cache

go test -v -count=1 -cover ./...

