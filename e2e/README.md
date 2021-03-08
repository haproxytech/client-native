# ![HAProxy](../assets/images/haproxy-weblogo-210x49.png "HAProxy")

## e2e Testing

testing can be done in two ways, locally if HAProxy is installed, or in docker image

## local machine

```bash
go test -tags integration ./...
```

or

```bash
make e2e
```

## docker environment

```bash
export HAPROXY_VERSION=2.4; sed -e "s/alpine:2.3/alpine:$HAPROXY_VERSION/g" e2e/Dockerfile-TestEnv | docker build -t test_env -f - .
docker build -f e2e/Dockerfile -t client-native-test .
docker run --entrypoint "go" client-native-test test -tags integration ./...
```

or

```bash
make e2e-docker
```

where HAPROXY_VERSION is set to desired version of HAProxy
