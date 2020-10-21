package controller

import (
	"reflect"
	"testing"

	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
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

func Test_sortUsername(t *testing.T) {
	username1 := "username1"
	username2 := "username2"
	username3 := "username3"
	username4 := "username4"
	username5 := "username5"

	type args struct {
		content *file.Content
	}
	tests := []struct {
		name string
		args args
		want []file.FConsumer
	}{
		{
			name: "Test Happy Path",
			args: args{
				content: &file.Content{
					Consumers: []file.FConsumer{
						{
							Consumer: kong.Consumer{
								Username: &username1,
							},
						}, {
							Consumer: kong.Consumer{
								Username: &username2,
							},
						}, {
							Consumer: kong.Consumer{
								Username: &username3,
							},
						}, {
							Consumer: kong.Consumer{
								Username: &username4,
							},
						}, {
							Consumer: kong.Consumer{
								Username: &username5,
							},
						},
					},
				},
			},
			want: []file.FConsumer{
				{
					Consumer: kong.Consumer{
						Username: &username5,
					},
				}, {
					Consumer: kong.Consumer{
						Username: &username4,
					},
				}, {
					Consumer: kong.Consumer{
						Username: &username3,
					},
				}, {
					Consumer: kong.Consumer{
						Username: &username2,
					},
				}, {
					Consumer: kong.Consumer{
						Username: &username1,
					},
				},
			},
		}, {
			name: "nil username at end",
			args: args{
				content: &file.Content{
					Consumers: []file.FConsumer{
						{
							Consumer: kong.Consumer{
								Username: &username1,
							},
						}, {
							Consumer: kong.Consumer{
								Username: nil,
							},
						}, {
							Consumer: kong.Consumer{
								Username: &username3,
							},
						}, {
							Consumer: kong.Consumer{
								Username: &username4,
							},
						}, {
							Consumer: kong.Consumer{
								Username: &username5,
							},
						},
					},
				},
			},
			want: []file.FConsumer{
				{
					Consumer: kong.Consumer{
						Username: &username5,
					},
				}, {
					Consumer: kong.Consumer{
						Username: &username4,
					},
				}, {
					Consumer: kong.Consumer{
						Username: &username3,
					},
				}, {
					Consumer: kong.Consumer{
						Username: &username1,
					},
				}, {
					Consumer: kong.Consumer{
						Username: nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortByUsername(tt.args.content)
			// Compare result of function with desired result
			got := tt.args.content.Consumers
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortByUsername() got %v, want %v", got, tt.want)
			}
		})
	}
}
