package helpers

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAllUpdatedReplicasRolledOutAndReady(t *testing.T) {
	tests := []struct {
		name       string
		generation int64
		status     appsv1.DeploymentStatus
		wantErr    bool
	}{
		{
			name:       "rolled out",
			generation: 2,
			status: appsv1.DeploymentStatus{
				ObservedGeneration: 2,
				Replicas:           1,
				UpdatedReplicas:    1,
				AvailableReplicas:  1,
			},
		},
		{
			name:       "status of the previous generation looks rolled out but the update is not observed yet",
			generation: 2,
			status: appsv1.DeploymentStatus{
				ObservedGeneration: 1,
				Replicas:           1,
				UpdatedReplicas:    1,
				AvailableReplicas:  1,
			},
			wantErr: true,
		},
		{
			name:       "old replicas pending termination",
			generation: 2,
			status: appsv1.DeploymentStatus{
				ObservedGeneration: 2,
				Replicas:           2,
				UpdatedReplicas:    1,
				AvailableReplicas:  1,
			},
			wantErr: true,
		},
		{
			name:       "updated replicas not available yet",
			generation: 2,
			status: appsv1.DeploymentStatus{
				ObservedGeneration: 2,
				Replicas:           1,
				UpdatedReplicas:    1,
				AvailableReplicas:  0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Generation: tt.generation},
				Spec:       appsv1.DeploymentSpec{Replicas: lo.ToPtr(int32(1))},
				Status:     tt.status,
			}
			err := allUpdatedReplicasRolledOutAndReady(d)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
