package reference

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var testRefService1 = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Namespace: "ns",
		Name:      "service1",
	},
}

var testRefService2 = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Namespace: "ns",
		Name:      "service2",
	},
}

var testRefSecret1 = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Namespace: "ns",
		Name:      "secret1",
	},
}

var testRefSecret2 = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Namespace: "ns",
		Name:      "secret2",
	},
}

func TestSetObjectReference(t *testing.T) {
	testCases := []struct {
		name           string
		addReferrer    client.Object
		addReferent    client.Object
		checkReferrer  client.Object
		checkReferrent client.Object
		found          bool
	}{
		{
			name:           "same_referrer_and_referent,should_be_found",
			addReferrer:    testRefService1,
			addReferent:    testRefSecret1,
			checkReferrer:  testRefService1,
			checkReferrent: testRefSecret1,
			found:          true,
		},
		{
			name:           "same_referrer_different_referent,should_not_be_found",
			addReferrer:    testRefService1,
			addReferent:    testRefSecret1,
			checkReferrer:  testRefService1,
			checkReferrent: testRefSecret2,
			found:          false,
		},
		{
			name:           "different_referrer_same_referent,should_not_be_found",
			addReferrer:    testRefService1,
			addReferent:    testRefSecret1,
			checkReferrer:  testRefService2,
			checkReferrent: testRefSecret1,
			found:          false,
		},
		{
			name:           "different_referrer_different_referent,should_not_be_found",
			addReferrer:    testRefService1,
			addReferent:    testRefSecret1,
			checkReferrer:  testRefService2,
			checkReferrent: testRefSecret2,
			found:          false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			c := NewCacheIndexers()
			err := c.SetObjectReference(tc.addReferrer, tc.addReferent)
			require.NoError(t, err, "should not return error on setting reference")
			item, exists, err := c.indexer.Get(&ObjectReference{Referrer: tc.checkReferrer, Referent: tc.checkReferrent})
			require.NoError(t, err)
			require.Equal(t, tc.found, exists)
			if tc.found {
				require.Truef(t,
					reflect.DeepEqual(item, &ObjectReference{Referrer: tc.addReferrer, Referent: tc.addReferent}),
					"reflect record got from cache should be equal to the one added")
			}
		})
	}
}
