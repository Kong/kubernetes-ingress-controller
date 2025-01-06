package config

type RouterFlavor string

const (
	RouterFlavorTraditional           RouterFlavor = "traditional"
	RouterFlavorTraditionalCompatible RouterFlavor = "traditional_compatible"
	RouterFlavorExpressions           RouterFlavor = "expressions"
)

func ShouldEnableExpressionRoutes(rf RouterFlavor) bool {
	return rf == RouterFlavorExpressions
}
