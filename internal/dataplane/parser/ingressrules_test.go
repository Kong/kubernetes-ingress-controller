package parser

import (
	"bytes"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

func TestMergeIngressRules(t *testing.T) {
	for _, tt := range []struct {
		name       string
		inputs     []ingressRules
		wantOutput *ingressRules
	}{
		{
			name: "empty list",
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{},
				ServiceNameToServices: map[string]kongstate.Service{},
			},
		},
		{
			name: "nil maps",
			inputs: []ingressRules{
				{}, {}, {},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{},
				ServiceNameToServices: map[string]kongstate.Service{},
			},
		},
		{
			name: "one input",
			inputs: []ingressRules{
				{
					SecretNameToSNIs:      map[string][]string{"a": {"b", "c"}, "d": {"e", "f"}},
					ServiceNameToServices: map[string]kongstate.Service{"1": {Namespace: "potato"}},
				},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{"a": {"b", "c"}, "d": {"e", "f"}},
				ServiceNameToServices: map[string]kongstate.Service{"1": {Namespace: "potato"}},
			},
		},
		{
			name: "three inputs",
			inputs: []ingressRules{
				{
					SecretNameToSNIs:      map[string][]string{"a": {"b", "c"}, "d": {"e", "f"}},
					ServiceNameToServices: map[string]kongstate.Service{"1": {Namespace: "potato"}},
				},
				{
					SecretNameToSNIs: map[string][]string{"g": {"h"}},
				},
				{
					ServiceNameToServices: map[string]kongstate.Service{"2": {Namespace: "carrot"}},
				},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{"a": {"b", "c"}, "d": {"e", "f"}, "g": {"h"}},
				ServiceNameToServices: map[string]kongstate.Service{"1": {Namespace: "potato"}, "2": {Namespace: "carrot"}},
			},
		},
		{
			name: "can merge SNI arrays",
			inputs: []ingressRules{
				{
					SecretNameToSNIs: map[string][]string{"a": {"b", "c"}},
				},
				{
					SecretNameToSNIs: map[string][]string{"a": {"d", "e"}},
				},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{"a": {"b", "c", "d", "e"}},
				ServiceNameToServices: map[string]kongstate.Service{},
			},
		},
		{
			name: "overwrites services",
			inputs: []ingressRules{
				{
					ServiceNameToServices: map[string]kongstate.Service{"svc-name": {Namespace: "old"}},
				},
				{
					ServiceNameToServices: map[string]kongstate.Service{"svc-name": {Namespace: "new"}},
				},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{},
				ServiceNameToServices: map[string]kongstate.Service{"svc-name": {Namespace: "new"}},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput := mergeIngressRules(tt.inputs...)
			assert.Equal(t, &gotOutput, tt.wantOutput)
		})
	}
}

