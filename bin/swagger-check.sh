#!/bin/sh
V=$(./swagger version | grep version | awk 'NF>1{print $NF}')

case "$V" in
  *$SWAGGER_VERSION*) echo "swagger version $V is OK" ;;
  *)          echo "Detected $V, Installing $SWAGGER_VERSION"; GOTOOLCHAIN=go1.24.6 GOBIN=$(pwd) go install "github.com/go-swagger/go-swagger/cmd/swagger@$SWAGGER_VERSION" ;;
esac
