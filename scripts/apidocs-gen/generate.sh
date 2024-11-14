#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(dirname "${BASH_SOURCE[0]}")/../.."
CRD_REF_DOCS_BIN="$1"
# Download the `kong/kubernetes-configuration` package and get its path.
KUBE_CONF_REPO=github.com/kong/kubernetes-configuration
KUBE_CONF_PATH=$(go mod download -json ${KUBE_CONF_REPO} | jq -rM .Dir)
echo "Dowloaded ${KUBE_CONF_REPO} in ${KUBE_CONF_PATH}"

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

generate "${KUBE_CONF_PATH}/api/configuration" "/docs/api-reference.md"
generate "${KUBE_CONF_PATH}/api/incubator" "/docs/incubator-api-reference.md"
