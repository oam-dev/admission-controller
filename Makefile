TAG ?= latest

PROJECT_NAME := "admission-controller"
PKG := "github.com/oam-dev/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: manifest
manifest:
	./hack/update-codegen.sh

.PHONY: build
build: manifest
	docker build -t oamdev/admission:${TAG} .

.PHONY: publish
publish: build
	docker push oamdev/admission:${TAG}

.PHONY: lint
lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

.PHONY: vet
vet: ## Run go vet
	@go vet ${PKG_LIST}

.PHONY: test
test: ## Run unittests
	@go test -short ${PKG_LIST}

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST}
	@cat cover.out

.PHONY: build-binary
build-binary: ## Build the binary file
	@go build -i $(PKG)/cmd/admission