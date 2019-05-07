/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package annotations

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ingressClassKey = "kubernetes.io/ingress.class"

	pluginsAnnotationKey = "plugins.konghq.com"

	configurationAnnotationKey = "configuration.konghq.com"

	// DefaultIngressClass defines the default class used
	// by Kong's ingres controller.
	DefaultIngressClass = "kong"
)

func validIngress(ingressAnnotationValue, ingressClass string) bool {
	// we have 2 valid combinations
	// 1 - ingress with default class | blank annotation on ingress
	// 2 - ingress with specific class | same annotation on ingress
	//
	// and 2 invalid combinations
	// 3 - ingress with default class | fixed annotation on ingress
	// 4 - ingress with specific class | different annotation on ingress
	if ingressAnnotationValue == "" && ingressClass == DefaultIngressClass {
		return true
	}
	return ingressAnnotationValue == ingressClass
}

// IngressClassValidatorFunc returns a function which can valid if an object
// belongs to an the ingressClass or not.
func IngressClassValidatorFunc(
	ingressClass string) func(obj metav1.Object) bool {

	return func(obj metav1.Object) bool {
		ingress := obj.GetAnnotations()[ingressClassKey]
		return validIngress(ingress, ingressClass)
	}
}

// IngressClassValidatorFuncFromObjectMeta returns a function which
// can valid if an ObjectMeta
// belongs to an the ingressClass or not.
func IngressClassValidatorFuncFromObjectMeta(
	ingressClass string) func(obj *metav1.ObjectMeta) bool {

	return func(obj *metav1.ObjectMeta) bool {
		ingress := obj.GetAnnotations()[ingressClassKey]
		return validIngress(ingress, ingressClass)
	}
}

// ExtractKongPluginsFromAnnotations extracts information about Kong
// Plugins configured using plugins.konghq.com annotation.
// This returns a list of KongPlugin resource names that should be applied.
func ExtractKongPluginsFromAnnotations(anns map[string]string) []string {
	var kongPluginCRs []string
	v, ok := anns[pluginsAnnotationKey]
	if !ok {
		return kongPluginCRs
	}
	for _, kongPlugin := range strings.Split(v, ",") {
		s := strings.TrimSpace(kongPlugin)
		if s != "" {
			kongPluginCRs = append(kongPluginCRs, s)
		}
	}
	return kongPluginCRs
}

// ExtractConfigurationName extracts the name of the KongIngress object that holds
// information about the configuration to use in Routes, Services and Upstreams
func ExtractConfigurationName(anns map[string]string) string {
	return anns[configurationAnnotationKey]
}
