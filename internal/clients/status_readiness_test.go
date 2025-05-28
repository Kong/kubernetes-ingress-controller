package clients

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
)

// Mock implementations for testing

type mockStatusClient struct {
	mock.Mock
}

func (m *mockStatusClient) IsReady(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockStatusClient) PodReference() (k8stypes.NamespacedName, bool) {
	args := m.Called()
	return args.Get(0).(k8stypes.NamespacedName), args.Bool(1)
}

func (m *mockStatusClient) BaseRootURL() string {
	args := m.Called()
	return args.String(0)
}

type mockStatusClientFactory struct {
	mock.Mock
}

func (m *mockStatusClientFactory) CreateStatusClient(ctx context.Context, discoveredStatusAPI adminapi.DiscoveredAdminAPI) (*adminapi.StatusClient, error) {
	args := m.Called(ctx, discoveredStatusAPI)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*adminapi.StatusClient), args.Error(1)
}

func TestDefaultStatusReadinessChecker_CheckStatusReadiness(t *testing.T) {
	tests := []struct {
		name                     string
		alreadyCreatedClients    []AlreadyCreatedStatusClient
		pendingClients           []adminapi.DiscoveredAdminAPI
		setupMocks               func(*mockStatusClientFactory, []*mockStatusClient)
		expectedTurnedReady      int
		expectedTurnedPending    int
	}{
		{
			name:                  "no clients",
			alreadyCreatedClients: []AlreadyCreatedStatusClient{},
			pendingClients:        []adminapi.DiscoveredAdminAPI{},
			setupMocks:            func(*mockStatusClientFactory, []*mockStatusClient) {},
			expectedTurnedReady:   0,
			expectedTurnedPending: 0,
		},
		{
			name:                  "pending client becomes ready",
			alreadyCreatedClients: []AlreadyCreatedStatusClient{},
			pendingClients: []adminapi.DiscoveredAdminAPI{
				{
					Address: "https://10.0.0.1:8100",
					PodRef: k8stypes.NamespacedName{
						Name:      "pod-1",
						Namespace: "default",
					},
				},
			},
			setupMocks: func(factory *mockStatusClientFactory, clients []*mockStatusClient) {
				factory.On("CreateStatusClient", mock.Anything, mock.Anything).Return(&adminapi.StatusClient{}, nil)
			},
			expectedTurnedReady:   1,
			expectedTurnedPending: 0,
		},
		{
			name:                  "pending client fails to become ready",
			alreadyCreatedClients: []AlreadyCreatedStatusClient{},
			pendingClients: []adminapi.DiscoveredAdminAPI{
				{
					Address: "https://10.0.0.1:8100",
					PodRef: k8stypes.NamespacedName{
						Name:      "pod-1",
						Namespace: "default",
					},
				},
			},
			setupMocks: func(factory *mockStatusClientFactory, clients []*mockStatusClient) {
				factory.On("CreateStatusClient", mock.Anything, mock.Anything).Return(nil, errors.New("connection failed"))
			},
			expectedTurnedReady:   0,
			expectedTurnedPending: 0,
		},
		{
			name: "ready client becomes pending",
			alreadyCreatedClients: []AlreadyCreatedStatusClient{
				func() AlreadyCreatedStatusClient {
					client := &mockStatusClient{}
					client.On("IsReady", mock.Anything).Return(errors.New("not ready"))
					client.On("PodReference").Return(k8stypes.NamespacedName{Name: "pod-1", Namespace: "default"}, true)
					client.On("BaseRootURL").Return("https://10.0.0.1:8100")
					return client
				}(),
			},
			pendingClients: []adminapi.DiscoveredAdminAPI{},
			setupMocks:     func(*mockStatusClientFactory, []*mockStatusClient) {},
			expectedTurnedReady:   0,
			expectedTurnedPending: 1,
		},
		{
			name: "ready client stays ready",
			alreadyCreatedClients: []AlreadyCreatedStatusClient{
				func() AlreadyCreatedStatusClient {
					client := &mockStatusClient{}
					client.On("IsReady", mock.Anything).Return(nil)
					client.On("BaseRootURL").Return("https://10.0.0.1:8100")
					return client
				}(),
			},
			pendingClients: []adminapi.DiscoveredAdminAPI{},
			setupMocks:     func(*mockStatusClientFactory, []*mockStatusClient) {},
			expectedTurnedReady:   0,
			expectedTurnedPending: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := &mockStatusClientFactory{}
			mockClients := make([]*mockStatusClient, len(tt.alreadyCreatedClients))
			for i := range mockClients {
				mockClients[i] = &mockStatusClient{}
			}

			tt.setupMocks(factory, mockClients)

			checker := NewDefaultStatusReadinessChecker(
				factory,
				time.Second*5,
				logr.Discard(),
			)

			ctx := context.Background()
			result := checker.CheckStatusReadiness(ctx, tt.alreadyCreatedClients, tt.pendingClients)

			assert.Len(t, result.ClientsTurnedReady, tt.expectedTurnedReady)
			assert.Len(t, result.ClientsTurnedPending, tt.expectedTurnedPending)

			// Verify that HasChanges works correctly
			expectedHasChanges := tt.expectedTurnedReady > 0 || tt.expectedTurnedPending > 0
			assert.Equal(t, expectedHasChanges, result.HasChanges())

			// Verify mock expectations
			factory.AssertExpectations(t)
			for _, client := range tt.alreadyCreatedClients {
				if mockClient, ok := client.(*mockStatusClient); ok {
					mockClient.AssertExpectations(t)
				}
			}
		})
	}
}

