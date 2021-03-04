#!/bin/bash
set -e

go run specification/build/build.go -file specification/haproxy-spec.yaml > specification/build/haproxy_spec.yaml

PROJECT_PATH=${PWD}
swagger generate model -f ${PROJECT_PATH}/specification/build/haproxy_spec.yaml -r ${PROJECT_PATH}/specification/copyright.txt -m models -t ${PROJECT_PATH}
