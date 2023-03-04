//go:build integration_tests
// +build integration_tests

package integration

import (
	"context"
	"crypto/tls"
	"fmt"

	pb "github.com/moul/pb/grpcbin/go-grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func grpcEchoResponds(ctx context.Context, url, hostname, input string) error {
	conn, err := grpc.DialContext(ctx, url,
		grpc.WithTransportCredentials(credentials.NewTLS(
			&tls.Config{
				ServerName:         hostname,
				InsecureSkipVerify: true, //nolint:gosec
			},
		)),
	)
	if err != nil {
		return fmt.Errorf("failed to dial GRPC server: %w", err)
	}
	defer conn.Close()

	client := pb.NewGRPCBinClient(conn)
	resp, err := client.DummyUnary(ctx, &pb.DummyMessage{
		FString: input,
	})
	if err != nil {
		return fmt.Errorf("failed to send GRPC request: %w", err)
	}

	if resp.FString != input {
		return fmt.Errorf("expected %q, got %q", input, resp.FString)
	}

	return nil
}
