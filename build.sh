#!/usr/bin/env bash
mkdir build
cd build || exit 233
go build -buildmode=c-archive -ldflags "-s -w" -o libqvb.a ../main.go
