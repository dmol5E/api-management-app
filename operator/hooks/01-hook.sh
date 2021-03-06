#!/usr/bin/env bash

if [[ $1 == "--config" ]]; then
  cat <<EOF
  configVersion: v1
  onStartup: 1
EOF
else
  kubectl apply -f /crd.yaml
  kubectl create clusterrole deployment-operator --verb=get,watch,list,create,update,patch --resource=deployments.apps
  kubectl create clusterrole custom-resource-manager --verb=get,watch,list,create,update,patch --resource=customresourcedefinitions.apiextensions.k8s.io
  kubectl create clusterrole api-config-manager --verb=get,watch,list,create,update,patch --resource=apiconfigs.apimanagement.cloud
  kubectl create clusterrolebinding deployment-operator --clusterrole=deployment-operator --serviceaccount=%namespace%:operator
  kubectl create clusterrolebinding custom-resource-manager --clusterrole=custom-resource-manager --serviceaccount=%namespace%:default
  kubectl create clusterrolebinding api-config-manager --clusterrole=api-config-manager --serviceaccount=%namespace%:default
  kubectl apply -f /deploy/api-gateway-deployment.yaml -n "${CLOUD_NAMESPACE}"
  kubectl apply -f /deploy/api-gateway-service.yaml -n "${CLOUD_NAMESPACE}"
  kubectl apply -f /deploy/api-publisher-deployment.yaml -n "${CLOUD_NAMESPACE}"
  kubectl apply -f /deploy/api-publisher-service.yaml -n "${CLOUD_NAMESPACE}"
fi