package store

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

func (c CacheStores) TakeSnapshot() (CacheStores, error) {
	// Create a fresh CacheStores instance to store the snapshot
	// in the c.takeSnapshot method. It happens here because it's
	// not required to be guarded by a lock.
	snapshot := NewCacheStores()
	listOfStores := c.listAllStores()

	c.l.RLock()
	defer c.l.RUnlock()

	err := takeSnapshot(&snapshot, listOfStores)
	return snapshot, err
}

// takeSnapshot iterates over all stores and add a deep copy of each object to the snapshot.
// It's up to the caller to ensure that the CacheStore from listOfStores has been derived is
// not modified while the snapshot is being taken, also supplying a pointer to the properly
// constructed CacheStore as an argument which when error is nil, will contain the snapshot.
func takeSnapshot(snapshot *CacheStores, listOfStores []cache.Store) error {
	for _, store := range listOfStores {
		for _, item := range store.List() {
			obj, ok := item.(runtime.Object)
			if !ok {
				return fmt.Errorf("expected runtime.Object, got %T", item)
			}

			copiedObj := obj.DeepCopyObject()
			if err := snapshot.Add(copiedObj); err != nil {
				return err
			}
		}
	}
	return nil
}
