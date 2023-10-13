package helpers

import (
	"context"
	"fmt"

	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

func CreateIngressClass(ctx context.Context, ingressClass string, client *kubernetes.Clientset) error {
	create := func() *netv1.IngressClass {
		return &netv1.IngressClass{
			ObjectMeta: metav1.ObjectMeta{
				Name: ingressClass,
			},
			Spec: netv1.IngressClassSpec{
				Controller: store.IngressClassKongController,
			},
		}
	}
	ingClasses := client.NetworkingV1().IngressClasses()

	_, err := ingClasses.Create(ctx, create(), metav1.CreateOptions{})
	if apierrors.IsAlreadyExists(err) {
		// If for some reason the ingress class is already in the cluster don't
		// fail the whole test suite but recreate it and continue.
		err = ingClasses.Delete(ctx, ingressClass, metav1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete ingress class %s: %w", ingressClass, err)
		}

		_, err = ingClasses.Create(ctx, create(), metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create ingress class %s: %w", ingressClass, err)
		}
	}
	return nil
}
