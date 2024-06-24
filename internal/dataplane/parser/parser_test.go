package parser

import (
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	netv1 "k8s.io/api/networking/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v2/test/helpers/certificate"
)

func TestGlobalPlugin(t *testing.T) {
	assert := assert.New(t)
	t.Run("global plugins are processed correctly", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			KongClusterPlugins: []*kongv1.KongClusterPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "bar-plugin",
						Labels: map[string]string{
							"global": "true",
						},
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Protocols:  kongv1.StringsToKongProtocols([]string{"http"}),
					PluginName: "basic-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"foo1": "bar1"}`),
					},
				},
			},
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		assert.Equal(1, len(state.Plugins),
			"expected one plugin to be rendered")

		sort.SliceStable(state.Plugins, func(i, j int) bool {
			return strings.Compare(*state.Plugins[i].Name, *state.Plugins[j].Name) > 0
		})

		assert.Equal("basic-auth", *state.Plugins[0].Name)
		assert.Equal(kong.Configuration{"foo1": "bar1"}, state.Plugins[0].Config)
	})
}

func TestSecretConfigurationPlugin(t *testing.T) {
	jwtPluginConfig := `{"run_on_preflight": false}`  // JSON
	basicAuthPluginConfig := "hide_credentials: true" // YAML
	assert := assert.New(t)
	stock := store.FakeObjects{
		Services: []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		},
		IngressesV1: []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "foo-plugin",
						annotations.IngressClassKey:                           annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "bar-plugin",
						annotations.IngressClassKey:                           annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.net",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	t.Run("plugins with secret configuration are processed correctly",
		func(t *testing.T) {
			objects := stock
			objects.KongPlugins = []*kongv1.KongPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-plugin",
						Namespace: "default",
					},
					PluginName: "jwt",
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Key:    "jwt-config",
							Secret: "conf-secret",
						},
					},
				},
			}
			objects.KongClusterPlugins = []*kongv1.KongClusterPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "global-bar-plugin",
						Labels: map[string]string{
							"global": "true",
						},
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Protocols:  kongv1.StringsToKongProtocols([]string{"http"}),
					PluginName: "basic-auth",
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "global-broken-bar-plugin",
						Labels: map[string]string{
							"global": "true",
						},
						Annotations: map[string]string{
							// explicitly none, this should not get rendered
						},
					},
					Protocols:  kongv1.StringsToKongProtocols([]string{"http"}),
					PluginName: "basic-auth",
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "bar-plugin",
					},
					Protocols:  kongv1.StringsToKongProtocols([]string{"http"}),
					PluginName: "basic-auth",
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
			}
			objects.Secrets = []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						UID:       k8stypes.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
						Name:      "conf-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"jwt-config":        []byte(jwtPluginConfig),
						"basic-auth-config": []byte(basicAuthPluginConfig),
					},
				},
			}
			store, err := store.NewFakeStore(objects)
			require.NoError(t, err)
			p := mustNewParser(t, store)
			result := p.BuildKongConfig()
			require.Empty(t, result.TranslationFailures)
			require.NoError(t, err)
			state := result.KongState
			require.NotNil(t, state)
			assert.Equal(3, len(state.Plugins),
				"expected three plugins to be rendered")

			sort.SliceStable(state.Plugins, func(i, j int) bool {
				return strings.Compare(*state.Plugins[i].Name,
					*state.Plugins[j].Name) > 0
			})
			assert.Equal("jwt", *state.Plugins[0].Name)
			assert.Equal(kong.Configuration{"run_on_preflight": false},
				state.Plugins[0].Config)

			assert.Equal("basic-auth", *state.Plugins[1].Name)
			assert.Equal(kong.Configuration{"hide_credentials": true},
				state.Plugins[2].Config)
			assert.Equal("basic-auth", *state.Plugins[2].Name)
			assert.Equal(kong.Configuration{"hide_credentials": true},
				state.Plugins[2].Config)
		})

	t.Run("plugins with missing secrets or keys are not constructed",
		func(t *testing.T) {
			objects := stock
			objects.KongPlugins = []*kongv1.KongPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "global-foo-plugin",
						Namespace: "default",
						Labels: map[string]string{
							"global": "true",
						},
					},
					PluginName: "jwt",
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Key:    "missing-key",
							Secret: "conf-secret",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-plugin",
						Namespace: "default",
					},
					PluginName: "jwt",
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Key:    "missing-key",
							Secret: "conf-secret",
						},
					},
				},
			}
			objects.KongClusterPlugins = []*kongv1.KongClusterPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "global-bar-plugin",
						Labels: map[string]string{
							"global": "true",
						},
					},
					Protocols:  kongv1.StringsToKongProtocols([]string{"http"}),
					PluginName: "basic-auth",
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "missing-secret",
							Namespace: "default",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "bar-plugin",
					},
					Protocols:  kongv1.StringsToKongProtocols([]string{"http"}),
					PluginName: "basic-auth",
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "missing-secret",
							Namespace: "default",
						},
					},
				},
			}
			objects.Secrets = []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						UID:       k8stypes.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
						Name:      "conf-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"jwt-config":        []byte(jwtPluginConfig),
						"basic-auth-config": []byte(basicAuthPluginConfig),
					},
				},
			}
			store, err := store.NewFakeStore(objects)
			require.NoError(t, err)
			p := mustNewParser(t, store)
			result := p.BuildKongConfig()
			require.Empty(t, result.TranslationFailures)
			require.NoError(t, err)
			state := result.KongState
			require.NotNil(t, state)
			assert.Equal(0, len(state.Plugins),
				"expected no plugins to be rendered")
		})

	t.Run("plugins with both config and configFrom are not constructed",
		func(t *testing.T) {
			objects := stock
			objects.KongPlugins = []*kongv1.KongPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "global-foo-plugin",
						Namespace: "default",
						Labels: map[string]string{
							"global": "true",
						},
					},
					PluginName: "jwt",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"fake": true}`),
					},
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Key:    "jwt-config",
							Secret: "conf-secret",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-plugin",
						Namespace: "default",
					},
					PluginName: "jwt",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"fake": true}`),
					},
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Key:    "jwt-config",
							Secret: "conf-secret",
						},
					},
				},
			}
			objects.KongClusterPlugins = []*kongv1.KongClusterPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "global-bar-plugin",
						Labels: map[string]string{
							"global": "true",
						},
					},
					Protocols:  kongv1.StringsToKongProtocols([]string{"http"}),
					PluginName: "basic-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"fake": true}`),
					},
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "bar-plugin",
					},
					Protocols:  kongv1.StringsToKongProtocols([]string{"http"}),
					PluginName: "basic-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"fake": true}`),
					},
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
			}
			objects.Secrets = []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						UID:       k8stypes.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
						Name:      "conf-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"jwt-config":        []byte(jwtPluginConfig),
						"basic-auth-config": []byte(basicAuthPluginConfig),
					},
				},
			}
			store, err := store.NewFakeStore(objects)
			require.NoError(t, err)
			p := mustNewParser(t, store)
			result := p.BuildKongConfig()
			require.Empty(t, result.TranslationFailures)
			state := result.KongState
			require.NotNil(t, state)
			assert.Equal(0, len(state.Plugins),
				"expected no plugins to be rendered")
		})

	t.Run("secretToConfiguration handles valid configuration and "+
		"discards invalid configuration", func(t *testing.T) {
		objects := stock
		jwtPluginConfig := `{"run_on_preflight": false}`  // JSON
		basicAuthPluginConfig := "hide_credentials: true" // YAML
		badJwtPluginConfig := "22222"                     // not JSON
		badBasicAuthPluginConfig := "111111"              // not YAML
		objects.Secrets = []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       k8stypes.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
					Name:      "conf-secret",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"jwt-config":            []byte(jwtPluginConfig),
					"basic-auth-config":     []byte(basicAuthPluginConfig),
					"bad-jwt-config":        []byte(badJwtPluginConfig),
					"bad-basic-auth-config": []byte(badBasicAuthPluginConfig),
				},
			},
		}
		references := []*kongv1.SecretValueFromSource{
			{
				Secret: "conf-secret",
				Key:    "jwt-config",
			},
			{
				Secret: "conf-secret",
				Key:    "basic-auth-config",
			},
		}
		badReferences := []*kongv1.SecretValueFromSource{
			{
				Secret: "conf-secret",
				Key:    "bad-basic-auth-config",
			},
			{
				Secret: "conf-secret",
				Key:    "bad-jwt-config",
			},
		}
		store, err := store.NewFakeStore(objects)
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		for _, testcase := range references {
			config, err := kongstate.SecretToConfiguration(store, *testcase, "default")
			assert.NotEmpty(config)
			require.NoError(t, err)
		}
		for _, testcase := range badReferences {
			config, err := kongstate.SecretToConfiguration(store, *testcase, "default")
			assert.Empty(config)
			assert.NotEmpty(err)
		}
	})
	t.Run("plugins with unparsable configuration are not constructed",
		func(t *testing.T) {
			jwtPluginConfig := "22222"        // not JSON
			basicAuthPluginConfig := "111111" // not YAML
			objects := stock
			objects.KongPlugins = []*kongv1.KongPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "global-foo-plugin",
						Namespace: "default",
						Labels: map[string]string{
							"global": "true",
						},
					},
					PluginName: "jwt",
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Key:    "missing-key",
							Secret: "conf-secret",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-plugin",
						Namespace: "default",
					},
					PluginName: "jwt",
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Key:    "missing-key",
							Secret: "conf-secret",
						},
					},
				},
			}
			objects.KongClusterPlugins = []*kongv1.KongClusterPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "global-bar-plugin",
						Labels: map[string]string{
							"global": "true",
						},
					},
					Protocols:  kongv1.StringsToKongProtocols([]string{"http"}),
					PluginName: "basic-auth",
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "missing-secret",
							Namespace: "default",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "bar-plugin",
					},
					Protocols:  kongv1.StringsToKongProtocols([]string{"http"}),
					PluginName: "basic-auth",
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "missing-secret",
							Namespace: "default",
						},
					},
				},
			}
			objects.Secrets = []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						UID:       k8stypes.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
						Name:      "conf-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"jwt-config":        []byte(jwtPluginConfig),
						"basic-auth-config": []byte(basicAuthPluginConfig),
					},
				},
			}
			store, err := store.NewFakeStore(objects)
			require.NoError(t, err)
			p := mustNewParser(t, store)
			result := p.BuildKongConfig()
			require.Empty(t, result.TranslationFailures)
			state := result.KongState
			require.NotNil(t, state)
			assert.Equal(0, len(state.Plugins),
				"expected no plugins to be rendered")
		})
}

