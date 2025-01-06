package subtranslator

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	incubatorv1alpha1 "github.com/kong/kubernetes-configuration/api/incubator/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// -----------------------------------------------------------------------------
// Ingress Translation - Public Functions
// -----------------------------------------------------------------------------

// TranslatedKubernetesObjectsCollector is an interface for collecting Kubernetes objects that have been translated
// successfully.
type TranslatedKubernetesObjectsCollector interface {
	Add(client.Object)
}

// FailuresCollector is an interface for collecting failures during translation.
type FailuresCollector interface {
	PushResourceFailure(reason string, causingObjects ...client.Object)
}

type TranslateIngressFeatureFlags struct {
	// ExpressionRoutes indicates whether to translate Kubernetes objects to expression based Kong Routes.
	ExpressionRoutes bool

	// KongServiceFacade indicates whether we should support KongServiceFacade as Ingress backends.
	KongServiceFacade bool
}

// TranslateIngresses receives a slice of Kubernetes Ingress objects and produces a translated set of kong.Services
// and kong.Routes which will come wrapped in a kongstate.Service object.
func TranslateIngresses(
	ingresses []*netv1.Ingress,
	icp kongv1alpha1.IngressClassParametersSpec,
	flags TranslateIngressFeatureFlags,
	translatedObjectsCollector TranslatedKubernetesObjectsCollector,
	failuresCollector FailuresCollector,
	storer store.Storer,
) map[string]kongstate.Service {
	index := newIngressTranslationIndex(flags, failuresCollector, storer)
	for _, ingress := range ingresses {
		prependRegexPrefix := MaybePrependRegexPrefixForIngressV1Fn(ingress, icp.EnableLegacyRegexDetection)
		index.Add(ingress, prependRegexPrefix)
		translatedObjectsCollector.Add(ingress)
	}

	return index.Translate()
}

// -----------------------------------------------------------------------------
// Ingress Translation - Private Consts & Vars
// -----------------------------------------------------------------------------

var defaultHTTPIngressPathType = netv1.PathTypeImplementationSpecific

const (
	defaultHTTPPort = 80
	defaultRetries  = 5
)

// defaultServiceTimeoutKongFormat returns the defaultServiceTimeout in format
// expected by Kong (pointer to an integer representing milliseconds).
//
// defaultServiceTimeout indicates the amount of time by default that we wait
// for connections to an underlying Kubernetes service to complete in the
// data-plane. The current value is based on a historical default that started
// in version 0 of the ingress controller.
func defaultServiceTimeoutInKongFormat() *int {
	const defaultServiceTimeout = time.Second * 60
	return kong.Int(int(defaultServiceTimeout.Milliseconds()))
}

// -----------------------------------------------------------------------------
// Ingress Translation - Private - Index
// -----------------------------------------------------------------------------

// ingressTranslationIndex is a de-duplicating index of the contents of a
// Kubernetes Ingress resource, where the key is a combination of data from
// that resource which makes it unique for the purpose of translating it into
// kong.Services and kong.Routes and the value is a combination of various
// metadata needed to configure and name those kong.Services and kong.Routes
// plus the URL paths which make those routes actionable. This index is used
// to enable compiling a minimal set of kong.Routes when translating into
// Kong resources for each rule in the ingress spec, where the combination of:
//
// - ingress.Namespace
// - ingress.Name
// - host for the Ingress rule
// - Kubernetes Service for the Ingress rule
// - the port for the Kubernetes Service
//
// are unique. For ingress spec rules which are not unique along those
// data-points, a separate kong.Service and separate kong.Routes will be created
// for each unique combination.
type ingressTranslationIndex struct {
	cache             map[string]*ingressTranslationMeta
	featureFlags      TranslateIngressFeatureFlags
	failuresCollector FailuresCollector
	storer            store.Storer
}

func newIngressTranslationIndex(flags TranslateIngressFeatureFlags, failuresCollector FailuresCollector, storer store.Storer) *ingressTranslationIndex {
	return &ingressTranslationIndex{
		cache:             make(map[string]*ingressTranslationMeta),
		featureFlags:      flags,
		failuresCollector: failuresCollector,
		storer:            storer,
	}
}

