package builder

import (
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/address"
)

// ListenerBuilder is a builder for gateway api Listener.
// Will set default values, as specified in the gateway API, for fields that are not set.
// Primarily used for testing.
type ListenerBuilder struct {
	listener gatewayv1beta1.Listener
}

func NewListener(name string) *ListenerBuilder {
	return &ListenerBuilder{
		listener: gatewayv1beta1.Listener{
			Name: gatewayv1beta1.SectionName(name),
		},
	}
}

// Build returns the configured Listener.
func (b *ListenerBuilder) Build() gatewayv1beta1.Listener {
	return b.listener
}

// IntoSlice returns the configured Listener in a slice.
func (b *ListenerBuilder) IntoSlice() []gatewayv1beta1.Listener {
	return []gatewayv1beta1.Listener{b.listener}
}

func (b *ListenerBuilder) WithPort(port int) *ListenerBuilder {
	b.listener.Port = gatewayv1beta1.PortNumber(port)
	return b
}

func (b *ListenerBuilder) HTTP() *ListenerBuilder {
	b.listener.Protocol = gatewayv1beta1.HTTPProtocolType
	return b
}

func (b *ListenerBuilder) HTTPS() *ListenerBuilder {
	b.listener.Protocol = gatewayv1beta1.HTTPSProtocolType
	return b
}

func (b *ListenerBuilder) TLS() *ListenerBuilder {
	b.listener.Protocol = gatewayv1beta1.TLSProtocolType
	return b
}

func (b *ListenerBuilder) TCP() *ListenerBuilder {
	b.listener.Protocol = gatewayv1beta1.TCPProtocolType
	return b
}

func (b *ListenerBuilder) UDP() *ListenerBuilder {
	b.listener.Protocol = gatewayv1beta1.UDPProtocolType
	return b
}

func (b *ListenerBuilder) WithHostname(hostname string) *ListenerBuilder {
	b.listener.Hostname = address.Of(gatewayv1beta1.Hostname(hostname))
	return b
}

func (b *ListenerBuilder) WithAllowedRoutes(routes *gatewayv1beta1.AllowedRoutes) *ListenerBuilder {
	b.listener.AllowedRoutes = routes
	return b
}

func (b *ListenerBuilder) WithTLSConfig(tlsConfig *gatewayv1beta1.GatewayTLSConfig) *ListenerBuilder {
	b.listener.TLS = tlsConfig
	return b
}
