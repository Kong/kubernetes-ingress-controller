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

	deprecatedAnnotationPrefix = "configuration.konghq.com"
	annotationPrefix           = "konghq.com"

	deprecatedPluginsKey       = "plugins.konghq.com"
	deprecatedConfigurationKey = deprecatedAnnotationPrefix

	configurationKey     = "/override"
	pluginsKey           = "/plugins"
	protocolKey          = "/protocol"
	protocolsKey         = "/protocols"
	clientCertKey        = "/client-cert"
	stripPathKey         = "/strip-path"
	pathKey              = "/path"
	httpsRedirectCodeKey = "/https-redirect-status-code"
	preserveHostKey      = "/preserve-host"
	regexPriorityKey     = "/regex-priority"
	hostHeaderKey        = "/host-header"
	methodsKey           = "/methods"

	// DefaultIngressClass defines the default class used
	// by Kong's ingress controller.
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

// IngressClassValidatorFunc returns a function which can validate if an Object
// belongs to an the ingressClass or not.
func IngressClassValidatorFunc(
	ingressClass string) func(obj metav1.Object) bool {

	return func(obj metav1.Object) bool {
		ingress := obj.GetAnnotations()[ingressClassKey]
		return validIngress(ingress, ingressClass)
	}
}

// IngressClassValidatorFuncFromObjectMeta returns a function which
// can validate if an ObjectMeta belongs to an the ingressClass or not.
func IngressClassValidatorFuncFromObjectMeta(
	ingressClass string) func(obj *metav1.ObjectMeta) bool {

	return func(obj *metav1.ObjectMeta) bool {
		ingress := obj.GetAnnotations()[ingressClassKey]
		return validIngress(ingress, ingressClass)
	}
}

// valueFromAnnotation returns the value of an annotation with key.
// key is without the annotation prefix of configuration.konghq.com or
// konghq.com.
// It first looks up key under the konghq.com group and if one doesn't
// exist then it looks up the configuration.konghq.com  annotation group.
func valueFromAnnotation(key string, anns map[string]string) string {
	value, exists := anns[annotationPrefix+key]
	if exists {
		return value
	}
	return anns[deprecatedAnnotationPrefix+key]
}

func pluginsFromAnnotations(anns map[string]string) (string, bool) {
	value, exists := anns[annotationPrefix+pluginsKey]
	if exists {
		return value, exists
	}
	value, exists = anns[deprecatedPluginsKey]
	return value, exists
}

// ExtractKongPluginsFromAnnotations extracts information about Kong
// Plugins configured using plugins.konghq.com annotation.
// This returns a list of KongPlugin resource names that should be applied.
func ExtractKongPluginsFromAnnotations(anns map[string]string) []string {
	var kongPluginCRs []string
	v, ok := pluginsFromAnnotations(anns)
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
	value, exists := anns[annotationPrefix+configurationKey]
	if exists {
		return value
	}
	return anns[deprecatedConfigurationKey]
}

// ExtractProtocolName extracts the protocol supplied in the annotation
func ExtractProtocolName(anns map[string]string) string {
	return valueFromAnnotation(protocolKey, anns)
}

// ExtractProtocolNames extracts the protocols supplied in the annotation
func ExtractProtocolNames(anns map[string]string) []string {
	val := valueFromAnnotation(protocolsKey, anns)
	return strings.Split(val, ",")
}

// ExtractClientCertificate extracts the secret name containing the
// client-certificate to use.
func ExtractClientCertificate(anns map[string]string) string {
	return valueFromAnnotation(clientCertKey, anns)
}

// ExtractStripPath extracts the strip-path annotations containing the
// the boolean string "true" or "false".
func ExtractStripPath(anns map[string]string) string {
	return valueFromAnnotation(stripPathKey, anns)
}

// ExtractPath extracts the path annotations containing the
// HTTP path.
func ExtractPath(anns map[string]string) string {
	return valueFromAnnotation(pathKey, anns)
}

// ExtractHTTPSRedirectStatusCode extracts the https redirect status
// code annotation value.
func ExtractHTTPSRedirectStatusCode(anns map[string]string) string {
	return valueFromAnnotation(httpsRedirectCodeKey, anns)
}

// ExtractPreserveHost extracts the preserve-host annotation value.
func ExtractPreserveHost(anns map[string]string) string {
	return valueFromAnnotation(preserveHostKey, anns)
}

// HasServiceUpstreamAnnotation returns true if the annotation
// ingress.kubernetes.io/service-upstream is set to "true" in anns.
func HasServiceUpstreamAnnotation(anns map[string]string) bool {
	return anns["ingress.kubernetes.io/service-upstream"] == "true"
}

// ExtractRegexPriority extracts the host-header annotation value.
func ExtractRegexPriority(anns map[string]string) string {
	return valueFromAnnotation(regexPriorityKey, anns)
}

// ExtractHostHeader extracts the regex-priority annotation value.
func ExtractHostHeader(anns map[string]string) string {
	return valueFromAnnotation(hostHeaderKey, anns)
}

// ExtractMethods extracts the methods annotation value.
func ExtractMethods(anns map[string]string) []string {
	val := valueFromAnnotation(methodsKey, anns)
	if val == "" {
		return []string{}
	}
	return strings.Split(val, ",")
}
