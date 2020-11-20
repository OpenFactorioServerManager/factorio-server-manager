#!/bin/bash
set -eou pipefail
(
  cd ..
  make build
  cp build/factorio-server-manager-linux.zip docker/factorio-server-manager-linux.zip
)
docker build -f Dockerfile-local -t factorio-server-manager:dev .

rm factorio-server-manager-linux.zip
