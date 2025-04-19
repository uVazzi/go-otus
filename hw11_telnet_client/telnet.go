package main

import (
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

func (client *telnetClient) Connect() (err error) {
	client.conn, err = net.DialTimeout("tcp", client.address, client.timeout)
	return
}

func (client *telnetClient) Close() error {
	if client.conn != nil {
		return client.conn.Close()
	}
	return nil
}

func (client *telnetClient) Send() error {
	_, err := io.Copy(client.conn, client.in)
	return err
}

func (client *telnetClient) Receive() error {
	_, err := io.Copy(client.out, client.conn)
	return err
}
