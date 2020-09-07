package utils

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
)

type fakeDiscoveryClient struct {
	discovery.ServerResourcesInterface

	results map[string]metav1.APIResourceList
	err     error
}

// ServerResourcesForGroupVersion returns the supported resources for a group and version.
func (fdc *fakeDiscoveryClient) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	resp := fdc.results[groupVersion]
	return &resp, fdc.err
}

func TestServerHasGVK(t *testing.T) {
	okClient := fakeDiscoveryClient{
		results: map[string]metav1.APIResourceList{
			"vegetables.k8s.io/v1": {APIResources: []metav1.APIResource{
				{Kind: "Potato"},
				{Kind: "Carrot"},
				{Kind: "Lettuce"},
			}},
			"fruits.k8s.io/v1": {APIResources: []metav1.APIResource{
				{Kind: "Apple"},
				{Kind: "Banana"},
				{Kind: "Pear"},
			}},
		},
	}

	errClient := fakeDiscoveryClient{
		err: errors.New("some fake error"),
	}

	for _, tt := range []struct {
		name   string
		client discovery.ServerResourcesInterface

		groupVersion, kind string

		wantResult bool
		wantErr    bool
	}{
		{
			name:         "positive case",
			client:       &okClient,
			groupVersion: "vegetables.k8s.io/v1",
			kind:         "Carrot",
			wantResult:   true,
		},
		{
			name:         "error",
			client:       &errClient,
			groupVersion: "vegetables.k8s.io/v1",
			kind:         "Carrot",
			wantErr:      true,
		},
		{
			name:         "gv has no such kind",
			client:       &okClient,
			groupVersion: "vegetables.k8s.io/v1",
			kind:         "Australia",
			wantResult:   false,
		},
		{
			name:         "has kind in another gv",
			client:       &okClient,
			groupVersion: "fruits.k8s.io/v1",
			kind:         "Potato",
			wantResult:   false,
		},
		{
			name:         "no such gv",
			client:       &okClient,
			groupVersion: "grains.k8s.io",
			kind:         "Wheat",
			wantResult:   false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotErr := serverHasGVK(tt.client, tt.groupVersion, tt.kind)

			if gotResult != tt.wantResult {
				t.Errorf("serverHasGVK result: got %t, want %t", gotResult, tt.wantResult)
			}
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("serverHasGVK: got error: %v, wanted error? %t", gotErr, tt.wantErr)
			}
		})
	}

}

func TestNegotiateResourceAPI(t *testing.T) {
	client := &fakeDiscoveryClient{
		results: map[string]metav1.APIResourceList{
			ExtensionsV1beta1.String(): {APIResources: []metav1.APIResource{{Kind: "Carrot"}}},
			NetworkingV1beta1.String(): {APIResources: []metav1.APIResource{{Kind: "Carrot"}, {Kind: "Potato"}}},
			NetworkingV1.String():      {APIResources: []metav1.APIResource{{Kind: "Potato"}}},
		},
	}

	for _, tt := range []struct {
		name            string
		allowedVersions []IngressAPI
		kind            string

		wantRes IngressAPI
		wantErr bool
	}{
		{
			name:    "no allowed versions",
			kind:    "Banana",
			wantRes: OtherAPI,
			wantErr: true,
		},
		{
			name:            "none of allowed versions has GVK",
			kind:            "Banana",
			allowedVersions: []IngressAPI{NetworkingV1, NetworkingV1beta1, ExtensionsV1beta1},
			wantRes:         OtherAPI,
			wantErr:         true,
		},
		{
			name:            "API gets deleted in latest version",
			kind:            "Carrot",
			allowedVersions: []IngressAPI{NetworkingV1, NetworkingV1beta1, ExtensionsV1beta1},
			wantRes:         NetworkingV1beta1,
		},
		{
			name:            "API gets introduced in version later than first",
			kind:            "Potato",
			allowedVersions: []IngressAPI{NetworkingV1, NetworkingV1beta1, ExtensionsV1beta1},
			wantRes:         NetworkingV1,
		},
		{
			name:            "Newest allowedVersion not in the allowed list",
			kind:            "Potato",
			allowedVersions: []IngressAPI{NetworkingV1beta1, ExtensionsV1beta1},
			wantRes:         NetworkingV1beta1,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, gotErr := NegotiateResourceAPI(client, tt.kind, tt.allowedVersions)
			if tt.wantErr {
				require.Error(t, gotErr)
			} else {
				require.NoError(t, gotErr)
			}

			require.Equal(t, tt.wantRes, gotRes)
		})
	}
}
