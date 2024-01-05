package controllers

import (
	"github.com/go-logr/logr"
	"github.com/samber/mo"
	k8stypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Reconciler interface {
	SetupWithManager(ctrl.Manager) error
	SetLogger(logr.Logger)
}

type OptionalNamespacedName = mo.Option[k8stypes.NamespacedName]
