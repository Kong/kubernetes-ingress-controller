package kongstate

import (
	"reflect"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser/util"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func Test_getPluginRelations(t *testing.T) {
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
							K8sService: corev1.Service{
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
							K8sService: corev1.Service{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns1",
									Annotations: map[string]string{
										annotations.AnnotationPrefix + annotations.PluginsKey: "foo,bar",
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
