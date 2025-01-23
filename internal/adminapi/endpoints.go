package adminapi

import (
	"context"
	"fmt"

	discoveryv1 "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DiscoveredAdminAPI represents an Admin API discovered from a Kubernetes Service.
// For field Address use format https://<podIP>:<port>, and for field TLSServerName
// use format <podIPDashes>.<serviceName>.<namespace>.svc, where podIPDashes
// is the IP address separated by dashes instead of dots.
type DiscoveredAdminAPI struct {
	// Address format is https://10.68.0.5:8444.
	Address string
	// TLSServerName format is pod.dataplane-admin-kong-rqwr9-sc49t.default.svc.
	TLSServerName string
	// PodRef is the reference to the Pod with the above IP address.
	PodRef k8stypes.NamespacedName
}

type Discoverer struct {
	// portNames is the set of port names that Admin API Service ports will be
	// matched against.
	portNames sets.Set[string]
}

func NewDiscoverer(
	adminAPIPortNames sets.Set[string],
) (*Discoverer, error) {
	if adminAPIPortNames.Len() == 0 {
		return nil, fmt.Errorf("no admin API port names provided")
	}

	return &Discoverer{
		portNames: adminAPIPortNames,
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

			adminAPI, err := adminAPIFromEndpoint(e, p, svc, endpoints.AddressType)
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
	addressFamily discoveryv1.AddressType,
) (DiscoveredAdminAPI, error) {
	podNN := k8stypes.NamespacedName{
		Name:      endpoint.TargetRef.Name,
		Namespace: endpoint.TargetRef.Namespace,
	}

	// NOTE: Endpoint's addresses are assumed to be fungible, therefore we pick
	// only the first one.
	// For the context please see the `Endpoint.Addresses` godoc.
	podIPAddr := endpoint.Addresses[0]
	if addressFamily == discoveryv1.AddressTypeIPv6 {
		podIPAddr = fmt.Sprintf("[%s]", podIPAddr)
	}

	if service.Name == "" {
		return DiscoveredAdminAPI{}, fmt.Errorf(
			"service name is empty for an endpoint with TargetRef %s/%s",
			endpoint.TargetRef.Namespace, endpoint.TargetRef.Name,
		)
	}

	// NOTE: We assume https below because the referenced Admin API
	// server will live in another Pod/elsewhere so allowing http would
	// not be considered best practice.
	return DiscoveredAdminAPI{
		// Address format:
		//   - ipv4 - https://10.244.0.16:8444
		//   - ipv6 - https://[fd00:10:244::d]:8444
		Address: fmt.Sprintf("https://%s:%d", podIPAddr, *port.Port),
		// TLSServerName format doesn't need to include the IP address part, it's the same for
		// ipv4 and ipv6: pod.dataplane-admin-kong-rqwr9-sc49t.default.svc.
		// Currently everywhere (KGO, Chart) certificates are generated like that
		// *.<service-name>.<namespace>.svc e.g.: *.dataplane-admin-kong-rqwr9-sc49t.default.svc
		// so we have to follow the same 4-parts pattern here, to satisfy wildcard. The first part
		// can be arbitral so let's use "pod". Ditching the first part (wildcard certificate) is
		// problematic, because this requires changes in the certificate generation logic and may
		// break existing users' setups, but it can be done one day.
		TLSServerName: fmt.Sprintf("pod.%s.%s.svc", service.Name, service.Namespace),
		PodRef:        podNN,
	}, nil
}
