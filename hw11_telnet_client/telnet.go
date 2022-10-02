package main

import (
	"context"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetConn struct {
	ctx     context.Context
	addr    string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(
	ctx context.Context,
	address string,
	timeout time.Duration,
	in io.ReadCloser,
	out io.Writer,
) TelnetClient {
	return &TelnetConn{
		ctx:     ctx,
		addr:    address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (tc *TelnetConn) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.addr, tc.timeout)
	tc.conn = conn
	return err
}

func (tc *TelnetConn) Close() error {
	return tc.conn.Close()
}

func (tc *TelnetConn) Send() error {
	_, err := io.Copy(tc.conn, tc.in)
	return err
}

func (tc *TelnetConn) Receive() error {
	_, err := io.Copy(tc.out, tc.conn)
	return err
}
