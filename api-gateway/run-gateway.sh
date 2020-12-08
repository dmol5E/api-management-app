#!/bin/sh

envoy --service-cluster "api-gateway" --service-node "${POD_HOSTNAME}" --drain-time-s 45 --parent-shutdown-time-s 60 -c '/envoy/envoy.yaml'