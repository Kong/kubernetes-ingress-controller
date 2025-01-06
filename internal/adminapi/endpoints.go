package adminapi

import (
	"context"
	"fmt"
	"strings"

	discoveryv1 "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cfgtypes "github.com/kong/kubernetes-ingress-controller/v3/internal/manager/config/types"
)

// DiscoveredAdminAPI represents an Admin API discovered from a Kubernetes Service.
type DiscoveredAdminAPI struct {
	Address string
	PodRef  k8stypes.NamespacedName
}

type Discoverer struct {
	// portNames is the set of port names that Admin API Service ports will be
	// matched against.
	portNames sets.Set[string]

	// dnsStrategy is the DNS strategy to use when resolving Admin API Service
	// addresses.
	dnsStrategy cfgtypes.DNSStrategy
}

func NewDiscoverer(
	adminAPIPortNames sets.Set[string],
	dnsStrategy cfgtypes.DNSStrategy,
) (*Discoverer, error) {
	if adminAPIPortNames.Len() == 0 {
		return nil, fmt.Errorf("no admin API port names provided")
	}
	if err := dnsStrategy.Validate(); err != nil {
		return nil, fmt.Errorf("invalid dns strategy: %w", err)
	}

	return &Discoverer{
		portNames:   adminAPIPortNames,
		dnsStrategy: dnsStrategy,
	}, nil
}

// GetAdminAPIsForService performs an endpoint lookup, using provided kubeClient
// to list provided Admin API Service EndpointSlices.
// The retrieved EndpointSlices' ports are compared with the provided portNames set.
func (d *Discoverer) GetAdminAPIsForService(
	ctx context.Context,
	kubeClient client.Client,
	service k8stypes.NamespacedName,
) (sets.Set[DiscoveredAdminAPI], error) {
	const (
		defaultEndpointSliceListPagingLimit = 100
	)

	// Get all the EndpointSlices assigned to the provided service.
	labelReq, err := labels.NewRequirement("kubernetes.io/service-name", selection.Equals, []string{service.Name})
	if err != nil {
		return nil, err
	}

	var (
		addresses     = sets.New[DiscoveredAdminAPI]()
		continueToken string
		labelSelector = labels.NewSelector().Add(*labelReq)
	)
	for {
		var endpointsList discoveryv1.EndpointSliceList
		if err := kubeClient.List(ctx, &endpointsList, &client.ListOptions{
			LabelSelector: labelSelector,
			Namespace:     service.Namespace,
			Continue:      continueToken,
			Limit:         defaultEndpointSliceListPagingLimit,
		}); err != nil {
			return nil, err
		}

		for _, es := range endpointsList.Items {
			adminAPI, err := d.AdminAPIsFromEndpointSlice(es)
			if err != nil {
				return nil, err
			}
			addresses = addresses.Union(adminAPI)
		}

		if endpointsList.Continue == "" {
			break
		}
		continueToken = endpointsList.Continue
	}
	return addresses, nil
}

// AdminAPIsFromEndpointSlice returns a list of Admin APIs when given
// an EndpointSlice.
func (d *Discoverer) AdminAPIsFromEndpointSlice(
	endpoints discoveryv1.EndpointSlice,
) (sets.Set[DiscoveredAdminAPI], error) {
	discoveredAdminAPIs := sets.New[DiscoveredAdminAPI]()
	for _, p := range endpoints.Ports {
		if p.Name == nil {
			continue
		}

		if !d.portNames.Has(*p.Name) {
			continue
		}

		var serviceName string
		for _, or := range endpoints.OwnerReferences {
			if or.Kind == "Service" && or.APIVersion == "v1" {
				serviceName = or.Name
				break
			}
		}

		for _, e := range endpoints.Endpoints {
			if e.Conditions.Terminating != nil && *e.Conditions.Terminating {
				continue
			}

			// We do not take into account endpoints that are not backed by a Pod.
			if e.TargetRef == nil || e.TargetRef.Kind != "Pod" {
				continue
			}

			if len(e.Addresses) < 1 {
				continue
			}

			svc := k8stypes.NamespacedName{
				Name:      serviceName,
				Namespace: endpoints.Namespace,
			}

			adminAPI, err := adminAPIFromEndpoint(e, p, svc, d.dnsStrategy, endpoints.AddressType)
			if err != nil {
				return nil, err
			}
			discoveredAdminAPIs = discoveredAdminAPIs.Insert(adminAPI)
		}
	}
	return discoveredAdminAPIs, nil
}

func adminAPIFromEndpoint(
	endpoint discoveryv1.Endpoint,
	port discoveryv1.EndpointPort,
	service k8stypes.NamespacedName,
	dnsStrategy cfgtypes.DNSStrategy,
	addressFamily discoveryv1.AddressType,
) (DiscoveredAdminAPI, error) {
	podNN := k8stypes.NamespacedName{
		Name:      endpoint.TargetRef.Name,
		Namespace: endpoint.TargetRef.Namespace,
	}

	// NOTE: Endpoint's addresses are assumed to be fungible, therefore we pick
	// only the first one.
	// For the context please see the `Endpoint.Addresses` godoc.
	eAddress := endpoint.Addresses[0]

	// NOTE: We assume https below because the referenced Admin API
	// server will live in another Pod/elsewhere so allowing http would
	// not be considered best practice.

	switch dnsStrategy {
	case cfgtypes.ServiceScopedPodDNSStrategy:
		if service.Name == "" {
			return DiscoveredAdminAPI{}, fmt.Errorf(
				"service name is empty for an endpoint with TargetRef %s/%s",
				endpoint.TargetRef.Namespace, endpoint.TargetRef.Name,
			)
		}

		ipAddr := strings.ReplaceAll(eAddress, ".", "-")
		address := fmt.Sprintf("%s.%s.%s.svc", ipAddr, service.Name, service.Namespace)

		return DiscoveredAdminAPI{
			Address: fmt.Sprintf("https://%s:%d", address, *port.Port),
			PodRef:  podNN,
		}, nil

	case cfgtypes.NamespaceScopedPodDNSStrategy:
		ipAddr := strings.ReplaceAll(eAddress, ".", "-")
		address := fmt.Sprintf("%s.%s.pod", ipAddr, service.Namespace)

		return DiscoveredAdminAPI{
			Address: fmt.Sprintf("https://%s:%d", address, *port.Port),
			PodRef:  podNN,
		}, nil

	case cfgtypes.IPDNSStrategy:
		bounded := eAddress
		if addressFamily == discoveryv1.AddressTypeIPv6 {
			bounded = fmt.Sprintf("[%s]", bounded)
		}
		return DiscoveredAdminAPI{
			Address: fmt.Sprintf("https://%s:%d", bounded, *port.Port),
			PodRef:  podNN,
		}, nil

	default:
		return DiscoveredAdminAPI{}, fmt.Errorf("unknown dns strategy: %s", dnsStrategy)
	}
}
