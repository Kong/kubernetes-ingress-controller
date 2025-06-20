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

type TelemetryServer struct {
	reportChan chan []byte
	listener   net.Listener
	started    bool
	cancel     context.CancelFunc
}

func NewTelemetryServer(t *testing.T) *TelemetryServer {
	t.Log("configuring TLS listener - server for telemetry data")
	cert := certificate.MustGenerateCert()
	telemetryServerListener, err := tls.Listen("tcp", "localhost:0", &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	})
	require.NoError(t, err)

	reportChan := make(chan []byte)

	return &TelemetryServer{
		reportChan: reportChan,
		listener:   telemetryServerListener,
	}
}

func (ts *TelemetryServer) Start(ctx context.Context, t *testing.T) {
	ctx, cancelFunc := context.WithCancel(ctx)

	handleConnection := func(ctx context.Context, t *testing.T, conn net.Conn, wg *sync.WaitGroup) {
		defer func() {
			if err := conn.Close(); err != nil {
				t.Logf("error closing connection: %v", err)
			}
			wg.Done()
		}()

		for {
			report := make([]byte, 2048) // Report is much shorter.
			n, err := conn.Read(report)
			if errors.Is(err, io.EOF) {
				break
			}
			if !assert.NoError(t, err) {
				return
			}
			t.Logf("received %d bytes of telemetry report", n)
			select {
			case ts.reportChan <- report[:n]:
			case <-ctx.Done():
				return
			}
		}
	}

	t.Logf("Starting telemetry server")
	go func() {
		// Any function return indicates that either the
		// report was sent or there was nothing to send.
		var wg sync.WaitGroup
		for {
			select {
			case <-ctx.Done():
				t.Logf("Context cancelled, stopping  telemetry server")
				wg.Wait()
				close(ts.reportChan)
				t.Logf("Telemetry server  stopped")
				return
			default:
				conn, err := ts.listener.Accept()
				if err != nil && errors.Is(err, net.ErrClosed) {
					break
				}
				if !assert.NoError(t, err) {
					break
				}
				wg.Add(1)
				go handleConnection(ctx, t, conn, &wg)
			}
		}
	}()

	ts.started = true
	ts.cancel = cancelFunc
}

func (ts *TelemetryServer) Endpoint() string {
	return ts.listener.Addr().String()
}

func (ts *TelemetryServer) ReportChan() chan []byte {
	return ts.reportChan
}

func (ts *TelemetryServer) Stop(t *testing.T) {
	t.Log("Stopping telemetry server")
	ts.cancel()
	assert.NoError(t, ts.listener.Close())
	ts.started = false
}
