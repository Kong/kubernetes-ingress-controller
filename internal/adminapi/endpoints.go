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
)

// DiscoveredAdminAPI represents an Admin API discovered from a Kubernetes Service.
type DiscoveredAdminAPI struct {
	Address string
	PodRef  k8stypes.NamespacedName
}

// GetAdminAPIsForService performs an endpoint lookup, using provided kubeClient
// to list provided Admin API Service EndpointSlices.
// The retrieved EndpointSlices' ports are compared with the provided portNames set.
func GetAdminAPIsForService(
	ctx context.Context,
	kubeClient client.Client,
	service k8stypes.NamespacedName,
	portNames sets.Set[string],
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
			addresses = addresses.Union(AdminAPIsFromEndpointSlice(es, portNames))
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
func AdminAPIsFromEndpointSlice(endpoints discoveryv1.EndpointSlice, portNames sets.Set[string]) sets.Set[DiscoveredAdminAPI] {
	discoveredAdminAPIs := sets.New[DiscoveredAdminAPI]()
	for _, p := range endpoints.Ports {
		if p.Name == nil {
			continue
		}

		if !portNames.Has(*p.Name) {
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
			if e.Conditions.Ready == nil || !*e.Conditions.Ready {
				continue
			}

			// We do not take into account endpoints that are not backed by a Pod.
			if e.TargetRef == nil || e.TargetRef.Kind != "Pod" {
				continue
			}
			podNN := k8stypes.NamespacedName{
				Name:      e.TargetRef.Name,
				Namespace: e.TargetRef.Namespace,
			}

			if len(e.Addresses) < 1 {
				continue
			}

			// Endpoint's addresses are assumed to be fungible, therefore we pick only the first one.
			// For the context please see the `Endpoint.Addresses` godoc.
			addr := strings.ReplaceAll(e.Addresses[0], ".", "-")

			var adminAPI DiscoveredAdminAPI
			// NOTE: We assume https here because the referenced Admin API
			// server will live in another Pod/elsewhere so allowing http would
			// not be considered best practice.
			if serviceName == "" {
				// If we couldn't find a service that's the owner of provided EndpointSlice
				// then fallback to providing a DNS name for the Pod only.
				adminAPI = DiscoveredAdminAPI{
					Address: fmt.Sprintf("https://%s.%s.pod:%d",
						addr, endpoints.Namespace, *p.Port,
					),
					PodRef: podNN,
				}
			} else {
				adminAPI = DiscoveredAdminAPI{
					Address: fmt.Sprintf("https://%s.%s.%s.svc:%d",
						addr, serviceName, endpoints.Namespace, *p.Port,
					),
					PodRef: podNN,
				}
			}
			discoveredAdminAPIs = discoveredAdminAPIs.Insert(adminAPI)
		}
	}
	return discoveredAdminAPIs
}