type addRegexPrefixFn func(string) *string

func (i *ingressTranslationIndex) Add(ingress *netv1.Ingress, addRegexPrefix addRegexPrefixFn) {
	for _, ingressRule := range ingress.Spec.Rules {
		if ingressRule.HTTP == nil || len(ingressRule.HTTP.Paths) < 1 {
			continue
		}

		for _, httpIngressPath := range ingressRule.HTTP.Paths {
			httpIngressPath.Path = flattenMultipleSlashes(httpIngressPath.Path)

			if httpIngressPath.Path == "" {
				httpIngressPath.Path = "/"
			}

			if httpIngressPath.PathType == nil {
				httpIngressPath.PathType = &defaultHTTPIngressPathType
			}

			backend, err := i.getIngressPathBackend(ingress.Namespace, httpIngressPath)
			if err != nil {
				i.failuresCollector.PushResourceFailure(fmt.Sprintf("failed to get backend for ingress path %q: %s", httpIngressPath.Path, err), ingress)
				continue
			}

			kongRouteName := backend.intoKongRouteName(k8stypes.NamespacedName{Namespace: ingress.Namespace, Name: ingress.Name}, ingressRule.Host)
			meta, ok := i.cache[kongRouteName]
			if !ok {
				meta = &ingressTranslationMeta{
					ingressNamespace: ingress.Namespace,
					ingressName:      ingress.Name,
					ingressUID:       string(ingress.UID),
					ingressHost:      ingressRule.Host,
					ingressTags:      util.GenerateTagsForObject(ingress),
					backend:          backend,
					addRegexPrefixFn: addRegexPrefix,
				}
			}

			meta.parentIngress = ingress
			meta.paths = append(meta.paths, httpIngressPath)
			i.cache[kongRouteName] = meta
		}
	}
}

func (i *ingressTranslationIndex) getIngressPathBackend(namespace string, httpIngressPath netv1.HTTPIngressPath) (ingressTranslationMetaBackend, error) {
	if service := httpIngressPath.Backend.Service; service != nil {
		return newIngressTranslationMetaBackendForKubernetesService(
			service.Name,
			PortDefFromServiceBackendPort(&service.Port),
		), nil
	}

	if resource := httpIngressPath.Backend.Resource; resource != nil {
		if !IsKongServiceFacade(resource) {
			gk := resource.Kind
			if resource.APIGroup != nil {
				gk = *resource.APIGroup + "/" + gk
			}
			return ingressTranslationMetaBackend{}, fmt.Errorf("unknown resource type %s", gk)
		}
		if !i.featureFlags.KongServiceFacade {
			return ingressTranslationMetaBackend{}, fmt.Errorf("KongServiceFacade is not enabled, please set the %q feature gate to 'true' to enable it", featuregates.KongServiceFacade)
		}

		serviceFacade, err := i.storer.GetKongServiceFacade(namespace, resource.Name)
		if err != nil {
			return ingressTranslationMetaBackend{}, fmt.Errorf("failed to get KongServiceFacade %q: %w", resource.Name, err)
		}
		return newIngressTranslationMetaBackendForKongServiceFacade(
			resource.Name,
			PortDefFromPortNumber(serviceFacade.Spec.Backend.Port),
			serviceFacade,
		), nil
	}

	// Should never happen since the Ingress API validation should catch this.
	return ingressTranslationMetaBackend{}, fmt.Errorf("no Service or Resource specified for Ingress path")
}

// IsKongServiceFacade returns true if the given resource reference is a KongServiceFacade.
func IsKongServiceFacade(resource *corev1.TypedLocalObjectReference) bool {
	return resource.Kind == incubatorv1alpha1.KongServiceFacadeKind &&
		resource.APIGroup != nil && *resource.APIGroup == incubatorv1alpha1.GroupVersion.Group
}

