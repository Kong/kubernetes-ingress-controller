package store

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured/unstructuredscheme"
	"k8s.io/apimachinery/pkg/runtime"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	yamlserializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/scheme"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

//go:generate go run ../../hack/generators/cache-stores

// NewCacheStoresFromObjYAML provides a new CacheStores object given any number of byte arrays containing
// YAML Kubernetes objects. An error is returned if any provided YAML was not a valid Kubernetes object.
func NewCacheStoresFromObjYAML(objs ...[]byte) (c CacheStores, err error) {
	kobjs := make([]runtime.Object, 0, len(objs))
	sr := serializer.NewYAMLSerializer(
		yamlserializer.DefaultMetaFactory,
		unstructuredscheme.NewUnstructuredCreator(),
		unstructuredscheme.NewUnstructuredObjectTyper(),
	)
	for _, yaml := range objs {
		kobj, _, decodeErr := sr.Decode(yaml, nil, nil)
		if err = decodeErr; err != nil {
			return
		}
		if err := util.PopulateTypeMeta(kobj, scheme.Get()); err != nil {
			return c, err
		}
		kobjs = append(kobjs, kobj)
	}
	return NewCacheStoresFromObjs(kobjs...)
}

// NewCacheStoresFromObjs provides a new CacheStores object given any number of Kubernetes
// objects that should be pre-populated. This function will sort objects into the appropriate
// sub-storage (e.g. IngressV1, TCPIngress, e.t.c.) but will produce an error if any of the
// input objects are erroneous or otherwise unusable as Kubernetes objects.
func NewCacheStoresFromObjs(objs ...runtime.Object) (CacheStores, error) {
	c := NewCacheStores()
	for _, obj := range objs {
		typedObj, err := mkObjFromGVK(obj.GetObjectKind().GroupVersionKind())
		if err != nil {
			return c, err
		}

		if err := convUnstructuredObj(obj, typedObj); err != nil {
			return c, err
		}

		if err := c.Add(typedObj); err != nil {
			return c, err
		}
	}
	return c, nil
}

// CacheStoresLockNotInitializedError is returned when the RW lock in the cache stores is nil.
// It indicates that the CacheStores may not be correctly initialized.
type CacheStoresLockNotInitializedError struct{}

var _ error = CacheStoresLockNotInitializedError{}

func (e CacheStoresLockNotInitializedError) Error() string {
	return "lock of cache stores not initialized"
}

// Available returns whether the cache is correctly initialized and available for loading/storing objects.
// When the lock or any of the stores in `nil`, it returns false.
func (c CacheStores) Available() bool {
	if c.l == nil {
		return false
	}
	for _, store := range c.ListAllStores() {
		if store == nil {
			return false
		}
	}
	return true
}
