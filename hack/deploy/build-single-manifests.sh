#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..
STAGE_DIR=$(mktemp -d)

mkdir $STAGE_DIR/base

cd $SCRIPT_ROOT

# TODO auth proxy service stuff?

cp -Lr ./config/base $STAGE_DIR
cp -Lr ./config/rbac $STAGE_DIR/base/
cp -Lr ./config/crd $STAGE_DIR/base/
cp -Lr ./config/variants $STAGE_DIR

kustomize build $STAGE_DIR/base > deploy/single-v2/all-in-one-dbless.yaml
kustomize build $STAGE_DIR/variants/postgres > deploy/single-v2/all-in-one-postgres.yaml
kustomize build $STAGE_DIR/variants/enterprise > deploy/single-v2/all-in-one-enterprise-dbless.yaml
kustomize build $STAGE_DIR/variants/enterprise-postgres > deploy/single-v2/all-in-one-enterprise-postgres.yaml

rm -rf $STAGE_DIR
