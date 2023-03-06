package manager

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

func TestHealthCheckServer(t *testing.T) {
	passChecker := func(req *http.Request) error {
		return nil
	}
	failChecker := func(req *http.Request) error {
		return errors.New("you shall not pass")
	}

	testCases := []struct {
		name               string
		healthzChecker     healthz.Checker
		readyzChecker      healthz.Checker
		livenessProbeCode  int
		readinessProbeCode int
	}{
		{
			name:               "healthz and readyz both pass should return both 200",
			healthzChecker:     passChecker,
			readyzChecker:      passChecker,
			livenessProbeCode:  http.StatusOK,
			readinessProbeCode: http.StatusOK,
		},
		{
			name:               "healthz fail, readyz not set should return 500 & 404",
			healthzChecker:     failChecker,
			livenessProbeCode:  http.StatusInternalServerError,
			readinessProbeCode: http.StatusNotFound,
		},
		{
			name:               "healthz pass, readyz fail should return 200 & 500",
			healthzChecker:     passChecker,
			readyzChecker:      failChecker,
			livenessProbeCode:  http.StatusOK,
			readinessProbeCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			h := &healthCheckHandler{}
			h.setHealthzCheck(tc.healthzChecker)
			h.setReadyzCheck(tc.readyzChecker)
			s := httptest.NewServer(h)
			defer s.Close()

			livenessResp, err := http.Get(s.URL + "/healthz")
			require.NoError(t, err)
			defer livenessResp.Body.Close()
			require.Equal(t, tc.livenessProbeCode, livenessResp.StatusCode)

			readinessResp, err := http.Get(s.URL + "/readyz")
			require.NoError(t, err)
			defer readinessResp.Body.Close()
			require.Equal(t, tc.readinessProbeCode, readinessResp.StatusCode)
		})
	}
}
