package translator

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"sigs.k8s.io/controller-runtime/pkg/client"

	incubatorv1alpha1 "github.com/kong/kubernetes-configuration/api/incubator/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
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
				ServiceNameToParent:   map[string]client.Object{},
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
				ServiceNameToParent:   map[string]client.Object{},
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
				// not all of the inputs have proper parents, so while they're not empty, this still is in the unit
				ServiceNameToParent: map[string]client.Object{},
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
					ServiceNameToParent:   map[string]client.Object{},
				},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      makeSecretToSNIs(map[string]testSNIs{"a": {hosts: []string{"b", "c"}}, "d": {hosts: []string{"e", "f"}}, "g": {hosts: []string{"h"}}}),
				ServiceNameToServices: map[string]kongstate.Service{"1": {Namespace: "potato"}, "2": {Namespace: "carrot"}},
				ServiceNameToParent:   map[string]client.Object{},
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
				ServiceNameToParent:   map[string]client.Object{},
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
				ServiceNameToParent:   map[string]client.Object{},
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
				ServiceNameToParent:   map[string]client.Object{},
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
				ServiceNameToParent:   map[string]client.Object{},
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

func TestAddFromIngressV1TLS(t *testing.T) {
	parentIngress := &netv1.Ingress{ObjectMeta: metav1.ObjectMeta{Namespace: "foo"}}

	type args struct {
		tlsSections []netv1.IngressTLS
	}
	tests := []struct {
		name string
		args args
		want map[string]testSNIs
	}{
		{
			name: "different secrets with no overlapping hosts",
			args: args{
				tlsSections: []netv1.IngressTLS{
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
				tlsSections: []netv1.IngressTLS{
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
			m.addFromIngressV1TLS(tt.args.tlsSections, parentIngress)

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
		expectedFailures    []string
	}{
		{
			name:      "if all backends have a service then all services will be returned and their annotations recorded",
			namespace: corev1.NamespaceDefault,
			backends: kongstate.ServiceBackends{
				builder.NewKongstateServiceBackend("test-service1").
					WithNamespace(corev1.NamespaceDefault).
					MustBuild(),
				builder.NewKongstateServiceBackend("test-service2").
					WithNamespace(corev1.NamespaceDefault).
					MustBuild(),
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
				builder.NewKongstateServiceBackend("test-service1").
					WithNamespace(corev1.NamespaceDefault).
					MustBuild(),
				builder.NewKongstateServiceBackend("test-service2").
					WithNamespace(corev1.NamespaceDefault).
					MustBuild(),
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
			expectedFailures: []string{
				"failed to resolve Kubernetes Service for backend: failed to fetch Service default/test-service2: Service default/test-service2 not found",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			parent := &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{Name: "ingress", Namespace: tt.namespace},
				TypeMeta:   metav1.TypeMeta{Kind: "Ingress", APIVersion: netv1.SchemeGroupVersion.String()},
			}
			storer, err := store.NewFakeStore(store.FakeObjects{Services: tt.services})
			require.NoError(t, err)

			failuresCollector := failures.NewResourceFailuresCollector(logr.Discard())
			translatedObjectsCollector := NewObjectsCollector()

			services, annotations := getK8sServicesForBackends(storer, tt.backends, translatedObjectsCollector, failuresCollector, parent)
			assert.Equal(t, tt.expectedServices, services)
			assert.Equal(t, tt.expectedAnnotations, annotations)
			var collectedFailures []string
			for _, failure := range failuresCollector.PopResourceFailures() {
				collectedFailures = append(collectedFailures, failure.Message())
			}
			assert.Equal(t, tt.expectedFailures, collectedFailures)
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
				"Service has inconsistent konghq.com/foo annotation and is used in multi-Service backend",
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
				"Service has inconsistent konghq.com/foo annotation and is used in multi-Service backend",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(zap.InfoLevel)
			logger := zapr.NewLogger(zap.New(core))

			failuresCollector := failures.NewResourceFailuresCollector(logger)
			assert.Equal(t, tt.expected, collectInconsistentAnnotations(tt.services, tt.annotations, failuresCollector, ""))
			assert.Len(t, failuresCollector.PopResourceFailures(), len(tt.expectedLogEntries), "expecting as many translation failures as log entries")
			for i := range tt.expectedLogEntries {
				assert.Contains(t, logs.All()[i].Entry.Message, tt.expectedLogEntries[i])
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
						builder.NewKongstateServiceBackend("k8s-service-to-skip1").
							WithNamespace("test-namespace").MustBuild(),
						builder.NewKongstateServiceBackend("k8s-service-to-skip2").
							WithNamespace("test-namespace").MustBuild(),
					},
				},
				"service-to-keep": {
					Service: kong.Service{
						Name: lo.ToPtr("service-to-skip"),
					},
					Namespace: "test-namespace",
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("k8s-service-to-keep1").
							WithNamespace("test-namespace").MustBuild(),
						builder.NewKongstateServiceBackend("k8s-service-to-keep2").
							WithNamespace("test-namespace").MustBuild(),
					},
				},
			},
			serviceNamesToSkip: map[string]interface{}{
				"service-to-skip": nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ingressRules := newIngressRules()
			fakeStore, err := store.NewFakeStore(store.FakeObjects{
				Services: tc.k8sServices,
			})
			require.NoError(t, err)
			ingressRules.ServiceNameToServices = tc.serviceNamesToServices
			logger := zapr.NewLogger(zap.NewNop())
			failuresCollector := failures.NewResourceFailuresCollector(logger)
			translatedObjectsCollector := NewObjectsCollector()
			servicesToBeSkipped := ingressRules.populateServices(logger, fakeStore, failuresCollector, translatedObjectsCollector)
			require.Equal(t, tc.serviceNamesToSkip, servicesToBeSkipped)
		})
	}
}

func TestResolveKubernetesServiceForBackend(t *testing.T) {
	testService := func(annotations map[string]string) *corev1.Service {
		return &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "test-service",
				Namespace:   "test-namespace",
				Annotations: annotations,
			},
		}
	}

	testCases := []struct {
		name                      string
		storerObjects             store.FakeObjects
		backend                   kongstate.ServiceBackend
		ingressNamespace          string
		expectedService           *corev1.Service
		expectErrorContains       string
		expectedTranslatedObjects []client.Object
	}{
		{
			name: "backend is an existing service",
			storerObjects: store.FakeObjects{
				Services: []*corev1.Service{testService(nil)},
			},
			backend: builder.NewKongstateServiceBackend("test-service").
				WithNamespace("test-namespace").
				WithPortNumber(80).
				MustBuild(),
			expectedService: testService(nil),
		},
		{
			name:          "backend is not an existing service",
			storerObjects: store.FakeObjects{},
			backend: builder.NewKongstateServiceBackend("test-service").
				WithNamespace("test-namespace").
				WithPortNumber(80).
				MustBuild(),
			expectErrorContains: "Service test-namespace/test-service not found",
		},
		{
			name: "backend is an existing KongServiceFacade with annotations",
			storerObjects: store.FakeObjects{
				KongServiceFacades: []*incubatorv1alpha1.KongServiceFacade{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-service-facade",
							Namespace: "test-namespace",
							Annotations: map[string]string{
								"common": "common-from-facade",
								"facade": "facade-from-facade",
							},
						},
						Spec: incubatorv1alpha1.KongServiceFacadeSpec{
							Backend: incubatorv1alpha1.KongServiceFacadeBackend{
								Name: "test-service",
								Port: 80,
							},
						},
					},
				},
				Services: []*corev1.Service{testService(map[string]string{
					"common":  "common-from-service",
					"service": "service-from-service",
				})},
			},
			backend: builder.NewKongstateServiceBackend("test-service-facade").
				WithType(kongstate.ServiceBackendTypeKongServiceFacade).
				WithNamespace("test-namespace").
				WithPortNumber(80).
				MustBuild(),
			expectedService: testService(map[string]string{
				"common":  "common-from-facade",
				"facade":  "facade-from-facade",
				"service": "service-from-service",
			}),
			expectedTranslatedObjects: []client.Object{
				&incubatorv1alpha1.KongServiceFacade{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service-facade",
						Namespace: "test-namespace",
					},
				},
			},
		},
		{
			name: "backend is an existing KongServiceFacade with no annotations",
			storerObjects: store.FakeObjects{
				KongServiceFacades: []*incubatorv1alpha1.KongServiceFacade{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-service-facade",
							Namespace: "test-namespace",
							Annotations: map[string]string{
								"common": "common-from-facade",
								"facade": "facade-from-facade",
							},
						},
						Spec: incubatorv1alpha1.KongServiceFacadeSpec{
							Backend: incubatorv1alpha1.KongServiceFacadeBackend{
								Name: "test-service",
								Port: 80,
							},
						},
					},
				},
				Services: []*corev1.Service{testService(nil)},
			},
			backend: builder.NewKongstateServiceBackend("test-service-facade").
				WithType(kongstate.ServiceBackendTypeKongServiceFacade).
				WithNamespace("test-namespace").
				WithPortNumber(80).
				MustBuild(),
			expectedService: testService(map[string]string{
				"common": "common-from-facade",
				"facade": "facade-from-facade",
			}),
			expectedTranslatedObjects: []client.Object{
				&incubatorv1alpha1.KongServiceFacade{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service-facade",
						Namespace: "test-namespace",
					},
				},
			},
		},
		{
			name: "backend is an existing KongServiceFacade referring not existing Service",
			storerObjects: store.FakeObjects{
				KongServiceFacades: []*incubatorv1alpha1.KongServiceFacade{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-service-facade",
							Namespace: "test-namespace",
							Annotations: map[string]string{
								"common": "common-from-facade",
								"facade": "facade-from-facade",
							},
						},
						Spec: incubatorv1alpha1.KongServiceFacadeSpec{
							Backend: incubatorv1alpha1.KongServiceFacadeBackend{
								Name: "not-existing-service",
								Port: 80,
							},
						},
					},
				},
			},
			backend: builder.NewKongstateServiceBackend("test-service-facade").
				WithType(kongstate.ServiceBackendTypeKongServiceFacade).
				WithNamespace("test-namespace").
				WithPortNumber(80).
				MustBuild(),
			expectErrorContains: "Service test-namespace/not-existing-service not found",
		},
		{
			name:          "backend is not existing KongServiceFacade",
			storerObjects: store.FakeObjects{},
			backend: builder.NewKongstateServiceBackend("not-existing-service-facade").
				WithType(kongstate.ServiceBackendTypeKongServiceFacade).
				WithNamespace("test-namespace").
				WithPortNumber(80).
				MustBuild(),
			expectErrorContains: "KongServiceFacade test-namespace/not-existing-service-facade not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeStore := lo.Must(store.NewFakeStore(tc.storerObjects))
			translatedObjectsCollector := NewObjectsCollector()
			service, err := resolveKubernetesServiceForBackend(fakeStore, tc.backend, translatedObjectsCollector)
			if tc.expectErrorContains != "" {
				require.ErrorContains(t, err, tc.expectErrorContains)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedService, service)
			gotTranslatedObjects := translatedObjectsCollector.Pop()
			for _, expectedObject := range tc.expectedTranslatedObjects {
				require.True(t, lo.ContainsBy(gotTranslatedObjects, func(obj client.Object) bool {
					return obj.GetNamespace()+"/"+obj.GetName() ==
						expectedObject.GetNamespace()+"/"+expectedObject.GetName()
				}), "expected translated object not found in actual translated objects")
			}
		})
	}
}

