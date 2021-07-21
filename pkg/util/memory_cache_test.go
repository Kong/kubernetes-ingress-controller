package util

import (
	"testing"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/stretchr/testify/assert"
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSetGetValue(t *testing.T) {
	err := InitCache()
	assert.Nil(t, err)

	ingress := networkingv1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "Ingress",
		},
		Spec:   networkingv1.IngressSpec{},
		Status: networkingv1.IngressStatus{},
	}

	hash, err := hashstructure.Hash(&ingress, hashstructure.FormatV2, nil)
	assert.Nil(t, err)
	key := "ingress-v1ingress"
	err = SetValue(key, hash)
	assert.Nil(t, err)

	value, err := GetValue(key)
	assert.Equal(t, value, hash)
}
