package configuration

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
)

func TestListCoreV1ServiceReferredSecrets(t *testing.T) {
	testCases := []struct {
		name          string
		service       *corev1.Service
		secretNum     int
		refSecretName k8stypes.NamespacedName
	}{
		{
			name: "service_has_no_annotations",
			service: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns1",
					Name:      "service1",
				},
			},
			secretNum: 0,
		},
		{
			name: "service_not_referring_secret_in_annotations",
			service: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns1",
					Name:      "service2",
				},
			},
			secretNum: 0,
		},
		{
			name: "service_referring_secret_in_annotations",
			service: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns1",
					Name:      "service2",
					Annotations: map[string]string{
						"konghq.com/client-cert": "secret1",
					},
				},
			},
			secretNum:     1,
			refSecretName: k8stypes.NamespacedName{Namespace: "ns1", Name: "secret1"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			secretNames := listCoreV1ServiceReferredSecrets(tc.service)
			require.Len(t, secretNames, tc.secretNum)
			if tc.secretNum > 0 {
				require.Contains(t, secretNames, tc.refSecretName)
			}
		})
	}
}

func TestListIngressReferredSecrets(t *testing.T) {
	testCases := []struct {
		name          string
		ingress       *netv1.Ingress
		secretNum     int
		refSecretName k8stypes.NamespacedName
	}{
		{
			name: "ingress_has_no_tls_should_refer_no_secrets",
			ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "ing1",
				},
			},
			secretNum: 0,
		},
		{
			name: "ingress_has_tls_should_refer_to_secrets",
			ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "ing1",
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{Hosts: []string{"example.com"}, SecretName: "secret1"},
					},
				},
			},
			secretNum:     1,
			refSecretName: k8stypes.NamespacedName{Namespace: "ns", Name: "secret1"},
		},
		{
			name: "ingress_has_tls_without_secretName_should_refer_no_secrets",
			ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "ing1",
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{Hosts: []string{"example.com"}},
					},
				},
			},
			secretNum: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			secretNames := listNetV1IngressReferredSecrets(tc.ingress)
			require.Len(t, secretNames, tc.secretNum)
			if tc.secretNum > 0 {
				require.Contains(t, secretNames, tc.refSecretName)
			}
		})
	}
}

func TestListKongPluginReferredSecrets(t *testing.T) {
	testCases := []struct {
		name          string
		plugin        *kongv1.KongPlugin
		secretNum     int
		refSecretName k8stypes.NamespacedName
	}{
		{
			name: "kong_plugin_refer_no_secrets",
			plugin: &kongv1.KongPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "plugin1",
				},
			},
			secretNum: 0,
		},
		{
			name: "kong_plugin_refer_secrets",
			plugin: &kongv1.KongPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "plugin1",
				},
				ConfigFrom: &kongv1.ConfigSource{
					SecretValue: kongv1.SecretValueFromSource{
						Secret: "secret1",
						Key:    "k",
					},
				},
			},
			secretNum: 1,
			refSecretName: k8stypes.NamespacedName{
				Namespace: "ns",
				Name:      "secret1",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			secretNames := listKongPluginReferredSecrets(tc.plugin)
			require.Len(t, secretNames, tc.secretNum)
			if tc.secretNum > 0 {
				require.Contains(t, secretNames, tc.refSecretName)
			}
		})
	}
}

func TestListKongClusterPluginReferredSecrets(t *testing.T) {
	testCases := []struct {
		name          string
		plugin        *kongv1.KongClusterPlugin
		secretNum     int
		refSecretName k8stypes.NamespacedName
	}{
		{
			name: "kong_cluster_plugin_refer_no_secrets",
			plugin: &kongv1.KongClusterPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name: "plugin1",
				},
			},
			secretNum: 0,
		},
		{
			name: "kong_cluster_plugin_refer_secrets",
			plugin: &kongv1.KongClusterPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name: "plugin1",
				},
				ConfigFrom: &kongv1.NamespacedConfigSource{
					SecretValue: kongv1.NamespacedSecretValueFromSource{
						Namespace: "ns",
						Secret:    "secret1",
						Key:       "k",
					},
				},
			},
			secretNum: 1,
			refSecretName: k8stypes.NamespacedName{
				Namespace: "ns",
				Name:      "secret1",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			secretNames := listKongClusterPluginReferredSecrets(tc.plugin)
			require.Len(t, secretNames, tc.secretNum)
			if tc.secretNum > 0 {
				require.Contains(t, secretNames, tc.refSecretName)
			}
		})
	}
}

func TestListKongConsumerReferredSecrets(t *testing.T) {
	testCases := []struct {
		name          string
		consumer      *kongv1.KongConsumer
		secretNum     int
		refSecretName k8stypes.NamespacedName
	}{
		{
			name: "consumer_refer_no_secrets",
			consumer: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "consumer1",
				},
			},
			secretNum: 0,
		},
		{
			name: "consumer_refer_secrets",
			consumer: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "consumer1",
				},
				Credentials: []string{"secret1", "secret2"},
			},
			secretNum: 2,
			refSecretName: k8stypes.NamespacedName{
				Namespace: "ns",
				Name:      "secret1",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			secretNames := listKongConsumerReferredSecrets(tc.consumer)
			require.Len(t, secretNames, tc.secretNum)
			if tc.secretNum > 0 {
				require.Contains(t, secretNames, tc.refSecretName)
			}
		})
	}
}

func TestListTCPIngressReferredSecrets(t *testing.T) {
	testCases := []struct {
		name          string
		tcpIngress    *kongv1beta1.TCPIngress
		secretNum     int
		refSecretName k8stypes.NamespacedName
	}{
		{
			name: "tcp_ingress_refer_no_secrets",
			tcpIngress: &kongv1beta1.TCPIngress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "ingress1",
				},
			},
			secretNum: 0,
		},
		{
			name: "tcp_ingress_refer_secrets",
			tcpIngress: &kongv1beta1.TCPIngress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "ingress1",
				},
				Spec: kongv1beta1.TCPIngressSpec{
					TLS: []kongv1beta1.IngressTLS{
						{Hosts: []string{"example.com"}, SecretName: "secret1"},
						{Hosts: []string{"konghq.com"}, SecretName: ""},
					},
				},
			},
			secretNum: 1,
			refSecretName: k8stypes.NamespacedName{
				Namespace: "ns",
				Name:      "secret1",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			secretNames := listTCPIngressReferredSecrets(tc.tcpIngress)
			require.Len(t, secretNames, tc.secretNum)
			if tc.secretNum > 0 {
				require.Contains(t, secretNames, tc.refSecretName)
			}
		})
	}
}
