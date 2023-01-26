package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// -----------------------------------------------------------------------------
// Main
// -----------------------------------------------------------------------------

const (
	outputFile = "../../internal/controllers/configuration/zz_generated_controllers.go"

	corev1     = "k8s.io/api/core/v1"
	netv1      = "k8s.io/api/networking/v1"
	netv1beta1 = "k8s.io/api/networking/v1beta1"
	extv1beta1 = "k8s.io/api/extensions/v1beta1"

	kongv1       = "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1  = "github.com/kong/kubernetes-ingress-controller/v2/api/configuration/v1beta1"
	kongv1alpha1 = "github.com/kong/kubernetes-ingress-controller/v2/api/configuration/v1alpha1"
)

// inputControllersNeeded is a list of the supported Types for the
// Kong Kubernetes Ingress Controller. If you need to add a new type
// for support, add it here and a new controller will be generated
// when you run `make controllers`.
var inputControllersNeeded = &typesNeeded{
	typeNeeded{
		Group:                             "\"\"",
		Version:                           "v1",
		Kind:                              "Service",
		PackageImportAlias:                "corev1",
		PackageAlias:                      "CoreV1",
		Package:                           corev1,
		Plural:                            "services",
		CacheType:                         "Service",
		NeedsStatusPermissions:            true,
		AcceptsIngressClassNameAnnotation: false,
		AcceptsIngressClassNameSpec:       false,
		NeedsUpdateReferences:             true,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		Group:                             "\"\"",
		Version:                           "v1",
		Kind:                              "Endpoints",
		PackageImportAlias:                "corev1",
		PackageAlias:                      "CoreV1",
		Package:                           corev1,
		Plural:                            "endpoints",
		CacheType:                         "Endpoint",
		NeedsStatusPermissions:            true,
		AcceptsIngressClassNameAnnotation: false,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"list", "watch"},
	},
	typeNeeded{
		Group:                             "networking.k8s.io",
		Version:                           "v1",
		Kind:                              "Ingress",
		PackageImportAlias:                "netv1",
		PackageAlias:                      "NetV1",
		Package:                           netv1,
		Plural:                            "ingresses",
		CacheType:                         "IngressV1",
		NeedsStatusPermissions:            true,
		CapableOfStatusUpdates:            true,
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       true,
		NeedsUpdateReferences:             true,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		Group:                             "networking.k8s.io",
		Version:                           "v1",
		Kind:                              "IngressClass",
		PackageImportAlias:                "netv1",
		PackageAlias:                      "NetV1",
		Package:                           netv1,
		Plural:                            "ingressclasses",
		CacheType:                         "IngressV1",
		NeedsStatusPermissions:            false,
		AcceptsIngressClassNameAnnotation: false,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		Group:                             "networking.k8s.io",
		Version:                           "v1beta1",
		Kind:                              "Ingress",
		PackageImportAlias:                "netv1beta1",
		PackageAlias:                      "NetV1Beta1",
		Package:                           netv1beta1,
		Plural:                            "ingresses",
		CacheType:                         "IngressV1beta1",
		NeedsStatusPermissions:            true,
		CapableOfStatusUpdates:            true,
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       true,
		NeedsUpdateReferences:             true,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		Group:                             "extensions",
		Version:                           "v1beta1",
		Kind:                              "Ingress",
		PackageImportAlias:                "extv1beta1",
		PackageAlias:                      "ExtV1Beta1",
		Package:                           extv1beta1,
		Plural:                            "ingresses",
		CacheType:                         "IngressV1beta1",
		NeedsStatusPermissions:            true,
		CapableOfStatusUpdates:            true,
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       true,
		NeedsUpdateReferences:             true,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		Group:                             "configuration.konghq.com",
		Version:                           "v1",
		Kind:                              "KongIngress",
		PackageImportAlias:                "kongv1",
		PackageAlias:                      "KongV1",
		Package:                           kongv1,
		Plural:                            "kongingresses",
		CacheType:                         "KongIngress",
		NeedsStatusPermissions:            true,
		AcceptsIngressClassNameAnnotation: false,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		Group:                             "configuration.konghq.com",
		Version:                           "v1",
		Kind:                              "KongPlugin",
		PackageImportAlias:                "kongv1",
		PackageAlias:                      "KongV1",
		Package:                           kongv1,
		Plural:                            "kongplugins",
		CacheType:                         "Plugin",
		NeedsStatusPermissions:            true,
		AcceptsIngressClassNameAnnotation: false,
		AcceptsIngressClassNameSpec:       false,
		NeedsUpdateReferences:             true,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		Group:                             "configuration.konghq.com",
		Version:                           "v1",
		Kind:                              "KongClusterPlugin",
		PackageImportAlias:                "kongv1",
		PackageAlias:                      "KongV1",
		Package:                           kongv1,
		Plural:                            "kongclusterplugins",
		CacheType:                         "ClusterPlugin",
		NeedsStatusPermissions:            true,
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       false,
		NeedsUpdateReferences:             true,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		Group:                             "configuration.konghq.com",
		Version:                           "v1",
		Kind:                              "KongConsumer",
		PackageImportAlias:                "kongv1",
		PackageAlias:                      "KongV1",
		Package:                           kongv1,
		Plural:                            "kongconsumers",
		CacheType:                         "Consumer",
		NeedsStatusPermissions:            true,
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       false,
		NeedsUpdateReferences:             true,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		Group:                             "configuration.konghq.com",
		Version:                           "v1beta1",
		Kind:                              "TCPIngress",
		PackageImportAlias:                "kongv1beta1",
		PackageAlias:                      "KongV1Beta1",
		Package:                           kongv1beta1,
		Plural:                            "tcpingresses",
		CacheType:                         "TCPIngress",
		NeedsStatusPermissions:            true,
		CapableOfStatusUpdates:            true,
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       false,
		NeedsUpdateReferences:             true,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		Group:                             "configuration.konghq.com",
		Version:                           "v1beta1",
		Kind:                              "UDPIngress",
		PackageImportAlias:                "kongv1beta1",
		PackageAlias:                      "KongV1Beta1",
		Package:                           kongv1beta1,
		Plural:                            "udpingresses",
		CacheType:                         "UDPIngress",
		NeedsStatusPermissions:            true,
		CapableOfStatusUpdates:            true,
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		Group:                             "configuration.konghq.com",
		Version:                           "v1alpha1",
		Kind:                              "IngressClassParameters",
		PackageImportAlias:                "kongv1alpha1",
		PackageAlias:                      "KongV1Alpha1",
		Package:                           kongv1alpha1,
		Plural:                            "ingressclassparameterses",
		CacheType:                         "IngressClassParameters",
		NeedsStatusPermissions:            false,
		CapableOfStatusUpdates:            false,
		AcceptsIngressClassNameAnnotation: false,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
}

