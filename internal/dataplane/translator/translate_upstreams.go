package translator

import (
	"fmt"
	"net"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
)

func (t *Translator) getUpstreams(serviceMap map[string]kongstate.Service) ([]kongstate.Upstream, map[string]kongstate.Service) {
	upstreamDedup := make(map[string]struct{}, len(serviceMap))
	var empty struct{}
	upstreams := make([]kongstate.Upstream, 0, len(serviceMap))
	for serviceName, service := range serviceMap {
		// the name of the Upstream for a service must match the service.Host
		// as the Gateway's internal DNS resolve mechanisms will fail to properly
		// resolve the host otherwise.
		name := *service.Host

		if _, exists := upstreamDedup[name]; !exists {
			// targetMap maps the target's Target field to the Target object. The field represents the string "target" field
			// of a Kong target entity. This field is either an IP address or hostname and a port separated by a colon. For
			// example: "10.0.0.1:80", "example.com:9000". We use a map because this field must be unique within an upstream,
			// and we may need to combine multiple Kubernetes backends into a single target for some configurations--some
			// routes may, for example, use the same Service twice or may use two Services with the same selector and same
			// endpoints.
			targetMap := map[string]kongstate.Target{}
			// populate all the kong targets for the upstream given all the backends
			for _, backend := range service.Backends {
				// gather the Kubernetes service for the backend
				backendNamespace := backend.Namespace()

				backendName := backend.Name()
				if backend.IsServiceFacade() {
					// In the case of KongServiceFacade we need to look it up to determine the backing Kubernetes Service.
					svcFacade, err := t.storer.GetKongServiceFacade(backend.Namespace(), backend.Name())
					if err != nil {
						t.registerTranslationFailure(
							fmt.Sprintf("couldn't get KongServiceFacade %s: %v", backend.Name(), err),
							service.Parent,
						)
						continue
					}
					backendName = svcFacade.Spec.Backend.Name
				}
				k8sService, ok := service.K8sServices[fmt.Sprintf("%s/%s", backendNamespace, backendName)]
				if !ok {
					t.registerTranslationFailure(
						fmt.Sprintf("can't add target for backend %s: no kubernetes service found", backendName),
						service.Parent,
					)
					continue
				}

				// determine the port for the backend
				port, err := findPort(k8sService, backend.PortDef())
				if err != nil {
					t.registerTranslationFailure(
						fmt.Sprintf("can't find port for backend kubernetes service: %v", err),
						k8sService, service.Parent,
					)
					continue
				}
				service.Port = lo.ToPtr(int(port.Port))
				serviceMap[serviceName] = service

				// get the new targets for this backend service
				newTargets := getServiceEndpoints(t.logger, t.storer, k8sService, port)

				if len(newTargets) == 0 {
					t.logger.V(util.InfoLevel).Info("No targets could be found for kubernetes service",
						"namespace", k8sService.Namespace, "name", k8sService.Name, "kong_service", *service.Name)
				}

				// if weights were set for the backend then that weight needs to be
				// distributed equally among all the targets.
				if weight, weightPresent := backend.Weight().Get(); weightPresent && len(newTargets) != 0 {
					// initialize the weight of the target based on the weight of the backend
					// which governs that target (and potentially more). If the weight of the
					// backend is 0 then this indicates an intention to drop all targets from
					// this backend from the load-balancer and is a special situation where
					// all derived targets will receive a weight of 0.
					targetWeight := weight

					// if the backend governing this target is not set to a weight of 0,
					// all targets derived from the backend split the weight, therefore
					// equally splitting the traffic load.
					if weight != 0 {
						targetWeight = weight / len(newTargets)
						// minimum weight of 1 if weight zero was not specifically set.
						if targetWeight == 0 {
							targetWeight = 1
						}
					}

					for i := range newTargets {
						newTargets[i].Weight = &targetWeight
					}
				}

				for _, t := range newTargets {
					targetMap = updateTargetMap(targetMap, t)
				}
			}

			targets := lo.Values(targetMap)
			// warn if an upstream was created with 0 targets
			if len(targets) == 0 {
				t.logger.V(util.InfoLevel).Info("No targets found to create upstream", "service_name", *service.Name)
			}

			// define the upstream including all the newly populated targets
			// to load-balance traffic to.
			upstream := kongstate.Upstream{
				Upstream: kong.Upstream{
					Name: kong.String(name),
					Tags: service.Tags, // populated by populateServices already
				},
				Service: service,
				Targets: targets,
			}
			upstreams = append(upstreams, upstream)
			upstreamDedup[name] = empty
		}
	}
	return upstreams, serviceMap
}

