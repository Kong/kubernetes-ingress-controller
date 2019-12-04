#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..

# k4k8s
kustomize build ./deploy/manifests/base > deploy/single/all-in-one-dbless.yaml
# k4k8s with DB
kustomize build ./deploy/manifests/postgres \
  > deploy/single/all-in-one-postgres.yaml
# k4k8s Enterprise
kustomize build ./deploy/manifests/enterprise-k8s \
  > deploy/single/all-in-one-dbless-k4k8s-enterprise.yaml
## k4k8s Enterprise with DB
#kustomize build ./deploy/manifests/enterprise-k8s \
#  > deploy/single/all-in-one-postgres-enterprise-k8s.yaml
# Kong Enterprise
kustomize build ./deploy/manifests/enterprise \
  > deploy/single/all-in-one-postgres-enterprise.yaml

# Kong Dev Config
cat  ./deploy/manifests/base/namespace.yaml \
     ./deploy/manifests/base/custom-types.yaml \
     ./deploy/manifests/base/rbac.yaml \
     ./deploy/manifests/base/service.yaml \
     ./deploy/manifests/base/custom-server-block.yaml \
     ./deploy/manifests/base/validation-service.yaml \
     ./deploy/manifests/dev-env/kong-ingress-dbless.yaml \
     > hack/dev-config.yaml
