package helpers

import (
	"context"
	"fmt"

	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

func CreateIngressClass(ctx context.Context, ingressClassName string, client *kubernetes.Clientset) error {
	ingressClass := &netv1.IngressClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: ingressClassName,
		},
		Spec: netv1.IngressClassSpec{
			Controller: store.IngressClassKongController,
		},
	}
	ingClasses := client.NetworkingV1().IngressClasses()

	if _, err := ingClasses.Create(ctx, ingressClass, metav1.CreateOptions{}); apierrors.IsAlreadyExists(err) {
		// If for some reason the ingress class is already in the cluster don't
		// fail the whole test suite but recreate it and continue.
		if err := ingClasses.Delete(ctx, ingressClassName, metav1.DeleteOptions{}); err != nil {
			return fmt.Errorf("failed to delete ingress class %s: %w", ingressClass, err)
		}

		if _, err := ingClasses.Create(ctx, ingressClass, metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("failed to create ingress class %s: %w", ingressClass, err)
		}
	}
	return nil
}
