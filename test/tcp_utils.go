package test

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

type Protocol string

const (
	ProtocolTCP Protocol = "tcp"
	ProtocolUDP Protocol = "udp"
)

// EchoResponds takes a TCP or UDP address URL and a Pod name and checks if
// a go-echo instance is running on that Pod at that address. It sends
// a message and checks if returned one matches. It returns an error with
// an explanation if it is not (typical network related errors like
// io.EOF or syscall.ECONNRESET are returned directly).
func EchoResponds(protocol Protocol, url string, podName string) error {
	dialer := net.Dialer{Timeout: RequestTimeout}
	conn, err := dialer.Dial(string(protocol), url)
	if err != nil {
		return err
	}

	header := []byte(fmt.Sprintf("Running on Pod %s.", podName))
	message := []byte(fmt.Sprintf("testing %sroute", protocol))

	wrote, err := conn.Write(message)
	if err != nil {
		return err
	}

	if wrote != len(message) {
		return fmt.Errorf("wrote message of size %d, expected %d", wrote, len(message))
	}

	if err := conn.SetDeadline(time.Now().Add(time.Second * 5)); err != nil {
		return err
	}

	headerResponse := make([]byte, len(header)+1)
	read, err := conn.Read(headerResponse)
	if err != nil {
		return err
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
		return err
	}

	if read != len(message) {
		return fmt.Errorf("read %d bytes but expected %d", read, len(message))
	}

	if !bytes.Contains(messageResponse, message) {
		return fmt.Errorf(`expected message response "%s", received: "%s"`, message, messageResponse)
	}

	return nil
}
