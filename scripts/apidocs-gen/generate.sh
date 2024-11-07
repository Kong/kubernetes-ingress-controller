#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(dirname "${BASH_SOURCE[0]}")/../.."
CRD_REF_DOCS_BIN="$1"
# Fetch the version of `kubernetes-configuration` from `go.mod`.
KUBE_CONF_VERSION=$(grep "github.com/kong/kubernetes-configuration" ${SCRIPT_ROOT}/go.mod   | awk '{print $2}')
GOPATH=$(go env GOPATH)

generate() {
  echo "INFO: generating API docs for ${1} package, output: ${2}"
  ${CRD_REF_DOCS_BIN} \
      --source-path="${1}" \
      --config="${SCRIPT_ROOT}/scripts/apidocs-gen/config.yaml" \
      --templates-dir="${SCRIPT_ROOT}/scripts/apidocs-gen/template" \
      --renderer=markdown \
      --output-path="${SCRIPT_ROOT}${2}" \
      --max-depth=10
}

generate "${GOPATH}/pkg/mod/github.com/kong/kubernetes-configuration@${KUBE_CONF_VERSION}/api/configuration" "/docs/api-reference.md"
generate "${GOPATH}/pkg/mod/github.com/kong/kubernetes-configuration@${KUBE_CONF_VERSION}/api/incubator" "/docs/incubator-api-reference.md"
