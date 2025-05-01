package main

import (
	"errors"
	"io"
	"net"
	"time"
)

var ErrCantConnect = errors.New("cant connect")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}
type telnet struct {
	connect net.Conn
	address string
	in      io.ReadCloser
	out     io.Writer
	timeout time.Duration
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnet{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *telnet) Connect() error {
	connect, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return ErrCantConnect
	}

	t.connect = connect

	return nil
}

func (t *telnet) Send() error {
	_, err := io.Copy(t.connect, t.in)

	return err
}

func (t *telnet) Receive() error {
	_, err := io.Copy(t.out, t.connect)

	return err
}

func (t *telnet) Close() error {
	return t.connect.Close()
}