func TestAddFromIngressV1beta1TLS(t *testing.T) {
	type args struct {
		tlsSections []netv1beta1.IngressTLS
		namespace   string
	}
	tests := []struct {
		name string
		args args
		want SecretNameToSNIs
	}{
		{
			args: args{
				tlsSections: []netv1beta1.IngressTLS{
					{
						Hosts: []string{
							"1.example.com",
							"2.example.com",
						},
						SecretName: "sooper-secret",
					},
					{
						Hosts: []string{
							"3.example.com",
							"4.example.com",
						},
						SecretName: "sooper-secret2",
					},
				},
				namespace: "foo",
			},
			want: SecretNameToSNIs{
				"foo/sooper-secret":  {"1.example.com", "2.example.com"},
				"foo/sooper-secret2": {"3.example.com", "4.example.com"},
			},
		},
		{
			args: args{
				tlsSections: []netv1beta1.IngressTLS{
					{
						Hosts: []string{
							"1.example.com",
						},
						SecretName: "sooper-secret",
					},
					{
						Hosts: []string{
							"3.example.com",
							"1.example.com",
							"4.example.com",
						},
						SecretName: "sooper-secret2",
					},
				},
				namespace: "foo",
			},
			want: SecretNameToSNIs{
				"foo/sooper-secret":  {"1.example.com"},
				"foo/sooper-secret2": {"3.example.com", "4.example.com"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newSecretNameToSNIs()
			m.addFromIngressV1beta1TLS(tt.args.tlsSections, tt.args.namespace)
			assert.Equal(t, m, tt.want)
		})
	}
}

func TestGetK8sServicesForBackends(t *testing.T) {
	for _, tt := range []struct {
		name                string
		namespace           string
		backends            kongstate.ServiceBackends
		services            []*corev1.Service
		expectedServices    []*corev1.Service
		expectedAnnotations map[string]string
		expectedLogEntries  []string
	}{
		{
			name:      "if all backends have a service then all services will be returned and their annotations recorded",
			namespace: corev1.NamespaceDefault,
			backends: kongstate.ServiceBackends{
				{
					Name: "test-service1",
				},
				{
					Name: "test-service2",
				},
			},
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service1",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							"konghq.com/foo": "bar",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service2",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							"konghq.com/foo": "baz",
						},
					},
				},
			},
			expectedServices: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service1",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							"konghq.com/foo": "bar",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service2",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							"konghq.com/foo": "baz",
						},
					},
				},
			},
			expectedAnnotations: map[string]string{
				"konghq.com/foo": "baz",
			},
		},
		{
			name:      "backends which have no corresponding services will fail to fetch",
			namespace: corev1.NamespaceDefault,
			backends: kongstate.ServiceBackends{
				{
					Name: "test-service1",
				},
				{
					Name: "test-service2",
				},
			},
			services: []*corev1.Service{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service1",
					Namespace: corev1.NamespaceDefault,
				},
			}},
			expectedServices: []*corev1.Service{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service1",
					Namespace: corev1.NamespaceDefault,
				},
			}},
			expectedAnnotations: map[string]string{},
			expectedLogEntries: []string{
				"failed to fetch service: Service default/test-service2 not found",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			storer, err := store.NewFakeStore(store.FakeObjects{Services: tt.services})
			require.NoError(t, err)

			stdout := new(bytes.Buffer)
			logger := logrus.New()
			logger.SetOutput(stdout)

			services, annotations := getK8sServicesForBackends(logger, storer, tt.namespace, tt.backends)
			assert.Equal(t, tt.expectedServices, services)
			assert.Equal(t, tt.expectedAnnotations, annotations)
			for _, expectedLogEntry := range tt.expectedLogEntries {
				assert.Contains(t, stdout.String(), expectedLogEntry)
			}
		})
	}
}

