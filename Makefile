.PHONY: test vet build check fmt

test:
	go test ./...

vet:
	go vet ./...

build:
	go build -o ./bin/pad ./cmd/pad

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

check:
	./scripts/check.sh
