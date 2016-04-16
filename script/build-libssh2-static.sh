#!/bin/sh

set -ex

INSTALL_PATH="$PWD/install"
SUBMODULE_PATH="$PWD/submodules/libssh2"

mkdir -p $INSTALL_PATH/lib &&
cd $SUBMODULE_PATH &&
./buildconf
./configure --prefix=$INSTALL_PATH --disable-shared --with-openssl CFLAGS="-fPIC" LDFLAGS="-m64 -L$INSTALL_PATH/lib -L$INSTALL_PATH/lib64" LIBS=-ldl
make
make install
