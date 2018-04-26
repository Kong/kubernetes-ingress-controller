1. Install [`minishift`](0)

Minishift is a tool that helps you run OpenShift locally by running a single-node OpenShift cluster inside a VM.

2. Start `minishift`

```bash
minishift start --memory 4096
```

It will take few minutes to get all resources provisioned.

```bash
kubectl get nodes
```

1. Download oc CLI from https://console.starter-us-west-2.openshift.com/console/command-line

2. Create a new project:

```bash
oc new-project kong-api
```

3. Deploy a postgresql database

```bash
oc create --namespace kong-api -f postgres.yaml
```

4. Deploy Kong

You need to execute the next command with `admin` permissions. The reason for this is the creation of a role cluster and the required [Custom Resource Definitions](1) 

**Example:** `oc login -u system:admin`

```bash
oc create --namespace kong-api -f kong.yaml
```

5. Expose Kong Admin API and Proxy ports

```bash
oc expose svc/kong-admin
oc expose svc/kong-proxy
```

[0]: https://github.com/minishift/minishift
[1]: https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/