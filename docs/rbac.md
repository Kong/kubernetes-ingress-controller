# RBAC permissions

This document outlines the permissions required for
Kong Ingress Controller to function correctly.

Please refer to [design](design.md) documentation to understand how the Ingress Controller
works.

## Configuration

Ingress Controller needs read permissions (get,list,watch)
on the following resources:

- Endpoints
- Nodes
- Pods
- Secrets
- Ingress
- KongPlugins
- KongConsumers
- KongCredentials
- KongIngress

By default, the controller listens for events and above resources across
all namespaces and will need access to these resources at the cluster level
(using ClusterRole and ClusterRoleBinding).

It needs update permission on the status subresource of Ingress as it updates
the it with the loadBalancer IP address once the Ingerss is satisfied.

## Leader election

Kong Ingress Controller, performs a leader-election when multiple
instances are running. This ensures that only one of the controller actually
configures Kong (when running in DB-mode)
and only of the controller actually updates the
Ingress status (when running in DB-less mode).

For this reason, it needs permission to create a ConfigMap.
By default, the permission is given at Cluster level but it can be narrowed
down to a single namespace (using Role and RoleBinding) if needed.

It also needs permission to read and update this ConfigMap.
This permission can be specific to the ConfigMap that is being used
for leader-election purposes.
The name of the ConfigMap is derived from the value of election-id
(default: `ingress-controller-leader`) and
ingress-class (default: `kong`) as: "<election-id>-<ingress-class>".
For example, the default ConfigMap that is used for leader election will
be "ingress-controller-leader-kong", and it will be present in the same
namespace that the controller is deployed in.

Controller also needs permission to read and update permission on this specific
ConfigMap.

[rbac.yaml](../deploy/manifests/rbac.yaml) contains the permissions needed for the Ingress Controller to
operate correctly.
