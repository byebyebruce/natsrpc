#!/bin/bash

set -ex

CURDIR=$(cd $(dirname $0); pwd)
PROTO_IMPORT=$CURDIR/../../..


#--go_out=plugins=grpc:$CURDIR/ \
#--go_opt=paths=source_relative \

protoc \
--proto_path=$PROTO_IMPORT \
--proto_path=$CURDIR/../../pb \
--proto_path=$CURDIR \
--go_out=paths=source_relative:$CURDIR \
--natsrpc_out=. \
$CURDIR/request.proto