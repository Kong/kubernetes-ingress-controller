package kongstate

import (
	"fmt"
	"sort"
	"testing"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/kong/go-kong/kong/custom"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
)

func TestExtractEntityFieldDefinitions(t *testing.T) {
	testCases := []struct {
		name           string
		schema         kong.Schema
		expectedFields map[string]EntityField
	}{
		{
			name: "absent fields should have a nil value",
			schema: map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"foo": map[string]interface{}{
							"type":     "string",
							"required": true,
						},
					},
					map[string]interface{}{
						"bar": map[string]interface{}{
							"type":      "foreign",
							"required":  true,
							"reference": "service",
						},
					},
				},
			},
			expectedFields: map[string]EntityField{
				"foo": {
					Name:     "foo",
					Type:     EntityFieldTypeString,
					Required: true,
					Auto:     false,
					UUID:     false,
				},
				"bar": {
					Name:      "bar",
					Type:      EntityFieldTypeForeign,
					Required:  true,
					Reference: "service",
				},
			},
		},
		{
			name: "irrelevant fields should be safely ignored",
			schema: map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"protocol": map[string]interface{}{
							"type":     "string",
							"required": true,
							"default":  "http",
							"one_of":   []string{"http", "https"},
						},
						"port": map[string]interface{}{
							"type":     "integer",
							"required": true,
							"default":  80,
							"min":      1,
							"max":      65535,
						},
					},
				},
				"checks": "some_check",
			},
			expectedFields: map[string]EntityField{
				"protocol": {
					Name:     "protocol",
					Type:     EntityFieldTypeString,
					Required: true,
					Default:  "http",
				},
				"port": {
					Name:     "port",
					Type:     EntityFieldTypeInteger,
					Required: true,
					Default:  int(80),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fields := ExtractEntityFieldDefinitions(tc.schema).Fields
			for fieldName, expectedField := range tc.expectedFields {
				actualField, ok := fields[fieldName]
				require.Truef(t, ok, "field %s should exist", fieldName)
				require.Equalf(t, expectedField, actualField, "field %s should be same as expected", fieldName)
			}
		})
	}
}

