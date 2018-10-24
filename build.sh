#!/usr/bin/env bash

mkdir -p dist

for GOOS in darwin linux windows; do
    for GOARCH in 386 amd64; do
        output_name=goglobls-$GOOS-$GOARCH
        if [ $GOOS = "windows" ]; then
            output_name+='.exe'
        fi
        go build -v -o dist/$output_name
    done
done