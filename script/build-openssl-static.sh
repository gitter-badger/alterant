#!/bin/sh

set -ex

INSTALL_PATH="$PWD/install"
SUBMODULE_PATH="$PWD/submodules/openssl"

cd $SUBMODULE_PATH &&
mkdir -p $INSTALL_PATH/lib &&
mkdir -p build &&

# Switch to a stable version
git checkout OpenSSL_1_0_2-stable &&
./config threads no-shared --prefix=$INSTALL_PATH -fPIC -DOPENSSL_PIC &&
make depend &&
make &&
make install_sw
