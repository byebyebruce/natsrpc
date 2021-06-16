#!/bin/bash

set -ex

cur=""$(cd $(dirname $0); pwd)
root="$cur/.."

cd "$root"
go run "$root/cmd" \
-ip=github.com/byebyebruce/natsrpc \
-s="testdata/helloworld.go" \
-d="testdata/autogen/helloworld.go" \
-op=service