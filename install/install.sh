#!/usr/bin/env bash

namespace=$1
kubectl apply ../deployment.yaml -n "${namespace}"

