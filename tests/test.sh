#!/bin/bash

for i in 1 2 3 4
do
  test_dir=$(pwd)/tests/$i
  docker run --rm -v $(pwd)/tests/device:/usr/share/sonic/device:ro -v $test_dir:/sonic sonic-configdb-utils:local generate -c /sonic/current-config_db.json -i /sonic/sonic-config.yaml -o /sonic/config_db.json
  diff --color=always $test_dir/expected.json $test_dir/config_db.json

  if [[ $? != 0 ]]; then
    echo TEST $i FAILED
    rm -f $(pwd)/tests/$i/config_db.json
    exit 1
  fi

  echo TEST $i PASSED
  rm -f $(pwd)/tests/$i/config_db.json
done
