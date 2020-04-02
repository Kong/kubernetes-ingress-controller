package v1

import (
	"reflect"
	"testing"
)

func TestConfiguration_DeepCopyInto(t *testing.T) {
	type args struct {
		out *Configuration
	}
	tests := []struct {
		name string
		in   *Configuration
		args args
	}{
		{
			in: &Configuration{},
			args: args{
				out: &Configuration{},
			},
		},
		{
			in: &Configuration{
				"foo": "bar",
			},
			args: args{
				out: &Configuration{
					"foo": "bar",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Configuration
			tt.in.DeepCopyInto(&got)
			if !reflect.DeepEqual(&got, tt.args.out) {
				t.Errorf("Configuration.DeepCopy() = %v, want %v", got, tt.args.out)
			}
		})
	}
}
