#!/bin/bash
export GOPATH=$HOME/go

if [[ -z "${CMD_PATH}" ]]; then
  CMD_PATH="./..."
  fi

if [ "${GITHUB_ACTIONS}" == "true" ]; then
    go test "${CMD_PATH}" 1> debug.out
else
    go test -v "${CMD_PATH}"
fi


BINARY="CharVomit"

if go get "${CMD_PATH}" 1>> debug.out; then
    if [ "${GOOS}" == "windows" ]; then
        if [ "${GITHUB_ACTIONS}" == "true" ]; then
            go build -v -ldflags="-X main.gitver=$(git describe --always --long --dirty)" -o ${BINARY}.exe ./cmd/CharVomit 1>> debug.out
            echo "${BINARY}.exe"
        else
            go build -v -ldflags="-X main.gitver=$(git describe --always --long --dirty)" -o ${BINARY}.exe ./cmd/CharVomit
        fi
    else
        if [ "${GITHUB_ACTIONS}" == "true" ]; then
            go build -v -ldflags="-X main.gitver=$(git describe --always --long --dirty)" -o ${BINARY} ./cmd/CharVomit 1>> debug.out
            echo "${BINARY}"
        else
            go build -v -ldflags="-X main.gitver=$(git describe --always --long --dirty)" -o ${BINARY} ./cmd/CharVomit
        fi
    fi
fi
