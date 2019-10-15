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
  This guides walks through how Kubernetes native declarative configuration
  can be used to dynamically provision credentials for authentication purposes
  in the Ingress layer.
- [Using cert-manager with Kong](cert-manager.md)  
  This guide walks through how to use cert-manager along with Kong Ingress
  Controller to automate TLS certificate provisioning and using them
  to encrypt your API traffic.
- [Configuring a fallback service](configuring-fallback-service.md)  
  This guides walks through how to setup a fallback service using Ingress
  resource. The fallback service will receive all requests that don't
  match against any of the defined Ingress rules.
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
- [Setting up custom plugin](setting-up-custom-plugin.md)  
  This guide walks through
  installation of a custom plugin into Kong using
  ConfigMaps and Volumes.
