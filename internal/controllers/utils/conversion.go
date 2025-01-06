package utils

import (
	"fmt"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
)

// UpdateLoadBalancerIngress updates any supported Ingress object with new []netv1.IngressLoadBalancerIngress
// in a backward-compatible fashion if needed. Update does not happen in case there are no changes detected.
func UpdateLoadBalancerIngress(
	ingress client.Object,
	newAddresses []netv1.IngressLoadBalancerIngress,
) (updateNeeded bool, err error) {
	// Convert to netv1 so that we can compare it with newAddresses.
	oldAddresses, err := ingressToNetV1LoadBalancerIngressStatus(ingress)
	if err != nil {
		return false, fmt.Errorf("failed to convert ingress to netv1.Ingress: %w", err)
	}

	updateNeeded = len(oldAddresses) != len(newAddresses) || !reflect.DeepEqual(oldAddresses, newAddresses)
	if !updateNeeded {
		return false, nil
	}

	switch obj := ingress.(type) {
	case *netv1.Ingress:
		obj.Status.LoadBalancer.Ingress = newAddresses
	case *kongv1beta1.TCPIngress:
		obj.Status.LoadBalancer.Ingress = netV1ToCoreV1LoadBalancerIngress(newAddresses)
	case *kongv1beta1.UDPIngress:
		obj.Status.LoadBalancer.Ingress = netV1ToCoreV1LoadBalancerIngress(newAddresses)
	default:
		return false, fmt.Errorf("unsupported ingress type: %T", obj)
	}

	return true, nil
}

func netV1ToCoreV1LoadBalancerIngress(in []netv1.IngressLoadBalancerIngress) []corev1.LoadBalancerIngress {
	out := make([]corev1.LoadBalancerIngress, 0, len(in))
	for _, i := range in {
		out = append(out, corev1.LoadBalancerIngress{
			IP:       i.IP,
			Hostname: i.Hostname,
			// consciously omitting ports as we do not populate them
		})
	}
	return out
}

func ingressToNetV1LoadBalancerIngressStatus(in any) ([]netv1.IngressLoadBalancerIngress, error) {
	switch obj := in.(type) {
	case *netv1.Ingress:
		return obj.Status.LoadBalancer.Ingress, nil
	case *kongv1beta1.TCPIngress:
		return coreV1ToNetV1LoadBalancerIngress(obj.Status.LoadBalancer.Ingress), nil
	case *kongv1beta1.UDPIngress:
		return coreV1ToNetV1LoadBalancerIngress(obj.Status.LoadBalancer.Ingress), nil
	default:
		return nil, fmt.Errorf("unsupported ingress type: %T", obj)
	}
}

func coreV1ToNetV1LoadBalancerIngress(in []corev1.LoadBalancerIngress) []netv1.IngressLoadBalancerIngress {
	out := make([]netv1.IngressLoadBalancerIngress, 0, len(in))
	for _, i := range in {
		out = append(out, netv1.IngressLoadBalancerIngress{
			IP:       i.IP,
			Hostname: i.Hostname,
			// consciously omitting ports as we do not populate them
		})
	}
	return out
}
