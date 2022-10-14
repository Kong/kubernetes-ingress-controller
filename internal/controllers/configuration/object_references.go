package configuration

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

// updateReferredObjects updates reference records where the referrer is the object in parameter obj.
// currently it only updates reference records to secrets, since we wanted to limit cache size of secrets:
// https://github.com/Kong/kubernetes-ingress-controller/issues/2868
func updateReferredObjects(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, obj client.Object) error {
	switch obj := obj.(type) {
	// functions update***ReferredSecrets first list the secrets referred by object,
	// then call UpdateReferencesToSecret to store referrence records between the object and referred secrets,
	// and also to remove the outdated reference records in cache where the secret is not referred by the obj specification anymore.
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
	referredSecretNames := make([]types.NamespacedName, 0, 1)

	anns := service.Annotations
	if service.Annotations == nil {
		anns = map[string]string{}
	}
	secretName := annotations.ExtractClientCertificate(anns)
	if secretName != "" {
		nsName := types.NamespacedName{
			Namespace: service.Namespace,
			Name:      secretName,
		}

		referredSecretNames = append(referredSecretNames, nsName)
	}

	return reference.UpdateReferencesToSecret(ctx, client, dataplaneClient, service, referredSecretNames)
}

func updateNetV1IngressReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, ingress *netv1.Ingress,
) error {
	referredSecretNames := make([]types.NamespacedName, 0, len(ingress.Spec.TLS))
	for _, tls := range ingress.Spec.TLS {
		nsName := types.NamespacedName{
			Namespace: ingress.Namespace,
			Name:      tls.SecretName,
		}
		referredSecretNames = append(referredSecretNames, nsName)
	}

	return reference.UpdateReferencesToSecret(ctx, client, dataplaneClient, ingress, referredSecretNames)
}

func updateNetV1beta1IngressReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, ingress *netv1beta1.Ingress,
) error {
	referredSecretNames := make([]types.NamespacedName, 0, len(ingress.Spec.TLS))
	for _, tls := range ingress.Spec.TLS {
		nsName := types.NamespacedName{
			Namespace: ingress.Namespace,
			Name:      tls.SecretName,
		}
		referredSecretNames = append(referredSecretNames, nsName)
	}

	return reference.UpdateReferencesToSecret(ctx, client, dataplaneClient, ingress, referredSecretNames)
}

func updateExtensionV1beta1IngressReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, ingress *extv1beta1.Ingress,
) error {
	referredSecretNames := make([]types.NamespacedName, 0, len(ingress.Spec.TLS))

	for _, tls := range ingress.Spec.TLS {
		nsName := types.NamespacedName{
			Namespace: ingress.Namespace,
			Name:      tls.SecretName,
		}
		referredSecretNames = append(referredSecretNames, nsName)
	}

	return reference.UpdateReferencesToSecret(ctx, client, dataplaneClient, ingress, referredSecretNames)
}

func updateKongPluginReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, plugin *kongv1.KongPlugin,
) error {
	referredSecretNames := make([]types.NamespacedName, 0, 1)

	if plugin.ConfigFrom != nil {
		nsName := types.NamespacedName{
			Namespace: plugin.Namespace,
			Name:      plugin.ConfigFrom.SecretValue.Secret,
		}
		referredSecretNames = append(referredSecretNames, nsName)
	}

	return reference.UpdateReferencesToSecret(ctx, client, dataplaneClient, plugin, referredSecretNames)
}

func updateKongClusterPluginReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, plugin *kongv1.KongClusterPlugin,
) error {
	referredSecretNames := make([]types.NamespacedName, 0, 1)

	if plugin.ConfigFrom != nil {
		nsName := types.NamespacedName{
			Namespace: plugin.ConfigFrom.SecretValue.Namespace,
			Name:      plugin.ConfigFrom.SecretValue.Secret,
		}
		referredSecretNames = append(referredSecretNames, nsName)

	}

	return reference.UpdateReferencesToSecret(ctx, client, dataplaneClient, plugin, referredSecretNames)
}

func updateKongConsumerReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, consumer *kongv1.KongConsumer,
) error {
	referredSecretNames := make([]types.NamespacedName, 0, len(consumer.Credentials))

	for _, secretName := range consumer.Credentials {
		nsName := types.NamespacedName{
			Namespace: consumer.Namespace,
			Name:      secretName,
		}
		referredSecretNames = append(referredSecretNames, nsName)
	}

	return reference.UpdateReferencesToSecret(ctx, client, dataplaneClient, consumer, referredSecretNames)
}

func updateTCPIngressReferredSecrets(
	ctx context.Context, client client.Client, dataplaneClient *dataplane.KongClient, tcpIngress *kongv1beta1.TCPIngress,
) error {
	referredSecretNames := make([]types.NamespacedName, len(tcpIngress.Spec.TLS))
	for _, tls := range tcpIngress.Spec.TLS {
		nsName := types.NamespacedName{
			Namespace: tcpIngress.Namespace,
			Name:      tls.SecretName,
		}
		referredSecretNames = append(referredSecretNames, nsName)
	}

	return reference.UpdateReferencesToSecret(ctx, client, dataplaneClient, tcpIngress, referredSecretNames)
}
