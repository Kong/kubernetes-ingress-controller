package adminapi

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/avast/retry-go/v4"
	deckutils "github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
)

type KonnectConfig struct {
	ConfigSynchronizationEnabled bool
	RuntimeGroup                 string
	Address                      string
	TLSClient                    TLSClientConfig
}

func NewKongClientForKonnectRuntimeGroup(ctx context.Context, c KonnectConfig) (*Client, error) {
	tlsClientCert, err := valueFromVariableOrFile(c.TLSClient.Cert, c.TLSClient.CertFile)
	if err != nil {
		return nil, fmt.Errorf("could not extract TLS client cert")
	}
	tlsClientKey, err := valueFromVariableOrFile(c.TLSClient.Key, c.TLSClient.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not extract TLS client key")
	}

	client, err := deckutils.GetKongClient(deckutils.KongClientConfig{
		Address:       fmt.Sprintf("%s/%s/%s", c.Address, "kic/api/runtime_groups", c.RuntimeGroup),
		TLSClientCert: string(tlsClientCert),
		TLSClientKey:  string(tlsClientKey),
	})
	if err != nil {
		return nil, err
	}
	// Konnect supports tags, we don't need to verify that.
	client.Tags = tagsStub{}

	if err := ensureKonnectConnection(ctx, client); err != nil {
		return nil, err
	}
	return NewKonnectClient(client, c.RuntimeGroup), nil
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

// valueFromVariableOrFile uses v value if it's not empty, and falls back to reading a file content when value is missing.
func valueFromVariableOrFile(v string, file string) ([]byte, error) {
	if v != "" {
		return []byte(v), nil
	}
	if file != "" {
		b, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		return b, nil
	}

	return nil, nil
}
