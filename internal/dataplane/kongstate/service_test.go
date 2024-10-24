package kongstate

import (
	"reflect"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
)

func TestOverrideService(t *testing.T) {
	testCases := []struct {
		name                  string
		inService             Service
		k8sServiceAnnotations map[string]string
		expectedService       Service
	}{
		{
			name: "no overrides",
			inService: Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			expectedService: Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			k8sServiceAnnotations: map[string]string{},
		},
		{
			name: "override protocol to https",
			inService: Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			expectedService: Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			k8sServiceAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.ProtocolKey: "https",
			},
		},
		{
			name: "override retries to 0",
			inService: Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			expectedService: Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
					Retries:  kong.Int(0),
				},
			},
			k8sServiceAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.RetriesKey: "0",
			},
		},
		{
			name: "override retries to 1",
			inService: Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			expectedService: Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
					Retries:  kong.Int(1),
				},
			},
			k8sServiceAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.RetriesKey: "1",
			},
		},
		{
			name: "override path",
			inService: Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			expectedService: Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/new-path"),
				},
			},
			k8sServiceAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.PathKey: "/new-path",
			},
		},
		{
			name: "override connect timeout, read timeout, write timeout",
			inService: Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			expectedService: Service{
				Service: kong.Service{
					Host:           kong.String("foo.com"),
					Port:           kong.Int(80),
					Name:           kong.String("foo"),
					Protocol:       kong.String("http"),
					Path:           kong.String("/"),
					ConnectTimeout: kong.Int(100),
					ReadTimeout:    kong.Int(100),
					WriteTimeout:   kong.Int(100),
				},
			},
			k8sServiceAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.ConnectTimeoutKey: "100",
				annotations.AnnotationPrefix + annotations.ReadTimeoutKey:    "100",
				annotations.AnnotationPrefix + annotations.WriteTimeoutKey:   "100",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := tc.inService
			for _, k8sSvc := range service.K8sServices {
				service.overrideByAnnotation(k8sSvc.Annotations)
				require.Equal(t, tc.expectedService.Service, service.Service)
			}
		})
	}
}

func TestNilServiceOverrideDoesntPanic(t *testing.T) {
	require.NotPanics(t, func() {
		var nilService *Service
		nilService.override() //nolint:errcheck
	})
}

