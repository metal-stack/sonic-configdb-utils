#!/bin/bash

function test() {
  test_dir=$1
  output=$test_dir/config_db.json

  docker run --rm --mac-address aa:aa:aa:aa:aa:aa -v $test_dir:/test_dir sonic-configdb-utils:local generate -e /test_dir/sonic-environment -o /test_dir/config_db.json -i /test_dir/sonic-config.yaml --device-dir /test_dir

  if [ $? -eq 1 ]; then
    echo TEST in $test_dir FAILED
    rm -f $output
    exit 1
  fi

  if [[ $(diff $test_dir/expected.json $output) ]]; then
    echo TEST in $test_dir FAILED
    diff --color=always $test_dir/expected.json $output
    rm -f $output
    exit 1
  fi

  rm -f $output
}

test $(pwd)/tests/1
test $(pwd)/tests/2
test $(pwd)/tests/3
