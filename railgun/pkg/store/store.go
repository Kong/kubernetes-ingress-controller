package store

import (
	"context"
	"fmt"
	"os"

	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1beta1"
	oldstr "github.com/kong/kubernetes-ingress-controller/pkg/store"
	"github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1"
	configurationv1alpha1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1"
	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// -----------------------------------------------------------------------------
// Secret Controller - Storer - Public Functions
// -----------------------------------------------------------------------------

// New produces a new oldstr.Storer which will house the provided kubernetes.Clientset
// and will provide parsing and translation of Kubernetes objects to the Kong Admin DSL.
//
// This implementation of Storer varies from the original because it using a controller
// runtime client instead of relying on client-go cache stores for Kubernetes objects.
// The upside of this implementation is that it's more straightforward, and it can natively
// work with our custom APIs without needing to use a typed client or other API tricks.
//
// TODO: there's significant technical debt associated with the storer implementations and
// interface as a whole, we need to determine if and how we want to continue using the store
// package in the future. Perhaps we can just provide a kubernetes.Clientset later on?
// If we continue to use storer, we should consider expanding the interface to include
// contexts, as all the relevant API calls made underneath the hood here use contexts.
func New(c client.Client) oldstr.Storer {
	return &store{c}
}

// -----------------------------------------------------------------------------
// Secret Controller - Storer - Private Types
// -----------------------------------------------------------------------------

type store struct {
	c client.Client
}

// -----------------------------------------------------------------------------
// Secret Controller - Storer - Public Get Methods
// -----------------------------------------------------------------------------

func (s *store) GetSecret(namespace, name string) (*apiv1.Secret, error) {
	secret := new(apiv1.Secret)
	if err := s.c.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: name}, secret); err != nil {
		return nil, err
	}
	return secret, nil
}

func (s *store) GetService(namespace, name string) (*apiv1.Service, error) {
	service := new(apiv1.Service)
	if err := s.c.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: name}, service); err != nil {
		return nil, err
	}
	return service, nil
}

func (s *store) GetEndpointsForService(namespace, name string) (*apiv1.Endpoints, error) {
	endpoints := new(apiv1.Endpoints)
	if err := s.c.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: name}, endpoints); err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (s *store) GetKongIngress(namespace, name string) (*configurationv1.KongIngress, error) {
	ingress := new(configurationv1.KongIngress)
	if err := s.c.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: name}, ingress); err != nil {
		return nil, err
	}
	return ingress, nil
}

func (s *store) GetKongPlugin(namespace, name string) (*configurationv1.KongPlugin, error) {
	plugin := new(configurationv1.KongPlugin)
	if err := s.c.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: name}, plugin); err != nil {
		return nil, err
	}
	return plugin, nil
}

func (s *store) GetKongClusterPlugin(name string) (*configurationv1.KongClusterPlugin, error) {
	clusterPlugin := new(configurationv1.KongClusterPlugin)
	if err := s.c.Get(context.Background(), client.ObjectKey{Name: name}, clusterPlugin); err != nil {
		return nil, err
	}
	return clusterPlugin, nil
}

func (s *store) GetKongConsumer(namespace, name string) (*configurationv1.KongConsumer, error) {
	consumer := new(configurationv1.KongConsumer)
	if err := s.c.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: name}, consumer); err != nil {
		return nil, err
	}
	return consumer, nil
}

// -----------------------------------------------------------------------------
// Secret Controller - Storer - Public List Methods
// -----------------------------------------------------------------------------

func (s *store) ListIngressesV1beta1() []*networkingv1beta1.Ingress {
	list := new(networkingv1beta1.IngressList)
	if err := s.c.List(context.Background(), list); err != nil {
		return nil
	}

	ingresses := make([]*networkingv1beta1.Ingress, 0, len(list.Items))
	for _, ingress := range list.Items {
		ingresses = append(ingresses, &ingress)
	}

	return ingresses
}

