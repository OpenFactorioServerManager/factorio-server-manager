#!/bin/sh

echo "Cloning ${FACTORIO_BRANCH}"
git clone -b ${FACTORIO_BRANCH} https://github.com/mroote/factorio-server-manager.git ${FACTORIO_ROOT}
echo "Creating build..."
make gen_release
echo "Copying build artifacts from ${PWD}"
mkdir -p /build
cp -v build/factorio-server-manager-linux.zip build/factorio-server-manager-windows.zip /build/
