#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..

COMMON="namespace.yaml custom-types.yaml rbac.yaml service.yaml"
COMMON+=" custom-server-block.yaml validation-service.yaml"
DB="postgres.yaml migration.yaml kong-ingress-postgres.yaml"
DBLESS="kong-ingress-dbless.yaml"

MANIFEST=$(cd ${SCRIPT_ROOT}/deploy/manifests; cat ${COMMON} ${DB})
echo "${MANIFEST}" > ${SCRIPT_ROOT}/deploy/single/all-in-one-postgres.yaml

MANIFEST=$(cd ${SCRIPT_ROOT}/deploy/manifests; cat ${COMMON} ${DBLESS})
echo "${MANIFEST}" > ${SCRIPT_ROOT}/deploy/single/all-in-one-dbless.yaml

