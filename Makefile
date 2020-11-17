GO ?= go
VERSION ?= $(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')

.PHONY: build
build:
	$(GO) build -ldflags '-s -w -X "go.jolheiser.com/tmpl/cmd.Version=$(VERSION)"'

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: test
test:
	$(GO) test -race ./...
