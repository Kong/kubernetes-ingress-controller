package test

const (
	// HTTPBinImage is the container image name we use for deploying the "httpbin" HTTP testing tool.
	// if you need a simple HTTP server for tests you're writing, use this and check the documentation.
	// See: https://github.com/kong/httpbin
	HTTPBinImage = "kong/httpbin:0.1.0"
	HTTPBinPort  = 80

	// EchoImage works with TCP, UDP, HTTP and responses with basic information about its environment and echo
	// read more about it here: http://github.com/kong/go-echo
	// e.g.
	// Welcome, you are connected to node kind-control-plane.
	// Running on Pod tcp-echo-58ccd6b78d-hn9t8.
	// In namespace foo.
	// With IP address 10.244.0.13.
	EchoImage    = "kong/go-echo:0.3.0"
	EchoTCPPort  = 1025
	EchoUDPPort  = 1026
	EchoHTTPPort = 1027
)
