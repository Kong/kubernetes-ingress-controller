package adminapi

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestAddressesFromEndpointSlice(t *testing.T) {
	endpointsSliceObjectMeta := metav1.ObjectMeta{
		Name:      uuid.NewString(),
		Namespace: "ns",
	}

	tests := []struct {
		name      string
		enspoints discoveryv1.EndpointSlice
		want      sets.Set[string]
	}{
		{
			enspoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2"},
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
			name: "basic",
			want: sets.New("https://10.0.0.1:8444", "https://10.0.0.2:8444"),
		},
		{
			enspoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(false),
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
			name: "not ready endpoints are not returned",
			want: sets.New[string](),
		},
		{
			enspoints: discoveryv1.EndpointSlice{
				ObjectMeta: metav1.ObjectMeta{
					Name:      uuid.NewString(),
					Namespace: "ns",
				},
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(false),
							Terminating: lo.ToPtr(true),
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
			name: "not ready and terminating endpoints are not returned",
			want: sets.New[string](),
		},
		{
			enspoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
					},
					{
						Addresses: []string{"10.0.1.1", "10.0.1.2"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
					},
					{
						Addresses: []string{"10.0.2.1"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(false),
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
			name: "multiple enpoints are concatenated properly",
			want: sets.New("https://10.0.0.1:8444", "https://10.0.0.2:8444", "https://10.0.0.3:8444", "https://10.0.1.1:8444", "https://10.0.1.2:8444"),
		},
		{
			enspoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
					},
					{
						Addresses: []string{"10.0.1.1", "10.0.1.2"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
					},
					{
						Addresses: []string{"10.0.2.1"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(false),
							Terminating: lo.ToPtr(false),
						},
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Name: lo.ToPtr("non-admin-port-name"),
						Port: lo.ToPtr(int32(8444)),
					},
				},
			},
			name: "ports not called 'admin' are not added",
			want: sets.New[string](),
		},
		{
			enspoints: discoveryv1.EndpointSlice{
				ObjectMeta:  endpointsSliceObjectMeta,
				AddressType: discoveryv1.AddressTypeIPv4,
				Endpoints: []discoveryv1.Endpoint{
					{
						Addresses: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
						Conditions: discoveryv1.EndpointConditions{
							Ready:       lo.ToPtr(true),
							Terminating: lo.ToPtr(false),
						},
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Port: lo.ToPtr(int32(8444)),
					},
				},
			},
			name: "ports without names are not taken into account ",
			want: sets.New[string](),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, AddressesFromEndpointSlice(tt.enspoints))
		})
	}
}

func TestGetURLsForService(t *testing.T) {
	var (
		serviceName                   = uuid.NewString()
		matchingServiceObjectMetaFunc = func() metav1.ObjectMeta {
			return metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: "ns",
				Labels: map[string]string{
					"kubernetes.io/service-name": serviceName,
				},
			}
		}
	)

	tests := []struct {
		name    string
		service types.NamespacedName
		objects []client.ObjectList
		want    sets.Set[string]
		wantErr bool
	}{
		{
			service: types.NamespacedName{
				Namespace: "ns",
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
			want: sets.New("https://10.0.0.1:8444", "https://10.0.0.2:8444", "https://9.0.0.1:8444"),
			name: "basic",
		},
		{
			service: types.NamespacedName{
				Namespace: "ns",
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
									Name: lo.ToPtr("non-admin-port"),
									Port: lo.ToPtr(int32(8444)),
								},
							},
						},
					},
				},
			},
			want: sets.New[string](),
			name: "port not called 'admin' are not taken into account",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := fakeclient.NewClientBuilder().
				WithLists(tt.objects...).
				Build()

			got, err := GetURLsForService(context.Background(), fakeClient, tt.service)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}
