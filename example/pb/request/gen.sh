#!/bin/bash

set -ex

CURDIR=$(cd $(dirname $0); pwd)
PROTO_IMPORT=$CURDIR/../../..

protoc \
--proto_path=$PROTO_IMPORT \
--proto_path=$CURDIR/../../pb \
--proto_path=$CURDIR \
--natsrpc_out=paths=source_relative:$CURDIR \
$CURDIR/request.proto