package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnablementStatus_Set(t *testing.T) {
	for _, tt := range []struct {
		input   string
		want    EnablementStatus
		wantErr bool
	}{
		{
			input: "enabled",
			want:  EnablementStatusEnabled,
		},
		{
			input: "disabled",
			want:  EnablementStatusDisabled,
		},
		{
			input: "auto",
			want:  EnablementStatusAuto,
		},
		{
			wantErr: true,
		},
		{
			input:   "blah",
			wantErr: true,
		},
	} {
		t.Run(tt.input, func(t *testing.T) {
			var got EnablementStatus
			err := got.Set(tt.input)
			require.Equal(t, tt.want, got)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
