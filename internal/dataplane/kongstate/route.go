package kongstate

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// Route represents a Kong Route and holds a reference to the Ingress
// rule.
type Route struct {
	kong.Route

	Ingress          util.K8sObjectInfo
	Plugins          []kong.Plugin
	ExpressionRoutes bool
}

var (
	validMethods      = regexp.MustCompile(`\A[A-Z]+$`)
	validPathHandling = regexp.MustCompile(`v\d`)

	// hostnames are complicated. shamelessly cribbed from https://stackoverflow.com/a/18494710
	// TODO if the Kong core adds support for wildcard SNI route match criteria, this should change.
	validSNIs  = regexp.MustCompile(`^([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*)+(\.([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*))*$`)
	validHosts = regexp.MustCompile(`^(\*\.)?([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*)+(\.([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*))*?(\.\*)?$`)
)

// normalizeProtocols prevents users from mismatching grpc/http.
func (r *Route) normalizeProtocols() {
	// skip updating protocols if expression routes enabled.
	if r.ExpressionRoutes {
		return
	}
	protocols := r.Protocols
	var http, grpc bool

	for _, protocol := range protocols {
		if strings.Contains(*protocol, "grpc") {
			grpc = true
		}
		if strings.Contains(*protocol, "http") {
			http = true
		}
		if !util.ValidateProtocol(*protocol) {
			http = true
		}
	}

	if grpc && http {
		r.Protocols = kong.StringSlice("http", "https")
	}

	if grpc {
		// grpc(s) doesn't accept strip_path
		r.StripPath = nil
	}
}

// useSSLProtocol updates the protocol of the route to either https or grpcs, or https and grpcs.
func (r *Route) useSSLProtocol() {
	var http, grpc bool
	var prots []*string

	for _, val := range r.Protocols {
		if strings.Contains(*val, "grpc") {
			grpc = true
		}

		if strings.Contains(*val, "http") {
			http = true
		}
	}

	if grpc {
		prots = append(prots, kong.String("grpcs"))
	}
	if http {
		prots = append(prots, kong.String("https"))
	}

	if !grpc && !http {
		prots = append(prots, kong.String("https"))
	}

	r.Protocols = prots
}

func (r *Route) overrideStripPath(anns map[string]string) {
	if r == nil {
		return
	}

	stripPathValue := annotations.ExtractStripPath(anns)
	if stripPathValue == "" {
		return
	}
	stripPathValue = strings.ToLower(stripPathValue)
	switch stripPathValue {
	case "true":
		r.StripPath = kong.Bool(true)
	case "false":
		r.StripPath = kong.Bool(false)
	default:
		return
	}
}

func (r *Route) overrideProtocols(anns map[string]string) {
	protocols := annotations.ExtractProtocolNames(anns)
	if len(protocols) == 0 {
		return
	}
	var prots []*string //nolint:prealloc
	for _, prot := range protocols {
		if !util.ValidateProtocol(prot) {
			return
		}
		prots = append(prots, kong.String(prot))
	}

	r.Protocols = prots
}

func (r *Route) overrideHTTPSRedirectCode(anns map[string]string) {
	if annotations.HasForceSSLRedirectAnnotation(anns) {
		r.HTTPSRedirectStatusCode = kong.Int(302)
		r.useSSLProtocol()
	}

	code := annotations.ExtractHTTPSRedirectStatusCode(anns)
	if code == "" {
		return
	}
	statusCode, err := strconv.Atoi(code)
	if err != nil {
		return
	}
	if statusCode != 426 &&
		statusCode != 301 &&
		statusCode != 302 &&
		statusCode != 307 &&
		statusCode != 308 {
		return
	}

	r.HTTPSRedirectStatusCode = kong.Int(statusCode)
}

func (r *Route) overridePreserveHost(anns map[string]string) {
	preserveHostValue := annotations.ExtractPreserveHost(anns)
	if preserveHostValue == "" {
		return
	}
	preserveHostValue = strings.ToLower(preserveHostValue)
	switch preserveHostValue {
	case "true":
		r.PreserveHost = kong.Bool(true)
	case "false":
		r.PreserveHost = kong.Bool(false)
	default:
		return
	}
}

func (r *Route) overrideRegexPriority(anns map[string]string) {
	priority := annotations.ExtractRegexPriority(anns)
	if priority == "" {
		return
	}
	regexPriority, err := strconv.Atoi(priority)
	if err != nil {
		return
	}

	r.RegexPriority = kong.Int(regexPriority)
}

func (r *Route) overrideMethods(logger logr.Logger, anns map[string]string) {
	annMethods := annotations.ExtractMethods(anns)
	if len(annMethods) == 0 {
		return
	}
	var methods []*string
	for _, method := range annMethods {
		if validMethods.MatchString(method) {
			methods = append(methods, kong.String(method))
		} else {
			// if any method is invalid (not an uppercase alpha string),
			// discard everything
			logger.Error(nil, "Invalid method", "route_name", r.Name, "method", method)
			return
		}
	}

	r.Methods = methods
}

