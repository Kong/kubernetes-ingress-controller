#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(dirname "${BASH_SOURCE[0]}")/../.."
CRD_REF_DOCS_BIN="$1"

generate() {
  echo "INFO: generating API docs for ${1} package, output: ${2}"
  ${CRD_REF_DOCS_BIN} \
      --source-path="${SCRIPT_ROOT}${1}" \
      --config="${SCRIPT_ROOT}/scripts/apidocs-gen/config.yaml" \
      --templates-dir="${SCRIPT_ROOT}/scripts/apidocs-gen/template" \
      --renderer=markdown \
      --output-path="${SCRIPT_ROOT}${2}" \
      --max-depth=10
}

# TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/6618
# Bring the generator back, relying on kong/kubernetes-configuration
# generate "/configuration" "/docs/api-reference.md"
# generate "/incubator" "/docs/incubator-api-reference.md"
