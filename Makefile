PROJECT_PATH=${PWD}
DOCKER_HAPROXY_VERSION?=2.4
SWAGGER_VERSION=v0.30.2
GO_VERSION:=${shell go mod edit -json | jq -r .Go}

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

.PHONY: models-native
models-native: spec
	swagger generate model -f ${PROJECT_PATH}/specification/build/haproxy_spec.yaml -r ${PROJECT_PATH}/specification/copyright.txt -m models -t ${PROJECT_PATH}

.PHONY: models
models: spec
	cd build/models;docker build \
		--build-arg SWAGGER_VERSION=${SWAGGER_VERSION} \
		--build-arg UID=$(shell id -u) \
		--build-arg GID=$(shell id -g) \
		-t client-native-models .
	docker run --rm -it -v "$(PWD)":/data client-native-models
	ls -lah models

.PHONY: lint
lint:
	docker run --rm -v $(pwd):/data cytopia/yamllint .
	docker run --rm ghcr.io/haproxytech/go-linter:1.46.2 --version
	docker run --rm -v ${PROJECT_PATH}:/app -w /app ghcr.io/haproxytech/go-linter:1.46.2 --timeout 5m --color always --max-issues-per-linter 0 --max-same-issues 0

.PHONY: gofumpt
gofumpt:
	gofumpt -l -w .
