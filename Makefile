PWD = $(shell pwd)
MODCACHE = $(shell go env GOMODCACHE)
GOCACHE = $(shell go env GOCACHE)

.PHONY: build
build: ## build
	#export PATH="$PATH:/home/go/lazygophers/codegen/dist/cli_linux_amd64_v3/:/root/go/bin/"
	#export PATH="$PATH:/Users/luoxin/persons/go/lazygophers/codegen/codegen/dist/cli_darwin_arm64/:/root/go/bin/"

	# cp ./codegen.cfg.yaml ~/Library/Application Support/lazygophers/codegen.cfg.yaml

	GOVERSION=$(shell go version | awk '{print $$3;}') \
	goreleaser build --clean --snapshot --single-target --config=debug.goreleaser.yaml
	#GOVERSION=$(shell go version | awk '{print $$3;}') goreleaser --clean --snapshot --skip=publish,validate --timeout=24h

.PHONY: lint
lint: ## lint go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run -v

help: ## Show this help
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | sed 's/Makefile://' | awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-z0-9A-Z_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 }'

# make build && ./dist/cli_darwin_arm64/codegen gen -i ../example/example.proto -d add-rpc -m user -a add:,admin
