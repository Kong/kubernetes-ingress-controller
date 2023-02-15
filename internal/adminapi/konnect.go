package adminapi

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"

	tlsutil "github.com/kong/kubernetes-ingress-controller/v2/internal/util/tls"
)

type KonnectConfig struct {
	ConfigSynchronizationEnabled bool
	RuntimeGroupID               string
	Address                      string
	RefreshNodePeriod            time.Duration
	TLSClient                    TLSClientConfig
}

func NewKongClientForKonnectRuntimeGroup(ctx context.Context, c KonnectConfig) (*Client, error) {
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
		return nil, fmt.Errorf("client ceritficate is missing")
	}

	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{*clientCertificate},
		MinVersion:   tls.VersionTLS12,
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tlsConfig
	client, err := kong.NewClient(
		lo.ToPtr(fmt.Sprintf("%s/%s/%s", c.Address, "kic/api/runtime_groups", c.RuntimeGroupID)),
		&http.Client{
			Transport: transport,
		},
	)
	if err != nil {
		return nil, err
	}
	// Konnect supports tags, we don't need to verify that.
	client.Tags = tagsStub{}

	if err := ensureKonnectConnection(ctx, client); err != nil {
		return nil, err
	}
	return NewKonnectClient(client, c.RuntimeGroupID), nil
}

func ensureKonnectConnection(ctx context.Context, client *kong.Client) error {
	const (
		retries = 60
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
	); err != nil {
		return fmt.Errorf("konnect client unhealthy, no success after %d retries: %w", retries, err)
	}

	return nil
}

// tagsStub replaces a default Tags service in the go-kong's Client for Konnect clients.
// It will always tell tags are supported, which is true for Konnect Runtime Group Admin API.
type tagsStub struct{}

func (t tagsStub) Exists(context.Context) (bool, error) {
	return true, nil
}
