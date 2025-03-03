package test

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"time"
)

type Protocol string

const (
	ProtocolTCP Protocol = "tcp"
	ProtocolUDP Protocol = "udp"
	ProtocolTLS Protocol = "tls"
)

type TLSOpt struct {
	CertPool    *x509.CertPool
	Hostname    string
	Passthrough bool
}

// WithTLSOption returns a tlsOpt struct with the provided values. Use it when ProtocolTLS is used.
// If passthrough is true, the go-echo service should return a message mentioning that it is listening
// in TLS mode.
func WithTLSOption(hostname string, certPool *x509.CertPool, passthrough bool) TLSOpt {
	return TLSOpt{
		Hostname:    hostname,
		CertPool:    certPool,
		Passthrough: passthrough,
	}
}

// EchoResponds takes a TCP, TLS or UDP address URL and a Pod name and checks if
// a go-echo instance is running on that Pod at that address. For TLS tlsOpt is
// required, otherwise it panics.
// It sends a message and checks if returned one matches. It returns an error with
// an explanation, wraps typical errors as io.EOF or syscall.ECONNRESET.
func EchoResponds(protocol Protocol, url string, podName string, tlsOpt ...TLSOpt) error {
	if protocol == ProtocolTLS && len(tlsOpt) != 1 {
		panic("TLS protocol requires TLS options (fix the code calling this function)")
	}

	dialer := net.Dialer{Timeout: RequestTimeout}
	var (
		tlsCfg TLSOpt
		conn   net.Conn
		err    error
	)
	if protocol == ProtocolTLS {
		tlsCfg = tlsOpt[0]
		conn, err = tls.DialWithDialer(
			&dialer,
			"tcp",
			url,
			&tls.Config{
				MinVersion: tls.VersionTLS12,
				ServerName: tlsCfg.Hostname,
				RootCAs:    tlsCfg.CertPool,
			},
		)
	} else {
		conn, err = dialer.Dial(string(protocol), url)
	}
	if err != nil {
		return fmt.Errorf("cannot dial %q: %w", protocol, err)
	}
	defer conn.Close()

	header := fmt.Appendf(nil, "Running on Pod %s.", podName)
	// If we are testing with passthrough, the go-echo service should return a message
	// mentioning that it is listening in TLS mode.
	if tlsCfg.Passthrough {
		header = append(header, []byte("\nThrough TLS connection.")...)
	}

	message := fmt.Appendf(nil, "testing %sroute", protocol)
	wrote, err := conn.Write(message)
	if err != nil {
		return fmt.Errorf("cannot write message: %w", err)
	}
	if wrote != len(message) {
		return fmt.Errorf("wrote message of size %d, expected %d", wrote, len(message))
	}

	if err := conn.SetDeadline(time.Now().Add(time.Second * 5)); err != nil {
		return fmt.Errorf("cannot set deadline: %w", err)
	}

	headerResponse := make([]byte, len(header)+1)
	read, err := conn.Read(headerResponse)
	if err != nil {
		return fmt.Errorf("cannot read header response: %w", err)
	}

	if read != len(header)+1 { // add 1 for newline
		return fmt.Errorf("read %d bytes but expected %d", read, len(header)+1)
	}

	if !bytes.Contains(headerResponse, header) {
		return fmt.Errorf(`expected header response "%s", received: "%s"`, header, headerResponse)
	}

	messageResponse := make([]byte, wrote+1)
	read, err = conn.Read(messageResponse)
	if err != nil {
		return fmt.Errorf("cannot read message response: %w", err)
	}

	if read != len(message) {
		return fmt.Errorf("read %d bytes but expected %d", read, len(message))
	}

	if !bytes.Contains(messageResponse, message) {
		return fmt.Errorf(`expected message response "%s", received: "%s"`, message, messageResponse)
	}

	return nil
}
