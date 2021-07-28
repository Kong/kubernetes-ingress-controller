#!/bin/bash -ex
DB=$1
VERSION=$2
SCRIPT_ROOT=$(dirname ${BASH_SOURCE})

# kill existing kubectl proxy or port-forwards
# TODO: make this less draconian
set +e
pkill -f kubectl
set -e

# helper function
function usage {
  echo "Usage: $0 <mode> <version>"
  echo "db: true or false"
  echo "Version is kong's version without dots.Examples: 11, 12, 13 etc"
  exit 1
}

# perform sanity checks
if [ -z "$VERSION" ];
then
  usage
fi

if [ "${DB}" != "true" ] && [ "${DB}" != "false" ];
then
  usage
fi

# select the kong pod
export POD_NAMESPACE="kong-dev"
export POD_NAME=$(kubectl get po -n $POD_NAMESPACE \
  -l=app=kong,version=$VERSION,db=${DB} \
  -o custom-columns=":metadata.name" --no-headers)

# set up port for communication with k8s api-server
kubectl proxy --port=8002 1> /dev/null &

# set up ports for communicating with kong
kubectl port-forward -n $POD_NAMESPACE $POD_NAME \
  8000:8000 8443:8443 8001:8001 1> /dev/null &

# run rabbit run
go run ./internal/ingress/controller/cli \
  --apiserver-host http://127.0.0.1:8002 \
  --publish-service kong-dev/kong-fake-publish-proxy \
  --kong-admin-url=http://127.0.0.1:8001
