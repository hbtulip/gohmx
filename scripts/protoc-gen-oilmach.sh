#!/usr/bin/env bash

protoc --proto_path=../api/ --go_out=plugins=grpc:../api/ transferMessage.proto

