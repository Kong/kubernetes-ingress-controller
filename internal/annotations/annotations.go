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
	"strconv"
	"strings"

	"github.com/samber/lo"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
)

type ClassMatching int

const (
	IgnoreClassMatch       ClassMatching = iota
	ExactOrEmptyClassMatch ClassMatching = iota
	ExactClassMatch        ClassMatching = iota
)

const (
	IngressClassKey = "kubernetes.io/ingress.class"

	AnnotationPrefix = "konghq.com"

	ConfigurationKey            = "/override"
	PluginsKey                  = "/plugins"
	ProtocolKey                 = "/protocol"
	ProtocolsKey                = "/protocols"
	ClientCertKey               = "/client-cert"
	StripPathKey                = "/strip-path"
	PathKey                     = "/path"
	HTTPSRedirectCodeKey        = "/https-redirect-status-code"
	PreserveHostKey             = "/preserve-host"
	RegexPriorityKey            = "/regex-priority"
	HostHeaderKey               = "/host-header"
	MethodsKey                  = "/methods"
	SNIsKey                     = "/snis"
	RequestBuffering            = "/request-buffering"
	ResponseBuffering           = "/response-buffering"
	HostAliasesKey              = "/host-aliases"
	RegexPrefixKey              = "/regex-prefix"
	ConnectTimeoutKey           = "/connect-timeout"
	WriteTimeoutKey             = "/write-timeout"
	ReadTimeoutKey              = "/read-timeout"
	RetriesKey                  = "/retries"
	HeadersKey                  = "/headers"
	HeadersSeparatorKey         = "/headers-separator"
	PathHandlingKey             = "/path-handling"
	UserTagKey                  = "/tags"
	RewriteURIKey               = "/rewrite"
	TLSVerifyKey                = "/tls-verify"
	TLSVerifyDepthKey           = "/tls-verify-depth"
	CACertificatesSecretsKey    = "/ca-certificates-secret"
	CACertificatesConfigMapsKey = "/ca-certificates-configmap"

	// GatewayClassUnmanagedKey is an annotation used on a Gateway resource to
	// indicate that the GatewayClass should be reconciled according to unmanaged
	// mode.
	//
	// NOTE: it's currently required that this annotation be present on all GatewayClass
	// resources: "unmanaged" mode is the only supported mode at this time.
	GatewayClassUnmanagedKey = "/gatewayclass-unmanaged"

	// GatewayPublishServiceKey is an annotation suffix used to indicate the Service(s) a Gateway's routes are
	// published to.
	GatewayPublishServiceKey = "/publish-service"

	// DefaultIngressClass defines the default class used
	// by Kong's ingress controller.
	DefaultIngressClass = "kong"

	// GatewayClassUnmanagedAnnotationValuePlaceholder is intended to be used as placeholder value for the
	// GatewayClassUnmanagedAnnotation annotation.
	GatewayClassUnmanagedAnnotationValuePlaceholder = "true"
)

// GatewayClassUnmanagedAnnotation is the complete annotations for unmanaged mode made by the konhq.com prefix
// followed by the gatewayclass-unmanaged GatewayClass suffix.
var GatewayClassUnmanagedAnnotation = fmt.Sprintf("%s%s", AnnotationPrefix, GatewayClassUnmanagedKey)

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
	return extractCommaDelimitedStrings(val)
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

func ExtractRegexPrefix(anns map[string]string) string {
	return anns[AnnotationPrefix+RegexPrefixKey]
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
	return extractCommaDelimitedStrings(val, strings.ToUpper)
}

