#! /usr/bin/bash

set -xe


WASM_EXE=2048.wasm

env GOOS=js GOARCH=wasm go build -o $WASM_EXE .
mv $WASM_EXE wasm/
cd ./wasm

(sleep 1 && kde-open http://localhost:8000/index.html) & python -m http.server
