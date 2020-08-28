package parser

import (
	"reflect"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_getPluginRelations(t *testing.T) {
	type args struct {
		state KongState
	}
	tests := []struct {
		name string
		args args
		want map[string]foreignRelations
	}{
		{
			name: "empty state",
			want: map[string]foreignRelations{},
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
							k8sKongConsumer: configurationv1.KongConsumer{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns1",
									Annotations: map[string]string{
										annotations.DeprecatedPluginsKey: "foo,bar",
									},
								},
							},
						},
					},
				},
			},
			want: map[string]foreignRelations{
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
										annotations.DeprecatedPluginsKey: "foo,bar",
									},
								},
							},
						},
					},
				},
			},
			want: map[string]foreignRelations{
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
									Ingress: networking.Ingress{
										ObjectMeta: metav1.ObjectMeta{
											Name:      "some-ingress",
											Namespace: "ns2",
											Annotations: map[string]string{
												annotations.DeprecatedPluginsKey: "foo,bar",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: map[string]foreignRelations{
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
									Ingress: networking.Ingress{
										ObjectMeta: metav1.ObjectMeta{
											Name:      "some-ingress",
											Namespace: "ns2",
											Annotations: map[string]string{
												annotations.DeprecatedPluginsKey: "foo,bar",
											},
										},
									},
								},
								{
									Route: kong.Route{
										Name: kong.String("bar-route"),
									},
									Ingress: networking.Ingress{
										ObjectMeta: metav1.ObjectMeta{
											Name:      "some-ingress",
											Namespace: "ns2",
											Annotations: map[string]string{
												annotations.DeprecatedPluginsKey: "bar,baz",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: map[string]foreignRelations{
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
							k8sKongConsumer: configurationv1.KongConsumer{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns1",
									Annotations: map[string]string{
										annotations.DeprecatedPluginsKey: "foo,bar",
									},
								},
							},
						},
						{
							Consumer: kong.Consumer{
								Username: kong.String("foo-consumer"),
							},
							k8sKongConsumer: configurationv1.KongConsumer{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns2",
									Annotations: map[string]string{
										annotations.DeprecatedPluginsKey: "foo,bar",
									},
								},
							},
						},
						{
							Consumer: kong.Consumer{
								Username: kong.String("bar-consumer"),
							},
							k8sKongConsumer: configurationv1.KongConsumer{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns1",
									Annotations: map[string]string{
										annotations.DeprecatedPluginsKey: "foobar",
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
										annotations.DeprecatedPluginsKey: "foo,bar",
									},
								},
							},
							Routes: []Route{
								{
									Route: kong.Route{
										Name: kong.String("foo-route"),
									},
									Ingress: networking.Ingress{
										ObjectMeta: metav1.ObjectMeta{
											Name:      "some-ingress",
											Namespace: "ns2",
											Annotations: map[string]string{
												annotations.DeprecatedPluginsKey: "foo,bar",
											},
										},
									},
								},
								{
									Route: kong.Route{
										Name: kong.String("bar-route"),
									},
									Ingress: networking.Ingress{
										ObjectMeta: metav1.ObjectMeta{
											Name:      "some-ingress",
											Namespace: "ns2",
											Annotations: map[string]string{
												annotations.DeprecatedPluginsKey: "bar,baz",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: map[string]foreignRelations{
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
			if got := getPluginRelations(tt.args.state); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPluginRelations() = %v, want %v", got, tt.want)
			}
		})
	}
}
