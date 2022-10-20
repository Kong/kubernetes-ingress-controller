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
	ctrlref "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

// updateReferredObjects updates reference records where the referrer is the object in parameter obj.
// currently it only updates reference records to secrets, since we wanted to limit cache size of secrets:
// https://github.com/Kong/kubernetes-ingress-controller/issues/2868
func updateReferredObjects(
	ctx context.Context, client client.Client, refIndexers ctrlref.CacheIndexers, dataplaneClient *dataplane.KongClient, obj client.Object) error {

	var referredSecretNames = make(map[types.NamespacedName]struct{})
	switch obj := obj.(type) {
	// functions update***ReferredSecrets first list the secrets referred by object,
	// then call UpdateReferencesToSecret to store referrence records between the object and referred secrets,
	// and also to remove the outdated reference records in cache where the secret is not referred by the obj specification anymore.
	case *corev1.Service:
		listCoreV1ServiceReferredSecrets(obj, referredSecretNames)
	case *netv1.Ingress:
		listNetV1IngressReferredSecrets(obj, referredSecretNames)
	case *netv1beta1.Ingress:
		listNetV1beta1IngressReferredSecrets(obj, referredSecretNames)
	case *extv1beta1.Ingress:
		listExtensionV1beta1IngressReferredSecrets(obj, referredSecretNames)
	case *kongv1.KongPlugin:
		listKongPluginReferredSecrets(obj, referredSecretNames)
	case *kongv1.KongClusterPlugin:
		listKongClusterPluginReferredSecrets(obj, referredSecretNames)
	case *kongv1.KongConsumer:
		listKongConsumerReferredSecrets(obj, referredSecretNames)
	case *kongv1beta1.TCPIngress:
		listTCPIngressReferredSecrets(obj, referredSecretNames)
	}

	return ctrlref.UpdateReferencesToSecret(ctx, client, refIndexers, dataplaneClient, obj, referredSecretNames)
}

func listCoreV1ServiceReferredSecrets(service *corev1.Service, referredSecretNames map[types.NamespacedName]struct{}) {

	if service.Annotations == nil {
		return
	}

	secretName := annotations.ExtractClientCertificate(service.Annotations)
	if secretName != "" {
		nsName := types.NamespacedName{
			Namespace: service.Namespace,
			Name:      secretName,
		}

		referredSecretNames[nsName] = struct{}{}
	}
}

func listNetV1IngressReferredSecrets(ingress *netv1.Ingress, referredSecretNames map[types.NamespacedName]struct{}) {
	for _, tls := range ingress.Spec.TLS {
		nsName := types.NamespacedName{
			Namespace: ingress.Namespace,
			Name:      tls.SecretName,
		}
		referredSecretNames[nsName] = struct{}{}
	}
}

func listNetV1beta1IngressReferredSecrets(ingress *netv1beta1.Ingress, referredSecretNames map[types.NamespacedName]struct{}) {
	for _, tls := range ingress.Spec.TLS {
		nsName := types.NamespacedName{
			Namespace: ingress.Namespace,
			Name:      tls.SecretName,
		}
		referredSecretNames[nsName] = struct{}{}
	}
}

func listExtensionV1beta1IngressReferredSecrets(ingress *extv1beta1.Ingress, referredSecretNames map[types.NamespacedName]struct{}) {
	for _, tls := range ingress.Spec.TLS {
		nsName := types.NamespacedName{
			Namespace: ingress.Namespace,
			Name:      tls.SecretName,
		}
		referredSecretNames[nsName] = struct{}{}
	}
}

func listKongPluginReferredSecrets(plugin *kongv1.KongPlugin, referredSecretNames map[types.NamespacedName]struct{}) {
	if plugin.ConfigFrom != nil {
		nsName := types.NamespacedName{
			Namespace: plugin.Namespace,
			Name:      plugin.ConfigFrom.SecretValue.Secret,
		}
		referredSecretNames[nsName] = struct{}{}
	}

}

func listKongClusterPluginReferredSecrets(plugin *kongv1.KongClusterPlugin, referredSecretNames map[types.NamespacedName]struct{}) {
	if plugin.ConfigFrom != nil {
		nsName := types.NamespacedName{
			Namespace: plugin.ConfigFrom.SecretValue.Namespace,
			Name:      plugin.ConfigFrom.SecretValue.Secret,
		}
		referredSecretNames[nsName] = struct{}{}
	}
}

func listKongConsumerReferredSecrets(consumer *kongv1.KongConsumer, referredSecretNames map[types.NamespacedName]struct{}) {
	for _, secretName := range consumer.Credentials {
		nsName := types.NamespacedName{
			Namespace: consumer.Namespace,
			Name:      secretName,
		}
		referredSecretNames[nsName] = struct{}{}
	}
}

func listTCPIngressReferredSecrets(tcpIngress *kongv1beta1.TCPIngress, referredSecretNames map[types.NamespacedName]struct{}) {
	for _, tls := range tcpIngress.Spec.TLS {
		nsName := types.NamespacedName{
			Namespace: tcpIngress.Namespace,
			Name:      tls.SecretName,
		}
		referredSecretNames[nsName] = struct{}{}
	}
}
