ARG SWAGGER_VERSION=v0.23.0
FROM quay.io/goswagger/swagger:$SWAGGER_VERSION

ARG SPEC_URL=https://raw.githubusercontent.com/haproxytech/client-native/master/specification/build/haproxy_spec.yaml
ARG UID
ARG GID
COPY script.sh .
VOLUME ["/data"]

RUN addgroup -g "$GID" -S docker && adduser -u "$UID" -S user -G docker && \
    chmod +x ./script.sh && mkdir -p /docker-generation/dataplaneapi && \
    echo "module github.com/haproxytech" | tee /docker-generation/go.mod && \
    chown -R "${UID}:${GID}" /docker-generation && \
    chown -R "${UID}:${GID}" /data && \
    wget -O /docker-generation/copyright.txt https://raw.githubusercontent.com/haproxytech/client-native/master/specification/copyright.txt && \
    wget -O /docker-generation/haproxy_spec.yaml "${SPEC_URL}"

USER user
ENTRYPOINT ["./script.sh"]
