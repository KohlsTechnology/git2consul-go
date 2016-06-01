#!/bin/sh
set -x
set -e

# Set temp environment vars
export GOPATH=/tmp/go
export PATH=${PATH}:${GOPATH}/bin
export BUILDPATH=${GOPATH}/src/github.com/Cimpress-MCP/go-git2consul
export PKG_CONFIG_PATH="/usr/lib/pkgconfig/:/usr/local/lib/pkgconfig/"

# Install build deps
apk --no-cache --no-progress --virtual build-deps add go gcc musl-dev make cmake openssl-dev libssh2-dev

# Install libgit2
/build/install-libgit2.sh

# Set up go environment
mkdir -p $(dirname ${BUILDPATH})
ln -s /app ${BUILDPATH}
