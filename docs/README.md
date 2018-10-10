# Kong Ingress Controller Documentation

## Design

- Check out our [Design][design] documentation to understand how
  Ingress Controller is
  designed to configure Kong in Kubernetes environment.

## Deployment

Kong Ingress Controller can be deployed on all types of Kubernetes clusters
using a wide variety of deployment options based on your use-case.
Check out the [Deployment Guide][deploymeent]

## Custom Resource Definitions

The controller can configure plugins and other Kong specific features
using Custom Resources. Please refer to our [custom resources][crd] and
[annotations][annotations] documentation for details.

## Annotations

Kong Ingress Controller supports some common annotations and
has specific annotations for Kong specific features like
plugins. Please check out the [annotations][annotations] guide.

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

## Roadmap

- Checkout our [Roadmap][roadmap] for features coming out in future.

[annotations]: annotations.md
[deployement]: deployment/
[crd]: custom-resources.md
[design]: design.md
[faqs]: faq.md
[roadmap]: roadmap.md
[troubleshooting]: troubleshooting.md
