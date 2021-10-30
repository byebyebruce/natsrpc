#!/bin/bash

set -ex

CURDIR=$(cd $(dirname $0); pwd)
PROTO_IMPORT=$CURDIR/../..


#--go_out=plugins=grpc:$CURDIR/ \
#--go_opt=paths=source_relative \

protoc \
-I $PROTO_IMPORT \
--proto_path=$CURDIR/ \
--natsrpc_out=. \
$CURDIR/*.proto