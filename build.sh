#!/user/bin/env bash
set -xe

# install package and dependencies

go get "github.com/akrylysov/algnhsa"
go get "github.com/gorilla/mux"

# build command
go build -o bin/application application.go