# Debug & Troubleshooting

## Debug

Using the flag `--v=XX` it is possible to increase the level of logging.
In particular:

- `--v=3` shows details about the service, Ingress rule and endpoint changes

## Authentication to the Kubernetes API Server

A number of components are involved in the authentication process and the first step is to narrow
down the source of the problem, namely whether it is a problem with service authentication or with the kubeconfig file.
Both authentications must work:

```text
+-------------+   service          +------------+
|             |   authentication   |            |
+  apiserver  +<-------------------+  ingress   |
|             |                    | controller |
+-------------+                    +------------+

```

## Service authentication

The Ingress controller needs information from API server to configure Kong.
Therefore, authentication is required, which can be achieved in two different ways:

1. **Service Account**: This is recommended,
   because nothing has to be configured.  
   The Ingress controller will use information provided by the system
   to communicate with the API server.  
   See 'Service Account' section for details.
1. **Kubeconfig file**: In some Kubernetes environments
   service accounts are not available.  
   In this case a manual configuration is required.
   The Ingress controller binary can be started with the `--kubeconfig` flag.
   The value of the flag is a path to a file specifying how
   to connect to the API server. Using the `--kubeconfig`
   does not requires the flag `--apiserver-host`.  
   The format of the file is identical to `~/.kube/config`
   which is used by `kubectl` to connect to the API server.  
   See 'kubeconfig' section for details.

## Discovering API-server

Using this flag `--apiserver-host=http://localhost:8080`,
it is possible to specify an unsecured API server or
reach a remote Kubernetes cluster using
[kubectl proxy](https://kubernetes.io/docs/user-guide/kubectl/kubectl_proxy/).
Please do not use this approach in production.

In the diagram below you can see the full authentication flow with all options, starting with the browser
on the lower left hand side.

```text

Kubernetes                                                  Workstation
+---------------------------------------------------+     +------------------+
|                                                   |     |                  |
|  +-----------+   apiserver        +------------+  |     |  +------------+  |
|  |           |   proxy            |            |  |     |  |            |  |
|  | apiserver |                    |  ingress   |  |     |  |  ingress   |  |
|  |           |                    | controller |  |     |  | controller |  |
|  |           |                    |            |  |     |  |            |  |
|  |           |                    |            |  |     |  |            |  |
|  |           |  service account/  |            |  |     |  |            |  |
|  |           |  kubeconfig        |            |  |     |  |            |  |
|  |           +<-------------------+            |  |     |  |            |  |
|  |           |                    |            |  |     |  |            |  |
|  +------+----+      kubeconfig    +------+-----+  |     |  +------+-----+  |
|         |<--------------------------------------------------------|        |
|                                                   |     |                  |
+---------------------------------------------------+     +------------------+
```

## Service Account

If using a service account to connect to the API server, Dashboard expects the file
`/var/run/secrets/kubernetes.io/serviceaccount/token` to be present. It provides a secret
token that is required to authenticate with the API server.

Verify with the following commands:

```shell
# start a container that contains curl
$ kubectl run test --image=tutum/curl -- sleep 10000

# check that container is running
$ kubectl get pods
NAME                   READY     STATUS    RESTARTS   AGE
test-701078429-s5kca   1/1       Running   0          16s

# check if secret exists
$ kubectl exec test-701078429-s5kca ls /var/run/secrets/kubernetes.io/serviceaccount/
ca.crt
namespace
token

# get service IP of master
$ kubectl get services
NAME         CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
kubernetes   10.0.0.1     <none>        443/TCP   1d

# check base connectivity from cluster inside
$ kubectl exec test-701078429-s5kca -- curl -k https://10.0.0.1
Unauthorized

# connect using tokens
$ TOKEN_VALUE=$(kubectl exec test-701078429-s5kca -- cat /var/run/secrets/kubernetes.io/serviceaccount/token)
$ echo $TOKEN_VALUE
eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3Mi....9A
$ kubectl exec test-701078429-s5kca -- curl --cacert /var/run/secrets/kubernetes.io/serviceaccount/ca.crt -H  "Authorization: Bearer $TOKEN_VALUE" https://10.0.0.1
{
  "paths": [
    "/api",
    "/api/v1",
    "/apis",
    "/apis/apps",
    "/apis/apps/v1alpha1",
    "/apis/authentication.k8s.io",
    "/apis/authentication.k8s.io/v1beta1",
    "/apis/authorization.k8s.io",
    "/apis/authorization.k8s.io/v1beta1",
    "/apis/autoscaling",
    "/apis/autoscaling/v1",
    "/apis/batch",
    "/apis/batch/v1",
    "/apis/batch/v2alpha1",
    "/apis/certificates.k8s.io",
    "/apis/certificates.k8s.io/v1alpha1",
    "/apis/extensions",
    "/apis/extensions/v1beta1",
    "/apis/policy",
    "/apis/policy/v1alpha1",
    "/apis/rbac.authorization.k8s.io",
    "/apis/rbac.authorization.k8s.io/v1alpha1",
    "/apis/storage.k8s.io",
    "/apis/storage.k8s.io/v1beta1",
    "/healthz",
    "/healthz/ping",
    "/logs",
    "/metrics",
    "/swaggerapi/",
    "/ui/",
    "/version"
  ]
}
```

If it is not working, there are two possible reasons:

1. The contents of the tokens are invalid.
   Find the secret name with `kubectl get secrets | grep service-account` and
  delete it with `kubectl delete secret <name>`.  
  It will automatically be recreated.
1. You have a non-standard Kubernetes installation
   and the file containing the token may not be present.  

The API server will mount a volume containing this file,
but only if the API server is configured to use
the ServiceAccount admission controller.  
If you experience this error,
verify that your API server is using the ServiceAccount admission controller.  
If you are configuring the API server by hand,
you can set this with the `--admission-control` parameter.  
Please note that you should use other admission controllers as well.
Before configuring this option, please read about admission controllers.

More information:

- [User Guide: Service Accounts](http://kubernetes.io/docs/user-guide/service-accounts/)
- [Cluster Administrator Guide: Managing Service Accounts](http://kubernetes.io/docs/admin/service-accounts-admin/)

## Kubeconfig

If you want to use a kubeconfig file for authentication,
follow the deploy procedure and
add the flag `--kubeconfig=/etc/kubernetes/kubeconfig.yaml` to the deployment
