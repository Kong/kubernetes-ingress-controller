# Kong Ingress on Openshift

## Install and Start minishift

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

## Install Kong as ingress controller

1. Download oc CLI from
   [Openshift Console CLI](https://console.starter-us-west-2.openshift.com/console/command-line)

2. Create a new project:

   ```bash
   oc new-project kong-api
   ```

3. Deploy a PostgreSQL database

   ```bash
   oc create --namespace kong-api \
   -f https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/single/postgres-openshift.yaml
   ```

4. Deploy Kong

   You need to execute the next command with `admin` permissions.
   The reason for this is the creation of a role cluster and the required [Custom Resource Definitions](1)

   **Example:** `oc login -u system:admin`

   ```bash
   oc create --namespace kong-api \
   -f https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/single/kong-resources-openshift.yaml
   ```

   *Note:* this process could take up to five minutes.

   **Create environment variables**

   ```bash
   export PROXY_IP=$(minishift ip)
   export PORTS=$(kubectl get service -n kong-api kong-proxy -o go-template='{{range .spec.ports}}{{.nodePort}} {{end}}')
   export HTTP_PORT=$(echo $PORTS  | cut -f 1 -d " ")
   export HTTPS_PORT=$(echo $PORTS | cut -f 2 -d " ")
   ```

   **Note:** using `oc get svc -n kong-api` it is also possible to get information about the ports assigned by Kubernetes.

### Test your deployment

1. Deploy a dummy application:

   ```bash
   curl https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/manifests/dummy-application-openshift.yaml \
   | kubectl create -f -
   ```

   This application just returns information about the pod and details from the HTTP request.
This is an example of the output:

   ```console
   Hostname: http-svc-7dd9588c5-gmbvh

   Pod Information:
     node name:  minishift
     pod name: http-svc-7dd9588c5-gmbvh
     pod namespace:  default
     pod IP: 172.17.0.7

   Server values:
     server_version=nginx: 1.13.3 - lua: 10008

   Request Information:
     client_address=127.0.0.1
     method=GET
     real path=/
     query=
     request_version=1.1
     request_uri=http://localhost:8080/

   Request Headers:
     accept=*/*
     host=localhost:8080
     user-agent=curl/7.47.0

   Request Body:
     -no body in request-
   ```

2. Create an Ingress rule

   ```bash
   echo "
   apiVersion: extensions/v1beta1
   kind: Ingress
   metadata:
     name: foo-bar
   spec:
     rules:
     - host: foo.bar
       http:
         paths:
         - path: /
           backend:
             serviceName: http-svc
             servicePort: 80
   " | kubectl create -f -
   ```

   Test the Ingress rule running:

   ```bash
   http ${PROXY_IP}:${HTTP_PORT} Host:foo.bar
   ```

## Accessing Kong's Admin API

The admin API is exposed using an Ingress to be
able easily to add some authentication plugin and protect the API.

```bash
echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: kong-admin-api
spec:
  rules:
  - host: kong-admin.api
    http:
      paths:
      - path: /
        backend:
          serviceName: kong-admin
          servicePort: 8001
" | kubectl create -f -
```

Using `oc get ingress kong-admin-api` -o yaml you can see the definition of the rule.

To interact with the API you can run

```bash
http ${PROXY_IP}:${HTTP_PORT} Host:kong-admin.api
```

[0]: https://github.com/minishift/minishift
[1]: https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/
