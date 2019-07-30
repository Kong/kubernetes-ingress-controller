# Security

This document explains the security aspects of Kong Ingress Controller.

Kong Ingress Controller communicates with Kubernetes API-server and Kong's
Admin API. APIs on both sides offer authentication/authorization features
and the controller integrates with them gracefully.

## Kubernetes RBAC

Kong Ingress Controller is deployed with RBAC permissions as explained in
[deployment](deployment.md) document.
It has read and list permissions on most resources but requires update
and create permission for a few resources to provide seamless integration.
The permissions can be locked down further if needed depending on the specific
use-case.
This RBAC policy is associated with a ServiceAccount and the ServiceAccount
is associated with Kong Ingress Controller.
The Controller uses the ServiceAccount credential to authenticate and
authorize itself against the Kubernetes API-server.

## Kong Admin API Protection

Kong's Admin API is used to control configuration of Kong and proxying behavior.
If an attacker happens to gain access to Kong's Admin API, they
will be able to perform all actions as an authorized user like
modifying or deleting Kong's configuration.
Hence, it is important that the deployment
ensures that the likelihood of this happening is as small as possible.

In the example deployements, the Controller and Kong's Admin API communicate
over the loopback (`lo`) interface of the pod.
Kong is no performing any kind of authorization or
authentication on the Admin API, hence the API is accessible only
on the loopback interface to limit the attack surface.
Although not ideal, this setup requires fewer steps
to get started and can be further hardened as required.

Please note that it is very important that Kong's Admin API is not accessible
inside the cluster as any malicious service can change Kong's configuration.
If you're exposing Kong's Admin API itself outside the cluster, please ensure
that you have the necessary authentication in place first.

### Authentication on Kong's Admin API

If Kong's Admin API is protected with one of the authentication plugins,
the Controller can authenticate itself against it to add another layer of
security.
The Controller comes with support for injecting arbitrary HTTP headers
in the requests it makes to Kong's Admin API, which can be used to inject
authentication credentials.
The headers can be specified using the CLI flag `--admin-header` in the Ingress
Controller.

The Ingress Controller will support mutual-TLS-based authentication on Kong's Admin
API in future.

### Kong Enterprise RBAC

Kong Enterprise comes with support for authentication and authorization on
Kong's Admin API.

Once an RBAC token is provisioned, Kong Ingress Controller can use the RBAC
token to authenticate against Kong Enterprise. Use the `--admin-header` CLI
flag to pass the RBAC token the Ingress Controller.
