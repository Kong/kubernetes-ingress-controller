package util

import (
	"fmt"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/client-go/discovery"
)

type IngressAPI int

const (
	OtherAPI          IngressAPI = iota
	NetworkingV1      IngressAPI = iota
	NetworkingV1beta1 IngressAPI = iota
	ExtensionsV1beta1 IngressAPI = iota
)

func (ia IngressAPI) String() string {
	switch ia {
	case NetworkingV1:
		return networkingv1.SchemeGroupVersion.String()
	case NetworkingV1beta1:
		return networkingv1beta1.SchemeGroupVersion.String()
	case ExtensionsV1beta1:
		return extensionsv1beta1.SchemeGroupVersion.String()
	case OtherAPI:
		return "unknown API"
	}
	return "unknown API"
}

// ServerHasGVK returns true iff the Kubernetes API server supports the given resource kind at the given group-version.
func ServerHasGVK(client discovery.ServerResourcesInterface, groupVersion, kind string) (bool, error) {
	list, err := client.ServerResourcesForGroupVersion(groupVersion)
	if err != nil {
		return false, err
	}

	for _, elem := range list.APIResources {
		if elem.Kind == kind {
			return true, nil
		}
	}
	return false, nil
}

func NegotiateResourceAPI(client discovery.ServerResourcesInterface, kind string, allowedVersions []IngressAPI,
) (IngressAPI, error) {
	// BUG: If the server does not know the `candidate` group/version at all, ServerHasGVK returns error, which interrupts
	// the search, instead of passing over that `candidate`. This impacts Kubernetes <1.14. Workaround is to remove APIs
	// newer than the newest one supported by apiserver from allowedVersions.
	// Context: https://github.com/Kong/kubernetes-ingress-controller/issues/918
	for _, candidate := range allowedVersions {
		if ok, err := ServerHasGVK(client, candidate.String(), kind); err != nil {
			return OtherAPI, err
		} else if ok {
			return candidate, nil
		}
	}
	return OtherAPI, fmt.Errorf("no suitable API for kind %q found, tried: %v", kind, allowedVersions)
}
