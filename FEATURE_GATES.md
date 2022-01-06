# Feature Gates

Upstream [Kubernetes][k8s] includes [Feature Gates][gates] to enable or disable features with flags and track the maturity of a feature using [Feature Stages][stages]. Here in the Kubernetes Ingress Controller (KIC) we use the same definitions of `Feature Gates` and `Feature Stages` from upstream Kubernetes, but with our own list of features.

Using `Feature Gates` enables contributors to add and manage new (and potentially) experimental functionality to the KIC in a controlled manner: the features will be "hidden" until generally available (GA) and the progress and maturity of features on their path to GA will be documented. Feature gates also create a clear path for deprecating features.

See below for current features and their statuses, and follow the links to the relevant feature documentation.

[k8s]:https://kubernetes.io
[gates]:https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/
[stages]:https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/#feature-stages

## Feature gates

Below you will find the overviews of features at various maturity levels:

- [Feature gates for graduated or deprecated features](/#feature-gates-for-graduated-or-deprecated-features)
- [Feature gates for Alpha or Beta features](/#feature-gates-for-alpha-or-beta-features)

Please read the [Important Notes](/#important-notes) section before using any `Alpha` or `Beta` features.

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

### Feature gates for graduated or deprecated features

{{< table caption="Feature Gates for Graduated or Deprecated Features" >}}

| Feature                    | Default | Stage      | Since | Until |
|----------------------------|---------|------------|-------|-------|

{{< /table >}}

Features that reach GA and over time become stable will be removed from this table, they can be found in the main [KIC CRD Documentation][specs] and [Guides][guides].

[specs]:https://docs.konghq.com/kubernetes-ingress-controller/latest/references/custom-resources/
[guides]:https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/overview/

### Feature gates for Alpha or Beta features

{{< table caption="Feature gates for features in Alpha or Beta states" >}}

| Feature | Default | Stage | Since | Until |
|---------|---------|-------|-------|-------|
| Knative | `true`  | Alpha | 0.8.0 | TBD   |
| Gateway | `false` | Alpha | TBD   | TBD   |

{{< /table > }}
