# FAQs

### Why endpoints and not services?

Kong ingress controller does not use
[Services][k8s-service] to route traffic
to the pods. Instead, it uses the Endpoints API
to bypass [kube-proxy][kube-proxy]
to allow Kong features like session affinity and
custom load balancing algorithms.
It also removes overhead
such as conntrack entries for iptables DNAT.

### Is it possible to create consumers using the Admin API?

From version 0.5.0 onwards, Kong Ingress Controller tags each entity
that it manages inside Kong's database and only manages the entities that
it creates.
This means that if consumers and credentials are created dynamically, they
won't be deleted by the Ingress Controller.

[k8s-service]: http://kubernetes.io/docs/user-guide/services
[kube-proxy]: http://kubernetes.io/docs/admin/kube-proxy