// findPort finds a port matching the specified definition in a Kubernetes Service.
func findPort(svc *corev1.Service, wantPort kongstate.PortDef) (*corev1.ServicePort, error) {
	switch wantPort.Mode {
	case kongstate.PortModeByNumber:
		// ExternalName Services have no port declaration of their own
		// We must assume that the user-requested port is valid and construct a ServicePort from it
		if svc.Spec.Type == corev1.ServiceTypeExternalName {
			return &corev1.ServicePort{
				Port:       wantPort.Number,
				TargetPort: intstr.FromInt(int(wantPort.Number)),
			}, nil
		}
		for _, port := range svc.Spec.Ports {
			port := port
			if port.Port == wantPort.Number {
				return &port, nil
			}
		}

	case kongstate.PortModeByName:
		if svc.Spec.Type == corev1.ServiceTypeExternalName {
			return nil, fmt.Errorf("rules with an ExternalName service must specify numeric ports")
		}
		for _, port := range svc.Spec.Ports {
			port := port
			if port.Name == wantPort.Name {
				return &port, nil
			}
			if port.TargetPort.Type == intstr.String && port.TargetPort.String() == wantPort.Name {
				return &port, nil
			}
		}

	case kongstate.PortModeImplicit:
		if svc.Spec.Type == corev1.ServiceTypeExternalName {
			return nil, fmt.Errorf("rules with an ExternalName service must specify numeric ports")
		}
		if len(svc.Spec.Ports) != 1 {
			return nil, fmt.Errorf("in implicit mode, service must have exactly 1 port, has %d", len(svc.Spec.Ports))
		}
		return &svc.Spec.Ports[0], nil

	default:
		return nil, fmt.Errorf("unknown mode %v", wantPort.Mode)
	}

	return nil, fmt.Errorf("no suitable port found")
}

func getServiceEndpoints(
	logger logr.Logger,
	s store.Storer,
	svc *corev1.Service,
	servicePort *corev1.ServicePort,
) []kongstate.Target {
	logger = logger.WithValues(
		"service_name", svc.Name,
		"service_namespace", svc.Namespace,
		"service_port", servicePort,
	)

	// In theory a Service could have multiple port protocols, we need to ensure we gather
	// endpoints based on all the protocols the service is configured for. We always check
	// for TCP as this is the default protocol for service ports.
	protocols := listProtocols(svc)

	// Check if the service is an upstream service through Ingress Class parameters.
	var isSvcUpstream bool
	ingressClassParameters, err := getIngressClassParametersOrDefault(s)
	if err != nil {
		logger.V(util.DebugLevel).Info("Unable to retrieve IngressClassParameters", "error", err)
	} else {
		isSvcUpstream = ingressClassParameters.ServiceUpstream
	}

	// Check all protocols for associated endpoints.
	endpoints := []util.Endpoint{}
	for protocol := range protocols {
		newEndpoints := getEndpoints(logger, svc, servicePort, protocol, s.GetEndpointSlicesForService, isSvcUpstream)
		endpoints = append(endpoints, newEndpoints...)
	}
	if len(endpoints) == 0 {
		logger.V(util.DebugLevel).Info("No active endpoints")
	}

	return targetsForEndpoints(endpoints)
}

// getIngressClassParametersOrDefault returns the parameters for the current ingress class.
// If the cluster operators have specified a set of parameters explicitly, it returns those.
// Otherwise, it returns a default set of parameters.
func getIngressClassParametersOrDefault(s store.Storer) (kongv1alpha1.IngressClassParametersSpec, error) {
	ingressClassName := s.GetIngressClassName()
	ingressClass, err := s.GetIngressClassV1(ingressClassName)
	if err != nil {
		return kongv1alpha1.IngressClassParametersSpec{}, err
	}

	params, err := s.GetIngressClassParametersV1Alpha1(ingressClass)
	if err != nil {
		return kongv1alpha1.IngressClassParametersSpec{}, err
	}

	return params.Spec, nil
}