func TestOverrideServicePath(t *testing.T) {
	type args struct {
		service Service
		anns    map[string]string
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{},
		{name: "basic empty service"},
		{
			name: "set to valid value",
			args: args{
				anns: map[string]string{
					"konghq.com/path": "/foo",
				},
			},
			want: Service{
				Service: kong.Service{
					Path: kong.String("/foo"),
				},
			},
		},
		{
			name: "does not set path if doesn't start with /",
			args: args{
				anns: map[string]string{
					"konghq.com/path": "foo",
				},
			},
			want: Service{},
		},
		{
			name: "overrides any other value",
			args: args{
				service: Service{
					Service: kong.Service{
						Path: kong.String("/foo"),
					},
				},
				anns: map[string]string{
					"konghq.com/path": "/bar",
				},
			},
			want: Service{
				Service: kong.Service{
					Path: kong.String("/bar"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.service.overridePath(tt.args.anns)
			if !reflect.DeepEqual(tt.args.service, tt.want) {
				t.Errorf("overrideServicePath() got = %v, want %v", tt.args.service, tt.want)
			}
		})
	}
}

func TestOverrideConnectTimeout(t *testing.T) {
	type args struct {
		service Service
		anns    map[string]string
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "set to valid value",
			args: args{
				anns: map[string]string{
					"konghq.com/connect-timeout": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					ConnectTimeout: kong.Int(3000),
				},
			},
		},
		{
			name: "value cannot parse to int",
			args: args{
				anns: map[string]string{
					"konghq.com/connect-timeout": "burranyi yedigei",
				},
			},
			want: Service{},
		},
		{
			name: "overrides any other value",
			args: args{
				service: Service{
					Service: kong.Service{
						ConnectTimeout: kong.Int(2000),
					},
				},
				anns: map[string]string{
					"konghq.com/connect-timeout": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					ConnectTimeout: kong.Int(3000),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.service.overrideConnectTimeout(tt.args.anns)
			if !reflect.DeepEqual(tt.args.service, tt.want) {
				t.Errorf("overrideConnectTimeout() got = %v, want %v", tt.args.service, tt.want)
			}
		})
	}
}

func TestOverrideWriteTimeout(t *testing.T) {
	type args struct {
		service Service
		anns    map[string]string
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "set to valid value",
			args: args{
				anns: map[string]string{
					"konghq.com/write-timeout": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					WriteTimeout: kong.Int(3000),
				},
			},
		},
		{
			name: "value cannot parse to int",
			args: args{
				anns: map[string]string{
					"konghq.com/write-timeout": "burranyi yedigei",
				},
			},
			want: Service{},
		},
		{
			name: "overrides any other value",
			args: args{
				service: Service{
					Service: kong.Service{
						WriteTimeout: kong.Int(2000),
					},
				},
				anns: map[string]string{
					"konghq.com/write-timeout": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					WriteTimeout: kong.Int(3000),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.service.overrideWriteTimeout(tt.args.anns)
			if !reflect.DeepEqual(tt.args.service, tt.want) {
				t.Errorf("overrideWriteTimeout() got = %v, want %v", tt.args.service, tt.want)
			}
		})
	}
}

func TestOverrideReadTimeout(t *testing.T) {
	type args struct {
		service Service
		anns    map[string]string
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "set to valid value",
			args: args{
				anns: map[string]string{
					"konghq.com/read-timeout": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					ReadTimeout: kong.Int(3000),
				},
			},
		},
		{
			name: "value cannot parse to int",
			args: args{
				anns: map[string]string{
					"konghq.com/read-timeout": "burranyi yedigei",
				},
			},
			want: Service{},
		},
		{
			name: "overrides any other value",
			args: args{
				service: Service{
					Service: kong.Service{
						ReadTimeout: kong.Int(2000),
					},
				},
				anns: map[string]string{
					"konghq.com/read-timeout": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					ReadTimeout: kong.Int(3000),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.service.overrideReadTimeout(tt.args.anns)
			if !reflect.DeepEqual(tt.args.service, tt.want) {
				t.Errorf("overrideReadTimeout() got = %v, want %v", tt.args.service, tt.want)
			}
		})
	}
}

func TestOverrideRetries(t *testing.T) {
	type args struct {
		service Service
		anns    map[string]string
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "set to valid value",
			args: args{
				anns: map[string]string{
					"konghq.com/retries": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					Retries: kong.Int(3000),
				},
			},
		},
		{
			name: "value cannot parse to int",
			args: args{
				anns: map[string]string{
					"konghq.com/retries": "burranyi yedigei",
				},
			},
			want: Service{},
		},
		{
			name: "overrides any other value",
			args: args{
				service: Service{
					Service: kong.Service{
						Retries: kong.Int(2000),
					},
				},
				anns: map[string]string{
					"konghq.com/retries": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					Retries: kong.Int(3000),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.service.overrideRetries(tt.args.anns)
			if !reflect.DeepEqual(tt.args.service, tt.want) {
				t.Errorf("overrideRetries() got = %v, want %v", tt.args.service, tt.want)
			}
		})
	}
}

func TestServiceOverride_DeterministicOrderWhenMoreThan1KubernetesService(t *testing.T) {
	service := Service{
		Service: kong.Service{},
		K8sServices: map[string]*corev1.Service{
			"default/service-3": {
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.RetriesKey: "3",
					},
				},
			},
			"default/service-1": {
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.RetriesKey: "1",
					},
				},
			},
			"default/service-2": {
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.RetriesKey: "2",
					},
				},
			},
		},
	}

	// We expect default/service-3 to be the last one to be processed effectively overriding the previous annotations.
	const expectedRetries = 3
	require.NoError(t, service.override())
	require.Equal(t, expectedRetries, *service.Service.Retries)
}
