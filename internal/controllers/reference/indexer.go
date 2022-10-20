package reference

import (
	"fmt"

	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
)

const (
	IndexNameReferrer = "referrer"
	IndexNameReferent = "referent"
)

// ErrTypeNotObjectReference is the error returned to caller to tell that type of the object stored
// in indexer is not ObjectReference.
// It should not happen in normal use, because only ObjectReference should be added to the indexer.
var ErrTypeNotObjectReference = fmt.Errorf("type of object in indexer is not ObjectReference")

type ObjectReference struct {
	Referrer client.Object
	Referent client.Object
}

// objectKeyFunc returns a k8s object key in the following format:
// group/version,Kind=kind/namespace/name.
// the combination is unique inside a kubernetes cluster.
func objectKeyFunc(obj client.Object) string {
	return obj.GetObjectKind().GroupVersionKind().String() + "/" +
		obj.GetNamespace() + "/" + obj.GetName()
}

// CacheIndexers implements a reference cache to store reference relationship between k8s objects
// provided by cache.Indexer. It could do CRUD on reference records when referrer and referent are
// both provided. It can also list reference records by referrer or by referent.
type CacheIndexers struct {
	indexer cache.Indexer
}

func NewCacheIndexers() CacheIndexers {
	return CacheIndexers{
		indexer: cache.NewIndexer(ObjectReferenceKeyFunc,
			cache.Indexers{
				IndexNameReferrer: ObjectReferenceIndexerReferrer,
				IndexNameReferent: ObjectReferenceIndexerReferent,
			}),
	}
}

// ObjectReferenceKeyFunc is the function to transfer the reference relataioships
// between k8s objects. The key is transferred to the following format:
// group/version,Kind=kind/namespace/name:group/version,Kind=kind/namespace/name
// The part before : is from referrer, and after : is from referent.
func ObjectReferenceKeyFunc(obj interface{}) (string, error) {
	ref, ok := obj.(*ObjectReference)
	if !ok {
		return "", ErrTypeNotObjectReference
	}
	referrerKey := objectKeyFunc(ref.Referrer)
	referentKey := objectKeyFunc(ref.Referent)

	return referrerKey + ":" + referentKey, nil
}

// ObjectReferenceIndexerReferrer is the index function to index by the referrer,
// which returns the index in the following format from referrer:
// group/version,Kind=kind/namespace/name.
func ObjectReferenceIndexerReferrer(obj interface{}) ([]string, error) {
	ref, ok := obj.(*ObjectReference)
	if !ok {
		return nil, ErrTypeNotObjectReference
	}
	return []string{objectKeyFunc(ref.Referrer)}, nil
}

// ObjectReferenceIndexerReferent is the index function to index by the referent,
// which returns index in the following format from referent:
// group/version,Kind=kind/namespace/name.
func ObjectReferenceIndexerReferent(obj interface{}) ([]string, error) {
	ref, ok := obj.(*ObjectReference)
	if !ok {
		return nil, ErrTypeNotObjectReference
	}
	return []string{objectKeyFunc(ref.Referent)}, nil
}

// SetObjectReference adds or updates a reference record between referrer and referent in reference cache.
func (c CacheIndexers) SetObjectReference(referrer client.Object, referent client.Object) error {
	ref := &ObjectReference{
		Referrer: referrer,
		Referent: referent,
	}
	return c.indexer.Add(ref)
}

// DeleteObjectReference deletes the reference record between referrer and referent from reference cache.
func (c CacheIndexers) DeleteObjectReference(referrer client.Object, referent client.Object) error {
	ref := &ObjectReference{
		Referrer: referrer,
		Referent: referent,
	}
	return c.indexer.Delete(ref)
}

// ObjectReferred returns true if an object is referenced (being the referent)
// in at least one reference record.
func (c CacheIndexers) ObjectReferred(obj client.Object) (bool, error) {
	refs, err := c.indexer.ByIndex(IndexNameReferent, objectKeyFunc(obj))
	if err != nil {
		return false, err
	}

	return len(refs) != 0, nil
}

// ListReferencesByReferrer lists all reference records where referrer has the same key
// (GroupVersionKind+NamespacedName, that means the same k8s object).
func (c CacheIndexers) ListReferencesByReferrer(referrer client.Object) ([]*ObjectReference, error) {
	refList, err := c.indexer.ByIndex(IndexNameReferrer, objectKeyFunc(referrer))
	if err != nil {
		return nil, err
	}
	returnRefList := make([]*ObjectReference, 0, len(refList))
	for _, ref := range refList {
		retRef, ok := ref.(*ObjectReference)
		if !ok {
			return nil, ErrTypeNotObjectReference
		}
		returnRefList = append(returnRefList, retRef)
	}
	return returnRefList, nil
}

// DeleteReferencesByReferrer deletes all reference records where referrer has the same key
// (GroupVersionKind+NamespacedName, that means the same k8s object).
// called when a k8s object deleted in cluster, or when we do not care about it anymore.
func (c CacheIndexers) DeleteReferencesByReferrer(referrer client.Object) error {
	key := objectKeyFunc(referrer)
	refs, err := c.indexer.ByIndex(IndexNameReferrer, key)
	if err != nil {
		return err
	}

	for _, ref := range refs {
		err = c.indexer.Delete(ref)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteObjectIfNotReferred deletes object from object cach by dataplaneClient
// the object is not referenced in reference cache.
func (c CacheIndexers) DeleteObjectIfNotReferred(obj client.Object, dataplaneClient *dataplane.KongClient) error {
	referred, err := c.ObjectReferred(obj)
	if err != nil {
		return err
	}
	if !referred {
		return dataplaneClient.DeleteObject(obj)
	}
	return nil
}

// ListReferredObjects lists all objects referred by referrer in reference cache.
func (c CacheIndexers) ListReferredObjects(referrer client.Object) ([]client.Object, error) {
	refs, err := c.ListReferencesByReferrer(referrer)
	if err != nil {
		return nil, err
	}
	objs := []client.Object{}
	for _, ref := range refs {
		objs = append(objs, ref.Referent)
	}
	return objs, nil
}
