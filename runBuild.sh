#!/usr/bin/env bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./build/crontab main.go
cp -rf ./configs ./build
cp -rf ./deploy/* ./build/