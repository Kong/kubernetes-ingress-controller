#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(dirname "${BASH_SOURCE[0]}")/../.."
CRD_REF_DOCS_BIN="$1"

${CRD_REF_DOCS_BIN} \
    --source-path="${SCRIPT_ROOT}/pkg/apis/configuration/" \
    --config="${SCRIPT_ROOT}/scripts/apidocs-gen/config.yaml" \
    --templates-dir="${SCRIPT_ROOT}/scripts/apidocs-gen/template" \
    --renderer=markdown \
    --output-path="${SCRIPT_ROOT}/docs/api-reference.md"
