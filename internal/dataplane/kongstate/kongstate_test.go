package kongstate

import (
	"reflect"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

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
							K8sKongConsumer: configurationv1.KongConsumer{
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
							K8sKongConsumer: configurationv1.KongConsumer{
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
							K8sKongConsumer: configurationv1.KongConsumer{
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
							K8sKongConsumer: configurationv1.KongConsumer{
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
	}
	consumers := []*configurationv1.KongConsumer{
		{
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
	}
	store, _ := store.NewFakeStore(store.FakeObjects{
		Secrets:       secrets,
		KongConsumers: consumers,
	})

	want := KongState{
		Consumers: []Consumer{{
			Consumer: kong.Consumer{
				Username: kong.String("foo"),
				CustomID: kong.String("foo"),
			},
			KeyAuths: []*KeyAuth{{kong.KeyAuth{
				Key: kong.String("whatever"),
				TTL: kong.Int(1024),
			}}},
			Oauth2Creds: []*Oauth2Credential{
				{
					kong.Oauth2Credential{
						Name:         kong.String("whatever"),
						ClientID:     kong.String("whatever"),
						ClientSecret: kong.String("whatever"),
						HashSecret:   kong.Bool(true),
						RedirectURIs: []*string{kong.String("http://example.com")},
					},
				},
			},
		}},
		Version: semver.MustParse("2.3.2"),
	}

	t.Run("parses consumer and credential from store into state", func(t *testing.T) {
		state := KongState{
			Version: semver.MustParse("2.3.2"),
		}
		state.FillConsumersAndCredentials(logrus.New(), store)
		assert.Equal(t, want.Consumers[0].Consumer.Username, state.Consumers[0].Consumer.Username)
		assert.Equal(t, want.Consumers[0].Consumer.CustomID, state.Consumers[0].Consumer.CustomID)
		assert.Equal(t, want.Consumers[0].KeyAuths[0].Key, state.Consumers[0].KeyAuths[0].Key)
		assert.Equal(t, want.Consumers[0].Oauth2Creds[0].ClientID, state.Consumers[0].Oauth2Creds[0].ClientID)
		assert.Equal(t, want.Consumers[0].Oauth2Creds[0].ClientSecret, state.Consumers[0].Oauth2Creds[0].ClientSecret)
		assert.Equal(t, want.Consumers[0].Oauth2Creds[0].HashSecret, state.Consumers[0].Oauth2Creds[0].HashSecret)
		assert.Equal(t, want.Consumers[0].Oauth2Creds[0].RedirectURIs, state.Consumers[0].Oauth2Creds[0].RedirectURIs)
	})
}
