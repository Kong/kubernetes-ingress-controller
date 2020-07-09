package controller

import (
	"reflect"
	"testing"

	"github.com/hbagdi/deck/file"
	"github.com/hbagdi/go-kong/kong"
	"github.com/sirupsen/logrus"
)

func Test_renderConfigWithCustomEntities(t *testing.T) {
	type args struct {
		state                   *file.Content
		customEntitiesJSONBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "basic sanity test for fast-path",
			args: args{
				state: &file.Content{
					FormatVersion: "1.1",
					Services: []file.FService{
						{
							Service: kong.Service{
								Name: kong.String("foo"),
								Host: kong.String("example.com"),
							},
						},
					},
				},
				customEntitiesJSONBytes: nil,
			},
			want:    []byte(`{"_format_version":"1.1","services":[{"host":"example.com","name":"foo"}]}`),
			wantErr: false,
		},
		{
			name: "does not break with random bytes in the custom entities",
			args: args{
				state: &file.Content{
					FormatVersion: "1.1",
					Services: []file.FService{
						{
							Service: kong.Service{
								Name: kong.String("foo"),
								Host: kong.String("example.com"),
							},
						},
					},
				},
				customEntitiesJSONBytes: []byte("random-bytes"),
			},
			want:    []byte(`{"_format_version":"1.1","services":[{"host":"example.com","name":"foo"}]}`),
			wantErr: false,
		},
		{
			name: "custom entities cannot hijack core entities",
			args: args{
				state: &file.Content{
					FormatVersion: "1.1",
					Services: []file.FService{
						{
							Service: kong.Service{
								Name: kong.String("foo"),
								Host: kong.String("example.com"),
							},
						},
					},
				},
				customEntitiesJSONBytes: []byte(`{"services":[{"host":"rogue.example.com","name":"rogue"}]}`),
			},
			want:    []byte(`{"_format_version":"1.1","services":[{"host":"example.com","name":"foo"}]}`),
			wantErr: false,
		},
		{
			name: "custom entities can be populated",
			args: args{
				state: &file.Content{
					FormatVersion: "1.1",
					Services: []file.FService{
						{
							Service: kong.Service{
								Name: kong.String("foo"),
								Host: kong.String("example.com"),
							},
						},
					},
				},
				customEntitiesJSONBytes: []byte(`{"my-custom-dao-name":` +
					`[{"name":"custom1","key1":"value1"},` +
					`{"name":"custom2","dumb":"test-value","boring-test-value-name":"really?"}]}`),
			},
			want: []byte(`{"_format_version":"1.1",` +
				`"my-custom-dao-name":[{"key1":"value1","name":"custom1"},` +
				`{"boring-test-value-name":"really?","dumb":"test-value","name":"custom2"}]` +
				`,"services":[{"host":"example.com","name":"foo"}]}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n KongController
			n.Logger = logrus.New()
			got, err := n.renderConfigWithCustomEntities(tt.args.state, tt.args.customEntitiesJSONBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderConfigWithCustomEntities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("renderConfigWithCustomEntities() = %v, want %v",
					string(got), string(tt.want))
			}
		})
	}
}
