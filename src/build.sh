#!/bin/sh

export GOPATH=$GOPATH:/Users/Zhangjd/Tests/ufop-demo
GOOS=linux GOARCH=amd64 go build -o ../deploy/qufop qufop.go
