package adminapi

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestAddressesFromEndpointSlice(t *testing.T) {
	const (
		namespaceName = "ns"
		serviceName   = "kong-admin"
	)

	var (
		endpointsSliceObjectMeta = metav1.ObjectMeta{
			Name:      uuid.NewString(),
			Namespace: namespaceName,
		}
		endpointsSliceWithOwnerReferenceObjectMeta = metav1.ObjectMeta{
			Name:      uuid.NewString(),
			Namespace: namespaceName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Name:       serviceName,
					Kind:       "Service",
				},
			},
		}
	)

	tests := []struct {
		name      string
		endpoints discoveryv1.EndpointSlice
		want      sets.Set[DiscoveredAdminAPI]
		portNames sets.Set[string]
	}{
		{
			name: "basic",
			endpoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-1"),
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Name: lo.ToPtr("admin"),
						Port: lo.ToPtr(int32(8444)),
					},
				},
			},
			portNames: sets.New("admin"),
			want: sets.New(
				DiscoveredAdminAPI{
					Address: "https://10-0-0-1.ns.pod:8444",
					PodRef: k8stypes.NamespacedName{
						Name: "pod-1", Namespace: namespaceName,
					},
				},
			),
		},
		{
			name: "basic with owner reference",
			endpoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceWithOwnerReferenceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-1"),
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Name: lo.ToPtr("admin"),
						Port: lo.ToPtr(int32(8444)),
					},
				},
			},
			portNames: sets.New("admin"),
			want: sets.New(
				DiscoveredAdminAPI{
					Address: "https://10-0-0-1.kong-admin.ns.svc:8444",
					PodRef: k8stypes.NamespacedName{
						Name: "pod-1", Namespace: namespaceName,
					},
				},
			),
		},
		{
			name: "not ready endpoints are not returned",
			endpoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(false),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-1"),
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Name: lo.ToPtr("admin"),
						Port: lo.ToPtr(int32(8444)),
					},
				},
			},
			portNames: sets.New("admin"),
			want:      sets.New[DiscoveredAdminAPI](),
		},
		{
			name: "not ready and terminating endpoints are not returned",
			endpoints: discoveryv1.EndpointSlice{
				ObjectMeta: metav1.ObjectMeta{
					Name:      uuid.NewString(),
					Namespace: namespaceName,
				},
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(false),
							Terminating: lo.ToPtr(true),
						},
						TargetRef: testPodReference(namespaceName, "pod-1"),
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Name: lo.ToPtr("admin"),
						Port: lo.ToPtr(int32(8444)),
					},
				},
			},
			portNames: sets.New("admin"),
			want:      sets.New[DiscoveredAdminAPI](),
		},
		{
			name: "multiple endpoints are concatenated properly",
			endpoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-1"),
					},
					{
						Addresses: []string{"10.0.1.1", "10.0.1.2"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-2"),
					},
					{
						Addresses: []string{"10.0.2.1"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(false),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-3"),
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Name: lo.ToPtr("admin"),
						Port: lo.ToPtr(int32(8444)),
					},
				},
			},
			portNames: sets.New("admin"),
			want: sets.New(
				DiscoveredAdminAPI{
					Address: "https://10-0-0-1.ns.pod:8444",
					PodRef: k8stypes.NamespacedName{
						Namespace: namespaceName,
						Name:      "pod-1",
					},
				},
				DiscoveredAdminAPI{
					Address: "https://10-0-1-1.ns.pod:8444",
					PodRef: k8stypes.NamespacedName{
						Namespace: namespaceName,
						Name:      "pod-2",
					},
				},
			),
		},
		{
			name: "multiple endpoints with owner reference are concatenated properly",
			endpoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceWithOwnerReferenceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-1"),
					},
					{
						Addresses: []string{"10.0.1.1", "10.0.1.2"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-2"),
					},
					{
						Addresses: []string{"10.0.2.1"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(false),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-3"),
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Name: lo.ToPtr("admin"),
						Port: lo.ToPtr(int32(8444)),
					},
				},
			},
			portNames: sets.New("admin"),
			want: sets.New(
				DiscoveredAdminAPI{
					Address: "https://10-0-0-1.kong-admin.ns.svc:8444",
					PodRef: k8stypes.NamespacedName{
						Namespace: namespaceName,
						Name:      "pod-1",
					},
				},
				DiscoveredAdminAPI{
					Address: "https://10-0-1-1.kong-admin.ns.svc:8444",
					PodRef: k8stypes.NamespacedName{
						Namespace: namespaceName,
						Name:      "pod-2",
					},
				},
			),
		},
		{
			name: "ports not called 'admin' are not added",
			endpoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-1"),
					},
					{
						Addresses: []string{"10.0.1.1", "10.0.1.2"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-2"),
					},
					{
						Addresses: []string{"10.0.2.1"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(false),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-3"),
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Name: lo.ToPtr("non-admin-port-name"),
						Port: lo.ToPtr(int32(8444)),
					},
				},
			},
			want: sets.New[DiscoveredAdminAPI](),
		},
		{
			name: "ports without names are not taken into account ",
			endpoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-1"),
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Port: lo.ToPtr(int32(8444)),
					},
				},
			},
			portNames: sets.New("admin"),
			want:      sets.New[DiscoveredAdminAPI](),
		},
		{
			name: "multiple ports names",
			endpoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-1"),
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Name: lo.ToPtr("admin-tls"),
						Port: lo.ToPtr(int32(8443)),
					},
					{
						Name: lo.ToPtr("admin"),
						Port: lo.ToPtr(int32(8444)),
					},
				},
			},
			portNames: sets.New("admin", "admin-tls"),
			want: sets.New(
				DiscoveredAdminAPI{
					Address: "https://10-0-0-1.ns.pod:8443",
					PodRef: k8stypes.NamespacedName{
						Namespace: namespaceName,
						Name:      "pod-1",
					},
				},
				DiscoveredAdminAPI{
					Address: "https://10-0-0-1.ns.pod:8444",
					PodRef: k8stypes.NamespacedName{
						Namespace: namespaceName,
						Name:      "pod-1",
					},
				},
			),
		},
		{
			name: "endpoints with no target ref are ignored",
			endpoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: nil,
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Name: lo.ToPtr("admin"),
						Port: lo.ToPtr(int32(8444)),
					},
				},
			},
			portNames: sets.New("admin"),
			want:      sets.New[DiscoveredAdminAPI](),
		},
		{
			name: "endpoints with target ref other than Pod are ignored",
			endpoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: &corev1.ObjectReference{Kind: "Node", Namespace: namespaceName, Name: "node-1"},
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Name: lo.ToPtr("admin"),
						Port: lo.ToPtr(int32(8444)),
					},
				},
			},
			portNames: sets.New("admin"),
			want:      sets.New[DiscoveredAdminAPI](),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, AdminAPIsFromEndpointSlice(tt.endpoints, tt.portNames))
		})
	}
}

