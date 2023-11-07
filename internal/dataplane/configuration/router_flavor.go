package configuration

// RouterFlavor is the type for Kong Gateway router flavors.
// Ref: https://docs.konghq.com/kubernetes-ingress-controller/latest/references/supported-router-flavors
type RouterFlavor string

const (
	// RouterFlavorTraditional is one of Kong Gateway router flavors.
	RouterFlavorTraditional RouterFlavor = "traditional"
	// RouterFlavorTraditionalCompatible is one of Kong Gateway router flavors.
	RouterFlavorTraditionalCompatible RouterFlavor = "traditional_compatible"
	// RouterFlavorExpressions is one of Kong Gateway router flavors.
	// Ref: https://docs.konghq.com/gateway/latest/reference/router-expressions-language/
	RouterFlavorExpressions RouterFlavor = "expressions"
)
