#!/bin/bash -e
path=$1
dir=$(dirname "$path")
filename=$(basename "$path")
extension="${filename##*.}"
nameonly="${filename%.*}"

goimports -w .
go build -o ~/bin/figo ./cmd/figo
go test -coverprofile /tmp/c.out ./...
#uncover /tmp/c.out
#go install ./cmd/...

