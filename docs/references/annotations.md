# Kong Ingress Controller annotations

Kong Ingress Controller uses some annotations to configure Ingress resources.

It supports the following annotations:

| Annotation name | Description | Guide |
|-----------------|-------------|-------|
| [`kubernetes.io/ingress.class`](#kubernetesioingressclass) | Restrict the Ingress rules that Kong should satisfy. | Coming soon |
| [`plugins.konghq.com`](#pluginskonghqcom) | Run plugins for specific service or Ingress. | [Using KongPlugin resource](../guides/using-kongplugin-resource.md) |
| [`configuration.konghq.com`](#configurationkonghqcom) | Fine grained routing and load-balancing. | [Using KongIngress resource](../guides/using-kongingress-resource.md)|
| [`configuration.konghq.com/protocol`](#configurationkonghqcom/protocol) | Set protocol on a Service. |
| [`configuration.konghq.com/protocols`](#configurationkonghqcom/protocols) | Set protocols on an Ingress. |
| [`ingress.kubernetes.io/service-upstream`](#ingresskubernetesioservice-upstream) | Offload load-balancing to kube-proxy or sidecar. | Coming soon |

## `kubernetes.io/ingress.class`

If you have multiple Ingress controllers in a single cluster,
you can pick one by specifying the `ingress.class`Â annotation.
Following is an example of
creating an Ingress with an annotation:

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: test-1
  annotations:
    kubernetes.io/ingress.class: "gce"
spec:
  rules:
  - host: example.com
    http:
      paths:
      - path: /test1
        backend:
          serviceName: echo
          servicePort: 80
```

This will target the GCE controller, forcing Kong Ingress Controller to ignore it.

On the other hand, an annotation such as

```yaml
metadata:
  name: test-1
  annotations:
    kubernetes.io/ingress.class: "kong"
```

will target Kong Ingress controller, forcing the GCE controller to ignore it.

> Deploying multiple ingress controller and not specifying the
annotation will cause both controllersÂ fighting to satisfy the Ingress
and will lead to unknown behaviour.

The ingress class used by Kong Ingress Controller to filter Ingress
resources can be changed using the `--ingress-class` CLI flag.

```yaml
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

### Multiple unrelated Kong Ingress Controllers

In some deployments, one might use multiple Kong Ingress Controller
in the same Kubernetes cluster
(e.g. one which serves public traffic, one which serves "internal" traffic).
For such deployments, please ensure that in addition to different
`ingress-class`, the `--election-id` is also different.

In such deployments, `kubernetes.io/ingress.class` annotation can be used on the
following custom resources as well:

- KongPlugin: To configure (global) plugins only in one of the Kong clusters.
- KongConsumer: To create different consumers in different Kong clusters.
- KongCredential: To create associated credentials for consumers.

## `plugins.konghq.com`

Kong's power comes from its plugin architecture, where plugins can modify
the request and response or impose certain policies on the requests being
proxied.

With Kong Ingress Controller, plugins can be configured by creating `KongPlugin`
Custom Resources and then associating them with an Ingress, Service or
KongConsumer resources.

Following is an example of how to use the annotation:

```yaml
plugins.konghq.com: high-rate-limit, docs-site-cors
```

Here, `high-rate-limit` and `docs-site-cors`
are the names of the KongPlugin resources which
should be to be applied to the Ingress rules defined in the
Ingress resource on which the annotation is being applied.

This annotation can also be applied to a Service resource in Kubernetes, which
will result in the plugin being executed at Service-level in Kong,
meaning the plugin will be
executed for every request that is proxied, no matter which Route it came from.

This annotation can be applied to a KongConsumer resource, which results in
plugin being executed whenever the specific consumer is accessing any of
the defined APIs.

Please follow the
[Using the KongPlugin resource](../guides/using-kongplugin-resource.md)
guide for details on how this annotation can be used.

## `configuration.konghq.com`

This annotation can associate a KongIngress resource with
an Ingress or a Service resource.
It serves as a way to bridge the gap between a sparse Ingress API in Kubernetes
with fine-grained controlled using the properties of Service, Route
and Upstream entities in Kong.

Please follow the
[Using the KongIngress resource](../guides/using-kongingress-resource.md)
guide for details on how to use this annotation.

## `configuration.konghq.com/protocol`

This annotation sets a protocol — `http`, `https`, `grpc`, or `grpcs` —
on a Service resource. The protocol is used for communication between a 
[Kong Service](https://docs.konghq.com/latest/admin-api/#service-object) and 
a Kubernetes Service, internally in the Kubernetes cluster.

## `configuration.konghq.com/protocols`

This annotation sets a pair of protocols (`http`,`https`) or (`grpc`,`grpcs`)
on an Ingress resource. The protocols are used for communication between the
Ingress point — in this case,
a [Kong Route](https://docs.konghq.com/latest/admin-api/#route-object) — and
the external user or service.

## `ingress.kubernetes.io/service-upstream`

By default, Kong Ingress Controller distributes traffic amongst all the Pods
of a Kubernetes `Service` by forwarding the requests directly to
Pod IP addresses. One can choose the load-balancing strategy to use
by specifying a KongIngress resource.

However, in some use-cases, the load-balancing should be left up
to `kube-proxy`, or a sidecar component in the case of Service Mesh deployments.

Setting this annotation to a Service resource in Kubernetes will configure
Kong Ingress Controller to directly forward
the traffic outbound for this Service
to the IP address of the service (usually the ClusterIP).

`kube-proxy` can then decide how it wants to handle the request and route the
traffic accordingly. If a sidecar intercepts the traffic from the controller,
it can also route traffic as it sees fit in this case.

Following is an example snippet you can use to configure this annotation
on a `Service` resource in Kubernetes, (please note the quotes around `true`):

```yaml
annotations:
  ingress.kubernetes.io/service-upstream: "true"
```

You need Kong Ingress Controller >= 0.6 for this annotation.

