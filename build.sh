#!/bin/bash

if [ ! -d "./bin" ]; then
    mkdir ./bin
fi

dep ensure

echo "build linux/amd64"
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -v -o ./bin/ssl-verify.linux.amd64

echo "build windows/amd64"
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -v -o ./bin/ssl-verify.exe

echo "build darwin/amd64"
GOOS=darwin GOARCH=amd64 go build -v -o ./bin/ssl-verify.darwin.amd64
