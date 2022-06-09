#!/usr/bin/env bash

TS_VERSION=$(head -1 version.txt)
TS_COMMIT=$(git rev-parse HEAD)

go run -ldflags "-X main.version=$TS_VERSION -X main.commit=$TS_COMMIT"  . "$@"
