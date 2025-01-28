#!/bin/bash

set -ue

DIR=$(dirname $0)
pushd "${DIR}/cmd/agent"
go build -o agent *.go
popd
pushd "${DIR}/cmd/server"
go build -o server *.go
popd
