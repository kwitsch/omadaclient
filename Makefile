.PHONY: generate fmt lint help
.DEFAULT_GOAL:=help

GOARCH?=$(shell go env GOARCH)
GOARM?=$(shell go env GOARM)


generate: ## run go generate and regenerate openapi client from api/openapi.yaml
	go generate ./...
	cd api && go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config oapi-codegen.yaml openapi.yaml

fmt: ## gofmt and goimports all go files
	go run mvdan.cc/gofumpt -l -w -extra .
	find . -name '*.go' -exec go run golang.org/x/tools/cmd/goimports -w {} +

lint: generate fmt ## run golangcli-lint checks
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 5m


help:  ## Shows help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'