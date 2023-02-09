package adminapi

import (
	"context"
	"fmt"

	discoveryv1 "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetURLsForService performs an endpoint lookup, using provided kubeClient
// to list provided Admin API Service EndpointSlices.
func GetURLsForService(ctx context.Context, kubeClient client.Client, service types.NamespacedName) (sets.Set[string], error) {
	const (
		defaultEndpointSliceListPagingLimit = 100
	)

	// Get all the EndpointSlices assigned to the provided service.
	labelReq, err := labels.NewRequirement("kubernetes.io/service-name", selection.Equals, []string{service.Name})
	if err != nil {
		return nil, err
	}

	var (
		addresses     = sets.New[string]()
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
			addresses = addresses.Union(AddressesFromEndpointSlice(es))
		}

		if endpointsList.Continue == "" {
			break
		}
	}
	return addresses, nil
}

// AddressesFromEndpointSlice returns a list of Admin API addresses when given
// an Endpointslice.
func AddressesFromEndpointSlice(endpoints discoveryv1.EndpointSlice) sets.Set[string] {
	addresses := sets.New[string]()
	for _, p := range endpoints.Ports {
		if p.Name == nil {
			continue
		}

		// NOTE: consider making this configurable.
		if *p.Name != "admin" {
			continue
		}

		for _, e := range endpoints.Endpoints {
			if e.Conditions.Ready == nil || !*e.Conditions.Ready {
				continue
			}

			for _, addr := range e.Addresses {
				// NOTE: We assume https here because the referenced Admin API
				// server will live in another Pod/elsewhere so allowing http would
				// not be considered best practice.
				addresses.Insert(fmt.Sprintf("https://%s:%d", addr, *p.Port))
			}
		}
	}
	return addresses
}
