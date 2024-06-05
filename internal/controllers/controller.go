package controllers

import (
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Reconciler interface {
	SetupWithManager(ctrl.Manager) error
	SetLogger(logr.Logger)
}
