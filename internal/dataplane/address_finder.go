package dataplane

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	netv1 "k8s.io/api/networking/v1"
	utilvalidation "k8s.io/apimachinery/pkg/util/validation"
)

// -----------------------------------------------------------------------------
// AddressFinder - Public Types
// -----------------------------------------------------------------------------

// AddressGetter is a function which can dynamically retrieve the list of IPs
// that the data-plane is listening on for ingress network traffic.
type AddressGetter func(ctx context.Context) ([]string, error)

// AddressFinder is a threadsafe metadata object which can provide the current
// live addresses in use by the dataplane at any point in time.
type AddressFinder struct {
	overrideAddresses    []string
	overrideAddressesUDP []string
	addressGetter        AddressGetter

	lock sync.RWMutex
}

// NewAddressFinder provides a new AddressFinder which can be used to find the
// current listening addresses of the data-plane for ingress network traffic.
func NewAddressFinder() *AddressFinder {
	return &AddressFinder{}
}

// -----------------------------------------------------------------------------
// AddressFinder - Public Methods
// -----------------------------------------------------------------------------

// SetGetter provides a callback function that the AddressFinder will use to
// dynamically retrieve the addresses of the data-plane.
func (a *AddressFinder) SetGetter(getter AddressGetter) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.addressGetter = getter
}

// SetOverrides hard codes a specific list of addresses to be the addresses
// that this finder produces for the data-plane. To disable overrides, call
// this method again with an empty list.
func (a *AddressFinder) SetOverrides(addrs []string) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.overrideAddresses = addrs
}

// SetUDPOverrides hard codes a specific list of addresses to be the UDP addresses
// that this finder produces for the data-plane. To disable overrides, call
// this method again with an empty list.
func (a *AddressFinder) SetUDPOverrides(addrs []string) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.overrideAddressesUDP = addrs
}

// GetAddresses provides a list of the addresses which the data-plane is
// listening on for ingress network traffic. Addresses can either be IP
// addresses or hostnames.
func (a *AddressFinder) GetAddresses(ctx context.Context) ([]string, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	if len(a.overrideAddresses) > 0 {
		return a.overrideAddresses, nil
	}

	if a.addressGetter != nil {
		return a.addressGetter(ctx)
	}

	return nil, fmt.Errorf("data-plane addresses can't be retrieved: no valid method available")
}

// GetUDPAddresses provides a list of the UDP addresses which the data-plane is
// listening on for ingress network traffic. Addresses can either be IP
// addresses or hostnames. If UDP settings are not configured, falls back to GetAddresses().
func (a *AddressFinder) GetUDPAddresses(ctx context.Context) ([]string, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	if len(a.overrideAddressesUDP) > 0 {
		return a.overrideAddressesUDP, nil
	}

	if len(a.overrideAddresses) > 0 && a.addressGetter == nil {
		return a.overrideAddresses, nil
	}

	if a.addressGetter != nil {
		return a.addressGetter(ctx)
	}

	return a.GetAddresses(ctx)
}

// GetLoadBalancerAddresses provides a list of the addresses which the
// data-plane is listening on for ingress network traffic, but provides the
// addresses in Kubernetes corev1.LoadBalancerIngress format. Addresses can be
// IP addresses or hostnames.
func (a *AddressFinder) GetLoadBalancerAddresses(ctx context.Context) ([]netv1.IngressLoadBalancerIngress, error) {
	addrs, err := a.GetAddresses(ctx)
	if err != nil {
		return nil, err
	}
	return getAddressHelper(addrs)
}

// getAddressHelper converts a string slice of addresses (IPs or hostnames) into an IngressLoadBalancerIngress
// (https://pkg.go.dev/k8s.io/api/networking/v1#IngressLoadBalancerIngress), or an error if one of the given strings
// is neither a valid IP nor a valid hostname.
func getAddressHelper(addrs []string) ([]netv1.IngressLoadBalancerIngress, error) {
	var loadBalancerAddresses []netv1.IngressLoadBalancerIngress
	for _, addr := range addrs {
		ing := netv1.IngressLoadBalancerIngress{}
		if net.ParseIP(addr) != nil {
			ing.IP = addr
		} else {
			if err := isValidHostname(addr); err != nil {
				return nil, err
			}
			ing.Hostname = addr
		}
		loadBalancerAddresses = append(loadBalancerAddresses, ing)
	}

	return loadBalancerAddresses, nil
}

// GetUDPLoadBalancerAddresses provides a list of the addresses which the
// data-plane is listening on for UDP network traffic, but provides the
// addresses in Kubernetes corev1.LoadBalancerIngress format. Addresses can be
// IP addresses or hostnames.
func (a *AddressFinder) GetUDPLoadBalancerAddresses(ctx context.Context) ([]netv1.IngressLoadBalancerIngress, error) {
	addrs, err := a.GetUDPAddresses(ctx)
	if err != nil {
		return nil, err
	}
	return getAddressHelper(addrs)
}

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

func isValidHostname(hostname string) error {
	if hostname == "" {
		return fmt.Errorf("empty address found")
	}

	var invalid bool
	for _, label := range strings.Split(hostname, ".") {
		if len(utilvalidation.IsDNS1123Label(label)) > 0 {
			invalid = true
		}
	}

	if invalid {
		return fmt.Errorf("%s is not a valid DNS hostname", hostname)
	}

	return nil
}
