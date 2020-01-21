#!/usr/bin/env bash
mkdir build
cd build || exit 233
go build -buildmode=c-shared -ldflags "-s -w" -o libqvb.so ../main.go
