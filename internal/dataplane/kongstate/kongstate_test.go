package kongstate

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/go-logr/zapr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

var kongConsumerTypeMeta = metav1.TypeMeta{
	APIVersion: kongv1.GroupVersion.String(),
	Kind:       "KongConsumer",
}

var serviceTypeMeta = metav1.TypeMeta{
	APIVersion: "v1",
	Kind:       "Service",
}

func TestKongState_SanitizedCopy(t *testing.T) {
	testedFields := sets.New[string]()
	for _, tt := range []struct {
		name string
		in   KongState
		want KongState
	}{
		{
			name: "sanitizes all consumers and certificates and copies all other fields",
			in: KongState{
				Services:       []Service{{Service: kong.Service{ID: kong.String("1")}}},
				Upstreams:      []Upstream{{Upstream: kong.Upstream{ID: kong.String("1")}}},
				Certificates:   []Certificate{{Certificate: kong.Certificate{ID: kong.String("1"), Key: kong.String("secret")}}},
				CACertificates: []kong.CACertificate{{ID: kong.String("1")}},
				Plugins:        []Plugin{{Plugin: kong.Plugin{ID: kong.String("1"), Config: map[string]interface{}{"key": "secret"}}}},
				Consumers: []Consumer{{
					KeyAuths: []*KeyAuth{{kong.KeyAuth{ID: kong.String("1"), Key: kong.String("secret")}}},
				}},
				Licenses: []License{{kong.License{ID: kong.String("1"), Payload: kong.String("secret")}}},
				ConsumerGroups: []ConsumerGroup{{
					ConsumerGroup: kong.ConsumerGroup{ID: kong.String("1"), Name: kong.String("consumer-group")},
				}},
				Vaults: []Vault{
					{
						Vault: kong.Vault{
							Name: kong.String("test-vault"), Prefix: kong.String("test-vault"),
						},
					},
				},
				CustomEntities: map[string]*KongCustomEntityCollection{
					"test_entities": {
						Schema: EntitySchema{
							Fields: map[string]EntityField{
								"name": {
									Type:     EntityFieldTypeString,
									Required: true,
								},
							},
						},
						Entities: []CustomEntity{
							{
								Object: map[string]interface{}{
									"name": "foo",
								},
							},
						},
					},
				},
			},
			want: KongState{
				Services:       []Service{{Service: kong.Service{ID: kong.String("1")}}},
				Upstreams:      []Upstream{{Upstream: kong.Upstream{ID: kong.String("1")}}},
				Certificates:   []Certificate{{Certificate: kong.Certificate{ID: kong.String("1"), Key: redactedString}}},
				CACertificates: []kong.CACertificate{{ID: kong.String("1")}},
				Plugins:        []Plugin{{Plugin: kong.Plugin{ID: kong.String("1"), Config: map[string]interface{}{"key": "secret"}}}}, // We don't redact plugins' config.
				Consumers: []Consumer{{
					KeyAuths: []*KeyAuth{{kong.KeyAuth{ID: kong.String("1"), Key: kong.String("{vault://52fdfc07-2182-454f-963f-5f0f9a621d72}")}}},
				}},
				Licenses: []License{{kong.License{ID: kong.String("1"), Payload: redactedString}}},
				ConsumerGroups: []ConsumerGroup{{
					ConsumerGroup: kong.ConsumerGroup{ID: kong.String("1"), Name: kong.String("consumer-group")},
				}},
				Vaults: []Vault{
					{
						Vault: kong.Vault{
							Name: kong.String("test-vault"), Prefix: kong.String("test-vault"),
						},
					},
				},
				CustomEntities: map[string]*KongCustomEntityCollection{
					"test_entities": {
						Schema: EntitySchema{
							Fields: map[string]EntityField{
								"name": {
									Type:     EntityFieldTypeString,
									Required: true,
								},
							},
						},
						Entities: []CustomEntity{
							{
								Object: map[string]interface{}{
									"name": "foo",
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			testedFields.Insert(extractNotEmptyFieldNames(tt.in)...)
			got := *tt.in.SanitizedCopy(mocks.StaticUUIDGenerator{UUID: "52fdfc07-2182-454f-963f-5f0f9a621d72"})
			assert.Equal(t, tt.want, got)
		})
	}

	ensureAllKongStateFieldsAreCoveredInTest(t, testedFields.UnsortedList())
}

func BenchmarkSanitizedCopy(b *testing.B) {
	const count = 1000
	ks := KongState{
		Certificates: func() []Certificate {
			certificates := make([]Certificate, 0, count)
			for i := 0; i < count; i++ {
				certificates = append(certificates,
					Certificate{kong.Certificate{ID: kong.String(strconv.Itoa(i)), Key: kong.String("secret")}},
				)
			}
			return certificates
		}(),
		Consumers: func() []Consumer {
			consumers := make([]Consumer, 0, count)
			for i := 0; i < count; i++ {
				consumers = append(consumers,
					Consumer{
						Consumer: kong.Consumer{ID: kong.String(strconv.Itoa(i))},
					},
				)
			}
			return consumers
		}(),
		Licenses: func() []License {
			licenses := make([]License, 0, count)
			for i := 0; i < count; i++ {
				licenses = append(licenses,
					License{kong.License{ID: kong.String(strconv.Itoa(i)), Payload: kong.String("secret")}},
				)
			}
			return licenses
		}(),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ret := ks.SanitizedCopy(mocks.StaticUUIDGenerator{UUID: "52fdfc07-2182-454f-963f-5f0f9a621d72"})
		_ = ret
	}
}

// extractNotEmptyFieldNames returns the names of all non-empty fields in the given KongState.
// This is to programmatically find out what fields are used in a test case.
func extractNotEmptyFieldNames(s KongState) []string {
	var fields []string
	typ := reflect.ValueOf(s).Type()
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		v := reflect.ValueOf(s).Field(i)
		if !f.Anonymous && f.IsExported() && !v.IsZero() {
			fields = append(fields, f.Name)
		}
	}
	return fields
}

// ensureAllKongStateFieldsAreCoveredInTest ensures that all fields in KongState are covered in a tests.
func ensureAllKongStateFieldsAreCoveredInTest(t *testing.T, testedFields []string) {
	allKongStateFields := func() []string {
		var fields []string
		typ := reflect.ValueOf(KongState{}).Type()
		for i := 0; i < typ.NumField(); i++ {
			fields = append(fields, typ.Field(i).Name)
		}
		return fields
	}()

	// Meta test - ensure we have testcases covering all fields in KongState.
	for _, field := range allKongStateFields {
		require.Containsf(t, testedFields, field, "field %s wasn't tested", field)
	}
}

func TestGetPluginRelations(t *testing.T) {
	type args struct {
		state KongState
	}
	tests := []struct {
		name string
		args args
		want map[string]util.ForeignRelations
	}{
		{
			name: "empty state",
			want: map[string]util.ForeignRelations{},
		},
		{
			name: "single consumer annotation",
			args: args{
				state: KongState{
					Consumers: []Consumer{
						{
							Consumer: kong.Consumer{
								Username: kong.String("foo-consumer"),
							},
							K8sKongConsumer: kongv1.KongConsumer{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns1",
									Annotations: map[string]string{
										annotations.AnnotationPrefix + annotations.PluginsKey: "foo,bar",
									},
								},
							},
						},
					},
				},
			},
			want: map[string]util.ForeignRelations{
				"ns1:foo": {Consumer: []string{"foo-consumer"}},
				"ns1:bar": {Consumer: []string{"foo-consumer"}},
			},
		},
		{
			name: "single consumer group annotation",
			args: args{
				state: KongState{
					ConsumerGroups: []ConsumerGroup{
						{
							ConsumerGroup: kong.ConsumerGroup{
								Name: kong.String("foo-consumer-group"),
							},
							K8sKongConsumerGroup: kongv1beta1.KongConsumerGroup{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns1",
									Annotations: map[string]string{
										annotations.AnnotationPrefix + annotations.PluginsKey: "foo,bar",
									},
								},
							},
						},
					},
				},
			},
			want: map[string]util.ForeignRelations{
				"ns1:foo": {ConsumerGroup: []string{"foo-consumer-group"}},
				"ns1:bar": {ConsumerGroup: []string{"foo-consumer-group"}},
			},
		},
		{
			name: "single service annotation",
			args: args{
				state: KongState{
					Services: []Service{
						{
							Service: kong.Service{
								Name: kong.String("foo-service"),
							},
							K8sServices: map[string]*corev1.Service{
								"foo-service": {
									ObjectMeta: metav1.ObjectMeta{
										Namespace: "ns1",
										Annotations: map[string]string{
											annotations.AnnotationPrefix + annotations.PluginsKey: "foo,bar",
										},
									},
								},
							},
						},
					},
				},
			},
			want: map[string]util.ForeignRelations{
				"ns1:foo": {Service: []string{"foo-service"}},
				"ns1:bar": {Service: []string{"foo-service"}},
			},
		},
		{
			name: "single Ingress annotation",
			args: args{
				state: KongState{
					Services: []Service{
						{
							Service: kong.Service{
								Name: kong.String("foo-service"),
							},
							Routes: []Route{
								{
									Route: kong.Route{
										Name: kong.String("foo-route"),
									},
									Ingress: util.K8sObjectInfo{
										Name:      "some-ingress",
										Namespace: "ns2",
										Annotations: map[string]string{
											annotations.AnnotationPrefix + annotations.PluginsKey: "foo,bar",
										},
									},
								},
							},
						},
					},
				},
			},
			want: map[string]util.ForeignRelations{
				"ns2:foo": {Route: []string{"foo-route"}},
				"ns2:bar": {Route: []string{"foo-route"}},
			},
		},
		{
			name: "multiple routes with annotation",
			args: args{
				state: KongState{
					Services: []Service{
						{
							Service: kong.Service{
								Name: kong.String("foo-service"),
							},
							Routes: []Route{
								{
									Route: kong.Route{
										Name: kong.String("foo-route"),
									},
									Ingress: util.K8sObjectInfo{
										Name:      "some-ingress",
										Namespace: "ns2",
										Annotations: map[string]string{
											annotations.AnnotationPrefix + annotations.PluginsKey: "foo,bar",
										},
									},
								},
								{
									Route: kong.Route{
										Name: kong.String("bar-route"),
									},
									Ingress: util.K8sObjectInfo{
										Name:      "some-ingress",
										Namespace: "ns2",
										Annotations: map[string]string{
											annotations.AnnotationPrefix + annotations.PluginsKey: "bar,baz",
										},
									},
								},
							},
						},
					},
				},
			},
			want: map[string]util.ForeignRelations{
				"ns2:foo": {Route: []string{"foo-route"}},
				"ns2:bar": {Route: []string{"foo-route", "bar-route"}},
				"ns2:baz": {Route: []string{"bar-route"}},
			},
		},
		{
			name: "multiple consumers, consumer groups, routes and services",
			args: args{
				state: KongState{
					Consumers: []Consumer{
						{
							Consumer: kong.Consumer{
								Username: kong.String("foo-consumer"),
							},
							K8sKongConsumer: kongv1.KongConsumer{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns1",
									Annotations: map[string]string{
										annotations.AnnotationPrefix + annotations.PluginsKey: "foo,bar",
									},
								},
							},
						},
						{
							Consumer: kong.Consumer{
								Username: kong.String("foo-consumer"),
							},
							K8sKongConsumer: kongv1.KongConsumer{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns2",
									Annotations: map[string]string{
										annotations.AnnotationPrefix + annotations.PluginsKey: "foo,bar",
									},
								},
							},
						},
						{
							Consumer: kong.Consumer{
								Username: kong.String("bar-consumer"),
							},
							K8sKongConsumer: kongv1.KongConsumer{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns1",
									Annotations: map[string]string{
										annotations.AnnotationPrefix + annotations.PluginsKey: "foobar",
									},
								},
							},
						},
					},
					ConsumerGroups: []ConsumerGroup{
						{
							ConsumerGroup: kong.ConsumerGroup{
								Name: kong.String("foo-consumer-group"),
							},
							K8sKongConsumerGroup: kongv1beta1.KongConsumerGroup{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns1",
									Annotations: map[string]string{
										annotations.AnnotationPrefix + annotations.PluginsKey: "foo,bar",
									},
								},
							},
						},
						{
							ConsumerGroup: kong.ConsumerGroup{
								Name: kong.String("foo-consumer-group"),
							},
							K8sKongConsumerGroup: kongv1beta1.KongConsumerGroup{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns2",
									Annotations: map[string]string{
										annotations.AnnotationPrefix + annotations.PluginsKey: "foo,bar",
									},
								},
							},
						},
						{
							ConsumerGroup: kong.ConsumerGroup{
								Name: kong.String("bar-consumer-group"),
							},
							K8sKongConsumerGroup: kongv1beta1.KongConsumerGroup{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns2",
									Annotations: map[string]string{
										annotations.AnnotationPrefix + annotations.PluginsKey: "bar,baz",
									},
								},
							},
						},
					},
					Services: []Service{
						{
							Service: kong.Service{
								Name: kong.String("foo-service"),
							},
							K8sServices: map[string]*corev1.Service{
								"foo-service": {
									ObjectMeta: metav1.ObjectMeta{
										Namespace: "ns1",
										Annotations: map[string]string{
											annotations.AnnotationPrefix + annotations.PluginsKey: "foo,bar",
										},
									},
								},
							},
							Routes: []Route{
								{
									Route: kong.Route{
										Name: kong.String("foo-route"),
									},
									Ingress: util.K8sObjectInfo{
										Name:      "some-ingress",
										Namespace: "ns2",
										Annotations: map[string]string{
											annotations.AnnotationPrefix + annotations.PluginsKey: "foo,bar",
										},
									},
								},
								{
									Route: kong.Route{
										Name: kong.String("bar-route"),
									},
									Ingress: util.K8sObjectInfo{
										Name:      "some-ingress",
										Namespace: "ns2",
										Annotations: map[string]string{
											annotations.AnnotationPrefix + annotations.PluginsKey: "bar,baz",
										},
									},
								},
							},
						},
					},
				},
			},
			want: map[string]util.ForeignRelations{
				"ns1:foo":    {Consumer: []string{"foo-consumer"}, ConsumerGroup: []string{"foo-consumer-group"}, Service: []string{"foo-service"}},
				"ns1:bar":    {Consumer: []string{"foo-consumer"}, ConsumerGroup: []string{"foo-consumer-group"}, Service: []string{"foo-service"}},
				"ns1:foobar": {Consumer: []string{"bar-consumer"}},
				"ns2:foo": {
					Consumer: []string{"foo-consumer"}, ConsumerGroup: []string{"foo-consumer-group"}, Route: []string{"foo-route"},
				},
				"ns2:bar": {
					Consumer: []string{"foo-consumer"}, ConsumerGroup: []string{"foo-consumer-group", "bar-consumer-group"}, Route: []string{"foo-route", "bar-route"},
				},
				"ns2:baz": {Route: []string{"bar-route"}, ConsumerGroup: []string{"bar-consumer-group"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, _ := store.NewFakeStore(store.FakeObjects{})
			if got := tt.args.state.getPluginRelations(store, logr.Discard()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPluginRelations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFillConsumersAndCredentials(t *testing.T) {
	secrets := []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fooCredSecret",
				Namespace: "default",
				Labels: map[string]string{
					labels.CredentialTypeLabel: "key-auth",
				},
			},
			Data: map[string][]byte{
				"key": []byte("whatever"),
				"ttl": []byte("1024"),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "barCredSecret",
				Namespace: "default",
				Labels: map[string]string{
					labels.CredentialTypeLabel: "oauth2",
				},
			},
			Data: map[string][]byte{
				"name":          []byte("whatever"),
				"client_id":     []byte("whatever"),
				"client_secret": []byte("whatever"),
				"redirect_uris": []byte("http://example.com"),
				"hash_secret":   []byte("true"),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "emptyCredSecret",
				Namespace: "default",
				Labels: map[string]string{
					labels.CredentialTypeLabel: "key-auth",
				},
			},
			Data: map[string][]byte{},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "unsupportedCredSecret",
				Namespace: "default",
				Labels: map[string]string{
					labels.CredentialTypeLabel: "unsupported",
				},
			},
			Data: map[string][]byte{
				"foo": []byte("bar"),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "labeledSecret",
				Namespace: "default",
				Labels: map[string]string{
					labels.CredentialTypeLabel: "key-auth",
				},
			},
			Data: map[string][]byte{
				"key": []byte("little-rabbits-be-good"),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "badTypeLabeledSecret",
				Namespace: "default",
				Labels: map[string]string{
					labels.CredentialTypeLabel: "bee-auth",
				},
			},
			Data: map[string][]byte{
				"foo": []byte("bar"),
			},
		},
	}

	testCases := []struct {
		name                               string
		k8sConsumers                       []*kongv1.KongConsumer
		expectedKongStateConsumers         []Consumer
		expectedTranslationFailureMessages map[k8stypes.NamespacedName]string
	}{
		{
			name: "KongConsumer with key-auth and oauth2",
			k8sConsumers: []*kongv1.KongConsumer{
				{
					TypeMeta: kongConsumerTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "default",
						Annotations: map[string]string{
							"kubernetes.io/ingress.class": annotations.DefaultIngressClass,
						},
					},
					Username: "foo",
					CustomID: "foo",
					Credentials: []string{
						"fooCredSecret",
						"barCredSecret",
					},
				},
			},
			expectedKongStateConsumers: []Consumer{
				{
					Consumer: kong.Consumer{
						Username: kong.String("foo"),
						CustomID: kong.String("foo"),
					},
					KeyAuths: []*KeyAuth{{kong.KeyAuth{
						Key: kong.String("whatever"),
						TTL: kong.Int(1024),
						Tags: util.GenerateTagsForObject(&corev1.Secret{
							ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "fooCredSecret"},
						}),
					}}},
					Oauth2Creds: []*Oauth2Credential{
						{
							kong.Oauth2Credential{
								Name:         kong.String("whatever"),
								ClientID:     kong.String("whatever"),
								ClientSecret: kong.String("whatever"),
								HashSecret:   kong.Bool(true),
								RedirectURIs: []*string{kong.String("http://example.com")},
								Tags: util.GenerateTagsForObject(&corev1.Secret{
									ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "barCredSecret"},
								}),
							},
						},
					},
				},
			},
		},
		{
			name: "missing username and custom_id",
			k8sConsumers: []*kongv1.KongConsumer{
				{
					TypeMeta: kongConsumerTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "default",
						Annotations: map[string]string{
							"kubernetes.io/ingress.class": annotations.DefaultIngressClass,
						},
					},
					Credentials: []string{
						"fooCredSecret",
						"barCredSecret",
					},
				},
			},
			expectedTranslationFailureMessages: map[k8stypes.NamespacedName]string{
				{Namespace: "default", Name: "foo"}: "no username or custom_id specified",
			},
		},
		{
			name: "referring to non-exist secret",
			k8sConsumers: []*kongv1.KongConsumer{
				{
					TypeMeta: kongConsumerTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "default",
						Annotations: map[string]string{
							"kubernetes.io/ingress.class": annotations.DefaultIngressClass,
						},
					},
					Username: "foo",
					Credentials: []string{
						"nonExistCredSecret",
					},
				},
			},
			expectedKongStateConsumers: []Consumer{
				{
					Consumer: kong.Consumer{
						Username: kong.String("foo"),
					},
				},
			},
			expectedTranslationFailureMessages: map[k8stypes.NamespacedName]string{
				{Namespace: "default", Name: "foo"}: "Failed to fetch secret",
			},
		},
		{
			name: "referring to secret with unsupported credType",
			k8sConsumers: []*kongv1.KongConsumer{
				{
					TypeMeta: kongConsumerTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "default",
						Annotations: map[string]string{
							"kubernetes.io/ingress.class": annotations.DefaultIngressClass,
						},
					},
					Username: "foo",
					Credentials: []string{
						"unsupportedCredSecret",
					},
				},
			},
			expectedKongStateConsumers: []Consumer{
				{
					Consumer: kong.Consumer{
						Username: kong.String("foo"),
					},
				},
			},
			expectedTranslationFailureMessages: map[k8stypes.NamespacedName]string{
				{Namespace: "default", Name: "foo"}: fmt.Sprintf("failed to provision credential: unsupported credential type: %q", "unsupported"),
			},
		},
		{
			name: "referring to secret with unsupported credential label",
			k8sConsumers: []*kongv1.KongConsumer{
				{
					TypeMeta: kongConsumerTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "default",
						Annotations: map[string]string{
							"kubernetes.io/ingress.class": annotations.DefaultIngressClass,
						},
					},
					Username: "foo",
					Credentials: []string{
						"badTypeLabeledSecret",
					},
				},
			},
			expectedKongStateConsumers: []Consumer{
				{
					Consumer: kong.Consumer{
						Username: kong.String("foo"),
					},
				},
			},
			expectedTranslationFailureMessages: map[k8stypes.NamespacedName]string{
				{Namespace: "default", Name: "foo"}: fmt.Sprintf("failed to provision credential: unsupported credential type: %q", "bee-auth"),
			},
		},
		{
			name: "KongConsumer with key-auth from label secret",
			k8sConsumers: []*kongv1.KongConsumer{
				{
					TypeMeta: kongConsumerTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "default",
						Annotations: map[string]string{
							"kubernetes.io/ingress.class": annotations.DefaultIngressClass,
						},
					},
					Username: "foo",
					CustomID: "foo",
					Credentials: []string{
						"labeledSecret",
					},
				},
			},
			expectedKongStateConsumers: []Consumer{
				{
					Consumer: kong.Consumer{
						Username: kong.String("foo"),
						CustomID: kong.String("foo"),
					},
					KeyAuths: []*KeyAuth{{kong.KeyAuth{
						Key: kong.String("little-rabbits-be-good"),
						Tags: util.GenerateTagsForObject(&corev1.Secret{
							ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "labeledSecret"},
						}),
					}}},
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			store, _ := store.NewFakeStore(store.FakeObjects{
				Secrets:       secrets,
				KongConsumers: tc.k8sConsumers,
			})
			logger := zapr.NewLogger(zap.NewNop())
			failuresCollector := failures.NewResourceFailuresCollector(logger)

			state := KongState{}
			state.FillConsumersAndCredentials(logger, store, failuresCollector)
			// compare translated consumers.
			require.Len(t, state.Consumers, len(tc.expectedKongStateConsumers))
			// compare fields. Since we only test for translating a single consumer, we only compare the first one if exists.
			if len(state.Consumers) > 0 && len(tc.expectedKongStateConsumers) > 0 {
				expectedConsumer := tc.expectedKongStateConsumers[0]
				kongStateConsumer := state.Consumers[0]
				assert.Equal(t, expectedConsumer.Consumer.Username, kongStateConsumer.Consumer.Username, "should have expected username")
				// compare credentials.
				assert.Equal(t, expectedConsumer.KeyAuths, kongStateConsumer.KeyAuths)
				assert.Equal(t, expectedConsumer.Oauth2Creds, kongStateConsumer.Oauth2Creds)
			}
			// check for expected translation failures.
			if len(tc.expectedTranslationFailureMessages) > 0 {
				translationFailures := failuresCollector.PopResourceFailures()
				for nsName, expectedMessage := range tc.expectedTranslationFailureMessages {
					relatedFailures := lo.Filter(translationFailures, func(f failures.ResourceFailure, _ int) bool {
						for _, obj := range f.CausingObjects() {
							if obj.GetNamespace() == nsName.Namespace && obj.GetName() == nsName.Name {
								return true
							}
						}
						return false
					})

					assert.Truef(t, lo.ContainsBy(relatedFailures, func(f failures.ResourceFailure) bool {
						return strings.Contains(f.Message(), expectedMessage)
					}), "should find expected translation failure caused by KongConsumer %s: should contain '%s'",
						nsName.String(), expectedMessage)
				}
			}
		})
	}
}

