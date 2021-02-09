package networking

import (
	"encoding/binary"
	"fmt"
	"net"
)

// -----------------------------------------------------------------------------
// Public Functions - Helper
// -----------------------------------------------------------------------------

// ConvertIPv4ToUint32 converts an IPv4 net.IP to a uint32
// FIXME: this does nothing to protect the caller from bad input yet
func ConvertIPv4ToUint32(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

// ConvertIPv4ToUint32 converts an IPv4 net.IP to a uint32
// FIXME: this does nothing to protect the caller from bad input yet
func ConvertUint32ToIPv4(i uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, i)
	return ip
}

// GetIPRangeStr provides a string range of IP address given two net.IPs.
// For example, "192.168.1.240-192.168.1.250".
func GetIPRangeStr(ip1, ip2 net.IP) string {
	return fmt.Sprintf("%s-%s", ip1, ip2)
}
