package controller

import (
	"context"
	"reflect"
	"testing"

	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser/kongstate"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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

func Test_toDeckContent(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   kongstate.KongState
		want file.Content
	}{
		{
			name: "sorts credentials consistently",
			in: kongstate.KongState{
				Consumers: []kongstate.Consumer{
					{
						KeyAuths: map[string]*kong.KeyAuth{
							"a": {Key: kong.String("key-22")},
							"b": {Key: kong.String("key-11")},
							"c": {Key: kong.String("key-33")},
						},
						HMACAuths: map[string]*kong.HMACAuth{
							"a": {Username: kong.String("hmac-22")},
							"b": {Username: kong.String("hmac-11")},
							"c": {Username: kong.String("hmac-33")},
						},
						JWTAuths: map[string]*kong.JWTAuth{
							"a": {Key: kong.String("jwt-22")},
							"b": {Key: kong.String("jwt-11")},
							"c": {Key: kong.String("jwt-33")},
						},
						BasicAuths: map[string]*kong.BasicAuth{
							"a": {Username: kong.String("basic-22")},
							"b": {Username: kong.String("basic-11")},
							"c": {Username: kong.String("basic-33")},
						},
						Oauth2Creds: map[string]*kong.Oauth2Credential{
							"a": {ClientID: kong.String("oauth2-22")},
							"b": {ClientID: kong.String("oauth2-11")},
							"c": {ClientID: kong.String("oauth2-33")},
						},
					},
				},
			},
			want: file.Content{
				FormatVersion: FormatVersion,
				Consumers: []file.FConsumer{
					{
						KeyAuths: []*kong.KeyAuth{
							{Key: kong.String("key-11")},
							{Key: kong.String("key-22")},
							{Key: kong.String("key-33")},
						},
						HMACAuths: []*kong.HMACAuth{
							{Username: kong.String("hmac-11")},
							{Username: kong.String("hmac-22")},
							{Username: kong.String("hmac-33")},
						},
						JWTAuths: []*kong.JWTAuth{
							{Key: kong.String("jwt-11")},
							{Key: kong.String("jwt-22")},
							{Key: kong.String("jwt-33")},
						},
						BasicAuths: []*kong.BasicAuth{
							{Username: kong.String("basic-11")},
							{Username: kong.String("basic-22")},
							{Username: kong.String("basic-33")},
						},
						Oauth2Creds: []*kong.Oauth2Credential{
							{ClientID: kong.String("oauth2-11")},
							{ClientID: kong.String("oauth2-22")},
							{ClientID: kong.String("oauth2-33")},
						},
					},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			n := KongController{
				cfg:    &Configuration{},
				Logger: logrus.New(),
			}
			got := n.toDeckContent(context.Background(), &tt.in)
			assert.Equal(t, tt.want, *got)
		})
	}
}
