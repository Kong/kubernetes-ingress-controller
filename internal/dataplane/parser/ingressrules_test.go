package parser

import (
	"bytes"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

type testSNIs struct {
	parents []client.Object
	hosts   []string
}

func makeSecretToSNIs(in map[string]testSNIs) SecretNameToSNIs {
	s := newSecretNameToSNIs()
	for k, v := range in {
		s.addUniqueParents(k, v.parents...)
		s.addUniqueHosts(k, v.hosts...)
	}
	return s
}

func TestMergeIngressRules(t *testing.T) {
	var (
		parent1 = &netv1.Ingress{ObjectMeta: metav1.ObjectMeta{UID: uuid.NewUUID()}}
		parent2 = &netv1.Ingress{ObjectMeta: metav1.ObjectMeta{UID: uuid.NewUUID()}}
	)

	for _, tt := range []struct {
		name                  string
		inputs                []ingressRules
		inputSecretNameToSNIs []map[string]testSNIs
		wantOutput            *ingressRules
	}{
		{
			name:                  "empty list",
			inputSecretNameToSNIs: []map[string]testSNIs{},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      newSecretNameToSNIs(),
				ServiceNameToServices: map[string]kongstate.Service{},
			},
		},
		{
			name: "nil maps",
			inputs: []ingressRules{
				{}, {}, {},
			},
			inputSecretNameToSNIs: []map[string]testSNIs{},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      newSecretNameToSNIs(),
				ServiceNameToServices: map[string]kongstate.Service{},
			},
		},
		{
			name:                  "one input",
			inputSecretNameToSNIs: []map[string]testSNIs{{"a": {hosts: []string{"b", "c"}}, "d": {hosts: []string{"e", "f"}}}},
			inputs: []ingressRules{
				{
					ServiceNameToServices: map[string]kongstate.Service{"1": {Namespace: "potato"}},
				},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      makeSecretToSNIs(map[string]testSNIs{"a": {hosts: []string{"b", "c"}}, "d": {hosts: []string{"e", "f"}}}),
				ServiceNameToServices: map[string]kongstate.Service{"1": {Namespace: "potato"}},
			},
		},
		{
			name: "three inputs",
			inputSecretNameToSNIs: []map[string]testSNIs{
				{"a": {hosts: []string{"b", "c"}}, "d": {hosts: []string{"e", "f"}}},
				{"g": {hosts: []string{"h"}}},
			},
			inputs: []ingressRules{
				{
					ServiceNameToServices: map[string]kongstate.Service{"1": {Namespace: "potato"}},
				},
				{
					ServiceNameToServices: map[string]kongstate.Service{"2": {Namespace: "carrot"}},
				},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      makeSecretToSNIs(map[string]testSNIs{"a": {hosts: []string{"b", "c"}}, "d": {hosts: []string{"e", "f"}}, "g": {hosts: []string{"h"}}}),
				ServiceNameToServices: map[string]kongstate.Service{"1": {Namespace: "potato"}, "2": {Namespace: "carrot"}},
			},
		},
		{
			name: "can merge SNI arrays",
			inputSecretNameToSNIs: []map[string]testSNIs{
				{"a": {hosts: []string{"b", "c"}}},
				{"a": {hosts: []string{"d", "e"}}},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      makeSecretToSNIs(map[string]testSNIs{"a": {hosts: []string{"b", "c", "d", "e"}}}),
				ServiceNameToServices: map[string]kongstate.Service{},
			},
		},
		{
			name: "can merge SNI arrays with parents",
			inputSecretNameToSNIs: []map[string]testSNIs{
				{"a": {parents: []client.Object{parent1}, hosts: []string{"b", "c"}}},
				{"a": {parents: []client.Object{parent2}, hosts: []string{"d", "e"}}},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      makeSecretToSNIs(map[string]testSNIs{"a": {parents: []client.Object{parent1, parent2}, hosts: []string{"b", "c", "d", "e"}}}),
				ServiceNameToServices: map[string]kongstate.Service{},
			},
		},
		{
			name: "can merge SNI arrays with repeating parents",
			inputSecretNameToSNIs: []map[string]testSNIs{
				{"a": {parents: []client.Object{parent1, parent2}, hosts: []string{"b", "c"}}},
				{"a": {parents: []client.Object{parent1, parent2}, hosts: []string{"d", "e"}}},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      makeSecretToSNIs(map[string]testSNIs{"a": {parents: []client.Object{parent1, parent2, parent1, parent2}, hosts: []string{"b", "c", "d", "e"}}}),
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
				SecretNameToSNIs:      newSecretNameToSNIs(),
				ServiceNameToServices: map[string]kongstate.Service{"svc-name": {Namespace: "new"}},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput := mergeIngressRules(tt.inputs...)
			for _, inputSecret := range tt.inputSecretNameToSNIs {
				gotOutput.SecretNameToSNIs.merge(makeSecretToSNIs(inputSecret))
			}
			assert.Equal(t, tt.wantOutput, &gotOutput)
		})
	}
}

