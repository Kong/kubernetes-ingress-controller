package utils

import (
	"fmt"
	"testing"

	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestUpdateLoadBalancerIngress(t *testing.T) {
	const (
		oldIP       = "1.2.3.4"
		oldHostname = "example"
	)

	oldIngress := func() []client.Object {
		return []client.Object{
			&netv1.Ingress{
				Status: netv1.IngressStatus{
					LoadBalancer: netv1.IngressLoadBalancerStatus{
						Ingress: []netv1.IngressLoadBalancerIngress{
							{
								IP:       oldIP,
								Hostname: oldHostname,
							},
						},
					},
				},
			},
			&netv1.Ingress{
				Status: netv1.IngressStatus{
					LoadBalancer: netv1.IngressLoadBalancerStatus{
						Ingress: []netv1.IngressLoadBalancerIngress{
							{
								IP:       oldIP,
								Hostname: oldHostname,
							},
						},
					},
				},
			},
			&kongv1beta1.TCPIngress{
				Status: kongv1beta1.TCPIngressStatus{
					LoadBalancer: corev1.LoadBalancerStatus{
						Ingress: []corev1.LoadBalancerIngress{
							{
								IP:       oldIP,
								Hostname: oldHostname,
							},
						},
					},
				},
			},
			&kongv1beta1.UDPIngress{
				Status: kongv1beta1.UDPIngressStatus{
					LoadBalancer: corev1.LoadBalancerStatus{
						Ingress: []corev1.LoadBalancerIngress{
							{
								IP:       oldIP,
								Hostname: oldHostname,
							},
						},
					},
				},
			},
		}
	}

	t.Run("update not needed as no changes detected", func(t *testing.T) {
		oldIngress := oldIngress()
		newAddresses := []netv1.IngressLoadBalancerIngress{
			{
				IP:       oldIP,
				Hostname: oldHostname,
			},
		}

		for _, old := range oldIngress {
			t.Run(fmt.Sprintf("%T", old), func(t *testing.T) {
				copiedOld := old.DeepCopyObject()
				updatedNeeded, err := UpdateLoadBalancerIngress(old, newAddresses)
				require.NoError(t, err)
				assert.False(t, updatedNeeded)
				assert.Equal(t, copiedOld, old, "when update not needed, the old object shouldn't be updated")
			})
		}
	})

	t.Run("update needed", func(t *testing.T) {
		oldIngress := oldIngress()
		const (
			newIP = "192.168.1.1"
		)
		newAddresses := []netv1.IngressLoadBalancerIngress{
			{
				IP:       newIP,
				Hostname: oldHostname,
			},
		}

		for _, old := range oldIngress {
			t.Run(fmt.Sprintf("%T", old), func(t *testing.T) {
				copiedOld := old.DeepCopyObject()
				updatedNeeded, err := UpdateLoadBalancerIngress(old, newAddresses)
				require.NoError(t, err)
				assert.True(t, updatedNeeded)
				assert.NotEqual(t, copiedOld, old, "when updated needed, the old object should be updated")
			})
		}
	})

	t.Run("unknown type passed should not panic", func(t *testing.T) {
		unknownObject := &corev1.Pod{}
		newAddresses := []netv1.IngressLoadBalancerIngress{
			{
				IP:       oldIP,
				Hostname: oldHostname,
			},
		}
		var err error
		assert.NotPanics(t, func() {
			_, err = UpdateLoadBalancerIngress(unknownObject, newAddresses)
		})
		require.ErrorContains(t, err, "unsupported ingress type")
	})
}
