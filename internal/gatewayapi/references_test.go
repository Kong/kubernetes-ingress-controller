package gatewayapi

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetPermittedForReferenceGrantFrom(t *testing.T) {
	grants := []*ReferenceGrant{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: "fitrat",
			},
			Spec: ReferenceGrantSpec{
				From: []ReferenceGrantFrom{
					{
						Group:     Group("gateway.networking.k8s.io"),
						Kind:      Kind("TCPRoute"),
						Namespace: Namespace("garbage"),
					},
					{
						Group:     Group("gateway.networking.k8s.io"),
						Kind:      Kind("TCPRoute"),
						Namespace: Namespace("behbudiy"),
					},
					{
						Group:     Group("gateway.networking.k8s.io"),
						Kind:      Kind("TCPRoute"),
						Namespace: Namespace("qodiriy"),
					},
				},
				To: []ReferenceGrantTo{
					{
						Group: Group(""),
						Kind:  Kind("GrantOne"),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: "cholpon",
			},
			Spec: ReferenceGrantSpec{
				From: []ReferenceGrantFrom{
					{
						Group:     Group("gateway.networking.k8s.io"),
						Kind:      Kind("UDPRoute"),
						Namespace: Namespace("behbudiy"),
					},
					{
						Group:     Group("gateway.networking.k8s.io"),
						Kind:      Kind("TCPRoute"),
						Namespace: Namespace("qodiriy"),
					},
				},
				To: []ReferenceGrantTo{
					{
						Group: Group(""),
						Kind:  Kind("GrantTwo"),
					},
				},
			},
		},
	}
	tests := []struct {
		msg    string
		from   ReferenceGrantFrom
		result map[Namespace][]ReferenceGrantTo
	}{
		{
			msg: "no matches whatsoever",
			from: ReferenceGrantFrom{
				Group:     Group("invalid.example"),
				Kind:      Kind("invalid"),
				Namespace: Namespace("invalid"),
			},
			result: map[Namespace][]ReferenceGrantTo{},
		},
		{
			msg: "non-matching namespace",
			from: ReferenceGrantFrom{
				Group:     Group("gateway.networking.k8s.io"),
				Kind:      Kind("UDPRoute"),
				Namespace: Namespace("niyazi"),
			},
			result: map[Namespace][]ReferenceGrantTo{},
		},
		{
			msg: "non-matching kind",
			from: ReferenceGrantFrom{
				Group:     Group("gateway.networking.k8s.io"),
				Kind:      Kind("TLSRoute"),
				Namespace: Namespace("behbudiy"),
			},
			result: map[Namespace][]ReferenceGrantTo{},
		},
		{
			msg: "non-matching group",
			from: ReferenceGrantFrom{
				Group:     Group("invalid.example"),
				Kind:      Kind("UDPRoute"),
				Namespace: Namespace("behbudiy"),
			},
			result: map[Namespace][]ReferenceGrantTo{},
		},
		{
			msg: "single match",
			from: ReferenceGrantFrom{
				Group:     Group("gateway.networking.k8s.io"),
				Kind:      Kind("UDPRoute"),
				Namespace: Namespace("behbudiy"),
			},
			result: map[Namespace][]ReferenceGrantTo{
				"cholpon": {
					{
						Group: Group(""),
						Kind:  Kind("GrantTwo"),
					},
				},
			},
		},
		{
			msg: "multiple matches",
			from: ReferenceGrantFrom{
				Group:     Group("gateway.networking.k8s.io"),
				Kind:      Kind("TCPRoute"),
				Namespace: Namespace("qodiriy"),
			},
			result: map[Namespace][]ReferenceGrantTo{
				"cholpon": {
					{
						Group: Group(""),
						Kind:  Kind("GrantTwo"),
					},
				},
				"fitrat": {
					{
						Group: Group(""),
						Kind:  Kind("GrantOne"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			result := GetPermittedForReferenceGrantFrom(logr.Discard(), tt.from, grants)
			assert.Equal(t, tt.result, result)
		})
	}
}