func TestCACertificate(t *testing.T) {
	assert := assert.New(t)
	caCert1, _ := certificate.MustGenerateSelfSignedCertPEMFormat(certificate.WithCATrue())
	t.Run("valid CACertificate is processed", func(t *testing.T) {
		secrets := []*corev1.Secret{
			{
				TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: corev1.SchemeGroupVersion.String()},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Data: map[string][]byte{
					"id":   []byte("8214a145-a328-4c56-ab72-2973a56d4eae"),
					"cert": caCert1,
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			Secrets: secrets,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		assert.Len(state.CACertificates, 1)
		// Parser tests do not check tags, these are tested independently.
		state.CACertificates[0].Tags = nil
		assert.Equal(kong.CACertificate{
			ID:   kong.String("8214a145-a328-4c56-ab72-2973a56d4eae"),
			Cert: kong.String(string(caCert1)),
		}, state.CACertificates[0])
	})

	caCert2, _ := certificate.MustGenerateSelfSignedCertPEMFormat(certificate.WithCATrue())
	t.Run("multiple CACertificates are processed", func(t *testing.T) {
		secrets := []*corev1.Secret{
			{
				TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: corev1.SchemeGroupVersion.String()},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Data: map[string][]byte{
					"id":   []byte("8214a145-a328-4c56-ab72-2973a56d4eae"),
					"cert": caCert1,
				},
			},
			{
				TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: corev1.SchemeGroupVersion.String()},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "non-default",
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Data: map[string][]byte{
					"id":   []byte("570c28aa-e784-43c1-8ec7-ae7f4ce40189"),
					"cert": caCert2,
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			Secrets: secrets,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		assert.Len(state.CACertificates, 2)
	})

	expiredCACert, _ := certificate.MustGenerateSelfSignedCertPEMFormat(certificate.WithCATrue(), certificate.WithAlreadyExpired())
	t.Run("invalid CACertificates are ignored", func(t *testing.T) {
		secrets := []*corev1.Secret{
			{
				TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: corev1.SchemeGroupVersion.String()},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "valid-cert",
					Namespace: "default",
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Data: map[string][]byte{
					"id":   []byte("8214a145-a328-4c56-ab72-2973a56d4eae"),
					"cert": caCert1,
				},
			},
			{
				TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: corev1.SchemeGroupVersion.String()},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "missing-cert-key",
					Namespace: "non-default",
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Data: map[string][]byte{
					"id": []byte("570c28aa-e784-43c1-8ec7-ae7f4ce40189"),
					// cert is missing
				},
			},
			{
				TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: corev1.SchemeGroupVersion.String()},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "missing-id-key",
					Namespace: "non-default",
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Data: map[string][]byte{
					// id is missing
					"cert": caCert2,
				},
			},
			{
				TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: corev1.SchemeGroupVersion.String()},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "expired-cert",
					Namespace: "non-default",
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Data: map[string][]byte{
					"id":   []byte("670c28aa-e784-43c1-8ec7-ae7f4ce40189"),
					"cert": expiredCACert,
				},
			},
			{
				TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: corev1.SchemeGroupVersion.String()},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "invalid-cert",
					Namespace: "non-default",
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Data: map[string][]byte{
					"id":   []byte("770c28aa-e784-43c1-8ec7-ae7f4ce40189"),
					"cert": []byte("invalid-cert"),
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			Secrets: secrets,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		assert.Len(result.TranslationFailures, 4)
		state := result.KongState
		require.NotNil(t, state)

		assert.Len(state.CACertificates, 1)
		// parser tests do not check tags, these are tested independently
		state.CACertificates[0].Tags = nil
		assert.Equal(kong.CACertificate{
			ID:   kong.String("8214a145-a328-4c56-ab72-2973a56d4eae"),
			Cert: kong.String(string(caCert1)),
		}, state.CACertificates[0])
	})
}

func TestServiceClientCertificate(t *testing.T) {
	assert := assert.New(t)
	t.Run("valid client-cert annotation", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
										{
											Path: "/bar",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "bar-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
		}
		crt, key := certificate.MustGenerateSelfSignedCertPEMFormat()
		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       k8stypes.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
					Name:      "secret1",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"tls.crt": crt,
					"tls.key": key,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       k8stypes.UID("ffaabbcc-180b-4702-a91f-61351a33c6e4"),
					Name:      "secret2",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"tls.crt": crt,
					"tls.key": key,
				},
			},
		}
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
					Annotations: map[string]string{
						"konghq.com/client-cert": "secret1",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar-svc",
					Namespace: "default",
					Annotations: map[string]string{
						"konghq.com/client-cert": "secret2",
						"konghq.com/protocol":    "https",
					},
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Secrets:     secrets,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		assert.Equal(1, len(state.Certificates),
			"expected one certificates to be rendered")
		assert.Equal("7428fb98-180b-4702-a91f-61351a33c6e4",
			*state.Certificates[0].ID)

		assert.Equal(2, len(state.Services))
		assert.Equal("7428fb98-180b-4702-a91f-61351a33c6e4",
			*state.Services[0].ClientCertificate.ID)
		assert.Equal("7428fb98-180b-4702-a91f-61351a33c6e4",
			*state.Services[1].ClientCertificate.ID)
	})
	t.Run("client-cert secret doesn't exist", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: corev1.SchemeGroupVersion.String()},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
					Annotations: map[string]string{
						"konghq.com/client-cert": "secret1",
					},
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Len(t, result.TranslationFailures, 1)
		state := result.KongState
		require.NotNil(t, state)
		assert.Equal(0, len(state.Certificates),
			"expected no certificates to be rendered")

		assert.Equal(1, len(state.Services))
		assert.Nil(state.Services[0].ClientCertificate)
	})
}

