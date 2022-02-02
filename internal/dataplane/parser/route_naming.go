package parser

import (
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc64"
	"regexp"

	"k8s.io/apimachinery/pkg/types"
)

// -----------------------------------------------------------------------------
// Route Naming - Public Vars & Consts
// -----------------------------------------------------------------------------

const (
	// IngressV1RoutePrefix indicates the string to expect as the prefix
	// for the name of a kong.Route object when that route object belongs
	// to a networking/v1 Ingress object.
	IngressV1RoutePrefix = "ingressv1"

	// IngressV1Beta1RoutePrefix indicates the string to expect as the prefix
	// for the name of a kong.Route object when that route object belongs
	// to a networking/v1beta1 Ingress object.
	IngressV1Beta1RoutePrefix = "ingressv1beta1"

	// TCPIngressV1Beta1RoutePrefix indicates the string to expect as the prefix
	// for the name of a kong.Route object when that route object belongs
	// to a kong/v1beta1 TCPIngress object.
	TCPIngressV1Beta1RoutePrefix = "tcpingressv1beta1"

	// UDPIngressV1Beta1RoutePrefix indicates the string to expect as the prefix
	// for the name of a kong.Route object when that route object belongs
	// to a kong/v1beta1 UDPIngress object.
	UDPIngressV1Beta1RoutePrefix = "udpingressv1beta1"

	// KnativeIngressV1Alpha1RoutePrefix indicates the string to expect as the
	// prefix for the name of a kong.Route object when that route object belongs
	// to a knative/v1alpha1 Ingress object.
	KnativeIngressV1Alpha1RoutePrefix = "knativeingressv1alpha1"

	// HTTPRouteV1Alpha2RoutePrefix indicates the string to expect as the
	// prefix for the name of a kong.Route object when that route object belongs
	// to a networking/v1alpha2 HTTPRoute object.
	HTTPRouteV1Alpha2RoutePrefix = "httproutev1alpha2"

	// UnknownRouteType is provided as the Kubernetes type for Kong Routes whose
	// names use the legacy naming convention and can not be dinstinctly linked
	// to the type of object they reference.
	UnknownRouteType = "unknown"
)

// -----------------------------------------------------------------------------
// Route Naming - Public Functions
// -----------------------------------------------------------------------------

// GetKubernetesObjectReferenceForKongRouteName produces an object type as well as a
// namespace and name reference to a Kubernetes object which a Kong Route belongs to
// give the name of the route.
//
// The object type can return as "unknown" when using the legacy route naming convention.
// This will mean it is either a net/v1 Ingress, a net/v1beta1 Ingress or a kong/v1beta1
// TCPIngress object which the caller will need to determine.
func GetKubernetesObjectReferenceForKongRouteName(routeName string) (string, types.NamespacedName, error) {
	// check for the default legacy naming convention first
	// as this is the default setting currently.
	usesLegacyNaming, objectType, nsn := isRouteUsingLegacyNaming(routeName)
	if usesLegacyNaming {
		return objectType, nsn, nil
	}

	// if the legacy naming convention isn't in use, process the route
	// name using the latest hashed based naming scheme.
	usesHashedNaming, objectType, nsn := isRouteUsingHashedNames(routeName)
	if !usesHashedNaming {
		return "", types.NamespacedName{}, fmt.Errorf("invalid route name %s", routeName)
	}
	return objectType, nsn, nil
}

// -----------------------------------------------------------------------------
// Route Naming - Naming Convention Regexes
// -----------------------------------------------------------------------------

// routeNameRegexp is a regex with 4 groups which is used to identify
//  - API type
//  - Namespace
//  - Name
//  - UID
// within a Kong Route name and can be used to identify the related
// Kubernetes Ingress objects to which the route belongs.
var routeNameRegexp = regexp.MustCompile(`^([a-zA-Z0-9]+)\.([^.]+)\.(.*)\.([a-zA-Z0-9]+)$`)

// legacyRouteNameRegexp is a regex which handles route names using
// the historical <namespace>.<name>.<number>.<type> format (where
// legacy Ingress, Ingress, and TCPIngress don't have the type).
var legacyRouteNameRegexp = regexp.MustCompile(`^([^.]+)\.(.*)\.[0-9]+(\.[a-zA-Z]+)?$`)

// -----------------------------------------------------------------------------
// Route Naming - Private Functions
// -----------------------------------------------------------------------------

var crcTable = crc64.MakeTable(crc64.ISO)

// getUniqIDForRouteConfig provides a deterministic and unique identifier
// string for the provided route configuration object (Ingress, HTTPRoute,
// e.t.c.). The provided config must be an object that can be serialized
// with json.Marshal().
//
// The resulting uniquely identifying string is a base32 hex encoded string
// from a crc64 hash of the objects serialized contents, without padding.
func getUniqIDForRouteConfig(config interface{}) (string, error) {
	// perform JSON serialization to get the unique string to encode
	routeHashInput, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("could not marshal input for route hash: %w", err)
	}

	// hash the serialized route configuration data
	routeHash := crc64.Checksum(routeHashInput, crcTable)

	// convert the uint64 bytes to []byte
	routeHashBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(routeHashBytes, routeHash)

	// return a string encoded version of the hashed match data
	return base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(routeHashBytes), nil
}

// isRouteUsingLegacyNaming indicates whether the provided route name uses the legacy
// kong Route naming convention (e.g. matches the legacyRouteNameRegexp) and if true
// returns the object type and namespace name reference to the related Kubernetes object.
//
// Note that legacy naming can only produce an object type for a few select types,
// all others will report UnknownRouteType.
func isRouteUsingLegacyNaming(routeName string) (matches bool, objectType string, nsn types.NamespacedName) {
	foundMatches := legacyRouteNameRegexp.FindAllStringSubmatch(routeName, -1)
	if len(foundMatches) != 1 {
		return
	}
	matches = true
	objectType = UnknownRouteType

	// determine the route type if possible, upgrade known type suffixes
	submatches := foundMatches[0]
	switch submatches[3] {
	case ".udp":
		objectType = UDPIngressV1Beta1RoutePrefix
	case ".httproute":
		objectType = HTTPRouteV1Alpha2RoutePrefix
	}

	nsn = types.NamespacedName{
		Namespace: submatches[1],
		Name:      submatches[2],
	}

	return
}

// isRouteUsingHashedNames indicates whether the provided route name uses the hashed
// kong Route naming convention (e.g. matches the routeNameRegexp) and if true returns
// the object type and namespace name referece to the related Kubernetes object.
func isRouteUsingHashedNames(routeName string) (matches bool, objectType string, nsn types.NamespacedName) {
	foundMatches := routeNameRegexp.FindAllStringSubmatch(routeName, -1)
	if len(foundMatches) != 1 {
		return
	}
	matches = true

	submatches := foundMatches[0]
	objectType = submatches[1]
	nsn = types.NamespacedName{
		Namespace: submatches[2],
		Name:      submatches[3],
	}

	return
}