var inputRBACPermissionsNeeded = &rbacsNeeded{
	rbacNeeded{
		Plural:    "nodes",
		Group:     `""`,
		RBACVerbs: []string{"list", "watch"},
	},
	rbacNeeded{
		Plural:    "pods",
		Group:     `""`,
		RBACVerbs: []string{"get", "list", "watch"},
	},
	rbacNeeded{
		Plural:    "events",
		Group:     `""`,
		RBACVerbs: []string{"create", "patch"},
	},
}

func main() {
	needed := necessary{
		types: inputControllersNeeded,
		rbacs: inputRBACPermissionsNeeded,
	}
	if err := needed.generate(); err != nil {
		fmt.Fprintf(os.Stderr, "could not generate input controllers: %v", err)
		os.Exit(1)
	}
}

// -----------------------------------------------------------------------------
// Private Functions - Helper
// -----------------------------------------------------------------------------

// header produces a skeleton of the controller file to be generated.
func header() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	boilerPlate, err := os.ReadFile("../../hack/boilerplate.go.txt")
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(boilerPlate)
	if err != nil {
		return nil, err
	}

	_, err = buf.WriteString(headerTemplate)
	return buf, err
}

// -----------------------------------------------------------------------------
// Generator
// -----------------------------------------------------------------------------

