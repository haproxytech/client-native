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

ARG HAPROXY_VERSION=2.3

FROM haproxytech/haproxy-alpine:${HAPROXY_VERSION}

RUN apk add socat openssl util-linux htop tzdata bash gcc openssl wget go && \
    rm -f /usr/local/bin/dataplaneapi /usr/bin/dataplaneapi && \
    rm -f /usr/local/etc/haproxy/haproxy.cfg && \
    chgrp -R haproxy /etc/haproxy /run /var &&\
    chmod -R g+w /etc/haproxy /run /var

ENV GODEBUG="madvdontneed=1" \
	PATH=$PATH:/usr/local/go/bin \
	HAPROXY_VERSION=$HAPROXY_VERSION

WORKDIR /src

COPY /go.mod /src/go.mod
COPY /go.sum /src/go.sum
RUN go mod download

COPY / /src

ENTRYPOINT ["go"]
CMD ["test", "-tags", "integration", "./..."]
