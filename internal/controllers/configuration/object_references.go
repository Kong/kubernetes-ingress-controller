package configuration

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

// updateReferredObjects updates reference records where the referrer is the object in parameter obj.
// currently it only updates reference records to secrets, since we wanted to limit cache size of secrets:
// https://github.com/Kong/kubernetes-ingress-controller/issues/2868
func updateReferredObjects(
	ctx context.Context, client client.Client, refIndexers ctrlref.CacheIndexers, dataplaneClient *dataplane.KongClient, obj client.Object,
) error {
	referredSecretNameMap := make(map[types.NamespacedName]struct{})
	var referredSecretList []types.NamespacedName
	switch obj := obj.(type) {
	// functions update***ReferredSecrets first list the secrets referred by object,
	// then call UpdateReferencesToSecret to store reference records between the object and referred secrets,
	// and also to remove the outdated reference records in cache where the secret is not referred by the obj specification anymore.
	case *corev1.Service:
		referredSecretList = listCoreV1ServiceReferredSecrets(obj)
	case *netv1.Ingress:
		referredSecretList = listNetV1IngressReferredSecrets(obj)
	case *kongv1.KongPlugin:
		referredSecretList = listKongPluginReferredSecrets(obj)
	case *kongv1.KongClusterPlugin:
		referredSecretList = listKongClusterPluginReferredSecrets(obj)
	case *kongv1.KongConsumer:
		referredSecretList = listKongConsumerReferredSecrets(obj)
	case *kongv1beta1.TCPIngress:
		referredSecretList = listTCPIngressReferredSecrets(obj)
	}

	for _, nsName := range referredSecretList {
		referredSecretNameMap[nsName] = struct{}{}
	}
	return ctrlref.UpdateReferencesToSecret(ctx, client, refIndexers, dataplaneClient, obj, referredSecretNameMap)
}

func listCoreV1ServiceReferredSecrets(service *corev1.Service) []types.NamespacedName {
	if service.Annotations == nil {
		return nil
	}

	referredSecretNames := make([]types.NamespacedName, 0, 1)
	secretName := annotations.ExtractClientCertificate(service.Annotations)
	if secretName != "" {
		nsName := types.NamespacedName{
			Namespace: service.Namespace,
			Name:      secretName,
		}

		referredSecretNames = append(referredSecretNames, nsName)
	}
	return referredSecretNames
}

func listNetV1IngressReferredSecrets(ingress *netv1.Ingress) []types.NamespacedName {
	referredSecretNames := make([]types.NamespacedName, 0, len(ingress.Spec.TLS))
	for _, tls := range ingress.Spec.TLS {
		if tls.SecretName == "" {
			continue
		}
		nsName := types.NamespacedName{
			Namespace: ingress.Namespace,
			Name:      tls.SecretName,
		}
		referredSecretNames = append(referredSecretNames, nsName)
	}
	return referredSecretNames
}

func listKongPluginReferredSecrets(plugin *kongv1.KongPlugin) []types.NamespacedName {
	referredSecretNames := make([]types.NamespacedName, 0, 1)
	if plugin.ConfigFrom != nil {
		nsName := types.NamespacedName{
			Namespace: plugin.Namespace,
			Name:      plugin.ConfigFrom.SecretValue.Secret,
		}
		referredSecretNames = append(referredSecretNames, nsName)
	}
	return referredSecretNames
}

func listKongClusterPluginReferredSecrets(plugin *kongv1.KongClusterPlugin) []types.NamespacedName {
	referredSecretNames := make([]types.NamespacedName, 0, 1)
	if plugin.ConfigFrom != nil {
		nsName := types.NamespacedName{
			Namespace: plugin.ConfigFrom.SecretValue.Namespace,
			Name:      plugin.ConfigFrom.SecretValue.Secret,
		}
		referredSecretNames = append(referredSecretNames, nsName)
	}
	return referredSecretNames
}

func listKongConsumerReferredSecrets(consumer *kongv1.KongConsumer) []types.NamespacedName {
	referredSecretNames := make([]types.NamespacedName, 0, len(consumer.Credentials))
	for _, secretName := range consumer.Credentials {
		nsName := types.NamespacedName{
			Namespace: consumer.Namespace,
			Name:      secretName,
		}
		referredSecretNames = append(referredSecretNames, nsName)
	}
	return referredSecretNames
}

func listTCPIngressReferredSecrets(tcpIngress *kongv1beta1.TCPIngress) []types.NamespacedName {
	referredSecretNames := make([]types.NamespacedName, 0, len(tcpIngress.Spec.TLS))
	for _, tls := range tcpIngress.Spec.TLS {
		if tls.SecretName == "" {
			continue
		}
		nsName := types.NamespacedName{
			Namespace: tcpIngress.Namespace,
			Name:      tls.SecretName,
		}
		referredSecretNames = append(referredSecretNames, nsName)
	}
	return referredSecretNames
}
