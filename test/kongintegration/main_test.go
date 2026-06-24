package kongintegration

import (
	"fmt"
	"os"
	"testing"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

func TestMain(m *testing.M) {
	if testenv.IsKongGatewayVersionEnterpriseOnly() && testenv.KongLicenseData() == "" {
		fmt.Println("ERROR: Kong 3.15+ used and no license provided")
		os.Exit(1)
	}

	os.Exit(m.Run())
}
