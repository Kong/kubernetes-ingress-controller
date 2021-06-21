package ctrlutils

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/deck/file"
	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"

	"github.com/sirupsen/logrus"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	"knative.dev/pkg/network"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	knativeversioned "knative.dev/networking/pkg/client/clientset/versioned"
	knativeinformerexternal "knative.dev/networking/pkg/client/informers/externalversions"
	knativeApis "knative.dev/pkg/apis"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// classSpec indicates the fieldName for objects which support indicating their Ingress Class by spec
const classSpec = "IngressClassName"

type PlugInConfig struct {
	body        []string
	headers     []string
	querystring []string
	uri         string
}

// CleanupFinalizer ensures that a deleted resource is no longer present in the object cache.
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

// returns false if Knative CRDs do not exist
func KnativeCRDExist(client client.Client) bool {
	knativeGVR := schema.GroupVersionResource{
		Group:    knative.SchemeGroupVersion.Group,
		Version:  knative.SchemeGroupVersion.Version,
		Resource: "ingresses",
	}
	_, err := client.RESTMapper().KindFor(knativeGVR)
	if meta.IsNoMatchError(err) {
		return false
	}
	return true
}

// update ingress status according to generated rules and specs
func UpdateIngress(targetContent *file.Content, log logrus.FieldLogger, ctx context.Context, kubeconfig *rest.Config) error {
	for _, svc := range targetContent.Services {
		for _, plugin := range svc.Plugins {
			if *plugin.Enabled == true {
				if config, ok := plugin.Config["add"]; ok {
					if add, ok := config.(PlugInConfig); ok {
						for _, header := range add.headers {
							if strings.HasPrefix(header, "Knative-Serving-") {
								fmt.Printf("located enabled knative service. update knative CR.")
								err := UpdateKnativeIngress(ctx, log, svc, kubeconfig)
								return fmt.Errorf("failed to update knative ingress err %v", err)
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func toKnativeLBStatus(coreLBStatus []apiv1.LoadBalancerIngress) []knative.LoadBalancerIngressStatus {
	var res []knative.LoadBalancerIngressStatus
	for _, status := range coreLBStatus {
		res = append(res, knative.LoadBalancerIngressStatus{
			IP:     status.IP,
			Domain: status.Hostname,
		})
	}
	return res
}

var ingressCondSet = knativeApis.NewLivingConditionSet()

func UpdateKnativeIngress(ctx context.Context, logger logrus.FieldLogger, svc file.FService, kubeCfg *rest.Config) error {
	// retrieve svc name/namespace
	var name, namespace string
	routeInf := strings.Split(*(svc.Routes[0].Name), ".")
	namespace = routeInf[0]
	name = routeInf[1]

	fmt.Printf("svc namespace %s name %s", namespace, name)
	if len(namespace) == 0 || len(name) == 0 {
		return fmt.Errorf("configured route information is not completed which should not happen.")
	}

	knativeCli, err := knativeversioned.NewForConfig(kubeCfg)
	if err != nil {
		return fmt.Errorf("failed to generate knative client. err %v", err)
	}
	knativeFactory := knativeinformerexternal.NewSharedInformerFactory(knativeCli, 0)
	knativeLister := knativeFactory.Networking().V1alpha1().Ingresses().Lister()

	currIng, err := knativeLister.Ingresses(namespace).Get(name)
	if err != nil {
		return fmt.Errorf("failed to fetch Knative Ingress %v/%v: %w", namespace, name, err)
	}

	// check if CR current status already updated
	curIPs := toCoreLBStatus(currIng.Status.PublicLoadBalancer)
	ips, err := RunningAddresses(ctx, kubeCfg)
	if err != nil {
		return fmt.Errorf("failed to retrieve cluster loadbalancer.")
	}
	status := SliceToStatus(ips)
	if IngressSliceEqual(status, curIPs) &&
		currIng.Status.ObservedGeneration == currIng.GetObjectMeta().GetGeneration() {
		logger.Debugf("no change in status, update skipped")
		return nil
	}

	// updating current custom status
	logger.WithField("ingress_status", status).Debugf("attempting to update Knative Ingress status")
	lbStatus := toKnativeLBStatus(status)

	clusterDomain := network.GetClusterDomainName()
	fmt.Println("cluster domain %s", clusterDomain)
	if err != nil {
		return err
	}

	for i := 0; i < len(lbStatus); i++ {
		lbStatus[i].DomainInternal = fmt.Sprintf("%s.%s.svc.%s",
			"ingress-controller-kong-proxy", "kong-system", clusterDomain)
	}

	currIng.Status.MarkLoadBalancerReady(lbStatus, lbStatus)
	ingressCondSet.Manage(&currIng.Status).MarkTrue(knative.IngressConditionReady)
	ingressCondSet.Manage(&currIng.Status).MarkTrue(knative.IngressConditionNetworkConfigured)
	currIng.Status.ObservedGeneration = currIng.GetObjectMeta().GetGeneration()

	_, err = knativeCli.NetworkingV1alpha1().Ingresses(namespace).UpdateStatus(ctx, currIng, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update ingress status: %v", err)
	} else {
		logger.WithField("ingress_status", status).Debugf("successfully updated ingress status")
	}
	return nil
}

func RunningAddresses(ctx context.Context, kubeCfg *rest.Config) ([]string, error) {
	addrs := []string{}

	namespace := "kong-system"
	CoreClient, _ := clientset.NewForConfig(kubeCfg)
	svc, err := CoreClient.CoreV1().Services(namespace).Get(ctx, "ingress-controller-kong-proxy ", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	switch svc.Spec.Type {
	case apiv1.ServiceTypeLoadBalancer:
		for _, ip := range svc.Status.LoadBalancer.Ingress {
			if ip.IP == "" {
				addrs = append(addrs, ip.Hostname)
			} else {
				addrs = append(addrs, ip.IP)
			}
		}

		addrs = append(addrs, svc.Spec.ExternalIPs...)
		return addrs, nil
	default:
		return addrs, nil
	}
}

// sliceToStatus converts a slice of IP and/or hostnames to LoadBalancerIngress
func SliceToStatus(endpoints []string) []apiv1.LoadBalancerIngress {
	lbi := []apiv1.LoadBalancerIngress{}
	for _, ep := range endpoints {
		if net.ParseIP(ep) == nil {
			lbi = append(lbi, apiv1.LoadBalancerIngress{Hostname: ep})
		} else {
			lbi = append(lbi, apiv1.LoadBalancerIngress{IP: ep})
		}
	}

	sort.SliceStable(lbi, func(a, b int) bool {
		return lbi[a].IP < lbi[b].IP
	})

	return lbi
}

func InSlice(e string, arr []string) bool {
	for _, v := range arr {
		if v == e {
			return true
		}
	}

	return false
}

func IngressSliceEqual(lhs, rhs []apiv1.LoadBalancerIngress) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for i := range lhs {
		if lhs[i].IP != rhs[i].IP {
			return false
		}
		if lhs[i].Hostname != rhs[i].Hostname {
			return false
		}
	}
	return true
}

func toCoreLBStatus(knativeLBStatus *knative.LoadBalancerStatus) []apiv1.LoadBalancerIngress {
	var res []apiv1.LoadBalancerIngress
	if knativeLBStatus == nil {
		return res
	}
	for _, status := range knativeLBStatus.Ingress {
		res = append(res, apiv1.LoadBalancerIngress{
			IP:       status.IP,
			Hostname: status.Domain,
		})
	}
	return res
}
