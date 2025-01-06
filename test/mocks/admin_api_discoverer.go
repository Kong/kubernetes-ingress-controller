package mocks

import (
	"context"
	"sync/atomic"

	discoveryv1 "k8s.io/api/discovery/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
)

// AdminAPIDiscoverer is a mock implementation of adminapi.Discoverer.
type AdminAPIDiscoverer struct {
	apisToReturn sets.Set[adminapi.DiscoveredAdminAPI]
	errToReturn  error

	getAdminAPIsForServiceCalledTimes     atomic.Int32
	adminAPIsFromEndpointSliceCalledTimes atomic.Int32
}

func NewAdminAPIDiscoverer(apisToReturn sets.Set[adminapi.DiscoveredAdminAPI], errToReturn error) *AdminAPIDiscoverer {
	return &AdminAPIDiscoverer{
		apisToReturn: apisToReturn,
		errToReturn:  errToReturn,
	}
}

func (m *AdminAPIDiscoverer) GetAdminAPIsForService(context.Context, client.Client, k8stypes.NamespacedName) (
	sets.Set[adminapi.DiscoveredAdminAPI],
	error,
) {
	m.getAdminAPIsForServiceCalledTimes.Add(1)
	if m.errToReturn != nil {
		return nil, m.errToReturn
	}
	return m.apisToReturn, nil
}

func (m *AdminAPIDiscoverer) AdminAPIsFromEndpointSlice(discoveryv1.EndpointSlice) (
	sets.Set[adminapi.DiscoveredAdminAPI],
	error,
) {
	m.adminAPIsFromEndpointSliceCalledTimes.Add(1)
	if m.errToReturn != nil {
		return nil, m.errToReturn
	}
	return m.apisToReturn, nil
}

func (m *AdminAPIDiscoverer) GetAdminAPIsForServiceCalledTimes() int {
	return int(m.getAdminAPIsForServiceCalledTimes.Load())
}

func (m *AdminAPIDiscoverer) AdminAPIsFromEndpointSliceCalledTimes() int {
	return int(m.adminAPIsFromEndpointSliceCalledTimes.Load())
}
