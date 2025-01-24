#!/bin/bash

for i in 1 2
do
  go run main.go generate -i tests/input_$i.yaml -o tests/result.json

  if [[ $(diff tests/expected_$i.json tests/result.json) ]]; then
    echo TEST $i FAILED
    diff --color=always tests/expected_$i.json tests/result.json
    rm tests/result.json
    exit 1
  fi

  rm tests/result.json
done