func TestStatusReadinessCheckResult_HasChanges(t *testing.T) {
	tests := []struct {
		name                 string
		clientsTurnedReady   int
		clientsTurnedPending int
		expectedHasChanges   bool
	}{
		{
			name:                 "no changes",
			clientsTurnedReady:   0,
			clientsTurnedPending: 0,
			expectedHasChanges:   false,
		},
		{
			name:                 "clients turned ready",
			clientsTurnedReady:   1,
			clientsTurnedPending: 0,
			expectedHasChanges:   true,
		},
		{
			name:                 "clients turned pending",
			clientsTurnedReady:   0,
			clientsTurnedPending: 1,
			expectedHasChanges:   true,
		},
		{
			name:                 "both changes",
			clientsTurnedReady:   1,
			clientsTurnedPending: 1,
			expectedHasChanges:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StatusReadinessCheckResult{
				ClientsTurnedReady:   make([]*adminapi.StatusClient, tt.clientsTurnedReady),
				ClientsTurnedPending: make([]adminapi.DiscoveredAdminAPI, tt.clientsTurnedPending),
			}

			assert.Equal(t, tt.expectedHasChanges, result.HasChanges())
		})
	}
}

func TestDefaultStatusReadinessChecker_checkPendingStatusClient(t *testing.T) {
	factory := &mockStatusClientFactory{}
	checker := NewDefaultStatusReadinessChecker(
		factory,
		time.Second*5,
		logr.Discard(),
	)

	pendingClient := adminapi.DiscoveredAdminAPI{
		Address: "https://10.0.0.1:8100",
		PodRef: k8stypes.NamespacedName{
			Name:      "pod-1",
			Namespace: "default",
		},
	}

	t.Run("successful creation", func(t *testing.T) {
		expectedClient := &adminapi.StatusClient{}
		factory.On("CreateStatusClient", mock.Anything, pendingClient).Return(expectedClient, nil).Once()

		ctx := context.Background()
		result := checker.checkPendingStatusClient(ctx, pendingClient)

		assert.Equal(t, expectedClient, result)
		factory.AssertExpectations(t)
	})

	t.Run("failed creation", func(t *testing.T) {
		factory.On("CreateStatusClient", mock.Anything, pendingClient).Return(nil, errors.New("creation failed")).Once()

		ctx := context.Background()
		result := checker.checkPendingStatusClient(ctx, pendingClient)

		assert.Nil(t, result)
		factory.AssertExpectations(t)
	})
}

func TestDefaultStatusReadinessChecker_checkAlreadyCreatedStatusClient(t *testing.T) {
	checker := NewDefaultStatusReadinessChecker(
		&mockStatusClientFactory{},
		time.Second*5,
		logr.Discard(),
	)

	t.Run("client is ready", func(t *testing.T) {
		client := &mockStatusClient{}
		client.On("IsReady", mock.Anything).Return(nil)
		client.On("BaseRootURL").Return("https://10.0.0.1:8100")

		ctx := context.Background()
		result := checker.checkAlreadyCreatedStatusClient(ctx, client)

		assert.True(t, result)
		client.AssertExpectations(t)
	})

	t.Run("client is not ready", func(t *testing.T) {
		client := &mockStatusClient{}
		client.On("IsReady", mock.Anything).Return(errors.New("not ready"))
		client.On("BaseRootURL").Return("https://10.0.0.1:8100")

		ctx := context.Background()
		result := checker.checkAlreadyCreatedStatusClient(ctx, client)

		assert.False(t, result)
		client.AssertExpectations(t)
	})
}