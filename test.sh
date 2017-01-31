#!/usr/bin/env bash

set -e -x

echo -n "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -v -coverprofile=profile.out $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
go tool cover -html=coverage.txt -o coverage.html
