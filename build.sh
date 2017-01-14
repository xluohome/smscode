#!/bin/bash

echo "building..."

export GOPATH=$GOPATH:$(pwd)

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o smscode .

echo "build Success!!!"