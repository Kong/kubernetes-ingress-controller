package translators

import (
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
)

func PortDefFromServiceBackendPort(sbp *netv1.ServiceBackendPort) kongstate.PortDef {
	switch {
	case sbp.Number != 0:
		return kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: sbp.Number}
	case sbp.Name != "":
		return kongstate.PortDef{Mode: kongstate.PortModeByName, Name: sbp.Name}
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
