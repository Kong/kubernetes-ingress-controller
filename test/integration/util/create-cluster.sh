#!/bin/bash

set -e

REGISTRY_PORT=5000
docker run -d --restart=always -p "$REGISTRY_PORT:5000" --name "$REGISTRY_NAME" registry:2

KIND_CONFIG="
apiVersion: kind.x-k8s.io/v1alpha4
kind: Cluster
name: $CLUSTER_NAME
containerdConfigPatches:
- |-
  [plugins.\"io.containerd.grpc.v1.cri\".registry.mirrors.\"$REGISTRY_NAME:$REGISTRY_PORT\"]
    endpoint = [\"http://${REGISTRY_NAME}:${REGISTRY_PORT}\"]
"

"$KIND_BINARY" create cluster --config=<(echo "$KIND_CONFIG")
docker network connect "kind" "$REGISTRY_NAME"
