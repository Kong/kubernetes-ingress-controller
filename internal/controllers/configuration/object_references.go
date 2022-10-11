package configuration

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

func updateReferredObjects(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, obj runtime.Object) error {
	switch obj := obj.(type) {
	case *corev1.Service:
		return updateCoreV1ServiceReferredSecrets(ctx, client, dataplaneClient, obj)
	case *netv1.Ingress:
		return updateNetV1IngressReferredSecrets(ctx, client, dataplaneClient, obj)
	case *netv1beta1.Ingress:
		return updateNetV1beta1IngressReferredSecrets(ctx, client, dataplaneClient, obj)
	case *extv1beta1.Ingress:
		return updateExtensionV1beta1IngressReferredSecrets(ctx, client, dataplaneClient, obj)
	case *kongv1.KongPlugin:
		return updateKongPluginReferredSecrets(ctx, client, dataplaneClient, obj)
	case *kongv1.KongClusterPlugin:
		return updateKongClusterPluginReferredSecrets(ctx, client, dataplaneClient, obj)
	case *kongv1.KongConsumer:
		return updateKongConsumerReferredSecrets(ctx, client, dataplaneClient, obj)
	case *kongv1beta1.TCPIngress:
		return updateTCPIngressReferredSecrets(ctx, client, dataplaneClient, obj)
	}

	return nil
}

func updateCoreV1ServiceReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, service *corev1.Service,
) error {
	if service.Annotations == nil {
		return nil
	}
	secretName := annotations.ExtractClientCertificate(service.Annotations)
	if secretName != "" {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: service.Namespace,
				Name:      secretName,
			},
		}

		err := client.Get(ctx, types.NamespacedName{
			Namespace: secret.Namespace,
			Name:      secret.Name,
		}, secret)
		if err != nil {
			return err
		}

		err = dataplaneClient.SetObjectReference(service.DeepCopy(), secret.DeepCopy())
		if err != nil {
			return err
		}

		if err := dataplaneClient.UpdateObject(secret); err != nil {
			return err
		}
	}
	// TODO: remove outdated reference records.

	return nil
}

func updateNetV1IngressReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, ingress *netv1.Ingress,
) error {
	for _, tls := range ingress.Spec.TLS {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ingress.Namespace,
				Name:      tls.SecretName,
			},
		}

		err := client.Get(ctx, types.NamespacedName{
			Namespace: secret.Namespace,
			Name:      secret.Name,
		}, secret)
		if err != nil {
			return err
		}

		err = dataplaneClient.SetObjectReference(ingress.DeepCopy(), secret.DeepCopy())
		if err != nil {
			return err
		}

		if err := dataplaneClient.UpdateObject(secret); err != nil {
			return err
		}

	}

	return nil
}

func updateNetV1beta1IngressReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, ingress *netv1beta1.Ingress,
) error {

	for _, tls := range ingress.Spec.TLS {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ingress.Namespace,
				Name:      tls.SecretName,
			},
		}

		err := client.Get(ctx, types.NamespacedName{
			Namespace: secret.Namespace,
			Name:      secret.Name,
		}, secret)
		if err != nil {
			return err
		}

		err = dataplaneClient.SetObjectReference(ingress.DeepCopy(), secret.DeepCopy())
		if err != nil {
			return err
		}

		if err := dataplaneClient.UpdateObject(secret); err != nil {
			return err
		}
	}

	return nil
}

func updateExtensionV1beta1IngressReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, ingress *extv1beta1.Ingress,
) error {

	for _, tls := range ingress.Spec.TLS {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ingress.Namespace,
				Name:      tls.SecretName,
			},
		}

		err := client.Get(ctx, types.NamespacedName{
			Namespace: secret.Namespace,
			Name:      secret.Name,
		}, secret)
		if err != nil {
			return err
		}

		err = dataplaneClient.SetObjectReference(ingress.DeepCopy(), secret.DeepCopy())
		if err != nil {
			return err
		}

		if err := dataplaneClient.UpdateObject(secret); err != nil {
			return err
		}
	}

	return nil
}

func updateKongPluginReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, plugin *kongv1.KongPlugin,
) error {
	if plugin.ConfigFrom != nil {
		secretName := plugin.ConfigFrom.SecretValue.Secret
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: plugin.Namespace,
				Name:      secretName,
			},
		}

		err := client.Get(ctx, types.NamespacedName{
			Namespace: secret.Namespace,
			Name:      secret.Name,
		}, secret)
		if err != nil {
			return err
		}

		err = dataplaneClient.SetObjectReference(plugin.DeepCopy(), secret.DeepCopy())
		if err != nil {
			return err
		}

		if err := dataplaneClient.UpdateObject(secret); err != nil {
			return err
		}
	}

	return nil
}

func updateKongClusterPluginReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, plugin *kongv1.KongClusterPlugin,
) error {
	if plugin.ConfigFrom != nil {
		secretNamespace := plugin.ConfigFrom.SecretValue.Namespace
		secretName := plugin.ConfigFrom.SecretValue.Secret
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: secretNamespace,
				Name:      secretName,
			},
		}

		err := client.Get(ctx, types.NamespacedName{
			Namespace: secret.Namespace,
			Name:      secret.Name,
		}, secret)
		if err != nil {
			return err
		}

		err = dataplaneClient.SetObjectReference(plugin.DeepCopy(), secret.DeepCopy())
		if err != nil {
			return err
		}

		if err := dataplaneClient.UpdateObject(secret); err != nil {
			return err
		}
	}

	return nil
}

func updateKongConsumerReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, consumer *kongv1.KongConsumer,
) error {
	for _, secretName := range consumer.Credentials {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: consumer.Namespace,
				Name:      secretName,
			},
		}

		err := client.Get(ctx, types.NamespacedName{
			Namespace: secret.Namespace,
			Name:      secret.Name,
		}, secret)
		if err != nil {
			return err
		}

		err = dataplaneClient.SetObjectReference(consumer.DeepCopy(), secret.DeepCopy())
		if err != nil {
			return err
		}

		if err := dataplaneClient.UpdateObject(secret); err != nil {
			return err
		}
	}

	return nil
}

func updateTCPIngressReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, tcpIngress *kongv1beta1.TCPIngress,
) error {
	for _, tls := range tcpIngress.Spec.TLS {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: tcpIngress.Namespace,
				Name:      tls.SecretName,
			},
		}

		err := client.Get(ctx, types.NamespacedName{
			Namespace: secret.Namespace,
			Name:      secret.Name,
		}, secret)
		if err != nil {
			return err
		}

		err = dataplaneClient.SetObjectReference(tcpIngress.DeepCopy(), secret.DeepCopy())
		if err != nil {
			return err
		}

		if err := dataplaneClient.UpdateObject(secret); err != nil {
			return err
		}
	}

	return nil
}
