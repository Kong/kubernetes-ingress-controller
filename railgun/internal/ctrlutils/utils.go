package ctrlutils

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// classSpec indicates the fieldName for objects which support indicating their Ingress Class by spec
const classSpec = "IngressClassName"

// CleanupFinalizer removes an object finalizer from an object which is currently being deleted.
func CleanupFinalizer(ctx context.Context, c client.Client, log logr.Logger, nsn types.NamespacedName, obj client.Object) (ctrl.Result, error) {
	if HasFinalizer(obj, KongIngressFinalizer) {
		log.Info("kong ingress finalizer needs to be removed from a resource which is deleting", "ingress", obj.GetName(), "finalizer", KongIngressFinalizer)
		finalizers := []string{}
		for _, finalizer := range obj.GetFinalizers() {
			if finalizer != KongIngressFinalizer {
				finalizers = append(finalizers, finalizer)
			}
		}
		obj.SetFinalizers(finalizers)
		if err := c.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("the kong ingress finalizer was removed from an a resource which is deleting", "ingress", obj.GetName(), "finalizer", KongIngressFinalizer)
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// HasFinalizer is a helper function to check whether a client.Object
// already has a specific finalizer set.
func HasFinalizer(obj client.Object, finalizer string) bool {
	hasFinalizer := false
	for _, foundFinalizer := range obj.GetFinalizers() {
		if foundFinalizer == finalizer {
			hasFinalizer = true
		}
	}
	return hasFinalizer
}

// HasAnnotation is a helper function to determine whether an object has a given annotation, and whether it's
// to the value provided.
func HasAnnotation(obj client.Object, key, expectedValue string) bool {
	foundValue, ok := obj.GetAnnotations()[key]
	return ok && foundValue == expectedValue
}

// MatchesIngressClassName indicates whether or not an object indicates that it's supported by the ingress class name provided.
func MatchesIngressClassName(obj client.Object, ingressClassName string) bool {
	if ing, ok := obj.(*netv1.Ingress); ok {
		if ing.Spec.IngressClassName != nil && *ing.Spec.IngressClassName == ingressClassName {
			return true
		}
	}

	if _, ok := obj.(*knative.Ingress); ok {
		return HasAnnotation(obj, annotations.KnativeIngressClassKey, ingressClassName)
	}

	return HasAnnotation(obj, annotations.IngressClassKey, ingressClassName)
}

type objWithIngressClassNameSpec struct {
	Spec struct{ IngressClassName *string }
}

// GeneratePredicateFuncsForIngressClassFilter builds a controller-runtime reconcilation predicate function which filters out objects
// which do not have the "kubernetes.io/ingress.class" annotation configured and set to the provided value or in their .spec.
func GeneratePredicateFuncsForIngressClassFilter(name string, specCheckEnabled, annotationCheckEnabled bool) predicate.Funcs {
	preds := predicate.NewPredicateFuncs(func(obj client.Object) bool {
		if annotationCheckEnabled && IsIngressClassAnnotationConfigured(obj, name) {
			return true
		}
		if specCheckEnabled && IsIngressClassSpecConfigured(obj, name) {
			return true
		}
		return false
	})
	preds.UpdateFunc = func(e event.UpdateEvent) bool {
		if annotationCheckEnabled && IsIngressClassAnnotationConfigured(e.ObjectOld, name) || IsIngressClassAnnotationConfigured(e.ObjectNew, name) {
			return true
		}
		if specCheckEnabled && IsIngressClassSpecConfigured(e.ObjectOld, name) || IsIngressClassSpecConfigured(e.ObjectNew, name) {
			return true
		}
		return false
	}
	return preds
}

// IsIngressClassAnnotationConfigured determines whether an object has an ingress.class annotation configured that
// matches the provide IngressClassName (and is therefore an object configured to be reconciled by that class).
//
// NOTE: keep in mind that the ingress.class annotation is deprecated and will be removed in a future release
//       of Kubernetes in favor of the .spec based implementation.
func IsIngressClassAnnotationConfigured(obj client.Object, expectedIngressClassName string) bool {
	if foundIngressClassName, ok := obj.GetAnnotations()[annotations.IngressClassKey]; ok {
		if foundIngressClassName == expectedIngressClassName {
			return true
		}
	}

	if foundIngressClassName, ok := obj.GetAnnotations()[annotations.KnativeIngressClassKey]; ok {
		if foundIngressClassName == expectedIngressClassName {
			return true
		}
	}

	return false
}

// IsIngressClassAnnotationConfigured determines whether an object has IngressClassName field in its spec and whether the value
// matches the provide IngressClassName (and is therefore an object configured to be reconciled by that class).
func IsIngressClassSpecConfigured(obj client.Object, expectedIngressClassName string) bool {
	switch obj := obj.(type) {
	case *netv1.Ingress:
		return obj.Spec.IngressClassName != nil && *obj.Spec.IngressClassName == expectedIngressClassName
	}
	return false
}

// CRDExists returns false if CRD does not exist
func CRDExists(client client.Client, gvr schema.GroupVersionResource) bool {
	_, err := client.RESTMapper().KindFor(gvr)
	if meta.IsNoMatchError(err) {
		return false
	}
	return true
}

// RetrieveKongAdminAPIURL retrieves the Kong Admin API URL from configured name/namespace service
func RetrieveKongAdminAPIURL(ctx context.Context, KongAdminAPI string, kubeCfg *rest.Config) (string, error) {
	namespace, name, err := util.ParseNameNS(KongAdminAPI)
	if err != nil {
		return "", fmt.Errorf("failed to parse kong admin api namespace and name: %w", err)
	}

	CoreClient, err := clientset.NewForConfig(kubeCfg)
	if err != nil || CoreClient == nil {
		return "", fmt.Errorf("failed creating k8s client %v", err)
	}

	svc, err := CoreClient.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed retrieve service object %s/%s: %w", namespace, name, err)
	}
	ingresses := svc.Status.LoadBalancer.Ingress
	adminIP := ""
	for _, ingress := range ingresses {
		if len(ingress.IP) > 0 {
			adminIP = ingress.IP
			break
		}
	}

	ports := svc.Spec.Ports
	var adminPort int32
	for _, port := range ports {
		if port.Name == "kong-admin" {
			adminPort = port.Port
		}
	}
	return fmt.Sprintf("http://%s:%d", adminIP, adminPort), nil
}
