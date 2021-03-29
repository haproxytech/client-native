ARG SWAGGER_VERSION
FROM quay.io/goswagger/swagger:$SWAGGER_VERSION

ARG UID
ARG GID
COPY script.sh .
VOLUME ["/data"]

RUN addgroup -g "$GID" -S docker && adduser -u "$UID" -S user -G docker && \
    chmod +x ./script.sh && mkdir -p /docker-generation/data && \
    echo "module github.com/haproxytech" | tee /docker-generation/go.mod && \
    chown -R "${UID}:${GID}" /docker-generation && \
    chown -R "${UID}:${GID}" /data

USER user
ENTRYPOINT ["./script.sh"]
