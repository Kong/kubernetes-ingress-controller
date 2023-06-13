package builder

import (
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ServicePortBuilder is a builder for core v1 ServicePort.
// Primarily used for testing.
type ServicePortBuilder struct {
	sp corev1.ServicePort
}

func NewServicePort() *ServicePortBuilder {
	return &ServicePortBuilder{
		sp: corev1.ServicePort{},
	}
}

// WithNodePort sets the target port on the service port.
func (b *ServicePortBuilder) WithNodePort(port int32) *ServicePortBuilder {
	b.sp.NodePort = port
	return b
}

// WithTargetPort sets the target port on the service port.
func (b *ServicePortBuilder) WithTargetPort(targetport intstr.IntOrString) *ServicePortBuilder {
	b.sp.TargetPort = targetport
	return b
}

// WithPort sets the port on the service port.
func (b *ServicePortBuilder) WithPort(port int32) *ServicePortBuilder {
	b.sp.Port = port
	return b
}

// WithAppProtocol sets the app protocol on the service port.
func (b *ServicePortBuilder) WithAppProtocol(appproto string) *ServicePortBuilder {
	b.sp.AppProtocol = lo.ToPtr(appproto)
	return b
}

// WithProtocol sets the protocol on the service port.
func (b *ServicePortBuilder) WithProtocol(proto corev1.Protocol) *ServicePortBuilder {
	b.sp.Protocol = proto
	return b
}

// WithName sets the name on the service port.
func (b *ServicePortBuilder) WithName(name string) *ServicePortBuilder {
	b.sp.Name = name
	return b
}

// Build returns the configured ServicePort.
func (b *ServicePortBuilder) Build() corev1.ServicePort {
	return b.sp
}

// IntoSlice returns the configured ServicePort in a slice.
func (b *ServicePortBuilder) IntoSlice() []corev1.ServicePort {
	return []corev1.ServicePort{b.sp}
}
