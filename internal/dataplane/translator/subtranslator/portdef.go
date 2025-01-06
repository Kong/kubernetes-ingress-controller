package subtranslator

import (
	netv1 "k8s.io/api/networking/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
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

func PortDefFromPortNumber(port int32) kongstate.PortDef {
	return kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: port}
}
