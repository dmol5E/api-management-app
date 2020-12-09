@echo off
set namespace=%1
kubectl create namespace %namespace%
kubectl create serviceaccount operator --namespace %namespace%
kubectl create clusterrole deployment-operator --verb=get,watch,list,create,update,patch --resource=deployments.apps
kubectl create clusterrolebinding deployment-operator --clusterrole=deployment-operator --serviceaccount=%namespace%:operator
kubectl apply -f .\operator\deployment.yaml -n %namespace%