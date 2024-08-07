package adminapi

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"

	tlsutil "github.com/kong/kubernetes-ingress-controller/v3/internal/util/tls"
)

type KonnectConfig struct {
	// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/3922
	// ConfigSynchronizationEnabled is the only toggle we had prior to the addition of the license agent.
	// We likely want to combine these into a single Konnect toggle or piggyback off other Konnect functionality.
	ConfigSynchronizationEnabled bool
	ControlPlaneID               string
	Address                      string
	UploadConfigPeriod           time.Duration
	RefreshNodePeriod            time.Duration
	TLSClient                    TLSClientConfig

	LicenseSynchronizationEnabled bool
	InitialLicensePollingPeriod   time.Duration
	LicensePollingPeriod          time.Duration
	ConsumersSyncDisabled         bool
}

func NewKongClientForKonnectControlPlane(c KonnectConfig) (*KonnectClient, error) {
	clientCertificate, err := tlsutil.ExtractClientCertificates(
		[]byte(c.TLSClient.Cert),
		c.TLSClient.CertFile,
		[]byte(c.TLSClient.Key),
		c.TLSClient.KeyFile,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to extract client certificates: %w", err)
	}
	if clientCertificate == nil {
		return nil, fmt.Errorf("client certificate is missing")
	}

	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{*clientCertificate},
		MinVersion:   tls.VersionTLS12,
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tlsConfig
	client, err := NewKongAPIClient(
		fmt.Sprintf("%s/%s/%s", c.Address, "kic/api/control-planes", c.ControlPlaneID),
		&http.Client{
			Transport: transport,
		},
	)
	if err != nil {
		return nil, err
	}
	return NewKonnectClient(client, c.ControlPlaneID, c.ConsumersSyncDisabled), nil
}

// EnsureKonnectConnection ensures that the client is able to connect to Konnect.
func EnsureKonnectConnection(ctx context.Context, client *kong.Client, logger logr.Logger) error {
	const (
		retries = 5
		delay   = time.Second
	)

	if err := retry.Do(
		func() error {
			// Call an arbitrary endpoint that should return no error.
			if _, _, err := client.Services.List(ctx, &kong.ListOpt{Size: 1}); err != nil {
				if errors.Is(err, context.Canceled) {
					return retry.Unrecoverable(err)
				}
				return err
			}
			return nil
		},
		retry.Attempts(retries),
		retry.Context(ctx),
		retry.Delay(delay),
		retry.DelayType(retry.FixedDelay),
		retry.LastErrorOnly(true),
		retry.OnRetry(func(n uint, err error) {
			logger.Info("Konnect Admin API client unhealthy, retrying", "retry", n, "error", err.Error())
		}),
	); err != nil {
		return fmt.Errorf("konnect client unhealthy, no success after %d retries: %w", retries, err)
	}

	return nil
}
