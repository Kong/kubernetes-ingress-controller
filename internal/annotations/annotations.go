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
	"fmt"
	"strings"

	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClassMatching int

const (
	IgnoreClassMatch       ClassMatching = iota
	ExactOrEmptyClassMatch ClassMatching = iota
	ExactClassMatch        ClassMatching = iota
)

const (
	IngressClassKey                  = "kubernetes.io/ingress.class"
	KnativeIngressClassKey           = "networking.knative.dev/ingress-class"
	KnativeIngressClassDeprecatedKey = "networking.knative.dev/ingress.class"

	AnnotationPrefix = "konghq.com"

	ConfigurationKey     = "/override"
	PluginsKey           = "/plugins"
	ProtocolKey          = "/protocol"
	ProtocolsKey         = "/protocols"
	ClientCertKey        = "/client-cert"
	StripPathKey         = "/strip-path"
	PathKey              = "/path"
	HTTPSRedirectCodeKey = "/https-redirect-status-code"
	PreserveHostKey      = "/preserve-host"
	RegexPriorityKey     = "/regex-priority"
	HostHeaderKey        = "/host-header"
	MethodsKey           = "/methods"
	SNIsKey              = "/snis"
	RequestBuffering     = "/request-buffering"
	ResponseBuffering    = "/response-buffering"
	HostAliasesKey       = "/host-aliases"
	RegexPrefixKey       = "/regex-prefix"
	ConnectTimeoutKey    = "/connect-timeout"
	WriteTimeoutKey      = "/write-timeout"
	ReadTimeoutKey       = "/read-timeout"
	RetriesKey           = "/retries"
	HeadersKey           = "/headers"
	PathHandlingKey      = "/path-handling"
	KonnectServiceKey    = "/konnect-service"

	// GatewayClassUnmanagedAnnotationSuffix is an annotation used on a Gateway resource to
	// indicate that the GatewayClass should be reconciled according to unmanaged
	// mode.
	//
	// NOTE: it's currently required that this annotation be present on all GatewayClass
	// resources: "unmanaged" mode is the only supported mode at this time.
	GatewayClassUnmanagedAnnotationSuffix = "gatewayclass-unmanaged"

	// DefaultIngressClass defines the default class used
	// by Kong's ingress controller.
	DefaultIngressClass = "kong"

	// GatewayClassUnmanagedAnnotationPlaceholder is intended to be used as placeholder value for the
	// GatewayClassUnmanagedAnnotation annotation.
	GatewayClassUnmanagedAnnotationValuePlaceholder = "true"
)

// GatewayClassUnmanagedAnnotation is the complete annotations for unmanaged mode made by the konhq.com prefix
// followed by the gatewayclass-unmanaged GatewayClass suffix.
var GatewayClassUnmanagedAnnotation = fmt.Sprintf("%s/%s", AnnotationPrefix, GatewayClassUnmanagedAnnotationSuffix)

func validIngress(ingressAnnotationValue, ingressClass string, handling ClassMatching) bool {
	switch handling {
	case IgnoreClassMatch:
		// class is not considered at all. any value, even a mismatch, is valid
		return true
	case ExactOrEmptyClassMatch:
		// aka lazy. exact match desired, but empty permitted
		return ingressAnnotationValue == "" || ingressAnnotationValue == ingressClass
	case ExactClassMatch:
		// what it says on the tin
		// this may be another place we want to return a warning, since an empty-class resource will never be valid
		return ingressAnnotationValue == ingressClass
	default:
		panic("invalid ingress class handling option received")
	}
}

// IngressClassValidatorFuncFromObjectMeta returns a function which
// can validate if an ObjectMeta belongs to an the ingressClass or not.
func IngressClassValidatorFuncFromObjectMeta(
	ingressClass string,
) func(obj *metav1.ObjectMeta, annotation string, handling ClassMatching) bool {
	return func(obj *metav1.ObjectMeta, annotation string, handling ClassMatching) bool {
		class := obj.GetAnnotations()[annotation]
		return validIngress(class, ingressClass, handling)
	}
}

