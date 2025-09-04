#! /usr/bin/bash

set -xe

go test -v -coverprofile cover.out ./...
go tool cover -html cover.out -o cover.html
kde-open cover.html