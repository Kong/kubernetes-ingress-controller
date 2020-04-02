#!/bin/bash -e
MODE=$1
VERSION=$2
SCRIPT_ROOT=$(dirname ${BASH_SOURCE})

function usage {
  echo "Usage: $0 <mode> <version>"
  echo "mode: db or dbless"
  echo "version: 1.4, 2.0,..."
  exit 1
}

if [ -z "$VERSION" ];
then
  usage
fi

if [ "${MODE}" != "db" ] && [ "${MODE}" != "dbless" ];
then
  usage
fi

kustomize build "${SCRIPT_ROOT}/${MODE}/${VERSION}" | kubectl apply -f -
