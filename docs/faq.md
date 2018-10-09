# FAQs

## Why is the default `kubernetes.io/ingress.class` is "nginx"?

This is to provide out of the box compatibility
with tools like [kube-lego][kube-lego].

## Why endpoints and not services?

Kong ingress controller does not use
[Services][k8s-service] to route traffic
to the pods. Instead, it uses the Endpoints API
to bypass [kube-proxy][kube-proxy]
to allow Kong features like session affinity and
custom load balancing algorithms.
It also removes overhead
such as conntrack entries for iptables DNAT.

[kube-lego]: https://github.com/jetstack/kube-lego
[k8s-service]: http://kubernetes.io/docs/user-guide/services
[kube-proxy]: http://kubernetes.io/docs/admin/kube-proxy