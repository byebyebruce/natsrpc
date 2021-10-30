#!/bin/bash

set -ex

CURDIR=$(cd $(dirname $0); pwd)
PROTO_IMPORT=$CURDIR/../..

protoc \
-I $PROTO_IMPORT \
--go_opt=paths=source_relative \
--proto_path=$CURDIR/ \
--go_out=plugins=grpc:. \
--natsrpc_out=. \
greet.proto


protoc \
-I $PROTO_IMPORT \
--proto_path=$CURDIR/ \
--go_out=plugins=grpc:$CURDIR/ \
--go_opt=paths=source_relative \
--natsrpc_out=. \
$CURDIR/*.proto