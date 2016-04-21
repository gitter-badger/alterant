#!/bin/sh

set -ex

INSTALL_PATH="../../../../install"
export BUILD="$PWD/vendor/libgit2/build"
export PKG_CONFIG_PATH="$PKG_CONFIG_PATH:$INSTALL_PATH/lib/pkgconfig:$INSTALL_PATH/lib64/pkgconfig"

FLAGS=$(pkg-config --static --libs --cflags libgit2 libssh2  openssl libcrypto) || exit 1

export CGO_LDFLAGS="-L$BUILD ${FLAGS}"
export CGO_CFLAGS="-I./$INSTALL_PATH/include"

$@
