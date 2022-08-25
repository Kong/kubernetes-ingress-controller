| title                        | status      |
|------------------------------|-------------|
| Dynamic Plugin Configuration | provisional |

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
  - [User Stories (Optional)](#user-stories-optional)
    - [Story 1](#story-1)
    - [Story 2](#story-2)
    - [Story 3](#story-3)
    - [Story 4](#story-4)
    - [Story 5](#story-5)
  - [Notes/Constraints/Caveats (Optional)](#notesconstraintscaveats-optional)
    - [Notes](#notes)
      - [Upgrade Strategy](#upgrade-strategy)
      - [Usage Examples](#usage-examples)
    - [Caveats](#caveats)
- [Drawbacks](#drawbacks)
- [Alternatives](#alternatives)
<!-- /toc -->

## Summary

Currently the configuration for the `KongPlugin` and `KongClusterPlugin` resources
is generally static, meaning that it is not possible to derive configuration values from
common configuration resources defined on same kubernetes cluster. This leads to a lot
duplication of hard coded values in numerous places making the deployment process
difficult to manage. While it is possible to load configuration from a secret using
the `configFrom` parameter, the entire configuration must be contained with in the
secret which is effectively the same problem.

It would be advantageous, and more in line with common Kubernetes behavior to be able
to derive individual values of plugin configuration from data already stored on the same
Kubernetes cluster. This proposal aims to define a change to the `KongPlugin` and `KongClusterPlugin`
Custom Resource Definitions to facilitate this behavior.

## Motivation

Allow `KongPlugin` and `KongClusterPlugin` resources to read configuration values from Kubernetes,
in addition to literal values.

Continue to allow `KongPlugin` and `KongClusterPlugin` resource to accept literal input values.

Reduce the amount of duplication of configuration value in the simple case. That is,
when using standard kubernetes objects and with out assumed templating systems.

### Goals

* Provide ability to read individual configuration values from `ConfigMap` resources
* Provide ability to read individual configuration values from `Secret` resources
* Provide ability to read individual configuration values from arbitrary or generic resource types (other custom resource types)
* KongPlugin controller should react to updates to ConfigMap and Secret objects referenced by managed plugins to re-apply the configuration to all attached KongPlugins
* Support configuration as an array of objects (`name`, [`value`, `valueFrom`]) similar to pod [envvar core][]
* Support value interpretation in the cases where literal `value` is used. (see [envvar core][])
  * References `$(config_name)` are expanded using previously defined configuration values
  * If a variable cannot be resolved, the input string should be unchanged.
  * Double `$$` are reduced to a single `$`, which allows for escaping ( `$$(foobar)` -> `$(foobar)`)

### Non-Goals

* We want to avoid introducing any new resource kinds, instead keeping with `KongPlugin`

## Proposal

In an effort to stay in line with current Kubernetes resource patterns it would be desirable
to mimic the way pods define environment variables using an array of objects where
the objects contain a `name` and one of `value` or `valueFrom`. `valueFrom`
would further supports `configMapKeyRef`, `secretKeyRef` or `genericKeyRef`

### User Stories (Optional)

#### Story 1

As a operator I want certain plugin configuration values to be dynamic in addition
to statically definable

#### Story 2

As a operator I want to be able to derive a literal value using value interpolation referencing
other values

#### Story 3

As an operator I want to be able to store common plugin options and configurations as
`ConfigMaps` or `Secrets` such that I can write common configurations once and not
have to re-write them across multiple plugins.

#### Story 4

As an operator I expect changes in `ConfigMap` or `Secret` data values to be reflected
in the resource and deployments that depend on them with out having to restart critical parts of the system

#### Story 5

As an operator I expect to have the option to provide default values for plugin configuration which
will be used in the case that a referent object is not available or readable.

### Notes/Constraints/Caveats (Optional)

Using the data defined from the existing `config` values first and utilizing the
`configFrom` values may provide a simple transition forward. This would allow existing
configuration to continue to work as it currently does and values defined by `configFrom`
override. Combined with a feature gate, it may provide means to move between styles

#### Notes

##### Upgrade Strategy

A feasible upgrade path would be to leverage conversion webhooks to translate existing
configuration formats to the proposed format. Translating back may be more complicated if
all of the configurations options are not configured to reside in the same `ConfigMap`
or `Secret`

##### Usage Examples
> `configFrom` is used as an example name in lieu of a better name

**Literal Value**

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: plugin-one
config:
  bar: hello-world
configFrom:
  - name: foo
    value: 100

  - name: bar
    value: goodbye-world # override value found in config.bar

  - name: widget
    value:
      object:
        values: supported
```

**ConfigMap Key Reference**

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: plugin-one
config:
  bar: hello-world
configFrom:
  - name: foo
    valueFrom:
      configMapKeyRef:
        name: my-config-map
        key: foo
        default: 100
        type: number
```

**Secret Key References**

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: plugin-one
config:
  bar: hello-world
configFrom:
  - name: FOO
    valueFrom:
      secretKeyRef:
        name: my-secret
        key: foo
        optional: true
        type: number
```

Pods also allow for previously read values to be interpolated to form new literal values.
This behavior is desirable

**Value Interpolation**

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: plugin-one
config:
  bar: hello-world
configFrom:
  - name: DB_USERNAME
    valueFrom:
      secretKeyRef:
        name: my-secret
        key: username

  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: my-secret
        key: password

  - name: DB_OPTIONS
    valueFrom:
      secretKeyRef:
        name: my-config-map
        key: options

  - name: DB_URL
    value: db://$(DB_USERNAME):$(DB_PASSWORD)@dbhost?$(DB_OPTIONS)
```

#### Caveats

* Many plugins have complex data structures. It is unclear as to the best way to handle defaults or overriding values. Blind override? strategic merge?
* * `configMapKeyRef` and `secretKeyRef` can reference resource in namespaces that differ from the namespace the plugin, or deployment live. It should be expected that, if not specified, the same namespace the plugin is defined in is used.

## Drawbacks

* Support for backwards compatibility between the current implementation, and this style
may not be entirely possible.
* Supporting both the existing literal syntax and the proposed syntax at the same time may be too difficult to maintain
* Dynamic resolution of values can place a fair amount of stress on kubernetes API servers.

## Alternatives

* It is additionally possible to maintain the existing `config` and `configFrom` keys as they exist
and introduce an additional key that supports the behaviors outlines in this proposal
* Introduce a v2 version of the plugin resources utilizing a conversion webhook to transparently migrate between them

[kubernetes.io]: https://kubernetes.io/
[kubernetes/enhancements]: https://git.k8s.io/enhancements
[kubernetes/kubernetes]: https://git.k8s.io/kubernetes
[kubernetes/website]: https://git.k8s.io/website
[envvar core]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#envvar-v1-core
