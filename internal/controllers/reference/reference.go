package reference

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
)

const (
	VersionV1      = "v1"
	KindSecret     = "Secret"
	KindConfigMap  = "ConfigMap"
	CACertLabelKey = "konghq.com/ca-cert"
)

type secretOrConfigMapT interface {
	client.Object

	*corev1.ConfigMap |
		*corev1.Secret
}

// UpdateReferencesToSecretOrConfigMap updates the reference records between referrer and each secret or configmap
// in namespacedNames in record cache.
func UpdateReferencesToSecretOrConfigMap[t secretOrConfigMapT](
	ctx context.Context,
	c client.Client,
	indexers CacheIndexers,
	dataplaneClient controllers.DataPlaneClient,
	referrer client.Object,
	referencedSecretOrConfigMapNameMap map[k8stypes.NamespacedName]struct{},
	referencedObject t,
) error {
	for nsName := range referencedSecretOrConfigMapNameMap {
		var obj client.Object
		switch (any)(referencedObject).(type) {
		case *corev1.Secret:
			obj = &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: VersionV1,
					Kind:       KindSecret,
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: nsName.Namespace,
					Name:      nsName.Name,
				},
			}
		case *corev1.ConfigMap:
			obj = &corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					APIVersion: VersionV1,
					Kind:       KindConfigMap,
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: nsName.Namespace,
					Name:      nsName.Name,
				},
			}
		}

		// Here we update the reference relationship even when the referred secret or configmap does not exist yet
		// If the referred secret or configmap is created, it could be reconciled in secret or configmap controller.
		referrerCopy := referrer.DeepCopyObject().(client.Object)
		if err := indexers.SetObjectReference(referrerCopy, obj); err != nil {
			return err
		}

		if err := c.Get(ctx, nsName, obj); err != nil {
			return err
		}

		if err := dataplaneClient.UpdateObject(obj); err != nil {
			return err
		}
	}

	return removeOutdatedReferencesToSecretOrConfigMap(ctx, indexers, c, dataplaneClient, referrer, referencedSecretOrConfigMapNameMap)
}

// removeOutdatedReferencesToSecretOrConfigMap removes outdated reference records to secrets or configmaps
// in reference indexer.
// objects that are referred by referrer are passed in referredSecretOrConfigMapNameMap parameter.
// If a secret or a configmap is not referenced by any other object after deleting outdated reference records,
// and it does not have label "konghq.com/ca-cert:true", it is not possible to be used in Kong gateway config
// and should be removed from the object cache inside KongClient.
func removeOutdatedReferencesToSecretOrConfigMap(
	ctx context.Context,
	indexers CacheIndexers,
	c client.Client,
	dataplaneClient controllers.DataPlaneClient,
	referrer client.Object,
	referredSecretOrConfigMapNameMap map[k8stypes.NamespacedName]struct{},
) error {
	referents, err := indexers.ListReferredObjects(referrer)
	if err != nil {
		return err
	}
	for _, obj := range referents {
		gvk := obj.GetObjectKind().GroupVersionKind()
		if gvk.Group != corev1.GroupName || gvk.Version != VersionV1 || (gvk.Kind != KindSecret && gvk.Kind != KindConfigMap) {
			continue
		}

		namespacedName := k8stypes.NamespacedName{
			Namespace: obj.GetNamespace(),
			Name:      obj.GetName(),
		}

		// if the secret or configmap is still referenced, no operations are taken so continue here.
		if _, ok := referredSecretOrConfigMapNameMap[namespacedName]; ok {
			continue
		}

		if err := indexers.DeleteObjectReference(referrer, obj); err != nil {
			return err
		}
		// remove the secret or configmap in object cache if it is not referred and does not have label "konghq.com/ca-cert:true".
		// Do this check and delete when the reference count may be reduced by 1.

		// retrieve the secret or configmap in k8s and check it has the label.
		getErr := c.Get(ctx, namespacedName, obj)
		// if the secret or configmap exists in k8s and has the label, we should not delete it in object cache.
		if getErr == nil {
			if obj.GetLabels() != nil && obj.GetLabels()[CACertLabelKey] == "true" {
				continue
			}
		} else {
			// if the secret or configmap does not exist in k8s, we ignore the error and continue the check and delete operation.
			// for other errors, we return the error and stop the operation.
			if !apierrors.IsNotFound(getErr) {
				return err
			}
		}

		if err := indexers.DeleteObjectIfNotReferred(obj, dataplaneClient); err != nil {
			return err
		}
	}
	return nil
}

// DeleteReferencesByReferrer deletes all reference records with specified referrer
// in reference cache.
// If the affected secret is not referred by any other objects, it deletes the secret in object cache.
func DeleteReferencesByReferrer(indexers CacheIndexers, dataplaneClient controllers.DataPlaneClient, referrer client.Object) error {
	referents, err := indexers.ListReferredObjects(referrer)
	if err != nil {
		indexers.logger.Error(err, "Failed to list referred objects",
			"referrer_kind", referrer.GetObjectKind().GroupVersionKind().String(),
			"referrer_namespace", referrer.GetNamespace(),
			"referrer_name", referrer.GetName(),
		)
		return err
	}

	// delete(gc) the reference record between referrer and referent.
	for _, referent := range referents {
		err := indexers.DeleteObjectReference(referrer, referent)
		if err != nil {
			return err
		}
	}

	// delete the referent in object cache if it is a secret and it is not referenced anymore.
	for _, referent := range referents {
		gvk := referent.GetObjectKind().GroupVersionKind()
		if !(gvk.Group == corev1.GroupName && gvk.Version == VersionV1 && gvk.Kind == KindSecret) {
			continue
		}
		err := indexers.DeleteObjectIfNotReferred(referent, dataplaneClient)
		if err != nil {
			return err
		}
	}

	return nil
}
