REVISION ?= $(shell git rev-parse --short=7 HEAD)
GO := CGO_ENABLED=0 GO111MODULE=on go

all: build

# openapi code generetation

generate: api

api: api/client/go api/client/javascript api/server/go

OPENAPI ?= docker run --rm \
		--user=$(shell id -u $(USER)):$(shell id -g $(USER)) \
		-v $(shell pwd):$(shell pwd) \
		openapitools/openapi-generator-cli:v5.0.0-beta3


api/client/go: api/api.yaml
	-rm -rf $@
	$(OPENAPI) generate -i $(shell pwd)/api/api.yaml -g go -o $(shell pwd)/api/client/go --additional-properties=withGoCodegenComment=true
	-rm -rf $@/{go.mod,main.go}
	touch $@

api/client/javascript: api/api.yaml
	-rm -rf $@
	$(OPENAPI) generate -i $(shell pwd)/api/api.yaml -g javascript -o $(shell pwd)/api/client/javascript --additional-properties=usePromises=true
	touch $@

api/server/go: api/api.yaml
	-rm -rf $@
	$(OPENAPI) generate -i $(shell pwd)/api/api.yaml -g go-server -o $(shell pwd)/api/server/go
	-rm -rf $@/{go.mod,main.go}
	touch $@

# build related

.PHONY: build
build: \
	cmd/api/api

.PHONY: cmd/api/api
cmd/api/api:
	$(GO) build -v -o ./cmd/api/api ./cmd/api