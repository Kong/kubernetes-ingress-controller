package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
)

func ingressWithClass(class string) *netv1.Ingress {
	if len(class) == 0 {
		return &netv1.Ingress{}
	}
	return &netv1.Ingress{
		Spec: netv1.IngressSpec{IngressClassName: &class},
	}
}

func ingressWithClassAnnotation(class string) *netv1.Ingress {
	if len(class) == 0 {
		return &netv1.Ingress{}
	}
	return &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				annotations.IngressClassKey: class,
			},
		},
	}
}

func knativeIngressWithClassAnnotation(class string) *knative.Ingress {
	if len(class) == 0 {
		return &knative.Ingress{}
	}
	return &knative.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				annotations.KnativeIngressClassKey: class,
			},
		},
	}
}

func TestMatchesIngressClass(t *testing.T) {
	type test struct {
		class           string
		controllerClass string
		isDefault       bool

		want bool
	}
	cases := []test{
		{class: "", isDefault: false, controllerClass: "foo", want: false},
		{class: "", isDefault: false, controllerClass: "kozel", want: false},
		{class: "", isDefault: true, controllerClass: "", want: true},
		{class: "", isDefault: true, controllerClass: "custom", want: true},
		{class: "", isDefault: true, controllerClass: "foo", want: true},
		{class: "", isDefault: true, controllerClass: "killer", want: true},
		{class: "", isDefault: true, controllerClass: "killer", want: true},
		{class: "custom", isDefault: false, controllerClass: "foo", want: false},
		{class: "custom", isDefault: false, controllerClass: "kozel", want: false},
		{class: "custom", isDefault: true, controllerClass: "custom", want: true},
		{class: "custom", isDefault: true, controllerClass: "foo", want: false},
		{class: "custom", isDefault: true, controllerClass: "kozel", want: false},
		{class: "foo", isDefault: false, controllerClass: "foo", want: true},
		{class: "foo", isDefault: true, controllerClass: "foo", want: true},
		{class: "kozel", isDefault: false, controllerClass: "kozel", want: true},
		{class: "kozel", isDefault: true, controllerClass: "kozel", want: true},
	}

	for idx, tt := range cases {
		t.Run(fmt.Sprintf("ingressWithClass test case %d", idx), func(t *testing.T) {
			got := MatchesIngressClass(ingressWithClass(tt.class), tt.controllerClass, tt.isDefault)
			require.Equal(t, tt.want, got)
		})
	}

	for idx, tt := range cases {
		t.Run(fmt.Sprintf("ingressWithClassAnnotation test case %d", idx), func(t *testing.T) {
			got := MatchesIngressClass(ingressWithClassAnnotation(tt.class), tt.controllerClass, tt.isDefault)
			require.Equal(t, tt.want, got)
		})
	}

	for idx, tt := range cases {
		t.Run(fmt.Sprintf("knativeIngressWithClassAnnotation test case %d", idx), func(t *testing.T) {
			got := MatchesIngressClass(knativeIngressWithClassAnnotation(tt.class), tt.controllerClass, tt.isDefault)
			require.Equal(t, tt.want, got)
		})
	}
}
