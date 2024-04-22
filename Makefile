PROGRAM_NAME = st-test

BUILD_VERSION=$(shell git describe --tags)
BUILD_DATE=$(shell date +%FT%T%z)
BUILD_COMMIT=$(shell git rev-parse --short HEAD)

PKG_PATH=st-test/cmd/util

LDFLAGS=-X ${PKG_PATH}.buildVersion=$(BUILD_VERSION) -X ${PKG_PATH}.buildDate=$(BUILD_DATE) -X ${PKG_PATH}.buildCommit=$(BUILD_COMMIT)

# TOOLS_PATH defines path to Golang-based utility binaries.
TOOLS_PATH=bin/tools

# Pattern to ignore while measuring test coverage.
COVERAGE_IGNORE_PATTERN="/mocks/"

golangci-lint=${TOOLS_PATH}/golangci-lint
gofumpt=${TOOLS_PATH}/gofumpt

.PHONY: help dep fmt test

$(gofumpt): Makefile
	GOBIN=`pwd`/$(TOOLS_PATH) go install mvdan.cc/gofumpt@v0.6.0

$(golangci-lint): Makefile
	GOBIN=`pwd`/$(TOOLS_PATH) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2

clean:
	rm -f ./bin/${PROGRAM_NAME}

setup: $(gofumpt) $(golangci-lint)

fmt: $(gofumpt) ## Format the source files
	$(gofumpt) -l -w .

dep: ## Get the dependencies
	go mod download

test: dep ## Run tests
	go test -timeout 5m -race -covermode=atomic -coverprofile=.coverage.out ./... && \
	grep -vE $(COVERAGE_IGNORE_PATTERN) .coverage.out > .coverage.filtered.out && \
	go tool cover -func=.coverage.out | tail -n1 | awk '{print "Total test coverage: " $$3}'
	@rm .coverage.out

cover: dep ## Run app tests with coverage report
	go test -timeout 5m -race -covermode=atomic -coverprofile=.coverage.out ./... && \
	grep -vE $(COVERAGE_IGNORE_PATTERN) .coverage.out > .coverage.filtered.out && \
	go tool cover -html=.coverage.out -o .coverage.html
	## Open coverage report in default system browser
	xdg-open .coverage.html
	## Remove coverage report
	sleep 2 && rm -f .coverage.out .coverage.html

build: clean
	go build -ldflags "${LDFLAGS}" -o ./bin/st-test ./cmd

lint: lint/sources lint/openapi ## Run all linters

lint/sources: ## Lint the source files
	$(golangci-lint) run --timeout 5m
	govulncheck ./...

lint/openapi: ## Lint openapi specifications
	@echo "Lint OpenAPI specifications"
	@for spec in $(OPENAPI_SPECS) ; do echo "* lint $$spec"; vacuum lint -t -q -x $$spec ; done