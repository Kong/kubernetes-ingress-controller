#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..
STAGE_DIR=$(mktemp -d)

mkdir $STAGE_DIR/rbac
mkdir $STAGE_DIR/base
mkdir $STAGE_DIR/prometheus

cd $SCRIPT_ROOT

# split RBAC roles from kubebuilder generated roles
kustomize build ./config/rbac/ -o $STAGE_DIR/rbac
cp $STAGE_DIR/rbac/rbac.authorization.k8s.io_v1_clusterrole_manager-role.yaml $STAGE_DIR/base/manager-role.yaml
cp $STAGE_DIR/rbac/rbac.authorization.k8s.io_v1_clusterrolebinding_manager-rolebinding.yaml $STAGE_DIR/base/manager-rolebinding.yaml
cp $STAGE_DIR/rbac/default_rbac.authorization.k8s.io_v1_role_leader-election-role.yaml $STAGE_DIR/base/leader-election-role.yaml
cp $STAGE_DIR/rbac/default_rbac.authorization.k8s.io_v1_rolebinding_leader-election-rolebinding.yaml $STAGE_DIR/base/leader-election-rolebinding.yaml
# TODO auth proxy service stuff?

cp -Lr ./config/base $STAGE_DIR
cp -Lr ./config/crd $STAGE_DIR/base/
cp -Lr ./config/variants $STAGE_DIR

kustomize build $STAGE_DIR/base > deploy/single-v2/all-in-one-dbless.yaml
kustomize build $STAGE_DIR/variants/postgres > deploy/single-v2/all-in-one-postgres.yaml
kustomize build $STAGE_DIR/variants/enterprise > deploy/single-v2/all-in-one-enterprise-dbless.yaml
kustomize build $STAGE_DIR/variants/enterprise-postgres > deploy/single-v2/all-in-one-enterprise-postgres.yaml
kustomize build $STAGE_DIR/prometheus > deploy/single-v2/all-in-one-prometheus.yaml

rm -rf $STAGE_DIR
