package configuration

import (
	"testing"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
)

func TestEnforceKongUpstreamPolicyStatus(t *testing.T) {
	testCases := []struct {
		name                             string
		kongupstreamPolicy               kongv1beta1.KongUpstreamPolicy
		inputObjects                     []client.Object
		expectedkongUpstreamPolicyStatus gatewayapi.Policystatus
	}{}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// test function here
		})
	}
}
