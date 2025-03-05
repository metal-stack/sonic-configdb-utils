#!/bin/bash

function test() {
  input=$1
  platform_file=$2
  expected=$3
  output=tests/result.json

  go run main.go generate -i $input -o $output -f $platform_file

  if [[ $(diff $expected $output) ]]; then
    echo TEST for $input FAILED
    diff --color=always $expected $output
    rm $output
    exit 1
  fi

  rm $output
}

test tests/input_1_x86_64-accton_as7726_32x-r0.yaml device/accton/x86_64-accton_as7726_32x-r0/platform.json tests/expected_1_x86_64-accton_as7726_32x-r0.json
test tests/input_2_x86_64-accton_as7726_32x-r0.yaml device/accton/x86_64-accton_as7726_32x-r0/platform.json tests/expected_2_x86_64-accton_as7726_32x-r0.json
test tests/input_x86_64-accton_as4630_54te-r0.yaml device/accton/x86_64-accton_as4630_54te-r0/platform.json tests/expected_x86_64-accton_as4630_54te-r0.json
