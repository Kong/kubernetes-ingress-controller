#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(dirname ${BASH_SOURCE})/../..

cd $REPO_ROOT

${REPO_ROOT}/bin/kustomize build config/base > deploy/single-v2/all-in-one-dbless.yaml
${REPO_ROOT}/bin/kustomize build config/variants/postgres > deploy/single-v2/all-in-one-postgres.yaml
${REPO_ROOT}/bin/kustomize build config/variants/enterprise > deploy/single-v2/all-in-one-enterprise-dbless.yaml
${REPO_ROOT}/bin/kustomize build config/variants/enterprise-postgres > deploy/single-v2/all-in-one-enterprise-postgres.yaml
