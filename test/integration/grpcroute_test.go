//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func grpcRequest(input string) string {
	return fmt.Sprintf(`{"greeting": "%s"}`, input)
}

func grpcResponse(input string) string {
	return fmt.Sprintf("{\n  \"reply\": \"hello %s\"\n}\n", input)
}

func grpcEchoResponds(ctx context.Context, url, hostname, input string) (bool, error) {
	args := []string{
		"-d",
		grpcRequest(input),
		"-insecure",
		"-servername",
		hostname,
		url,
		"hello.HelloService.SayHello",
	}
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)

	cmd := exec.CommandContext(ctx, "grpcurl", args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("failed to echo GRPC server STDOUT=(%s) STDERR=(%s): %w", strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err)
	}

	return stdout.String() == grpcResponse(input), nil
}
