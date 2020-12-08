#!/usr/bin/env bash

if [[ $1 == "--config" ]]; then
  cat <<EOF
  configVersion: v1
  onStartup: 1
EOF
else
  kubectl apply -f /deploy/api-gateway-deployment.yaml -n "${CLOUD_NAMESPACE}"
  kubectl apply -f /deploy/api-gateway-service.yaml -n "${CLOUD_NAMESPACE}"
  kubectl apply -f /deploy/api-publisher-deployment.yaml -n "${CLOUD_NAMESPACE}"
  kubectl apply -f /deploy/api-publisher-service.yaml -n "${CLOUD_NAMESPACE}"
fi