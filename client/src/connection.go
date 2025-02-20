package main

import (
	"io"
	"log"
	"net"
)

type Connect interface {
	Reader() io.Reader
	Writer() io.Writer
	Close() error
}

func NewConnection(addr string) (Connect, error) {
	c := &connect{}

	err := c.dial(addr)

	return c, err
}

type connect struct {
	conn   net.Conn
	reader io.Reader
	writer io.Writer
}

func (c *connect) dial(addr string) error {
	var conn net.Conn
	var err error

	log.Println("trying to connect:", addr)
	conn, err = net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

func (c *connect) Reader() io.Reader {
	return c.conn
}

func (c *connect) Writer() io.Writer {
	return c.conn
}

func (c *connect) Close() error {
	return c.conn.Close()
}