func TestAddFromIngressV1beta1TLS(t *testing.T) {
	parentIngress := &netv1beta1.Ingress{ObjectMeta: metav1.ObjectMeta{Namespace: "foo"}}

	type args struct {
		tlsSections []netv1beta1.IngressTLS
	}
	tests := []struct {
		name string
		args args
		want map[string]testSNIs
	}{
		{
			name: "different secrets with no overlapping hosts",
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
			},
			want: map[string]testSNIs{
				"foo/sooper-secret":  {hosts: []string{"1.example.com", "2.example.com"}, parents: []client.Object{parentIngress}},
				"foo/sooper-secret2": {hosts: []string{"3.example.com", "4.example.com"}, parents: []client.Object{parentIngress}},
			},
		},
		{
			name: "different secrets with one overlapping host",
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
			},
			want: map[string]testSNIs{
				"foo/sooper-secret":  {hosts: []string{"1.example.com"}, parents: []client.Object{parentIngress}},
				"foo/sooper-secret2": {hosts: []string{"3.example.com", "4.example.com"}, parents: []client.Object{parentIngress}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newSecretNameToSNIs()
			m.addFromIngressV1TLS(v1beta1toV1TLS(tt.args.tlsSections), parentIngress)

			for k, v := range tt.want {
				assert.ElementsMatch(t, v.parents, m.Parents(k))
				assert.ElementsMatch(t, v.hosts, m.Hosts(k))
			}
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
					TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: corev1.SchemeGroupVersion.String()},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service1",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							"konghq.com/bar": "foo",
							"konghq.com/baz": "foo",
						},
					},
				},
				{
					TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: corev1.SchemeGroupVersion.String()},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service2",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							"konghq.com/baz": "foo",
							"konghq.com/foo": "bar",
							"konghq.com/bar": "foo",
						},
					},
				},
				{
					TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: corev1.SchemeGroupVersion.String()},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service3",
						Namespace: corev1.NamespaceDefault,
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
				"Service has inconsistent konghq.com/foo annotation and is used in multi-Service backend",
			},
		},
		{
			name: "validation fails if all services have the same annotations, but not the same value",
			services: []*corev1.Service{
				{
					TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: corev1.SchemeGroupVersion.String()},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service1",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							"konghq.com/foo": "bar",
						},
					},
				},
				{
					TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: corev1.SchemeGroupVersion.String()},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service2",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							"konghq.com/foo": "baz",
						},
					},
				},
				{
					TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: corev1.SchemeGroupVersion.String()},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service3",
						Namespace: corev1.NamespaceDefault,
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
				"Service has inconsistent konghq.com/foo annotation and is used in multi-Service backend",
				"Service has inconsistent konghq.com/foo annotation and is used in multi-Service backend",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			logger, loggerHook := test.NewNullLogger()
			failuresCollector, err := failures.NewResourceFailuresCollector(logger)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, servicesAllUseTheSameKongAnnotations(tt.services, tt.annotations, failuresCollector, ""))
			assert.Len(t, failuresCollector.PopResourceFailures(), len(tt.expectedLogEntries), "expecting as many translation failures as log entries")
			for i := range tt.expectedLogEntries {
				assert.Contains(t, loggerHook.AllEntries()[i].Message, tt.expectedLogEntries[i])
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
					TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: corev1.SchemeGroupVersion.String()},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "k8s-service-to-skip1",
						Namespace: "test-namespace",
						Annotations: map[string]string{
							"konghq.com/foo": "bar",
						},
					},
				},
				{
					TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: corev1.SchemeGroupVersion.String()},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "k8s-service-to-skip2",
						Namespace: "test-namespace",
					},
				},
				{
					TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: corev1.SchemeGroupVersion.String()},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "k8s-service-to-keep1",
						Namespace: "test-namespace",
						Annotations: map[string]string{
							"konghq.com/foo": "bar",
						},
					},
				},
				{
					TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: corev1.SchemeGroupVersion.String()},
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
						Name: lo.ToPtr("service-to-skip"),
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
						Name: lo.ToPtr("service-to-skip"),
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
			logger, _ := test.NewNullLogger()
			failuresCollector, err := failures.NewResourceFailuresCollector(logger)
			require.NoError(t, err)
			servicesToBeSkipped := ingressRules.populateServices(logrus.New(), fakeStore, failuresCollector)
			require.Equal(t, tc.serviceNamesToSkip, servicesToBeSkipped)
			require.Len(t, failuresCollector.PopResourceFailures(), len(servicesToBeSkipped), "expecting as many translation failures as services to skip")
		})
	}
}
