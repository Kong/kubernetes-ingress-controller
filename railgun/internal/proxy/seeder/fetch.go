package seeder

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/ctrlutils"
)

// -----------------------------------------------------------------------------
// Seeder - Private Methods
// -----------------------------------------------------------------------------

// defaultListOptions is a convenience var around the default List options.
var defaultListOptions = metav1.ListOptions{}

// fetchCore fetches a fresh list of all Kubernetes objects and filters them according
// to the *Seeder's configured s.ingressClassName and returns the result.
func (s *Seeder) fetchCore(ctx context.Context) ([]client.Object, error) {
	objs := make([]client.Object, 0)

	// -------------------------------------------------------------------------
	// Kubernetes Services
	// -------------------------------------------------------------------------

	for _, namespace := range s.namespaces {
		services, err := s.kc.CoreV1().Services(namespace).List(ctx, defaultListOptions)
		if err != nil {
			return nil, err
		}
		for _, obj := range services.Items {
			copyObj := obj
			if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
				objs = append(objs, &copyObj)
			}
		}
	}

	for _, namespace := range s.namespaces {
		endpoints, err := s.kc.CoreV1().Endpoints(namespace).List(ctx, defaultListOptions)
		if err != nil {
			return nil, err
		}
		for _, obj := range endpoints.Items {
			copyObj := obj
			if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
				objs = append(objs, &copyObj)
			}
		}
	}

	// -------------------------------------------------------------------------
	// Kubernetes Core Ingress
	// -------------------------------------------------------------------------

	// de-duplicate converted ingress objects (e.g. converted v1beta1.Ingress from extv1beta1.Ingress)
	dedup := map[types.UID]client.Object{}
	if s.controllerConfig == nil || s.controllerConfig.IngressExtV1beta1Enabled == util.EnablementStatusEnabled {
		for _, namespace := range s.namespaces {
			list, err := s.kc.ExtensionsV1beta1().Ingresses(namespace).List(ctx, defaultListOptions)
			if err != nil {
				return nil, err
			}
			for _, obj := range list.Items {
				copyObj := obj
				if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
					copyObj.GetUID()
					dedup[copyObj.GetUID()] = &copyObj
				}
			}
		}
	}

	if s.controllerConfig == nil || s.controllerConfig.IngressNetV1beta1Enabled == util.EnablementStatusEnabled {
		for _, namespace := range s.namespaces {
			v1beta1Ingresses, err := s.kc.NetworkingV1beta1().Ingresses(namespace).List(ctx, defaultListOptions)
			if err != nil {
				return nil, err
			}
			for _, obj := range v1beta1Ingresses.Items {
				copyObj := obj
				if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
					copyObj.GetUID()
					dedup[copyObj.GetUID()] = &copyObj
				}
			}
		}
	}

	if s.controllerConfig == nil || s.controllerConfig.IngressNetV1Enabled == util.EnablementStatusEnabled {
		for _, namespace := range s.namespaces {
			ingresses, err := s.kc.NetworkingV1().Ingresses(namespace).List(ctx, defaultListOptions)
			if err != nil {
				return nil, err
			}
			for _, obj := range ingresses.Items {
				copyObj := obj
				if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
					copyObj.GetUID()
					dedup[copyObj.GetUID()] = &copyObj
				}
			}
		}
	}

	// extract the de-duplicated list of ingress objects
	for _, obj := range dedup {
		objs = append(objs, obj)
	}

	return objs, nil
}

// fetchKong fetches a fresh list of all Kong Kubernetes objects and filters them
// according to the *Seeder's configured s.ingressClassName and returns the result.
func (s *Seeder) fetchKong(ctx context.Context) ([]client.Object, error) {
	objs := make([]client.Object, 0)

	// -------------------------------------------------------------------------
	// Kong Ingress
	// -------------------------------------------------------------------------

	for _, namespace := range s.namespaces {
		kongIngresses, err := s.kongc.ConfigurationV1().KongIngresses(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, obj := range kongIngresses.Items {
			copyObj := obj
			if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
				objs = append(objs, &copyObj)
			}
		}
	}

	for _, namespace := range s.namespaces {
		tcpIngresses, err := s.kongc.ConfigurationV1beta1().TCPIngresses(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, obj := range tcpIngresses.Items {
			copyObj := obj
			if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
				objs = append(objs, &copyObj)
			}
		}
	}

	for _, namespace := range s.namespaces {
		udpIngresses, err := s.kongc.ConfigurationV1beta1().UDPIngresses(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, obj := range udpIngresses.Items {
			copyObj := obj
			if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
				objs = append(objs, &copyObj)
			}
		}
	}

	// -------------------------------------------------------------------------
	// Kong Plugins
	// -------------------------------------------------------------------------

	kongPlugins, err := s.kongc.ConfigurationV1().KongClusterPlugins().List(ctx, defaultListOptions)
	if err != nil {
		return nil, err
	}
	for _, obj := range kongPlugins.Items {
		copyObj := obj
		if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
			objs = append(objs, &copyObj)
		}
	}

	if s.controllerConfig == nil || s.controllerConfig.KongClusterPluginEnabled == util.EnablementStatusEnabled {
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
	}

	// -------------------------------------------------------------------------
	// Kong Consumers
	// -------------------------------------------------------------------------

	for _, namespace := range s.namespaces {
		kongConsumers, err := s.kongc.ConfigurationV1().KongConsumers(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, obj := range kongConsumers.Items {
			copyObj := obj
			if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
				objs = append(objs, &copyObj)
			}
		}
	}

	return objs, nil
}

func (s *Seeder) fetchOther(ctx context.Context) ([]client.Object, error) {
	objs := make([]client.Object, 0)

	if s.controllerConfig == nil || s.controllerConfig.KnativeIngressEnabled == util.EnablementStatusEnabled {
		for _, namespace := range s.namespaces {
			knativeIngresses, err := s.knativec.NetworkingV1alpha1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				return nil, err
			}
			for _, obj := range knativeIngresses.Items {
				copyObj := obj
				if ctrlutils.IsObjectSupported(&copyObj, s.ingressClassName) {
					objs = append(objs, &copyObj)
				}
			}
		}
	}

	return objs, nil
}
