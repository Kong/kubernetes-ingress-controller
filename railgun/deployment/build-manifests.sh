#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..

# k4k8s crds
kustomize build ./kong-ingress-controller > ./kong-ingress-controller.yaml
