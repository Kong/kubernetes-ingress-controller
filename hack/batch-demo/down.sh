#!/bin/bash

NAMESPACE=batch-demo

kubectl delete ing --namespace $NAMESPACE --all
kubectl delete svc --namespace $NAMESPACE --all
