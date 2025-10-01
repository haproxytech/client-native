PROJECT_PATH=${PWD}
DOCKER_HAPROXY_VERSION?=3.2
SWAGGER_VERSION=v0.32.3
GO_VERSION:=${shell go mod edit -json | jq -r .Go}
GOLANGCI_LINT_VERSION=1.64.5
CHECK_COMMIT=5.2.0

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
models: gentypes spec swagger-check go-method-gen-check
	rm -rf models/
	./bin/swagger generate model --additional-initialism=FCGI -f ${PROJECT_PATH}/specification/build/haproxy_spec.yaml -r ${PROJECT_PATH}/specification/copyright.txt -m models -t ${PROJECT_PATH}
	./bin/go-method-gen --header-file=specification/copyright.txt --scan=models --debug --overrides=models/funcs/overrides.yaml && find ./generated -name "*.go" -exec cp {} ./models \; && rm -rf generated
	go run cmd/struct_equal_generator/*.go -l ${PROJECT_PATH}/specification/copyright.txt ${PROJECT_PATH}/models
	go run cmd/struct_tags_checker/*.go ${PROJECT_PATH}/models
	go run cmd/kubebuilder_marker_generator/*.go  ${PROJECT_PATH}/models
	go run cmd/server_params_runtime/*.go ${PROJECT_PATH}/models
	go run cmd/defaults-setter/main.go ${PROJECT_PATH}/specification/build/haproxy_spec.yaml ${PROJECT_PATH}/models
	$(MAKE) gofumpt

.PHONY: go-method-gen-check
go-method-gen-check:
	@GO_METHOD_GEN_BIN_NAME="go-method-gen"; \
	GO_METHOD_GEN_GITHUB="github.com/haproxytech/go-method-gen/cmd/go-method-gen@latest"; \
	if [ -f "$$GO_METHOD_GEN_BIN_NAME" ]; then \
		echo "✅ $$GO_METHOD_GEN_BIN_NAME already installed"; \
	else \
		GOBIN=$(PWD)/bin go install $$GO_METHOD_GEN_GITHUB && \
		echo "✅ $$GO_METHOD_GEN_BIN_NAME installed"; \
	fi

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

.PHONY: check-commit
check-commit:
	cd bin;CHECK_COMMIT=${CHECK_COMMIT} sh check-commit.sh
	bin/check-commit
