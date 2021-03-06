@echo off
set namespace=%1
kubectl create serviceaccount operator --namespace %namespace%
kubectl apply -f .\roles.yaml -n %namespace%
kubectl create clusterrolebinding crd-manage-operator --clusterrole=custom-resource-manager --serviceaccount=%namespace%:operator
kubectl apply -f .\operator\deployment.yaml -n %namespace%