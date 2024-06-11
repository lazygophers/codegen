PWD = $(shell pwd)
MODCACHE = $(shell go env GOMODCACHE)
GOCACHE = $(shell go env GOCACHE)

.PHONY: run
run: ## 本地运行
	go mod tidy
	CGO_ENABLED=0 \
	go run --tags debug --ldflags '-s -w --extldflags "-static -fpic"'-v ./cmd/

.PHONY: build
build: ## 打包
	GOVERSION=$(shell go version | awk '{print $$3;}') \
	goreleaser --clean --snapshot --skip=publish,validate

.PHONY: lint
lint: ## 校验
	go build -v -o /dev/null ./cmd/
	golangci-lint run -v

help: ## help
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | sed 's/Makefile://' | awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-z0-9A-Z_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 }'
