.PHONY: all build test lint fmt clean run deps

SHELL=/bin/bash -o pipefail

VERSION=0.1.0
GOBUILD=env GOOS=linux GOARCH=amd64 go build -v
GOCLEAN=go clean
LINT_VERSION=v1.56.2

# Name of the binary to be built
BINARY_NAME=packs-api

# Reports
VET_REPORT="vet_report.json"
SEC_REPORT="sec_report.json"

build: generate
	$(GOBUILD) -o $(BINARY_NAME) cmd/*.go

test:
	go test -v -cover ./...

install_lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin $(LINT_VERSION)

lint: install_lint
	golangci-lint run --tests=false --skip-dirs mocks --out-format=checkstyle

fmt:
	go fmt ./...

clean:
	$(GOCLEAN) ./...

deps:
	go -v ./...

tools:
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go mod tidy

vet: tools
	go vet -json ./... 2> $(VET_REPORT)

sec: tools
	gosec -fmt sonarqube -out $(SEC_REPORT) ./...

generate:
	go generate ./...

version:
	@echo $(VERSION)