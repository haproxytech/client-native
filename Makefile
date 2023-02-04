PROJECT_PATH=${PWD}
DOCKER_HAPROXY_VERSION?=2.7
SWAGGER_VERSION=v0.30.2
GO_VERSION:=${shell go mod edit -json | jq -r .Go}
GOLANGCI_LINT_VERSION=1.51.1

.PHONY: test
test:
	go test ./...

.PHONY: e2e
e2e:
	go test -tags integration ./...

.PHONY: e2e-docker
e2e-docker:
	docker build -f e2e/Dockerfile-base --build-arg HAPROXY_VERSION=${DOCKER_HAPROXY_VERSION} --build-arg GO_VERSION=${GO_VERSION} -t test_env:${DOCKER_HAPROXY_VERSION} .
	docker build -f e2e/Dockerfile -t client-native-test:${DOCKER_HAPROXY_VERSION} .
	docker run --rm -it client-native-test:${DOCKER_HAPROXY_VERSION}

.PHONY: spec
spec:
	go run specification/build/build.go -file specification/haproxy-spec.yaml > specification/build/haproxy_spec.yaml

.PHONY: models
models: spec swagger-check
	./bin/swagger generate model --additional-initialism=FCGI -f ${PROJECT_PATH}/specification/build/haproxy_spec.yaml -r ${PROJECT_PATH}/specification/copyright.txt -m models -t ${PROJECT_PATH}

.PHONY: swagger-check
swagger-check:
	cd bin; SWAGGER_VERSION=${SWAGGER_VERSION} sh swagger-check.sh

.PHONY: lint
lint:
	cd bin;GOLANGCI_LINT_VERSION=${GOLANGCI_LINT_VERSION} sh lint-check.sh
	bin/golangci-lint run --timeout 5m --color always --max-issues-per-linter 0 --max-same-issues 0

.PHONY: lint-yaml
lint-yaml:
	docker run --rm -v $(pwd):/data cytopia/yamllint .

.PHONY: gofumpt
gofumpt:
	gofumpt -l -w .
