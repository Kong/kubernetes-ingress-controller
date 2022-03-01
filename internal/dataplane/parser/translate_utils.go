package parser

import (
	"fmt"
	"strings"

	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// kongHeaderRegexPrefix is a reserved prefix string that Kong uses to determine if it should parse a header value
// as a regex
const kongHeaderRegexPrefix = "~*"

// -----------------------------------------------------------------------------
// Translate Utilities - Gateway
// -----------------------------------------------------------------------------

// convertGatewayMatchHeadersToKongRouteMatchHeaders takes an input list of Gateway APIs HTTPHeaderMatch
// and converts these header matching rules to the format expected by go-kong.
func convertGatewayMatchHeadersToKongRouteMatchHeaders(headers []gatewayv1alpha2.HTTPHeaderMatch) (map[string][]string, error) {
	// iterate through each provided header match checking for invalid
	// options and otherwise converting to kong type format.
	convertedHeaders := make(map[string][]string)
	for _, header := range headers {
		if header.Type != nil && *header.Type == gatewayv1alpha2.HeaderMatchRegularExpression {
			convertedHeaders[string(header.Name)] = []string{kongHeaderRegexPrefix + header.Value}
		} else if header.Type == nil || *header.Type == gatewayv1alpha2.HeaderMatchExact {

			// split the header values by comma
			values := strings.Split(header.Value, ",")
			convertedHeaders[string(header.Name)] = values
		} else {
			return nil, fmt.Errorf("unknown/unsupported header match type: %s", string(*header.Type))
		}
	}

	return convertedHeaders, nil
}
