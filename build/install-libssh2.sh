#!/bin/sh
set -x

# Set temp environment vars
export REPO=https://github.com/libssh2/libssh2
export BRANCH=libssh2-1.7.0
export REPO_PATH=/tmp/libssh2
export PKG_CONFIG_PATH="/usr/lib/pkgconfig/:/usr/local/lib/pkgconfig/"
export PKG_CONFIG_PATH="${PKG_CONFIG_PATH}:/tmp/libgit2/install/lib/pkgconfig:/tmp/openssl/install/lib/pkgconfig:/tmp/libssh2/build/src"

# Compile & Install libgit2 (v0.23)
git clone -b ${BRANCH} --depth 1 -- ${REPO} ${REPO_PATH}

mkdir -p ${REPO_PATH}/build
cd ${REPO_PATH}/build
cmake -DTHREADSAFE=ON \
      -DBUILD_CLAR=OFF \
      -DBUILD_SHARED_LIBS=OFF \
      -DCMAKE_C_FLAGS=-fPIC \
      -DCMAKE_BUILD_TYPE="RelWithDebInfo" \
      -DCMAKE_INSTALL_PREFIX=../install \
      ..
cmake --build . --target install

# Cleanup
# rm -r ${LIBGIT2PATH}
