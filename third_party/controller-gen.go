//go:build third_party
// +build third_party

package third_party

import (
	_ "sigs.k8s.io/controller-tools/pkg/crd"
	_ "sigs.k8s.io/controller-tools/pkg/rbac"
	_ "sigs.k8s.io/controller-tools/pkg/version"
	_ "sigs.k8s.io/controller-tools/pkg/webhook"
)

// For some reason controller-tools is the only package that complains when
// imported. Proabably because it's a main package but not sure why other main
// package modules don't complain like this.
//
// Hence we resort to manually importing several of its subpackages to
// make sigs.k8s.io/controller-tools show up in go.mod with a proper version
// so that it can be managed via go tools when udating tools' dependencies.
//
// go: finding module for package sigs.k8s.io/controller-tools
// github.com/kong/gateway-operator/third_party imports
// sigs.k8s.io/controller-tools: module sigs.k8s.io/controller-tools@latest found (v0.9.2), but does not contain package sigs.k8s.io/controller-tools

//go:generate go install -modfile go.mod sigs.k8s.io/controller-tools/cmd/controller-gen
