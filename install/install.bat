@echo off
set namespace=%1
kubectl apply -f .\operator\deployment.yaml -n %namespace%