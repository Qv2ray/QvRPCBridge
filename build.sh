#!/usr/bin/env bash
cd build || exit 233
go build -buildmode=c-archive -ldflags "-s -w" -o libqvb.a ../main.go
