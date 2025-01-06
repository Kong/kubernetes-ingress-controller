package kongstate

import (
	"github.com/kong/go-kong/kong"
)

// License represents the license object in Kong.
type License struct {
	kong.License
}

// SanitizedCopy returns a shallow copy with sensitive values redacted best-effort.
func (l License) SanitizedCopy() License {
	return License{
		License: kong.License{
			ID:        l.ID,
			Payload:   redactedString,
			CreatedAt: l.CreatedAt,
			UpdatedAt: l.UpdatedAt,
		},
	}
}
