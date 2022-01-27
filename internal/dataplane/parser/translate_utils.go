package parser

import (
	"fmt"
	"strings"

	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

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
		// TODO: implement regex header matching
		if header.Type != nil && *header.Type == gatewayv1alpha2.HeaderMatchRegularExpression {
			return nil, fmt.Errorf("regular expression header matches are not yet supported")
		}

		// split the header values by comma
		values := strings.Split(header.Value, ",")
		convertedHeaders[string(header.Name)] = values
	}

	return convertedHeaders, nil
}
