#### Development Environment
If you want to develop the ingress controller locally, take the following steps.

1. This guide assumes you are running in GKE.
2. From this directory, apply the dbless config below to get kong running in your k8s cluster:

    `kubectl apply -f dev-config.yaml`

3. Run Kong Ingress Controller locally:

    `bash start.sh kong`
