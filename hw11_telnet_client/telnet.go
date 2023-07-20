package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (tc *telnetClient) Connect() error {
	dialer := &net.Dialer{
		Timeout: tc.timeout,
	}

	conn, err := dialer.Dial("tcp", tc.address)
	if err != nil {
		log.Printf("Failed to connect %v: %v", tc.address, err)
		return err
	}
	tc.conn = conn
	return nil
}

func (tc *telnetClient) Close() error {
	return tc.conn.Close()
}

func (tc *telnetClient) Send() error {
	if tc.conn == nil {
		return fmt.Errorf("error send to %v, connection closed", tc.address)
	}
	if _, err := io.Copy(tc.conn, tc.in); err != nil {
		return fmt.Errorf("error send to %v: %w", tc.address, err)
	}
	return nil
}

func (tc *telnetClient) Receive() error {
	if tc.conn == nil {
		return fmt.Errorf("error receive from %v, connection closed", tc.address)
	}

	if _, err := io.Copy(tc.out, tc.conn); err != nil {
		return fmt.Errorf("error send to %v: %w", tc.address, err)
	}
	return nil
}
