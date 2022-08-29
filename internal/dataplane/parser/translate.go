package parser

import (
	"fmt"
	"strings"

	"github.com/kong/go-kong/kong"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
)

func serviceBackendPortToStr(port netv1.ServiceBackendPort) string {
	if port.Name != "" {
		return fmt.Sprintf("pname-%s", port.Name)
	}
	return fmt.Sprintf("pnum-%d", port.Number)
}

// TODO decide where we want this to live. was previously not exported.

// PathsFromK8s takes a path and Ingress path type and returns a set of Kong route paths that satisfy that path type.
// It optionally adds the Kong 3.x regex path prefix for path types that require a regex path. It rejects unknown path
// types with an error.
func PathsFromK8s(path string, pathType netv1.PathType, addRegexPrefix bool) ([]*string, error) {
	routePaths := []string{}
	routeRegexPaths := []string{}
	switch pathType {
	case netv1.PathTypePrefix:
		base := strings.Trim(path, "/")
		if base == "" {
			routePaths = append(routePaths, "/")
		} else {
			routePaths = append(routePaths, "/"+base+"/")
			routeRegexPaths = append(routeRegexPaths, "/"+base+"$")
		}
	case netv1.PathTypeExact:
		relative := strings.TrimLeft(path, "/")
		routeRegexPaths = append(routeRegexPaths, "/"+relative+"$")
	case netv1.PathTypeImplementationSpecific:
		if path == "" {
			routePaths = append(routePaths, "/")
		} else {
			routePaths = append(routePaths, path)
		}
	default:
		return nil, fmt.Errorf("unknown pathType %v", pathType)
	}

	if addRegexPrefix {
		for i, orig := range routeRegexPaths {
			routeRegexPaths[i] = kongPathRegexPrefix + orig
		}
	}
	routePaths = append(routePaths, routeRegexPaths...)
	return kong.StringSlice(routePaths...), nil
}

var priorityForPath = map[netv1.PathType]int{
	netv1.PathTypeExact:                  300,
	netv1.PathTypePrefix:                 200,
	netv1.PathTypeImplementationSpecific: 100,
}

func PortDefFromServiceBackendPort(sbp *netv1.ServiceBackendPort) kongstate.PortDef {
	switch {
	case sbp.Name != "":
		return kongstate.PortDef{Mode: kongstate.PortModeByName, Name: sbp.Name}
	case sbp.Number != 0:
		return kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: sbp.Number}
	default:
		return kongstate.PortDef{Mode: kongstate.PortModeImplicit}
	}
}

func PortDefFromIntStr(is intstr.IntOrString) kongstate.PortDef {
	if is.Type == intstr.String {
		return kongstate.PortDef{Mode: kongstate.PortModeByName, Name: is.StrVal}
	}
	return kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: is.IntVal}
}
