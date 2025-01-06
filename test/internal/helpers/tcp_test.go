package helpers_test

import (
	"bytes"
	"context"
	"io"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

func TestTCPProxy(t *testing.T) {
	ls, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	var (
		destAcceptedConn atomic.Bool
		destDroppedConn  atomic.Bool
		destReceivedData bytes.Buffer
	)
	go func() {
		for {
			c, err := ls.Accept()
			if err != nil {
				return
			}
			destAcceptedConn.Store(true)
			go func() {
				_, _ = io.Copy(&destReceivedData, c)
				destDroppedConn.Store(true)
			}()
		}
	}()

	proxy, err := helpers.NewTCPProxy(ls.Addr().String())
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := proxy.Run(ctx)
		assert.NoError(t, err)
	}()

	t.Log("Ensuring proxy is accepting connections by default")
	conn, err := net.Dial("tcp", proxy.Address())
	require.NoError(t, err)
	require.Eventually(t, destAcceptedConn.Load, time.Second, time.Millisecond, "destination didn't accept connection")

	t.Log("Ensuring proxy forwards data from the source to the destination")
	_, err = conn.Write([]byte("hello"))
	require.NoError(t, err)
	require.Eventually(t, func() bool { return destReceivedData.String() == "hello" }, time.Second, time.Millisecond)

	t.Log("Ensuring proxy dropped connection after StopHandlingConnections")
	proxy.StopHandlingConnections()
	require.Eventually(t, destDroppedConn.Load, time.Second, time.Millisecond)

	t.Log("Ensuring proxy dropped existing connection after StopHandlingConnections")
	require.Eventually(t, func() bool {
		_, err = conn.Write([]byte("hello"))
		return err != nil
	}, time.Second, time.Millisecond)

	t.Log("Ensuring proxy handles connections after StartHandlingConnections")
	proxy.StartHandlingConnections()
	require.Eventually(t, func() bool {
		conn, err := net.Dial("tcp", proxy.Address())
		if err != nil {
			return false
		}
		_, err = conn.Write([]byte("hello"))
		return err == nil
	}, time.Second, time.Millisecond)
}
