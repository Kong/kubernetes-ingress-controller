package util

import "fmt"

// EnablementStatus can be 'enabled', 'disabled' or 'auto'.
type EnablementStatus int

const (
	// EnablementStatusDisabled says that the resource it controls is disabled.
	EnablementStatusDisabled EnablementStatus = iota
	// EnablementStatusEnabled says that the resource it controls is enabled.
	EnablementStatusEnabled EnablementStatus = iota
	// EnablementStatusAuto says that whether the resource it controls is enabled
	// or disabled should be decided upon by automation.
	EnablementStatusAuto EnablementStatus = iota
)

// String converts EnablementStatus to a lowercase word.
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

// Set sets the value of the EnablementStatus to match the provided string value.
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

func (e *EnablementStatus) Type() string {
	return "EnablementStatus"
}
