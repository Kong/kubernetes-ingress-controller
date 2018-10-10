# Kong Ingress controller Roadmap

- Introduce support for upcoming Kong 1.0.
- Use tags in resources created in Kong.
  Kong Ingress Controller will then manage entities created by the controller
  itself and not delete other entities, like the ones created manually.
- Detail CRD specs for self-documenting purpose and also use CRD validations
- Support seamless integration with Kong Enterprise features
  like RBAC, Workspaces.
- **Smarter Sync Logic**: For every event, the Ingress Controller currently
  performs a sync across all entities which is expensive and unnecessary.
  In most cases, controller can make decisions to update a specific entity
  without doing a sync across all entities.
- **Metrics**: Export Kong Ingress Controller performance and error metrics
  in Prometheus Exposition Format for observability.
- **Status sub-resource for k8s > 1.11**: Populate more meta-data in status sub-resource
  for Custom Resources.
- Use finalizers on Custom Resources(CR) to make sure that
  resources on the Kong side are deleted if a CR is deleted.
- Use admission controller to validate custom types and default values
