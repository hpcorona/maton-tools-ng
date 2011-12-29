#!/bin/sh

export GOOS=windows
export GOARCH=386

8g -o _go_.8 ng.go
8l -o ng.exe _go_.8