func (i *ingressTranslationIndex) Translate() map[string]kongstate.Service {
	kongStateServiceCache := make(map[string]kongstate.Service)
	for _, meta := range i.cache {
		kongServiceName := meta.generateKongServiceName()
		kongStateService, ok := kongStateServiceCache[kongServiceName]
		if !ok {
			var err error
			kongStateService, err = meta.translateIntoKongStateService(kongServiceName, meta.backend.port)
			if err != nil {
				i.failuresCollector.PushResourceFailure(fmt.Sprintf("failed to translate Ingress into Kong Service: %s", err), meta.parentIngress)
				continue
			}
		}

		if i.featureFlags.ExpressionRoutes {
			route := meta.translateIntoKongExpressionRoute()
			kongStateService.Routes = append(kongStateService.Routes, *route)
		} else {
			route := meta.translateIntoKongRoute()
			kongStateService.Routes = append(kongStateService.Routes, *route)
		}
		sort.SliceStable(kongStateService.Routes, func(i, j int) bool {
			return *kongStateService.Routes[i].Name < *kongStateService.Routes[j].Name
		})
		kongStateServiceCache[kongServiceName] = kongStateService
	}

	return kongStateServiceCache
}

// -----------------------------------------------------------------------------
// Ingress Translation - Private - Metadata
// -----------------------------------------------------------------------------

type ingressTranslationMeta struct {
	parentIngress    client.Object
	ingressNamespace string
	ingressName      string
	ingressUID       string
	ingressHost      string
	ingressTags      []*string
	backend          ingressTranslationMetaBackend
	paths            []netv1.HTTPIngressPath
	addRegexPrefixFn addRegexPrefixFn
}

type ingressPathBackendType string

const (
	ingressPathBackendTypeKongServiceFacade ingressPathBackendType = "KongServiceFacade"
	ingressPathBackendTypeKubernetesService ingressPathBackendType = "KubernetesService"
)

type ingressTranslationMetaBackend struct {
	// backendType is the type of backend. If left empty, it's assumed to be a Kubernetes Service.
	backendType ingressPathBackendType

	// name is the Kubernetes object name of the backend.
	name string

	// port is the port of the backend.
	port kongstate.PortDef

	// parentKongServiceFacade is the parent KongServiceFacade object if the backend is a KongServiceFacade. Otherwise, it's nil.
	parentKongServiceFacade *incubatorv1alpha1.KongServiceFacade
}

func newIngressTranslationMetaBackendForKongServiceFacade(
	name string,
	port kongstate.PortDef,
	parentKongServiceFacade *incubatorv1alpha1.KongServiceFacade,
) ingressTranslationMetaBackend {
	return ingressTranslationMetaBackend{
		backendType:             ingressPathBackendTypeKongServiceFacade,
		name:                    name,
		port:                    port,
		parentKongServiceFacade: parentKongServiceFacade,
	}
}

func newIngressTranslationMetaBackendForKubernetesService(
	name string,
	port kongstate.PortDef,
) ingressTranslationMetaBackend {
	return ingressTranslationMetaBackend{
		backendType: ingressPathBackendTypeKubernetesService,
		name:        name,
		port:        port,
	}
}

// intoKongRouteName constructs a Kong Route name for the ingressTranslationMeta object.
func (b ingressTranslationMetaBackend) intoKongRouteName(ingress k8stypes.NamespacedName, host string) string {
	// For KongServiceFacade backends, Kong Routes are created separately for the following combination
	// `<ingress-namespace>.<ingress-name>.<host>.<service-facade-name>.svc.facade`.

	// Note: in the case of KongServiceFacade, we don't use the port in the Kong Route name
	// because KongServiceFacade may only specify one port (unlike Kubernetes Service) so it's not necessary.
	if b.backendType == ingressPathBackendTypeKongServiceFacade {
		return fmt.Sprintf("%s.%s.%s.%s.svc.facade", ingress.Namespace, ingress.Name, host, b.name)
	}

	// Otherwise, we assume it's a Kubernetes Service and create Kong Routes for the following combination
	// `<ingress-namespace>.<ingress-name>.<host>.<service-name>.<service-port>`.
	return fmt.Sprintf("%s.%s.%s.%s.%s", ingress.Namespace, ingress.Name, b.name, host, b.port.CanonicalString())
}

// isServiceFacade returns true if the backend is a KongServiceFacade.
func (b ingressTranslationMetaBackend) isServiceFacade() bool {
	return b.backendType == ingressPathBackendTypeKongServiceFacade
}

