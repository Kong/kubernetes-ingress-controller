package adminapi

import (
	"context"
	"errors"
	"fmt"
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
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	cfgtypes "github.com/kong/kubernetes-ingress-controller/v3/internal/manager/config/types"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
)

func TestDiscoverer_AddressesFromEndpointSlice(t *testing.T) {
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
		name        string
		endpoints   discoveryv1.EndpointSlice
		want        sets.Set[DiscoveredAdminAPI]
		portNames   sets.Set[string]
		dnsStrategy cfgtypes.DNSStrategy
		expectedErr error
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
				Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
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
			dnsStrategy: cfgtypes.NamespaceScopedPodDNSStrategy,
		},
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
				Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
			},
			portNames: sets.New("admin"),
			want: sets.New(
				DiscoveredAdminAPI{
					Address: "https://10.0.0.1:8444",
					PodRef: k8stypes.NamespacedName{
						Name: "pod-1", Namespace: namespaceName,
					},
				},
			),
			dnsStrategy: cfgtypes.IPDNSStrategy,
		},
		{
			name: "basic IPDNSStrategy IPv6",
			endpoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv6,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"fe80::cae2:65ff:fe7b:2852", "fe80::cae2:65ff:fe7b:2853"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
						TargetRef: testPodReference(namespaceName, "pod-1"),
					},
				},
				Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
			},
			portNames: sets.New("admin"),
			want: sets.New(
				DiscoveredAdminAPI{
					Address: "https://[fe80::cae2:65ff:fe7b:2852]:8444",
					PodRef: k8stypes.NamespacedName{
						Name: "pod-1", Namespace: namespaceName,
					},
				},
			),
			dnsStrategy: cfgtypes.IPDNSStrategy,
		},
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
				Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
			},
			portNames: sets.New("admin"),
			want: sets.New(
				DiscoveredAdminAPI{
					Address: "https://10.0.0.1:8444",
					PodRef: k8stypes.NamespacedName{
						Name: "pod-1", Namespace: namespaceName,
					},
				},
			),
			dnsStrategy: cfgtypes.ServiceScopedPodDNSStrategy,
			expectedErr: errors.New("service name is empty for an endpoint with TargetRef ns/pod-1"),
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
				Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
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
			dnsStrategy: cfgtypes.ServiceScopedPodDNSStrategy,
		},
		{
			name: "not ready endpoints are returned",
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
				Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
			},
			portNames: sets.New("admin"),
			want: sets.New[DiscoveredAdminAPI](
				DiscoveredAdminAPI{
					Address: "https://10.0.0.1:8444",
					PodRef: k8stypes.NamespacedName{
						Name: "pod-1", Namespace: namespaceName,
					},
				}),
			dnsStrategy: cfgtypes.IPDNSStrategy,
		},
		{
			name: "ready and terminating endpoints are not returned",
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
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(true),
						},
						TargetRef: testPodReference(namespaceName, "pod-1"),
					},
				},
				Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
			},
			portNames:   sets.New("admin"),
			want:        sets.New[DiscoveredAdminAPI](),
			dnsStrategy: cfgtypes.IPDNSStrategy,
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
							Terminating: lo.ToPtr(true),
						},
						TargetRef: testPodReference(namespaceName, "pod-3"),
					},
				},
				Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
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
			dnsStrategy: cfgtypes.NamespaceScopedPodDNSStrategy,
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
							Terminating: lo.ToPtr(true),
						},
						TargetRef: testPodReference(namespaceName, "pod-3"),
					},
				},
				Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
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
			dnsStrategy: cfgtypes.ServiceScopedPodDNSStrategy,
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
				Ports: builder.NewEndpointPort(8444).WithName("non-admin-port-name").IntoSlice(),
			},
			want:        sets.New[DiscoveredAdminAPI](),
			portNames:   sets.New("admin"),
			dnsStrategy: cfgtypes.IPDNSStrategy,
		},
		{
			name: "ports without names are not taken into account",
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
				Ports: builder.NewEndpointPort(8444).IntoSlice(),
			},
			portNames:   sets.New("admin"),
			want:        sets.New[DiscoveredAdminAPI](),
			dnsStrategy: cfgtypes.IPDNSStrategy,
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
					builder.NewEndpointPort(8443).WithName("admin-tls").Build(),
					builder.NewEndpointPort(8444).WithName("admin").Build(),
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
			dnsStrategy: cfgtypes.NamespaceScopedPodDNSStrategy,
		},
		{
			name: "endpoints with no target ref return error for service scopec dns strategy",
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
				Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
			},
			portNames:   sets.New("admin"),
			want:        sets.New[DiscoveredAdminAPI](),
			dnsStrategy: cfgtypes.IPDNSStrategy,
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
				Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
			},
			portNames:   sets.New("admin"),
			want:        sets.New[DiscoveredAdminAPI](),
			dnsStrategy: cfgtypes.IPDNSStrategy,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("dnsstrategy_%s/%s", tt.dnsStrategy, tt.name), func(t *testing.T) {
			discoverer, err := NewDiscoverer(tt.portNames, tt.dnsStrategy)
			require.NoError(t, err)

			got, err := discoverer.AdminAPIsFromEndpointSlice(tt.endpoints)
			if tt.expectedErr != nil {
				require.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDiscoverer_GetAdminAPIsForService(t *testing.T) {
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
		name        string
		service     k8stypes.NamespacedName
		objects     []client.ObjectList
		dnsStrategy cfgtypes.DNSStrategy
		want        sets.Set[DiscoveredAdminAPI]
		wantErr     bool
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
							Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
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
							Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
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
										Terminating: lo.ToPtr(true),
									},
									TargetRef: testPodReference(namespaceName, "pod-3"),
								},
							},
							Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
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
			dnsStrategy: cfgtypes.NamespaceScopedPodDNSStrategy,
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
							Ports: builder.NewEndpointPort(8444).WithName("non-admin-port").IntoSlice(),
						},
					},
				},
			},
			want:        sets.New[DiscoveredAdminAPI](),
			dnsStrategy: cfgtypes.IPDNSStrategy,
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
							Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
						},
					},
				},
			},
			want:        sets.New[DiscoveredAdminAPI](),
			dnsStrategy: cfgtypes.IPDNSStrategy,
		},
		{
			name: "terminating Endpoints are not matched",
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
										Ready:       lo.ToPtr(false),
										Terminating: lo.ToPtr(true),
									},
									TargetRef: testPodReference(namespaceName, "pod-1"),
								},
							},
							Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
						},
					},
				},
			},
			want:        sets.New[DiscoveredAdminAPI](),
			dnsStrategy: cfgtypes.IPDNSStrategy,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("dnsstrategy_%s/%s", tt.dnsStrategy, tt.name), func(t *testing.T) {
			require.NoError(t, tt.dnsStrategy.Validate())

			fakeClient := fake.NewClientBuilder().
				WithLists(tt.objects...).
				Build()

			portNames := sets.New("admin")
			discoverer, err := NewDiscoverer(portNames, tt.dnsStrategy)
			require.NoError(t, err)

			got, err := discoverer.GetAdminAPIsForService(context.Background(), fakeClient, tt.service)
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
