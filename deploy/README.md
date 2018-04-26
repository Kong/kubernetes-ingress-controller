This document contains several guides to install the Kong ingress controller in a Kubernetes cluster

1. [Using minikube][0]:

This guide installs:

* PostgreSQL (statefulset)
* RBAC permissions
* Kong deployment in data-plane mode
* Ingress controller deployment with Kong in control-plane mode and init container to run migrations

*Notes:*

- This setup does not provide HA for PostgreSQL

1. [Using openshift/minishift][1]:

This guide installs:

* PostgreSQL (deployment)
* RBAC permissions
* Kong deployment and job to run migrations
* Ingress controller deployment

*Notes:*

- This setup does not provide HA for PostgreSQL
- Because of CPU/RAM requirements, this does not work in OpenShift Online (free account)

[0]: minikube.md
[1]: openshift.md
