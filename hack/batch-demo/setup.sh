#!/bin/bash

NAMESPACE=batch-demo

echo "Installing the echoheaders application in namespace $NAMESPACE"

kubectl create namespace $NAMESPACE
kubectl apply --namespace $NAMESPACE -f https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/manifests/dummy-application.yaml
