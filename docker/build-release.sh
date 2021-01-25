#!/bin/sh

go_version=$(go version)

echo "Go Version: ${go_version}"
echo "Creating build..."
make gen_release
echo "Copying build artifacts from ${PWD}"
mkdir -p /build
cp -v build/factorio-server-manager-linux.zip build/factorio-server-manager-windows.zip /build/
