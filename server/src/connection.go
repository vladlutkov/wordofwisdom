package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
)

type Listener interface {
	Listen() <-chan net.Conn
	Close() error
}

type listener struct {
	l net.Listener
	c chan net.Conn
}

func NewListener(ctx context.Context, port int) (Listener, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	// зачем аксептить соедиенение, если из канала не читают, но должны?
	// поэтому без буффера
	c := make(chan net.Conn)
	lstnr := &listener{
		l: l,
		c: c,
	}
	go lstnr.start(ctx)
	return lstnr, nil
}

func (l *listener) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			conn, err := l.l.Accept()
			// чтобы не печатать лишних ошибок
			if err != nil && strings.HasSuffix(err.Error(), "use of closed network connection") {
				return
			}
			if err != nil {
				log.Println("client connection failure:", err)
				continue
			}

			l.c <- conn
		}
	}
}

func (l *listener) Listen() <-chan net.Conn {
	return l.c
}

func (l *listener) Close() error {
	return l.l.Close()
}
