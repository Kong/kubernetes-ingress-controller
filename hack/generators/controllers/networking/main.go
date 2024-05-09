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

	corev1      = "k8s.io/api/core/v1"
	discoveryv1 = "k8s.io/api/discovery/v1"
	netv1       = "k8s.io/api/networking/v1"

	kongv1       = "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1beta1  = "github.com/kong/kubernetes-ingress-controller/v3/api/configuration/v1beta1"
	kongv1alpha1 = "github.com/kong/kubernetes-ingress-controller/v3/api/configuration/v1alpha1"

	incubatorv1alpha1 = "github.com/kong/kubernetes-ingress-controller/v3/api/incubator/v1alpha1"
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
		Group:                             "discovery.k8s.io",
		Version:                           "v1",
		Kind:                              "EndpointSlice",
		PackageImportAlias:                "discoveryv1",
		PackageAlias:                      "DiscoveryV1",
		Package:                           discoveryv1,
		Plural:                            "endpointslices",
		CacheType:                         "EndpointSlice",
		NeedsStatusPermissions:            false,
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
		ConfigStatusNotificationsEnabled:  true,
		IngressAddressUpdatesEnabled:      true,
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
		Group:                            "configuration.konghq.com",
		Version:                          "v1",
		Kind:                             "KongPlugin",
		PackageImportAlias:               "kongv1",
		PackageAlias:                     "KongV1",
		Package:                          kongv1,
		Plural:                           "kongplugins",
		CacheType:                        "Plugin",
		NeedsStatusPermissions:           true,
		ConfigStatusNotificationsEnabled: false, // TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/4578
		ProgrammedCondition: ProgrammedConditionConfiguration{
			UpdatesEnabled: false, // TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/4578
		},
		AcceptsIngressClassNameAnnotation: false,
		AcceptsIngressClassNameSpec:       false,
		NeedsUpdateReferences:             true,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		Group:                            "configuration.konghq.com",
		Version:                          "v1",
		Kind:                             "KongClusterPlugin",
		PackageImportAlias:               "kongv1",
		PackageAlias:                     "KongV1",
		Package:                          kongv1,
		Plural:                           "kongclusterplugins",
		CacheType:                        "ClusterPlugin",
		NeedsStatusPermissions:           true,
		ConfigStatusNotificationsEnabled: false, // TODO true after https://github.com/Kong/kubernetes-ingress-controller/issues/4578
		ProgrammedCondition: ProgrammedConditionConfiguration{
			UpdatesEnabled: false, // TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/4578
		},
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
		ConfigStatusNotificationsEnabled:  true,
		ProgrammedCondition: ProgrammedConditionConfiguration{
			UpdatesEnabled: true,
		},
	},
	typeNeeded{
		Group:                            "configuration.konghq.com",
		Version:                          "v1beta1",
		Kind:                             "KongConsumerGroup",
		PackageImportAlias:               "kongv1beta1",
		PackageAlias:                     "KongV1Beta1",
		Package:                          kongv1beta1,
		Plural:                           "kongconsumergroups",
		CacheType:                        "ConsumerGroup",
		NeedsStatusPermissions:           true,
		ConfigStatusNotificationsEnabled: true,
		ProgrammedCondition: ProgrammedConditionConfiguration{
			UpdatesEnabled: true,
		},
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
		ConfigStatusNotificationsEnabled:  true,
		IngressAddressUpdatesEnabled:      true,
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
		ConfigStatusNotificationsEnabled:  true,
		IngressAddressUpdatesEnabled:      true,
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
		AcceptsIngressClassNameAnnotation: false,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		Group:                            "incubator.ingress-controller.konghq.com",
		Version:                          "v1alpha1",
		Kind:                             "KongServiceFacade",
		PackageImportAlias:               "incubatorv1alpha1",
		PackageAlias:                     "IncubatorV1Alpha1",
		Package:                          incubatorv1alpha1,
		Plural:                           "kongservicefacades",
		CacheType:                        "KongServiceFacade",
		NeedsStatusPermissions:           true,
		ConfigStatusNotificationsEnabled: true,
		ProgrammedCondition: ProgrammedConditionConfiguration{
			UpdatesEnabled:       true,
			CustomUnknownMessage: "Found no references to this resource in Ingress or similar resources.",
		},
		AcceptsIngressClassNameAnnotation: true,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		Group:                            "configuration.konghq.com",
		Version:                          "v1alpha1",
		Kind:                             "KongVault",
		PackageImportAlias:               "kongv1alpha1",
		PackageAlias:                     "KongV1Alpha1",
		Package:                          kongv1alpha1,
		Plural:                           "kongvaults",
		CacheType:                        "KongVault",
		NeedsStatusPermissions:           true,
		ConfigStatusNotificationsEnabled: true,
		ProgrammedCondition: ProgrammedConditionConfiguration{
			UpdatesEnabled: true,
		},
		AcceptsIngressClassNameAnnotation: true,
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

	// ConfigStatusNotificationsEnabled indicates that the controller should receive updates via the StatusQueue when the
	// configuration status of the resource changes.
	ConfigStatusNotificationsEnabled bool

	// IngressAddressUpdatesEnabled indicates that the controllers should update the address of the ingress for the
	// resource.
	IngressAddressUpdatesEnabled bool

	// ProgrammedCondition contains the configuration for the Programmed condition for the resource.
	ProgrammedCondition ProgrammedConditionConfiguration

	// NeedUpdateReferences is true if we need to update the reference relationships
	// between reconciled object and other objects.
	NeedsUpdateReferences bool
}

