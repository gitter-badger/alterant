#!/bin/sh

set -ex

INSTALL_PATH="$PWD/install"
SUBMODULE_PATH="$PWD/submodules/openssl"

cd $SUBMODULE_PATH
mkdir -p $INSTALL_PATH/lib

FLAGS="--prefix=$INSTALL_PATH no-shared -fPIC"
if [ $(uname) == "Darwin" ]; then
    ./Configure darwin64-x86_64-cc $FLAGS
elif [ $(uname) == "Linux" ]; then
    ./config threads $FLAGS -DOPENSSL_PIC
fi

make depend
make
make install_sw
