# Kong Ingress Controller

## Description

This repository contains a Kubernetes controller built around the [Kubernetes Ingress resource](http://kubernetes.io/docs/user-guide/ingress/)

Learn more about using Ingress on [k8s.io](http://kubernetes.io/docs/user-guide/ingress/)

### What is an Ingress Controller?

Configuring a webserver or loadbalancer is harder than it should be. Most webserver configuration files are very similar. There are some applications that have weird little quirks that tend to throw a wrench in things, but for the most part you can apply the same logic to them and achieve a desired result.

The Ingress resource embodies this idea, and an Ingress controller is meant to handle all the quirks associated with a specific "class" of Ingress.

An Ingress Controller is a daemon, deployed as a Kubernetes Pod, that watches the apiserver's `/ingresses` endpoint for updates to the [Ingress resource](https://kubernetes.io/docs/concepts/services-networking/ingress/). Its job is to satisfy requests for Ingresses.

## Annotation ingress.class

If you have multiple Ingress controllers in a single cluster, you can pick one by specifying the `ingress.class`
annotation, eg creating an Ingress with an annotation like

```yaml
metadata:
  name: foo
  annotations:
    kubernetes.io/ingress.class: "gce"
```

will target the GCE controller, forcing the nginx controller to ignore it, while an annotation like

```yaml
metadata:
  name: foo
  annotations:
    kubernetes.io/ingress.class: "nginx"
```

will target the nginx controller, forcing the GCE controller to ignore it.

__Note__: Deploying multiple ingress controller and not specifying the annotation will result in both controllers fighting to satisfy the Ingress.

### Running multiple ingress controllers

If you're running multiple ingress controllers, or running on a cloud provider that natively handles ingress, you need to specify the annotation `kubernetes.io/ingress.class: "nginx"` in all ingresses that you would like this controller to claim.  This mechanism also provides users the ability to run _multiple_ NGINX ingress controllers (e.g. one which serves public traffic, one which serves "internal" traffic).  When utilizing this functionality the option `--ingress-class` should be changed to a value unique for the cluster within the definition of the replication controller. Here is a partial example:

```
spec:
  template:
     spec:
       containers:
         - name: kong-ingress-internal-controller
           args:
             - /kong-ingress-controller
             - '--election-id=ingress-controller-leader-internal'
             - '--ingress-class=kong-internal'
```

Not specifying the annotation will lead to multiple ingress controllers claiming the same ingress. Specifying a value which does not match the class of any existing ingress controllers will result in all ingress controllers ignoring the ingress.

The use of multiple ingress controllers in a single cluster is supported in Kubernetes versions >= 1.3.

### Limitations

- Ingress rules for TLS require the definition of the field `host`

### Why endpoints and not services

The NGINX ingress controller does not use [Services](http://kubernetes.io/docs/user-guide/services) to route traffic to the pods. Instead it uses the Endpoints API in order to bypass [kube-proxy](http://kubernetes.io/docs/admin/kube-proxy/) to allow NGINX features like session affinity and custom load balancing algorithms. It also removes some overhead, such as conntrack entries for iptables DNAT.
