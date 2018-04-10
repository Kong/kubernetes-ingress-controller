#!/bin/bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail
set -x

GOPATH=$(go env GOPATH)
PACKAGE_NAME=github.com/kong/kubernetes-ingress-controller
REPO_ROOT="$GOPATH/src/$PACKAGE_NAME"


SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
${CODEGEN_PKG}/generate-groups.sh "deepcopy" \
  ${PACKAGE_NAME}/internal ${PACKAGE_NAME}/internal \
  .:ingress \
  --output-base "$GOPATH/src" \
  --go-header-file ${SCRIPT_ROOT}/hack/boilerplate/boilerplate.go.txt

${CODEGEN_PKG}/generate-groups.sh "deepcopy" \
  ${PACKAGE_NAME}/internal/client/plugin ${PACKAGE_NAME}/internal/apis \
  admin:v1 \
  --output-base "$GOPATH/src" \
  --go-header-file ${SCRIPT_ROOT}/hack/boilerplate/boilerplate.go.txt

${CODEGEN_PKG}/generate-groups.sh "all" \
  ${PACKAGE_NAME}/internal/client/plugin ${PACKAGE_NAME}/internal/apis \
  plugin:v1 \
  --go-header-file ${SCRIPT_ROOT}/hack/boilerplate/boilerplate.go.txt

${CODEGEN_PKG}/generate-groups.sh "all" \
  ${PACKAGE_NAME}/internal/client/consumer ${PACKAGE_NAME}/internal/apis \
  consumer:v1 \
  --go-header-file ${SCRIPT_ROOT}/hack/boilerplate/boilerplate.go.txt

${CODEGEN_PKG}/generate-groups.sh "all" \
  ${PACKAGE_NAME}/internal/client/credential ${PACKAGE_NAME}/internal/apis \
  credential:v1 \
  --go-header-file ${SCRIPT_ROOT}/hack/boilerplate/boilerplate.go.txt

${CODEGEN_PKG}/generate-groups.sh "all" \
  ${PACKAGE_NAME}/internal/client/configuration ${PACKAGE_NAME}/internal/apis \
  configuration:v1 \
  --go-header-file ${SCRIPT_ROOT}/hack/boilerplate/boilerplate.go.txt
