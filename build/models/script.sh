#!/bin/sh

swagger version
swagger generate model -f /data/specification/build/haproxy_spec.yaml -r /data/specification/copyright.txt -m models -t /docker-generation/data
rm -rf /data/models/*
cp /docker-generation/data/models/* /data/models/
