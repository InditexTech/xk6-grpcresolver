#!/bin/bash
rm -rf ./grpc
git clone https://github.com/grpc/grpc-go ./grpc/
cd ./grpc/examples && go run ./helloworld/greeter_server/main.go