func TestKongRouteAnnotations(t *testing.T) {
	t.Run("strip-path annotation is correctly processed (true)", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						"konghq.com/strip-path":     "trUe",
						annotations.IngressClassKey: "kong",
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		require.Len(t, state.Services, 1, "expected one service to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Service.Tags = nil
		assert.Equal(t, kong.Service{
			Name:           kong.String("default.foo-svc.pnum-80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
			ID:             kong.String("cce764ba-6924-5519-92ab-17c15418bbec"),
		}, state.Services[0].Service)

		require.Len(t, state.Services[0].Routes, 1, "expected one route to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.bar.00"),
			StripPath:         kong.Bool(true),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			RegexPriority:     kong.Int(100),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			ID:                kong.String("0977b208-1478-5280-90cf-037af9ffc3ef"),
		}, state.Services[0].Routes[0].Route)
	})
	t.Run("strip-path annotation is correctly processed (false)", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: "kong",
						"konghq.com/strip-path":     "false",
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		assert.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		require.Len(t, state.Services, 1, "expected one service to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Service.Tags = nil
		assert.Equal(t, kong.Service{
			Name:           kong.String("default.foo-svc.pnum-80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
			ID:             kong.String("cce764ba-6924-5519-92ab-17c15418bbec"),
		}, state.Services[0].Service)

		require.Len(t, state.Services[0].Routes, 1, "expected one route to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.bar.00"),
			StripPath:         kong.Bool(false),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			RegexPriority:     kong.Int(100),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			ID:                kong.String("0977b208-1478-5280-90cf-037af9ffc3ef"),
		}, state.Services[0].Routes[0].Route)
	})
	t.Run("https-redirect-status-code annotation is correctly processed", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey:             "kong",
						"konghq.com/https-redirect-status-code": "301",
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		require.Len(t, state.Services, 1, "expected one service to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Service.Tags = nil
		assert.Equal(t, kong.Service{
			Name:           kong.String("default.foo-svc.pnum-80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
			ID:             kong.String("cce764ba-6924-5519-92ab-17c15418bbec"),
		}, state.Services[0].Service)

		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Len(t, state.Services[0].Routes, 1, "expected one route to be rendered")
		assert.Equal(t, kong.Route{
			Name:                    kong.String("default.bar.00"),
			StripPath:               kong.Bool(false),
			HTTPSRedirectStatusCode: kong.Int(301),
			Hosts:                   kong.StringSlice("example.com"),
			PreserveHost:            kong.Bool(true),
			Paths:                   kong.StringSlice("/"),
			Protocols:               kong.StringSlice("http", "https"),
			RegexPriority:           kong.Int(100),
			ResponseBuffering:       kong.Bool(true),
			RequestBuffering:        kong.Bool(true),
			ID:                      kong.String("0977b208-1478-5280-90cf-037af9ffc3ef"),
		}, state.Services[0].Routes[0].Route)
	})

	t.Run("bad https-redirect-status-code annotation is ignored",
		func(t *testing.T) {
			ingresses := []*netv1.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							annotations.IngressClassKey:             "kong",
							"konghq.com/https-redirect-status-code": "whoops",
						},
					},
					Spec: netv1.IngressSpec{
						Rules: []netv1.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: netv1.IngressRuleValue{
									HTTP: &netv1.HTTPIngressRuleValue{
										Paths: []netv1.HTTPIngressPath{
											{
												Path: "/",
												Backend: netv1.IngressBackend{
													Service: &netv1.IngressServiceBackend{
														Name: "foo-svc",
														Port: netv1.ServiceBackendPort{
															Number: 80,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}

			services := []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-svc",
						Namespace: "default",
					},
				},
			}
			store, err := store.NewFakeStore(store.FakeObjects{
				IngressesV1: ingresses,
				Services:    services,
			})
			require.NoError(t, err)
			p := mustNewParser(t, store)
			result := p.BuildKongConfig()
			require.Empty(t, result.TranslationFailures)
			state := result.KongState
			require.NotNil(t, state)

			require.Len(t, state.Services, 1, "expected one service to be rendered")
			// parser tests do not check tags, these are tested independently
			state.Services[0].Service.Tags = nil
			assert.Equal(t, kong.Service{
				Name:           kong.String("default.foo-svc.pnum-80"),
				Host:           kong.String("foo-svc.default.80.svc"),
				Path:           kong.String("/"),
				Port:           kong.Int(80),
				ConnectTimeout: kong.Int(60000),
				ReadTimeout:    kong.Int(60000),
				WriteTimeout:   kong.Int(60000),
				Retries:        kong.Int(5),
				Protocol:       kong.String("http"),
				ID:             kong.String("cce764ba-6924-5519-92ab-17c15418bbec"),
			}, state.Services[0].Service)

			require.Len(t, state.Services[0].Routes, 1, "expected one route to be rendered")
			// parser tests do not check tags, these are tested independently
			state.Services[0].Routes[0].Route.Tags = nil
			assert.Equal(t, kong.Route{
				Name:              kong.String("default.bar.00"),
				StripPath:         kong.Bool(false),
				Hosts:             kong.StringSlice("example.com"),
				PreserveHost:      kong.Bool(true),
				Paths:             kong.StringSlice("/"),
				Protocols:         kong.StringSlice("http", "https"),
				RegexPriority:     kong.Int(100),
				ResponseBuffering: kong.Bool(true),
				RequestBuffering:  kong.Bool(true),
				ID:                kong.String("0977b208-1478-5280-90cf-037af9ffc3ef"),
			}, state.Services[0].Routes[0].Route)
		})
	t.Run("preserve-host annotation is correctly processed", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						"konghq.com/preserve-host":  "faLsE",
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		require.Len(t, state.Services, 1, "expected one service to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Service.Tags = nil
		assert.Equal(t, kong.Service{
			Name:           kong.String("default.foo-svc.pnum-80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
			ID:             kong.String("cce764ba-6924-5519-92ab-17c15418bbec"),
		}, state.Services[0].Service)

		require.Len(t, state.Services[0].Routes, 1, "expected one route to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.bar.00"),
			StripPath:         kong.Bool(false),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(false),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			RegexPriority:     kong.Int(100),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			ID:                kong.String("0977b208-1478-5280-90cf-037af9ffc3ef"),
		}, state.Services[0].Routes[0].Route)
	})

	t.Run("preserve-host annotation with random string is correctly processed", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
						"konghq.com/preserve-host":  "wiggle wiggle wiggle",
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		require.Len(t, state.Services, 1, "expected one service to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Service.Tags = nil
		assert.Equal(t, kong.Service{
			Name:           kong.String("default.foo-svc.pnum-80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
			ID:             kong.String("cce764ba-6924-5519-92ab-17c15418bbec"),
		}, state.Services[0].Service)

		require.Len(t, state.Services[0].Routes, 1, "expected one route to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.bar.00"),
			StripPath:         kong.Bool(false),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			RegexPriority:     kong.Int(100),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			ID:                kong.String("0977b208-1478-5280-90cf-037af9ffc3ef"),
		}, state.Services[0].Routes[0].Route)
	})

	t.Run("regex-priority annotation is correctly processed", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						"konghq.com/regex-priority": "10",
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		require.Len(t, state.Services, 1, "expected one service to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Service.Tags = nil
		assert.Equal(t, kong.Service{
			Name:           kong.String("default.foo-svc.pnum-80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
			ID:             kong.String("cce764ba-6924-5519-92ab-17c15418bbec"),
		}, state.Services[0].Service)

		require.Len(t, state.Services[0].Routes, 1, "expected one route to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.bar.00"),
			StripPath:         kong.Bool(false),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			RegexPriority:     kong.Int(10),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			ID:                kong.String("0977b208-1478-5280-90cf-037af9ffc3ef"),
		}, state.Services[0].Routes[0].Route)
	})

	t.Run("non-integer regex-priority annotation is ignored", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						"konghq.com/regex-priority": "IAmAString",
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		require.Len(t, state.Services, 1, "expected one service to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Service.Tags = nil
		assert.Equal(t, kong.Service{
			Name:           kong.String("default.foo-svc.pnum-80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
			ID:             kong.String("cce764ba-6924-5519-92ab-17c15418bbec"),
		}, state.Services[0].Service)

		require.Len(t, state.Services[0].Routes, 1, "expected one route to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.bar.00"),
			StripPath:         kong.Bool(false),
			RegexPriority:     kong.Int(100),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			ID:                kong.String("0977b208-1478-5280-90cf-037af9ffc3ef"),
		}, state.Services[0].Routes[0].Route)
	})

	t.Run("route buffering options are processed (true)", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route-buffering-test",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey:     annotations.DefaultIngressClass,
						"konghq.com/request-buffering":  "True",
						"konghq.com/response-buffering": "True",
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		assert.Len(t, state.Services, 1, "expected one service to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Service.Tags = nil
		assert.Equal(t, kong.Service{
			Name:           kong.String("default.foo-svc.pnum-80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
			ID:             kong.String("cce764ba-6924-5519-92ab-17c15418bbec"),
		}, state.Services[0].Service)

		assert.Len(t, state.Services[0].Routes, 1, "expected one route to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.route-buffering-test.00"),
			StripPath:         kong.Bool(false),
			RegexPriority:     kong.Int(100),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			RequestBuffering:  kong.Bool(true),
			ResponseBuffering: kong.Bool(true),
			ID:                kong.String("7baaf482-722d-5f33-802f-098d31491846"),
		}, state.Services[0].Routes[0].Route)
	})
	t.Run("route buffering options are processed (false)", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route-buffering-test",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey:     annotations.DefaultIngressClass,
						"konghq.com/request-buffering":  "False",
						"konghq.com/response-buffering": "False",
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		assert.Len(t, state.Services, 1, "expected one service to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Service.Tags = nil
		assert.Equal(t, kong.Service{
			Name:           kong.String("default.foo-svc.pnum-80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
			ID:             kong.String("cce764ba-6924-5519-92ab-17c15418bbec"),
		}, state.Services[0].Service)

		assert.Len(t, state.Services[0].Routes, 1, "expected one route to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.route-buffering-test.00"),
			StripPath:         kong.Bool(false),
			RegexPriority:     kong.Int(100),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			RequestBuffering:  kong.Bool(false),
			ResponseBuffering: kong.Bool(false),
			ID:                kong.String("7baaf482-722d-5f33-802f-098d31491846"),
		}, state.Services[0].Routes[0].Route)
	})
	t.Run("route buffering options are not processed with bad annotation values", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route-buffering-test",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey:     annotations.DefaultIngressClass,
						"konghq.com/request-buffering":  "invalid-value",
						"konghq.com/response-buffering": "invalid-value",
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		kongTrue := kong.Bool(true)
		assert.Len(t, state.Services, 1, "expected one service to be rendered")
		assert.Len(t, state.Services[0].Routes, 1, "expected one route to be rendered")
		assert.Equal(t, kongTrue, state.Services[0].Routes[0].Route.RequestBuffering)
		assert.Equal(t, kongTrue, state.Services[0].Routes[0].Route.ResponseBuffering)
	})
}

func TestKongProcessClasslessIngress(t *testing.T) {
	assert := assert.New(t)
	t.Run("Kong classless ingress evaluated (true)", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		require.Len(t, state.Services, 1, "expected one service to be rendered")
	})
	t.Run("Kong classless ingress evaluated (false)", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		assert.Equal(0, len(state.Services),
			"expected zero service to be rendered")
	})
}

