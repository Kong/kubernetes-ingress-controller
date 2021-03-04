package admission

import (
	"testing"

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
