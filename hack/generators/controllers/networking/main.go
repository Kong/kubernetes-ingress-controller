package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
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

	kongv1          = "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1     = "github.com/kong/kubernetes-ingress-controller/v2/api/configuration/v1beta1"
	knativev1alpha1 = "knative.dev/networking/pkg/apis/networking/v1alpha1"
)

// inputControllersNeeded is a list of the supported Types for the
// Kong Kubernetes Ingress Controller. If you need to add a new type
// for support, add it here and a new controller will be generated
// when you run `make controllers`.
var inputControllersNeeded = &typesNeeded{
	typeNeeded{
		PackageImportAlias:                "corev1",
		PackageAlias:                      "CoreV1",
		Package:                           corev1,
		Type:                              "Service",
		Plural:                            "services",
		URL:                               "\"\"",
		CacheType:                         "Service",
		AcceptsIngressClassNameAnnotation: false,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		PackageImportAlias:                "corev1",
		PackageAlias:                      "CoreV1",
		Package:                           corev1,
		Type:                              "Endpoints",
		Plural:                            "endpoints",
		URL:                               "\"\"",
		CacheType:                         "Endpoint",
		AcceptsIngressClassNameAnnotation: false,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"list", "watch"},
	},
	typeNeeded{
		PackageImportAlias:                "corev1",
		PackageAlias:                      "CoreV1",
		Package:                           corev1,
		Type:                              "Secret",
		Plural:                            "secrets",
		URL:                               "\"\"",
		CacheType:                         "Secret",
		AcceptsIngressClassNameAnnotation: false,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"list", "watch"},
	},
	typeNeeded{
		PackageImportAlias:                "netv1",
		PackageAlias:                      "NetV1",
		Package:                           netv1,
		Type:                              "Ingress",
		Plural:                            "ingresses",
		URL:                               "networking.k8s.io",
		CacheType:                         "IngressV1",
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       true,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		PackageImportAlias:                "netv1beta1",
		PackageAlias:                      "NetV1Beta1",
		Package:                           netv1beta1,
		Type:                              "Ingress",
		Plural:                            "ingresses",
		URL:                               "networking.k8s.io",
		CacheType:                         "IngressV1beta1",
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       true,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		PackageImportAlias:                "extv1beta1",
		PackageAlias:                      "ExtV1Beta1",
		Package:                           extv1beta1,
		Type:                              "Ingress",
		Plural:                            "ingresses",
		URL:                               "extensions",
		CacheType:                         "IngressV1beta1",
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       true,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		PackageImportAlias:                "kongv1",
		PackageAlias:                      "KongV1",
		Package:                           kongv1,
		Type:                              "KongIngress",
		Plural:                            "kongingresses",
		URL:                               "configuration.konghq.com",
		CacheType:                         "KongIngress",
		AcceptsIngressClassNameAnnotation: false,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		PackageImportAlias:                "kongv1",
		PackageAlias:                      "KongV1",
		Package:                           kongv1,
		Type:                              "KongPlugin",
		Plural:                            "kongplugins",
		URL:                               "configuration.konghq.com",
		CacheType:                         "Plugin",
		AcceptsIngressClassNameAnnotation: false,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		PackageImportAlias:                "kongv1",
		PackageAlias:                      "KongV1",
		Package:                           kongv1,
		Type:                              "KongClusterPlugin",
		Plural:                            "kongclusterplugins",
		URL:                               "configuration.konghq.com",
		CacheType:                         "ClusterPlugin",
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		PackageImportAlias:                "kongv1",
		PackageAlias:                      "KongV1",
		Package:                           kongv1,
		Type:                              "KongConsumer",
		Plural:                            "kongconsumers",
		URL:                               "configuration.konghq.com",
		CacheType:                         "Consumer",
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		PackageImportAlias:                "kongv1beta1",
		PackageAlias:                      "KongV1Beta1",
		Package:                           kongv1beta1,
		Type:                              "TCPIngress",
		Plural:                            "tcpingresses",
		URL:                               "configuration.konghq.com",
		CacheType:                         "TCPIngress",
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		PackageImportAlias:                "kongv1beta1",
		PackageAlias:                      "KongV1Beta1",
		Package:                           kongv1beta1,
		Type:                              "UDPIngress",
		Plural:                            "udpingresses",
		URL:                               "configuration.konghq.com",
		CacheType:                         "UDPIngress",
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
	typeNeeded{
		PackageImportAlias:                "knativev1alpha1",
		PackageAlias:                      "Knativev1alpha1",
		Package:                           knativev1alpha1,
		Type:                              "Ingress",
		Plural:                            "ingresses",
		URL:                               "networking.internal.knative.dev",
		CacheType:                         "KnativeIngress",
		AcceptsIngressClassNameAnnotation: true,
		AcceptsIngressClassNameSpec:       false,
		RBACVerbs:                         []string{"get", "list", "watch"},
	},
}

