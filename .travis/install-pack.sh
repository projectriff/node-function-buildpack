#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

wget -qO- https://github.com/buildpack/pack/releases/download/v0.4.1/pack-v0.4.1-linux.tgz | tar xvz -C $HOME/bin
export PATH="$HOME/bin:$PATH"
