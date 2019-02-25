# Kong Ingress on Docker-for-mac Edge

## Prerequisites

1. Ensure the Kubernetes cluster is installed

`Preferences->Kubernetes->Enable Kubernetes`

1. Install [httpie][0] which provides the `http` command
 
 `$ brew install httpie`

## Deploy Kong Ingress Controller

3. Deploy Kong Ingress Controller using `kubectl`:

    ```bash
    curl https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/single/all-in-one-postgres.yaml \
      | kubectl create -f -
    ```

    This command creates:

    ```bash

    namespace "kong" created
    customresourcedefinition "kongplugins.configuration.konghq.com" created
    customresourcedefinition "kongconsumers.configuration.konghq.com" created
    customresourcedefinition "kongcredentials.configuration.konghq.com" created
    service "postgres" created
    statefulset "postgres" created
    serviceaccount "kong-serviceaccount" created
    clusterrole "kong-ingress-clusterrole" created
    role "kong-ingress-role" created
    rolebinding "kong-ingress-role-nisa-binding" created
    clusterrolebinding "kong-ingress-clusterrole-nisa-binding" created
    service "kong-ingress-controller" created
    deployment "kong-ingress-controller" created
    service "kong-proxy" created
    deployment "kong" created
    ```

    *Note:* this process could take up to five minutes the first time

## Deploy a dummy application

5. Setup a dummy service to proxy via Ingress rule:

    ```bash
    curl https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/manifests/dummy-application.yaml \
      | kubectl create -n kong -f -
    ```

	This application must be in the same namespace as the ingress controller.

    This application just returns information about the pod and details from the HTTP request.
    This is an example of the output:
    
    Find its NodePort:
    
    ```bash
    kubectl get svc http-svc -n kong
    NAME       TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
    http-svc   NodePort   10.101.46.155   <none>        80:30942/TCP   44m
    ```
    In this case, it's 30942. Connect via a browser, or http
    
    ```bash
    http localhost:30942
    ```
    Which serves something like this (redacted)

    ```console
	Connection: keep-alive
	Content-Type: text/plain
	Date: Thu, 21 Feb 2019 03:30:55 GMT
	Proxy-Connection: keep-alive
	Server: echoserver
	Transfer-Encoding: chunked
	
	Hostname: http-svc-699c679d79-4mz24
	
	Pod Information:
		node name:	docker-desktop
		pod name:	http-svc-699c679d79-4mz24
		pod namespace:	kong
		pod IP:	x.x.x.x
	
	Server values:
		server_version=nginx: 1.13.3 - lua: 10008
	
	Request Information:
		client_address=x.x.x.x
		method=GET
		real path=/
		query=
		request_version=1.1
		request_uri=http://localhost:8080/
	
	Request Headers:
		accept=*/*
		accept-encoding=gzip, deflate
		connection=keep-alive
		host=localhost:30942
		user-agent=HTTPie/1.0.2
	
	Request Body:
		-no body in request-
    ```