func TestKongState_FillIDs(t *testing.T) {
	testCases := []struct {
		name   string
		state  KongState
		expect func(t *testing.T, s KongState)
	}{
		{
			name: "fills service IDs",
			state: KongState{
				Services: []Service{
					{
						Service: kong.Service{
							Name: kong.String("service.foo"),
						},
					},
					{
						Service: kong.Service{
							Name: kong.String("service.bar"),
						},
					},
				},
			},
			expect: func(t *testing.T, s KongState) {
				require.NotEmpty(t, s.Services[0].ID)
				require.NotEmpty(t, s.Services[1].ID)
			},
		},
		{
			name: "fills route IDs",
			state: KongState{
				Services: []Service{
					{
						Service: kong.Service{
							Name: kong.String("service.foo"),
						},
						Routes: []Route{
							{
								Route: kong.Route{
									Name: kong.String("route.foo"),
								},
							},
							{
								Route: kong.Route{
									Name: kong.String("route.bar"),
								},
							},
						},
					},
				},
			},
			expect: func(t *testing.T, s KongState) {
				require.NotEmpty(t, s.Services[0].ID)
				require.NotEmpty(t, s.Services[0].Routes[0].ID)
				require.NotEmpty(t, s.Services[0].Routes[1].ID)
			},
		},
		{
			name: "fills consumer IDs",
			state: KongState{
				Consumers: []Consumer{
					{
						Consumer: kong.Consumer{
							Username: kong.String("user.foo"),
						},
					},
					{
						Consumer: kong.Consumer{
							Username: kong.String("user.bar"),
						},
					},
				},
			},
			expect: func(t *testing.T, s KongState) {
				require.NotEmpty(t, s.Consumers[0].ID)
				require.NotEmpty(t, s.Consumers[1].ID)
			},
		},
		{
			name: "fills services, routes, and consumer IDs",
			state: KongState{
				Services: []Service{
					{
						Service: kong.Service{
							Name: kong.String("service.foo"),
						},
						Routes: []Route{
							{
								Route: kong.Route{
									Name: kong.String("route.bar"),
								},
							},
						},
					},
				},
				Consumers: []Consumer{
					{
						Consumer: kong.Consumer{
							Username: kong.String("user.baz"),
						},
					},
				},
			},
			expect: func(t *testing.T, s KongState) {
				require.NotEmpty(t, s.Services[0].ID)
				require.NotEmpty(t, s.Services[0].Routes[0].ID)
				require.NotEmpty(t, s.Consumers[0].ID)
			},
		},
		{
			name: "fills consumer, consumer group, vault IDs",
			state: KongState{
				Consumers: []Consumer{
					{
						Consumer: kong.Consumer{
							Username: kong.String("user.0"),
						},
					},
				},
				ConsumerGroups: []ConsumerGroup{
					{
						ConsumerGroup: kong.ConsumerGroup{
							Name: kong.String("cg.0"),
						},
					},
				},
				Vaults: []Vault{
					{
						Vault: kong.Vault{
							Prefix: kong.String("vault.0"),
						},
					},
				},
			},
			expect: func(t *testing.T, s KongState) {
				require.NotEmpty(t, s.Consumers[0].ID)
				require.NotEmpty(t, s.ConsumerGroups[0].ID)
				require.NotEmpty(t, s.Vaults[0].ID)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.state.FillIDs(zapr.NewLogger(zap.NewNop()), "")
			tc.expect(t, tc.state)
		})
	}
}

