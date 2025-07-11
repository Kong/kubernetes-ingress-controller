---
title: Configuring Sticky Sessions with Drain Support
description: |
  Learn how to implement sticky sessions with graceful draining of terminating pods using Kong Ingress Controller.
  
breadcrumbs:
  - /kubernetes-ingress-controller/
  - index: kubernetes-ingress-controller
    section: Guides

content_type: guide
layout: guide

products:
  - kic

works_on:
  - on-prem
  - konnect
---

## Overview

Sticky sessions ensure that repeat client requests are routed to the same backend pod, which is essential for maintaining user session state. When combined with drain support, you can implement graceful pod termination during deployments or scaling events, allowing existing sessions to complete while preventing new traffic from being routed to pods that are shutting down.

This guide covers:
1. How sticky sessions work in Kong Ingress Controller
2. How to enable drain support for graceful terminations
3. Configuring sticky sessions with a `KongUpstreamPolicy`
4. Best practices for implementing these features together

## Prerequisites

* A Kubernetes cluster with Kong Ingress Controller installed
* `kubectl` configured to communicate with your cluster
* Basic understanding of Kong Gateway and Kubernetes concepts

## Understanding Sticky Sessions

Kong's sticky sessions feature uses browser-managed cookies to route repeat requests from the same client to the same backend target. When a client first connects, Kong sets a cookie in the response. On subsequent requests, if the cookie is present and valid, Kong routes the client to the same target.

Sticky sessions are useful for:
- Session persistence across multiple requests
- Applications that store session state locally
- Improving cache hit rates

## Understanding Drain Support

Drain support is a feature that allows Kong to gracefully handle terminating pods. When a pod begins the termination process in Kubernetes:

1. The pod is marked for termination but continues running
2. The pod's status changes to `Terminating`
3. With drain support enabled, Kong:
   - Identifies these terminating pods
   - Adds them to the upstream with a weight of 0
   - Allows existing connections to complete
   - Prevents new connections from being routed to these pods

This ensures a smooth transition during deployments, scaling events, or node maintenance.

## Enabling Drain Support

To enable drain support, you must start Kong Ingress Controller with the `--enable-drain-support` flag set to `true`. This can be done in your deployment YAML:

```yaml
containers:
- name: ingress-controller
  image: kong/kubernetes-ingress-controller:2.11
  args:
  - /kong-ingress-controller
  - --enable-drain-support=true
  # other arguments...
```

Alternatively, you can set the corresponding environment variable:

```yaml
env:
- name: CONTROLLER_ENABLE_DRAIN_SUPPORT
  value: "true"
```

## Configuring Sticky Sessions with KongUpstreamPolicy

To implement sticky sessions, you'll need to create a `KongUpstreamPolicy` resource that specifies the `sticky-sessions` algorithm and configure your service to use it.

1. Create a `KongUpstreamPolicy` with sticky sessions:

```yaml
apiVersion: configuration.konghq.com/v1beta1
kind: KongUpstreamPolicy
metadata:
  name: sticky-session-policy
spec:
  algorithm: sticky-sessions
  hashOn:
    input: "none"
  stickySessions:
    cookie: "session-id"
    cookiePath: "/"
```

2. Annotate your service to use this policy:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-app-service
  annotations:
    konghq.com/upstream-policy: sticky-session-policy
spec:
  # service configuration...
```

## Complete Example

Here's a complete example that implements both sticky sessions and drain support:

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example-app
  labels:
    app: example
spec:
  replicas: 3
  selector:
    matchLabels:
      app: example
  template:
    metadata:
      labels:
        app: example
    spec:
      terminationGracePeriodSeconds: 60  # Give pods time to drain
      containers:
        - name: example
          image: my-example-app:latest
          ports:
            - containerPort: 8080
          lifecycle:
            preStop:
              exec:
                command: ["sh", "-c", "sleep 30"]  # Allow time for connections to drain
---
apiVersion: v1
kind: Service
metadata:
  name: example-service
  annotations:
    konghq.com/upstream-policy: sticky-session-policy
spec:
  ports:
    - port: 80
      targetPort: 8080
  selector:
    app: example
---
apiVersion: configuration.konghq.com/v1beta1
kind: KongUpstreamPolicy
metadata:
  name: sticky-session-policy
spec:
  # Use sticky-sessions algorithm for session persistence
  algorithm: sticky-sessions
  
  # Configure cookie settings
  stickySessions:
    cookie: session-id
    cookiePath: "/"
    
  # Optional: Configure health checks
  healthchecks:
    active:
      type: http
      httpPath: /healthz
      concurrency: 10
      healthy:
        interval: 10
        successes: 3
      unhealthy:
        timeouts: 3
        interval: 10
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-ingress
  annotations:
    konghq.com/strip-path: "true"
spec:
  ingressClassName: kong
  rules:
    - http:
        paths:
          - path: /example
            pathType: Prefix
            backend:
              service:
                name: example-service
                port:
                  number: 80
```

## How It Works Together

With the configuration above:

1. When a client first connects to your application, Kong will:
   - Route the request to one of the available pods
   - Set a `session-id` cookie in the response

2. For subsequent requests from the same client:
   - Kong checks for the `session-id` cookie
   - Routes the request to the same backend pod

3. During a deployment or pod termination:
   - Kubernetes marks the pod as `Terminating`
   - Kong Ingress Controller (with drain support enabled) identifies terminating pods
   - These pods remain in the upstream with weight set to 0
   - Existing sessions can complete their work
   - New sessions are routed to healthy pods

## Testing Your Configuration

To test if sticky sessions are working:

1. Make a request to your service and inspect the response headers for the `session-id` cookie
2. Make additional requests and verify they're being routed to the same pod

To test drain support:

1. Scale down your deployment: `kubectl scale deployment example-app --replicas=2`
2. If you have an active session with a pod that's terminating, your session should continue to work
3. New sessions should be directed only to the remaining healthy pods

## Configuration Options

### Sticky Sessions Settings

| Option | Description | Default |
|--------|-------------|---------|
| `cookie` | Name of the cookie used for tracking | - |
| `cookiePath` | Path attribute of the cookie | / |

### Drain Support

Drain support has no additional configuration options beyond enabling it via the controller flag `--enable-drain-support=true`.

## Limitations

- Sticky sessions rely on cookies, which may not be supported in all client environments
- Drain support requires the `--enable-drain-support` flag to be explicitly enabled

## Conclusion

Combining sticky sessions with drain support provides a powerful way to maintain session affinity while ensuring graceful handling of pod terminations. By following the configuration outlined in this guide, you can improve user experience during deployments and scaling events in your Kubernetes environment.
