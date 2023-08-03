package kongstate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

var kongConsumerTypeMeta = metav1.TypeMeta{
	APIVersion: kongv1.GroupVersion.String(),
	Kind:       "KongConsumer",
}

func TestKongState_SanitizedCopy(t *testing.T) {
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
				Plugins:        []Plugin{{Plugin: kong.Plugin{ID: kong.String("1")}}},
				Consumers: []Consumer{{
					KeyAuths: []*KeyAuth{{kong.KeyAuth{ID: kong.String("1"), Key: kong.String("secret")}}},
				}},
			},
			want: KongState{
				Services:       []Service{{Service: kong.Service{ID: kong.String("1")}}},
				Upstreams:      []Upstream{{Upstream: kong.Upstream{ID: kong.String("1")}}},
				Certificates:   []Certificate{{Certificate: kong.Certificate{ID: kong.String("1"), Key: redactedString}}},
				CACertificates: []kong.CACertificate{{ID: kong.String("1")}},
				Plugins:        []Plugin{{Plugin: kong.Plugin{ID: kong.String("1")}}},
				Consumers: []Consumer{{
					KeyAuths: []*KeyAuth{{kong.KeyAuth{ID: kong.String("1"), Key: redactedString}}},
				}},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := *tt.in.SanitizedCopy()
			assert.Equal(t, tt.want, got)
		})
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
			name: "multiple consumers, routes and services",
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
				"ns1:foo":    {Consumer: []string{"foo-consumer"}, Service: []string{"foo-service"}},
				"ns1:bar":    {Consumer: []string{"foo-consumer"}, Service: []string{"foo-service"}},
				"ns1:foobar": {Consumer: []string{"bar-consumer"}},
				"ns2:foo":    {Consumer: []string{"foo-consumer"}, Route: []string{"foo-route"}},
				"ns2:bar":    {Consumer: []string{"foo-consumer"}, Route: []string{"foo-route", "bar-route"}},
				"ns2:baz":    {Route: []string{"bar-route"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.state.getPluginRelations(); !reflect.DeepEqual(got, tt.want) {
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
			},
			Data: map[string][]byte{
				"kongCredType": []byte("key-auth"),
				"key":          []byte("whatever"),
				"ttl":          []byte("1024"),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "barCredSecret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"kongCredType":  []byte("oauth2"),
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
			},
			Data: map[string][]byte{
				"kongCredType": []byte("key-auth"),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "unsupportedCredSecret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"kongCredType": []byte("unsupported"),
				"foo":          []byte("bar"),
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
				{Namespace: "default", Name: "foo"}: "failed to fetch secret",
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
				{Namespace: "default", Name: "foo"}: fmt.Sprintf("failed to provision credential: unsupported kongCredType: %q", "unsupported"),
			},
		},
	}

	for i, tc := range testCases {
		indexStr := strconv.Itoa(i)
		tc := tc
		t.Run(indexStr+"-"+tc.name, func(t *testing.T) {
			store, _ := store.NewFakeStore(store.FakeObjects{
				Secrets:       secrets,
				KongConsumers: tc.k8sConsumers,
			})
			logger := logrus.New()
			failureCollector, err := failures.NewResourceFailuresCollector(logger)
			require.NoError(t, err)

			state := KongState{}
			state.FillConsumersAndCredentials(store, failureCollector, semver.MustParse("2.3.2"))
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
				translationFailures := failureCollector.PopResourceFailures()
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.state.FillIDs(logrus.New())
			tc.expect(t, tc.state)
		})
	}
}
