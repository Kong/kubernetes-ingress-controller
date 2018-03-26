#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..

WITH_POSTGRESQL="manifests/namespace.yaml manifests/custom-types.yaml manifests/postgres.yaml manifests/rbac.yaml manifests/ingress-controller.yaml provider/baremetal/kong-proxy-nodeport.yaml manifests/kong.yaml"
MANIFEST=$(cd ${SCRIPT_ROOT}/deploy; cat ${WITH_POSTGRESQL})

echo "${MANIFEST}" > ${SCRIPT_ROOT}/deploy/single/all-in-one-postgres.yaml

# TODO: add cassandra deployment
