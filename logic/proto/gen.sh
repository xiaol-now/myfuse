#!/usr/bin/bash

protoc --go_out=plugins=grpc:. -I. interface.proto