func TestKnativeIngressAndPlugins(t *testing.T) {
	t.Run("knative ingress annotated with konghq.com/override", func(t *testing.T) {
		ingresses := []*knative.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "knative-ingress-with-override",
					Namespace: "foo-ns",
					Annotations: map[string]string{
						"networking.knative.dev/ingress-class":                      annotations.DefaultIngressClass,
						annotations.AnnotationPrefix + annotations.ConfigurationKey: "https-only",
					},
				},
				Spec: knative.IngressSpec{
					Rules: []knative.IngressRule{
						{
							Hosts: []string{"my-func.example.com"},
							HTTP: &knative.HTTPIngressRuleValue{
								Paths: []knative.HTTPIngressPath{
									{
										Path: "/",
										AppendHeaders: map[string]string{
											"foo": "bar",
										},
										Splits: []knative.IngressBackendSplit{
											{
												IngressBackend: knative.IngressBackend{
													ServiceNamespace: "foo-ns",
													ServiceName:      "foo-svc",
													ServicePort:      intstr.FromInt(42),
												},
												Percent: 100,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		kongIngresses := []*kongv1.KongIngress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "https-only",
					Namespace: "foo-ns",
				},
				Route: &kongv1.KongIngressRoute{
					Protocols:               kongv1.ProtocolSlice("https"),
					HTTPSRedirectStatusCode: kong.Int(308),
				},
			},
		}
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "foo-ns",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			KnativeIngresses: ingresses,
			KongIngresses:    kongIngresses,
			Services:         services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		require.Len(t, state.Services, 1, "expected one knative service")
		svc := state.Services[0]

		assert.Len(t, svc.Routes, 1, "expected one route in knative service")
		route := svc.Routes[0]

		assert.Equal(t, kong.StringSlice("https"), route.Protocols, "expected https protocol after override")
		assert.Equal(t, kong.Int(308), route.HTTPSRedirectStatusCode, "expected 308 after override")
	})
	t.Run("knative ingress without konghq.com/override", func(t *testing.T) {
		ingresses := []*knative.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "knative-ingress-without-override",
					Namespace: "foo-ns",
					Annotations: map[string]string{
						"networking.knative.dev/ingress-class": annotations.DefaultIngressClass,
					},
				},
				Spec: knative.IngressSpec{
					Rules: []knative.IngressRule{
						{
							Hosts: []string{"my-func.example.com"},
							HTTP: &knative.HTTPIngressRuleValue{
								Paths: []knative.HTTPIngressPath{
									{
										Path: "/",
										AppendHeaders: map[string]string{
											"foo": "bar",
										},
										Splits: []knative.IngressBackendSplit{
											{
												IngressBackend: knative.IngressBackend{
													ServiceNamespace: "foo-ns",
													ServiceName:      "foo-svc",
													ServicePort:      intstr.FromInt(42),
												},
												Percent: 100,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		kongIngresses := []*kongv1.KongIngress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "https-only",
					Namespace: "foo-ns",
				},
				Route: &kongv1.KongIngressRoute{
					Protocols:               kongv1.ProtocolSlice("https"),
					HTTPSRedirectStatusCode: kong.Int(308),
				},
			},
		}
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "foo-ns",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			KnativeIngresses: ingresses,
			KongIngresses:    kongIngresses,
			Services:         services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		require.Len(t, state.Services, 1, "expected one knative service")
		svc := state.Services[0]

		assert.Len(t, svc.Routes, 1, "expected one route in knative service")
		route := svc.Routes[0]

		assert.NotEqual(t, kong.StringSlice("https"), route.Protocols)
		assert.Nil(t, route.HTTPSRedirectStatusCode)
	})
	t.Run("knative ingress with multiple konghq.com annotations", func(t *testing.T) {
		ingresses := []*knative.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "knative-ingress-with-annotations",
					Namespace: "foo-ns",
					Annotations: map[string]string{
						"networking.knative.dev/ingress-class":                          annotations.DefaultIngressClass,
						annotations.AnnotationPrefix + annotations.ProtocolsKey:         "https",
						annotations.AnnotationPrefix + annotations.HTTPSRedirectCodeKey: "308",
						annotations.AnnotationPrefix + annotations.StripPathKey:         "true",
						annotations.AnnotationPrefix + annotations.MethodsKey:           "POST,PUT",
					},
				},
				Spec: knative.IngressSpec{
					Rules: []knative.IngressRule{
						{
							Hosts: []string{"my-func.example.com"},
							HTTP: &knative.HTTPIngressRuleValue{
								Paths: []knative.HTTPIngressPath{
									{
										Path: "/",
										AppendHeaders: map[string]string{
											"foo": "bar",
										},
										Splits: []knative.IngressBackendSplit{
											{
												IngressBackend: knative.IngressBackend{
													ServiceNamespace: "foo-ns",
													ServiceName:      "foo-svc",
													ServicePort:      intstr.FromInt(42),
												},
												Percent: 100,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "foo-ns",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			KnativeIngresses: ingresses,
			Services:         services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		require.Len(t, state.Services, 1, "expected one knative service")
		svc := state.Services[0]

		require.Len(t, svc.Routes, 1, "expected one route in knative service")
		route := svc.Routes[0]

		assert.Equal(t, kong.StringSlice("https"), route.Protocols, "expected https after konghq.com/protocols")
		assert.Equal(t, kong.Int(308), route.HTTPSRedirectStatusCode, "expected 308 after konghq.com/https-redirect-status-code")
		assert.Equal(t, kong.Bool(true), route.StripPath, "expected true after konghq.com/strip-path")
		assert.Equal(t, kong.StringSlice("POST", "PUT"), route.Methods, "expected POST and PUT after konghq.com/methods")
	})
	t.Run("knative ingress rule and service-level plugin", func(t *testing.T) {
		ingresses := []*knative.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "knative-ingress",
					Namespace: "foo-ns",
					Annotations: map[string]string{
						"networking.knative.dev/ingress-class": annotations.DefaultIngressClass,
					},
				},
				Spec: knative.IngressSpec{
					Rules: []knative.IngressRule{
						{
							Hosts: []string{"my-func.example.com"},
							HTTP: &knative.HTTPIngressRuleValue{
								Paths: []knative.HTTPIngressPath{
									{
										Path: "/",
										AppendHeaders: map[string]string{
											"foo": "bar",
										},
										Splits: []knative.IngressBackendSplit{
											{
												IngressBackend: knative.IngressBackend{
													ServiceNamespace: "foo-ns",
													ServiceName:      "foo-svc",
													ServicePort:      intstr.FromInt(42),
												},
												Percent: 100,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "foo-ns",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "knative-key-auth",
						"networking.knative.dev/ingress-class":                annotations.DefaultIngressClass,
					},
				},
			},
		}
		plugins := []*kongv1.KongPlugin{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "knative-key-auth",
					Namespace: "foo-ns",
				},
				PluginName: "key-auth",
				Protocols:  kongv1.StringsToKongProtocols([]string{"http"}),
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{"foo": "bar", "knative": "yo"}`),
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			KnativeIngresses: ingresses,
			Services:         services,
			KongPlugins:      plugins,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		require.Len(t, state.Services, 1, "expected one service to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Service.Tags = nil
		svc := state.Services[0]

		assert.Equal(t, kong.Service{
			Name:           kong.String("foo-ns.foo-svc.42"),
			Host:           kong.String("foo-svc.foo-ns.42.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
			ID:             kong.String("c2ef7d3d-bb2a-5ae1-bc89-8b626ebab19e"),
		}, svc.Service)

		assert.Len(t, svc.Plugins, 1, "expected one request-transformer plugin")
		assert.Equal(t, kong.Plugin{
			Name: kong.String("request-transformer"),
			Config: kong.Configuration{
				"add": map[string]interface{}{
					"headers": []string{"foo:bar"},
				},
			},
		}, svc.Plugins[0])

		assert.Len(t, svc.Routes, 1, "expected one route to be rendered")
		// parser tests do not check tags, these are tested independently
		svc.Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("foo-ns.knative-ingress.00"),
			StripPath:         kong.Bool(false),
			Hosts:             kong.StringSlice("my-func.example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			RegexPriority:     kong.Int(0),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			ID:                kong.String("464fe8ea-acc9-58fd-8bb2-00f046287a2c"),
		}, svc.Routes[0].Route)

		assert.Len(t, state.Plugins, 1, "expected one key-auth plugin")
		// parser tests do not check tags, these are tested independently
		state.Plugins[0].Plugin.Tags = nil
		assert.Equal(t, kong.Plugin{
			Name: kong.String("key-auth"),
			Config: kong.Configuration{
				"foo":     "bar",
				"knative": "yo",
			},
			Service: &kong.Service{
				ID: kong.String("foo-ns.foo-svc.42"),
			},
			Protocols: kong.StringSlice("http"),
		}, state.Plugins[0].Plugin)
	})
}

func TestKongServiceAnnotations(t *testing.T) {
	t.Run("path annotation is correctly processed", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
					Annotations: map[string]string{
						"konghq.com/path": "/baz",
					},
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		assert.Len(t, state.Services, 1, "expected one service to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Service.Tags = nil
		assert.Equal(t, kong.Service{
			Name:           kong.String("default.foo-svc.pnum-80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/baz"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
			ID:             kong.String("cce764ba-6924-5519-92ab-17c15418bbec"),
		}, state.Services[0].Service)

		assert.Len(t, state.Services[0].Routes, 1, "expected one route to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.bar.00"),
			StripPath:         kong.Bool(false),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			RegexPriority:     kong.Int(100),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			ID:                kong.String("0977b208-1478-5280-90cf-037af9ffc3ef"),
		}, state.Services[0].Routes[0].Route)
	})

	t.Run("host-header annotation is correctly processed", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
					Annotations: map[string]string{
						"konghq.com/host-header": "example.com",
					},
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		assert.Len(t, state.Services, 1, "expected one service to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Service.Tags = nil
		assert.Equal(t, kong.Service{
			Name:           kong.String("default.foo-svc.pnum-80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
			ID:             kong.String("cce764ba-6924-5519-92ab-17c15418bbec"),
		}, state.Services[0].Service)

		assert.Len(t, state.Upstreams, 1, "expected one upstream to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Upstreams[0].Upstream.Tags = nil
		assert.Equal(t, kong.Upstream{
			Name:       kong.String("foo-svc.default.80.svc"),
			HostHeader: kong.String("example.com"),
		}, state.Upstreams[0].Upstream)

		assert.Len(t, state.Services[0].Routes, 1, "expected one route to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.bar.00"),
			StripPath:         kong.Bool(false),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			RegexPriority:     kong.Int(100),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			ID:                kong.String("0977b208-1478-5280-90cf-037af9ffc3ef"),
		}, state.Services[0].Routes[0].Route)
	})

	t.Run("methods annotation is correctly processed", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						"konghq.com/methods":        "POST,GET",
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)

		require.Len(t, state.Services, 1, "expected one service to be rendered")
		// parser tests do not check tags, these are tested independently
		state.Services[0].Service.Tags = nil
		assert.Equal(t, kong.Service{
			Name:           kong.String("default.foo-svc.pnum-80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
			ID:             kong.String("cce764ba-6924-5519-92ab-17c15418bbec"),
		}, state.Services[0].Service)

		require.Len(t, state.Services[0].Routes, 1, "expected one route to be rendered")
		// parser tests do not check tags, these are tested independently
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.bar.00"),
			StripPath:         kong.Bool(false),
			RegexPriority:     kong.Int(100),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			Methods:           kong.StringSlice("POST", "GET"),
			ID:                kong.String("0977b208-1478-5280-90cf-037af9ffc3ef"),
			Tags: []*string{
				kong.String("k8s-name:bar"),
				kong.String("k8s-namespace:default"),
			},
		}, state.Services[0].Routes[0].Route)
	})
}

func TestDefaultBackend(t *testing.T) {
	assert := assert.New(t)
	t.Run("default backend is processed correctly", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ing-with-default-backend",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					DefaultBackend: &netv1.IngressBackend{
						Service: &netv1.IngressServiceBackend{
							Name: "default-svc",
							Port: netv1.ServiceBackendPort{
								Number: 80,
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "default-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		require.Len(t, state.Services, 1, "expected one service to be rendered")
		assert.Equal("default.default-svc.80", *state.Services[0].Name)
		assert.Equal("default-svc.default.80.svc", *state.Services[0].Host)
		assert.Equal(1, len(state.Services[0].Routes),
			"expected one routes to be rendered")
		assert.Equal("default.ing-with-default-backend", *state.Services[0].Routes[0].Name)
		assert.Equal("/", *state.Services[0].Routes[0].Paths[0])
	})

	t.Run("client-cert secret doesn't exist", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: corev1.SchemeGroupVersion.String()},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
					Annotations: map[string]string{
						"konghq.com/client-cert": "secret1",
					},
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Len(t, result.TranslationFailures, 1)
		state := result.KongState
		require.NotNil(t, state)
		assert.Equal(0, len(state.Certificates),
			"expected no certificates to be rendered")

		assert.Equal(1, len(state.Services))
		assert.Nil(state.Services[0].ClientCertificate)
	})
}

func TestParserSecret(t *testing.T) {
	assert := assert.New(t)
	t.Run("invalid TLS secret", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"bar.com"},
						},
					},
				},
			},
		}

		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret1",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(""),
					"tls.key": []byte(""),
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Secrets:     secrets,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		assert.Equal(0, len(state.Certificates),
			"expected no certificates to be rendered with empty secret")
	})

	crt, key := certificate.MustGenerateSelfSignedCertPEMFormat()
	t.Run("duplicate certificates order by time", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "ns1",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret2",
							Hosts:      []string{"bar.com"},
						},
					},
				},
			},
		}

		t1, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
		t2, _ := time.Parse(time.RFC3339, "2006-01-02T15:05:05Z")
		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       "3e8edeca-7d23-4e02-84c9-437d11b746a6",
					Name:      "secret1",
					Namespace: "default",
					CreationTimestamp: metav1.Time{
						Time: t1,
					},
				},
				Data: map[string][]byte{
					"tls.crt": crt,
					"tls.key": key,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       "fc28a22c-41e1-4cd6-9099-fd7756ffe58e",
					Name:      "secret2",
					Namespace: "ns1",
					CreationTimestamp: metav1.Time{
						Time: t2,
					},
				},
				Data: map[string][]byte{
					"tls.crt": crt,
					"tls.key": key,
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Secrets:     secrets,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		assert.Len(state.Certificates, 1, "certificates are de-duplicated")

		sort.SliceStable(state.Certificates[0].SNIs, func(i, j int) bool {
			return strings.Compare(*state.Certificates[0].SNIs[i],
				*state.Certificates[0].SNIs[j]) > 0
		})
		// parser tests do not check tags, these are tested independently
		state.Certificates[0].Tags = nil
		assert.Equal(kongstate.Certificate{
			Certificate: kong.Certificate{
				ID:   kong.String("3e8edeca-7d23-4e02-84c9-437d11b746a6"),
				Cert: kong.String(strings.TrimSpace(string(crt))),
				Key:  kong.String(strings.TrimSpace(string(key))),
				SNIs: kong.StringSlice("foo.com", "bar.com"),
			},
		}, state.Certificates[0])
	})
	t.Run("duplicate certificates order by uid", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "ns1",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret2",
							Hosts:      []string{"bar.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "baz",
					Namespace: "ns2",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret3",
							Hosts:      []string{"baz.com"},
						},
					},
				},
			},
		}

		t1, _ := time.Parse(time.RFC3339, "2006-01-02T15:05:05Z")
		t2, _ := time.Parse(time.RFC3339, "2006-01-02T15:05:05Z")
		t3, _ := time.Parse(time.RFC3339, "2006-01-02T15:06:05Z")
		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       "3c28a22c-41e1-4cd6-9099-fd7756ffe58e",
					Name:      "secret1",
					Namespace: "default",
					CreationTimestamp: metav1.Time{
						Time: t1,
					},
				},
				Data: map[string][]byte{
					"tls.crt": crt,
					"tls.key": key,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       "2c28a22c-41e1-4cd6-9099-fd7756ffe58e",
					Name:      "secret2",
					Namespace: "ns1",
					CreationTimestamp: metav1.Time{
						Time: t2,
					},
				},
				Data: map[string][]byte{
					"tls.crt": crt,
					"tls.key": key,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       "1c28a22c-41e1-4cd6-9099-fd7756ffe58e",
					Name:      "secret3",
					Namespace: "ns2",
					CreationTimestamp: metav1.Time{
						Time: t3,
					},
				},
				Data: map[string][]byte{
					"tls.crt": crt,
					"tls.key": key,
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Secrets:     secrets,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		assert.Equal(1, len(state.Certificates),
			"certificates are de-duplicated")

		sort.SliceStable(state.Certificates[0].SNIs, func(i, j int) bool {
			return strings.Compare(*state.Certificates[0].SNIs[i],
				*state.Certificates[0].SNIs[j]) > 0
		})
		// parser tests do not check tags, these are tested independently
		state.Certificates[0].Tags = nil
		assert.Equal(kongstate.Certificate{
			Certificate: kong.Certificate{
				ID:   kong.String("2c28a22c-41e1-4cd6-9099-fd7756ffe58e"),
				Cert: kong.String(strings.TrimSpace(string(crt))),
				Key:  kong.String(strings.TrimSpace(string(key))),
				SNIs: kong.StringSlice("foo.com", "baz.com", "bar.com"),
			},
		}, state.Certificates[0])
	})
	t.Run("duplicate SNIs", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "ns1",
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret2",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
		}

		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret1",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"tls.crt": crt,
					"tls.key": key,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret2",
					Namespace: "ns1",
				},
				Data: map[string][]byte{
					"tls.crt": crt,
					"tls.key": key,
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Secrets:     secrets,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		assert.Len(state.Certificates, 1, "SNIs are de-duplicated")
	})
}

func TestParserSNI(t *testing.T) {
	t.Run("route includes SNI when TLS info present, but not for wildcard hostnames", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"example.com", "*.example.com"},
						},
					},
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Host: "*.example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		crt, key := certificate.MustGenerateSelfSignedCertPEMFormat()
		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret1",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"tls.crt": crt,
					"tls.key": key,
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Secrets:     secrets,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		state.Services[0].Routes[1].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.foo.00"),
			StripPath:         kong.Bool(false),
			RegexPriority:     kong.Int(100),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			ID:                kong.String("95b19369-8def-5255-a40f-dd459cd689e8"),
		}, state.Services[0].Routes[0].Route)
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.foo.10"),
			StripPath:         kong.Bool(false),
			RegexPriority:     kong.Int(100),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			Hosts:             kong.StringSlice("*.example.com"),
			SNIs:              nil,
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			ID:                kong.String("9c78e74b-423f-5d44-abad-29d1cd137520"),
		}, state.Services[0].Routes[1].Route)
	})
	t.Run("route does not include SNI when TLS info absent", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.foo.00"),
			StripPath:         kong.Bool(false),
			RegexPriority:     kong.Int(100),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			Hosts:             kong.StringSlice("example.com"),
			SNIs:              nil,
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			ID:                kong.String("95b19369-8def-5255-a40f-dd459cd689e8"),
		}, state.Services[0].Routes[0].Route)
	})
}

func TestParserHostAliases(t *testing.T) {
	annHostAliasesKey := annotations.AnnotationPrefix + annotations.HostAliasesKey
	t.Run("route Hosts includes Host-Aliases when Host-Aliases are present", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
						annHostAliasesKey:           "*.example.com,*.sample.com,*.illustration.com",
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.foo.00"),
			StripPath:         kong.Bool(false),
			RegexPriority:     kong.Int(100),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			Hosts:             kong.StringSlice("example.com", "*.example.com", "*.sample.com", "*.illustration.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			ID:                kong.String("95b19369-8def-5255-a40f-dd459cd689e8"),
		}, state.Services[0].Routes[0].Route)
	})
	t.Run("route Hosts remain unmodified when Host-Aliases are not present", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.foo.00"),
			StripPath:         kong.Bool(false),
			RegexPriority:     kong.Int(100),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			ID:                kong.String("95b19369-8def-5255-a40f-dd459cd689e8"),
		}, state.Services[0].Routes[0].Route)
	})
	t.Run("route Hosts will not contain duplicates when Host-Aliases duplicates the host", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
						annHostAliasesKey:           "example.com,*.example.com",
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		// parser tests do not check tags, these are tested independently
		state.Services[0].Routes[0].Route.Tags = nil
		assert.Equal(t, kong.Route{
			Name:              kong.String("default.foo.00"),
			StripPath:         kong.Bool(false),
			RegexPriority:     kong.Int(100),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			Hosts:             kong.StringSlice("example.com", "*.example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			ID:                kong.String("95b19369-8def-5255-a40f-dd459cd689e8"),
		}, state.Services[0].Routes[0].Route)
	})
}

func TestPluginAnnotations(t *testing.T) {
	assert := assert.New(t)
	t.Run("simple association", func(t *testing.T) {
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "foo-plugin",
						annotations.IngressClassKey:                           annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		plugins := []*kongv1.KongPlugin{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-plugin",
					Namespace: "default",
				},
				PluginName: "key-auth",
				Protocols:  []kongv1.KongProtocol{"grpc"},
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{
					"foo": "bar",
					"add": {
						"headers": [
							"header1:value1",
							"header2:value2"
							]
						}
					}`),
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Services:    services,
			KongPlugins: plugins,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		assert.Equal(1, len(state.Plugins),
			"expected no plugins to be rendered with missing plugin")
		pl := state.Plugins[0].Plugin
		pl.Route = nil
		// parser tests do not check tags, these are tested independently
		pl.Tags = nil
		assert.Equal(pl, kong.Plugin{
			Name:      kong.String("key-auth"),
			Protocols: kong.StringSlice("grpc"),
			Config: kong.Configuration{
				"foo": "bar",
				"add": map[string]interface{}{
					"headers": []interface{}{
						"header1:value1",
						"header2:value2",
					},
				},
			},
		})
	})
	t.Run("KongPlugin takes precedence over KongPlugin", func(t *testing.T) {
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "foo-plugin",
						annotations.IngressClassKey:                           annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		clusterPlugins := []*kongv1.KongClusterPlugin{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-plugin",
					Namespace: "default",
				},
				PluginName: "basic-auth",
				Protocols:  []kongv1.KongProtocol{"grpc"},
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{"foo": "bar"}`),
				},
			},
		}
		plugins := []*kongv1.KongPlugin{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-plugin",
					Namespace: "default",
				},
				PluginName: "key-auth",
				Protocols:  []kongv1.KongProtocol{"grpc"},
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{"foo": "bar"}`),
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1:        ingresses,
			Services:           services,
			KongPlugins:        plugins,
			KongClusterPlugins: clusterPlugins,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		assert.Equal(1, len(state.Plugins),
			"expected no plugins to be rendered with missing plugin")
		assert.Equal("key-auth", *state.Plugins[0].Name)
		assert.Equal("grpc", *state.Plugins[0].Protocols[0])
	})
	t.Run("KongClusterPlugin association", func(t *testing.T) {
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "foo-plugin",
						annotations.IngressClassKey:                           annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		clusterPlugins := []*kongv1.KongClusterPlugin{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-plugin",
					Namespace: "default",
				},
				PluginName: "basic-auth",
				Protocols:  []kongv1.KongProtocol{"grpc"},
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{"foo": "bar"}`),
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1:        ingresses,
			Services:           services,
			KongClusterPlugins: clusterPlugins,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		assert.Equal(1, len(state.Plugins),
			"expected no plugins to be rendered with missing plugin")
		assert.Equal("basic-auth", *state.Plugins[0].Name)
		assert.Equal("grpc", *state.Plugins[0].Protocols[0])
	})
	t.Run("missing plugin", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "does-not-exist",
						annotations.IngressClassKey:                           annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "foo-svc",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		assert.Equal(0, len(state.Plugins),
			"expected no plugins to be rendered with missing plugin")
	})
}

func TestGetEndpoints(t *testing.T) {
	tests := []struct {
		name              string
		svc               *corev1.Service
		port              *corev1.ServicePort
		proto             corev1.Protocol
		fn                func(string, string) ([]*discoveryv1.EndpointSlice, error)
		result            []util.Endpoint
		isServiceUpstream bool
	}{
		{
			name:  "no service should return 0 endpoints",
			svc:   nil,
			port:  nil,
			proto: corev1.ProtocolTCP,
			fn: func(string, string) ([]*discoveryv1.EndpointSlice, error) {
				return nil, nil
			},
			result: []util.Endpoint{},
		},
		{
			name:  "no service port should return 0 endpoints",
			svc:   &corev1.Service{},
			port:  nil,
			proto: corev1.ProtocolTCP,
			fn: func(string, string) ([]*discoveryv1.EndpointSlice, error) {
				return nil, nil
			},
			result: []util.Endpoint{},
		},
		{
			name:  "a service without endpoints should return 0 endpoints",
			svc:   &corev1.Service{},
			port:  &corev1.ServicePort{Name: "default"},
			proto: corev1.ProtocolTCP,
			fn: func(string, string) ([]*discoveryv1.EndpointSlice, error) {
				return []*discoveryv1.EndpointSlice{}, nil
			},
			result: []util.Endpoint{},
		},
		{
			name: "a service type ServiceTypeExternalName with a valid port should return one endpoint",
			svc: &corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:         corev1.ServiceTypeExternalName,
					ExternalName: "10.0.0.1.xip.io",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			port: &corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			proto: corev1.ProtocolTCP,
			fn: func(string, string) ([]*discoveryv1.EndpointSlice, error) {
				return []*discoveryv1.EndpointSlice{}, nil
			},
			result: []util.Endpoint{
				{
					Address: "10.0.0.1.xip.io",
					Port:    "80",
				},
			},
		},
		{
			name: "a service with ingress.kubernetes.io/service-upstream annotation should return one endpoint",
			svc: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
					Annotations: map[string]string{
						"ingress.kubernetes.io/service-upstream": "true",
					},
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeClusterIP,
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
							Port:       2080,
						},
					},
				},
			},
			port: &corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
				Port:       2080,
			},
			proto: corev1.ProtocolTCP,
			fn: func(string, string) ([]*discoveryv1.EndpointSlice, error) {
				return []*discoveryv1.EndpointSlice{}, nil
			},
			result: []util.Endpoint{
				{
					Address: "foo.bar.svc",
					Port:    "2080",
				},
			},
		},
		{
			name: "a service with configured IngressClassParameters as ServiceUpstream should return one endpoint",
			svc: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeClusterIP,
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
							Port:       2080,
						},
					},
				},
			},
			port: &corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
				Port:       2080,
			},
			proto: corev1.ProtocolTCP,
			fn: func(string, string) ([]*discoveryv1.EndpointSlice, error) {
				return []*discoveryv1.EndpointSlice{}, nil
			},
			result: []util.Endpoint{
				{
					Address: "foo.bar.svc",
					Port:    "2080",
				},
			},
			isServiceUpstream: true,
		},
		{
			name: "should return no endpoints when there is an error searching for endpoints",
			svc: &corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			port: &corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
				Port:       2080,
			},
			proto: corev1.ProtocolTCP,
			fn: func(string, string) ([]*discoveryv1.EndpointSlice, error) {
				return nil, fmt.Errorf("unexpected error")
			},
			result: []util.Endpoint{},
		},
		{
			name: "should return no endpoints when the protocol does not match",
			svc: &corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			port: &corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			proto: corev1.ProtocolTCP,
			fn: func(string, string) ([]*discoveryv1.EndpointSlice, error) {
				return []*discoveryv1.EndpointSlice{
					{
						Endpoints: []discoveryv1.Endpoint{
							{
								Addresses: []string{"1.1.1.1"},
								NodeName:  lo.ToPtr("dummy"),
							},
						},
						Ports: builder.NewEndpointPort(80).WithProtocol(corev1.ProtocolUDP).IntoSlice(),
					},
				}, nil
			},
			result: []util.Endpoint{},
		},
		{
			name: "should return no endpoints when there is no ready Addresses",
			svc: &corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			port: &corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			proto: corev1.ProtocolTCP,
			fn: func(string, string) ([]*discoveryv1.EndpointSlice, error) {
				return []*discoveryv1.EndpointSlice{
					{
						Endpoints: []discoveryv1.Endpoint{
							{
								Addresses: []string{"1.1.1.1"},
								NodeName:  lo.ToPtr("dummy"),
								Conditions: discoveryv1.EndpointConditions{
									Ready: lo.ToPtr(false),
								},
							},
						},
						Ports: []discoveryv1.EndpointPort{
							{
								Protocol: lo.ToPtr(corev1.ProtocolUDP),
							},
						},
					},
				}, nil
			},
			result: []util.Endpoint{},
		},
		{
			name: "should return no endpoints when the name of the port name do not match any port in the endpoint Subsets",
			svc: &corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			port: &corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			proto: corev1.ProtocolTCP,
			fn: func(string, string) ([]*discoveryv1.EndpointSlice, error) {
				return []*discoveryv1.EndpointSlice{
					{
						Endpoints: []discoveryv1.Endpoint{
							{
								Addresses: []string{"1.1.1.1"},
								NodeName:  lo.ToPtr("dummy"),
							},
						},
						Ports: builder.NewEndpointPort(80).WithName("another-name").WithProtocol(corev1.ProtocolTCP).IntoSlice(),
					},
				}, nil
			},
			result: []util.Endpoint{},
		},
		{
			name: "should return one endpoint when the name of the port name match a port in the EndpointSlices",
			svc: &corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			port: &corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			proto: corev1.ProtocolTCP,
			fn: func(string, string) ([]*discoveryv1.EndpointSlice, error) {
				return []*discoveryv1.EndpointSlice{
					{
						Endpoints: []discoveryv1.Endpoint{
							{
								Addresses: []string{"1.1.1.1"},
								NodeName:  lo.ToPtr("dummy"),
							},
						},
						Ports: []discoveryv1.EndpointPort{
							{
								Protocol: lo.ToPtr(corev1.ProtocolTCP),
								Port:     lo.ToPtr(int32(80)),
								Name:     lo.ToPtr("default"),
							},
						},
					},
				}, nil
			},
			result: []util.Endpoint{
				{
					Address: "1.1.1.1",
					Port:    "80",
				},
			},
		},
		{
			name: "should return one endpoint when the name of the port name match more than one port in the endpointSlice",
			svc: &corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromString("port-1"),
						},
					},
				},
			},
			port: &corev1.ServicePort{
				Name:       "port-1",
				TargetPort: intstr.FromString("port-1"),
			},
			proto: corev1.ProtocolTCP,
			fn: func(string, string) ([]*discoveryv1.EndpointSlice, error) {
				return []*discoveryv1.EndpointSlice{
					{
						Endpoints: []discoveryv1.Endpoint{
							{
								Addresses: []string{"1.1.1.1"},
								NodeName:  lo.ToPtr("dummy"),
							},
						},

						Ports: []discoveryv1.EndpointPort{
							builder.NewEndpointPort(80).WithName("port-1").WithProtocol(corev1.ProtocolTCP).Build(),
							builder.NewEndpointPort(80).WithName("port-1").WithProtocol(corev1.ProtocolTCP).Build(),
						},
					},
				}, nil
			},
			result: []util.Endpoint{
				{
					Address: "1.1.1.1",
					Port:    "80",
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := getEndpoints(logrus.New(), testCase.svc, testCase.port, testCase.proto, testCase.fn, testCase.isServiceUpstream)
			require.Equal(t, testCase.result, result)
		})
	}
}

func TestKnativeSelectSplit(t *testing.T) {
	type args struct {
		splits []knative.IngressBackendSplit
	}
	tests := []struct {
		name string
		args args
		want knative.IngressBackendSplit
	}{
		{
			name: "empty ingress",
		},
		{
			name: "no split",
			args: args{
				splits: []knative.IngressBackendSplit{
					{
						IngressBackend: knative.IngressBackend{
							ServiceNamespace: "foo-ns",
							ServiceName:      "foo-svc",
							ServicePort:      intstr.FromInt(42),
						},
						Percent: 100,
					},
				},
			},
			want: knative.IngressBackendSplit{
				IngressBackend: knative.IngressBackend{
					ServiceNamespace: "foo-ns",
					ServiceName:      "foo-svc",
					ServicePort:      intstr.FromInt(42),
				},
				Percent: 100,
			},
		},
		{
			name: "less than 100%% but one split only",
			args: args{
				splits: []knative.IngressBackendSplit{
					{
						IngressBackend: knative.IngressBackend{
							ServiceNamespace: "foo-ns",
							ServiceName:      "foo-svc",
							ServicePort:      intstr.FromInt(42),
						},
						Percent: 42,
					},
				},
			},
			want: knative.IngressBackendSplit{
				IngressBackend: knative.IngressBackend{
					ServiceNamespace: "foo-ns",
					ServiceName:      "foo-svc",
					ServicePort:      intstr.FromInt(42),
				},
				Percent: 42,
			},
		},
		{
			name: "multiple splits with unequal splits",
			args: args{
				splits: []knative.IngressBackendSplit{
					{
						IngressBackend: knative.IngressBackend{
							ServiceNamespace: "bar-ns",
							ServiceName:      "bar-svc",
							ServicePort:      intstr.FromInt(42),
						},
						Percent: 42,
					},
					{
						IngressBackend: knative.IngressBackend{
							ServiceNamespace: "foo-ns",
							ServiceName:      "foo-svc",
							ServicePort:      intstr.FromInt(42),
						},
						Percent: 58,
					},
				},
			},
			want: knative.IngressBackendSplit{
				IngressBackend: knative.IngressBackend{
					ServiceNamespace: "foo-ns",
					ServiceName:      "foo-svc",
					ServicePort:      intstr.FromInt(42),
				},
				Percent: 58,
			},
		},
		{
			name: "multiple splits with unequal splits",
			args: args{
				splits: []knative.IngressBackendSplit{
					{
						IngressBackend: knative.IngressBackend{
							ServiceNamespace: "bar-ns",
							ServiceName:      "bar-svc",
							ServicePort:      intstr.FromInt(42),
						},
						Percent: 40,
					},
					{
						IngressBackend: knative.IngressBackend{
							ServiceNamespace: "baz-ns",
							ServiceName:      "baz-svc",
							ServicePort:      intstr.FromInt(42),
						},
						Percent: 20,
					},
					{
						IngressBackend: knative.IngressBackend{
							ServiceNamespace: "foo-ns",
							ServiceName:      "foo-svc",
							ServicePort:      intstr.FromInt(42),
						},
						Percent: 40,
					},
				},
			},
			want: knative.IngressBackendSplit{
				IngressBackend: knative.IngressBackend{
					ServiceNamespace: "bar-ns",
					ServiceName:      "bar-svc",
					ServicePort:      intstr.FromInt(42),
				},
				Percent: 40,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := knativeSelectSplit(tt.args.splits)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestPickPort(t *testing.T) {
	svc0 := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-0",
			Namespace: "foo-namespace",
			Annotations: map[string]string{
				annotations.IngressClassKey: annotations.DefaultIngressClass,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "port1", Port: 111, TargetPort: intstr.FromInt(1111)},
				{Name: "port2", Port: 222, TargetPort: intstr.FromString("port1")},
				{Name: "port3", Port: 333, TargetPort: intstr.FromString("potato")},
				{Port: 444},
			},
		},
	}

	svc1 := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-1",
			Namespace: "foo-namespace",
			Annotations: map[string]string{
				annotations.IngressClassKey: annotations.DefaultIngressClass,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "port1", Port: 9999},
			},
		},
	}

	svc2 := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-2",
			Namespace: "foo-namespace",
			Annotations: map[string]string{
				annotations.IngressClassKey: annotations.DefaultIngressClass,
			},
		},
		Spec: corev1.ServiceSpec{
			Type:         corev1.ServiceTypeExternalName,
			ExternalName: "external.example.com",
		},
	}

	endpointSliceList := []*discoveryv1.EndpointSlice{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-0-1",
				Namespace: "foo-namespace",
				Labels: map[string]string{
					discoveryv1.LabelServiceName: "service-0",
				},
			},
			Endpoints: []discoveryv1.Endpoint{
				{
					Addresses: []string{"1.1.1.1"},
				},
			},
			Ports: []discoveryv1.EndpointPort{
				builder.NewEndpointPort(111).WithName("port1").WithProtocol(corev1.ProtocolTCP).Build(),
				builder.NewEndpointPort(222).WithName("port2").WithProtocol(corev1.ProtocolTCP).Build(),
				builder.NewEndpointPort(333).WithName("port3").WithProtocol(corev1.ProtocolTCP).Build(),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-1-1",
				Namespace: "foo-namespace",
				Labels: map[string]string{
					discoveryv1.LabelServiceName: "service-1",
				},
			},
			Endpoints: []discoveryv1.Endpoint{
				{
					Addresses: []string{"2.2.2.2"},
				},
			},
			Ports: builder.NewEndpointPort(9999).WithName("port1").WithProtocol(corev1.ProtocolTCP).IntoSlice(),
		},
	}

	for _, tt := range []struct {
		name string
		objs store.FakeObjects
		port netv1.ServiceBackendPort

		wantTarget string
	}{
		{
			name: "port by number",
			objs: store.FakeObjects{
				Services:       []*corev1.Service{&svc0},
				EndpointSlices: endpointSliceList,

				IngressesV1: []*netv1.Ingress{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "foo",
							Namespace:   "foo-namespace",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
						Spec: netv1.IngressSpec{
							Rules: []netv1.IngressRule{
								{
									Host: "example.com",
									IngressRuleValue: netv1.IngressRuleValue{
										HTTP: &netv1.HTTPIngressRuleValue{
											Paths: []netv1.HTTPIngressPath{
												{
													Path: "/",
													Backend: netv1.IngressBackend{
														Service: &netv1.IngressServiceBackend{
															Name: "service-0",
															Port: netv1.ServiceBackendPort{Number: 111},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantTarget: "1.1.1.1:111",
		},
		{
			name: "port by number external name",
			objs: store.FakeObjects{
				Services:       []*corev1.Service{&svc2},
				EndpointSlices: endpointSliceList,

				IngressesV1: []*netv1.Ingress{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "foo",
							Namespace:   "foo-namespace",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
						Spec: netv1.IngressSpec{
							Rules: []netv1.IngressRule{
								{
									Host: "example.com",
									IngressRuleValue: netv1.IngressRuleValue{
										HTTP: &netv1.HTTPIngressRuleValue{
											Paths: []netv1.HTTPIngressPath{
												{
													Path: "/externalname",
													Backend: netv1.IngressBackend{
														Service: &netv1.IngressServiceBackend{
															Name: "service-2",
															Port: netv1.ServiceBackendPort{Number: 222},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantTarget: "external.example.com:222",
		},
		{
			name: "port by name",
			objs: store.FakeObjects{
				Services:       []*corev1.Service{&svc0},
				EndpointSlices: endpointSliceList,

				IngressesV1: []*netv1.Ingress{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "foo",
							Namespace:   "foo-namespace",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
						Spec: netv1.IngressSpec{
							Rules: []netv1.IngressRule{
								{
									Host: "example.com",
									IngressRuleValue: netv1.IngressRuleValue{
										HTTP: &netv1.HTTPIngressRuleValue{
											Paths: []netv1.HTTPIngressPath{
												{
													Path: "/",
													Backend: netv1.IngressBackend{
														Service: &netv1.IngressServiceBackend{
															Name: "service-0",
															Port: netv1.ServiceBackendPort{Name: "port3"},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantTarget: "1.1.1.1:333",
		},
		{
			name: "port implicit",
			objs: store.FakeObjects{
				Services:       []*corev1.Service{&svc1},
				EndpointSlices: endpointSliceList,

				IngressesV1: []*netv1.Ingress{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "foo",
							Namespace:   "foo-namespace",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
						Spec: netv1.IngressSpec{
							Rules: []netv1.IngressRule{
								{
									Host: "example.com",
									IngressRuleValue: netv1.IngressRuleValue{
										HTTP: &netv1.HTTPIngressRuleValue{
											Paths: []netv1.HTTPIngressPath{
												{
													Path: "/",
													Backend: netv1.IngressBackend{
														Service: &netv1.IngressServiceBackend{
															Name: "service-1",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantTarget: "2.2.2.2:9999",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			store, err := store.NewFakeStore(tt.objs)
			require.NoError(t, err)

			p := mustNewParser(t, store)
			result := p.BuildKongConfig()
			require.Empty(t, result.TranslationFailures)

			require.Equal(t, tt.wantTarget, *result.KongState.Upstreams[0].Targets[0].Target.Target)
		})
	}
}

func TestCertificate(t *testing.T) {
	assert := assert.New(t)

	crt1, key1 := certificate.MustGenerateSelfSignedCertPEMFormat()
	crt2, key2 := certificate.MustGenerateSelfSignedCertPEMFormat()
	crt3, key3 := certificate.MustGenerateSelfSignedCertPEMFormat()
	t.Run("same host with multiple namespace return the first namespace/secret by asc", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "ns3",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "ns2",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "ns1",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
		}

		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       k8stypes.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
					Name:      "secret1",
					Namespace: "ns1",
				},
				Data: map[string][]byte{
					"tls.crt": crt1,
					"tls.key": key1,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       k8stypes.UID("6392jz73-180b-4702-a91f-61351a33c6e4"),
					Name:      "secret1",
					Namespace: "ns2",
				},
				Data: map[string][]byte{
					"tls.crt": crt2,
					"tls.key": key2,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       k8stypes.UID("72x2j56k-180b-4702-a91f-61351a33c6e4"),
					Name:      "secret1",
					Namespace: "ns3",
				},
				Data: map[string][]byte{
					"tls.crt": crt3,
					"tls.key": key3,
				},
			},
		}
		fooCertificate := kongstate.Certificate{
			Certificate: kong.Certificate{
				ID:   kong.String("7428fb98-180b-4702-a91f-61351a33c6e4"),
				Cert: kong.String(strings.TrimSpace(string(crt1))),
				Key:  kong.String(strings.TrimSpace(string(key1))),
				SNIs: []*string{kong.String("foo.com")},
				Tags: []*string{
					kong.String("k8s-name:secret1"),
					kong.String("k8s-namespace:ns1"),
					kong.String("k8s-uid:7428fb98-180b-4702-a91f-61351a33c6e4"),
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Secrets:     secrets,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		assert.Len(state.Certificates, 3)
		// foo.com with cert should be fixed
		assert.Contains(state.Certificates, fooCertificate)
	})
	t.Run("SNIs slice with same certificate should be ordered by asc", func(t *testing.T) {
		ingresses := []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo3",
					Namespace: "ns1",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret",
							Hosts:      []string{"foo3.xxx.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo2",
					Namespace: "ns1",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret",
							Hosts:      []string{"foo2.xxx.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo1",
					Namespace: "ns1",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					TLS: []netv1.IngressTLS{
						{
							SecretName: "secret",
							Hosts:      []string{"foo1.xxx.com"},
						},
					},
				},
			},
		}

		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       k8stypes.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
					Name:      "secret",
					Namespace: "ns1",
				},
				Data: map[string][]byte{
					"tls.crt": crt1,
					"tls.key": key1,
				},
			},
		}
		fooCertificate := kongstate.Certificate{
			Certificate: kong.Certificate{
				ID:   kong.String("7428fb98-180b-4702-a91f-61351a33c6e4"),
				Cert: kong.String(strings.TrimSpace(string(crt1))),
				Key:  kong.String(strings.TrimSpace(string(key1))),
				SNIs: []*string{
					kong.String("foo1.xxx.com"),
					kong.String("foo2.xxx.com"),
					kong.String("foo3.xxx.com"),
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: ingresses,
			Secrets:     secrets,
		})
		require.NoError(t, err)
		p := mustNewParser(t, store)
		result := p.BuildKongConfig()
		require.Empty(t, result.TranslationFailures)
		state := result.KongState
		require.NotNil(t, state)
		assert.Len(state.Certificates, 1)
		// parser tests do not check tags, these are tested independently
		state.Certificates[0].Tags = nil
		assert.Equal(state.Certificates[0], fooCertificate)
	})
}

func TestParser_FillsEntitiesIDs(t *testing.T) {
	s, err := store.NewFakeStore(store.FakeObjects{
		Services: []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "svc.foo",
					Namespace: "ns",
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name: "foo",
							Port: 80,
						},
					},
				},
			},
		},
		IngressesV1: []*netv1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ingress.foo",
					Namespace: "ns",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "foo.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/foo",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "svc.foo",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		KongConsumers: []*kongv1.KongConsumer{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user.foo",
					Namespace: "ns",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Username: "user.foo",
			},
		},
	})
	require.NoError(t, err)
	p := mustNewParser(t, s)

	result := p.BuildKongConfig()
	require.Empty(t, result.TranslationFailures)
	state := result.KongState
	require.NotNil(t, state)

	require.Len(t, state.Services, 1)
	require.NotNil(t, state.Services[0].ID)
	assert.Equal(t, "10157742-edd5-51f6-8141-f21dc8017e87", *state.Services[0].ID, "expected deterministic ID")

	require.Len(t, state.Services[0].Routes, 1)
	require.NotNil(t, state.Services[0].Routes[0].ID)
	assert.Equal(t, "1474bf56-c263-5919-b3e6-e9bc6e4b237f", *state.Services[0].Routes[0].ID, "expected deterministic ID")

	require.Len(t, state.Consumers, 1)
	require.NotNil(t, state.Consumers[0].ID)
	assert.Equal(t, "93c4b796-7cc1-5f86-834c-3bbdf00a806c", *state.Consumers[0].ID, "expected deterministic ID")
}

func TestNewFeatureFlags(t *testing.T) {
	testCases := []struct {
		name string

		featureGates     map[string]bool
		kongVersion      semver.Version
		routerFlavor     string
		updateStatusFlag bool

		expectedFeatureFlags FeatureFlags
		expectInfoLog        string
	}{
		{
			name:             "default",
			featureGates:     map[string]bool{},
			routerFlavor:     "traditional",
			updateStatusFlag: true,
			expectedFeatureFlags: FeatureFlags{
				ReportConfiguredKubernetesObjects: true,
			},
		},
		{
			name: "expression routes feature gate enabled and router flavor matches",
			featureGates: map[string]bool{
				featuregates.CombinedRoutesFeature:   true,
				featuregates.ExpressionRoutesFeature: true,
			},
			routerFlavor: kongRouterFlavorExpressions,
			expectedFeatureFlags: FeatureFlags{
				CombinedServiceRoutes: true,
				ExpressionRoutes:      true,
			},
			expectInfoLog: "expression routes mode enabled",
		},
		{
			name: "expression routes feature gate enabled and router flavor does not match",
			featureGates: map[string]bool{
				featuregates.CombinedRoutesFeature:   true,
				featuregates.ExpressionRoutesFeature: true,
			},
			routerFlavor: "any_other_router_mode",
			expectedFeatureFlags: FeatureFlags{
				CombinedServiceRoutes: true,
			},
			expectInfoLog: "ExpressionRoutes feature gate enabled but Gateway is running with \"any_other_router_mode\" router flavor, using that instead",
		},
		{
			name: "expression routes feature gate enabled and combined routes disabled",
			featureGates: map[string]bool{
				featuregates.ExpressionRoutesFeature: true,
			},
			routerFlavor:         kongRouterFlavorExpressions,
			expectedFeatureFlags: FeatureFlags{},
			expectInfoLog:        "ExpressionRoutes feature gate is enabled but CombinedRoutes feature gate is disabled, do not enable expression routes",
		},
		{
			name: "combined routes feature gate enabled",
			featureGates: map[string]bool{
				featuregates.CombinedRoutesFeature: true,
			},
			expectedFeatureFlags: FeatureFlags{
				CombinedServiceRoutes: true,
			},
			expectInfoLog: "combined routes mode has been enabled",
		},
		{
			name:        "kong version >= 3.0 enables regex path prefix",
			kongVersion: semver.Version{Major: 3, Minor: 0},
			expectedFeatureFlags: FeatureFlags{
				RegexPathPrefix: true,
			},
		},
		{
			name: "combined services and routes enabled",
			featureGates: map[string]bool{
				featuregates.CombinedRoutesFeature:   true,
				featuregates.CombinedServicesFeature: true,
			},
			expectedFeatureFlags: FeatureFlags{
				CombinedServiceRoutes: true,
				CombinedServices:      true,
			},
		},
		{
			name: "combined services enabled, but routes disabled",
			featureGates: map[string]bool{
				featuregates.CombinedRoutesFeature:   false,
				featuregates.CombinedServicesFeature: true,
			},
			expectedFeatureFlags: FeatureFlags{
				CombinedServiceRoutes: false,
				CombinedServices:      false,
			},
		},
		{
			name: "fill ids enabled",
			featureGates: map[string]bool{
				featuregates.FillIDsFeature: true,
			},
			expectedFeatureFlags: FeatureFlags{
				FillIDs: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log, logHook := test.NewNullLogger()
			actualFlags := NewFeatureFlags(log, tc.featureGates, tc.kongVersion, tc.routerFlavor, tc.updateStatusFlag)

			require.Equal(t, tc.expectedFeatureFlags, actualFlags)

			if tc.expectInfoLog != "" {
				lastEntry := logHook.LastEntry()
				require.NotNil(t, lastEntry)
				require.Equal(t, logrus.InfoLevel, lastEntry.Level)
				require.Equal(t, tc.expectInfoLog, lastEntry.Message)
			}
		})
	}
}

type mockLicenseGetter struct {
	license mo.Option[kong.License]
}

func (m *mockLicenseGetter) GetLicense() mo.Option[kong.License] {
	return m.license
}

func TestParser_License(t *testing.T) {
	s, _ := store.NewFakeStore(store.FakeObjects{})
	p := mustNewParser(t, s)

	t.Run("no license is populated by default", func(t *testing.T) {
		result := p.BuildKongConfig()
		require.Empty(t, result.KongState.Licenses)
	})

	t.Run("no license is populated when license getter returns no license", func(t *testing.T) {
		p.InjectLicenseGetter(&mockLicenseGetter{})
		result := p.BuildKongConfig()
		require.Empty(t, result.KongState.Licenses)
	})

	t.Run("license is populated when license getter returns a license", func(t *testing.T) {
		licenseGetterWithLicense := &mockLicenseGetter{
			license: mo.Some(kong.License{
				ID:      lo.ToPtr("license-id"),
				Payload: lo.ToPtr("license-payload"),
			}),
		}
		p.InjectLicenseGetter(licenseGetterWithLicense)
		result := p.BuildKongConfig()
		require.Len(t, result.KongState.Licenses, 1)
		license := result.KongState.Licenses[0]
		require.Equal(t, "license-id", *license.ID)
		require.Equal(t, "license-payload", *license.Payload)
	})
}

func TestParser_ConfiguredKubernetesObjects(t *testing.T) {
	testCases := []struct {
		name                          string
		objectsInStore                store.FakeObjects
		expectedObjectsToBeConfigured []k8stypes.NamespacedName
	}{
		{
			name:                          "no objects in cache",
			objectsInStore:                store.FakeObjects{},
			expectedObjectsToBeConfigured: []k8stypes.NamespacedName{},
		},
		{
			name: "KongConsumers",
			objectsInStore: store.FakeObjects{
				KongConsumers: []*kongv1.KongConsumer{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "consumer1",
							Namespace:   "bar",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
						Username: "consumer1",
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "consumer2",
							Namespace:   "bar",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
						Username: "consumer2",
					},
				},
			},
			expectedObjectsToBeConfigured: []k8stypes.NamespacedName{
				{Name: "consumer1", Namespace: "bar"},
				{Name: "consumer2", Namespace: "bar"},
			},
		},
		{
			name: "KongConsumerGroup",
			objectsInStore: store.FakeObjects{
				KongConsumerGroups: []*kongv1beta1.KongConsumerGroup{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "consumer-group1",
							Namespace:   "bar",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "consumer-group2",
							Namespace:   "bar",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
					},
				},
			},
			expectedObjectsToBeConfigured: []k8stypes.NamespacedName{
				{Name: "consumer-group1", Namespace: "bar"},
				{Name: "consumer-group2", Namespace: "bar"},
			},
		},
		{
			name: "KongPlugins with KongConsumer",
			objectsInStore: store.FakeObjects{
				KongPlugins: []*kongv1.KongPlugin{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "plugin1",
							Namespace:   "bar",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
						PluginName: "plugin1",
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "plugin2",
							Namespace:   "bar",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
						PluginName: "plugin2",
					},
				},
				KongConsumers: []*kongv1.KongConsumer{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "consumer",
							Namespace: "bar",
							Annotations: map[string]string{
								annotations.IngressClassKey:                           annotations.DefaultIngressClass,
								annotations.AnnotationPrefix + annotations.PluginsKey: "plugin1,plugin2",
							},
						},
						Username: "foo",
					},
				},
			},
			expectedObjectsToBeConfigured: []k8stypes.NamespacedName{
				{Name: "plugin1", Namespace: "bar"},
				{Name: "plugin2", Namespace: "bar"},
				{Name: "consumer", Namespace: "bar"},
			},
		},
		{
			name: "KongClusterPlugins with KongConsumer",
			objectsInStore: store.FakeObjects{
				KongClusterPlugins: []*kongv1.KongClusterPlugin{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "plugin1",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
						PluginName: "plugin2",
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "plugin2",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
						PluginName: "plugin2",
					},
				},
				KongConsumers: []*kongv1.KongConsumer{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "consumer",
							Namespace: "bar",
							Annotations: map[string]string{
								annotations.IngressClassKey:                           annotations.DefaultIngressClass,
								annotations.AnnotationPrefix + annotations.PluginsKey: "plugin1,plugin2",
							},
						},
						Username: "foo",
					},
				},
			},
			expectedObjectsToBeConfigured: []k8stypes.NamespacedName{
				{Name: "plugin1"},
				{Name: "plugin2"},
				{Name: "consumer", Namespace: "bar"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, _ := store.NewFakeStore(tc.objectsInStore)
			p := mustNewParser(t, s)

			result := p.BuildKongConfig()
			require.Len(t, result.ConfiguredKubernetesObjects, len(tc.expectedObjectsToBeConfigured))

			for _, expectedObj := range tc.expectedObjectsToBeConfigured {
				assert.True(t, lo.ContainsBy(result.ConfiguredKubernetesObjects, func(obj client.Object) bool {
					return expectedObj.Name == obj.GetName() && expectedObj.Namespace == obj.GetNamespace()
				}), "configured objects do not contain the expected %s, actual: %v", expectedObj, result.ConfiguredKubernetesObjects)
			}
		})
	}
}

func mustNewParser(t *testing.T, storer store.Storer) *Parser {
	const testKongVersion = "3.2.0"

	v, err := semver.Parse(testKongVersion)
	require.NoError(t, err)

	p, err := NewParser(logrus.New(), storer,
		FeatureFlags{
			// We'll assume these are true for all tests.
			FillIDs:                           true,
			ReportConfiguredKubernetesObjects: true,
		},
		v,
	)
	require.NoError(t, err)
	return p
}

func TestTargetsForEndpoints(t *testing.T) {
	// targetsForEndpoints should generate expected output for each type of input Endpoint: hostname, IPv4, and IPv6.
	// Addresses are joined to the Port with a : character, and IPv6 Addresses are additionally surrounded in brackets
	// before joining.
	input := []util.Endpoint{
		{
			Address: "hostname.example",
			Port:    "1111",
		},
		{
			Address: "127.0.0.1",
			Port:    "2222",
		},
		{
			Address: "fe80::cae2:65ff:fe7b:2852",
			Port:    "3333",
		},
	}

	wantTargets := []kongstate.Target{
		{
			Target: kong.Target{
				Target: kong.String("hostname.example:1111"),
			},
		},
		{
			Target: kong.Target{
				Target: kong.String("127.0.0.1:2222"),
			},
		},
		{
			Target: kong.Target{
				Target: kong.String("[fe80::cae2:65ff:fe7b:2852]:3333"),
			},
		},
	}

	targets := targetsForEndpoints(input)

	require.Equal(t, wantTargets, targets)
}
