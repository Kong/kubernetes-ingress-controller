package license

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
)

func TestCompareLicense(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name           string
		license1       *kongv1alpha1.KongLicense
		license2       *kongv1alpha1.KongLicense
		expectedResult bool
	}{
		{
			name: "The newer one should win",
			license1: &kongv1alpha1.KongLicense{
				ObjectMeta: metav1.ObjectMeta{
					CreationTimestamp: metav1.NewTime(now.Add(-10 * time.Second)),
					Name:              "alpha",
				},
			},
			license2: &kongv1alpha1.KongLicense{
				ObjectMeta: metav1.ObjectMeta{
					CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
					Name:              "beta",
				},
			},
			expectedResult: false,
		},
		{
			name: "If the creationTimestamp equals, the one with lexical smaller name should win",
			license1: &kongv1alpha1.KongLicense{
				ObjectMeta: metav1.ObjectMeta{
					CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
					Name:              "alpha",
				},
			},
			license2: &kongv1alpha1.KongLicense{
				ObjectMeta: metav1.ObjectMeta{
					CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
					Name:              "beta",
				},
			},
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedResult,
				compareKongLicense(tc.license1, tc.license2),
				"Should return expected compare results between two licenses")
		})
	}
}

func TestKongLicenseController_pickLicense(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name              string
		licenses          []*kongv1alpha1.KongLicense
		expectedNil       bool
		chosenLicenseName string
	}{
		{
			name:        "No licenses in cache - should return nil",
			licenses:    []*kongv1alpha1.KongLicense{},
			expectedNil: true,
		},
		{
			name: "Should choose the newest one",
			licenses: []*kongv1alpha1.KongLicense{
				{
					ObjectMeta: metav1.ObjectMeta{
						CreationTimestamp: metav1.NewTime(now.Add(-10 * time.Second)),
						Name:              "older",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						Name:              "newer",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						CreationTimestamp: metav1.NewTime(now.Add(-2 * time.Second)),
						Name:              "newest",
					},
				},
			},
			expectedNil:       false,
			chosenLicenseName: "newest",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := &KongV1Alpha1KongLicenseReconciler{
				LicenseCache: NewLicenseCache(),
			}
			for _, l := range tc.licenses {
				err := r.LicenseCache.Add(l)
				require.NoError(t, err, "Should have no error in adding KongLicense to cache")
			}
			chosenLicense := r.pickLicenseInCache()
			if tc.expectedNil {
				require.Nil(t, chosenLicense, "Should get no license")
			} else {
				require.NotNil(t, chosenLicense, "Should return an available license")
				require.Equal(t, tc.chosenLicenseName, chosenLicense.Name,
					"Should choose expected KongLicense")
			}
		})
	}
}
