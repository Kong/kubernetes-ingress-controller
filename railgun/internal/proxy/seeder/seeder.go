package seeder

import (
	"context"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	knative "knative.dev/networking/pkg/client/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/railgun/internal/ctrlutils"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/proxy"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
)

// -----------------------------------------------------------------------------
// Seeder
// -----------------------------------------------------------------------------

// Seeder is an object which can perform pre-fetch seed rounds to cache
// supported objects into the proxy cache server.
type Seeder struct {
	ingressClassName string

	logger logrus.FieldLogger
	prx    proxy.Proxy

	kc       *kubernetes.Clientset
	kongc    *clientset.Clientset
	knativec *knative.Clientset
}

// New provides a need *Seeder object.
func New(restCFG *rest.Config, prx proxy.Proxy) (*Seeder, error) {
	return NewBuilder(restCFG, prx).Build()
}

// Seed lists all supported API types, filters them to make sure the
// object is supported (e.g. uses the managers' ingress.class) and then
// pulls fresh copies of the object from the Kubernetes API to (re)seed
// the proxy cache. This is required to deal with situations where events
// get lost by controllers due to networking failures, poorly timed controller
// pod restarts, e.t.c.
func (s *Seeder) Seed(ctx context.Context) error {
	// FIXME - namespace filtration
	// FIXME - optionality/enablement for apis
	objs := make([]client.Object, 0)

	s.logger.Info("fetching supported core kubernetes objects")
	coreObjs, err := s.fetchCore(ctx)
	if err != nil {
		return err
	}
	objs = append(objs, coreObjs...)

	s.logger.Info("fetching supported kong kubernetes objects")
	kongObjs, err := s.fetchKong(ctx)
	if err != nil {
		return err
	}
	objs = append(objs, kongObjs...)

	s.logger.Info("fetching 3rd party kubernetes objects")
	otherObjs, err := s.fetchOther(ctx)
	if err != nil {
		return err
	}
	objs = append(objs, otherObjs...)

	if len(objs) < 1 {
		s.logger.Info("seed round successful: there were no objects that needed to be cached")
		return nil
	}

	if err := s.prx.UpdateObjects(objs...); err != nil {
		return err
	}
	s.logger.Infof("seed round successful: %d objects were added to the cache", len(objs))

	return nil
}

// -----------------------------------------------------------------------------
// Seeder - Controller Runtime - Runnable Implementation
// -----------------------------------------------------------------------------

func (s *Seeder) Start(ctx context.Context) error {
	return s.Seed(ctx)
}

// -----------------------------------------------------------------------------
// Seeder - Private Methods
// -----------------------------------------------------------------------------

// fetchCore fetches a fresh list of all Kubernetes objects and filters them according
// to the *Seeder's configured s.ingressClassName and returns the result.
func (s *Seeder) fetchCore(ctx context.Context) ([]client.Object, error) {
	objs := make([]client.Object, 0)

	deprecatedIngresses, err := s.kc.ExtensionsV1beta1().Ingresses(corev1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, obj := range deprecatedIngresses.Items {
		copyObj := obj
		if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
			objs = append(objs, &copyObj)
		}
	}

	v1beta1Ingresses, err := s.kc.NetworkingV1beta1().Ingresses(corev1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, obj := range v1beta1Ingresses.Items {
		copyObj := obj
		if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
			objs = append(objs, &copyObj)
		}
	}

	ingresses, err := s.kc.NetworkingV1().Ingresses(corev1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, obj := range ingresses.Items {
		copyObj := obj
		if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
			objs = append(objs, &copyObj)
		}
	}

	services, err := s.kc.CoreV1().Services(corev1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, obj := range services.Items {
		copyObj := obj
		if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
			objs = append(objs, &copyObj)
		}
	}

	endpoints, err := s.kc.CoreV1().Endpoints(corev1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, obj := range endpoints.Items {
		copyObj := obj
		if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
			objs = append(objs, &copyObj)
		}
	}

	return objs, nil
}

// fetchKong fetches a fresh list of all Kong Kubernetes objects and filters them
// according to the *Seeder's configured s.ingressClassName and returns the result.
func (s *Seeder) fetchKong(ctx context.Context) ([]client.Object, error) {
	objs := make([]client.Object, 0)

	kongPlugins, err := s.kongc.ConfigurationV1().KongClusterPlugins().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, obj := range kongPlugins.Items {
		copyObj := obj
		if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
			objs = append(objs, &copyObj)
		}
	}

	kongClusterPlugins, err := s.kongc.ConfigurationV1().KongClusterPlugins().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, obj := range kongClusterPlugins.Items {
		copyObj := obj
		if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
			objs = append(objs, &copyObj)
		}
	}

	kongConsumers, err := s.kongc.ConfigurationV1().KongConsumers(corev1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, obj := range kongConsumers.Items {
		copyObj := obj
		if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
			objs = append(objs, &copyObj)
		}
	}

	kongIngresses, err := s.kongc.ConfigurationV1().KongIngresses(corev1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, obj := range kongIngresses.Items {
		copyObj := obj
		if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
			objs = append(objs, &copyObj)
		}
	}

	tcpIngresses, err := s.kongc.ConfigurationV1beta1().TCPIngresses(corev1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, obj := range tcpIngresses.Items {
		copyObj := obj
		if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
			objs = append(objs, &copyObj)
		}
	}

	udpIngresses, err := s.kongc.ConfigurationV1beta1().UDPIngresses(corev1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, obj := range udpIngresses.Items {
		copyObj := obj
		if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
			objs = append(objs, &copyObj)
		}
	}

	return objs, nil
}

func (s *Seeder) fetchOther(ctx context.Context) ([]client.Object, error) {
	objs := make([]client.Object, 0)

	knativeIngresses, err := s.knativec.NetworkingV1alpha1().Ingresses(corev1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, obj := range knativeIngresses.Items {
		copyObj := obj
		if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
			objs = append(objs, &copyObj)
		}
	}

	return objs, nil
}
