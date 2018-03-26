This document contains several guides to install the Kong ingress controller in a Kubernetes cluster

1. [Using minikube][0]:

This guide installs:

* PostgresQL (statefulset)
* RBAC permissions
* Kong deployment in data-plane mode
* Ingress controller deployment with Kong in control-plane mode and init container to run migrations

*Notes:*

- This setup do not provides HA for PostgresQL

[0]: minikube.md