func (m *ingressTranslationMeta) translateIntoKongStateService(
	kongServiceName string,
	portDef kongstate.PortDef,
) (kongstate.Service, error) {
	if m.backend.isServiceFacade() {
		serviceBackend, err := kongstate.NewServiceBackendForServiceFacade(
			k8stypes.NamespacedName{
				Namespace: m.parentIngress.GetNamespace(),
				Name:      m.backend.name,
			},
			portDef,
		)
		if err != nil {
			return kongstate.Service{}, fmt.Errorf("failed to create ServiceBackend for KongServiceFacade %q: %w", m.backend.name, err)
		}

		return kongstate.Service{
			Namespace: m.parentIngress.GetNamespace(),
			Service: kong.Service{
				Name:           kong.String(kongServiceName),
				Host:           kong.String(fmt.Sprintf("%s.%s.svc.facade", m.parentIngress.GetNamespace(), m.backend.name)),
				Port:           kong.Int(defaultHTTPPort),
				Protocol:       kong.String("http"),
				Path:           kong.String("/"),
				ConnectTimeout: defaultServiceTimeoutInKongFormat(),
				ReadTimeout:    defaultServiceTimeoutInKongFormat(),
				WriteTimeout:   defaultServiceTimeoutInKongFormat(),
				Retries:        kong.Int(defaultRetries),
			},
			Backends: []kongstate.ServiceBackend{serviceBackend},
			Parent:   m.backend.parentKongServiceFacade,
		}, nil
	}

	serviceBackend, err := kongstate.NewServiceBackendForService(
		k8stypes.NamespacedName{
			Namespace: m.parentIngress.GetNamespace(),
			Name:      m.backend.name,
		},
		portDef,
	)
	if err != nil {
		return kongstate.Service{}, fmt.Errorf("failed to create ServiceBackend for Kubernetes Service %q: %w", m.backend.name, err)
	}

	// Otherwise, we assume it's a Kubernetes Service.
	return kongstate.Service{
		Namespace: m.parentIngress.GetNamespace(),
		Service: kong.Service{
			Name:           kong.String(kongServiceName),
			Host:           kong.String(fmt.Sprintf("%s.%s.%s.svc", m.backend.name, m.parentIngress.GetNamespace(), portDef.CanonicalString())),
			Port:           kong.Int(defaultHTTPPort),
			Protocol:       kong.String("http"),
			Path:           kong.String("/"),
			ConnectTimeout: defaultServiceTimeoutInKongFormat(),
			ReadTimeout:    defaultServiceTimeoutInKongFormat(),
			WriteTimeout:   defaultServiceTimeoutInKongFormat(),
			Retries:        kong.Int(defaultRetries),
		},
		Backends: []kongstate.ServiceBackend{serviceBackend},
		Parent:   m.parentIngress,
	}, nil
}

func (m *ingressTranslationMeta) generateKongServiceName() string {
	if m.backend.isServiceFacade() {
		// For KongServiceFacade we create one Kong Service per KongServiceFacade.
		// The naming pattern is `<facade-namespace>.<facade-name>.svc.facade`.
		return fmt.Sprintf("%s.%s.svc.facade", m.parentIngress.GetNamespace(), m.backend.name)
	}

	// For Kubernetes Services, we create one Kong Service per Kubernetes Service + port combination.
	// The naming pattern is `<service-namespace>.<service-name>.<service-port>`.
	return fmt.Sprintf(
		"%s.%s.%s",
		m.parentIngress.GetNamespace(),
		m.backend.name,
		m.backend.port.CanonicalString(),
	)
}