// typesNeeded is a list of Kubernetes API types which are supported
// by the Kong Kubernetes Ingress Controller and need to have "input"
// controllers generated for them.
type typesNeeded []typeNeeded

// rbacsNeeded is a list of Kubernetes API objects which the Kong
// Kubernetes Ingress Controller interacts with, but does not need a
// controller for, only permissions.
type rbacsNeeded []rbacNeeded

type necessary struct {
	types *typesNeeded
	rbacs *rbacsNeeded
}

// generate generates a controller/input/<controller>.go Kubernetes controller
// for every supported type populated in the list.
func (needed necessary) generate() error {
	contents, err := header()
	if err != nil {
		return err
	}

	for _, t := range *needed.types {
		if err := t.generate(contents); err != nil {
			return err
		}
	}

	for _, r := range *needed.rbacs {
		if err := r.generate(contents); err != nil {
			return err
		}
	}

	return os.WriteFile(outputFile, contents.Bytes(), 0o600)
}

type typeNeeded struct {
	Group   string
	Version string
	Kind    string

	PackageImportAlias string
	PackageAlias       string
	Package            string
	Plural             string
	CacheType          string
	RBACVerbs          []string

	// AcceptsIngressClassNameAnnotation indicates that the object accepts (and the controller will listen to)
	// the "kubernetes.io/ingress.class" annotation to decide whether or not the object is supported.
	AcceptsIngressClassNameAnnotation bool

	// AcceptsIngressClassNameSpec indicates the the object indicates the ingress.class that should support it via
	// an attribute in its specification named .IngressClassName
	AcceptsIngressClassNameSpec bool

	// NeedsStatusPermissions indicates whether permissions for the object should also include permissions to update
	// its status
	NeedsStatusPermissions bool

	// CapableOfStatusUpdates indicates that the controllers should manage status
	// updates for the resource.
	CapableOfStatusUpdates bool

	// NeedUpdateReferences is true if we need to update the reference relationships
	// between reconciled object and other objects.
	NeedsUpdateReferences bool
}

