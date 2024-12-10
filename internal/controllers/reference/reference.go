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

// UpdateReferencesToSecret updates the reference records between referrer and each secret
// in namespacedNames in record cache.
func UpdateReferencesToSecret(
	ctx context.Context,
	c client.Client, indexers CacheIndexers, dataplaneClient controllers.DataPlaneClient,
	referrer client.Object, referencedSecretNameMap map[k8stypes.NamespacedName]struct{},
) error {
	for nsName := range referencedSecretNameMap {
		secret := &corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: VersionV1,
				Kind:       KindSecret,
			},
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

	return removeOutdatedReferencesToSecret(ctx, indexers, c, dataplaneClient, referrer, referencedSecretNameMap)
}

// UpdateReferencesToConfigMap updates the reference records between referrer and each ConfigMap
// in namespacedNames in record cache.
func UpdateReferencesToConfigMap(
	ctx context.Context,
	c client.Client, indexers CacheIndexers, dataplaneClient controllers.DataPlaneClient,
	referrer client.Object, referencedConfigMapNameMap map[k8stypes.NamespacedName]struct{},
) error {
	for nsName := range referencedConfigMapNameMap {
		configMap := &corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				APIVersion: VersionV1,
				Kind:       KindConfigMap,
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: nsName.Namespace,
				Name:      nsName.Name,
			},
		}

		// Here we update the reference relationship even when the referred configmap does not exist yet
		// If the referred ConfiMap is created, it could be reconciled in ConfigMap controller.
		referrerCopy := referrer.DeepCopyObject().(client.Object)
		if err := indexers.SetObjectReference(
			referrerCopy, configMap.DeepCopy()); err != nil {
			return err
		}

		if err := c.Get(ctx, nsName, configMap); err != nil {
			return err
		}

		if err := dataplaneClient.UpdateObject(configMap); err != nil {
			return err
		}
	}

	return removeOutdatedReferencesToConfigMap(ctx, indexers, c, dataplaneClient, referrer, referencedConfigMapNameMap)
}

// removeOutdatedReferenceToSecret removes outdated reference records to secrets in reference indexer.
// secrets that are referred by referrer are passed in referredSecretNames parameter.
// If a secret is not referenced by any other object after deleting outdated reference records,
// and it does not have label "konghq.com/ca-cert:true", it is not possible to be used in Kong gateway config
// and should be removed from the object cache inside KongClient.
func removeOutdatedReferencesToSecret(
	ctx context.Context,
	indexers CacheIndexers, c client.Client, dataplaneClient controllers.DataPlaneClient,
	referrer client.Object, referredSecretNameMap map[k8stypes.NamespacedName]struct{},
) error {
	referents, err := indexers.ListReferredObjects(referrer)
	if err != nil {
		return err
	}
	for _, obj := range referents {
		gvk := obj.GetObjectKind().GroupVersionKind()
		// delete the reference record if the secret is not referred by the service.
		if gvk.Group == corev1.GroupName && gvk.Version == VersionV1 && gvk.Kind == KindSecret {
			namespacedName := k8stypes.NamespacedName{
				Namespace: obj.GetNamespace(),
				Name:      obj.GetName(),
			}

			// if the secret is still referenced, no operations are taken so continue here.
			if _, ok := referredSecretNameMap[namespacedName]; ok {
				continue
			}

			if err := indexers.DeleteObjectReference(referrer, obj); err != nil {
				return err
			}
			// remove the secret in object cache if it is not referred and does not have label "konghq.com/ca-cert:true".
			// Do this check and delete when the reference count may be reduced by 1.

			// retrieve the secret in k8s and check it has the label.
			secret := &corev1.Secret{}
			getErr := c.Get(ctx, namespacedName, secret)
			// if the secret exists in k8s and has the label, we should not delete it in object cache.
			if getErr == nil {
				if secret.Labels != nil && secret.Labels[CACertLabelKey] == "true" {
					continue
				}
			} else {
				// if the secret does not exist in k8s, we ignore the error and continue the check and delete operation.
				// for other errors, we return the error and stop the operation.
				if !apierrors.IsNotFound(getErr) {
					return err
				}
			}

			if err := indexers.DeleteObjectIfNotReferred(obj, dataplaneClient); err != nil {
				return err
			}
		}
	}
	return nil
}

// removeOutdatedReferencesToConfigMap removes outdated reference records to configMaps in reference indexer.
// ConfigMaps that are referred by referrer are passed in referredConfigMapNameMap parameter.
// If a ConfigMap is not referenced by any other object after deleting outdated reference records,
// and it does not have label "konghq.com/ca-cert:true", it is not possible to be used in Kong gateway config
// and should be removed from the object cache inside KongClient.
func removeOutdatedReferencesToConfigMap(
	ctx context.Context,
	indexers CacheIndexers, c client.Client, dataplaneClient controllers.DataPlaneClient,
	referrer client.Object, referredConfigMapNameMap map[k8stypes.NamespacedName]struct{},
) error {
	referents, err := indexers.ListReferredObjects(referrer)
	if err != nil {
		return err
	}
	for _, obj := range referents {
		gvk := obj.GetObjectKind().GroupVersionKind()
		// Delete the reference record if the secret is not referred by the service.
		if gvk.Group == corev1.GroupName && gvk.Version == VersionV1 && gvk.Kind == KindConfigMap {
			namespacedName := k8stypes.NamespacedName{
				Namespace: obj.GetNamespace(),
				Name:      obj.GetName(),
			}

			// if the configmap is still referenced, no operations are taken so continue here.
			if _, ok := referredConfigMapNameMap[namespacedName]; ok {
				continue
			}

			if err := indexers.DeleteObjectReference(referrer, obj); err != nil {
				return err
			}
			// remove the configMap in object cache if it is not referred and does not have label "konghq.com/ca-cert:true".
			// Do this check and delete when the reference count may be reduced by 1.

			// retrieve the configMap in k8s and check it has the label.
			configMap := &corev1.ConfigMap{}
			getErr := c.Get(ctx, namespacedName, configMap)
			// if the configMap exists in k8s and has the label, we should not delete it in object cache.
			if getErr != nil {
				// if the configMap does not exist in k8s, we ignore the error and continue the check and delete operation.
				// for other errors, we return the error and stop the operation.
				if !apierrors.IsNotFound(getErr) {
					return err
				}
			}
			if configMap.Labels != nil && configMap.Labels[CACertLabelKey] == "true" {
				continue
			}

			if err := indexers.DeleteObjectIfNotReferred(obj, dataplaneClient); err != nil {
				return err
			}
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
