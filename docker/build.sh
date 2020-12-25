#!/bin/sh

echo "Cloning ${FAC_BRANCH}"
git clone -b ${FAC_BRANCH} https://github.com/mroote/factorio-server-manager.git ${FAC_ROOT}
echo "Creating build..."
make gen_release
echo "Copying build artifacts..."
cp -v build/* /build/