func (t *typeNeeded) generate(contents *bytes.Buffer) error {
	tmpl, err := template.New("controller").Funcs(sprig.TxtFuncMap()).Parse(controllerTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(contents, t)
}

// rbacNeeded represents a resource that we only require RBAC permissions for.
type rbacNeeded struct {
	Group     string
	Plural    string
	RBACVerbs []string
}

func (r *rbacNeeded) generate(contents *bytes.Buffer) error {
	tmpl, err := template.New("rbac").Funcs(sprig.TxtFuncMap()).Parse(rbacTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(contents, r)
}

// -----------------------------------------------------------------------------
// Templates
// -----------------------------------------------------------------------------

var headerTemplate = `
// Code generated by Kong; DO NOT EDIT.

package configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	ctrlref "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/reference"
	ctrlutils "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/utils"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object/status"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
)
`

var rbacTemplate = `
// -----------------------------------------------------------------------------
// API Group {{.Group}} resource {{.Plural}}
// -----------------------------------------------------------------------------

//+kubebuilder:rbac:groups={{.Group}},resources={{.Plural}},verbs={{ .RBACVerbs | join ";" }}
`

var controllerTemplate = `
// -----------------------------------------------------------------------------
// {{.PackageAlias}} {{.Kind}} - Reconciler
// -----------------------------------------------------------------------------

// {{.PackageAlias}}{{.Kind}}Reconciler reconciles {{.Kind}} resources
type {{.PackageAlias}}{{.Kind}}Reconciler struct {
	client.Client

	Log             logr.Logger
	Scheme          *runtime.Scheme
	DataplaneClient *dataplane.KongClient
	CacheSyncTimeout time.Duration
{{- if .CapableOfStatusUpdates }}

	DataplaneAddressFinder *dataplane.AddressFinder
	StatusQueue            *status.Queue
{{- end}}
{{- if or .AcceptsIngressClassNameSpec .AcceptsIngressClassNameAnnotation}}

	IngressClassName string
	DisableIngressClassLookups bool
{{- end}}
{{- if .NeedsUpdateReferences}}
	ReferenceIndexers ctrlref.CacheIndexers
{{- end}}
}

// SetupWithManager sets up the controller with the Manager.
func (r *{{.PackageAlias}}{{.Kind}}Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("{{.PackageAlias}}{{.Kind}}", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}

{{- if .CapableOfStatusUpdates}}
	// if configured, start the status updater controller
	if r.StatusQueue != nil {
		if err := c.Watch(
			&source.Channel{Source: r.StatusQueue.Subscribe(schema.GroupVersionKind{
				Group:   "{{.Group}}",
				Version: "{{.Version}}",
				Kind:    "{{.Kind}}",
			})},
			&handler.EnqueueRequestForObject{},
		); err != nil {
			return err
		}
	}
{{- end}}
{{- if .AcceptsIngressClassNameAnnotation}}
	if !r.DisableIngressClassLookups {
		err = c.Watch(
			&source.Kind{Type: &netv1.IngressClass{}},
			handler.EnqueueRequestsFromMapFunc(r.listClassless),
			predicate.NewPredicateFuncs(ctrlutils.IsDefaultIngressClass),
		)
		if err != nil {
			return err
		}
	}
	preds := ctrlutils.GeneratePredicateFuncsForIngressClassFilter(r.IngressClassName)
{{- end}}
	return c.Watch(
		&source.Kind{Type: &{{.PackageImportAlias}}.{{.Kind}}{}},
		&handler.EnqueueRequestForObject{},
{{- if .AcceptsIngressClassNameAnnotation}}
		preds,
{{- end}}
	)
}

{{- if .AcceptsIngressClassNameAnnotation}}
// listClassless finds and reconciles all objects without ingress class information
func (r *{{.PackageAlias}}{{.Kind}}Reconciler) listClassless(obj client.Object) []reconcile.Request {
	resourceList := &{{.PackageImportAlias}}.{{.Kind}}List{}
	if err := r.Client.List(context.Background(), resourceList); err != nil {
		r.Log.Error(err, "failed to list classless {{.Plural}}")
		return nil
	}
	var recs []reconcile.Request
	for i, resource := range resourceList.Items {
		if ctrlutils.IsIngressClassEmpty(&resourceList.Items[i]) {
			recs = append(recs, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: resource.Namespace,
					Name:      resource.Name,
				},
			})
		}
	}
	return recs
}
{{- end}}

//+kubebuilder:rbac:groups={{.Group}},resources={{.Plural}},verbs={{ .RBACVerbs | join ";" }}
{{- if .NeedsStatusPermissions}}
//+kubebuilder:rbac:groups={{.Group}},resources={{.Plural}}/status,verbs=get;update;patch
{{- end}}

// Reconcile processes the watched objects
func (r *{{.PackageAlias}}{{.Kind}}Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("{{.PackageAlias}}{{.Kind}}", req.NamespacedName)

	// get the relevant object
	obj := new({{.PackageImportAlias}}.{{.Kind}})
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		if apierrors.IsNotFound(err) {
			obj.Namespace = req.Namespace
			obj.Name = req.Name
			{{if .NeedsUpdateReferences}}
			// remove reference record where the {{.Kind}} is the referrer
			if err := ctrlref.DeleteReferencesByReferrer(r.ReferenceIndexers, r.DataplaneClient, obj); err != nil {
				return ctrl.Result{}, err
			}
			{{end}}
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(obj)
		}
		return ctrl.Result{}, err
	}
	log.V(util.DebugLevel).Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.V(util.DebugLevel).Info("resource is being deleted, its configuration will be removed", "type", "{{.Kind}}", "namespace", req.Namespace, "name", req.Name)
		{{if .NeedsUpdateReferences}}
		// remove reference record where the {{.Kind}} is the referrer
		if err := ctrlref.DeleteReferencesByReferrer(r.ReferenceIndexers, r.DataplaneClient, obj); err != nil {
			return ctrl.Result{}, err
		}
		{{end}}
		objectExistsInCache, err := r.DataplaneClient.ObjectExists(obj)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			if err := r.DataplaneClient.DeleteObject(obj); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, nil
	}
{{if .AcceptsIngressClassNameAnnotation}}
	class := new(netv1.IngressClass)
	if !r.DisableIngressClassLookups {
		if err := r.Get(ctx, types.NamespacedName{Name: r.IngressClassName}, class); err != nil {
			// we log this without taking action to support legacy configurations that only set ingressClassName or
			// used the class annotation and did not create a corresponding IngressClass. We only need this to determine
			// if the IngressClass is default or to configure default settings, and can assume no/no additional defaults
			// if none exists.
			log.V(util.DebugLevel).Info("could not retrieve IngressClass", "ingressclass", r.IngressClassName)
		}
	}
	// if the object is not configured with our ingress.class, then we need to ensure it's removed from the cache
	if !ctrlutils.MatchesIngressClass(obj, r.IngressClassName, ctrlutils.IsDefaultIngressClass(class)) {
		log.V(util.DebugLevel).Info("object missing ingress class, ensuring it's removed from configuration",
		"namespace", req.Namespace, "name", req.Name, "class", r.IngressClassName)
		return ctrl.Result{}, r.DataplaneClient.DeleteObject(obj)
	} else {
		log.V(util.DebugLevel).Info("object has matching ingress class", "namespace", req.Namespace, "name", req.Name,
		"class", r.IngressClassName)
	}
{{end}}
	// update the kong Admin API with the changes
	if err := r.DataplaneClient.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}

{{- if .NeedsUpdateReferences }}
	// update reference relationship from the {{.Kind}} to other objects.
	if err := updateReferredObjects(ctx, r.Client, r.ReferenceIndexers, r.DataplaneClient, obj); err != nil {
		if apierrors.IsNotFound(err) {
			// reconcile again if the secret does not exist yet
			return ctrl.Result{
				Requeue: true,
			}, nil
		}
		return ctrl.Result{}, err
	}
{{- end }}

{{- if .CapableOfStatusUpdates}}
	// if status updates are enabled report the status for the object
	if r.DataplaneClient.AreKubernetesObjectReportsEnabled() {
		log.V(util.DebugLevel).Info("determining whether data-plane configuration has succeeded", "namespace", req.Namespace, "name", req.Name)
		
		if  !r.DataplaneClient.KubernetesObjectIsConfigured(obj) {
			log.V(util.DebugLevel).Info("resource not yet configured in the data-plane", "namespace", req.Namespace, "name", req.Name)
			return ctrl.Result{Requeue: true}, nil // requeue until the object has been properly configured
		}

		log.V(util.DebugLevel).Info("determining gateway addresses for object status updates", "namespace", req.Namespace, "name", req.Name)
		addrs, err := r.DataplaneAddressFinder.GetLoadBalancerAddresses(ctx)
		if err != nil {
			return ctrl.Result{}, err
		}

		log.V(util.DebugLevel).Info("found addresses for data-plane updating object status", "namespace", req.Namespace, "name", req.Name)
		updateNeeded, err := ctrlutils.UpdateLoadBalancerIngress(obj, addrs)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to update load balancer address: %w", err)
		}
		if updateNeeded {
			return ctrl.Result{}, r.Status().Update(ctx, obj)
		}
		log.V(util.DebugLevel).Info("status update not needed", "namespace", req.Namespace, "name", req.Name)
	}
{{- end}}

	return ctrl.Result{}, nil
}
`
