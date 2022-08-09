#!/bin/bash

# Fail fast.
set -euo pipefail

env

# Login.
echo "$INPUT_PASSWORD"  | docker login scan.connect.redhat.com -u "${INPUT_USERNAME}" --password-stdin

# Run checks.
docker run                                                          \
    -v "/var/run/docker.sock":"/var/run/docker.sock"                \
    -v "/home/runner/.docker/":"/docker"                            \
    quay.io/opdev/preflight:stable                                  \
    check container "${INPUT_IMAGE}"                                \
    --docker-config=/docker/config.json                             \
    --submit                                                        \
    --certification-project-id="${INPUT_PROJECTID}"                 \
    --pyxis-api-token="${INPUT_APITOKEN}" --loglevel=debug
