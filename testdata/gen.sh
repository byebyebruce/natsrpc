#!/bin/bash

set -ex

cur=""$(cd $(dirname $0); pwd)
root="$cur/.."

cd "$root"
go run "$root/cmd" \
-ip=github.com/byebyebruce/natsrpc \
-s="testdata/greeter.go" \
-d="testdata/autogen/greeter.go" \
-op=service