func IngressClassValidatorFuncFromV1Ingress(
	ingressClass string,
) func(ingress *netv1.Ingress, handling ClassMatching) bool {
	return func(ingress *netv1.Ingress, handling ClassMatching) bool {
		class := ingress.Spec.IngressClassName
		className := ""
		if class != nil {
			className = *class
		}
		return validIngress(className, ingressClass, handling)
	}
}

func pluginsFromAnnotations(anns map[string]string) string {
	return anns[AnnotationPrefix+PluginsKey]
}

// ExtractKongPluginsFromAnnotations extracts information about Kong
// Plugins configured using konghq.com/plugins annotation.
// This returns a list of KongPlugin resource names that should be applied.
func ExtractKongPluginsFromAnnotations(anns map[string]string) []string {
	var kongPluginCRs []string
	v := pluginsFromAnnotations(anns)
	if v == "" {
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
// information about the configuration to use in Routes, Services and Upstreams.
func ExtractConfigurationName(anns map[string]string) string {
	return anns[AnnotationPrefix+ConfigurationKey]
}

// ExtractProtocolName extracts the protocol supplied in the annotation.
func ExtractProtocolName(anns map[string]string) string {
	return anns[AnnotationPrefix+ProtocolKey]
}

// ExtractProtocolNames extracts the protocols supplied in the annotation.
func ExtractProtocolNames(anns map[string]string) []string {
	val := anns[AnnotationPrefix+ProtocolsKey]
	return strings.Split(val, ",")
}

// ExtractClientCertificate extracts the secret name containing the
// client-certificate to use.
func ExtractClientCertificate(anns map[string]string) string {
	return anns[AnnotationPrefix+ClientCertKey]
}

// ExtractStripPath extracts the strip-path annotations containing the
// the boolean string "true" or "false".
func ExtractStripPath(anns map[string]string) string {
	return anns[AnnotationPrefix+StripPathKey]
}

// ExtractPath extracts the path annotations containing the
// HTTP path.
func ExtractPath(anns map[string]string) string {
	return anns[AnnotationPrefix+PathKey]
}

// ExtractHTTPSRedirectStatusCode extracts the https redirect status
// code annotation value.
func ExtractHTTPSRedirectStatusCode(anns map[string]string) string {
	return anns[AnnotationPrefix+HTTPSRedirectCodeKey]
}

// HasForceSSLRedirectAnnotation returns true if the annotation
// ingress.kubernetes.io/force-ssl-redirect is set to "true" in anns.
func HasForceSSLRedirectAnnotation(anns map[string]string) bool {
	return anns["ingress.kubernetes.io/force-ssl-redirect"] == "true"
}

// ExtractPreserveHost extracts the preserve-host annotation value.
func ExtractPreserveHost(anns map[string]string) string {
	return anns[AnnotationPrefix+PreserveHostKey]
}

// HasServiceUpstreamAnnotation returns true if the annotation
// ingress.kubernetes.io/service-upstream is set to "true" in anns.
func HasServiceUpstreamAnnotation(anns map[string]string) bool {
	return anns["ingress.kubernetes.io/service-upstream"] == "true"
}

// ExtractRegexPriority extracts the regex-priority annotation value.
func ExtractRegexPriority(anns map[string]string) string {
	return anns[AnnotationPrefix+RegexPriorityKey]
}

// ExtractHostHeader extracts the host-header annotation value.
func ExtractHostHeader(anns map[string]string) string {
	return anns[AnnotationPrefix+HostHeaderKey]
}

// ExtractMethods extracts the methods annotation value.
func ExtractMethods(anns map[string]string) []string {
	val := anns[AnnotationPrefix+MethodsKey]
	if val == "" {
		return nil
	}
	return strings.Split(val, ",")
}

// ExtractSNIs extracts the route SNI match criteria annotation value.
func ExtractSNIs(anns map[string]string) ([]string, bool) {
	val, exists := anns[AnnotationPrefix+SNIsKey]
	if val == "" {
		return nil, exists
	}
	return strings.Split(val, ","), exists
}

// ExtractRequestBuffering extracts the boolean annotation indicating
// whether or not a route should buffer requests.
func ExtractRequestBuffering(anns map[string]string) (string, bool) {
	s, ok := anns[AnnotationPrefix+RequestBuffering]
	return s, ok
}

// ExtractResponseBuffering extracts the boolean annotation indicating
// whether or not a route should buffer responses.
func ExtractResponseBuffering(anns map[string]string) (string, bool) {
	s, ok := anns[AnnotationPrefix+ResponseBuffering]
	return s, ok
}

// ExtractHostAliases extracts the host-aliases annotation value.
func ExtractHostAliases(anns map[string]string) ([]string, bool) {
	val, exists := anns[AnnotationPrefix+HostAliasesKey]
	if !exists {
		return nil, false
	}
	if val == "" {
		return nil, false
	}
	return strings.Split(val, ","), true
}

// ExtractConnectTimeout extracts the connection timeout annotation value.
func ExtractConnectTimeout(anns map[string]string) (string, bool) {
	val, exists := anns[AnnotationPrefix+ConnectTimeoutKey]
	if !exists {
		return "", false
	}
	return val, true
}

// ExtractWriteTimeout extracts the write timeout annotation value.
func ExtractWriteTimeout(anns map[string]string) (string, bool) {
	val, exists := anns[AnnotationPrefix+WriteTimeoutKey]
	if !exists {
		return "", false
	}
	return val, true
}

// ExtractReadTimeout extracts the read timeout annotation value.
func ExtractReadTimeout(anns map[string]string) (string, bool) {
	val, exists := anns[AnnotationPrefix+ReadTimeoutKey]
	if !exists {
		return "", false
	}
	return val, true
}

// ExtractRetries extracts the retries annotation value.
func ExtractRetries(anns map[string]string) (string, bool) {
	val, exists := anns[AnnotationPrefix+RetriesKey]
	if !exists {
		return "", false
	}
	return val, true
}

// ExtractHeaders extracts the parsed headers annotations values. It returns a map of header names to slices of values.
func ExtractHeaders(anns map[string]string) (map[string][]string, bool) {
	headers := make(map[string][]string)
	prefix := AnnotationPrefix + HeadersKey + "."
	for key, val := range anns {
		if strings.HasPrefix(key, prefix) {
			header := strings.TrimPrefix(key, prefix)
			if len(header) == 0 || len(val) == 0 {
				continue
			}
			headers[header] = strings.Split(val, ",")
		}
	}
	if len(headers) == 0 {
		return headers, false
	}
	return headers, true
}

// ExtractPathHandling extracts the path handling annotation value.
func ExtractPathHandling(anns map[string]string) (string, bool) {
	val, exists := anns[AnnotationPrefix+PathHandlingKey]
	if !exists {
		return "", false
	}
	return val, true
}

// ExtractUnmanagedGatewayClassMode extracts the value of the unmanaged gateway
// mode annotation.
func ExtractUnmanagedGatewayClassMode(anns map[string]string) string {
	if anns == nil {
		return ""
	}
	return anns[GatewayClassUnmanagedAnnotation]
}

// UpdateUnmanagedAnnotation updates the value of the annotation konghq.com/gatewayclass-unmanaged.
func UpdateUnmanagedAnnotation(anns map[string]string, annotationValue string) {
	anns[GatewayClassUnmanagedAnnotation] = annotationValue
}

func ExtractKonnectService(anns map[string]string) string {
	return anns[AnnotationPrefix+KonnectServiceKey]
}
