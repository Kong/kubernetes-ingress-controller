//+build integration_tests

package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
	"github.com/kong/kubernetes-testing-framework/pkg/runbooks"
	"github.com/stretchr/testify/assert"
)

func TestIngress(t *testing.T) {
	assert.NoError(t, runbooks.DeployIngressForContainer(kc, "kong", "/nginx", k8s.NewContainer("nginx", "nginx", 80)))
	assert.Eventually(t, func() bool {
		resp, err := http.Get(fmt.Sprintf("%s/nginx", proxyURL().String()))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, time.Minute*5, time.Second*1)
}