// getEndpoints returns a list of <endpoint ip>:<port> for a given service/target port combination.
// It also checks if the service is an upstream service either by its annotations
// of by IngressClassParameters configuration provided as a flag.
func getEndpoints(
	logger logr.Logger,
	service *corev1.Service,
	port *corev1.ServicePort,
	proto corev1.Protocol,
	getEndpointSlices func(string, string) ([]*discoveryv1.EndpointSlice, error),
	isSvcUpstream bool,
) []util.Endpoint {
	if service == nil || port == nil {
		return []util.Endpoint{}
	}

	// If service is an upstream service...
	if isSvcUpstream || annotations.HasServiceUpstreamAnnotation(service.Annotations) {
		// ... return its address as the only endpoint.
		return []util.Endpoint{
			{
				Address: service.Name + "." + service.Namespace + ".svc",
				Port:    fmt.Sprint(port.Port),
			},
		}
	}

	logger = logger.WithValues(
		"service_name", service.Name,
		"service_namespace", service.Namespace,
		"service_port", port.String(),
	)

	// ExternalName services
	if service.Spec.Type == corev1.ServiceTypeExternalName {
		logger.V(util.DebugLevel).Info("Found service of type=ExternalName")
		return []util.Endpoint{
			{
				Address: service.Spec.ExternalName,
				Port:    port.TargetPort.String(),
			},
		}
	}

	logger.V(util.DebugLevel).Info("Fetching EndpointSlices")
	endpointSlices, err := getEndpointSlices(service.Namespace, service.Name)
	if err != nil {
		logger.Error(err, "Error fetching EndpointSlices")
		return []util.Endpoint{}
	}
	logger.V(util.DebugLevel).Info("Fetched EndpointSlices", "count", len(endpointSlices))

	// Avoid duplicated upstream servers when the service contains
	// multiple port definitions sharing the same target port.
	uniqueUpstream := make(map[util.Endpoint]struct{})
	upstreamServers := make([]util.Endpoint, 0)
	for _, endpointSlice := range endpointSlices {
		for _, p := range endpointSlice.Ports {
			if p.Port == nil || *p.Port < 0 || *p.Protocol != proto || *p.Name != port.Name {
				continue
			}
			upstreamPort := fmt.Sprint(*p.Port)
			for _, endpoint := range endpointSlice.Endpoints {
				// Ready indicates that this endpoint is prepared to receive traffic, according to whatever
				// system is managing the endpoint. A nil value indicates an unknown state.
				// In most cases consumers should interpret this unknown state as ready.
				// Field Ready has the same semantic as Endpoints from CoreV1 in Addresses.
				// https://kubernetes.io/docs/concepts/services-networking/endpoint-slices/#conditions
				if endpoint.Conditions.Ready != nil && !*endpoint.Conditions.Ready {
					continue
				}
				// One address per endpoint is rather expected (allowing multiple is due to historical reasons)
				// read more https://github.com/kubernetes/kubernetes/issues/106267#issuecomment-978770401.
				// These are all assumed to be fungible and clients may choose to only use the first element.
				upstreamServer := util.Endpoint{
					Address: endpoint.Addresses[0],
					Port:    upstreamPort,
				}
				if _, exists := uniqueUpstream[upstreamServer]; !exists {
					upstreamServers = append(upstreamServers, upstreamServer)
					uniqueUpstream[upstreamServer] = struct{}{}
				}
			}
		}
	}
	logger.V(util.DebugLevel).Info("Found endpoints", "endpoints", upstreamServers)
	return upstreamServers
}

// targetWeightOrDefault returns the effective value of a target weight pointer. If the pointer is non-nil, it returns
// the pointee. If the pointer is nil, it returns 100, the default Kong target weight. This allows us to sum
// deduplicated targets' weights if one happens to be unset in the controller.
func targetWeightOrDefault(in *int) int {
	if in != nil {
		return *in
	}
	return 100
}

func updateTargetMap(targetMap map[string]kongstate.Target, t kongstate.Target) map[string]kongstate.Target {
	// See https://github.com/Kong/kubernetes-ingress-controller/issues/5761:
	// Duplicate targets will appear in configurations that use Services with the same selector, which are used
	// by some rollout systems. We need to deduplicate them while honoring the total weight.
	//
	// Because kongstate.Target is a nested kong.Target and the target IP is also a field named Target, the
	// key names are a bit silly: while fields like t.Weight and t.Upstream resolve fine, t.Target does not, and
	// instead requires t.Target.Target. For consistency, everything below explicitly includes the nested object
	// name, so t.Target.Weight instead of t.Weight.
	if existing, ok := targetMap[*t.Target.Target]; ok {
		sum := targetWeightOrDefault(existing.Target.Weight) + targetWeightOrDefault(t.Target.Weight)
		existing.Target.Weight = &sum
		targetMap[*t.Target.Target] = existing
	} else {
		targetMap[*t.Target.Target] = t
	}
	return targetMap
}

// targetsForEndpoints generates kongstate.Target objects for each util.Endpoint provided.
func targetsForEndpoints(endpoints []util.Endpoint) []kongstate.Target {
	targets := []kongstate.Target{}
	for _, endpoint := range endpoints {
		addr := endpoint.Address
		parsed := net.ParseIP(endpoint.Address)
		if parsed != nil {
			if parsed.To4() == nil {
				// If we have an IPv6 endpoint, we need to surround it with brackets, else the port concat after this will
				// treat the port as part of the address.
				addr = fmt.Sprintf("[%s]", endpoint.Address)
			}
		}
		target := kongstate.Target{
			Target: kong.Target{
				Target: kong.String(addr + ":" + endpoint.Port),
			},
		}
		targets = append(targets, target)
	}
	return targets
}

// listProtocols is a helper function to map out all the in-use corev1.Protocols
// for a service given a corev1.Service object.
//
// TODO: due to historical logic this function defaults to assuming TCP protocol
// is valid for the Service and its endpoints, however we need to follow up
// on this as this is not technically correct and causes waste.
// See: https://github.com/Kong/kubernetes-ingress-controller/issues/1429
func listProtocols(svc *corev1.Service) map[corev1.Protocol]bool {
	protocols := map[corev1.Protocol]bool{corev1.ProtocolTCP: true}
	for _, port := range svc.Spec.Ports {
		if port.Protocol != "" {
			protocols[port.Protocol] = true
		}
	}
	return protocols
}
