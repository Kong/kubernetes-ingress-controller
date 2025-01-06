package store

import (
	"testing"

	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeyFunc(t *testing.T) {
	type args struct {
		obj interface{}
	}

	type F struct {
		Name      string
		Namespace string
	}
	type B struct {
		F
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			want: "Bar/Foo",
			args: args{
				obj: &F{
					Name:      "Foo",
					Namespace: "Bar",
				},
			},
		},
		{
			want: "Bar/Fu",
			args: args{
				obj: B{
					F: F{
						Name:      "Fu",
						Namespace: "Bar",
					},
				},
			},
		},
		{
			want: "default/foo",
			args: args{
				obj: netv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "default",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := namespacedKeyFunc(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("namespacedKeyFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("namespacedKeyFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
