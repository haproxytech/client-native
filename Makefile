PROJECT_PATH=${PWD}
DOCKER_HAPROXY_VERSION?=2.4

.PHONY: test
test:
	go test ./...

.PHONY: e2e
e2e:
	go test -tags integration ./...

.PHONY: e2e-docker
e2e-docker:
	sed -e "s/alpine:2.3/alpine:${DOCKER_HAPROXY_VERSION}/g" e2e/Dockerfile-TestEnv | docker build -t test_env -f - .
	docker build -f e2e/Dockerfile -t client-native-test .
	docker run --entrypoint "go" client-native-test test -tags integration ./...

.PHONY: spec
spec:
	go run specification/build/build.go -file specification/haproxy-spec.yaml > specification/build/haproxy_spec.yaml

.PHONY: models
models: spec
	swagger generate model -f ${PROJECT_PATH}/specification/build/haproxy_spec.yaml -r ${PROJECT_PATH}/specification/copyright.txt -m models -t ${PROJECT_PATH}

.PHONY: lint
lint:
	docker run --rm -v $(pwd):/data cytopia/yamllint .
	golangci-lint run --color always --timeout 120s
