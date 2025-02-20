package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	p := os.Getenv("PORT")
	port, err := strconv.Atoi(p)
	if err != nil {
		log.Fatal("failed to parse PORT env variable:", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)

	lstnr, err := NewListener(ctx, port)
	if err != nil {
		log.Fatal("failed to start TCP server:", err)
	}
	log.Printf("server started at 0.0.0.0:%d\n", port)
	defer lstnr.Close()

	usrs := NewUsers(ctx, 5, 5*time.Second)
	for {
		select {
		case conn := <-lstnr.Listen():
			usr := usrs.Register(conn.RemoteAddr().String())
			go handle(ctx, conn, usr)
			log.Println("new connection:", conn.RemoteAddr().String())
		case <-sigterm:
			cancel()
			return
		}
	}
}
