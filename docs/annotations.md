# Kubernetes annotations

To configure Kong plugins, credentials and consumers the ingress controller uses annotations to create a mapping between the Ingress and the [custom types][0].
The prefix of the annotation shows which plugin we are trying to set up. For instance, the next code shows we want to configure the `rate-limiting` plugin:

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

Setting annotations in Ingress rules set ups plugins in `Kong Routes`. Sometimes, we could need to apply plugins in `Kong Services`. To achieve this, we can use the same annotations but applied to the Kubernetes service itself.

*Please check the [Kong 0.13 release notes][1] to learn about Routes and Services*

**Rules:**

- If the Ingress and Kubernetes service contains the same annotation, only the defined in the service will be used.
- When there is no overlap of plugins in Ingress and Services annotations both plugins will be configured.

[0]: custom-types.md
[1]: https://konghq.com/blog/kong-ce-0-13-0-released/
