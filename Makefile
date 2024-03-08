PWD = $(shell pwd)
MODCACHE = $(shell go env GOMODCACHE)
GOCACHE = $(shell go env GOCACHE)

.PHONY: build
build: ## build
	GOVERSION=$(shell go version | awk '{print $$3;}') goreleaser --clean --snapshot --skip=publish,validate --timeout=24h

.PHONY: lint
lint: ## lint go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run -v

help: ## Show this help
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | sed 's/Makefile://' | awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-z0-9A-Z_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 }'
