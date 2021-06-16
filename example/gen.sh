#!/bin/bash

set -ex

cur=$(cd $(dirname $0); pwd)
root="$cur/.."

go run "$root/cmd" \
-ip=github.com/byebyebruce/natsrpc \
-s="$root/example/api/example.go" \
-d="$root/example/api/service/example.go" \
-op=service