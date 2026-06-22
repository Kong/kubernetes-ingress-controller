package kongintegration

import (
	"fmt"
	"os"
	"testing"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

func TestMain(m *testing.M) {
	if testenv.IsKongGatewayVersionEnterpriseOnly() && !testenv.KongEnterpriseEnabled() {
		fmt.Println("INFO: skipping suite, because Kong Gateway >= 3.15 is enterprise only")
		os.Exit(0)
	}
	os.Exit(m.Run())
}
