package main

import (
	"errors"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
)

func TestFixVersion(t *testing.T) {
	validVersions := map[string]string{
		"0.14.1":                          "0.14.1",
		"0.14.2rc":                        "0.14.2-rc",
		"0.14.2rc1":                       "0.14.2-rc1",
		"0.14.2preview":                   "0.14.2-preview",
		"0.14.2preview1":                  "0.14.2-preview1",
		"0.33-enterprise-edition":         "0.33.0-enterprise",
		"0.33-1-enterprise-edition":       "0.33.1-enterprise",
		"1.3.0.0-enterprise-edition-lite": "1.3.0-0-enterprise-lite",
		"1.3.0.0-enterprise-lite":         "1.3.0-0-enterprise-lite",
	}
	for inputVersion, expectedVersion := range validVersions {
		v, err := getSemVerVer(inputVersion)
		if err != nil {
			t.Errorf("error converting %s: %v", inputVersion, err)
		} else if v.String() != expectedVersion {
			t.Errorf("converting %s, expecting %s, getting %s", inputVersion, expectedVersion, v.String())
		}
	}

	invalidVersions := []string{
		"",
		"0-1-1",
	}
	for _, inputVersion := range invalidVersions {
		_, err := getSemVerVer(inputVersion)
		if err == nil {
			t.Errorf("expecting error converting %s, getting no errors", inputVersion)
		}
	}
}

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
