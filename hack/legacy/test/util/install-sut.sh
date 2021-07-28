#!/bin/bash

set -e

REMOTE_IMAGE="localhost:5000/kic:local"

docker tag "$KIC_IMAGE" "$REMOTE_IMAGE"
docker push "$REMOTE_IMAGE"

SUT_ROOT="$(dirname "$BASH_SOURCE")/../sut"
kustomize build --load-restrictor LoadRestrictionsNone "$SUT_ROOT" | kubectl apply -f -

kubectl wait --for=condition=Available --namespace=kong deploy/ingress-kong --timeout=300s
