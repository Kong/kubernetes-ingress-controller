#!/bin/bash

export GO111MODULE=on

VERSION="kubernetes-1.19.0-rc.4"
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
${PACKAGE_NAME}/pkg/client/configuration \
${PACKAGE_NAME}/pkg/apis \
"configuration:v1,v1beta1" \
--go-header-file ${SCRIPT_ROOT}/hack/boilerplate.go.txt
