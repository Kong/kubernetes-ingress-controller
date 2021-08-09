#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..
CONFIG_ROOT=config/

cd $SCRIPT_ROOT

kustomize build $CONFIG_ROOT/base > deploy/single-v2/all-in-one-dbless.yaml
kustomize build $CONFIG_ROOT/variants/postgres > deploy/single-v2/all-in-one-postgres.yaml
kustomize build $CONFIG_ROOT/variants/enterprise > deploy/single-v2/all-in-one-enterprise-dbless.yaml
kustomize build $CONFIG_ROOT/variants/enterprise-postgres > deploy/single-v2/all-in-one-enterprise-postgres.yaml
