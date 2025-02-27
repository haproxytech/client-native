PROJECT_PATH=${PWD}
DOCKER_HAPROXY_VERSION?=3.0
SWAGGER_VERSION=v0.31.0
GO_VERSION:=${shell go mod edit -json | jq -r .Go}
GOLANGCI_LINT_VERSION=1.64.5

.PHONY: test
test:
	go test ./...

.PHONY: test-equal
test-equal:
	go test -tags equal ./...

.PHONY: e2e
e2e:
	go install github.com/oktalz/gotest@latest
	gotest -t integration

.PHONY: e2e-docker
e2e-docker:
	docker build -f e2e/Dockerfile --build-arg GO_VERSION=${GO_VERSION} --build-arg HAPROXY_VERSION=${DOCKER_HAPROXY_VERSION} -t client-native-test:${DOCKER_HAPROXY_VERSION} .
	docker run --rm -t client-native-test:${DOCKER_HAPROXY_VERSION}

# config-parser auto-generated types
.PHONY: gentypes
gentypes:
	cd config-parser && go run generate/*.go ${PROJECT_PATH}/config-parser

.PHONY: spec
spec:
	go run cmd/specification/*.go -file specification/haproxy-spec.yaml > specification/build/haproxy_spec.yaml

.PHONY: models
models: gentypes spec swagger-check
	./bin/swagger generate model --additional-initialism=FCGI -f ${PROJECT_PATH}/specification/build/haproxy_spec.yaml -r ${PROJECT_PATH}/specification/copyright.txt -m models -t ${PROJECT_PATH}
	rm -rf models/server_params_prepare_for_runtime.go
	rm -rf models/*_compare.go
	rm -rf models/*_compare_test.go
	go run cmd/struct_equal_generator/*.go -l ${PROJECT_PATH}/specification/copyright.txt ${PROJECT_PATH}/models
	go run cmd/struct_tags_checker/*.go ${PROJECT_PATH}/models
	go run cmd/kubebuilder_marker_generator/*.go  ${PROJECT_PATH}/models
	go run cmd/server_params_runtime/*.go ${PROJECT_PATH}/models

.PHONY: swagger-check
swagger-check:
	cd bin; SWAGGER_VERSION=${SWAGGER_VERSION} sh swagger-check.sh

.PHONY: lint
lint:
	cd bin;GOLANGCI_LINT_VERSION=${GOLANGCI_LINT_VERSION} sh lint-check.sh
	bin/golangci-lint run --timeout 5m --color always --max-issues-per-linter 0 --max-same-issues 0

.PHONY: lint-yaml
lint-yaml:
	docker run --rm -v ${PROJECT_PATH}:/data cytopia/yamllint .

.PHONY: gofumpt
gofumpt:
	go install mvdan.cc/gofumpt@latest
	gofumpt -l -w .
