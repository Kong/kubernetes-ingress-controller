package meshdetect

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

// MeshKind defines names of different service meshes.
type MeshKind string

const (
	MeshKindIstio      MeshKind = "istio"
	MeshKindLinkerd    MeshKind = "linkerd"
	MeshKindKuma       MeshKind = "kuma"
	MeshKindKongMesh   MeshKind = "kong-mesh"
	MeshKindConsul     MeshKind = "consul"
	MeshKindTraefik    MeshKind = "traefik"
	MeshKindAWSAppMesh MeshKind = "aws-app-mesh"
)

// MeshesToDetect is the list of meshes to detect.
var MeshesToDetect = []MeshKind{
	MeshKindIstio,
	MeshKindLinkerd,
	MeshKindKuma,
	MeshKindKongMesh,
	MeshKindConsul,
	MeshKindTraefik,
	MeshKindAWSAppMesh,
}

// meshServiceName is one of the services of system components in each mesh.
var meshServiceName = map[MeshKind]string{
	MeshKindIstio:      "istiod",
	MeshKindLinkerd:    "linkerd-proxy-injector",
	MeshKindKuma:       "kuma-control-plane",
	MeshKindKongMesh:   "kong-mesh-control-plane",
	MeshKindConsul:     "consul-server",
	MeshKindTraefik:    "traefik-mesh-controller",
	MeshKindAWSAppMesh: "appmesh-controller-webhook-service",
}

// mustMakeRequirement is a tool function to make labels.Requirement in map initializations.
// it runs labels.NewRequirement without opts and panics on error.
func mustMakeRequirement(key string, op selection.Operator, vals []string) *labels.Requirement {
	req, err := labels.NewRequirement(key, op, vals)
	if err != nil {
		panic("failed to make requirement:" + err.Error())
	}
	return req
}

// meshPodAnnotations is the annotation of pod indicating that the pod should be injected with sidecars.
var meshPodAnnotations = map[MeshKind]*labels.Requirement{
	MeshKindIstio:    mustMakeRequirement("sidecar.istio.io/status", selection.Exists, []string{}),
	MeshKindLinkerd:  mustMakeRequirement("linkerd.io/proxy-version", selection.Exists, []string{}),
	MeshKindKuma:     mustMakeRequirement("kuma.io/sidecar-injected", selection.Equals, []string{"true"}),
	MeshKindKongMesh: mustMakeRequirement("kuma.io/sidecar-injected", selection.Equals, []string{"true"}),
	MeshKindConsul:   mustMakeRequirement("consul.hashicorp.com/connect-inject-status", selection.Equals, []string{"injected"}),
}

// meshServiceAnnotations is the annotation of service indicating that
// the service in managed by the service mesh. (currently only for traefik)
var meshServiceAnnotations = map[MeshKind]*labels.Requirement{
	MeshKindTraefik: mustMakeRequirement("mesh.traefik.io/traffic-type", selection.In, []string{"HTTP", "TCP"}),
}

// meshSidecarContainerName is the name of sidecar container injected by each service mesh.
var meshSidecarContainerName = map[MeshKind]string{
	MeshKindIstio:      "istio-proxy",
	MeshKindLinkerd:    "linkerd-proxy",
	MeshKindKuma:       "kuma-sidecar",
	MeshKindKongMesh:   "kuma-sidecar",
	MeshKindConsul:     "envoy-sidecar",
	MeshKindAWSAppMesh: "envoy",
}

const (
	// awsAppMeshEnvoyImageName is the image used for aws appmesh sidecars.
	// reference from AWS: https://docs.aws.amazon.com/app-mesh/latest/userguide/envoy.html
	awsAppMeshEnvoyImageName = "aws-appmesh-envoy"
)

// meshInitContainerName is the name of init container injected by each service mesh.
var meshInitContainerName = map[MeshKind]string{
	MeshKindIstio:      "istio-init",
	MeshKindLinkerd:    "linkerd-init",
	MeshKindKuma:       "kuma-init",
	MeshKindKongMesh:   "kuma-init",
	MeshKindConsul:     "consul-connect-inject-init",
	MeshKindAWSAppMesh: "proxyinit",
}