var inputRBACPermissionsNeeded = &rbacsNeeded{
	rbacNeeded{
		Plural:    "nodes",
		URL:       `""`,
		RBACVerbs: []string{"list", "watch"},
	},
	rbacNeeded{
		Plural:    "pods",
		URL:       `""`,
		RBACVerbs: []string{"get", "list", "watch"},
	},
	rbacNeeded{
		Plural:    "events",
		URL:       `""`,
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
// controller for, only permissions
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

	return os.WriteFile(outputFile, contents.Bytes(), 0600)
}

type typeNeeded struct {
	PackageImportAlias string
	PackageAlias       string
	Package            string
	Type               string
	Plural             string
	URL                string
	CacheType          string
	RBACVerbs          []string

	// AcceptsIngressClassNameAnnotation indicates that the object accepts (and the controller will listen to)
	// the "kubernetes.io/ingress.class" annotation to decide whether or not the object is supported.
	AcceptsIngressClassNameAnnotation bool

	// AcceptsIngressClassNameSpec indicates the the object indicates the ingress.class that should support it via
	// an attribute in its specification named .IngressClassName
	AcceptsIngressClassNameSpec bool
}

func (t *typeNeeded) generate(contents *bytes.Buffer) error {
	tmpl, err := template.New("controller").Funcs(sprig.TxtFuncMap()).Parse(controllerTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(contents, t)
}

// rbacNeeded represents a resource that we only require RBAC permissions for
type rbacNeeded struct {
	Plural    string
	URL       string
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
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	knativev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/ctrlutils"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/proxy"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)
`

var rbacTemplate = `
// -----------------------------------------------------------------------------
// API Group {{.URL}} resource {{.Plural}}
// -----------------------------------------------------------------------------

//+kubebuilder:rbac:groups={{.URL}},resources={{.Plural}},verbs={{ .RBACVerbs | join ";" }}
`

var controllerTemplate = `
// -----------------------------------------------------------------------------
// {{.PackageAlias}} {{.Type}}
// -----------------------------------------------------------------------------

// {{.PackageAlias}}{{.Type}} reconciles {{.Type}} resources
type {{.PackageAlias}}{{.Type}}Reconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
	Proxy  proxy.Proxy
{{- if or .AcceptsIngressClassNameSpec .AcceptsIngressClassNameAnnotation}}

	IngressClassName string
{{- end}}
}

// SetupWithManager sets up the controller with the Manager.
func (r *{{.PackageAlias}}{{.Type}}Reconciler) SetupWithManager(mgr ctrl.Manager) error {
{{- if .AcceptsIngressClassNameAnnotation}}
	preds := ctrlutils.GeneratePredicateFuncsForIngressClassFilter(r.IngressClassName, {{.AcceptsIngressClassNameSpec}}, true)
	return ctrl.NewControllerManagedBy(mgr).For(&{{.PackageImportAlias}}.{{.Type}}{}, builder.WithPredicates(preds)).Complete(r)
{{- else}}
	return ctrl.NewControllerManagedBy(mgr).For(&{{.PackageImportAlias}}.{{.Type}}{}).Complete(r)
{{- end}}
}

//+kubebuilder:rbac:groups={{.URL}},resources={{.Plural}},verbs={{ .RBACVerbs | join ";" }}
//+kubebuilder:rbac:groups={{.URL}},resources={{.Plural}}/status,verbs=get;update;patch

// Reconcile processes the watched objects
func (r *{{.PackageAlias}}{{.Type}}Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("{{.PackageAlias}}{{.Type}}", req.NamespacedName)

	// get the relevant object
	obj := new({{.PackageImportAlias}}.{{.Type}})
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		obj.Namespace = req.Namespace
		obj.Name = req.Name
		objectExistsInCache, err := r.Proxy.ObjectExists(obj)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			log.V(util.DebugLevel).Info("deleted {{.Type}} object remains in proxy cache, removing", "namespace", req.Namespace, "name", req.Name)
			if err := r.Proxy.DeleteObject(obj); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.V(util.DebugLevel).Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.V(util.DebugLevel).Info("resource is being deleted, its configuration will be removed", "type", "{{.Type}}", "namespace", req.Namespace, "name", req.Name)
		objectExistsInCache, err := r.Proxy.ObjectExists(obj)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			if err := r.Proxy.DeleteObject(obj); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, nil
	}
{{if .AcceptsIngressClassNameAnnotation}}
	// if the object is not configured with our ingress.class, then we need to ensure it's removed from the cache
	if !ctrlutils.MatchesIngressClassName(obj, r.IngressClassName) {
		log.V(util.DebugLevel).Info("object missing ingress class, ensuring it's removed from configuration", "namespace", req.Namespace, "name", req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
{{end}}
	// update the kong Admin API with the changes
	if err := r.Proxy.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
`
