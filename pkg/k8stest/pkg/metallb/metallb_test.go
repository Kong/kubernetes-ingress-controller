package metallb_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	ktfmetal "github.com/kong/kubernetes-ingress-controller/pkg/k8stest/pkg/metallb"
	ktfnet "github.com/kong/kubernetes-ingress-controller/pkg/k8stest/pkg/networking"
)

func TestHelperFunctions(t *testing.T) {
	network := net.IPNet{
		IP:   net.IPv4(192, 168, 1, 0),
		Mask: net.IPv4Mask(0, 0, 0, 255),
	}
	ip1, ip2 := ktfmetal.GetIPRangeForMetallb(network)
	assert.Equal(t, ip1.String(), net.IPv4(192, 168, 1, 240).String())
	assert.Equal(t, ip2.String(), net.IPv4(192, 168, 1, 250).String())
	assert.Equal(t, ktfnet.GetIPRangeStr(ip1, ip2), "192.168.1.240-192.168.1.250")
}
