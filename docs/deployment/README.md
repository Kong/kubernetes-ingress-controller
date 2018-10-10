# Deploying Kong Ingress Controller

Kong ingress controller can be installed on a local or managed
Kubernetes cluster. Here are some guides to get you started:

1. [Using minikube][0]:

   If you have a local Minikube instance running,
   this guide will help you deploy the Ingress Controller.

   *Notes:*
     - This setup does not provide HA for PostgreSQL

1. [Using openshift/minishift][1]:

    Openshift is a Kubernetes distribution by Redhat and
    has few minor differences in how a user logs in using
    `oc` CLI.

   *Notes:*
     - This setup does not provide HA for PostgreSQL
     - Because of CPU/RAM requirements,
       this does not work in OpenShift Online (free account)

1. [Goolge Kubernetes Engine(GKE)][2]:

   [GKE](https://cloud.google.com/kubernetes-engine/)
   is a managed Kubernetes cluster service.
   This guide is a walk through to setup Kong Ingress
   Controller on GKE alongwith TLS certs from
   Let's Encrypt.

[0]: minikube.md
[1]: openshift.md
[2]: gke.md
