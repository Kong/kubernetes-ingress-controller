# Configuring Custom Entities

This is an **advanced-level** guide for users using custom entities in Kong.
Most users do not need to use this feature.

Kong has in-built extensibility with its plugin architecture.
Plugins in Kong have a `config` property where users can store configuration
for any custom plugin and this suffices in most use cases.
In some use cases, plugins define custom entities to store additional
configuration outside the plugin instance itself.
This guide elaborates on how such custom entities can be used with the Kong
Ingress Controller.

> Note: All entities shipped with Kong are supported by Kong Ingress Controller
out of the box. This guide applies only if you have a custom entity in your
plugin. To check if your plugin contains a custom entity, the source code
will usually contain a `daos.lua` file.
Custom plugins have first-class support in Kong Ingress Controller
via the `KongPlugin` CRD.
Please read [the custom plugin guide](../setting-up-custom-plugins.md) instead
if you are only using Custom plugins.

## Caveats

- The feature discussed in this guide apply for DB-less deployments of Kong.
  The feature is not supported for deployments where Kong is used with a
  database or Kong is used in hybrid mode.
  For these deployments, configure custom entities directly using Kong's Admin
  API.
- Custom entities which have a foreign relation with other core entities in Kong
  are not supported. Only entities which can exist by themselves and then
  be referenced via plugin configuration are supported.

## Creating a JSON representation of the custom entity

In this section, we will learn how to create a JSON representation of
a custom entity.

Suppose you have a custom entity with the following schema in your plugin source:

```lua 
{
  name                = "xkcds",
  primary_key         = { "id" },
  cache_key           = { "name" },
  endpoint_key        = "name",
  fields = {
    { id = typedefs.uuid },
    {
      name = {
        type= "string",
        required = true,
        unique = true,
      },
    },
    {
      url = {
        type = "string",
        required = true,
      },
    },
    { created_at = typedefs.auto_timestamp_s },
    { updated_at = typedefs.auto_timestamp_s },
  },
}
```

An instance of such an entity would look like:

```json
{
  "id": "385def6e-3059-4929-bb12-d205e97284c5",
  "name": "Bobby Drop Tables",
  "url": "https://xkcd.com/327/"
}
```

Multiple instances of such an entity are represented as follows:

```json
{
  "xkcds": [
    {
      "id": "385def6e-3059-4929-bb12-d205e97284c5",
      "name": "bobby_tables",
      "url": "https://xkcd.com/327/"
    },
    {
      "id": "d079a632-ac8d-4a9a-860c-71de82e8fc11",
      "name": "compiling",
      "url": "https://xkcd.com/303/"
    }
  ]
}
```

If you have more than one custom entities that you would like to configure
then you can create other entities by specifying the entity name at the root
level of the JSON as the key and then a JSON array containing the
custom entities as the value of the key.

To configure custom entities in a DB-less instance of Kong,
you first need to create such a JSON representation of your entities.

## Configuring the custom entity secret

Once you have the JSON representation, we need to store the configuration
inside a Kubernetes Secret.
The following command assumes the filename to be `entities.json` but you can
use any other filename as well:

```bash
$ kubectl create secret generic -n kong kong-custom-entities --from-file=config=entities.json
secret/kong-custom-entities created
```

Some things to note:
- The key inside the secret must be `config`. This is not configurable at the
  moment.
- The secret must be accessible by the Ingress Controller. The recommended
  practice here is to install the secret in the same namespace in which Kong
  is running.
 
## Configure the Ingress Controller

Once you have the secret containing the custom entities configured,
you need to instruct the controller to read the secret and sync the custom
entities to Kong.

To do this, you need to add the following environment variable to the
`ingress-ccontroller` container:

```yaml
env:
- name: CONTROLLER_KONG_CUSTOM_ENTITIES_SECRET
  value: kong/kong-custom-entities
```

This value of the environment variable takes the form of `<namespace>/<name>`.
You need to configure this only once.

This instructs the controller to watch the above secret and configure Kong
with any custom entities present inside the secret.
If you change the configuration and update the secret with different entities,
the controller will dynamically fetch the updated secret and configure Kong.

## Verification

You can verify that the custom entity was actually created in Kong's memory
using the `GET /xkcds` (endpoint will differ based on the name of the entity)
on Kong's Admin API.
You can forward traffic from your local machine to the Kong Pod to access it:

```bash
$ kubectl port-forward kong/kong-pod-name 8444:8444
```

and in a separate terminal:

```bash
 $ curl -k https://localhost:8444/<entity-name>
```

## Using the custom entity

You can now use reference the custom entity in any of your custom plugin's
`config` object:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: random-xkcd-header
config:
  xkcds:
  - d079a632-ac8d-4a9a-860c-71de82e8fc11
plugin: xkcd-header
```
