#!/user/bin/env bash
set -xe

# install package and dependencies

go get ./...

# build command
go build -o bin/application application.go