package admission

import (
	"context"
	"testing"

	"github.com/kong/go-kong/kong"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestKongHTTPValidator_ValidateCredential(t *testing.T) {
	type args struct {
		secret corev1.Secret
	}
	tests := []struct {
		name        string
		args        args
		wantOK      bool
		wantMessage string
		wantErr     bool
	}{
		{
			name: "valid key-auth credential",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"key":          []byte("foo"),
						"kongCredType": []byte("key-auth"),
					},
				},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name: "valid jwt credential",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"algorithm":      []byte("foo-algorithm"),
						"key":            []byte("foo-key"),
						"secret":         []byte("foo-secret"),
						"rsa_public_key": []byte("my-key"),
						"kongCredType":   []byte("key-auth"),
					},
				},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name: "valid jwt credential without rsa_public_key",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"algorithm":    []byte("foo-algorithm"),
						"key":          []byte("foo-key"),
						"secret":       []byte("foo-secret"),
						"kongCredType": []byte("key-auth"),
					},
				},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name: "valid keyauth_credential",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"key":          []byte("foo"),
						"kongCredType": []byte("keyauth_credential"),
					},
				},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name: "invalid key-auth credential",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"key-wrong":    []byte("foo"),
						"kongCredType": []byte("key-auth"),
					},
				},
			},
			wantOK:      false,
			wantMessage: "missing required field(s): key",
			wantErr:     false,
		},
		{
			name: "valid mtls-auth credential",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"subject_name": []byte("foo"),
						"kongCredType": []byte("mtls-auth"),
					},
				},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name: "invalid mtls-auth credential",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"kongCredType": []byte("mtls-auth"),
					},
				},
			},
			wantOK:      false,
			wantMessage: "missing required field(s): subject_name",
			wantErr:     false,
		},
		{
			name: "invalid credential type",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"kongCredType": []byte("foo"),
					},
				},
			},
			wantOK:      false,
			wantMessage: "invalid credential type: foo",
			wantErr:     false,
		},
		{
			name: "non-kong secrets are passed",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"key": []byte("foo"),
					},
				},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := KongHTTPValidator{}
			got, got1, err := validator.ValidateCredential(tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("KongHTTPValidator.ValidateCredential() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantOK {
				t.Errorf("KongHTTPValidator.ValidateCredential() got = %v, want %v", got, tt.wantOK)
			}
			if got1 != tt.wantMessage {
				t.Errorf("KongHTTPValidator.ValidateCredential() got1 = %v, want %v", got1, tt.wantMessage)
			}
		})
	}
}

type fakeConsumerSvc struct {
	kong.AbstractConsumerService

	consumer *kong.Consumer
	err      error
}

func (f *fakeConsumerSvc) Get(ctx context.Context, usernameOrID *string) (*kong.Consumer, error) {
	return f.consumer, f.err
}

func TestKongHTTPValidator_ValidateConsumer(t *testing.T) {
	for _, tt := range []struct {
		name        string
		ConsumerSvc kong.AbstractConsumerService

		in configurationv1.KongConsumer

		wantSuccess   bool
		wantErrorText string
		wantErr       bool
	}{
		{
			name:          "empty username",
			in:            configurationv1.KongConsumer{},
			wantSuccess:   false,
			wantErrorText: ErrTextUsernameEmpty,
		},
		{
			name:        "kong says consumer not found",
			ConsumerSvc: &fakeConsumerSvc{err: kong.NewAPIError(404, "")},
			in:          configurationv1.KongConsumer{Username: "something"},
			wantSuccess: true,
		},
		{
			name:        "kong says HTTP 500",
			ConsumerSvc: &fakeConsumerSvc{err: kong.NewAPIError(500, "")},
			in:          configurationv1.KongConsumer{Username: "something"},
			wantSuccess: false,
			wantErr:     true,
		},
		{
			name:          "consumer already exists",
			ConsumerSvc:   &fakeConsumerSvc{consumer: &kong.Consumer{}},
			in:            configurationv1.KongConsumer{Username: "something"},
			wantSuccess:   false,
			wantErrorText: ErrTextConsumerExists,
		},
		{
			name:        "validation successful",
			ConsumerSvc: &fakeConsumerSvc{},
			in:          configurationv1.KongConsumer{Username: "something"},
			wantSuccess: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			v := KongHTTPValidator{
				ConsumerSvc: tt.ConsumerSvc,
				Logger:      logrus.New(),
			}
			gotSuccess, gotErrorText, gotErr := v.ValidateConsumer(context.Background(), tt.in)

			require.Equal(t, tt.wantSuccess, gotSuccess)
			require.Equal(t, tt.wantErrorText, gotErrorText)
			if tt.wantErr {
				require.Error(t, gotErr)
			} else {
				require.NoError(t, gotErr)
			}
		})
	}
}
