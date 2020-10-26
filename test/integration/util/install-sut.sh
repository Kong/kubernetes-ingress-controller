#!/bin/bash

set -e

REMOTE_IMAGE="localhost:5000/kic:local"

docker tag "$KIC_IMAGE" "$REMOTE_IMAGE"
docker push "$REMOTE_IMAGE"

REPO_ROOT="$(dirname "$BASH_SOURCE")/../../.."
IMAGE="test-local-registry:5000/kic:local"
kubectl apply -f "$REPO_ROOT/deploy/single/all-in-one-dbless.yaml"
kubectl patch -n kong deploy/ingress-kong --patch "
{
	\"spec\": {
		\"template\": {
			\"spec\": {
				\"containers\": [
					{
						\"name\": \"ingress-controller\",
						\"image\": \"$IMAGE\",
						\"env\": [
							{
								\"name\": \"CONTROLLER_ANONYMOUS_REPORTS\",
								\"value\": \"false\"
							}
						]
					}
				]
			}
		}
	}
}"


kubectl rollout status --watch --namespace=kong deploy/ingress-kong --timeout=0
kubectl wait --for=condition=Available --namespace=kong deploy/ingress-kong --timeout=120s
