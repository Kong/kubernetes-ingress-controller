package builder

import (
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// ListenerBuilder is a builder for gateway api Listener.
// Will set default values, as specified in the gateway API, for fields that are not set.
// Primarily used for testing.
type ListenerBuilder struct {
	listener gatewayapi.Listener
}

func NewListener(name string) *ListenerBuilder {
	return &ListenerBuilder{
		listener: gatewayapi.Listener{
			Name: gatewayapi.SectionName(name),
		},
	}
}

// Build returns the configured Listener.
func (b *ListenerBuilder) Build() gatewayapi.Listener {
	return b.listener
}

// IntoSlice returns the configured Listener in a slice.
func (b *ListenerBuilder) IntoSlice() []gatewayapi.Listener {
	return []gatewayapi.Listener{b.listener}
}

func (b *ListenerBuilder) WithPort(port int) *ListenerBuilder {
	b.listener.Port = gatewayapi.PortNumber(port)
	return b
}

func (b *ListenerBuilder) HTTP() *ListenerBuilder {
	b.listener.Protocol = gatewayapi.HTTPProtocolType
	return b
}

func (b *ListenerBuilder) HTTPS() *ListenerBuilder {
	b.listener.Protocol = gatewayapi.HTTPSProtocolType
	return b
}

func (b *ListenerBuilder) TLS() *ListenerBuilder {
	b.listener.Protocol = gatewayapi.TLSProtocolType
	return b
}

func (b *ListenerBuilder) TCP() *ListenerBuilder {
	b.listener.Protocol = gatewayapi.TCPProtocolType
	return b
}

func (b *ListenerBuilder) UDP() *ListenerBuilder {
	b.listener.Protocol = gatewayapi.UDPProtocolType
	return b
}

func (b *ListenerBuilder) WithHostname(hostname string) *ListenerBuilder {
	b.listener.Hostname = lo.ToPtr(gatewayapi.Hostname(hostname))
	return b
}

func (b *ListenerBuilder) WithAllowedRoutes(routes *gatewayapi.AllowedRoutes) *ListenerBuilder {
	b.listener.AllowedRoutes = routes
	return b
}

func (b *ListenerBuilder) WithTLSConfig(tlsConfig *gatewayapi.GatewayTLSConfig) *ListenerBuilder {
	b.listener.TLS = tlsConfig
	return b
}
