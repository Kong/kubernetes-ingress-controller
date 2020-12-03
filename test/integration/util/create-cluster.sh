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
kubeadmConfigPatches:
- |
  apiVersion: kubeadm.k8s.io/v1beta2
  kind: ClusterConfiguration
  metadata:
    name: config
  apiServer:
    extraArgs:
      "feature-gates": "IPv6DualStack=true"
      "service-cluster-ip-range": "10.96.0.0/16,fd00::/108"
  controllerManager:
    extraArgs:
      "feature-gates": "IPv6DualStack=true"
      "service-cluster-ip-range": "10.96.0.0/16,fd00::/108"
      "cluster-cidr": "10.244.0.0/16,fc00::/48"
      "node-cidr-mask-size": \"0\"
- |
  apiVersion: kubeadm.k8s.io/v1beta2
  kind: InitConfiguration
  metadata:
    name: config
  nodeRegistration:
    kubeletExtraArgs:
      "feature-gates": "IPv6DualStack=true"
- |
  apiVersion: kubeproxy.config.k8s.io/v1alpha1
  kind: KubeProxyConfiguration
  featureGates:
    IPv6DualStack: true
  clusterCIDR: "10.244.0.0/16,fc00::/48"
"

"$KIND_BINARY" create cluster --config=<(echo "$KIND_CONFIG")
docker network connect "kind" "$REGISTRY_NAME"