func (m *ingressTranslationMeta) translateIntoKongRoute() *kongstate.Route {
	ingressHost := m.ingressHost
	if strings.Contains(ingressHost, "*") {
		// '_' is not allowed in host, so we use '_' to replace '*' since '*' is not allowed in Kong.
		ingressHost = strings.ReplaceAll(ingressHost, "*", "_")
	}
	routeName := m.backend.intoKongRouteName(k8stypes.NamespacedName{Namespace: m.ingressNamespace, Name: m.ingressName}, ingressHost)

	route := &kongstate.Route{
		Ingress: util.FromK8sObject(m.parentIngress),
		Route: kong.Route{
			Name:              kong.String(routeName),
			StripPath:         kong.Bool(false),
			PreserveHost:      kong.Bool(true),
			Protocols:         kong.StringSlice("http", "https"),
			RegexPriority:     kong.Int(0),
			RequestBuffering:  kong.Bool(true),
			ResponseBuffering: kong.Bool(true),
			Tags:              m.ingressTags,
		},
	}

	if m.ingressHost != "" {
		route.Route.Hosts = append(route.Route.Hosts, kong.String(m.ingressHost))
	}

	for _, httpIngressPath := range m.paths {
		paths := PathsFromIngressPaths(httpIngressPath)
		for i, path := range paths {
			paths[i] = m.addRegexPrefixFn(*path)
		}
		route.Paths = append(route.Paths, paths...)
	}

	return route
}

// -----------------------------------------------------------------------------
// Ingress Translation - Private - Helper Functions
// -----------------------------------------------------------------------------

// TODO this is exported because most of the translator translate functions are still in the translator package. if/when we
// refactor to move them here, this should become private.

// PathsFromIngressPaths takes a path and Ingress path type and returns a set of Kong route paths that satisfy that path
// type. It optionally adds the Kong 3.x regex path prefix for path types that require a regex path. It rejects
// unknown path types with an error.
func PathsFromIngressPaths(httpIngressPath netv1.HTTPIngressPath) []*string {
	routePaths := []string{}
	routeRegexPaths := []string{}
	if httpIngressPath.PathType == nil {
		return nil
	}

	switch *httpIngressPath.PathType {
	case netv1.PathTypePrefix:
		base := strings.Trim(httpIngressPath.Path, "/")
		if base == "" {
			routePaths = append(routePaths, "/")
		} else {
			routePaths = append(routePaths, "/"+base+"/")
			routeRegexPaths = append(routeRegexPaths, KongPathRegexPrefix+"/"+base+"$")
		}
	case netv1.PathTypeExact:
		relative := strings.TrimLeft(httpIngressPath.Path, "/")
		routeRegexPaths = append(routeRegexPaths, KongPathRegexPrefix+"/"+relative+"$")
	case netv1.PathTypeImplementationSpecific:
		if httpIngressPath.Path == "" {
			routePaths = append(routePaths, "/")
		} else {
			routePaths = append(routePaths, httpIngressPath.Path)
		}
	default:
		// the default case here is mostly to provide a home for this comment: we explicitly do not handle unknown
		// PathTypes, and leave it up to the callers if they want to handle empty responses. barring spec changes,
		// however, this should not be a concern: Kubernetes rejects any Ingress with an unknown PathType already, so
		// none should ever end up here. prior versions of this function returned an error in this case, but it
		// should be unnecessary in practice and not returning one simplifies the call chain above (this would be the
		// only part of translation that can error)
		return nil
	}

	routePaths = append(routePaths, routeRegexPaths...)
	return kong.StringSlice(routePaths...)
}

func flattenMultipleSlashes(path string) string {
	out := make([]rune, 0, len(path))
	in := []rune(path)
	for i := 0; i < len(in); i++ {
		c := in[i]
		if c == '/' {
			for i < len(in)-1 && in[i+1] == '/' {
				i++
			}
		}
		out = append(out, c)
	}
	return string(out)
}

// legacyRegexPathExpression is the regular expression used by Kong <3.0 to determine if a path is not a regex.
var legacyRegexPathExpression = regexp.MustCompile(`^[a-zA-Z0-9\.\-_~/%]*$`)

// MaybePrependRegexPrefix takes a path, controller regex prefix, and a legacy heuristic toggle. It returns the path
// with the Kong regex path prefix if it either began with the controller prefix or did not, but matched the legacy
// heuristic, and the heuristic was enabled.
func MaybePrependRegexPrefix(path, controllerPrefix string, applyLegacyHeuristic bool) string {
	if strings.HasPrefix(path, controllerPrefix) {
		path = strings.Replace(path, controllerPrefix, KongPathRegexPrefix, 1)
	} else if applyLegacyHeuristic {
		// this regex matches if the path _is not_ considered a regex by Kong 2.x
		if legacyRegexPathExpression.FindString(path) == "" {
			if !strings.HasPrefix(path, KongPathRegexPrefix) {
				path = KongPathRegexPrefix + path
			}
		}
	}
	return path
}

