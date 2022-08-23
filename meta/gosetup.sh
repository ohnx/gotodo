#!/bin/bash

export GOPATH=$PWD
export PATH=$PATH:$GOPATH/bin

# dependencies
go get github.com/julienschmidt/httprouter
go get github.com/mattn/go-sqlite3
go get github.com/rs/cors

# own code
mkdir -p $PWD/src/github.com/ohnx
git clone https://github.com/ohnx/gotodo.git $PWD/src/github.com/ohnx/