func TestKongState_BuildPluginsCollisions(t *testing.T) {
	for _, tt := range []struct {
		name       string
		in         []*kongv1.KongPlugin
		pluginRels map[string]util.ForeignRelations
		want       []string
	}{
		{
			name: "collision test",
			in: []*kongv1.KongPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-plugin",
						Namespace: "default",
					},
					PluginName:   "jwt",
					InstanceName: "test",
				},
			},
			pluginRels: map[string]util.ForeignRelations{
				"default:foo-plugin": {
					// this shouldn't happen in practice, as all generated route names are unique
					// however, it's hard to find a SHA256 collision with two different inputs
					Route: []string{"collision", "collision"},
				},
			},
			want: []string{"test-bae3267aa", "test-bae3267aafead3adb6031bc1c732516336e7f7b324baf61bb68a39cc89112741"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			log := zapr.NewLogger(zap.NewNop())
			store, _ := store.NewFakeStore(store.FakeObjects{
				KongPlugins: tt.in,
			})
			// this is not testing the kongPluginFromK8SPlugin failure cases, so there is no failures collector
			got := buildPlugins(log, store, nil, tt.pluginRels)
			require.Len(t, got, 2)
			require.Equal(t, tt.want, []string{*got[0].InstanceName, *got[1].InstanceName})
		})
	}
}

