#!/usr/bin/env bash

TS_VERSION=$(head -1 version.txt)
TS_COMMIT=$(git rev-parse HEAD)

# build for MacOS AMD64
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$TS_VERSION -X main.commit=$TS_COMMIT"  -o terrasolver_mac_amd64

# build for MacOS ARM64
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$TS_VERSION -X main.commit=$TS_COMMIT"  -o terrasolver_mac_arm64

# build for Linux AMD64
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$TS_VERSION -X main.commit=$TS_COMMIT"  -o terrasolver_linux_amd64
