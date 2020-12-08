@echo off
set namespace=%1
kubectl apply -f ../deployment.yaml -n %namespace%