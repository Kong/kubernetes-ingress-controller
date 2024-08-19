/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func TestParseNameNS(t *testing.T) {
	tests := []struct {
		title  string
		input  string
		ns     string
		name   string
		expErr bool
	}{
		{"empty string", "", "", "", true},
		{"demo", "demo", "", "", true},
		{"kube-system", "kube-system", "", "", true},
		{"default/kube-system", "default/kube-system", "default", "kube-system", false},
	}

	for _, test := range tests {
		ns, name, err := ParseNameNS(test.input)
		if test.expErr {
			if err == nil {
				t.Errorf("%v: expected error but returned nil", test.title)
			}
			continue
		}
		if test.ns != ns {
			t.Errorf("%v: expected %v but returned %v", test.title, test.ns, ns)
		}
		if test.name != name {
			t.Errorf("%v: expected %v but returned %v", test.title, test.name, name)
		}
	}
}

func TestGetNodeIP(t *testing.T) {
	ctx := context.Background()

	fKNodes := []struct {
		cs *testclient.Clientset
		n  string
		ea string
	}{
		// empty node list
		{testclient.NewSimpleClientset(), "demo", ""},

		// node not exist
		{testclient.NewSimpleClientset(&corev1.NodeList{Items: []corev1.Node{{
			ObjectMeta: metav1.ObjectMeta{
				Name: "demo",
			},
			Status: corev1.NodeStatus{
				Addresses: []corev1.NodeAddress{
					{
						Type:    corev1.NodeInternalIP,
						Address: "10.0.0.1",
					},
				},
			},
		}}}), "notexistnode", ""},

		// node  exist
		{testclient.NewSimpleClientset(&corev1.NodeList{Items: []corev1.Node{{
			ObjectMeta: metav1.ObjectMeta{
				Name: "demo",
			},
			Status: corev1.NodeStatus{
				Addresses: []corev1.NodeAddress{
					{
						Type:    corev1.NodeInternalIP,
						Address: "10.0.0.1",
					},
				},
			},
		}}}), "demo", "10.0.0.1"},

		// search the correct node
		{testclient.NewSimpleClientset(&corev1.NodeList{Items: []corev1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "demo1",
				},
				Status: corev1.NodeStatus{
					Addresses: []corev1.NodeAddress{
						{
							Type:    corev1.NodeInternalIP,
							Address: "10.0.0.1",
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "demo2",
				},
				Status: corev1.NodeStatus{
					Addresses: []corev1.NodeAddress{
						{
							Type:    corev1.NodeInternalIP,
							Address: "10.0.0.2",
						},
					},
				},
			},
		}}), "demo2", "10.0.0.2"},

		// get NodeExternalIP
		{testclient.NewSimpleClientset(&corev1.NodeList{Items: []corev1.Node{{
			ObjectMeta: metav1.ObjectMeta{
				Name: "demo",
			},
			Status: corev1.NodeStatus{
				Addresses: []corev1.NodeAddress{
					{
						Type:    corev1.NodeInternalIP,
						Address: "10.0.0.1",
					}, {
						Type:    corev1.NodeExternalIP,
						Address: "10.0.0.2",
					},
				},
			},
		}}}), "demo", "10.0.0.2"},

		// get NodeInternalIP
		{testclient.NewSimpleClientset(&corev1.NodeList{Items: []corev1.Node{{
			ObjectMeta: metav1.ObjectMeta{
				Name: "demo",
			},
			Status: corev1.NodeStatus{
				Addresses: []corev1.NodeAddress{
					{
						Type:    corev1.NodeExternalIP,
						Address: "",
					}, {
						Type:    corev1.NodeInternalIP,
						Address: "10.0.0.2",
					},
				},
			},
		}}}), "demo", "10.0.0.2"},
	}

	for _, fk := range fKNodes {
		address := GetNodeIPOrName(ctx, fk.cs, fk.n)
		if address != fk.ea {
			t.Errorf("expected %s, but returned %s", fk.ea, address)
		}
	}
}

func TestGetPodDetails(t *testing.T) {
	ctx := context.Background()
	// POD_NAME & POD_NAMESPACE not exist
	t.Setenv("POD_NAME", "")
	t.Setenv("POD_NAMESPACE", "")
	_, err1 := GetPodDetails(ctx, testclient.NewSimpleClientset())
	assert.Error(t, err1)

	// POD_NAME not exist
	t.Setenv("POD_NAME", "")
	t.Setenv("POD_NAMESPACE", corev1.NamespaceDefault)
	_, err2 := GetPodDetails(ctx, testclient.NewSimpleClientset())
	assert.Error(t, err2)

	// POD_NAMESPACE not exist
	t.Setenv("POD_NAME", "testpod")
	t.Setenv("POD_NAMESPACE", "")
	_, err3 := GetPodDetails(ctx, testclient.NewSimpleClientset())
	assert.Error(t, err3)

	// POD exists
	t.Setenv("POD_NAME", "testpod")
	t.Setenv("POD_NAMESPACE", corev1.NamespaceDefault)
	_, err4 := GetPodDetails(ctx, testclient.NewSimpleClientset())
	assert.NoError(t, err4)

	// success to get PodInfo
	fkClient := testclient.NewSimpleClientset(
		&corev1.PodList{Items: []corev1.Pod{{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testpod",
				Namespace: corev1.NamespaceDefault,
				Labels: map[string]string{
					"first":  "first_label",
					"second": "second_label",
				},
			},
		}}},
		&corev1.NodeList{Items: []corev1.Node{{
			ObjectMeta: metav1.ObjectMeta{
				Name: "demo",
			},
			Status: corev1.NodeStatus{
				Addresses: []corev1.NodeAddress{
					{
						Type:    corev1.NodeInternalIP,
						Address: "10.0.0.1",
					},
				},
			},
		}}})

	epi, err5 := GetPodDetails(ctx, fkClient)
	assert.NoError(t, err5)
	assert.NotNil(t, epi)
}

func TestGenerateTagsForObject(t *testing.T) {
	expectedTagSet := []*string{
		lo.ToPtr(K8sNameTagPrefix + "yedigei"),
		lo.ToPtr(K8sNamespaceTagPrefix + "aitmatov"),
		lo.ToPtr(K8sKindTagPrefix + "HTTPRoute"),
		lo.ToPtr(K8sUIDTagPrefix + "buryani"),
		lo.ToPtr(K8sGroupTagPrefix + "gateway.networking.k8s.io"),
		lo.ToPtr(K8sVersionTagPrefix + "v1"),
		lo.ToPtr("temir-jol"),
		lo.ToPtr("snaryad-soqqısı"),
	}

	// In memory kubernetes objects do not have GVK filled in.
	// Relevant kubernetes issue: https://github.com/kubernetes/kubernetes/issues/80609
	testObj := &gatewayapi.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "gateway.networking.k8s.io/v1",
			Kind:       "HTTPRoute",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "yedigei",
			Namespace: "aitmatov",
			UID:       "buryani",
			Annotations: map[string]string{
				annotations.AnnotationPrefix + annotations.UserTagKey: "temir-jol,snaryad-soqqısı,temir-jol,temir-jol",
			},
		},
	}

	tags := GenerateTagsForObject(testObj)
	if diff := cmp.Diff(expectedTagSet, tags); diff != "" {
		t.Fatalf("generated tags are not as expected, diff:\n%s", diff)
	}
}
