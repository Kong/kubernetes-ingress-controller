#!/bin/bash -x
export POD_NAMESPACE=$1
if [[ -z "$POD_NAMESPACE" ]];
then
  echo "Usage: $0 <namespace>"
  exit 1
fi

export POD_NAME=`kubectl get po -n $POD_NAMESPACE -l=app=ingress-kong  -o custom-columns=":metadata.name" --no-headers`

set +e
pkill -f kubectl
set -e

kubectl proxy --port=8002 1> /dev/null &
kubectl port-forward -n $POD_NAMESPACE $POD_NAME 8001:8001 1> /dev/null &

go run ./cli/ingress-controller \
  --apiserver-host http://127.0.0.1:8002 \
  --publish-service kong/kong-proxy \
  --kong-url=http://127.0.0.1:8001