func TestKongState_FillUpstreamOverrides(t *testing.T) {
	const (
		kongIngressName        = "kongIngress"
		kongUpstreamPolicyName = "policy"
	)
	serviceAnnotatedWithKongUpstreamPolicy := func() *corev1.Service {
		return &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service",
				Namespace: "default",
				Annotations: map[string]string{
					kongv1beta1.KongUpstreamPolicyAnnotationKey: kongUpstreamPolicyName,
				},
			},
		}
	}

	serviceAnnotatedWithKongUpstreamPolicyAndKongIngress := func() *corev1.Service {
		s := serviceAnnotatedWithKongUpstreamPolicy()
		s.Annotations[annotations.AnnotationPrefix+annotations.ConfigurationKey] = kongIngressName
		return s
	}

	testCases := []struct {
		name                 string
		upstream             Upstream
		kongUpstreamPolicies []*kongv1beta1.KongUpstreamPolicy
		kongIngresses        []*kongv1.KongIngress
		expectedUpstream     kong.Upstream
		expectedFailures     []failures.ResourceFailure
	}{
		{
			name: "upstream with no overrides",
			upstream: Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo-upstream"),
				},
			},
			expectedUpstream: kong.Upstream{
				Name: kong.String("foo-upstream"),
			},
		},
		{
			name: "upstream backed by service annotated with KongUpstreamPolicy",
			upstream: Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo-upstream"),
				},
				Service: Service{
					K8sServices: map[string]*corev1.Service{"": serviceAnnotatedWithKongUpstreamPolicy()},
				},
			},
			kongUpstreamPolicies: []*kongv1beta1.KongUpstreamPolicy{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      kongUpstreamPolicyName,
						Namespace: "default",
					},
					Spec: kongv1beta1.KongUpstreamPolicySpec{
						Algorithm: lo.ToPtr("least-connections"),
					},
				},
			},
			expectedUpstream: kong.Upstream{
				Name:      kong.String("foo-upstream"),
				Algorithm: kong.String("least-connections"),
			},
		},
		{
			name: "upstream backed by service annotated with KongUpstreamPolicy that doesn't exist",
			upstream: Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo-upstream"),
				},
				Service: Service{
					K8sServices: map[string]*corev1.Service{"": serviceAnnotatedWithKongUpstreamPolicy()},
				},
			},
			expectedUpstream: kong.Upstream{
				Name: kong.String("foo-upstream"),
			},
			expectedFailures: []failures.ResourceFailure{
				lo.Must(failures.NewResourceFailure(
					"failed fetching KongUpstreamPolicy: KongUpstreamPolicy default/policy not found",
					serviceAnnotatedWithKongUpstreamPolicy(),
				)),
			},
		},
		{
			name: "KongUpstreamPolicy is applied even if KongIngress is not found",
			upstream: Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo-upstream"),
				},
				Service: Service{
					K8sServices: map[string]*corev1.Service{"": serviceAnnotatedWithKongUpstreamPolicyAndKongIngress()},
				},
			},
			kongUpstreamPolicies: []*kongv1beta1.KongUpstreamPolicy{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      kongUpstreamPolicyName,
						Namespace: "default",
					},
					Spec: kongv1beta1.KongUpstreamPolicySpec{
						Algorithm: lo.ToPtr("least-connections"),
					},
				},
			},
			expectedUpstream: kong.Upstream{
				Name:      kong.String("foo-upstream"),
				Algorithm: kong.String("least-connections"),
			},
			expectedFailures: []failures.ResourceFailure{
				lo.Must(failures.NewResourceFailure(
					"failed to get KongIngress: KongIngress kongIngress not found",
					serviceAnnotatedWithKongUpstreamPolicyAndKongIngress(),
				)),
			},
		},
		{
			name: "KongUpstreamPolicy overwrites KongIngress",
			upstream: Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo-upstream"),
				},
				Service: Service{
					K8sServices: map[string]*corev1.Service{"": serviceAnnotatedWithKongUpstreamPolicyAndKongIngress()},
				},
			},
			kongUpstreamPolicies: []*kongv1beta1.KongUpstreamPolicy{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      kongUpstreamPolicyName,
						Namespace: "default",
					},
					Spec: kongv1beta1.KongUpstreamPolicySpec{
						Algorithm: lo.ToPtr("least-connections"),
					},
				},
			},
			kongIngresses: []*kongv1.KongIngress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      kongIngressName,
						Namespace: "default",
					},
					Upstream: &kongv1.KongIngressUpstream{
						Algorithm: lo.ToPtr("round-robin"),
					},
				},
			},
			expectedUpstream: kong.Upstream{
				Name:      kong.String("foo-upstream"),
				Algorithm: kong.String("least-connections"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := store.NewFakeStore(store.FakeObjects{
				KongUpstreamPolicies: tc.kongUpstreamPolicies,
				KongIngresses:        tc.kongIngresses,
			})
			require.NoError(t, err)
			failuresCollector := failures.NewResourceFailuresCollector(logr.Discard())

			kongState := KongState{Upstreams: []Upstream{tc.upstream}}
			kongState.FillUpstreamOverrides(s, logr.Discard(), failuresCollector)
			require.Equal(t, tc.expectedUpstream, kongState.Upstreams[0].Upstream)
			require.ElementsMatch(t, tc.expectedFailures, failuresCollector.PopResourceFailures())
		})
	}
}

