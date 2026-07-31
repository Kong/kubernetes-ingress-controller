package kongintegration

import (
	"fmt"
	"os"
	"testing"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

func TestMain(m *testing.M) {
	enterpriseOnly, err := testenv.IsKongGatewayVersionEnterpriseOnly()
	if err != nil {
		fmt.Printf("ERROR: failed to determine Kong Gateway version: %v\n", err)
		os.Exit(1)
	}
	if enterpriseOnly && testenv.KongLicenseData() == "" {
		fmt.Println("ERROR: Kong 3.15+ used and no license provided")
		os.Exit(1)
	}

	os.Exit(m.Run())
}
