# Feature Gates

Upstream [Kubernetes][k8s] includes [Feature Gates][gates] to enable or disable features with flags and track the maturity of a feature using [Feature Stages][stages]. Here in the Kubernetes Ingress Controller (KIC) we use the same definitions of `Feature Gates` and `Feature Stages` from upstream Kubernetes, but with our own list of features.

Using `Feature Gates` enables contributors to add and manage new (and potentially) experimental functionality to the KIC in a controlled manner: the features will be "hidden" until generally available (GA) and the progress and maturity of features on their path to GA will be documented. Feature gates also create a clear path for deprecating features.

See below for current features and their statuses, and follow the links to the relevant feature documentation.

[k8s]:https://kubernetes.io
[gates]:https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/
[stages]:https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/#feature-stages

## Feature gates

Below you will find the overviews of features at various maturity levels:

- [Feature gates for graduated or deprecated features](#feature-gates-for-graduated-or-deprecated-features)
- [Feature gates for Alpha or Beta features](#feature-gates-for-alpha-or-beta-features)

Please read the [Important Notes](#important-notes) section before using any `Alpha` or `Beta` features.

### Important notes

- Most features will be planned and detailed using [Kubernetes Enhancement Proposals (KEP)][k8s-kep]: If you're interested in the development side of features familiarize yourself with our [KEPs][kic-keps]
- The `Since` and `Until` rows in below tables refer to [KIC Releases][releases]
- For `GA` features the documentation exists in the main [Kong Documentation][kong-docs], see the [API reference][api-ref] and [Guides][kic-guides]

An additional **warning** for end-users who are reading this documentation and trying to enable `Alpha` or `Beta` features: it is **very important** to understand that features that are currently in an `Alpha` or `Beta` state may **become `Deprecated` at any time** and **may be removed as part of the next consecutive minor release**. This is especially true for `Alpha` maturity features. In other words, **until a feature becomes GA there are no guarantees that it's going to continue being available**. To avoid disruption to your services engage with the community and read the [CHANGELOG](/CHANGELOG.md) carefully to track progress. Alternatively **do not use features until they have reached a GA status**.

[k8s-keps]:https://github.com/kubernetes/enhancements
[kic-keps]:https://github.com/Kong/kubernetes-ingress-controller/tree/main/keps
[releases]:https://github.com/Kong/kubernetes-ingress-controller/releases
[kong-docs]:https://github.com/Kong/docs.konghq.com
[api-ref]:https://docs.konghq.com/kubernetes-ingress-controller/latest/references/custom-resources/
[kic-guides]:https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/overview/

### Documentation

If you're looking for documentation for `Alpha` maturity features you can find feature preview documentation [here in this repo](/FEATURE_PREVIEW_DOCUMENTATION.md).

Once a feature graduates from `Alpha` to `Beta` maturity these preview docs will be moved to the main [Kong Documentation][kong-docs].

[kong-docs]:https://github.com/kong/docs.konghq.com

### Feature gates for graduated or deprecated features

| Feature          | Default | Stage | Since  | Until |
|------------------|---------|-------|--------|-------|
| Gateway          | `false` | Alpha | 2.2.0  | 2.6.0 |
| Gateway          | `true`  | Beta  | 2.6.0  | 3.0.0 |
| CombinedRoutes   | `false` | Alpha | 2.4.0  | 3.0.0 |
| CombinedRoutes   | `true`  | Beta  | 2.8.0  | 3.0.0 |
| CombinedServices | `false` | Alpha | 2.10.0 | 3.0.0 |
| CombinedServices | `true`  | Beta  | 2.11.0 | 3.0.0 |
| ExpressionRoutes | `false` | Alpha | 2.10.0 | 3.0.0 |
| Knative          | `false` | Alpha | 0.8.0  | 3.0.0 |

Features that reach GA and over time become stable will be removed from this table, they can be found in the main [KIC CRD Documentation][specs] and [Guides][guides].

[specs]:https://docs.konghq.com/kubernetes-ingress-controller/latest/references/custom-resources/
[guides]:https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/overview/

### Feature gates for Alpha or Beta features

| Feature      | Default | Stage | Since  | Until |
|--------------|---------|-------|--------|-------|
| GatewayAlpha | `false` | Alpha | 2.6.0  | TBD   |
| FillIDs      | `false` | Alpha | 2.10.0 | 3.0.0 |
| FillIDs      | `true`  | Beta  | 3.0.0  | TBD   |
| RewriteURIs  | `false` | Alpha | 2.12.0 | TBD   |

**NOTE**: The `Gateway` feature gate refers to [Gateway
 API](https://github.com/kubernetes-sigs/gateway-api) APIs which are in
 `v1beta1` or later. `GatewayAlpha` refers to APIs which are still in alpha.
 These are separated to make a clear distinction in the support stage for these
 APIs.

### Differences between traditional and combined routes

Ingress and HTTPRoute resources use a different approach to configuration layout
than Kong routes. In Ingress and HTTPRoute the upstream service is associated
with individual rules or rule paths, whereas Kong routes associate the upstream
service with the entire route, which may allow multiple different paths.

This difference in layout means that a single Kong route often cannot represent
an entire Ingress or HTTPRoute. For example, the Ingress

```
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example
spec:
  ingressClassName: kong
  rules:
  - host: "ingress.example"
    http:
      paths:
      - path: /one
        pathType: Prefix
        backend:
          service:
            name: red
            port:
              number: 80
      - path: /two
        pathType: Prefix
        backend:
          service:
            name: red
            port:
              number: 80
      - path: /three
        pathType: Prefix
        backend:
          service:
            name: blue
            port:
              number: 80
```

requires more than one Kong route. A single route only supports a single
service, so placing all three paths on the same route would incorrectly direct
traffic destined for the `blue` service to the `red` service.

To account for this, KIC's traditional route generation strategy created a route
for every individual Ingress path. The Ingress above would generate three
routes, for paths `/one`, `/two`, and `/three`. This simple strategy ensures
that requests always route to the correct service, but results in a large number
of routes. Larger configurations can affect performance, especially when loading
configuration updates.

Combined routes reduces configuration size by consolidating routes that share
the same service and hostname. The above example (where all paths use the same
hostname) would instead result in two routes, one matching either `/one` or
`/two` for service `red` and another matching `/three` for service `blue`.

Combined routes thus does not change which requests are routed to which
services. It has no effect on request routing, only the arrangement of
configuration. The number of routes and route names _will_ change, however, and
you should expect to see a disconnect in monitoring information (Prometheus
metrics, logging plugin output, etc.) and third-party tools that rely on route
names or IDs.

Combined routes use a revised naming scheme. Traditional Ingress routes used a
`<namespace>.<name>.<rule index><path index>` name format (e.g.
`default.httpbin.00` for the first (0-indexed) rule's first path on the
`httpbin` Ingress in the `default` namespace), whereas combined routes use a
`<namespace>.<name>.<service>.<hostname>.<port>` scheme (e.g.
`default.httpbin.httpbin.ing.example.80` for all paths in the `default/httpbin`
Ingress for the `httpbin` service on port `80` with the hostname `ing.example`).
HTTPRoutes use the same `httproute.<namespace>.<name>.<rule>.<match>` scheme as
before, but the indices are the _first_ rule and match with a given backendRef,
whereas traditional would generate routes for _every_ match. If rule 1 match 2
has the same backendRef as rule 3 match 1, you'll see a single `.1.2` route with
paths from both.

HTTPRoutes have more combination rules than Ingresses because their rules are
more expressive. Rules cannot be combined with others if they use different
filters, header matches, or query parameter matches, since these are implemented
using Kong settings that apply to an entire route.

HTTPRoute backendRefs can target multiple Services. In traditional mode, KIC
generates a Kong service for every backendRef, labeled with the rule and match
indices. In combined mode, KIC generates a Kong service for every unique set of
services: if two HTTPRoute rules use both serviceA and serviceB in their
backendRefs, KIC will generate a single Kong service with endpoints from both
serviceA and serviceB, named for the first rule and match indices with that
combination of Services.
