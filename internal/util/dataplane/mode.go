package dataplane

import "fmt"

type DBMode string

const (
	DBModeOff      DBMode = "off"
	DBModePostgres DBMode = "postgres"
)

func NewDBMode(mode string) (DBMode, error) {
	switch mode {
	case "", string(DBModeOff):
		return DBModeOff, nil
	case string(DBModePostgres):
		return DBModePostgres, nil
	}
	return "", fmt.Errorf("unsupported db mode: %q", mode)
}

// IsDBLessMode can be used to detect the proxy mode (db or dbless).
func IsDBLessMode(mode DBMode) bool {
	return mode == "" || mode == DBModeOff
}

// DBBacked returns true if the gateway is DB backed.
// reverse of IsDBLessMode for readability.
func DBBacked(mode DBMode) bool {
	return !IsDBLessMode(mode)
}

type RouterFlavor string

const (
	RouterFlavorTraditional           RouterFlavor = "traditional"
	RouterFlavorTraditionalCompatible RouterFlavor = "traditional_compatible"
	RouterFlavorExpressions           RouterFlavor = "expressions"
)

func ShouldEnableExpressionRoutes(rf RouterFlavor) bool {
	return rf == RouterFlavorExpressions
}