func (r *Route) overrideSNIs(logger logr.Logger, anns map[string]string) {
	var annSNIs []string
	var exists bool
	annSNIs, exists = annotations.ExtractSNIs(anns)
	// this is not a length check because we specifically want to provide a means
	// to set "no SNI criteria", by providing the annotation with an empty string value
	if !exists {
		return
	}
	var snis []*string
	for _, sni := range annSNIs {
		if validSNIs.MatchString(sni) {
			snis = append(snis, kong.String(sni))
		} else {
			// SNI is not a valid hostname
			logger.Error(nil, "Invalid SNI", "route_name", r.Name, "sni", sni)
			return
		}
	}

	r.SNIs = snis
}

// overrideByAnnotation sets Route protocols via annotation.
func (r *Route) overrideByAnnotation(logger logr.Logger) {
	r.overrideStripPath(r.Ingress.Annotations)
	r.overrideHTTPSRedirectCode(r.Ingress.Annotations)
	r.overridePreserveHost(r.Ingress.Annotations)
	r.overrideRequestBuffering(logger, r.Ingress.Annotations)
	r.overrideResponseBuffering(logger, r.Ingress.Annotations)
	r.overrideProtocols(r.Ingress.Annotations)
	// skip the fields that are not supported when kong is using expression router:
	// `regexPriority`, `methods`, `snis`, `hosts`, `headers`, `pathHandling`,
	if !r.ExpressionRoutes {
		r.overrideRegexPriority(r.Ingress.Annotations)
		r.overrideMethods(logger, r.Ingress.Annotations)
		r.overrideSNIs(logger, r.Ingress.Annotations)
		r.overrideHosts(logger, r.Ingress.Annotations)
		r.overrideHeaders(r.Ingress.Annotations)
		r.overridePathHandling(logger, r.Ingress.Annotations)
	}
}

// override sets Route fields by KongIngress first, then by annotation.
func (r *Route) override(logger logr.Logger) {
	if r == nil {
		return
	}

	r.overrideByAnnotation(logger)
	r.normalizeProtocols()
}

// overrideRequestBuffering ensures defaults for the request_buffering option.
func (r *Route) overrideRequestBuffering(logger logr.Logger, anns map[string]string) {
	annotationValue, ok := annotations.ExtractRequestBuffering(anns)
	if !ok {
		// the annotation is not set, quit
		return
	}

	isEnabled, err := strconv.ParseBool(strings.ToLower(annotationValue))
	if err != nil {
		// the value provided is not a parseable boolean, quit
		logger.Error(err, "Invalid request_buffering value", "kongroute", r.Name)
		return
	}

	r.RequestBuffering = kong.Bool(isEnabled)
}

// overrideResponseBuffering ensures defaults for the response_buffering option.
func (r *Route) overrideResponseBuffering(logger logr.Logger, anns map[string]string) {
	annotationValue, ok := annotations.ExtractResponseBuffering(anns)
	if !ok {
		// the annotation is not set, quit
		return
	}

	isEnabled, err := strconv.ParseBool(strings.ToLower(annotationValue))
	if err != nil {
		// the value provided is not a parseable boolean, quit
		logger.Error(err, "Invalid response_buffering values", "kongroute", r.Name)
		return
	}

	r.ResponseBuffering = kong.Bool(isEnabled)
}

// overrideHosts appends Host-Aliases to Hosts.
func (r *Route) overrideHosts(logger logr.Logger, anns map[string]string) {
	var hosts []*string
	var annHostAliases []string
	var exists bool
	annHostAliases, exists = annotations.ExtractHostAliases(anns)
	if !exists {
		// the annotation is not set, quit
		return
	}

	// avoid allowing duplicate hosts or host-aliases from being added
	appendIfMissing := func(hosts []*string, host string) []*string {
		for _, uniqueHost := range hosts {
			if *uniqueHost == host {
				return hosts
			}
		}
		return append(hosts, kong.String(host))
	}

	// Merge hosts and host-aliases
	hosts = append(hosts, r.Hosts...)
	for _, hostAlias := range annHostAliases {
		if validHosts.MatchString(hostAlias) {
			hosts = appendIfMissing(hosts, hostAlias)
		} else {
			// Host Alias is not a valid hostname
			logger.Error(nil, "Invalid host alias", "value", hostAlias, "kongroute", r.Name)
			return
		}
	}

	r.Hosts = hosts
}

func (r *Route) overrideHeaders(anns map[string]string) {
	headers, exists := annotations.ExtractHeaders(anns)
	if !exists {
		return
	}
	r.Headers = headers
}

func (r *Route) overridePathHandling(logger logr.Logger, anns map[string]string) {
	val, ok := annotations.ExtractPathHandling(anns)
	if !ok {
		return
	}

	if !validPathHandling.MatchString(val) {
		logger.Error(nil, "Invalid path_handling", "value", val, "kongroute", r.Name)
		return
	}

	r.PathHandling = kong.String(val)
}
