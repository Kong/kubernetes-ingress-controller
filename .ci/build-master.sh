#!/bin/bash -e

if [[ "$TRAVIS_BRANCH" != "master" ]];
then
  exit 0
fi

REPO="kong-docker-kubernetes-ingress-controller.bintray.io"
export IMGNAME="master"
export TAG=`git rev-parse --short HEAD`

make container
docker tag ${REPO}/${IMGNAME}:${TAG} ${REPO}/${IMGNAME}:latest
echo "${BINTRAY_KEY}" | docker login -u "${BINTRAY_USER}" ${REPO} --password-stdin

docker push ${REPO}/${IMGNAME}:${TAG}
docker push ${REPO}/${IMGNAME}:latest
