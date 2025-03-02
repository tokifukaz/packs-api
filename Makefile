.PHONY: all build test fmt clean version

SHELL=/bin/bash -o pipefail

VERSION=0.1.0
GOBUILD=env GOOS=linux GOARCH=amd64 go build -v
GOCLEAN=go clean

# Name of the binary to be built
BINARY_NAME=packs-api

build: generate
	$(GOBUILD) -o $(BINARY_NAME) cmd/*.go

test:
	go test -v -cover ./...

fmt:
	go fmt ./...

clean:
	$(GOCLEAN) ./...

generate:
	go generate ./...

version:
	@echo $(VERSION)