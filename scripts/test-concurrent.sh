#!/bin/bash

trap 'exit' INT

for i in $(seq 1 1400); do
  if [[ -f "testdata/${i}.txt" ]]; then
    curl -X POST --data-binary "@testdata/${i}.txt" localhost:8080/learn -H 'Content-Type: text/plain' &
  fi
  curl localhost:8080/generate &
done
