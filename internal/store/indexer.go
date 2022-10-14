package store

import (
	"fmt"

	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	IndexNameReferrer = "referrer"
	IndexNameReferent = "referent"
)

type ObjectReference struct {
	Referrer client.Object
	Referent client.Object
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

// ObjectReferenceKeyFunc is the function to trasfer the reference relataioship
// between k8s objects. The key is transferred to the following format:
// group/version,Kind=kind/namespace/name:group/version,Kind=kind/namespace/name
// The part before : is from referrer, and after : is from referent.
func ObjectReferenceKeyFunc(obj interface{}) (string, error) {
	ref, ok := obj.(*ObjectReference)
	if !ok {
		return "", fmt.Errorf("object is type %T, not ObjectReference", obj)
	}
	referrerKey := fmt.Sprintf("%s/%s/%s",
		ref.Referrer.GetObjectKind().GroupVersionKind().String(),
		ref.Referrer.GetNamespace(), ref.Referrer.GetName(),
	)
	referentKey := fmt.Sprintf("%s/%s/%s",
		ref.Referent.GetObjectKind().GroupVersionKind().String(),
		ref.Referent.GetNamespace(), ref.Referent.GetName(),
	)

	return referrerKey + ":" + referentKey, nil
}

// ObjectReferenceIndexerReferrer is the index function to index by the referrer,
// which returns index in the follwig format from referrer:
// group/version,Kind=kind/namespace/name.
func ObjectReferenceIndexerReferrer(obj interface{}) ([]string, error) {
	ref, ok := obj.(*ObjectReference)
	if !ok {
		return nil, fmt.Errorf("object is type %T, not ObjectReference", obj)
	}
	return []string{fmt.Sprintf("%s/%s/%s",
		ref.Referrer.GetObjectKind().GroupVersionKind().String(),
		ref.Referrer.GetNamespace(), ref.Referrer.GetName(),
	)}, nil
}

// ObjectReferenceIndexerReferent is the index function to index by the referent,
// which returns index in the following format from referent:
// group/version,Kind=kind/namespace/name.
func ObjectReferenceIndexerReferent(obj interface{}) ([]string, error) {
	ref, ok := obj.(*ObjectReference)
	if !ok {
		return nil, fmt.Errorf("object is type %T, not ObjectReference", obj)
	}
	return []string{fmt.Sprintf("%s/%s/%s",
		ref.Referent.GetObjectKind().GroupVersionKind().String(),
		ref.Referent.GetNamespace(), ref.Referent.GetName(),
	)}, nil
}

// SetReference stores a reference record in the cache.
func (c CacheIndexers) SetReference(ref *ObjectReference) error {
	return c.indexer.Add(ref)
}

// DeleteReference deletes a reference record (with given key of referrer and referent).
func (c CacheIndexers) DeleteReference(ref *ObjectReference) error {
	return c.indexer.Delete(ref)
}

// ObjectReferred returns true if an object is referenced (being the referent)
// in at least one reference record.
func (c CacheIndexers) ObjectReferred(obj client.Object) (bool, error) {
	key := fmt.Sprintf("%s/%s/%s",
		obj.GetObjectKind().GroupVersionKind().String(),
		obj.GetNamespace(), obj.GetName(),
	)

	refs, err := c.indexer.ByIndex(IndexNameReferent, key)
	if err != nil {
		return false, err
	}

	return len(refs) != 0, nil
}

// ListReferencesByReferrer list all reference records where referrer has the same key
// (GroupVersionKind+NamespacedName, that means the same k8s object).
func (c CacheIndexers) ListReferencesByReferrer(referrer client.Object) ([]*ObjectReference, error) {
	key := fmt.Sprintf("%s/%s/%s",
		referrer.GetObjectKind().GroupVersionKind().String(),
		referrer.GetNamespace(), referrer.GetName(),
	)
	refList, err := c.indexer.ByIndex(IndexNameReferrer, key)
	if err != nil {
		return nil, err
	}
	returnRefList := []*ObjectReference{}
	for _, ref := range refList {
		returnRefList = append(returnRefList, ref.(*ObjectReference))
	}
	return returnRefList, nil
}

// DeleteReferencesByReferrer deletes all reference records where referrer has the same key
// (GroupVersionKind+NamespacedName, that means the same k8s object).
// called when a k8s object deleted in cluster, or when we do not care about it anymore.
func (c CacheIndexers) DeleteReferencesByReferrer(referrer client.Object) error {
	key := fmt.Sprintf("%s/%s/%s",
		referrer.GetObjectKind().GroupVersionKind().String(),
		referrer.GetNamespace(), referrer.GetName(),
	)
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
