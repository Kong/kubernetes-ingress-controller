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
- [Using KongConsumer and KongCredential resources](using-consumer-credential-resource)
  This guides walks through how Kubernetes native declarative configuration
  can be used to dynamically provision credentials for authentication purposes
  in the Ingress layer.
- [Using cert-manager with Kong](cert-manager.md)
  This guide walks through how to use cert-manager along with Kong Ingress
  Controller to automate TLS certificate provisioning and using them
  to encrypt your API traffic.
