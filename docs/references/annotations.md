# Kong Ingress Controller annotations

Kong Ingress Controller uses some annotations to configure Ingress resources.

It supports the following annotations:

| Annotation name | Description | Guide |
|-----------------|-------------|-------|
| [`kubernetes.io/ingress.class`](#kubernetesioingressclass) | Restrict the Ingress rules that Kong should satisfy. | TODO |
| [`plugins.konghq.com`](#pluginskonghqcom) | Run plugins for specific service or Ingress. | [Using KongPlugin resource](../guides/using-kongplugin-resource.md) |
| [`configuration.konghq.com`](#configurationkonghqcom) | Fine grained routing and load-balancing. | [Using KongIngress resource](../guides/using-kongingress-resource.md)|

## `kubernetes.io/ingress.class`

If you have multiple Ingress controllers in a single cluster,
you can pick one by specifying the `ingress.class` annotation.
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
annotation will cause both controllers fighting to satisfy the Ingress
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
plugin being executed whenver the specific consumer is accessing any of
the defined APIs.

Please follow the
[Using the KongPlugin resource](../guides/using-kongplugin-resource.md)
guide for details on how this annoatation can be used.

## `configuration.konghq.com`

This annotation can associate a KongIngress resource with
an Ingress or a Service resource.
It serves as a way to bridge the gap between a sparse Ingress API in Kubernetes
with fine-grained controlled using the properties of Service, Route
and Upstream entities in Kong.

Please follow the
[Using the KongIngress resource](../guides/using-kongingress-resource.md)
guide for details on how to use this annotation.
