#!/bin/bash

for num in $(seq 1 1399); do
  pushd testdata || exit
  url="http://www.gutenberg.org/files/${num}/${num}-0.txt"
  wget "$url"
  popd || exit
done