6. Create an Ingress rule

    ```bash
    echo "
    apiVersion: extensions/v1beta1
    kind: Ingress
    metadata:
      name: foo-bar
      namespace: kong
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

## Test the http Ingress rule:

Get the kong proxy port - as its a loadbalancer, you need the first in the pair (I'm not sure what happens if port 80 is already allocated)

```bash
$ kubectl get svc kong-proxy  -n kong
NAME         TYPE           CLUSTER-IP    EXTERNAL-IP   PORT(S)                      AGE
kong-proxy   LoadBalancer   x.x.x.x       localhost     80:32222/TCP,443:31223/TCP   179m
```
    
Fetch the page
    
```bash
http http://localhost:80 Host:foo.bar
```

## Setup TLS

Setup an Ingress rule with TLS section (HTTPS)

7. Create a secret with an SSL certificate

    ```bash
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=secure-foo-bar/O=konghq.org"
    ```

    ```bash
    kubectl create secret tls tls-secret --key tls.key --cert tls.crt -n kong
    ```

8. Create an Ingress rule with a TLS section

    ```bash
    echo "
    apiVersion: extensions/v1beta1
    kind: Ingress
    metadata:
      name: secure-foo-bar
      namespace: kong
    spec:
      tls:
      - hosts:
        - secure.foo.bar
        secretName: tls-secret
      rules:
      - host: secure.foo.bar
        http:
          paths:
          - path: /
            backend:
              serviceName: http-svc
              servicePort: 80
    " | kubectl create -f -
    ```

## Test the https Ingress rule:

```bash
http --verify=no https://localhost Host:secure.foo.bar
```

## Access logs

9. Using `kubectl logs` it is possible to
  follow the interaction of the Ingress controller and Kong.

    ```bash
    kubectl logs -n kong --selector="app=ingress-kong" -c ingress-controller
    ```

## Scaling a deployment

10. One of the main features provided by an Ingress controller
    is the ability to react to changes in the Kubernetes cluster.
    This means if we scale a deployment or a pod dies
    we need to update the Kong configuration (a Target in this case)

    ```bash
	$ kubectl get pods  -n kong |grep http-svc
    ```

    To see this we can execute:

    ```bash
    kubectl scale deployment http-svc --replicas=5
    ```

    After a second we should have more targets (pods)

    ```bash
    $ kubectl get pods  -n kong |grep http-svc
    ```

    If we scale down the number of replicas or a pod is
    killed the Kong target list will be updated

## Configure a Kong plugin using annotations in Ingress

11. Create a Kong plugin running:

    ```bash
    echo "
    apiVersion: configuration.konghq.com/v1
    kind: KongPlugin
    metadata:
      name: rl-by-ip
      namespace: kong
    config:
      hour: 100
      limit_by: ip
      second: 10
    plugin: rate-limiting
    " | kubectl create -f -
    ```

12. Check the plugins was created correctly running:

    ```bash
    kubectl get kongplugins -n kong
    NAME                        AGE
    rl-by-ip                    16s
    ```

13. Patch the Ingress to add the annotation for the rate-limiting plugin

    ```bash
    kubectl patch ingress foo-bar \
      -p '{"metadata":{"annotations":{"plugins.konghq.com":"rl-by-ip\n"}}}' -n kong
    ```

14. Check the plugin was created in Kong:

    Get the port for the `kong-ingress-controller`
    
    ```bash
    kubectl get svc kong-ingress-controller -n kong
    ```

    ```bash
    http http://localhost:<port>/plugins
    ```

15. Verify the plugin was configured making a request and
    checking the `X-RateLimit-Limit-*` headers are present in the response

    ```bash
    http --headers http://localhost Host:foo.bar
    .
    .
    .
    X-RateLimit-Limit-hour: 100
    X-RateLimit-Limit-second: 10
    X-RateLimit-Remaining-hour: 91
    X-RateLimit-Remaining-second: 9
    ```

## Configure a Kong plugin using annotations on Kubernetes Services

16. Using annotations we can decide where we want to apply plugins.
    We have two options:

    - Annotations in an Ingress Kong will apply the plugin/s to the Route
    - Annotations in a Kubernetes Service will apply the plugin/s to the Kong Service

    *Note:* Annotations on a Service takes precedence over annotations on an Ingress

    ```bash
    kubectl patch svc http-svc \
      -p '{"metadata":{"annotations":{"plugins.konghq.com": "rl-by-ip\n"}}}' -n kong
    ```

    This change should add a new plugin on the corresponding Kong Service:

    ```bash
    http http://localhost:<port>/plugins
    ```

## Cleanup

17. Clean up all the resources created:

    ```bash
      kubectl delete crd KongPlugin
      kubectl delete crd KongConsumer
      kubectl delete crd KongCredential
      kubectl delete crd KongIngress
      kubectl delete ing foo-bar
      Kubectl delete ing secure-foo-bar
      kubectl delete secret tls-secret
      kubectl delete namespace kong
    ```
[0]: https://httpie.org/