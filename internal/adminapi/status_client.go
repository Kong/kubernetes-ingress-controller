package adminapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type StatusClient struct {
	httpClient *http.Client
}

func NewStatusClient() *StatusClient {
	return &StatusClient{
		httpClient: &http.Client{
			Transport: http.DefaultTransport.(*http.Transport).Clone(),
		},
	}
}

// AdminAPIReady checks if the Gateway's Admin API is ready to accept requests.
func (s *StatusClient) AdminAPIReady(ctx context.Context, address string) error {
	const adminAPIStatusEndpoint = "/status"
	u, err := url.JoinPath(address, adminAPIStatusEndpoint)
	if err != nil {
		return fmt.Errorf("failed to join URL path: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
