package proxy

import (
	"testing"

	"github.com/kong/kubernetes-ingress-controller/pkg/store"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_networkingIngressV1Beta1(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		secret string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "valid secret",
			args: args{
				secret: "default/validCustomEntities",
			},
			want:    []byte("carp"),
			wantErr: true,
		},
		{
			name: "incorrect name format",
			args: args{
				secret: "!",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "non-existent secret",
			args: args{
				secret: "default/nope",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "secret lacks config key",
			args: args{
				secret: "default/invalidCustomEntities",
			},
			want:    nil,
			wantErr: true,
		},
	}
	store, err := store.NewFakeStore(store.FakeObjects{
		Secrets: []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "validCustomEntities",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"config": []byte("carp"),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "invalidCustomEntities",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"ohno": []byte("carp"),
				},
			},
		},
	})
	assert.Nil(err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchCustomEntities(tt.args.secret, store)
			if err != nil && !tt.wantErr {
				t.Errorf("kongPluginFromK8SClusterPlugin error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(tt.want, got)
		})
	}
}
