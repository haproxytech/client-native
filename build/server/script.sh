#!/bin/sh
set -e

swagger version
cp /data/configure_data_plane.go /docker-generation/dataplaneapi/configure_data_plane.go
swagger generate server -f /docker-generation/haproxy_spec.yaml \
    -A "Data Plane" \
    -t /docker-generation/ \
    --existing-models github.com/haproxytech/client-native/v2/models \
    --exclude-main \
    --skip-models \
    -s dataplaneapi \
    -r /docker-generation/copyright.txt
rm /data/doc.go
rm /data/embedded_spec.go
rm /data/server.go
rm -rf /data/operations/*
cp -a /docker-generation/dataplaneapi/. /data
