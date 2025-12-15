package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConfig_validateClientSideThrottling(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "throttling disabled - no validation",
			config: Config{
				EnableClientSideThrottling: false,
				APIServerQPS:               0,  // Invalid but should not trigger error
				APIServerBurst:             -1, // Invalid but should not trigger error
			},
			wantErr: false,
		},
		{
			name: "throttling enabled with valid positive values",
			config: Config{
				EnableClientSideThrottling: true,
				APIServerQPS:               100,
				APIServerBurst:             300,
			},
			wantErr: false,
		},
		{
			name: "throttling enabled with QPS = 0",
			config: Config{
				EnableClientSideThrottling: true,
				APIServerQPS:               0,
				APIServerBurst:             300,
			},
			wantErr: true,
			errMsg:  "apiserver-qps must be positive when client-side throttling is enabled, got 0",
		},
		{
			name: "throttling enabled with negative QPS",
			config: Config{
				EnableClientSideThrottling: true,
				APIServerQPS:               -10,
				APIServerBurst:             300,
			},
			wantErr: true,
			errMsg:  "apiserver-qps must be positive when client-side throttling is enabled, got -10",
		},
		{
			name: "throttling enabled with Burst = 0",
			config: Config{
				EnableClientSideThrottling: true,
				APIServerQPS:               100,
				APIServerBurst:             0,
			},
			wantErr: true,
			errMsg:  "apiserver-burst must be positive when client-side throttling is enabled, got 0",
		},
		{
			name: "throttling enabled with negative Burst",
			config: Config{
				EnableClientSideThrottling: true,
				APIServerQPS:               100,
				APIServerBurst:             -50,
			},
			wantErr: true,
			errMsg:  "apiserver-burst must be positive when client-side throttling is enabled, got -50",
		},
		{
			name: "throttling enabled with QPS=1 and Burst=1",
			config: Config{
				EnableClientSideThrottling: true,
				APIServerQPS:               1,
				APIServerBurst:             1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateClientSideThrottling()
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConfig_validateLeaderElection(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid leader election timing",
			config: Config{
				LeaderElectionLeaseDuration: 15 * time.Second,
				LeaderElectionRenewDeadline: 10 * time.Second,
				LeaderElectionRetryPeriod:   2 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "zero lease duration",
			config: Config{
				LeaderElectionLeaseDuration: 0,
				LeaderElectionRenewDeadline: 10 * time.Second,
				LeaderElectionRetryPeriod:   2 * time.Second,
			},
			wantErr: true,
			errMsg:  "leader-election-lease-duration must be positive",
		},
		{
			name: "negative lease duration",
			config: Config{
				LeaderElectionLeaseDuration: -15 * time.Second,
				LeaderElectionRenewDeadline: 10 * time.Second,
				LeaderElectionRetryPeriod:   2 * time.Second,
			},
			wantErr: true,
			errMsg:  "leader-election-lease-duration must be positive",
		},
		{
			name: "zero renew deadline",
			config: Config{
				LeaderElectionLeaseDuration: 15 * time.Second,
				LeaderElectionRenewDeadline: 0,
				LeaderElectionRetryPeriod:   2 * time.Second,
			},
			wantErr: true,
			errMsg:  "leader-election-renew-deadline must be positive",
		},
		{
			name: "zero retry period",
			config: Config{
				LeaderElectionLeaseDuration: 15 * time.Second,
				LeaderElectionRenewDeadline: 10 * time.Second,
				LeaderElectionRetryPeriod:   0,
			},
			wantErr: true,
			errMsg:  "leader-election-retry-period must be positive",
		},
		{
			name: "renew deadline >= lease duration",
			config: Config{
				LeaderElectionLeaseDuration: 15 * time.Second,
				LeaderElectionRenewDeadline: 15 * time.Second,
				LeaderElectionRetryPeriod:   2 * time.Second,
			},
			wantErr: true,
			errMsg:  "leader-election-renew-deadline (15s) must be less than leader-election-lease-duration (15s)",
		},
		{
			name: "renew deadline > lease duration",
			config: Config{
				LeaderElectionLeaseDuration: 10 * time.Second,
				LeaderElectionRenewDeadline: 15 * time.Second,
				LeaderElectionRetryPeriod:   2 * time.Second,
			},
			wantErr: true,
			errMsg:  "leader-election-renew-deadline (15s) must be less than leader-election-lease-duration (10s)",
		},
		{
			name: "retry period >= renew deadline",
			config: Config{
				LeaderElectionLeaseDuration: 15 * time.Second,
				LeaderElectionRenewDeadline: 10 * time.Second,
				LeaderElectionRetryPeriod:   10 * time.Second,
			},
			wantErr: true,
			errMsg:  "leader-election-retry-period (10s) must be less than leader-election-renew-deadline (10s)",
		},
		{
			name: "retry period > renew deadline",
			config: Config{
				LeaderElectionLeaseDuration: 15 * time.Second,
				LeaderElectionRenewDeadline: 5 * time.Second,
				LeaderElectionRetryPeriod:   10 * time.Second,
			},
			wantErr: true,
			errMsg:  "leader-election-retry-period (10s) must be less than leader-election-renew-deadline (5s)",
		},
		{
			name: "minimal valid timing (1ms each, properly ordered)",
			config: Config{
				LeaderElectionLeaseDuration: 3 * time.Millisecond,
				LeaderElectionRenewDeadline: 2 * time.Millisecond,
				LeaderElectionRetryPeriod:   1 * time.Millisecond,
			},
			wantErr: false,
		},
		{
			name: "all timing parameters equal",
			config: Config{
				LeaderElectionLeaseDuration: 10 * time.Second,
				LeaderElectionRenewDeadline: 10 * time.Second,
				LeaderElectionRetryPeriod:   10 * time.Second,
			},
			wantErr: true,
			errMsg:  "leader-election-renew-deadline (10s) must be less than leader-election-lease-duration (10s)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateLeaderElection()
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConfig_Validate_ClientSideThrottlingAndLeaderElection(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config with throttling disabled",
			config: Config{
				EnableClientSideThrottling:  false,
				APIServerQPS:                100,
				APIServerBurst:              300,
				LeaderElectionLeaseDuration: 15 * time.Second,
				LeaderElectionRenewDeadline: 10 * time.Second,
				LeaderElectionRetryPeriod:   2 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid config with throttling enabled",
			config: Config{
				EnableClientSideThrottling:  true,
				APIServerQPS:                100,
				APIServerBurst:              300,
				LeaderElectionLeaseDuration: 15 * time.Second,
				LeaderElectionRenewDeadline: 10 * time.Second,
				LeaderElectionRetryPeriod:   2 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "invalid throttling config",
			config: Config{
				EnableClientSideThrottling:  true,
				APIServerQPS:                -1,
				APIServerBurst:              300,
				LeaderElectionLeaseDuration: 15 * time.Second,
				LeaderElectionRenewDeadline: 10 * time.Second,
				LeaderElectionRetryPeriod:   2 * time.Second,
			},
			wantErr: true,
			errMsg:  "invalid client-side throttling configuration",
		},
		{
			name: "invalid leader election config",
			config: Config{
				EnableClientSideThrottling:  true,
				APIServerQPS:                100,
				APIServerBurst:              300,
				LeaderElectionLeaseDuration: 15 * time.Second,
				LeaderElectionRenewDeadline: 20 * time.Second,
				LeaderElectionRetryPeriod:   2 * time.Second,
			},
			wantErr: true,
			errMsg:  "invalid leader election configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
