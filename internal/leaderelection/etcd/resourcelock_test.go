package etcd_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/leaderelection/resourcelock"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/leaderelection/etcd"
)

func TestNewEtcdLock(t *testing.T) {
	cfg := etcd.EtcdLockConfig{
		Client:        nil, // nil for unit test
		Key:           "/test/leader-election/test-id",
		LeaseDuration: 15 * time.Second,
		Identity:      "test-pod-123_abc-def-ghi",
	}

	lock := etcd.NewEtcdLock(cfg)

	require.NotNil(t, lock)
	require.Equal(t, cfg.Key, lock.Key)
	require.Equal(t, cfg.LeaseDuration, lock.LeaseDuration)
	require.Equal(t, cfg.Identity, lock.Identity())
}

func TestEtcdLockIdentity(t *testing.T) {
	tests := []struct {
		name     string
		identity string
	}{
		{
			name:     "simple identity",
			identity: "pod-abc123",
		},
		{
			name:     "identity with UUID format",
			identity: "kong-ingress-controller-7b9f8_550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:     "identity with namespace prefix",
			identity: "kube-system/leader-election-pod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lock := etcd.NewEtcdLock(etcd.EtcdLockConfig{
				Key:      "/test/key",
				Identity: tt.identity,
			})

			require.Equal(t, tt.identity, lock.Identity())
		})
	}
}

func TestEtcdLockDescribe(t *testing.T) {
	tests := []struct {
		name             string
		key              string
		expectedDescribe string
	}{
		{
			name:             "simple key",
			key:              "/test/leader",
			expectedDescribe: "etcd lock: /test/leader",
		},
		{
			name:             "full key path",
			key:              "/kong-ingress-controller/leader-election/default/5b374a9e.konghq.com",
			expectedDescribe: "etcd lock: /kong-ingress-controller/leader-election/default/5b374a9e.konghq.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lock := etcd.NewEtcdLock(etcd.EtcdLockConfig{
				Key:      tt.key,
				Identity: "test-identity",
			})

			require.Equal(t, tt.expectedDescribe, lock.Describe())
		})
	}
}

func TestNewLeaderElectionRecord(t *testing.T) {
	identity := "test-pod-123_abc-def-ghi"
	leaseDurationSeconds := 15

	record := etcd.NewLeaderElectionRecord(identity, leaseDurationSeconds)

	require.Equal(t, identity, record.HolderIdentity)
	require.Equal(t, leaseDurationSeconds, record.LeaseDurationSeconds)
	require.Equal(t, 0, record.LeaderTransitions)
	// AcquireTime and RenewTime should be set to current time (approximately).
	require.False(t, record.AcquireTime.IsZero())
	require.False(t, record.RenewTime.IsZero())
	require.WithinDuration(t, record.AcquireTime.Time, record.RenewTime.Time, time.Second)
}

func TestEtcdLockRecordEvent(t *testing.T) {
	// Test that RecordEvent doesn't panic when EventRecorder is nil.
	lock := etcd.NewEtcdLock(etcd.EtcdLockConfig{
		Key:           "/test/key",
		Identity:      "test-identity",
		EventRecorder: nil,
	})

	// This should not panic.
	lock.RecordEvent("test event")
}

func TestEtcdLockImplementsResourceLockInterface(t *testing.T) {
	// Compile-time check that EtcdLock implements resourcelock.Interface.
	var _ resourcelock.Interface = (*etcd.EtcdLock)(nil)
}

func TestEtcdResourceLockConstant(t *testing.T) {
	require.Equal(t, "etcd", etcd.EtcdResourceLock)
}
