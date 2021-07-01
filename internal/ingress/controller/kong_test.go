package controller

import (
	"reflect"
	"testing"

	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
)

func Test_getIngressControllerTags(t *testing.T) {
	type args struct {
		config sendconfig.Kong
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "configuration with tag support and filter tags",
			args: args{
				config: sendconfig.Kong{
					DeprecatedHasTagSupport: true,
					FilterTags:              []string{"foo-tag", "bar-tag"},
				},
			},
			want: []string{"foo-tag", "bar-tag"},
		},
		{
			name: "configuratiion with tag support and no filter tags",
			args: args{
				config: sendconfig.Kong{
					DeprecatedHasTagSupport: true,
					FilterTags:              []string{},
				},
			},
			want: nil,
		}, {
			name: "configuration with no tag support",
			args: args{
				config: sendconfig.Kong{
					DeprecatedHasTagSupport: false,
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getIngressControllerTags(tt.args.config)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getIngressControllerTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
