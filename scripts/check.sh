#!/usr/bin/env sh
set -eu

test -z "$(gofmt -l .)"
./scripts/check-template.sh
go test ./...
go vet ./...
go build -o /tmp/pad ./cmd/pad
