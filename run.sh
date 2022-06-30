#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
./server -port=8011 &
./server -port=8012 &
./server -port=8013 -api=1 &

wait