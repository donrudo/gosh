#!/usr/bin/env bash

pushd plugins > /dev/null || exit

for f in *; {
  go build -buildmode=plugin  -o "../${DIR_PKG_PLUGIN}/${f}.cmd.so" "${f}"
  echo "BUILD: $f.cmd.so";
}

popd > /dev/null || exit

