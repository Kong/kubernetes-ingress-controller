package test

import "time"

const (
	// HTTPBinImage is the container image name we use for deploying the "httpbin" HTTP testing tool.
	// If you need a simple HTTP server for tests you're writing, use this and check the documentation.
	// See: https://github.com/kong/httpbin
	HTTPBinImage = "kong/httpbin:0.1.0"
	HTTPBinPort  = 80

	// EchoImage works with TCP, UDP, HTTP, TLS and responds with basic information about its environment and echo
	// Sample response:
	// Welcome, you are connected to node kind-control-plane.
	// Running on Pod tcp-echo-58ccd6b78d-hn9t8.
	// In namespace foo.
	// With IP address 10.244.0.13.
	// See: http://github.com/kong/go-echo
	EchoImage    = "kong/go-echo:0.3.0"
	EchoTCPPort  = 1025
	EchoUDPPort  = 1026
	EchoHTTPPort = 1027

	// GRPCBinImage is the container image name we use for deploying the "grpcbin" GRPC testing tool.
	// See: https://github.com/Kong/grpcbin
	GRPCBinImage = "kong/grpcbin:latest"
	GRPCBinPort  = 9001

	// EnvironmentCleanupTimeout is the amount of time that will be given by the test suite to the
	// testing environment to perform its cleanup when the test suite is shutting down.
	EnvironmentCleanupTimeout = 3 * time.Minute

	// RequestTimeout is the amount of time that will be given to any request to complete.
	RequestTimeout = 10 * time.Second
)
