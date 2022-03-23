package parser

import (
	"fmt"
	"strings"

	"github.com/kong/go-kong/kong"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
)

func serviceBackendPortToStr(port networkingv1.ServiceBackendPort) string {
	if port.Name != "" {
		return fmt.Sprintf("pname-%s", port.Name)
	}
	return fmt.Sprintf("pnum-%d", port.Number)
}

func pathsFromK8s(path string, pathType networkingv1.PathType) ([]*string, error) {
	switch pathType {
	case networkingv1.PathTypePrefix:
		base := strings.Trim(path, "/")
		if base == "" {
			return kong.StringSlice("/"), nil
		}
		return kong.StringSlice(
			"/"+base+"$",
			"/"+base+"/",
		), nil
	case networkingv1.PathTypeExact:
		relative := strings.TrimLeft(path, "/")
		return kong.StringSlice("/" + relative + "$"), nil
	case networkingv1.PathTypeImplementationSpecific:
		if path == "" {
			return kong.StringSlice("/"), nil
		}
		return kong.StringSlice(path), nil
	}

	return nil, fmt.Errorf("unknown pathType %v", pathType)
}

var priorityForPath = map[networkingv1.PathType]int{
	networkingv1.PathTypeExact:                  300,
	networkingv1.PathTypePrefix:                 200,
	networkingv1.PathTypeImplementationSpecific: 100,
}

func PortDefFromServiceBackendPort(sbp *networkingv1.ServiceBackendPort) kongstate.PortDef {
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