func TestSortCustomEntities(t *testing.T) {
	tesCases := []struct {
		name                    string
		customEntityCollections map[string]*KongCustomEntityCollection
		sortedCollections       map[string]*KongCustomEntityCollection
	}{
		{
			name: "custom entities in multiple namespaces",
			customEntityCollections: map[string]*KongCustomEntityCollection{
				"foo": {
					Entities: []CustomEntity{
						{
							Object: custom.Object{
								"name": "e1",
								"key":  "value1",
							},
							K8sKongCustomEntity: &kongv1alpha1.KongCustomEntity{
								ObjectMeta: metav1.ObjectMeta{
									Name:      "aab",
									Namespace: "bbb",
								},
							},
						},
						{
							Object: custom.Object{
								"name": "e2",
								"key":  "value2",
							},
							K8sKongCustomEntity: &kongv1alpha1.KongCustomEntity{
								ObjectMeta: metav1.ObjectMeta{
									Name:      "abc",
									Namespace: "bbb",
								},
							},
						},
						{
							Object: custom.Object{
								"name": "e3",
								"key":  "value3",
							},
							K8sKongCustomEntity: &kongv1alpha1.KongCustomEntity{
								ObjectMeta: metav1.ObjectMeta{
									Name:      "abc",
									Namespace: "aaa",
								},
							},
						},
					},
				},
			},
			sortedCollections: map[string]*KongCustomEntityCollection{
				"foo": {
					Entities: []CustomEntity{
						{
							Object: custom.Object{
								"name": "e3",
								"key":  "value3",
							},
							K8sKongCustomEntity: &kongv1alpha1.KongCustomEntity{
								ObjectMeta: metav1.ObjectMeta{
									Name:      "abc",
									Namespace: "aaa",
								},
							},
						},
						{
							Object: custom.Object{
								"name": "e1",
								"key":  "value1",
							},
							K8sKongCustomEntity: &kongv1alpha1.KongCustomEntity{
								ObjectMeta: metav1.ObjectMeta{
									Name:      "aab",
									Namespace: "bbb",
								},
							},
						},
						{
							Object: custom.Object{
								"name": "e2",
								"key":  "value2",
							},
							K8sKongCustomEntity: &kongv1alpha1.KongCustomEntity{
								ObjectMeta: metav1.ObjectMeta{
									Name:      "abc",
									Namespace: "bbb",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tesCases {
		t.Run(tc.name, func(t *testing.T) {
			ks := &KongState{
				CustomEntities: tc.customEntityCollections,
			}
			ks.sortCustomEntities()
			require.Equal(t, tc.sortedCollections, ks.CustomEntities)
		})
	}
}

func TestFindustomEntityForeignFields(t *testing.T) {
	testCustomEntity := &kongv1alpha1.KongCustomEntity{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "fake-entity",
		},
		Spec: kongv1alpha1.KongCustomEntitySpec{
			EntityType:     "fake_entities",
			ControllerName: annotations.DefaultIngressClass,
			Fields: apiextensionsv1.JSON{
				Raw: []byte(`{"uri":"/api/me"}`),
			},
			ParentRef: &kongv1alpha1.ObjectReference{
				Group: kong.String(kongv1.GroupVersion.Group),
				Kind:  kong.String("KongPlugin"),
				Name:  "fake-plugin",
			},
		},
	}
	kongService1 := kong.Service{
		Name: kong.String("service1"),
		ID:   kong.String("service1"),
	}
	kongService2 := kong.Service{
		Name: kong.String("service2"),
		ID:   kong.String("service2"),
	}
	kongRoute1 := kong.Route{
		Name: kong.String("route1"),
		ID:   kong.String("route1"),
	}
	kongRoute2 := kong.Route{
		Name: kong.String("route2"),
		ID:   kong.String("route2"),
	}
	kongConsumer1 := kong.Consumer{
		Username: kong.String("consumer1"),
		ID:       kong.String("consumer1"),
	}
	kongConsumer2 := kong.Consumer{
		Username: kong.String("consumer2"),
		ID:       kong.String("consumer2"),
	}
	testCases := []struct {
		name                     string
		customEntity             *kongv1alpha1.KongCustomEntity
		schema                   EntitySchema
		pluginRelEntities        PluginRelatedEntitiesRefs
		foreignFieldCombinations [][]entityForeignFieldValue
	}{
		{
			name:         "attached to single entity: service",
			customEntity: testCustomEntity,
			schema: EntitySchema{
				Fields: map[string]EntityField{
					"foo":     {Name: "foo", Type: EntityFieldTypeString, Required: true},
					"service": {Name: "service", Type: EntityFieldTypeForeign, Reference: "services"},
				},
			},
			pluginRelEntities: PluginRelatedEntitiesRefs{
				RelatedEntities: map[string]RelatedEntitiesRef{
					"default:fake-plugin": {
						Services: []*Service{
							{
								Service: kongService1,
							},
							{
								Service: kongService2,
							},
						},
					},
				},
			},
			foreignFieldCombinations: [][]entityForeignFieldValue{
				{
					{fieldName: "service", foreignEntityType: kong.EntityTypeServices, foreignEntityID: "service1"},
				},
				{
					{fieldName: "service", foreignEntityType: kong.EntityTypeServices, foreignEntityID: "service2"},
				},
			},
		},
		{
			name:         "attached to routes and consumers",
			customEntity: testCustomEntity,
			schema: EntitySchema{
				Fields: map[string]EntityField{
					"foo":      {Name: "foo", Type: EntityFieldTypeString, Required: true},
					"route":    {Name: "route", Type: EntityFieldTypeForeign, Reference: "routes"},
					"consumer": {Name: "consumer", Type: EntityFieldTypeForeign, Reference: "consumers"},
				},
			},
			pluginRelEntities: PluginRelatedEntitiesRefs{
				RelatedEntities: map[string]RelatedEntitiesRef{
					"default:fake-plugin": {
						Routes: []*Route{
							{
								Route: kongRoute1,
							},
							{
								Route: kongRoute2,
							},
						},
						Consumers: []*Consumer{
							{
								Consumer: kongConsumer1,
							},
							{
								Consumer: kongConsumer2,
							},
						},
					},
				},
				RouteAttachedService: map[string]*Service{
					"route1": {Service: kongService1},
					"route2": {Service: kongService2},
				},
			},
			foreignFieldCombinations: [][]entityForeignFieldValue{
				{
					{fieldName: "consumer", foreignEntityType: kong.EntityTypeConsumers, foreignEntityID: "consumer1"},
					{fieldName: "route", foreignEntityType: kong.EntityTypeRoutes, foreignEntityID: "route1"},
				},
				{
					{fieldName: "consumer", foreignEntityType: kong.EntityTypeConsumers, foreignEntityID: "consumer1"},
					{fieldName: "route", foreignEntityType: kong.EntityTypeRoutes, foreignEntityID: "route2"},
				},
				{
					{fieldName: "consumer", foreignEntityType: kong.EntityTypeConsumers, foreignEntityID: "consumer2"},
					{fieldName: "route", foreignEntityType: kong.EntityTypeRoutes, foreignEntityID: "route1"},
				},
				{
					{fieldName: "consumer", foreignEntityType: kong.EntityTypeConsumers, foreignEntityID: "consumer2"},
					{fieldName: "route", foreignEntityType: kong.EntityTypeRoutes, foreignEntityID: "route2"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			combinations := findCustomEntityForeignFields(
				logr.Discard(),
				tc.customEntity,
				tc.schema,
				tc.pluginRelEntities,
				"",
			)
			for index, combination := range combinations {
				sort.SliceStable(combination, func(i, j int) bool {
					return combination[i].fieldName < combination[j].fieldName
				})
				combinations[index] = combination
			}
			for _, expectedCombination := range tc.foreignFieldCombinations {
				require.Contains(t, combinations, expectedCombination)
			}
		})
	}
}

func TestKongState_FillCustomEntities(t *testing.T) {
	customEntityTypeMeta := metav1.TypeMeta{
		APIVersion: kongv1alpha1.GroupVersion.Group + "/" + kongv1alpha1.GroupVersion.Version,
		Kind:       "KongCustomEntity",
	}
	kongService1 := kong.Service{
		Name: kong.String("service1"),
		ID:   kong.String("service1"),
	}
	kongService2 := kong.Service{
		Name: kong.String("service2"),
		ID:   kong.String("service2"),
	}
	ksService1 := Service{
		Service: kongService1,
		K8sServices: map[string]*corev1.Service{
			"default/service1": {
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "service1",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "degraphql-1",
					},
				},
			},
		},
	}
	ksService2 := Service{
		Service: kongService2,
		K8sServices: map[string]*corev1.Service{
			"default/service2": {
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "service2",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "degraphql-1",
					},
				},
			},
		},
	} // Service: service2

	testCases := []struct {
		name                        string
		initialState                *KongState
		customEntities              []*kongv1alpha1.KongCustomEntity
		plugins                     []*kongv1.KongPlugin
		schemas                     map[string]kong.Schema
		expectedCustomEntities      map[string][]custom.Object
		expectedTranslationFailures map[k8stypes.NamespacedName]string
	}{
		{
			name:         "single custom entity",
			initialState: &KongState{},
			customEntities: []*kongv1alpha1.KongCustomEntity{
				{
					TypeMeta: customEntityTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "session-foo",
					},
					Spec: kongv1alpha1.KongCustomEntitySpec{
						EntityType:     "sessions",
						ControllerName: annotations.DefaultIngressClass,
						Fields: apiextensionsv1.JSON{
							Raw: []byte(`{"name":"session1"}`),
						},
					},
				},
			},
			schemas: map[string]kong.Schema{
				"sessions": {
					"fields": []interface{}{
						map[string]interface{}{
							"name": map[string]interface{}{
								"type":     "string",
								"required": true,
							},
						},
					},
				},
			},
			expectedCustomEntities: map[string][]custom.Object{
				"sessions": {
					{
						"name": "session1",
					},
				},
			},
		},
		{
			name:         "custom entity with unknown type",
			initialState: &KongState{},
			customEntities: []*kongv1alpha1.KongCustomEntity{
				{
					TypeMeta: customEntityTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "session-foo",
					},
					Spec: kongv1alpha1.KongCustomEntitySpec{
						EntityType:     "sessions",
						ControllerName: annotations.DefaultIngressClass,
						Fields: apiextensionsv1.JSON{
							Raw: []byte(`{"name":"session1"}`),
						},
					},
				},
			},
			expectedTranslationFailures: map[k8stypes.NamespacedName]string{
				{
					Namespace: "default",
					Name:      "session-foo",
				}: "failed to fetch entity schema for entity type sessions: schema not found",
			},
		},
		{
			name:         "multiple custom entities with same type",
			initialState: &KongState{},
			customEntities: []*kongv1alpha1.KongCustomEntity{
				{
					TypeMeta: customEntityTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "session-foo",
					},
					Spec: kongv1alpha1.KongCustomEntitySpec{
						EntityType:     "sessions",
						ControllerName: annotations.DefaultIngressClass,
						Fields: apiextensionsv1.JSON{
							Raw: []byte(`{"name":"session-foo"}`),
						},
					},
				},
				{
					TypeMeta: customEntityTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "session-bar",
					},
					Spec: kongv1alpha1.KongCustomEntitySpec{
						EntityType:     "sessions",
						ControllerName: annotations.DefaultIngressClass,
						Fields: apiextensionsv1.JSON{
							Raw: []byte(`{"name":"session-bar"}`),
						},
					},
				},
				{
					TypeMeta: customEntityTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default-1",
						Name:      "session-foo",
					},
					Spec: kongv1alpha1.KongCustomEntitySpec{
						EntityType:     "sessions",
						ControllerName: annotations.DefaultIngressClass,
						Fields: apiextensionsv1.JSON{
							Raw: []byte(`{"name":"session-foo-1"}`),
						},
					},
				},
			},
			schemas: map[string]kong.Schema{
				"sessions": {
					"fields": []interface{}{
						map[string]interface{}{
							"name": map[string]interface{}{
								"type":     "string",
								"required": true,
							},
						},
					},
				},
			},
			expectedCustomEntities: map[string][]custom.Object{
				// Should be sorted by original KCE namespace/name.
				"sessions": {
					{
						// from default/bar
						"name": "session-bar",
					},
					{
						// from default/foo
						"name": "session-foo",
					},
					{
						// from default-1/foo
						"name": "session-foo-1",
					},
				},
			},
		},
		{
			name: "custom entities with reference to other entities (services)",
			initialState: &KongState{
				Services: []Service{ksService1}, // Services
			},
			customEntities: []*kongv1alpha1.KongCustomEntity{
				{
					TypeMeta: customEntityTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "degraphql-1",
					},
					Spec: kongv1alpha1.KongCustomEntitySpec{
						EntityType:     "degraphql_routes",
						ControllerName: annotations.DefaultIngressClass,
						Fields: apiextensionsv1.JSON{
							Raw: []byte(`{"uri":"/api/me"}`),
						},
						ParentRef: &kongv1alpha1.ObjectReference{
							Group: kong.String(kongv1.GroupVersion.Group),
							Kind:  kong.String("KongPlugin"),
							Name:  "degraphql-1",
						},
					},
				},
			},
			plugins: []*kongv1.KongPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "degraphql-1",
					},
					PluginName: "degraphql",
				},
			},
			schemas: map[string]kong.Schema{
				"degraphql_routes": {
					"fields": []interface{}{
						map[string]interface{}{
							"uri": map[string]interface{}{
								"type":     "string",
								"required": true,
							},
						},
						map[string]interface{}{
							"service": map[string]interface{}{
								"type":      "foreign",
								"reference": "services",
							},
						},
					},
				},
			},
			expectedCustomEntities: map[string][]custom.Object{
				"degraphql_routes": {
					{
						"uri": "/api/me",
						"service": map[string]interface{}{
							// ID of Kong service "service1" in workspace "".
							"id": "service1",
						},
					},
				},
			},
		},
		{
			name: "custom entity attached to multiple services via plugin",
			initialState: &KongState{
				Services: []Service{
					ksService1,
					ksService2,
				},
			},
			customEntities: []*kongv1alpha1.KongCustomEntity{
				{
					TypeMeta: customEntityTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "degraphql-1",
					},
					Spec: kongv1alpha1.KongCustomEntitySpec{
						EntityType:     "degraphql_routes",
						ControllerName: annotations.DefaultIngressClass,
						Fields: apiextensionsv1.JSON{
							Raw: []byte(`{"uri":"/api/me"}`),
						},
						ParentRef: &kongv1alpha1.ObjectReference{
							Group: kong.String(kongv1.GroupVersion.Group),
							Kind:  kong.String("KongPlugin"),
							Name:  "degraphql-1",
						},
					},
				},
			},
			plugins: []*kongv1.KongPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "degraphql-1",
					},
					PluginName: "degraphql",
				},
			},
			schemas: map[string]kong.Schema{
				"degraphql_routes": {
					"fields": []interface{}{
						map[string]interface{}{
							"uri": map[string]interface{}{
								"type":     "string",
								"required": true,
							},
						},
						map[string]interface{}{
							"service": map[string]interface{}{
								"type":      "foreign",
								"reference": "services",
							},
						},
					},
				},
			},
			expectedCustomEntities: map[string][]custom.Object{
				"degraphql_routes": {
					{
						"uri": "/api/me",
						"service": map[string]interface{}{
							// ID of Kong service "service1" in workspace "".
							"id": "service1",
						},
					},
					{
						"uri": "/api/me",
						"service": map[string]interface{}{
							// ID of Kong service "service2" in workspace "".
							"id": "service2",
						},
					},
				},
			},
		},
		{
			name: "custom entity attached to route",
			initialState: &KongState{
				Services: []Service{
					{
						Service: kongService1,
						Routes: []Route{
							{
								Route: kong.Route{
									Name: kong.String("route1"),
									ID:   kong.String("route1"),
								},
								Ingress: util.K8sObjectInfo{
									Name:      "ingerss1",
									Namespace: "default",
									Annotations: map[string]string{
										annotations.AnnotationPrefix + annotations.PluginsKey: "degraphql-1",
									},
								},
							},
						},
					},
				},
			},
			customEntities: []*kongv1alpha1.KongCustomEntity{
				{
					TypeMeta: customEntityTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "degraphql-1",
					},
					Spec: kongv1alpha1.KongCustomEntitySpec{
						EntityType:     "degraphql_routes",
						ControllerName: annotations.DefaultIngressClass,
						Fields: apiextensionsv1.JSON{
							Raw: []byte(`{"uri":"/api/me"}`),
						},
						ParentRef: &kongv1alpha1.ObjectReference{
							Group: kong.String(kongv1.GroupVersion.Group),
							Kind:  kong.String("KongPlugin"),
							Name:  "degraphql-1",
						},
					},
				},
			},
			plugins: []*kongv1.KongPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "degraphql-1",
					},
					PluginName: "degraphql",
				},
			},
			schemas: map[string]kong.Schema{
				"degraphql_routes": {
					"fields": []interface{}{
						map[string]interface{}{
							"uri": map[string]interface{}{
								"type":     "string",
								"required": true,
							},
						},
						map[string]interface{}{
							"route": map[string]interface{}{
								"type":      "foreign",
								"reference": "routes",
							},
						},
					},
				},
			},
			expectedCustomEntities: map[string][]custom.Object{
				"degraphql_routes": {
					{
						"uri": "/api/me",
						"route": map[string]interface{}{
							// ID of Kong route "route1".
							"id": "route1",
						},
					},
				},
			},
		},
		{
			name: "custom entity attached to two services and one consumer",
			initialState: &KongState{
				Services: []Service{
					ksService1,
					ksService2,
				},
				Consumers: []Consumer{
					{
						Consumer: kong.Consumer{
							ID:       kong.String("consumer1"),
							Username: kong.String("consumer1"),
						},
						K8sKongConsumer: kongv1.KongConsumer{
							ObjectMeta: metav1.ObjectMeta{
								Namespace: "default",
								Name:      "consumer1",
								Annotations: map[string]string{
									annotations.AnnotationPrefix + annotations.PluginsKey: "degraphql-1",
								},
							},
						},
					},
				},
			},
			customEntities: []*kongv1alpha1.KongCustomEntity{
				{
					TypeMeta: customEntityTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "fake-entity-1",
					},
					Spec: kongv1alpha1.KongCustomEntitySpec{
						EntityType:     "fake_entities",
						ControllerName: annotations.DefaultIngressClass,
						Fields: apiextensionsv1.JSON{
							Raw: []byte(`{"foo":"bar"}`),
						},
						ParentRef: &kongv1alpha1.ObjectReference{
							Group: kong.String(kongv1.GroupVersion.Group),
							Kind:  kong.String("KongPlugin"),
							Name:  "degraphql-1",
						},
					},
				},
			},
			plugins: []*kongv1.KongPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "degraphql-1",
					},
					PluginName: "degraphql",
				},
			},
			schemas: map[string]kong.Schema{
				"fake_entities": {
					"fields": []interface{}{
						map[string]interface{}{
							"foo": map[string]interface{}{
								"type":     "string",
								"required": true,
							},
						},
						map[string]interface{}{
							"service": map[string]interface{}{
								"type":      "foreign",
								"reference": "services",
							},
						},
						map[string]interface{}{
							"consumer": map[string]interface{}{
								"type":      "foreign",
								"reference": "consumers",
							},
						},
					},
				},
			},
			expectedCustomEntities: map[string][]custom.Object{
				"fake_entities": {
					{
						"foo": "bar",
						"service": map[string]interface{}{
							// ID of Kong service "service1".
							"id": "service1",
						},
						"consumer": map[string]interface{}{
							// ID of Kong consumer "consumer1".
							"id": "consumer1",
						},
					},
					{
						"foo": "bar",
						"service": map[string]interface{}{
							// ID of Kong service "service2".
							"id": "service2",
						},
						"consumer": map[string]interface{}{
							// ID of Kong consumer "consumer1".
							"id": "consumer1",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := store.NewFakeStore(store.FakeObjects{
				KongCustomEntities: tc.customEntities,
				KongPlugins:        tc.plugins,
			})
			require.NoError(t, err)
			failuresCollector := failures.NewResourceFailuresCollector(logr.Discard())

			ks := tc.initialState
			ks.FillCustomEntities(
				logr.Discard(), s,
				failuresCollector,
				&fakeSchemaGetter{schemas: tc.schemas}, "",
			)
			for entityType, expectedObjectList := range tc.expectedCustomEntities {
				require.NotNil(t, ks.CustomEntities[entityType])
				objectList := lo.Map(ks.CustomEntities[entityType].Entities, func(e CustomEntity, _ int) custom.Object {
					return e.Object
				})
				require.Equal(t, expectedObjectList, objectList)
			}

			translationFailures := failuresCollector.PopResourceFailures()
			for nsName, message := range tc.expectedTranslationFailures {
				hasError := lo.ContainsBy(translationFailures, func(f failures.ResourceFailure) bool {
					fmt.Println(f.Message())
					return f.Message() == message && lo.ContainsBy(f.CausingObjects(), func(o client.Object) bool {
						return o.GetNamespace() == nsName.Namespace && o.GetName() == nsName.Name
					})
				})
				require.Truef(t, hasError, "translation error for KongCustomEntity %s not found", nsName)
			}
		})
	}
}
