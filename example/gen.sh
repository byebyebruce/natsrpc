#!/bin/bash

set -ex

cur=$(cd $(dirname $0); pwd)
root="$cur/.."

cd "$root"
go run "$root/cmd" \
-ip=github.com/byebyebruce/natsrpc \
-s="example/api/example.go" \
-d="example/api/service/example.go" \
-op=service