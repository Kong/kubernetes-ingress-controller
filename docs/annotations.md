# Kubernetes annotations

`KongPlugin` and `KongIngress` resources need to be associated with an Ingress resource
for it to take effect, since these resources add value to your routing.

# KongPlugin

## `plugins.konghq.com` Annotation
`KongPlugin` resource can be configured using the `plugins.konghq.com` annotation.
This annotation was introduced in Kong Ingress Controller version 0.2.0.

Following is an example on how to use the annotation:

```yaml
plugins.konghq.com: high-rate-limit, docs-site-cors
```

Here, `high-rate-limit` and `docs-site-cors` are the name of the KongPlugin resources which
need to be applied to the Ingress.

This annotation can be applied to a Service in Kubernetes as well, which
will result in the plugin being executed at Service in Kong, meaning the plugin will be
executed for every request that is proxied, no matter which Route it came from.

## DEPRECATED `<name>.plugin.konghq.com` Annotation

Before version 0.2.0, a different annotation was used to configure plugins,
which is now deprecated.

The annotation can be used as follows:

```yaml
rate-limiting.plugin.konghq.com: |
  add-ratelimiting-to-route
```

The content of the annotation, in this case, `add-ratelimiting-to-route` indicates the name of the `KongPlugin` containing the configuration to be used.

**Rules:**

- the prefix must be a valid plugin name.
- the suffix must be `.plugin.konghq.com`
- the end of the line must be `|` if we want to add multiple plugins.
- each line should contain a valid `KongPlugin` in the Kubernetes cluster.
- `KongPlugin` k8s resources must be unique to each service/ ingress that use any kong plugin

Setting annotations in Ingress rules set ups plugins in `Kong Routes`. Sometimes, we could need to apply plugins in `Kong Services`. To achieve this, we can use the same annotations but applied to the Kubernetes service itself.

*Please check the [Kong 0.13 release notes][1] to learn about Routes and Services*

**Rules:**

[0]: custom-types.md
[1]: https://konghq.com/blog/kong-ce-0-13-0-released/
