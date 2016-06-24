#!/bin/sh
set -x
set -e

# Set temp environment vars
export GOPATH=/tmp/go
export PATH=${PATH}:${GOPATH}/bin
export BUILDPATH=${GOPATH}/src/github.com/Cimpress-MCP/go-git2consul
export PKG_CONFIG_PATH="/usr/lib/pkgconfig/:/usr/local/lib/pkgconfig/"

# Build git2consul
cd ${BUILDPATH}
go get -v
GOOS=linux GOARCH=amd64 go build -o /build/bin/git2consul.linux.amd64 .
# GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o /build/bin/git2consul.darwin.amd64 .