// ExtractSNIs extracts the route SNI match criteria annotation value.
func ExtractSNIs(anns map[string]string) ([]string, bool) {
	val, exists := anns[AnnotationPrefix+SNIsKey]
	return extractCommaDelimitedStrings(val), exists
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
	return extractCommaDelimitedStrings(val), true
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
	const prefix = AnnotationPrefix + HeadersKey + "."
	separator, ok := anns[AnnotationPrefix+HeadersSeparatorKey]
	if !ok {
		separator = ","
	}
	for key, val := range anns {
		if strings.HasPrefix(key, prefix) {
			header := strings.TrimPrefix(key, prefix)
			if len(header) == 0 || len(val) == 0 {
				continue
			}
			headers[header] = lo.Map(strings.Split(val, separator), func(hv string, _ int) string {
				return strings.TrimSpace(hv)
			})
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

// ExtractGatewayPublishService extracts the value of the gateway publish service annotation.
func ExtractGatewayPublishService(anns map[string]string) []string {
	if anns == nil {
		return []string{}
	}
	publish := anns[AnnotationPrefix+GatewayPublishServiceKey]
	return extractCommaDelimitedStrings(publish)
}

// UpdateGatewayPublishService updates the value of the annotation konghq.com/gatewayclass-unmanaged.
func UpdateGatewayPublishService(anns map[string]string, services []string) {
	anns[AnnotationPrefix+GatewayPublishServiceKey] = strings.Join(services, ",")
}

// ExtractUserTags extracts a set of tags from a comma-separated string.
func ExtractUserTags(anns map[string]string) []string {
	val := anns[AnnotationPrefix+UserTagKey]
	return extractCommaDelimitedStrings(val)
}

// ExtractRewriteURI extracts the rewrite annotation value.
func ExtractRewriteURI(anns map[string]string) (string, bool) {
	s, ok := anns[AnnotationPrefix+RewriteURIKey]
	return s, ok
}

// ExtractUpstreamPolicy extracts the upstream policy annotation value.
func ExtractUpstreamPolicy(anns map[string]string) (string, bool) {
	s, ok := anns[kongv1beta1.KongUpstreamPolicyAnnotationKey]
	return s, ok
}

// ExtractTLSVerify extracts the tls-verify annotation value.
func ExtractTLSVerify(anns map[string]string) (value bool, ok bool) {
	s, ok := anns[AnnotationPrefix+TLSVerifyKey]
	if !ok {
		// If the annotation is not present, we consider it not set.
		return false, false
	}
	verify, err := strconv.ParseBool(s)
	if err != nil {
		// If the annotation is present but not a valid boolean string, we consider it not set.
		return false, false
	}
	// If the annotation is present and a valid boolean string, we return the value.
	return verify, true
}

// ExtractTLSVerifyDepth extracts the tls-verify-depth annotation value.
func ExtractTLSVerifyDepth(anns map[string]string) (int, bool) {
	s, ok := anns[AnnotationPrefix+TLSVerifyDepthKey]
	if !ok {
		// If the annotation is not present, we consider it not set.
		return 0, false
	}
	depth, err := strconv.Atoi(s)
	if err != nil {
		// If the annotation is present but not a valid integer string, we consider it not set.
		return 0, false
	}
	// If the annotation is present and a valid integer string, we return the value.
	return depth, true
}

// ExtractCACertificatesFromSecrets extracts the ca-certificates secret names from the annotation.
// It expects a comma-separated list of certificate names.
func ExtractCACertificatesFromSecrets(anns map[string]string) []string {
	s, ok := anns[AnnotationPrefix+CACertificatesSecretsKey]
	if !ok {
		return nil
	}
	return extractCommaDelimitedStrings(s)
}

func ExtractCACertificatesFromConfigMap(anns map[string]string) []string {
	s, ok := anns[AnnotationPrefix+CACertificatesConfigMapsKey]
	if !ok {
		return nil
	}
	return extractCommaDelimitedStrings(s)
}

// extractCommaDelimitedStrings extracts a list of non-empty strings from a comma-separated string.
// It trims spaces from the strings.
// It accepts optional sanitization functions to apply to each string.
func extractCommaDelimitedStrings(s string, sanitizeFns ...func(string) string) []string {
	// If it's an empty string, return nil.
	if strings.TrimSpace(s) == "" {
		return nil
	}

	// Split values by comma.
	values := strings.Split(s, ",")

	// Allocate an output slice with the same capacity as the input slice.
	// This may be a bit more than needed as we'll filter out empty strings later.
	out := make([]string, 0, len(values))

	// Trim and sanitize each value.
	for _, v := range values {
		sanitized := strings.TrimSpace(v)
		if sanitized == "" {
			// Discard empty strings.
			continue
		}

		// Apply optional sanitization functions (e.g. upper-casing).
		for _, sanitizeFn := range sanitizeFns {
			sanitized = sanitizeFn(sanitized)
		}

		// Append to the output slice.
		out = append(out, sanitized)
	}

	return out
}

// SetTLSVerify sets the tls-verify annotation value.
func SetTLSVerify(anns map[string]string, value bool) {
	anns[AnnotationPrefix+TLSVerifyKey] = strconv.FormatBool(value)
}

// SetCACertificates merge the ca-certificates secret names into the already existing CA certificates set via annotation.
func SetCACertificates(anns map[string]string, certificates []string) {
	existingCACerts := anns[AnnotationPrefix+CACertificatesConfigMapsKey]
	if existingCACerts == "" {
		anns[AnnotationPrefix+CACertificatesConfigMapsKey] = strings.Join(certificates, ",")
	} else {
		anns[AnnotationPrefix+CACertificatesConfigMapsKey] = existingCACerts + "," + strings.Join(certificates, ",")
	}
}

// SetHostHeader sets the host-header annotation value.
func SetHostHeader(anns map[string]string, value string) {
	anns[AnnotationPrefix+HostHeaderKey] = value
}

// SetProtocol sets the protocol annotation value.
func SetProtocol(anns map[string]string, value string) {
	anns[AnnotationPrefix+ProtocolKey] = value
}

// SetTLSVerifyDepth sets the tls-verify-depth annotation value.
func SetTLSVerifyDepth(anns map[string]string, depth int) {
	anns[AnnotationPrefix+TLSVerifyDepthKey] = strconv.Itoa(depth)
}
