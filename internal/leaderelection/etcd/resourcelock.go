// Package etcd provides an etcd-based implementation of the resourcelock.Interface
// for use with Kubernetes leader election. This allows the Kong Ingress Controller
// to use etcd directly for leader election instead of the Kubernetes Lease API.
package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

const (
	// EtcdResourceLock is the resource lock type identifier for etcd.
	EtcdResourceLock = "etcd"
)

// EtcdLock implements resourcelock.Interface using etcd as the backend.
// This provides the same interface as the Kubernetes Lease-based lock
// but stores the leader election record directly in etcd.
type EtcdLock struct {
	// Client is the etcd client used for operations.
	Client *clientv3.Client

	// LockConfig contains the identity and event recorder.
	LockConfig resourcelock.ResourceLockConfig

	// Key is the etcd key where the leader election record is stored.
	Key string

	// LeaseDuration is the duration that non-leader candidates will
	// wait to force acquire leadership.
	LeaseDuration time.Duration

	// Internal state for optimistic concurrency.
	observedRecord    *resourcelock.LeaderElectionRecord
	observedRawRecord []byte
	observedRevision  int64
}

// EtcdLockConfig contains configuration for creating an EtcdLock.
type EtcdLockConfig struct {
	// Client is the etcd client.
	Client *clientv3.Client

	// Key is the etcd key for the lock (e.g., "/kong-ingress-controller/leader-election/5b374a9e.konghq.com").
	Key string

	// LeaseDuration is how long the lock is held before expiring.
	LeaseDuration time.Duration

	// Identity is the unique identifier for this lock holder.
	Identity string

	// EventRecorder is used to record events (can be nil).
	EventRecorder resourcelock.EventRecorder
}

// NewEtcdLock creates a new EtcdLock.
func NewEtcdLock(cfg EtcdLockConfig) *EtcdLock {
	return &EtcdLock{
		Client: cfg.Client,
		Key:    cfg.Key,
		LockConfig: resourcelock.ResourceLockConfig{
			Identity:      cfg.Identity,
			EventRecorder: cfg.EventRecorder,
		},
		LeaseDuration: cfg.LeaseDuration,
	}
}

// Get returns the leader election record from etcd.
// It returns the record, the raw bytes, and any error.
func (el *EtcdLock) Get(ctx context.Context) (*resourcelock.LeaderElectionRecord, []byte, error) {
	resp, err := el.Client.Get(ctx, el.Key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get leader election record from etcd: %w", err)
	}

	// If the key doesn't exist, return a NotFound error.
	// This is required by the leader election library - it checks errors.IsNotFound(err)
	// to determine whether to call Create().
	if len(resp.Kvs) == 0 {
		return nil, nil, apierrors.NewNotFound(
			schema.GroupResource{Group: "", Resource: "etcdlock"},
			el.Key,
		)
	}

	kv := resp.Kvs[0]
	record := &resourcelock.LeaderElectionRecord{}
	if err := json.Unmarshal(kv.Value, record); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal leader election record: %w", err)
	}

	// Store the observed state for optimistic concurrency.
	el.observedRecord = record
	el.observedRawRecord = kv.Value
	el.observedRevision = kv.ModRevision

	return record, kv.Value, nil
}

// Create attempts to create a new leader election record in etcd.
// It fails if a record already exists.
func (el *EtcdLock) Create(ctx context.Context, ler resourcelock.LeaderElectionRecord) error {
	recordBytes, err := json.Marshal(ler)
	if err != nil {
		return fmt.Errorf("failed to marshal leader election record: %w", err)
	}

	// Use a transaction to ensure the key doesn't already exist.
	txnResp, err := el.Client.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(el.Key), "=", 0)).
		Then(clientv3.OpPut(el.Key, string(recordBytes))).
		Commit()
	if err != nil {
		return fmt.Errorf("failed to create leader election record in etcd: %w", err)
	}

	if !txnResp.Succeeded {
		return fmt.Errorf("leader election record already exists in etcd")
	}

	// Update observed state.
	el.observedRecord = &ler
	el.observedRawRecord = recordBytes
	if len(txnResp.Responses) > 0 && txnResp.Responses[0].GetResponsePut() != nil {
		el.observedRevision = txnResp.Header.Revision
	}

	return nil
}

// Update attempts to update the leader election record in etcd.
// It uses optimistic concurrency based on the last observed revision.
func (el *EtcdLock) Update(ctx context.Context, ler resourcelock.LeaderElectionRecord) error {
	recordBytes, err := json.Marshal(ler)
	if err != nil {
		return fmt.Errorf("failed to marshal leader election record: %w", err)
	}

	// Use a transaction with compare-and-swap based on ModRevision.
	txnResp, err := el.Client.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(el.Key), "=", el.observedRevision)).
		Then(clientv3.OpPut(el.Key, string(recordBytes))).
		Commit()
	if err != nil {
		return fmt.Errorf("failed to update leader election record in etcd: %w", err)
	}

	if !txnResp.Succeeded {
		return fmt.Errorf("leader election record was modified by another process (optimistic lock failure)")
	}

	// Update observed state.
	el.observedRecord = &ler
	el.observedRawRecord = recordBytes
	el.observedRevision = txnResp.Header.Revision

	return nil
}

// RecordEvent records an event for the leader election.
// This is used for debugging and audit purposes.
func (el *EtcdLock) RecordEvent(s string) {
	if el.LockConfig.EventRecorder == nil {
		return
	}
	// Note: We don't have a Kubernetes object to attach events to when using etcd,
	// so we skip event recording. In a production environment, you might want to
	// log these events or send them to a monitoring system.
}

// Identity returns the identity of this lock holder.
func (el *EtcdLock) Identity() string {
	return el.LockConfig.Identity
}

// Describe returns a description of this resource lock.
func (el *EtcdLock) Describe() string {
	return fmt.Sprintf("etcd lock: %s", el.Key)
}

// Helper functions for creating the etcd client and lock.

// NewEtcdClient creates a new etcd client from the configuration.
func NewEtcdClient(cfg *Config) (*clientv3.Client, error) {
	clientCfg, err := cfg.ToClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client config: %w", err)
	}

	client, err := clientv3.New(clientCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return client, nil
}

// NewEtcdLockFromConfig creates a new EtcdLock from the etcd configuration.
func NewEtcdLockFromConfig(cfg *Config, electionID, identity string, leaseDuration time.Duration) (*EtcdLock, error) {
	client, err := NewEtcdClient(cfg)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%s/%s", cfg.ElectionPrefix, electionID)

	return NewEtcdLock(EtcdLockConfig{
		Client:        client,
		Key:           key,
		LeaseDuration: leaseDuration,
		Identity:      identity,
	}), nil
}

// LeaderElectionRecordToMetav1 converts times in the record to metav1.Time format.
// This is a helper for creating records with proper timestamps.
func NewLeaderElectionRecord(identity string, leaseDurationSeconds int) resourcelock.LeaderElectionRecord {
	now := metav1.Now()
	return resourcelock.LeaderElectionRecord{
		HolderIdentity:       identity,
		LeaseDurationSeconds: leaseDurationSeconds,
		AcquireTime:          now,
		RenewTime:            now,
		LeaderTransitions:    0,
	}
}
