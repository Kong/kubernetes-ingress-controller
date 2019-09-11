#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..

kustomize build ./deploy/manifests/base > deploy/single/all-in-one-dbless.yaml
kustomize build ./deploy/manifests/postgres \
  > deploy/single/all-in-one-postgres.yaml
kustomize build ./deploy/manifests/enterprise \
  > deploy/single/all-in-one-postgres-enterprise.yaml

