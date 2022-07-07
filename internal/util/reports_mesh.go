package util

import (
	"sort"
	"strconv"
	"strings"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/meshdetect"
)

var meshKindShortNames = map[meshdetect.MeshKind]string{
	meshdetect.MeshKindIstio:      "i",
	meshdetect.MeshKindLinkerd:    "l",
	meshdetect.MeshKindKuma:       "k",
	meshdetect.MeshKindKongMesh:   "km",
	meshdetect.MeshKindConsul:     "c",
	meshdetect.MeshKindTraefik:    "t",
	meshdetect.MeshKindAWSAppMesh: "a",
}

func serializeMeshDeploymentResults(
	results map[meshdetect.MeshKind]*meshdetect.DeploymentResults,
) string {
	if results == nil {
		return ""
	}

	signals := []string{}
	for _, meshKind := range meshdetect.MeshesToDetect {
		result := results[meshKind]
		if result == nil {
			continue
		}
		// signal3: service exists
		if result.ServiceExists {
			signals = append(signals, meshKindShortNames[meshKind]+"3")
		}
	}

	if len(signals) > 0 {
		// sort the signals (in alphabetical order),
		// then join them together to produce a consistent output for same results.
		sort.Strings(signals)
		return "mdep=\"" + strings.Join(signals, ",") + "\""
	}

	return ""
}

func serializeMeshRunUnderResults(
	results map[meshdetect.MeshKind]*meshdetect.RunUnderResults,
) string {
	if results == nil {
		return ""
	}

	signals := []string{}
	for _, meshKind := range meshdetect.MeshesToDetect {
		result := results[meshKind]
		if result == nil {
			continue
		}

		// signal2: pod/service has annotation
		if result.PodOrServiceAnnotation {
			signals = append(signals, meshKindShortNames[meshKind]+"2")
		}
		// signal3: sidecar injected
		if result.SidecarContainerInjected {
			signals = append(signals, meshKindShortNames[meshKind]+"3")
		}
		// signal4: init container injected
		if result.InitContainerInjected {
			signals = append(signals, meshKindShortNames[meshKind]+"4")
		}
	}

	if len(signals) > 0 {
		// sort the signals to produce a constistent output.
		sort.Strings(signals)
		return "kinm=\"" + strings.Join(signals, ",") + "\""
	}
	return ""
}

func serializeMeshServiceDistribution(
	result *meshdetect.ServiceDistributionResults,
) string {
	if result == nil {
		return ""
	}

	// format: mdist="all100,a10,i20,k50,km50"
	serializedStr := "mdist=\""
	serializedStr = serializedStr + "all" + strconv.Itoa(result.TotalServices)
	if result.MeshDistribution != nil {
		// append number of services running in the mesh, if there are any.
		var signals []string
		for _, meshKind := range meshdetect.MeshesToDetect {
			num := result.MeshDistribution[meshKind]
			if num > 0 {
				signals = append(signals, meshKindShortNames[meshKind]+strconv.Itoa(num))
			}
		}

		if len(signals) > 0 {
			// sort the signals to produce a constistent output.
			sort.Strings(signals)
			serializedStr = serializedStr + "," + strings.Join(signals, ",")
		}
	}
	return serializedStr + "\""
}
