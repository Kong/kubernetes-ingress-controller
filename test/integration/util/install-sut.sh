#!/bin/bash

set -e

REMOTE_IMAGE="localhost:5000/kic:local"

docker tag "$KIC_IMAGE" "$REMOTE_IMAGE"
docker push "$REMOTE_IMAGE"

REPO_ROOT="$(dirname "$BASH_SOURCE")/../../.."
MANIFEST="$(
	sed 's!^\( *image: \)kong-docker-kubernetes-ingress-controller.*$!\1test-local-registry:5000/kic:local!' \
		"$REPO_ROOT/deploy/single/all-in-one-dbless.yaml"
)"

kubectl apply -f <(echo "$MANIFEST")
kubectl wait --for=condition=Available --namespace=kong deploy/ingress-kong --timeout=120s