func (s *store) ListIngressesV1() []*networkingv1.Ingress {
	list := new(networkingv1.IngressList)
	if err := s.c.List(context.Background(), list); err != nil {
		return nil
	}

	ingresses := make([]*networkingv1.Ingress, 0, len(list.Items))
	for _, ingress := range list.Items {
		ingresses = append(ingresses, &ingress)
	}

	return ingresses
}

func (s *store) ListTCPIngresses() ([]*configurationv1beta1.TCPIngress, error) {
	list := new(configurationv1beta1.TCPIngressList)
	if err := s.c.List(context.Background(), list); err != nil {
		return nil, err
	}

	ingresses := make([]*configurationv1beta1.TCPIngress, 0, len(list.Items))
	for _, ingress := range list.Items {
		ingresses = append(ingresses, &ingress)
	}

	return ingresses, nil
}

func (s *store) ListUDPIngresses() ([]*v1alpha1.UDPIngress, error) {
	list := new(configurationv1alpha1.UDPIngressList)
	if err := s.c.List(context.Background(), list); err != nil {
		return nil, err
	}

	ingresses := make([]*configurationv1alpha1.UDPIngress, 0, len(list.Items))
	for _, ingress := range list.Items {
		ingresses = append(ingresses, &ingress)
	}

	return ingresses, nil
}

func (s *store) ListKnativeIngresses() ([]*knative.Ingress, error) {
	list := new(knative.IngressList)
	if err := s.c.List(context.Background(), list); err != nil {
		return nil, err
	}

	ingresses := make([]*knative.Ingress, 0, len(list.Items))
	for _, ingress := range list.Items {
		ingresses = append(ingresses, &ingress)
	}

	return ingresses, nil
}

func (s *store) ListGlobalKongPlugins() ([]*configurationv1.KongPlugin, error) {
	list := new(configurationv1.KongPluginList)
	if err := s.c.List(context.Background(), list); err != nil {
		return nil, err
	}

	ingresses := make([]*configurationv1.KongPlugin, 0, len(list.Items))
	for _, ingress := range list.Items {
		ingresses = append(ingresses, &ingress)
	}

	return ingresses, nil
}

func (s *store) ListGlobalKongClusterPlugins() ([]*configurationv1.KongClusterPlugin, error) {
	list := new(configurationv1.KongClusterPluginList)
	if err := s.c.List(context.Background(), list); err != nil {
		return nil, err
	}

	ingresses := make([]*configurationv1.KongClusterPlugin, 0, len(list.Items))
	for _, ingress := range list.Items {
		ingresses = append(ingresses, &ingress)
	}

	return ingresses, nil

}

func (s *store) ListKongConsumers() []*configurationv1.KongConsumer {
	list := new(configurationv1.KongConsumerList)
	if err := s.c.List(context.Background(), list); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return nil
	}

	ingresses := make([]*configurationv1.KongConsumer, 0, len(list.Items))
	for _, ingress := range list.Items {
		ingresses = append(ingresses, &ingress)
	}

	return ingresses
}

func (s *store) ListCACerts() ([]*apiv1.Secret, error) {
	req1, err := labels.NewRequirement("konghq.com/ca-cert", selection.Exists, []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to list CACerts due to label error: %w", err)
	}

	req2, err := labels.NewRequirement("konghq.com/ca-cert", selection.Equals, []string{"true"})
	if err != nil {
		return nil, fmt.Errorf("failed to list CACerts due to label error: %w", err)
	}

	selector := labels.NewSelector()
	selector = selector.Add(*req1)
	selector = selector.Add(*req2)
	opts := &client.ListOptions{LabelSelector: selector}

	list := new(apiv1.SecretList)
	if err := s.c.List(context.Background(), list, opts); err != nil {
		return nil, err
	}

	secrets := make([]*apiv1.Secret, 0, len(list.Items))
	for _, secret := range list.Items {
		secrets = append(secrets, &secret)
	}

	return secrets, nil
}
