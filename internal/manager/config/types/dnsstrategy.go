package types

import "fmt"

// DNSStrategy defines the strategy which KIC will use to create Pod addresses.
type DNSStrategy string

const (
	// IPDNSStrategy defines a strategy where instead of DNS names KIC creates
	// addresses from IP addresses.
	IPDNSStrategy DNSStrategy = "ip"
	// ServiceScopedPodDNSStrategy defines a strategy where KIC creates addresses
	// using the following template:
	// pod-ip-address.service-name.my-namespace.svc.cluster.local
	// Ref: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#a-aaaa-records-1
	//
	// Note: this is known to not work on GKE because it uses kube-dns instead
	// of coredns. GKE docs explicitly mention that:
	// > kube-dns only creates DNS records for Services that have Endpoints.
	//
	// Ref: https://cloud.google.com/kubernetes-engine/docs/how-to/kube-dns#service-dns-records
	ServiceScopedPodDNSStrategy DNSStrategy = "service"
	// NamespaceScopedPodDNSStrategy defines a strategy where KIC creates addresses
	// using the following template:
	// pod-ip-address.my-namespace.pod.cluster-domain.example
	// Ref: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#a-aaaa-records-1
	NamespaceScopedPodDNSStrategy DNSStrategy = "pod"
)

func (d DNSStrategy) Validate() error {
	switch d {
	case IPDNSStrategy:
		return nil
	case ServiceScopedPodDNSStrategy:
		return nil
	case NamespaceScopedPodDNSStrategy:
		return nil
	default:
		return fmt.Errorf("unknown dns strategy: %s", d)
	}
}

func (d DNSStrategy) String() string {
	return string(d)
}
