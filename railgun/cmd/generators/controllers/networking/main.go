package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

// -----------------------------------------------------------------------------
// Main
// -----------------------------------------------------------------------------

const outputFile = "controllers/configuration/zz_generated_controllers.go"

// inputControllersNeeded is a list of the supported Types for the
// Kong Kubernetes Ingress Controller. If you need to add a new type
// for support, add it here and a new controller will be generated
// when you run `make controllers`.
var inputControllersNeeded = &typesNeeded{
	typeNeeded{
		PackageImportAlias: "netv1",
		PackageAlias:       "NetV1",
		Package:            "k8s.io/api/networking/v1",
		Type:               "Ingress",
		Plural:             "ingresses",
		URL:                "networking.k8s.io",
	},
	typeNeeded{
		PackageImportAlias: "netv1beta1",
		PackageAlias:       "NetV1Beta1",
		Package:            "k8s.io/api/networking/v1beta1",
		Type:               "Ingress",
		Plural:             "ingresses",
		URL:                "networking.k8s.io",
	},
	typeNeeded{
		PackageImportAlias: "extv1beta1",
		PackageAlias:       "ExtV1Beta1",
		Package:            "k8s.io/api/extensions/v1beta1",
		Type:               "Ingress",
		Plural:             "ingresses",
		URL:                "apiextensions.k8s.io",
	},
	typeNeeded{
		PackageImportAlias: "kongv1",
		PackageAlias:       "KongV1",
		Package:            "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1",
		Type:               "KongIngress",
		Plural:             "kongingresses",
		URL:                "configuration.konghq.com",
	},
	typeNeeded{
		PackageImportAlias: "kongv1",
		PackageAlias:       "KongV1",
		Package:            "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1",
		Type:               "KongPlugin",
		Plural:             "kongplugins",
		URL:                "configuration.konghq.com",
	},
	typeNeeded{
		PackageImportAlias: "kongv1",
		PackageAlias:       "KongV1",
		Package:            "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1",
		Type:               "KongClusterPlugin",
		Plural:             "kongclusterplugins",
		URL:                "configuration.konghq.com",
	},
	typeNeeded{
		PackageImportAlias: "kongv1",
		PackageAlias:       "KongV1",
		Package:            "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1",
		Type:               "KongConsumer",
		Plural:             "kongconsumers",
		URL:                "configuration.konghq.com",
	},
	typeNeeded{
		PackageImportAlias: "kongv1alpha1",
		PackageAlias:       "KongV1",
		Package:            "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1",
		Type:               "UDPIngress",
		Plural:             "udpingresses",
		URL:                "configuration.konghq.com",
	},
}

func main() {
	if err := inputControllersNeeded.generate(); err != nil {
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

	boilerPlate, err := ioutil.ReadFile("hack/boilerplate.go.txt")
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

// generate generates a controller/input/<controller>.go Kubernetes controller
// for every supported type populated in the list.
func (types typesNeeded) generate() error {
	contents, err := header()
	if err != nil {
		return err
	}

	for _, t := range types {
		if err := t.generate(contents); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(outputFile, contents.Bytes(), 0644)
}

type typeNeeded struct {
	PackageImportAlias string
	PackageAlias       string
	Package            string
	Type               string
	Plural             string
	URL                string
}

func (t *typeNeeded) generate(contents *bytes.Buffer) error {
	tmpl, err := template.New("controller").Parse(controllerTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(contents, t)
}

// -----------------------------------------------------------------------------
// Templates
// -----------------------------------------------------------------------------

// TODO: we should switch to autogenerating the imports, for now they're added manually.
var headerTemplate = `
// Code generated by Kong; DO NOT EDIT.

package configuration

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrl "sigs.k8s.io/controller-runtime"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	kongv1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
)
`

var controllerTemplate = `
// -----------------------------------------------------------------------------
// {{.PackageAlias}} {{.Type}}
// -----------------------------------------------------------------------------

// {{.PackageAlias}}{{.Type}} reconciles a Ingress object
type {{.PackageAlias}}{{.Type}}Reconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	TargetNamespacedName *types.NamespacedName
}

// SetupWithManager sets up the controller with the Manager.
func (r *{{.PackageAlias}}{{.Type}}Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&{{.PackageImportAlias}}.{{.Type}}{}).Complete(r)
}

//+kubebuilder:rbac:groups={{.URL}},resources={{.Plural}},verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups={{.URL}},resources={{.Plural}}/status,verbs=get;update;patch
//+kubebuilder:rbac:groups={{.URL}},resources={{.Plural}}/finalizers,verbs=update

// Reconcile processes the watched objects
func (r *{{.PackageAlias}}{{.Type}}Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("{{.PackageAlias}}{{.Type}}", req.NamespacedName)

	obj := new({{.PackageImportAlias}}.{{.Type}})
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.Info("resource is being deleted, its configuration will be removed", "type", "{{.Type}}", "namespace", req.Namespace, "name", req.Name)
		return cleanupObj(ctx, r.Client, log, *r.TargetNamespacedName, req.NamespacedName, obj)
	}

	return storeIngressObj(ctx, r.Client, log, *r.TargetNamespacedName, req.NamespacedName, obj)
}
`
