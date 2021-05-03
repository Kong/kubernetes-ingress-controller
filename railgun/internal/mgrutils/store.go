package mgrutils

import (
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
)

var (
	// CacheStores is the global cache for all controllers to store and retrieve Kubernetes objects from cache.
	CacheStores *store.CacheStores
)

func init() {
	newCacheStores := store.NewCacheStores()
	CacheStores = &newCacheStores
}
