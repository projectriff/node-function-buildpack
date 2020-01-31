#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

wget -qO- https://github.com/buildpack/pack/releases/download/v0.8.1/pack-v0.8.1-linux.tgz | tar xvz -C .
sudo mv pack /usr/local/bin/
