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

Using this Kubernetes feature allows us to add additional commands to [kubectl][4] which improves the user experience:

```bash
$ kubectl get kongplugins
NAME                             AGE
add-ratelimiting-to-route        5h
http-svc-consumer-ratelimiting   5h
```

### KongPlugin

This object allows the configuration of Kong plugins in the same way we [add plugins using the admin API][3]

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

[0]: https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/
[1]: https://getkong.org/docs/0.13.x/admin-api/#consumer-object
[2]: https://getkong.org/docs/0.13.x/admin-api/#plugin-object
[3]: https://getkong.org/docs/0.13.x/admin-api/#add-plugin
[4]: https://kubernetes.io/docs/reference/kubectl/overview/