type ProgrammedConditionConfiguration struct {
	// UpdatesEnabled indicates that the controllers should manage the Programmed condition for the
	// resource.
	UpdatesEnabled bool

	// CustomUnknownMessage is the message to use for the Programmed condition when the configuration status is Unknown.
	CustomUnknownMessage string
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
	discoveryv1 "k8s.io/api/discovery/v1"
	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8stypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/reference"
	ctrlutils "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/utils"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object/status"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
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
	DataplaneClient controllers.DataPlane
	CacheSyncTimeout time.Duration
{{- if .IngressAddressUpdatesEnabled }}

	DataplaneAddressFinder *dataplane.AddressFinder
{{- end}}
{{- if .ConfigStatusNotificationsEnabled }}
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

var _ controllers.Reconciler = &{{.PackageAlias}}{{.Kind}}Reconciler{}

// SetupWithManager sets up the controller with the Manager.
func (r *{{.PackageAlias}}{{.Kind}}Reconciler) SetupWithManager(mgr ctrl.Manager) error {
    blder := ctrl.NewControllerManagedBy(mgr).
		// set the controller name
		Named("{{.PackageAlias}}{{.Kind}}").
		WithOptions(controller.Options{
			Reconciler: r,
			LogConstructor: func(_ *reconcile.Request) logr.Logger {
				return r.Log
			},
			CacheSyncTimeout: r.CacheSyncTimeout,
	    })

{{- if .ConfigStatusNotificationsEnabled }}
	// if configured, start the status updater controller
	if r.StatusQueue != nil {
		blder.WatchesRawSource(
			source.Channel(
				r.StatusQueue.Subscribe(schema.GroupVersionKind{
					Group:   "{{.Group}}",
					Version: "{{.Version}}",
					Kind:    "{{.Kind}}",
				}),
				&handler.EnqueueRequestForObject{},
			),
		)
	}
{{- end}}
{{- if .AcceptsIngressClassNameAnnotation}}
	if !r.DisableIngressClassLookups {
		blder.Watches(&netv1.IngressClass{},
			handler.EnqueueRequestsFromMapFunc(r.listClassless),
			builder.WithPredicates(predicate.NewPredicateFuncs(ctrlutils.IsDefaultIngressClass)),
		)
	}
	preds := ctrlutils.GeneratePredicateFuncsForIngressClassFilter(r.IngressClassName)
{{- end}}
    return blder.Watches(&{{.PackageImportAlias}}.{{.Kind}}{},
		&handler.EnqueueRequestForObject{},
{{- if .AcceptsIngressClassNameAnnotation}}
		builder.WithPredicates(preds),
{{- end}}
	).
		Complete(r)
}

{{- if .AcceptsIngressClassNameAnnotation}}
// listClassless finds and reconciles all objects without ingress class information
func (r *{{.PackageAlias}}{{.Kind}}Reconciler) listClassless(ctx context.Context, obj client.Object) []reconcile.Request {
	resourceList := &{{.PackageImportAlias}}.{{.Kind}}List{}
	if err := r.Client.List(ctx, resourceList); err != nil {
		r.Log.Error(err, "Failed to list classless {{.Plural}}")
		return nil
	}
	var recs []reconcile.Request
	for i, resource := range resourceList.Items {
		if ctrlutils.IsIngressClassEmpty(&resourceList.Items[i]) {
			recs = append(recs, reconcile.Request{
				NamespacedName: k8stypes.NamespacedName{
					Namespace: resource.Namespace,
					Name:      resource.Name,
				},
			})
		}
	}
	return recs
}
{{- end}}

// SetLogger sets the logger.
func (r *{{.PackageAlias}}{{.Kind}}Reconciler) SetLogger(l logr.Logger) {
	r.Log = l
}

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
	log.V(util.DebugLevel).Info("Reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.V(util.DebugLevel).Info("Resource is being deleted, its configuration will be removed", "type", "{{.Kind}}", "namespace", req.Namespace, "name", req.Name)
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
		if err := r.Get(ctx, k8stypes.NamespacedName{Name: r.IngressClassName}, class); err != nil {
			// we log this without taking action to support legacy configurations that only set ingressClassName or
			// used the class annotation and did not create a corresponding IngressClass. We only need this to determine
			// if the IngressClass is default or to configure default settings, and can assume no/no additional defaults
			// if none exists.
			log.V(util.DebugLevel).Info("Could not retrieve IngressClass", "ingressclass", r.IngressClassName)
		}
	}
	// if the object is not configured with our ingress.class, then we need to ensure it's removed from the cache
	if !ctrlutils.MatchesIngressClass(obj, r.IngressClassName, ctrlutils.IsDefaultIngressClass(class)) {
		log.V(util.DebugLevel).Info("Object missing ingress class, ensuring it's removed from configuration",
		"namespace", req.Namespace, "name", req.Name, "class", r.IngressClassName)
		return ctrl.Result{}, r.DataplaneClient.DeleteObject(obj)
	} else {
		log.V(util.DebugLevel).Info("Object has matching ingress class", "namespace", req.Namespace, "name", req.Name,
		"class", r.IngressClassName)
	}
{{end}}
	// update the kong Admin API with the changes
	if err := r.DataplaneClient.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}

