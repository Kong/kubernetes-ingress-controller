package util

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/meshdetect"
)

func TestSerializeMeshDeploymentResults(t *testing.T) {
	testCases := []struct {
		caseName      string
		results       map[meshdetect.MeshKind]*meshdetect.DeploymentResults
		serializedStr string
	}{
		{
			caseName: "deployment:kong-mesh",
			results: map[meshdetect.MeshKind]*meshdetect.DeploymentResults{
				meshdetect.MeshKindKongMesh: {
					ServiceExists: true,
				},
			},
			serializedStr: "mdep=\"km3\"",
		},
		{
			caseName: "deployment:traefik",
			results: map[meshdetect.MeshKind]*meshdetect.DeploymentResults{
				meshdetect.MeshKindTraefik: {
					ServiceExists: true,
				},
			},
			serializedStr: "mdep=\"t3\"",
		},
		{
			caseName: "deployment:consul,aws-app-mesh",
			results: map[meshdetect.MeshKind]*meshdetect.DeploymentResults{
				meshdetect.MeshKindConsul: {
					ServiceExists: true,
				},
				meshdetect.MeshKindAWSAppMesh: {
					ServiceExists: true,
				},
			},
			serializedStr: "mdep=\"a3,c3\"",
		},
		{
			caseName:      "deployment:nil results should produce empty string",
			results:       nil,
			serializedStr: "",
		},
	}

	for _, tc := range testCases {
		// t.Run runs function in separate goroutine, so we need to assign tc to a local variable.
		tc := tc
		t.Run(tc.caseName, func(t *testing.T) {
			serialized := serializeMeshDeploymentResults(tc.results)
			require.Equalf(t, tc.serializedStr, serialized,
				"case %s: serialized message should be the same as expected", tc.caseName)
		})

	}
}

func TestSerializeMeshRunUnderResult(t *testing.T) {
	testCases := []struct {
		caseName      string
		results       map[meshdetect.MeshKind]*meshdetect.RunUnderResults
		serializedStr string
	}{
		{
			caseName: "run_under:istio,linkerd",
			results: map[meshdetect.MeshKind]*meshdetect.RunUnderResults{
				meshdetect.MeshKindIstio: {
					PodOrServiceAnnotation: true,
				},
				meshdetect.MeshKindLinkerd: {
					PodOrServiceAnnotation:   true,
					SidecarContainerInjected: true,
					InitContainerInjected:    true,
				},
			},
			serializedStr: "kinm=\"i2,l2,l3,l4\"",
		},
		{
			caseName: "run_under:kuma,kong-mesh",
			results: map[meshdetect.MeshKind]*meshdetect.RunUnderResults{
				meshdetect.MeshKindKuma: {
					PodOrServiceAnnotation:   true,
					SidecarContainerInjected: true,
				},
				meshdetect.MeshKindKongMesh: {
					PodOrServiceAnnotation:   true,
					SidecarContainerInjected: true,
				},
			},
			serializedStr: "kinm=\"k2,k3,km2,km3\"",
		},
		{
			caseName: "run_under:traefik,aws-app-mesh",
			results: map[meshdetect.MeshKind]*meshdetect.RunUnderResults{
				meshdetect.MeshKindTraefik: {
					PodOrServiceAnnotation: true,
				},
				meshdetect.MeshKindAWSAppMesh: {
					PodOrServiceAnnotation:   true,
					SidecarContainerInjected: true,
					InitContainerInjected:    true,
				},
			},
			serializedStr: "kinm=\"a2,a3,a4,t2\"",
		},
		{
			caseName:      "run_under:should return empty string for nil results",
			results:       nil,
			serializedStr: "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.caseName, func(t *testing.T) {
			serialized := serializeMeshRunUnderResults(tc.results)
			require.Equalf(t, tc.serializedStr, serialized,
				"case %s: serialized message should be the same as expected", tc.caseName)
		})
	}
}

func TestSerializeMeshServiceDistribution(t *testing.T) {
	testCases := []struct {
		caseName      string
		results       *meshdetect.ServiceDistributionResults
		serializedStr string
	}{
		{
			caseName: "service_distribution:istio=32,kuma=50,kong-mesh=50,traefik=20",
			results: &meshdetect.ServiceDistributionResults{
				TotalServices: 234,
				MeshDistribution: map[meshdetect.MeshKind]int{
					meshdetect.MeshKindIstio:    32,
					meshdetect.MeshKindKuma:     50,
					meshdetect.MeshKindKongMesh: 50,
					meshdetect.MeshKindTraefik:  20,
				},
			},
			serializedStr: "mdist=\"all234,i32,k50,km50,t20\"",
		},
		{
			caseName:      "service_distribution:should return empty string for nil results",
			results:       nil,
			serializedStr: "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.caseName, func(t *testing.T) {
			serialized := serializeMeshServiceDistribution(tc.results)
			require.Equalf(t, tc.serializedStr, serialized,
				"case %s: serialized message should be the same as expected", tc.caseName)
		})

	}
}
