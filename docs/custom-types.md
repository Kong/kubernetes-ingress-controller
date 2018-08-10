# Custom Resource Definitions

Kong relies on several [Custom Resource Definition object][0] to declare the additional information to Ingress rules and synchronize configuration with the Kong admin API

This new types are:

- kongconsumer
- kongcredential
- kongplugin

Each one of this new object in Kubernetes have a one-to-one relation with a Kong resource:

- [Consumer][1]
- [Plugin][2]
- Credential created in each authentication plugin.

Using this Kubernetes feature allows us to add additional commands to [kubectl][3] which improves the user experience:

```bash
$ kubectl get kongplugins
NAME                             AGE
add-ratelimiting-to-route        5h
http-svc-consumer-ratelimiting   5h
```

### KongPlugin

This object allows the configuration of Kong plugins in the same way we [add plugins using the admin API][4]

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: <object name>
  namespace: <object namespace>
consumerRef: <name of an existing consumer>
disabled: <boolean>
config:
    key: value
```

- The field `consumerRef` implies the plugin will be used for a particular consumer.
- The value of the field must reference an existing consumer in the same namespace.
- When `consumerRef` is empty it implies the plugin is global. This means, all the requests will use the plugin.
- The field `config` contains a list of `key` and `value` required to configure the plugin.
- The field `disabled` allows us to change the state of the plugin in Kong.

**Important:** the validation of the fields is left to the user. Setting invalid fields avoid the plugin configuration.

*Example:*

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: http-svc-consumer-ratelimiting
  namespace: default
config:
  key: value
```

### KongConsumer

Definition:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: <object name>
  namespace: <object namespace>
username: <user name>
customId: <custom ID>
```

*Example:*

To set up a consumer, first we need a plugin and a credential:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: consumer-team-x
username: team-X

---

apiVersion: configuration.konghq.com/v1
kind: KongCredential
metadata:
  name: credential-team-x
consumerRef: consumer-team-x
type: key-auth
config:
  key: 62eb165c070a41d5c1b58d9d3d725ca1

---

apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: http-svc-consumer-ratelimiting
consumerRef: consumer-team-x
config:
  hour: 1000
  limit_by: ip
  second: 100
```

### KongIngress

This option allows us to configure settings from kong related to the [Upstream][5], [Service][6] and [routes][7] that are defined in the Kubernetes Ingress specification. All `KongIngress` objects must be in the same namespace as the Ingress rule using them.

*There are three ways of activating this feature:*
- Create a `KongIngress` object with the same name as the Ingress rule and it will be applied to all paths in the Ingress rule. This is the default object that will be used.

- Use a `KongIngress` object in one or more Ingress rules by using the annotation `configuration.konghq.com: <name>` in the Ingress rule. This `KongIngress` object will be used instead of a `KongIngress` object that has the same name as the Ingress rule

- Create a `KongIngress` object with the same name as a backend service in an Ingress rule. This is applied to the individual backend and will be be used instead of any of the previous `KongIngress` objects. This allows us to have alternate configurations for individual backend services in the Ingress rule. While using a different `KongIngress` object listed above for the other backends in the Ingress rule. 

In the following example our `clients-svc` service requires a different configuration than the rest of the backends. We would create a new `KongIngress` object named `clients-svc` that will be used for that backend service while `products-svc` and `services-svc` will use the `default` object as it is specified in the `configuration.konghq.com` annotation.
```yaml
apiVersion: configuration.konghq.com/v1
kind: KongIngress
metadata:
  name: default
route:
  preserve_host: false
  strip_path: true
---
apiVersion: configuration.konghq.com/v1
kind: KongIngress
metadata:
  name: clients-svc
proxy:
  retries: 1
route:
  preserve_host: true
  strip_path: false
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: example-services
  annotations:
    configuration.konghq.com: default
spec:
  rules:
  - host: "example.com"
    http:
      paths:
      - path: /products
        backend:
          serviceName: products-svc
          servicePort: 80
      - path: /clients
        backend:
          serviceName: clients-svc
          servicePort: 80
      - path: /services
        backend:
          serviceName: services-svc
          servicePort: 80
```

*Note:* Is not required to define the complete object, we can define the `upstream`, `proxy` or `route` sections

Example:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongIngress
metadata:
  name: configuration-demo
upstream:
  hash_on: none
  hash_fallback: none
  healthchecks:
    active:
      concurrency: 10
      healthy:
      http_statuses:
        - 200
        - 302
      interval: 0
      successes: 0
      http_path: "/"
      timeout: 1
      unhealthy:
        http_failures: 0
        http_statuses:
        - 429
        interval: 0
        tcp_failures: 0
        timeouts: 0
    passive:
      healthy:
      http_statuses:
        - 200
      successes: 0
      unhealthy:
        http_failures: 0
        http_statuses:
        - 429
        - 503
        tcp_failures: 0
        timeouts: 0
    slots: 10
proxy:
  path: /
  connect_timeout: 10000
  retries: 10
  read_timeout: 10000
  write_timeout: 10000
route:
  methods:
  - POST
  - GET
  regex_priority: 0
  strip_path: false
  preserve_host: true
```

[0]: https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/
[1]: https://getkong.org/docs/0.13.x/admin-api/#consumer-object
[2]: https://getkong.org/docs/0.13.x/admin-api/#plugin-object
[3]: https://kubernetes.io/docs/reference/kubectl/overview/
[4]: https://getkong.org/docs/0.13.x/admin-api/#add-plugin
[5]: https://getkong.org/docs/0.13.x/admin-api/#upstream-objects
[6]: https://getkong.org/docs/0.13.x/admin-api/#service-object
[7]: https://getkong.org/docs/0.13.x/admin-api/#route-object
