package kongstate

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
)

func TestCertificate_SanitizedCopy(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   Certificate
		want Certificate
	}{
		{
			name: "fills all fields but Consumer and sanitizes key",
			in: Certificate{kong.Certificate{
				ID:        kong.String("1"),
				Cert:      kong.String("2"),
				Key:       kong.String("3"),
				CreatedAt: int64Ptr(4),
				SNIs:      []*string{kong.String("5.1"), kong.String("5.2")},
				Tags:      []*string{kong.String("6.1"), kong.String("6.2")},
			}},
			want: Certificate{kong.Certificate{
				ID:        kong.String("1"),
				Cert:      kong.String("2"),
				Key:       redactedString,
				CreatedAt: int64Ptr(4),
				SNIs:      []*string{kong.String("5.1"), kong.String("5.2")},
				Tags:      []*string{kong.String("6.1"), kong.String("6.2")},
			}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := *tt.in.SanitizedCopy()
			assert.Equal(t, tt.want, got)
		})
	}
}
