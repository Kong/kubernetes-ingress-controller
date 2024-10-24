package reference

import (
	"reflect"
	"testing"

	"github.com/go-logr/logr"
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
		name          string
		addReferrer   client.Object
		addReferent   client.Object
		checkReferrer client.Object
		checkReferent client.Object
		found         bool
	}{
		{
			name:          "same_referrer_and_referent,should_be_found",
			addReferrer:   testRefService1,
			addReferent:   testRefSecret1,
			checkReferrer: testRefService1,
			checkReferent: testRefSecret1,
			found:         true,
		},
		{
			name:          "same_referrer_different_referent,should_not_be_found",
			addReferrer:   testRefService1,
			addReferent:   testRefSecret1,
			checkReferrer: testRefService1,
			checkReferent: testRefSecret2,
			found:         false,
		},
		{
			name:          "different_referrer_same_referent,should_not_be_found",
			addReferrer:   testRefService1,
			addReferent:   testRefSecret1,
			checkReferrer: testRefService2,
			checkReferent: testRefSecret1,
			found:         false,
		},
		{
			name:          "different_referrer_different_referent,should_not_be_found",
			addReferrer:   testRefService1,
			addReferent:   testRefSecret1,
			checkReferrer: testRefService2,
			checkReferent: testRefSecret2,
			found:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCacheIndexers(logr.Discard())
			err := c.SetObjectReference(tc.addReferrer, tc.addReferent)
			require.NoError(t, err, "should not return error on setting reference")
			item, exists, err := c.indexer.Get(&ObjectReference{Referrer: tc.checkReferrer, Referent: tc.checkReferent})
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

func TestDeleteObjectReference(t *testing.T) {
	testCases := []struct {
		name           string
		deleteReferrer client.Object
		deleteReferent client.Object
		checkReferrer  client.Object
		checkReferent  client.Object
		found          bool
	}{
		{
			name:           "deleted_object_should_not_be_found",
			deleteReferrer: testRefService1,
			deleteReferent: testRefSecret2,
			checkReferrer:  testRefService1,
			checkReferent:  testRefSecret2,
			found:          false,
		},
		{
			name:           "delete_on_non_exist_object_should_not_return_error",
			deleteReferrer: testRefService1,
			deleteReferent: testRefSecret2,
			checkReferrer:  testRefService1,
			checkReferent:  testRefSecret1,
			found:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCacheIndexers(logr.Discard())
			err := c.SetObjectReference(testRefService1, testRefSecret1)
			require.NoError(t, err, "should not return error on setting reference")
			err = c.SetObjectReference(testRefService1, testRefSecret2)
			require.NoError(t, err, "should not return error on setting reference")

			err = c.DeleteObjectReference(tc.deleteReferrer, tc.deleteReferent)
			require.NoError(t, err, "should not return error on setting reference")
			_, exists, err := c.indexer.Get(&ObjectReference{Referrer: tc.checkReferrer, Referent: tc.checkReferent})
			require.NoError(t, err)
			require.Equal(t, tc.found, exists)
		})
	}
}

func TestObjectReferred(t *testing.T) {
	testCases := []struct {
		name          string
		addReferrer   client.Object
		addReferent   client.Object
		checkReferent client.Object
		referred      bool
	}{
		{
			name:          "object_referred",
			addReferrer:   testRefService1,
			addReferent:   testRefSecret1,
			checkReferent: testRefSecret1,
			referred:      true,
		},
		{
			name:          "object_not_referred",
			addReferrer:   testRefService1,
			addReferent:   testRefSecret1,
			checkReferent: testRefSecret2,
			referred:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCacheIndexers(logr.Discard())
			err := c.SetObjectReference(tc.addReferrer, tc.addReferent)
			require.NoError(t, err, "should not return error on setting reference")

			referred, err := c.ObjectReferred(tc.checkReferent)
			require.NoError(t, err)
			require.Equal(t, tc.referred, referred)
		})
	}
}

func TestListReferredObjects(t *testing.T) {
	testCases := []struct {
		name          string
		addReferrer   client.Object
		addReferent   client.Object
		checkReferrer client.Object
		objectNum     int
	}{
		{
			name:          "has_referred_objects",
			addReferrer:   testRefService1,
			addReferent:   testRefSecret1,
			checkReferrer: testRefService1,
			objectNum:     1,
		},
		{
			name:          "has_no_referred_objects",
			addReferrer:   testRefService1,
			addReferent:   testRefSecret1,
			checkReferrer: testRefService2,
			objectNum:     0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCacheIndexers(logr.Discard())
			err := c.SetObjectReference(tc.addReferrer, tc.addReferent)
			require.NoError(t, err, "should not return error on setting reference")

			referents, err := c.ListReferredObjects(tc.checkReferrer)
			require.NoError(t, err)
			require.Len(t, referents, tc.objectNum)
		})
	}
}

func TestDeleteReferencesByReferrer(t *testing.T) {
	testCases := []struct {
		name           string
		deleteReferrer client.Object
		checkReferrer  client.Object
		checkReferent  client.Object
		found          bool
	}{
		{
			name:           "should_delete_references_with_referrer_correctly",
			deleteReferrer: testRefService1,
			checkReferrer:  testRefService1,
			checkReferent:  testRefSecret1,
			found:          false,
		},
		{
			name:           "should_not_delete_unrelated_references",
			deleteReferrer: testRefService1,
			checkReferrer:  testRefService2,
			checkReferent:  testRefSecret2,
			found:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCacheIndexers(logr.Discard())
			err := c.SetObjectReference(testRefService1, testRefSecret1)
			require.NoError(t, err, "should not return error on setting reference")
			err = c.SetObjectReference(testRefService2, testRefSecret2)
			require.NoError(t, err, "should not return error on setting reference")

			err = c.DeleteReferencesByReferrer(tc.deleteReferrer)
			require.NoError(t, err)
			_, exists, err := c.indexer.Get(&ObjectReference{Referrer: tc.checkReferrer, Referent: tc.checkReferent})
			require.NoError(t, err)
			require.Equal(t, tc.found, exists)
		})
	}
}

func TestListReferrerObjectsByReferent(t *testing.T) {
	testCases := []struct {
		name          string
		addReferrer   client.Object
		addReferent   client.Object
		checkReferent client.Object
		objectNum     int
	}{
		{
			name:          "has_referring_objects",
			addReferrer:   testRefService1,
			addReferent:   testRefSecret1,
			checkReferent: testRefSecret1,
			objectNum:     1,
		},
		{
			name:          "has_no_referring_objects",
			addReferrer:   testRefService1,
			addReferent:   testRefSecret1,
			checkReferent: testRefSecret2,
			objectNum:     0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCacheIndexers(logr.Discard())
			err := c.SetObjectReference(tc.addReferrer, tc.addReferent)
			require.NoError(t, err, "should not return error on setting reference")

			referrers, err := c.ListReferrerObjectsByReferent(tc.checkReferent)
			require.NoError(t, err)
			require.Len(t, referrers, tc.objectNum)
		})
	}
}
