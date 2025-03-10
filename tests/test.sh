#!/bin/bash

function test() {
  config_dir=$1

  go run main.go generate -i $config_dir/sonic-config.yaml --sonic-config-dir $config_dir --device-dir $config_dir/../device

  if [[ $(diff $config_dir/expected.json $config_dir/config_db.json) ]]; then
    echo TEST in $config_dir FAILED
    diff --color=always $expected $output
    rm $output
    exit 1
  fi

  rm $config_dir/config_db.json
}

test $(pwd)/tests/1
test $(pwd)/tests/2
test $(pwd)/tests/3