func TestResolveKubernetesServiceForBackend_DoesNotModifyCache(t *testing.T) {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "test-namespace",
			Annotations: map[string]string{
				"service": "from-service",
			},
		},
	}
	// Preserve a copy to compare against later.
	svcCopy := svc.DeepCopy()

	kongServiceFacade := &incubatorv1alpha1.KongServiceFacade{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service-facade",
			Namespace: "test-namespace",
			Annotations: map[string]string{
				"facade": "from-facade",
			},
		},
		Spec: incubatorv1alpha1.KongServiceFacadeSpec{
			Backend: incubatorv1alpha1.KongServiceFacadeBackend{
				Name: "test-service",
				Port: 80,
			},
		},
	}
	fakeStore := lo.Must(store.NewFakeStore(store.FakeObjects{
		Services:           []*corev1.Service{svc},
		KongServiceFacades: []*incubatorv1alpha1.KongServiceFacade{kongServiceFacade},
	}))
	backend := builder.NewKongstateServiceBackend("test-service-facade").
		WithNamespace("test-namespace").
		WithPortNumber(80).
		WithType(kongstate.ServiceBackendTypeKongServiceFacade).
		MustBuild()

	translatedObjectsCollector := NewObjectsCollector()
	resolvedService, err := resolveKubernetesServiceForBackend(fakeStore, backend, translatedObjectsCollector)
	require.NoError(t, err)
	require.Equal(t, svcCopy, svc, "service stored in cache should not be modified")
	require.Equal(t, resolvedService.Annotations, map[string]string{
		"service": "from-service",
		"facade":  "from-facade",
	}, "annotations should be merged in the returned service")
}
