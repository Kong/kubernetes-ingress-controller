package envtest

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
)

// TelemetryServer represents a server that listens for telemetry data over a TLS connection.
type TelemetryServer struct {
	// Channel to receive telemetry reports.
	reportChan chan []byte
	// TLS listener for incoming telemetry connections.
	listener net.Listener
	// Indicates whether the server is running.
	started bool
	// Function to cancel the server's context.
	cancel context.CancelFunc
	// WaitGroup to track handlers.
	wg sync.WaitGroup
}

// NewTelemetryServer creates and configures a new TelemetryServer instance.
// It generates a TLS listener using a self-signed certificate.
func NewTelemetryServer(t *testing.T) *TelemetryServer {
	t.Helper()
	t.Log("configuring TLS listener - server for telemetry data")
	telemetryServerListener, err := tls.Listen("tcp", "localhost:0", &tls.Config{
		Certificates: []tls.Certificate{certificate.MustGenerateCert()},
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	})
	require.NoError(t, err)

	return &TelemetryServer{
		reportChan: make(chan []byte),
		listener:   telemetryServerListener,
	}
}

// Start begins the telemetry server, accepting incoming connections and processing telemetry data.
// It does not block.
func (ts *TelemetryServer) Start(ctx context.Context, t *testing.T) {
	t.Helper()
	if ts.started {
		t.Log("telemetry server already started, failing")
		t.FailNow()
	}
	ctx, cancelFunc := context.WithCancel(ctx)

	// handleConnection processes incoming telemetry data from a single connection.
	handleConnection := func(ctx context.Context, t *testing.T, conn net.Conn, wg *sync.WaitGroup) {
		t.Helper()
		defer func() {
			if err := conn.Close(); err != nil {
				t.Logf("error closing connection: %v", err)
			}
			wg.Done()
		}()

		for {
			report := make([]byte, 2048) // Buffer for telemetry data.
			n, err := conn.Read(report)
			if errors.Is(err, io.EOF) {
				break
			}
			if !assert.NoError(t, err) {
				return
			}
			t.Logf("received %d bytes of telemetry report", n)
			select {
			case ts.reportChan <- report[:n]: // Send the report to the channel.
			case <-ctx.Done(): // Exit if the context is canceled.
				return
			}
		}
	}

	t.Logf("Starting telemetry server")
	go func() {
		// Main loop to accept and handle incoming connections.
		for {
			select {
			case <-ctx.Done():
				t.Logf("Context cancelled, stopping telemetry server")
				ts.wg.Wait()
				close(ts.reportChan)
				t.Logf("Telemetry server stopped")
				return
			default:
				conn, err := ts.listener.Accept()
				if err != nil && errors.Is(err, net.ErrClosed) {
					break
				}
				if !assert.NoError(t, err) {
					break
				}
				ts.wg.Add(1)
				go handleConnection(ctx, t, conn, &ts.wg)
			}
		}
	}()

	ts.started = true
	ts.cancel = cancelFunc
}

// Endpoint returns the address of the telemetry server.
func (ts *TelemetryServer) Endpoint() string {
	return ts.listener.Addr().String()
}

// ReportChan provides access to the telemetry report channel.
func (ts *TelemetryServer) ReportChan() <-chan []byte {
	return ts.reportChan
}

// Stop shuts down the telemetry server, closing the listener and canceling the context.
func (ts *TelemetryServer) Stop(t *testing.T) {
	t.Helper()
	if !ts.started {
		t.Log("telemetry server already stopped, doing nothing")
		return
	}
	t.Log("Stopping telemetry server")
	ts.cancel()
	assert.NoError(t, ts.listener.Close())
	ts.wg.Wait()
	t.Log("Telemetry server stopped")
	ts.started = false
}
