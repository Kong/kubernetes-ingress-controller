# Kong Ingress Controller Documentation

## Table of contents

- [Concepts](#concepts)
  - [Architecture](#architecture)
  - [Custom Resources](#custom-resources)
  - [Deployment methods](#deployment-methods)
  - [High-availability and scaling](#high-availability-and-scaling)
  - [Security](#security)
- [Guides and Tutorials](#guides-and-tutorials)
- [Configuration reference](#configuration-reference)
- [FAQs](#faqs)
- [Troubleshooting](#troubleshooting)

## Concepts

### Architecture

The [design][design] document explains how Kong Ingress Controller works
inside a Kubernetes cluster and configures Kong to proxy traffic as per
rules defined in the Ingress resources.

### Custom Resources

The Ingress resource in Kubernetes is a fairly narrow and ambiguous API, and
doesn't offer resources to describe the specifics of proxying.
To overcome this limitation, the `KongIngress` Custom resource is used as an
"extension" to the existing Ingress API.

A few custom resources are bundled with Kong Ingress Controller to configure
settings that are specific to Kong and provide fine-grained control over
the proxying behavior.

Please refer to [custom resources][crd] concept document for more details.

### Deployment Methods

Kong Ingress Controller can be deployed in a variety of deployment patterns.
Please refer to the [deployment](concepts/deployment.md) documentation,
which explains all the components
involved and different ways of deploying them based on the use-case.

### High-availability and Scaling

The Kong Ingress Controller is designed to scale with your traffic
and infrastructure.
Please refer to [this document](concepts/ha-and-scaling.md) to understand
failures scenarios, recovery methods, as well as scaling considerations.

### Security

Please refer to [this document](concepts/security.md) to understand the
default security settings and how to further secure the Ingress Controller.

## Guides and Tutorials

Please browse through [guides][guides] to get started or understand how to configure
a specific setting with Kong Ingress Controller.

## Configuration Reference

The configurations in the Kong Ingress Controller can be tweaked using
Custom Resources and annotations.
Please refer to the following documents detailing this process:

- [Custom Resource Definitions](references/custom-resources.md)
- [Annotations](references/annotations.md)
- [CLI arguments](references/cli-arguments.md)

## FAQs

[FAQs][faqs] will help find answers to common problems quickly.
Please feel free to open Pull Requests to contribute to the list.

## Troubleshooting

Please read through our [deployment guide][deployment] for a detailed
understanding of how Ingress Controller is designed and deployed
along alongside Kong.

- [FAQs][faqs] might help as well.
- [Troubleshooting][troubleshooting] guide can help
  resolve some issues.  
  Please contribute back if you feel your experience can help
  the larger community.

[annotations]: annotations.md
[crd]: concepts/custom-resources.md
[deployment]: deployment/
[design]: concepts/design.md
[faqs]: faq.md
[troubleshooting]: troubleshooting.md
[guides]: guides/

[Back to top](#kong-ingress-controller-documentation)