func TestGetAdminAPIsForService(t *testing.T) {
	const namespaceName = "ns"

	var (
		serviceName                   = uuid.NewString()
		matchingServiceObjectMetaFunc = func() metav1.ObjectMeta {
			return metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: namespaceName,
				Labels: map[string]string{
					"kubernetes.io/service-name": serviceName,
				},
			}
		}
	)

	tests := []struct {
		name    string
		service k8stypes.NamespacedName
		objects []client.ObjectList
		want    sets.Set[DiscoveredAdminAPI]
		wantErr bool
	}{
		{
			name: "basic",
			service: k8stypes.NamespacedName{
				Namespace: namespaceName,
				Name:      serviceName,
			},
			objects: []client.ObjectList{
				&discoveryv1.EndpointSliceList{
					Items: []discoveryv1.EndpointSlice{
						{
							ObjectMeta:  matchingServiceObjectMetaFunc(),
							AddressType: discoveryv1.AddressTypeIPv4,
							Endpoints: []discoveryv1.Endpoint{
								{
									Addresses: []string{"10.0.0.1", "10.0.0.2"},
									Conditions: discoveryv1.EndpointConditions{
										Ready:       lo.ToPtr(true),
										Terminating: lo.ToPtr(false),
									},
									TargetRef: testPodReference(namespaceName, "pod-1"),
								},
							},
							Ports: []discoveryv1.EndpointPort{
								{
									Name: lo.ToPtr("admin"),
									Port: lo.ToPtr(int32(8444)),
								},
							},
						},
					},
				},
				&discoveryv1.EndpointSliceList{
					Items: []discoveryv1.EndpointSlice{
						{
							ObjectMeta:  matchingServiceObjectMetaFunc(),
							AddressType: discoveryv1.AddressTypeIPv4,
							Endpoints: []discoveryv1.Endpoint{
								{
									Addresses: []string{"9.0.0.1"},
									Conditions: discoveryv1.EndpointConditions{
										Ready:       lo.ToPtr(true),
										Terminating: lo.ToPtr(false),
									},
									TargetRef: testPodReference(namespaceName, "pod-2"),
								},
							},
							Ports: []discoveryv1.EndpointPort{
								{
									Name: lo.ToPtr("admin"),
									Port: lo.ToPtr(int32(8444)),
								},
							},
						},
					},
				},
				&discoveryv1.EndpointSliceList{
					Items: []discoveryv1.EndpointSlice{
						{
							ObjectMeta:  matchingServiceObjectMetaFunc(),
							AddressType: discoveryv1.AddressTypeIPv4,
							Endpoints: []discoveryv1.Endpoint{
								{
									Addresses: []string{"8.0.0.1"},
									Conditions: discoveryv1.EndpointConditions{
										Ready:       lo.ToPtr(false),
										Terminating: lo.ToPtr(false),
									},
									TargetRef: testPodReference(namespaceName, "pod-3"),
								},
							},
							Ports: []discoveryv1.EndpointPort{
								{
									Name: lo.ToPtr("admin"),
									Port: lo.ToPtr(int32(8444)),
								},
							},
						},
					},
				},
			},
			want: sets.New(
				DiscoveredAdminAPI{
					Address: "https://10-0-0-1.ns.pod:8444", PodRef: k8stypes.NamespacedName{
						Namespace: namespaceName,
						Name:      "pod-1",
					},
				},
				DiscoveredAdminAPI{
					Address: "https://9-0-0-1.ns.pod:8444", PodRef: k8stypes.NamespacedName{
						Namespace: namespaceName,
						Name:      "pod-2",
					},
				},
			),
		},
		{
			name: "ports not matching the specified port names are not taken into account",
			service: k8stypes.NamespacedName{
				Namespace: namespaceName,
				Name:      serviceName,
			},
			objects: []client.ObjectList{
				&discoveryv1.EndpointSliceList{
					Items: []discoveryv1.EndpointSlice{
						{
							ObjectMeta:  matchingServiceObjectMetaFunc(),
							AddressType: discoveryv1.AddressTypeIPv4,
							Endpoints: []discoveryv1.Endpoint{
								{
									Addresses: []string{"7.0.0.1"},
									Conditions: discoveryv1.EndpointConditions{
										Ready:       lo.ToPtr(true),
										Terminating: lo.ToPtr(false),
									},
									TargetRef: testPodReference(namespaceName, "pod-1"),
								},
							},
							Ports: []discoveryv1.EndpointPort{
								{
									Name: lo.ToPtr("non-admin-port"),
									Port: lo.ToPtr(int32(8444)),
								},
							},
						},
					},
				},
			},
			want: sets.New[DiscoveredAdminAPI](),
		},
		{
			name: "Endpoints without a TargetRef are not matched",
			service: k8stypes.NamespacedName{
				Namespace: namespaceName,
				Name:      serviceName,
			},
			objects: []client.ObjectList{
				&discoveryv1.EndpointSliceList{
					Items: []discoveryv1.EndpointSlice{
						{
							ObjectMeta:  matchingServiceObjectMetaFunc(),
							AddressType: discoveryv1.AddressTypeIPv4,
							Endpoints: []discoveryv1.Endpoint{
								{
									Addresses: []string{"7.0.0.1"},
									Conditions: discoveryv1.EndpointConditions{
										Ready:       lo.ToPtr(true),
										Terminating: lo.ToPtr(false),
									},
								},
							},
							Ports: []discoveryv1.EndpointPort{
								{
									Name: lo.ToPtr("admin"),
									Port: lo.ToPtr(int32(8444)),
								},
							},
						},
					},
				},
			},
			want: sets.New[DiscoveredAdminAPI](),
		},
		{
			name: "not Ready Endpoints are not matched",
			service: k8stypes.NamespacedName{
				Namespace: namespaceName,
				Name:      serviceName,
			},
			objects: []client.ObjectList{
				&discoveryv1.EndpointSliceList{
					Items: []discoveryv1.EndpointSlice{
						{
							ObjectMeta:  matchingServiceObjectMetaFunc(),
							AddressType: discoveryv1.AddressTypeIPv4,
							Endpoints: []discoveryv1.Endpoint{
								{
									Addresses: []string{"7.0.0.1"},
									Conditions: discoveryv1.EndpointConditions{
										Ready: lo.ToPtr(false),
									},
									TargetRef: testPodReference(namespaceName, "pod-1"),
								},
							},
							Ports: []discoveryv1.EndpointPort{
								{
									Name: lo.ToPtr("admin"),
									Port: lo.ToPtr(int32(8444)),
								},
							},
						},
					},
				},
			},
			want: sets.New[DiscoveredAdminAPI](),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := fakeclient.NewClientBuilder().
				WithLists(tt.objects...).
				Build()

			portNames := sets.New("admin")
			got, err := GetAdminAPIsForService(context.Background(), fakeClient, tt.service, portNames)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func testPodReference(namespace, name string) *corev1.ObjectReference { //nolint:unparam
	return &corev1.ObjectReference{
		Kind:      "Pod",
		Namespace: namespace,
		Name:      name,
	}
}
