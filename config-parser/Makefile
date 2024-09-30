PROJECT_PATH=$(shell pwd)
GOLANGCI_LINT_VERSION=1.54.1

.PHONY: generate
generate:
	go run generate/*.go ${PROJECT_PATH}
	$(MAKE) format

.PHONY: format
format:
	go install mvdan.cc/gofumpt@latest
	gofumpt -l -w .

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	cd bin;GOLANGCI_LINT_VERSION=${GOLANGCI_LINT_VERSION} sh lint-check.sh
	bin/golangci-lint run --timeout 5m --color always --max-issues-per-linter 0 --max-same-issues 0
