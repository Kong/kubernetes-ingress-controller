#!/bin/bash -e

# build
export IMGNAME="master"
export TAG=`git rev-parse --short HEAD`

make container

# push
echo "${BINTRAY_KEY}" | docker login -u "${BINTRAY_USER}" ${REPO} --password-stdin
docker push ${REPO}/${IMGNAME}:${TAG}

