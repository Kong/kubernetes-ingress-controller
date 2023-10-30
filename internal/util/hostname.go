package util

import (
	"strings"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// HostnamesIntersect checks if the hostnameA and hostnameB have an intersection.
// To perform this check, the function HostnamesMatch is called twice swapping the
// parameters and using first hostnameA as a mask, then hostnameB.
// If there is at least one match, the hostnames intersect.
func HostnamesIntersect[H1, H2 gatewayapi.HostnameT](hostnameA H1, hostnameB H2) bool {
	var (
		a = (string)(hostnameA)
		b = (string)(hostnameB)
	)
	return HostnamesMatch(a, b) || HostnamesMatch(b, a)
}

// HostnamesMatch checks that the hostnameB matches the hostnameA. HostnameA is treated as mask
// to be checked against the hostnameB.
func HostnamesMatch(hostnameA, hostnameB string) bool {
	// the hostnames are in the form of "foo.bar.com"; split them
	// in a slice of substrings
	hostnameALabels := strings.Split(hostnameA, ".")
	hostnameBLabels := strings.Split(hostnameB, ".")

	var a, b int
	var wildcard bool

	// iterate over the parts of both the hostnames
	for a, b = 0, 0; a < len(hostnameALabels) && b < len(hostnameBLabels); a, b = a+1, b+1 {
		var matchFound bool

		// if the current part of B is a wildcard, we need to find the first
		// A part that matches with the following B part
		if wildcard {
			for ; b < len(hostnameBLabels); b++ {
				if hostnameALabels[a] == hostnameBLabels[b] {
					matchFound = true
					break
				}
			}
		}

		// if no match was found, the hostnames don't match
		if wildcard && !matchFound {
			return false
		}

		// check if at least on of the current parts are a wildcard; if so, continue
		if hostnameALabels[a] == "*" {
			wildcard = true
			continue
		}
		// reset the wildcard  variables
		wildcard = false

		// if the current a part is different from the b part, the hostnames are incompatible
		if hostnameALabels[a] != hostnameBLabels[b] {
			return false
		}
	}
	return len(hostnameBLabels)-b == len(hostnameALabels)-a
}
