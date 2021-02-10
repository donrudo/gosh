#!/usr/bin/env bash

pushd plugins > /dev/null

echo ""

for f in *; {
    go build -buildmode=plugin  -o ../pkg/plugins/$f.cmd.so $f
    echo "BUILD: $f.cmd.so";}

echo ""
popd > /dev/null

