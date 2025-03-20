#!/bin/bash

function test() {
  test_dir=$1
  cp $test_dir/current-config_db.json $test_dir/config_db.json

  docker run --rm -v $(pwd)/tests/device:/usr/share/sonic/device:ro -v $test_dir:/etc/sonic -v $test_dir:/sonic sonic-configdb-utils:local generate

  if [ $? -eq 1 ]; then
    rm -f $test_dir/config_db.json
    exit 1
  fi

  if [[ $(diff $test_dir/expected.json $test_dir/config_db.json) ]]; then
    echo TEST in $test_dir FAILED
    diff --color=always $test_dir/expected.json $test_dir/config_db.json
    rm $test_dir/config_db.json
    exit 1
  fi

  rm -f $test_dir/config_db.json
}

test $(pwd)/tests/1
test $(pwd)/tests/2
test $(pwd)/tests/3
