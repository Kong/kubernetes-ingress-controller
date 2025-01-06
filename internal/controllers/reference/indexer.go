package reference

import (
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/scheme"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
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
func objectKeyFunc(obj client.Object) (string, error) {
	// TypeMeta is necessary to generate the correct key for references, but we can't use the original object.
	// controller-runtime's client provides the same object to both predicates and the admission webhook, and can result
	// in a race condition if this uses the original
	o := obj.DeepCopyObject()
	metaObj, ok := o.(client.Object)
	if !ok {
		return "", fmt.Errorf("could not convert %s/%s back to client Object", obj.GetNamespace(), obj.GetName())
	}
	s, err := scheme.Get()
	if err != nil {
		return "", fmt.Errorf("could not get scheme for %s/%s metadata: %w", obj.GetNamespace(), obj.GetName(), err)
	}
	err = util.PopulateTypeMeta(o, s)
	if err != nil {
		return "", fmt.Errorf("could not populate %s/%s metadata: %w", obj.GetNamespace(), obj.GetName(), err)
	}
	return metaObj.GetObjectKind().GroupVersionKind().String() + "/" +
		metaObj.GetNamespace() + "/" + metaObj.GetName(), nil
}

// CacheIndexers implements a reference cache to store reference relationship between k8s objects
// provided by cache.Indexer. It could do CRUD on reference records when referrer and referent are
// both provided. It can also list reference records by referrer or by referent.
type CacheIndexers struct {
	logger  logr.Logger
	indexer cache.Indexer
}

func NewCacheIndexers(logger logr.Logger) CacheIndexers {
	return CacheIndexers{
		logger: logger,
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
	referrerKey, err := objectKeyFunc(ref.Referrer)
	if err != nil {
		return "", err
	}
	referentKey, err := objectKeyFunc(ref.Referent)
	if err != nil {
		return "", err
	}

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
	referrerKey, err := objectKeyFunc(ref.Referrer)
	if err != nil {
		return []string{}, err
	}
	return []string{referrerKey}, nil
}

// ObjectReferenceIndexerReferent is the index function to index by the referent,
// which returns index in the following format from referent:
// group/version,Kind=kind/namespace/name.
func ObjectReferenceIndexerReferent(obj interface{}) ([]string, error) {
	ref, ok := obj.(*ObjectReference)
	if !ok {
		return nil, ErrTypeNotObjectReference
	}
	referentKey, err := objectKeyFunc(ref.Referent)
	if err != nil {
		return []string{}, err
	}
	return []string{referentKey}, nil
}

// SetObjectReference adds or updates a reference record between referrer and referent in reference cache.
func (c CacheIndexers) SetObjectReference(referrer client.Object, referent client.Object) error {
	c.logger.V(logging.DebugLevel).Info("Set reference relation",
		"referrer_kind", referrer.GetObjectKind().GroupVersionKind().String(),
		"referrer_namespace", referrer.GetNamespace(),
		"referrer_name", referrer.GetName(),
		"referent_kind", referent.GetObjectKind().GroupVersionKind().String(),
		"referent_namespace", referent.GetNamespace(),
		"referent_name", referent.GetName(),
	)
	ref := &ObjectReference{
		Referrer: referrer,
		Referent: referent,
	}
	return c.indexer.Add(ref)
}

// DeleteObjectReference deletes the reference record between referrer and referent from reference cache.
func (c CacheIndexers) DeleteObjectReference(referrer client.Object, referent client.Object) error {
	c.logger.V(logging.DebugLevel).Info("Delete reference relation",
		"referrer_kind", referrer.GetObjectKind().GroupVersionKind().String(),
		"referrer_namespace", referrer.GetNamespace(),
		"referrer_name", referrer.GetName(),
		"referent_kind", referent.GetObjectKind().GroupVersionKind().String(),
		"referent_namespace", referent.GetNamespace(),
		"referent_name", referent.GetName(),
	)
	ref := &ObjectReference{
		Referrer: referrer,
		Referent: referent,
	}
	return c.indexer.Delete(ref)
}

// ObjectReferred returns true if an object is referenced (being the referent)
// in at least one reference record.
func (c CacheIndexers) ObjectReferred(obj client.Object) (bool, error) {
	objKey, err := objectKeyFunc(obj)
	if err != nil {
		return false, err
	}
	refs, err := c.indexer.ByIndex(IndexNameReferent, objKey)
	if err != nil {
		return false, err
	}

	return len(refs) != 0, nil
}

// ListReferencesByReferrer lists all reference records where referrer has the same key
// (GroupVersionKind+NamespacedName, that means the same k8s object).
func (c CacheIndexers) ListReferencesByReferrer(referrer client.Object) ([]*ObjectReference, error) {
	referrerKey, err := objectKeyFunc(referrer)
	if err != nil {
		return nil, err
	}
	refList, err := c.indexer.ByIndex(IndexNameReferrer, referrerKey)
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

// ListReferencesByReferent lists all reference records referring to the same referent.
func (c CacheIndexers) ListReferencesByReferent(referent client.Object) ([]*ObjectReference, error) {
	referentKey, err := objectKeyFunc(referent)
	if err != nil {
		return nil, err
	}
	refList, err := c.indexer.ByIndex(IndexNameReferent, referentKey)
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
	key, err := objectKeyFunc(referrer)
	if err != nil {
		return err
	}
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
func (c CacheIndexers) DeleteObjectIfNotReferred(obj client.Object, dataplaneClient controllers.DataPlaneClient) error {
	referred, err := c.ObjectReferred(obj)
	if err != nil {
		return err
	}
	if !referred {
		c.logger.V(logging.DebugLevel).Info("Delete object from cache because it is no longer referred",
			"kind", obj.GetObjectKind(),
			"namespace", obj.GetNamespace(),
			"name", obj.GetName(),
		)
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

// ListReferrerObjectsByReferent lists all objects that refers to the same referent.
func (c CacheIndexers) ListReferrerObjectsByReferent(referent client.Object) ([]client.Object, error) {
	refs, err := c.ListReferencesByReferent(referent)
	if err != nil {
		return nil, err
	}
	objs := []client.Object{}
	for _, ref := range refs {
		objs = append(objs, ref.Referrer)
	}
	return objs, nil
}
