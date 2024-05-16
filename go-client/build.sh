#!/bin/bash
cd example
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o runner .
mv runner ..
cd ..

docker build . -t go-example:1.0.0
