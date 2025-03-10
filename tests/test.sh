#!/bin/bash

function test() {
  test_dir=$1

  go run main.go generate -i $test_dir/sonic-config.yaml -o $test_dir/config_db.json -p $test_dir/platform.json

  if [[ $(diff $test_dir/expected.json $test_dir/config_db.json) ]]; then
    echo TEST in $test_dir FAILED
    diff --color=always $expected $output
    rm $output
    exit 1
  fi

  rm $test_dir/config_db.json
}

test $(pwd)/tests/1
test $(pwd)/tests/2
test $(pwd)/tests/3
