REVISION ?= $(shell git rev-parse --short=7 HEAD)
GO := CGO_ENABLED=0 GO111MODULE=on go

all: build

.PHONY: build
build: \
	cmd/api/api

.PHONY: cmd/api/api
cmd/api/api:
	$(GO) build -v -o ./cmd/api/api ./cmd/api