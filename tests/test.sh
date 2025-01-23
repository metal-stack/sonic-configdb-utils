#!/bin/bash

go run main.go generate -i tests/input.yaml -o tests/result.json

if [[ $(diff tests/result.json tests/expected.json) ]]; then
  diff --color=always tests/result.json tests/expected.json
  rm tests/result.json
  exit 1
fi

rm tests/result.json
