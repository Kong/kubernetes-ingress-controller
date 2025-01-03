package store

import (
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"hash"
	"slices"

	"github.com/samber/lo"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

// TakeSnapshot takes a snapshot of the CacheStores.
func (c CacheStores) TakeSnapshot() (CacheStores, error) {
	if c.l == nil {
		return CacheStores{}, CacheStoresLockNotInitializedError{}
	}
	// Create a fresh CacheStores instance to store the snapshot
	// in the c.takeSnapshot method. It happens here because it's
	// not required to be guarded by a lock.
	snapshot := NewCacheStores()
	listOfStores := c.ListAllStores()

	c.l.RLock()
	defer c.l.RUnlock()

	err := takeSnapshot(&snapshot, listOfStores)
	return snapshot, err
}

// SnapshotHash is type that represents a hash of the snapshot.
// It's a base32 encoded string of the sha256 hash of the snapshot.
// It's printable and human-readable.
type SnapshotHash string

// SnapshotHashEmpty is a constant that represents an empty snapshot hash.
const SnapshotHashEmpty SnapshotHash = ""

func newHashCalculator() hashCalculator {
	return hashCalculator{calc: sha256.New()}
}

type hashCalculator struct {
	calc hash.Hash
}

func (hc hashCalculator) Write(s string) {
	hc.calc.Write([]byte(s))
}

func (hc hashCalculator) Get() SnapshotHash {
	return SnapshotHash(base32.StdEncoding.EncodeToString(hc.calc.Sum(nil)))
}

// TakeSnapshotIfChanged takes a snapshot of the CacheStores if the hash of the current state
// differs from the hash of the previous snapshot supplied as an argument (to make and initial
// just pass empty string). When error is not nil discard all other return values.
// When newHash is empty it means that the snapshot hasn't been taken - returned snapshot is
// meaningless. This is a situation when hash of the current state is the same as the hash of
// the previous snapshot supplied as an argument.
func (c CacheStores) TakeSnapshotIfChanged(previousSnapshotHash SnapshotHash) (
	snapshot CacheStores,
	newHash SnapshotHash,
	err error,
) {
	if c.l == nil {
		return CacheStores{}, "", CacheStoresLockNotInitializedError{}
	}
	// Initialize all variables that don't need to be guarded by a lock.
	snapshot = NewCacheStores()
	listOfStores := c.ListAllStores()
	accessor := meta.NewAccessor()
	hashCalculator := newHashCalculator()

	c.l.RLock()
	defer c.l.RUnlock()

	// Compute the hash of the current store.
	for _, store := range listOfStores {
		// Underlying store is implemented a thread-safe map so for method List() it doesn't maintain order of items.
		// To successfully calculate hash we need to sort the items.
		var capturedErr error
		valuesForHashComputation := lo.Map(store.List(), func(item interface{}, _ int) string {
			obj, ok := item.(runtime.Object)
			if !ok {
				capturedErr = fmt.Errorf("expected runtime.Object, got %T", item)
				return ""
			}
			uid, err := accessor.UID(obj)
			if err != nil {
				capturedErr = fmt.Errorf("failed to get UID: %w", err)
				return ""
			}
			resourceVer, err := accessor.ResourceVersion(obj)
			if err != nil {
				capturedErr = fmt.Errorf("failed to get ResourceVersion: %w", err)
				return ""
			}
			// UID is unique for each object in Kubernetes and ResourceVersion reflects the version of the object.
			return string(uid) + resourceVer
		})
		if capturedErr != nil {
			return CacheStores{}, SnapshotHashEmpty, capturedErr
		}
		// Strings have to be used instead of byte slices, because Cmp.Ordered has to be satisfied.
		slices.Sort(valuesForHashComputation)
		for _, v := range valuesForHashComputation {
			hashCalculator.Write(v)
		}
	}
	// Encode the hash to base32 string to make it human-readable.
	newHash = hashCalculator.Get()

	// If the hash of the current state is the same as the hash of the previous snapshot, return an empty snapshot.
	if newHash == previousSnapshotHash {
		return CacheStores{}, SnapshotHashEmpty, nil
	}

	// Take a snapshot of the current state as the hash of the current state differs from the previous one.
	if err := takeSnapshot(&snapshot, listOfStores); err != nil {
		return CacheStores{}, SnapshotHashEmpty, fmt.Errorf("failed to take snapshot: %w", err)
	}
	return snapshot, newHash, nil
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
