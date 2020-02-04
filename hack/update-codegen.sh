#!/bin/bash

export GO111MODULE=on

VERSION="kubernetes-1.15.3"
PACKAGE_NAME=github.com/kong/kubernetes-ingress-controller
SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..

if [[ ! -d /tmp/code-generator ]];
then
  git clone https://github.com/kubernetes/code-generator.git  /tmp/code-generator
  pushd /tmp/code-generator
  git checkout $VERSION
  go get ./...
  popd
fi

/tmp/code-generator/generate-groups.sh \
all \
${PACKAGE_NAME}/internal/client/configuration \
${PACKAGE_NAME}/internal/apis \
"configuration:v1,v1beta1" \
--go-header-file ${SCRIPT_ROOT}/hack/boilerplate.go.txt
