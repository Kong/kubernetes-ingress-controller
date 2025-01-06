package dataplane

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	netv1 "k8s.io/api/networking/v1"
)

func TestAddressFinder(t *testing.T) {
	t.Log("generating a new AddressFinder")
	finder := NewAddressFinder()
	require.NotNil(t, finder)
	require.Nil(t, finder.addressGetter)

	t.Log("verifying that a finder with no overrides or getter produces an error")
	ctx := context.Background()
	addrs, err := finder.GetAddresses(ctx)
	require.Error(t, err)
	require.Empty(t, addrs)
	require.Equal(t, "data-plane addresses can't be retrieved: no valid method available", err.Error())

	t.Log("generating a fake AddressGetter")
	defaultAddrs := []string{"127.0.0.1", "127.0.0.2"}
	overrideAddrs := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}
	fakeGetter := func(_ context.Context) ([]string, error) { return defaultAddrs, nil }

	t.Log("verifying getting a list of addresses from the finder after a getter function is provided")
	finder.SetGetter(fakeGetter)
	addrs, err = finder.GetAddresses(ctx)
	require.NoError(t, err)
	require.Equal(t, defaultAddrs, addrs)

	t.Log("verifying that overrides take precedent over the getter")
	finder.SetOverrides(overrideAddrs)
	addrs, err = finder.GetAddresses(ctx)
	require.NoError(t, err)
	require.Equal(t, overrideAddrs, addrs)

	t.Log("verifying disabling overrides")
	finder.SetOverrides(nil)
	addrs, err = finder.GetAddresses(ctx)
	require.NoError(t, err)
	require.Equal(t, defaultAddrs, addrs)

	t.Log("verifying k8s load balancer formatted version of the addresses")
	lbs, err := finder.GetLoadBalancerAddresses(ctx)
	require.NoError(t, err)
	require.Equal(t, []netv1.IngressLoadBalancerIngress{{IP: defaultAddrs[0]}, {IP: defaultAddrs[1]}}, lbs)

	t.Log("verifying valid DNS names are formatting properly")
	dnsAddrs := []string{"127.0.0.1", "example1.konghq.com", "example2.konghq.com"}
	finder.SetOverrides(dnsAddrs)
	lbs, err = finder.GetLoadBalancerAddresses(ctx)
	require.NoError(t, err)
	require.Equal(t, []netv1.IngressLoadBalancerIngress{
		{IP: dnsAddrs[0]},
		{Hostname: dnsAddrs[1]},
		{Hostname: dnsAddrs[2]},
	}, lbs)

	t.Log("verifying empty addresses return an error")
	finder.SetOverrides([]string{""})
	lbs, err = finder.GetLoadBalancerAddresses(ctx)
	require.Error(t, err)
	require.Empty(t, lbs)
	require.Equal(t, "empty address found", err.Error())

	t.Log("verifying invalid DNS names return an error")
	invalidDNSAddrs := []string{"support@konghq.com"}
	finder.SetOverrides(invalidDNSAddrs)
	lbs, err = finder.GetLoadBalancerAddresses(ctx)
	require.Error(t, err)
	require.Empty(t, lbs)
	require.Equal(t, fmt.Sprintf("%s is not a valid DNS hostname", invalidDNSAddrs[0]), err.Error())
}
