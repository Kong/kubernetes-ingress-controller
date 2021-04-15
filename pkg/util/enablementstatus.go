package util

import "fmt"

// EnablementStatus can be 'enabled', 'disabled' or 'auto'.
type EnablementStatus int

const (
	EnablementStatusDisabled EnablementStatus = iota
	EnablementStatusEnabled  EnablementStatus = iota
	EnablementStatusAuto     EnablementStatus = iota
)

const EnablementStatusUsageString = "can be one of [enabled, disabled, auto]"

func (e *EnablementStatus) String() string {
	switch *e {
	case EnablementStatusDisabled:
		return "disabled"
	case EnablementStatusEnabled:
		return "enabled"
	case EnablementStatusAuto:
		return "auto"
	default:
		panic(fmt.Sprintf("unknown EnablementStatus value %v", *e))
	}
}

func (e *EnablementStatus) Set(s string) error {
	for _, val := range []EnablementStatus{
		EnablementStatusDisabled, EnablementStatusEnabled, EnablementStatusAuto,
	} {
		if s == val.String() {
			*e = val
			return nil
		}
	}

	return fmt.Errorf("%q is not a valid EnablementStatus", s)
}
