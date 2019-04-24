# Kong Ingress Controller annotations

Kong Ingress Controller uses some annotations to configure Ingress resources.

It supports the following annotations:

- [`kubernetes.io/ingress.class`](#kubernetesioingressclass)
- [`plugins.konghq.com`](#pluginskonghqcom)
- [`configuration.konghq.com`](#configurationkonghqcom)

## `kubernetes.io/ingress.class`

If you have multiple Ingress controllers in a single cluster,
you can pick one by specifying the `ingress.class` annotation.
Following is an example of
creating an Ingress with an annotation:

```yaml
metadata:
  name: foo
  annotations:
    kubernetes.io/ingress.class: "gce"
```

will target the GCE controller, forcing Kong Ingress Controller to ignore it.

On the other hand, an annotation such as

```yaml
metadata:
  name: foo
  annotations:
    kubernetes.io/ingress.class: "kong"
```

will target Kong Ingress controller, forcing the GCE controller to ignore it.

> Deploying multiple ingress controller and not specifying the
annotation will cause both controllers fighting to satisfy the Ingress
and will lead to unknown behaviour.

The ingress class used by Kong Ingress Controller is customizable as well
using the `--ingress-class` flag as follows:

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

In some deployments, one might wish to use multiple Kong clusters in the same
k8s cluster
(e.g. one which serves public traffic, one which serves "internal" traffic).
For such deployments, please ensure that in addition to different
`ingress-class`, the `--election-id` also needs to be different.

In such deployments, `ingress.class` annotation can be used on the
following custom resources as well:

- KongPlugin: To configure (global) plugins only in one of the Kong clusters.
- KongConsumer: To create different consumers in different Kong clusters.
- KongCredential: To create associated credentials for consumers.

## `plugins.konghq.com`

`KongPlugin` custom resource can be configured using the
`plugins.konghq.com` annotation.

Following is an example of how to use the annotation:

```yaml
plugins.konghq.com: high-rate-limit, docs-site-cors
```

Here, `high-rate-limit` and `docs-site-cors`
are the names of the KongPlugin resources which
should be to be applied to the Ingress rules defined in the
Ingress resource on which the annotation is applied.

This annotation can be applied to a Service Object in Kubernetes as well, which
will result in the plugin being executed at Service-level in Kong,
meaning the plugin will be
executed for every request that is proxied, no matter which Route it came from.

See [KongPlugin](#kongplugin) for examples of how to apply a plugin to service
or ingress.

## `configuration.konghq.com`

This annotation can associate a KongIngress custom resource with
an Ingress or a Service resource.
Only a single KongIngress resource can be specified and
it will override the properties of Service, Route and Upstream objects that
are specified in the referenced `KongIngress` object.
Please read [KongIngress](#kongingress) for more details.

[kongplugin]: ../custom-resources.md#KongPlugin
[kongingress]: ../custom-resources.md#KongIngress