func TestFillVaults(t *testing.T) {
	kongVaultTypeMeta := metav1.TypeMeta{
		APIVersion: kongv1alpha1.GroupVersion.String(),
		Kind:       "KongVault",
	}
	now := time.Now()
	testCases := []struct {
		name                     string
		kongVaults               []*kongv1alpha1.KongVault
		expectedTranslatedVaults []Vault
		// name of KongVault -> failure message
		expectedTranslationFailures map[string]string
	}{
		{
			name: "single valid KongVault",
			kongVaults: []*kongv1alpha1.KongVault{
				{
					TypeMeta: kongVaultTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "",
						Name:      "vault-1",
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: kongv1alpha1.KongVaultSpec{
						Backend: "env",
						Prefix:  "env-1",
					},
				},
			},
			expectedTranslatedVaults: []Vault{
				{
					Vault: kong.Vault{
						Name:   kong.String("env"),
						Prefix: kong.String("env-1"),
					},
					K8sKongVault: &kongv1alpha1.KongVault{
						ObjectMeta: metav1.ObjectMeta{
							Name: "vault-1",
						},
					},
				},
			},
		},
		{
			name: "one valid KongVault with correct ingress class, and one KongVault with other ingress class",
			kongVaults: []*kongv1alpha1.KongVault{
				{
					TypeMeta: kongVaultTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name: "vault-1",
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: kongv1alpha1.KongVaultSpec{
						Backend: "env",
						Prefix:  "env-1",
					},
				},
				{
					TypeMeta: kongVaultTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name: "vault-2",
						Annotations: map[string]string{
							annotations.IngressClassKey: "other-ingress-class",
						},
					},
					Spec: kongv1alpha1.KongVaultSpec{
						Backend: "env",
						Prefix:  "env-2",
					},
				},
			},
			expectedTranslatedVaults: []Vault{
				{
					Vault: kong.Vault{
						Name:   kong.String("env"),
						Prefix: kong.String("env-1"),
					},
					K8sKongVault: &kongv1alpha1.KongVault{
						ObjectMeta: metav1.ObjectMeta{
							Name: "vault-1",
						},
					},
				},
			},
		},
		{
			name: "KongVault with invalid configuration is rejected",
			kongVaults: []*kongv1alpha1.KongVault{
				{
					TypeMeta: kongVaultTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name: "vault-invalid",
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: kongv1alpha1.KongVaultSpec{
						Backend: "env",
						Prefix:  "env-1",
						Config: apiextensionsv1.JSON{
							Raw: []byte(`{{}`),
						},
					},
				},
			},
			expectedTranslationFailures: map[string]string{
				"vault-invalid": `failed to parse configuration of vault "vault-invalid" to JSON`,
			},
		},
		{
			name: "multiple KongVaults with same spec.prefix, only one translated and translation failure for the other",
			kongVaults: []*kongv1alpha1.KongVault{
				{
					TypeMeta: kongVaultTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:              "vault-0-newer",
						CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: kongv1alpha1.KongVaultSpec{
						Backend: "env",
						Prefix:  "env-1",
					},
				},
				{
					TypeMeta: kongVaultTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:              "vault-1",
						CreationTimestamp: metav1.NewTime(now.Add(-10 * time.Second)),
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: kongv1alpha1.KongVaultSpec{
						Backend: "env",
						Prefix:  "env-1",
					},
				},
				{
					TypeMeta: kongVaultTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:              "vault-2",
						CreationTimestamp: metav1.NewTime(now.Add(-10 * time.Second)),
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: kongv1alpha1.KongVaultSpec{
						Backend: "env",
						Prefix:  "env-1",
					},
				},
			},
			expectedTranslatedVaults: []Vault{
				{
					Vault: kong.Vault{
						Name:   kong.String("env"),
						Prefix: kong.String("env-1"),
					},
					K8sKongVault: &kongv1alpha1.KongVault{
						ObjectMeta: metav1.ObjectMeta{
							Name: "vault-1",
						},
					},
				},
			},
			expectedTranslationFailures: map[string]string{
				"vault-0-newer": `spec.prefix "env-1" is duplicate`,
				"vault-2":       `spec.prefix "env-1" is duplicate`,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := store.NewFakeStore(store.FakeObjects{
				KongVaults: tc.kongVaults,
			})

			require.NoError(t, err)
			logger := testr.New(t)
			f := failures.NewResourceFailuresCollector(logger)
			ks := &KongState{}
			ks.FillVaults(logger, s, f)

			assert.Len(t, ks.Vaults, len(tc.expectedTranslatedVaults), "should have expected number of translated vaults")
			for _, expectedVault := range tc.expectedTranslatedVaults {
				assert.Truef(t, lo.ContainsBy(ks.Vaults, func(v Vault) bool {
					return (v.Name != nil && *v.Name == *expectedVault.Name) &&
						(v.Prefix != nil && *v.Prefix == *expectedVault.Prefix) &&
						(v.K8sKongVault != nil && v.K8sKongVault.Name == expectedVault.K8sKongVault.Name)
				}),
					"cannot find translated vault for KongVault %q", expectedVault.K8sKongVault.Name,
				)
			}

			translationFailures := f.PopResourceFailures()
			for name, message := range tc.expectedTranslationFailures {
				assert.Truef(t, lo.ContainsBy(translationFailures, func(failure failures.ResourceFailure) bool {
					return strings.Contains(failure.Message(), message) &&
						lo.ContainsBy(failure.CausingObjects(), func(obj client.Object) bool {
							return obj.GetName() == name
						})
				}),
					"cannot find expected translation failure for KongVault %s", name,
				)
			}
		})
	}
}