// MaybePrependRegexPrefixForIngressV1Fn returns a function that prepends a regex prefix to a path for a given netv1.Ingress.
func MaybePrependRegexPrefixForIngressV1Fn(ingress *netv1.Ingress, applyLegacyHeuristic bool) func(path string) *string {
	// If the ingress has a regex prefix annotation, use that, otherwise use the controller default.
	regexPrefix := ControllerPathRegexPrefix
	if prefix, ok := ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.RegexPrefixKey]; ok {
		regexPrefix = prefix
	}

	return func(path string) *string {
		return lo.ToPtr(MaybePrependRegexPrefix(path, regexPrefix, applyLegacyHeuristic))
	}
}

type runeType int

const (
	runeTypeEscape runeType = iota
	runeTypeMark
	runeTypeDigit
	runeTypePlain
)

// generateRewriteURIConfig parses uri with SM of four states.
// `runeTypeEscape` indicates `\` encountered and `$` expected, the SM state will transfer
// to `runeTypePlain`.
// `runeTypeMark` indicates `$` encountered and digit expected, the SM state will transfer
// to `runeTypeDigit`.
// `runeTypeDigit` indicates digit encountered and digit expected. If the following
// character is still digit, the SM state will remain unchanged. Otherwise, a new capture
// group will be created and the SM state will transfer to `runeTypePlain`.
// `runeTypePlain` indicates the following character is plain text other than `$` and `\`.
// The former will cause the SM state to transfer to `runeTypeMark` and the latter will
// cause the SM state to transfer to `runeTypeEscape`.
func generateRewriteURIConfig(uri string) (string, error) {
	out := strings.Builder{}
	lastRuneType := runeTypePlain
	for i, char := range uri {
		switch lastRuneType {
		case runeTypeEscape:
			if char != '$' {
				return "", fmt.Errorf("unexpected %c at pos %d", char, i)
			}

			out.WriteRune(char)
			lastRuneType = runeTypePlain

		case runeTypeMark:
			if !unicode.IsDigit(char) {
				return "", fmt.Errorf("unexpected %c at pos %d", char, i)
			}

			out.WriteString("$(uri_captures[")
			out.WriteRune(char)
			lastRuneType = runeTypeDigit

		case runeTypeDigit:
			if unicode.IsDigit(char) {
				out.WriteRune(char)
			} else {
				out.WriteString("])")
				switch char {
				case '$':
					lastRuneType = runeTypeMark
				case '\\':
					lastRuneType = runeTypeEscape
				default:
					out.WriteRune(char)
					lastRuneType = runeTypePlain
				}
			}

		case runeTypePlain:
			switch char {
			case '$':
				lastRuneType = runeTypeMark
			case '\\':
				lastRuneType = runeTypeEscape
			default:
				out.WriteRune(char)
			}
		}
	}

	if lastRuneType == runeTypeDigit {
		out.WriteString("])")
		lastRuneType = runeTypePlain
	}

	if lastRuneType != runeTypePlain {
		return "", fmt.Errorf("unexpected end of string")
	}

	return out.String(), nil
}

// MaybeRewriteURI appends a request-transformer plugin to Kong routes based on
// the value of konghq.com/rewrite annotation configured on related K8s Ingresses.
func MaybeRewriteURI(service *kongstate.Service, rewriteURIEnable bool) error {
	for i := range service.Routes {
		route := &service.Routes[i]

		rewriteURI, exists := annotations.ExtractRewriteURI(route.Ingress.Annotations)
		if !exists {
			continue
		}
		if !rewriteURIEnable {
			return fmt.Errorf("konghq.com/rewrite annotation not supported when rewrite uris disabled")
		}
		if rewriteURI == "" {
			rewriteURI = "/"
		}

		config, err := generateRewriteURIConfig(rewriteURI)
		if err != nil {
			return err
		}
		route.Plugins = append(route.Plugins, kong.Plugin{
			Name: kong.String("request-transformer"),
			Config: kong.Configuration{
				"replace": map[string]string{
					"uri": config,
				},
			},
		})
	}
	return nil
}
