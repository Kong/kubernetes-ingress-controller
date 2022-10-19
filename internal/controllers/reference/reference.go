package reference

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
)

const (
	VersionV1  = "v1"
	KindSecret = "Secret"
)

// UpdateReferencesToSecret update the reference records between referrer and each secret
// in namespacedNames in record cache.
func UpdateReferencesToSecret(
	ctx context.Context,
	c client.Client, indexers CacheIndexers, dataplaneClient *dataplane.KongClient,
	referrer client.Object, namespacedNames []types.NamespacedName,
) error {
	for _, nsName := range namespacedNames {

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: nsName.Namespace,
				Name:      nsName.Name,
			},
		}

		// Here we update the reference relationship even when the referred secret does not exist yet
		// If the referred secret is created, it could be reconciled in secret controller.
		referrerCopy := referrer.DeepCopyObject().(client.Object)
		if err := indexers.SetObjectReference(
			referrerCopy, secret.DeepCopy()); err != nil {
			return err
		}

		if err := c.Get(ctx, nsName, secret); err != nil {
			return err
		}

		if err := dataplaneClient.UpdateObject(secret); err != nil {
			return err
		}
	}

	return removeOutdatedReferencesToSecret(indexers, dataplaneClient, referrer, namespacedNames)
}

// removeOutdatedReferenceToSecret removes outdated reference records to secrets in reference indexer.
// secrets that are really referred by referrer is passed in parameter referredSecretNames.
// if the secret is not actually referrenced by any object after deleting outdated reference records,
// the secret will be removed from object cache in dataplaneClient.
func removeOutdatedReferencesToSecret(
	indexers CacheIndexers, dataplaneClient *dataplane.KongClient,
	referrer client.Object, referredSecretNames []types.NamespacedName,
) error {
	referredSecretNameMap := make(map[types.NamespacedName]bool, len(referredSecretNames))
	for _, nsName := range referredSecretNames {
		referredSecretNameMap[nsName] = true
	}

	referents, err := indexers.ListReferredObjects(referrer)
	if err != nil {
		return err
	}
	for _, obj := range referents {
		gvk := obj.GetObjectKind().GroupVersionKind()
		// delete the reference record if the secret is not referred by the service.
		if gvk.Group == corev1.GroupName && gvk.Version == VersionV1 && gvk.Kind == KindSecret {
			namespacedName := types.NamespacedName{
				Namespace: obj.GetNamespace(),
				Name:      obj.GetName(),
			}
			if !referredSecretNameMap[namespacedName] {
				if err := indexers.DeleteObjectReference(referrer, obj); err != nil {
					return err
				}
				// remove the secret in cache if it is not referred.
				// Do this check and delete when the reference count may be reduced by 1.
				if err := indexers.DeleteObjectIfNotReferred(obj, dataplaneClient); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// DeleteReferencesByReferrer deletes all reference records with specified referrer
// in reference cache.
// If the affected secret is not referred by any other objects, it deletes the secret in object cache.
func DeleteReferencesByReferrer(indexers CacheIndexers, dataplaneClient *dataplane.KongClient, referrer client.Object) error {
	referents, err := indexers.ListReferredObjects(referrer)
	if err != nil {
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
