// Code generated by hack/generators/cache/snapshot.go; DO NOT EDIT.
package store

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

func (c CacheStores) TakeSnapshot() (CacheStores, error) {
	// Create a fresh CacheStores instance to store the snapshot.
	snapshot := NewCacheStores()

	// Gather cache store fields for all kinds.
	allStores := []cache.Store{
		c.IngressV1,
		c.IngressClassV1,
		c.Service,
		c.Secret,
		c.EndpointSlice,
		c.HTTPRoute,
		c.UDPRoute,
		c.TCPRoute,
		c.TLSRoute,
		c.GRPCRoute,
		c.ReferenceGrant,
		c.Gateway,
		c.Plugin,
		c.ClusterPlugin,
		c.Consumer,
		c.ConsumerGroup,
		c.KongIngress,
		c.TCPIngress,
		c.UDPIngress,
		c.KongUpstreamPolicy,
		c.IngressClassParametersV1alpha1,
		c.KongServiceFacade,
		c.KongVault,
	}

	c.l.RLock()
	defer c.l.RUnlock()

	// Iterate over all stores and add a deep copy of each object to the snapshot.
	for _, store := range allStores {
		for _, item := range store.List() {
			obj, ok := item.(runtime.Object)
			if !ok {
				return CacheStores{}, fmt.Errorf("expected runtime.Object, got %T", item)
			}

			copiedObj := obj.DeepCopyObject()
			if err := snapshot.Add(copiedObj); err != nil {
				return CacheStores{}, err
			}
		}
	}

	return snapshot, nil
}