func TestDoK8sServicesMatchAnnotations(t *testing.T) {
	for _, tt := range []struct {
		name               string
		services           []*corev1.Service
		annotations        map[string]string
		expected           bool
		expectedLogEntries []string
	}{
		{
			name:        "if no services are provided, then there's no validation failure",
			annotations: map[string]string{"foo": "bar"},
			expected:    true,
		},
		{
			name: "validation passes for a group of services with no annotations expected, even if they all have different annotations",
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-service1",
						Annotations: map[string]string{
							"konghq.com/foo": "bar",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-service2",
						Annotations: map[string]string{
							"konghq.com/bar": "foo",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-service3",
						Annotations: map[string]string{
							"konghq.com/baz": "foo",
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "validation passes for a group of services all have the expected annotations",
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-service1",
						Annotations: map[string]string{
							"konghq.com/foo": "bar",
							"konghq.com/bar": "foo",
							"konghq.com/baz": "foo",
							"example.com":    "foo",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-service2",
						Annotations: map[string]string{
							"konghq.com/baz": "foo",
							"konghq.com/foo": "bar",
							"konghq.com/bar": "foo",
							"example.com":    "bar",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-service3",
						Annotations: map[string]string{
							"konghq.com/bar": "foo",
							"konghq.com/foo": "bar",
							"konghq.com/baz": "foo",
						},
					},
				},
			},
			annotations: map[string]string{
				"konghq.com/foo": "bar",
				"konghq.com/bar": "foo",
				"konghq.com/baz": "foo",
			},
			expected: true,
		},
		{
			name: "validation fails if one service does not have all expected annotations",
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-service1",
						Annotations: map[string]string{
							"konghq.com/bar": "foo",
							"konghq.com/baz": "foo",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-service2",
						Annotations: map[string]string{
							"konghq.com/baz": "foo",
							"konghq.com/foo": "bar",
							"konghq.com/bar": "foo",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-service3",
						Annotations: map[string]string{
							"konghq.com/bar": "foo",
							"konghq.com/foo": "bar",
							"konghq.com/baz": "foo",
						},
					},
				},
			},
			annotations: map[string]string{
				"konghq.com/foo": "bar",
				"konghq.com/bar": "foo",
				"konghq.com/baz": "foo",
			},
			expected: false,
			expectedLogEntries: []string{
				"in the backend group of 3 kubernetes services some have the konghq.com/foo annotation while others don't",
			},
		},
		{
			name: "validation fails if all services have the same annotations, but not the same value",
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-service1",
						Annotations: map[string]string{
							"konghq.com/foo": "bar",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-service2",
						Annotations: map[string]string{
							"konghq.com/foo": "baz",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-service3",
						Annotations: map[string]string{
							"konghq.com/foo": "buzz",
						},
					},
				},
			},
			annotations: map[string]string{
				"konghq.com/foo": "bar",
			},
			expected: false,
			expectedLogEntries: []string{
				"the value of annotation konghq.com/foo is different between the 3 services which comprise this backend.",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stdout := new(bytes.Buffer)
			logger := logrus.New()
			logger.SetOutput(stdout)
			assert.Equal(t, tt.expected, servicesAllUseTheSameKongAnnotations(logger, tt.services, tt.annotations))
			for _, expectedLogEntry := range tt.expectedLogEntries {
				assert.Contains(t, stdout.String(), expectedLogEntry)
			}
		})
	}
}

func TestPopulateServices(t *testing.T) {
	testCases := []struct {
		name                   string
		k8sServices            []*corev1.Service
		serviceNamesToServices map[string]kongstate.Service
		serviceNamesToSkip     map[string]interface{}
	}{
		{
			name: "one service to skip, one service to keep",
			k8sServices: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "k8s-service-to-skip1",
						Namespace: "test-namespace",
						Annotations: map[string]string{
							"konghq.com/foo": "bar",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "k8s-service-to-skip2",
						Namespace: "test-namespace",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "k8s-service-to-keep1",
						Namespace: "test-namespace",
						Annotations: map[string]string{
							"konghq.com/foo": "bar",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "k8s-service-to-keep2",
						Namespace: "test-namespace",
						Annotations: map[string]string{
							"konghq.com/foo": "bar",
						},
					},
				},
			},
			serviceNamesToServices: map[string]kongstate.Service{
				"service-to-skip": {
					Service: kong.Service{
						Name: pointer.StringPtr("service-to-skip"),
					},
					Namespace: "test-namespace",
					Backends: []kongstate.ServiceBackend{
						{
							Name:      "k8s-service-to-skip1",
							Namespace: "test-namespace",
						},
						{
							Name:      "k8s-service-to-skip2",
							Namespace: "test-namespace",
						},
					},
				},
				"service-to-keep": {
					Service: kong.Service{
						Name: pointer.StringPtr("service-to-skip"),
					},
					Namespace: "test-namespace",
					Backends: []kongstate.ServiceBackend{
						{
							Name:      "k8s-service-to-keep1",
							Namespace: "test-namespace",
						},
						{
							Name:      "k8s-service-to-keep2",
							Namespace: "test-namespace",
						},
					},
				},
			},
			serviceNamesToSkip: map[string]interface{}{
				"service-to-skip": nil,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			ingressRules := newIngressRules()
			fakeStore, err := store.NewFakeStore(store.FakeObjects{
				Services: tc.k8sServices,
			})
			require.NoError(t, err)
			ingressRules.ServiceNameToServices = tc.serviceNamesToServices
			servicesToBeSkipped := ingressRules.populateServices(logrus.New(), fakeStore)
			require.Equal(t, tc.serviceNamesToSkip, servicesToBeSkipped)
		})
	}
}