func TestFillOverrides_ServiceFailures(t *testing.T) {
	tests := []struct {
		name                               string
		state                              *KongState
		want                               Service
		expectedTranslationFailureMessages map[k8stypes.NamespacedName]string
	}{
		{
			name: "service protocol set to valid value",
			state: &KongState{
				Services: []Service{
					{
						Namespace: "default",
						K8sServices: map[string]*corev1.Service{
							"test": {
								TypeMeta: serviceTypeMeta,
								ObjectMeta: metav1.ObjectMeta{
									Name:      "test",
									Namespace: "default",
									Annotations: map[string]string{
										"konghq.com/protocol": "wss",
									},
								},
							},
						},
					},
				},
			},
			want: Service{
				Service: kong.Service{
					Protocol: kong.String("wss"),
				},
			},
		},
		{
			name: "service protocol set to invalid value",
			state: &KongState{
				Services: []Service{
					{
						Namespace: "default",
						K8sServices: map[string]*corev1.Service{
							"test": {
								TypeMeta: serviceTypeMeta,
								ObjectMeta: metav1.ObjectMeta{
									Name:      "test",
									Namespace: "default",
									Annotations: map[string]string{
										"konghq.com/protocol": "djnfkgjfgn",
									},
								},
							},
						},
					},
				},
			},
			want: Service{
				Service: kong.Service{
					Protocol: kong.String("http"),
				},
			},
			expectedTranslationFailureMessages: map[k8stypes.NamespacedName]string{
				{Namespace: "default", Name: "test"}: "konghq.com/protocol annotation has invalid value: djnfkgjfgn",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, err := store.NewFakeStore(store.FakeObjects{})
			require.NoError(t, err)
			logger := zapr.NewLogger(zap.NewNop())
			failuresCollector := failures.NewResourceFailuresCollector(logger)
			tt.state.FillOverrides(logger, store, failuresCollector)
			if len(tt.expectedTranslationFailureMessages) > 0 {
				translationFailures := failuresCollector.PopResourceFailures()
				for nsName, expectedMessage := range tt.expectedTranslationFailureMessages {
					relatedFailures := lo.Filter(translationFailures, func(f failures.ResourceFailure, _ int) bool {
						for _, obj := range f.CausingObjects() {
							if obj.GetNamespace() == nsName.Namespace && obj.GetName() == nsName.Name {
								return true
							}
						}
						return false
					})

					assert.Truef(t, lo.ContainsBy(relatedFailures, func(f failures.ResourceFailure) bool {
						return strings.Contains(f.Message(), expectedMessage)
					}), "should find expected translation failure caused by Service %s: should contain '%s'",
						nsName.String(), expectedMessage)
				}
			}
		})
	}
}