{{- define "updateReferences" }}
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

{{- /* For ProgrammedCondition.UpdatesEnabled we do not update references before status is updated because in case of
       a reference to non-existing object, the status update would never happen. */ -}}
{{- if and .NeedsUpdateReferences (not .ProgrammedCondition.UpdatesEnabled) }}
	{{- template "updateReferences" . }}
{{- end }}

{{- if .ConfigStatusNotificationsEnabled }}
	// if status updates are enabled report the status for the object
	if r.DataplaneClient.AreKubernetesObjectReportsEnabled() {
		{{- if .IngressAddressUpdatesEnabled }}
		log.V(util.DebugLevel).Info("Determining whether data-plane configuration has succeeded", "namespace", req.Namespace, "name", req.Name)

		if  !r.DataplaneClient.KubernetesObjectIsConfigured(obj) {
			log.V(util.DebugLevel).Info("Resource not yet configured in the data-plane", "namespace", req.Namespace, "name", req.Name)
			return ctrl.Result{Requeue: true}, nil // requeue until the object has been properly configured
		}

		log.V(util.DebugLevel).Info("Determining gateway addresses for object status updates", "namespace", req.Namespace, "name", req.Name)
		addrs, err := r.DataplaneAddressFinder.GetLoadBalancerAddresses(ctx)
		if err != nil {
			return ctrl.Result{}, err
		}

		log.V(util.DebugLevel).Info("Found addresses for data-plane updating object status", "namespace", req.Namespace, "name", req.Name)
		updateNeeded, err := ctrlutils.UpdateLoadBalancerIngress(obj, addrs)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to update load balancer address: %w", err)
		}
		{{- end }}

		{{- if .ProgrammedCondition.UpdatesEnabled }}
		log.V(util.DebugLevel).Info("Updating programmed condition status", "namespace", req.Namespace, "name", req.Name)
		configurationStatus := r.DataplaneClient.KubernetesObjectConfigurationStatus(obj)
		conditions, updateNeeded := ctrlutils.EnsureProgrammedCondition(
			configurationStatus, 
			obj.Generation, 
			obj.Status.Conditions,
		{{- if .ProgrammedCondition.CustomUnknownMessage }}
			ctrlutils.WithUnknownMessage("{{ .ProgrammedCondition.CustomUnknownMessage }}"),
		{{- end }}
		)
		obj.Status.Conditions = conditions
		{{- end }}
		if updateNeeded {
			return ctrl.Result{}, r.Status().Update(ctx, obj)
		}
		log.V(util.DebugLevel).Info("Status update not needed", "namespace", req.Namespace, "name", req.Name)
	}
{{- end}}

{{- /* For ProgrammedCondition.UpdatesEnabled we update references after status is updated because otherwise in case of
       a reference to non-existing object, the status update would never happen. */ -}}
{{- if and .NeedsUpdateReferences .ProgrammedCondition.UpdatesEnabled }}
	{{- template "updateReferences" . }}
{{- end }}

	return ctrl.Result{}, nil
}
`
