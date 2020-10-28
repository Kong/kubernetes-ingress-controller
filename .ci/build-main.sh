#!/bin/bash -e

# build
export IMGNAME="master"
export TAG=`git rev-parse --short HEAD`

# Docker login
REPO="kong-docker-kubernetes-ingress-controller.bintray.io"
echo "${BINTRAY_KEY}" | docker login -u "${BINTRAY_USER}" ${REPO} --password-stdin

# Build and push
make multi-arch 
