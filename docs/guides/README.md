# Guides

Follow one of the guides to learn more about how to use
Kong Ingress controller:

- [Getting started](getting-started.md) with Kong Ingress Controller
- [Using KongPlugin resource](using-kongplugin-resource.md)  
  This guide walks through setting up plugins in Kong using a declarative
  approach.
- [Using KongIngress resource](using-kongingress-resource.md)  
  This guide explains how the KongIngress resource can be used to change Kong
  specific settings like load-balancing, health-checking and proxy behaviour.
- [Using KongConsumer and KongCredential resources](using-consumer-credential-resource.md)  
  This guide walks through how Kubernetes native declarative configuration
  can be used to dynamically provision credentials for authentication purposes
  in the Ingress layer.
- [Using JWT and ACL KongPlugin resources](configure-acl-plugin.md)  
  This guides walks you through configuring the JWT plugin and ACL plugin for
  authentication purposes at the Ingress layer
- [Using cert-manager with Kong](cert-manager.md)  
  This guide walks through how to use cert-manager along with Kong Ingress
  Controller to automate TLS certificate provisioning and using them
  to encrypt your API traffic.
- [Configuring a fallback service](configuring-fallback-service.md)  
  This guide walks through how to setup a fallback service using Ingress
  resource. The fallback service will receive all requests that don't
  match against any of the defined Ingress rules.
- [Using external service](using-external-service.md)  
  This guide shows how to expose services running outside Kubernetes via Kong,
  using [External Name](https://kubernetes.io/docs/concepts/services-networking/service/#externalname)
  Services in Kubernetes.
- [Configuring HTTPS redirects for your services](configuring-https-redirect.md)  
  This guide walks through how to configure Kong Ingress Controller to
  redirect HTTP request to HTTPS so that all communication
  from the external world to your APIs and microservices is encrypted.
- [Using Redis for rate-limiting](redis-rate-limiting.md)  
  This guide walks through how to use Redis for storing rate-limit information
  in a multi-node Kong deployment.
- [Integrate Kong Ingress Controller with Prometheus/Grafana](prometheus-grafana.md)  
  This guide walks through the steps of how to deploy Kong Ingress Controller
  and Prometheus to obtain metrics for the traffic flowing into your
  Kubernetes cluster.
- [Configuring circuit-breaker and health-checking](configuring-health-checks.md)  
  This guide walks through the usage of Circuit-breaking and health-checking
  features of Kong Ingress Controller.
- [Setting up custom plugin](setting-up-custom-plugins.md)  
  This guide walks through
  installation of a custom plugin into Kong using
  ConfigMaps and Volumes.
- [Using ingress with gRPC](using-ingress-with-grpc.md)  
  This guide walks through how to use Kong Ingress Controller with gRPC.
- [Setting up upstream mTLS](upstream-mtls.md)  
  This guide gives an overview of how to setup mutual TLS authentication
  between Kong and your upstream server.
- [Preserveing Client IP address](preserve-client-ip.md)  
  This guide gives an overview of different methods to preserve the Client
  IP address.
- [Using KongClusterPlugin resource](using-kongclusterplugin-resource.md)  
  This guide walks through setting up plugins that can be shared across
  Kubernetes namespaces.
- [Using Kong with Knative](using-kong-with-knative.md)  
  This guide gives an overview of how to setup Kong as the Ingress point
  for Knative workloads.
- [Exposing TCP-based service](using-tcpingress.md)  
  This guide gives an overview of how to use TCPIngress resource to expose
  non-HTTP based services outside a Kubernetes cluster.
- [Using mtls-auth plugin](using-mtls-auth-plugin.md)  
  This guide gives an overview of how to use `mtls-auth` plugin and CA
  certificates to authenticate requests using client certificates.
- [Configuring custom entities in Kong](configuring-custom-entities.md)  
  This guide gives an overview of how to configure custom entities for
  deployments of Kong Ingress Controller running without a database.
- [Using OpenID-connect plugin](using-oidc-plugin.md)  
  This guide walks through steps necessary to set up OIDC authentication.

