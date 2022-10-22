# !/bin/bash

# Fail fast.
set -euo pipefail


# Required environment variables can not be empty.
require_env() {
    local env_vars_names=("$@")
    local missing_vars=0

    for name in "${env_vars_names[@]}"; do
        if [[ -z "${!name:-}" ]]; then
            printf "Require non empty environment variable: %s\n" ${name}
            missing_vars=1
        fi
    done

    if [[ ${missing_vars} -gt 0 ]]; then
        exit 1
    fi
}
require_env INPUT_PASSWORD INPUT_USERNAME INPUT_IMAGE INPUT_SUBMIT


# Login.
echo "$INPUT_PASSWORD"  | docker login quay.io -u "${INPUT_USERNAME}" --password-stdin


# Run checks, do not submit the results.
if [[ "${INPUT_SUBMIT}" == "false" ]]; then
    printf "Skipping submission to connect portal"
    docker run                                                          \
        -v "/var/run/docker.sock":"/var/run/docker.sock"                \
        -v "/home/runner/.docker/":"/docker"                            \
        quay.io/opdev/preflight:stable                                  \
        check container "${INPUT_IMAGE}"                                \
        --docker-config=/docker/config.json
    exit 0
fi


# Run checks and submit the results.
printf "Submitting to connect portal"
require_env INPUT_APITOKEN INPUT_CERTIFICATIONID
docker run                                                          \
    -v "/var/run/docker.sock":"/var/run/docker.sock"                \
    -v "/home/runner/.docker/":"/docker"                            \
    quay.io/opdev/preflight:stable                                  \
    check container "${INPUT_IMAGE}"                                \
    --docker-config=/docker/config.json                             \
    --submit                                                        \
    --certification-project-id="${INPUT_CERTIFICATIONID}"           \
    --pyxis-api-token="${INPUT_APITOKEN}"
