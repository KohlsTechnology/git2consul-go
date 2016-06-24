#!/bin/sh
set -x
set -e

# Set temp environment vars
export LIBGIT2REPO=https://github.com/libgit2/libgit2.git
export LIBGIT2BRANCH=v0.24.0
export LIBGIT2PATH=/tmp/libgit2

# Compile & Install libgit2 (v0.23)
git clone -b ${LIBGIT2BRANCH} --depth 1 -- ${LIBGIT2REPO} ${LIBGIT2PATH}

mkdir -p ${LIBGIT2PATH}/build
cd ${LIBGIT2PATH}/build
cmake -DBUILD_CLAR=OFF ..
cmake --build . --target install

# Cleanup
rm -r ${LIBGIT2PATH}