type fakeSchemaGetter struct {
	schemas map[string]kong.Schema
}

var _ SchemaGetter = &fakeSchemaGetter{}

func (s *fakeSchemaGetter) Get(_ context.Context, entityType string) (kong.Schema, error) {
	schema, ok := s.schemas[entityType]
	if !ok {
		return nil, fmt.Errorf("schema not found")
	}
	return schema, nil
}

func TestIsRemotePluginReferenceAllowed(t *testing.T) {
	serviceTypeMeta := metav1.TypeMeta{
		Kind: "Service",
	}

	testCases := []struct {
		name            string
		referrer        client.Object
		pluginNamespace string
		pluginName      string
		referenceGrants []*gatewayapi.ReferenceGrant
		shouldAllow     bool
	}{
		{
			name: "no reference grant",
			referrer: &corev1.Service{
				TypeMeta: serviceTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo",
					Name:      "service-foo",
				},
			},
			pluginNamespace: "bar",
			pluginName:      "plugin-bar",
			shouldAllow:     false,
		},
		{
			name: "have reference grant",
			referrer: &corev1.Service{
				TypeMeta: serviceTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo",
					Name:      "service-foo",
				},
			},
			pluginNamespace: "bar",
			pluginName:      "plugin-bar",
			referenceGrants: []*gatewayapi.ReferenceGrant{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "bar", // same namespace as plugin
						Name:      "grant-1",
					},
					Spec: gatewayapi.ReferenceGrantSpec{
						From: []gatewayapi.ReferenceGrantFrom{
							{
								Kind:      gatewayapi.Kind("Service"),
								Namespace: "foo",
							},
						},
						To: []gatewayapi.ReferenceGrantTo{
							{
								Group: gatewayapi.Group(kongv1.GroupVersion.Group),
								Kind:  gatewayapi.Kind("KongPlugin"),
							},
						},
					},
				},
			},
			shouldAllow: true,
		},
		{
			name: "reference grant created but in different namespace",
			referrer: &corev1.Service{
				TypeMeta: serviceTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo",
					Name:      "service-foo",
				},
			},
			pluginNamespace: "bar",
			pluginName:      "plugin-bar",
			referenceGrants: []*gatewayapi.ReferenceGrant{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "foo", // Not same namespace as plugin
						Name:      "grant-1",
					},
					Spec: gatewayapi.ReferenceGrantSpec{
						From: []gatewayapi.ReferenceGrantFrom{
							{
								Kind:      gatewayapi.Kind("Service"),
								Namespace: "foo",
							},
						},
						To: []gatewayapi.ReferenceGrantTo{
							{
								Group: gatewayapi.Group(kongv1.GroupVersion.Group),
								Kind:  gatewayapi.Kind("KongPlugin"),
							},
						},
					},
				},
			},
			shouldAllow: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := store.NewFakeStore(store.FakeObjects{
				ReferenceGrants: tc.referenceGrants,
			})
			require.NoError(t, err)
			err = isRemotePluginReferenceAllowed(logr.Discard(), s, pluginReference{
				Referrer:  tc.referrer,
				Namespace: tc.pluginNamespace,
				Name:      tc.pluginName,
			})
			if tc.shouldAllow {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
