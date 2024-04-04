# Copyright 2019 HAProxy Technologies LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG HAPROXY_VERSION=2.7
ARG DOCKER_URL=haproxytech
ARG GO_VERSION=1.22

FROM golang:${GO_VERSION}-alpine as go_image

FROM ${DOCKER_URL}/haproxy-alpine:${HAPROXY_VERSION}

RUN apk add --no-cache -U socat openssl util-linux htop tzdata bash gcc openssl wget git curl make bash jq git && \
    rm -f /usr/local/bin/dataplaneapi /usr/bin/dataplaneapi && \
    rm -f /usr/local/etc/haproxy/haproxy.cfg && \
    chgrp -R haproxy /etc/haproxy /run /var &&\
    chmod -R g+w /etc/haproxy /run /var &&\
    mkdir -p /root/go
COPY --from=go_image /usr/local/go/ /usr/local/go/

ENV PATH=$PATH:/usr/local/go/bin \
    CGO_ENABLED=0 \
    HAPROXY_VERSION=$HAPROXY_VERSION
RUN GOBIN=/usr/local/go/bin go install github.com/oktalz/gotest@latest

WORKDIR /src

COPY /go.mod /src/go.mod
COPY /go.sum /src/go.sum
RUN go mod download

COPY / /src
ENTRYPOINT ["gotest","-t", "integration"]
