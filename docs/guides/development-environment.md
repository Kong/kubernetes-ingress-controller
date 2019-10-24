If you want to develop locally, take the following steps to set up your development environment.

1. This guide assumes you are running in GKE.
2. From this directory, run apply the dbless config to get kong running in your k8s cluster.
    `k apply -f dev-config.yaml`
3. Run Kong Ingress Controller locally.
    `bash start.sh kong`
