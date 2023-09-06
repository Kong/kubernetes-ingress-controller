package helpers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

// TCPProxy is a simple server that forwards TCP connections to a given destination.
// It can be used to simulate network failures by stopping accepting new connections and interrupting existing ones.
type TCPProxy struct {
	destination string
	address     string
	listener    net.Listener

	// interruptSignalChs is a list of channels that are used to interrupt connections.
	interruptSignalChs []chan struct{}
	// shouldHandleNewConnections is a flag that indicates whether new connections should be accepted.
	shouldHandleNewConnections bool

	mu sync.RWMutex
}

func NewTCPProxy(destination string) (*TCPProxy, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}
	return &TCPProxy{
		destination:                destination,
		address:                    listener.Addr().String(),
		listener:                   listener,
		shouldHandleNewConnections: true,
	}, nil
}

// Run starts connections accepting loop and blocks until the context is canceled.
func (p *TCPProxy) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		_ = p.listener.Close()
	}()

	for {
		c, err := p.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return fmt.Errorf("failed to accept connection: %w", err)
		}

		if !p.shouldHandleConnections() {
			_ = c.Close()
			continue
		}

		go p.handleConnection(c, p.newInterruptSignalCh())
	}
}

// Address returns the address of the proxy.
func (p *TCPProxy) Address() string {
	return p.address
}

// StopHandlingConnections stops handling connections by interrupting all existing connections and immediately closing
// new connections.
func (p *TCPProxy) StopHandlingConnections() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Interrupt all connections.
	for _, ch := range p.interruptSignalChs {
		close(ch)
	}
	p.interruptSignalChs = nil

	// Ensure no new connections are accepted.
	p.shouldHandleNewConnections = false
}

// StartHandlingConnections starts handling new connections.
func (p *TCPProxy) StartHandlingConnections() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.shouldHandleNewConnections = true
}

func (p *TCPProxy) handleConnection(src net.Conn, interruptSignalCh chan struct{}) {
	defer func() {
		_ = src.Close()
	}()

	dst, err := net.Dial("tcp", p.destination)
	if err != nil {
		return
	}
	defer func() {
		_ = dst.Close()
	}()

	copyDoneCh := make(chan struct{}, 2)
	go p.copy(dst, src, copyDoneCh)
	go p.copy(src, dst, copyDoneCh)

	select {
	case <-copyDoneCh:
	case <-interruptSignalCh:
	}
}

func (p *TCPProxy) copy(dst, src net.Conn, doneCh chan struct{}) {
	_, _ = io.Copy(dst, src)
	doneCh <- struct{}{}
}

func (p *TCPProxy) shouldHandleConnections() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.shouldHandleNewConnections
}

func (p *TCPProxy) newInterruptSignalCh() chan struct{} {
	ch := make(chan struct{})
	p.interruptSignalChs = append(p.interruptSignalChs, ch)
	return ch
}
