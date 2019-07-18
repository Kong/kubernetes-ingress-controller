# Deploying Kong Ingress Controller

Kong ingress controller can be installed on a local, managed
or any Kubernetes cluster which supports a service of type `LoadBalancer`.

Here are some guides to get you started:

1. [Using minikube][0]:

   This guide helps you get Kong Ingress Controller on a local
   Kubernetes cluster.

1. [Google Kubernetes Engine(GKE)][2]:

   [GKE](https://cloud.google.com/kubernetes-engine/)
   is a managed Kubernetes cluster offering by Google.
   This guide is a walk through to setup Kong Ingress Controller on GKE.
   If you've access to GKE, please use this guide.

1. [Azure Kubernetes Service(AKS))][3]:

   [AKS](https://azure.microsoft.com/en-us/services/kubernetes-service/)
   is a managed Kubernetes cluster offering by Microsoft Azure.
   This guide is a walk through to setup Kong Ingress
   Controller on AKS.

Once you've Kong Ingress Controlled installed, please follow our
[getting started](../guides/getting-started.md) tutorial to learn
about how to use the Ingress Controller.

[0]: minikube.md
[2]: gke.md
[3]: aks.md
