@echo off
set namespace=%1
kubectl create namespace %namespace%
kubectl create serviceaccount operator --namespace %namespace%
kubectl create clusterrole deployment-operator --verb=get,watch,list,create,update,patch --resource=deployments.apps
kubectl create clusterrolebinding deployment-operator --clusterrole=deployment-operator --serviceaccount=%namespace%:operator
kubectl create clusterrole custom-resource-manager --verb=get,watch,list,create,update,patch --resource=customresourcedefinitions.apiextensions.k8s.io
kubectl create clusterrole route-config-manager --verb=get,watch,list,create,update,patch --resource=routeconfigs.apimanagement.cloud
kubectl create clusterrolebinding custom-resource-manager --clusterrole=custom-resource-manager --serviceaccount=api-management-app:default
kubectl create clusterrolebinding route-config-manager --clusterrole=route-config-manager --serviceaccount=api-management-app:default
kubectl apply -f .\operator\deployment.yaml -n %namespace%