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

// key format:
// group/version,Kind=kind/namespace/name:group/version,Kind=kind/namespace/name.
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

func (c CacheIndexers) SetReference(ref *ObjectReference) error {
	return c.indexer.Add(ref)
}

func (c CacheIndexers) RemoveReference(ref *ObjectReference) error {
	return c.indexer.Delete(ref)
}

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
