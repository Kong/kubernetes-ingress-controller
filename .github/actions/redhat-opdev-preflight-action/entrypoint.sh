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
require_env INPUT_PASSWORD INPUT_USERNAME INPUT_CERTIFICATIONID INPUT_IMAGE


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
    --certification-project-id="${INPUT_CERTIFICATIONID}"           \
    --pyxis-api-token="${INPUT_APITOKEN}"
