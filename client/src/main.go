package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func readUInt8(reader io.Reader) (int, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return 0, err
	}

	return int(buf[0]), nil
}

func readUInt32(reader io.Reader) (int, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return 0, err
	}
	return int(binary.BigEndian.Uint32(buf)), nil
}

func quotes(addr string) error {
	conn, err := NewConnection(addr)
	if err != nil {
		return fmt.Errorf("connection failed: %v\n", err)
	}
	defer conn.Close()

	for {
		reader := conn.Reader()

		difficulty, err := readUInt8(reader)
		if err != nil {
			return fmt.Errorf("failed to read difficulty: %v\n", err)
		}

		challenge := make([]byte, 4)
		_, err = io.ReadFull(reader, challenge)
		if err != nil {
			return fmt.Errorf("failed to read challenge: %v\n", err)
		}
		log.Printf("difficulty: %d, challenge: %X\n", difficulty, challenge)

		start := time.Now()
		nonce := solve(difficulty, challenge)
		duration := time.Since(start)
		log.Printf("nonce: %X, calculation time: %v", nonce, duration)

		writer := conn.Writer()
		_, err = writer.Write(nonce)
		if err != nil {
			return fmt.Errorf("failed to write nonce: %v\n", err)
		}

		quoteLen, err := readUInt32(reader)
		if err != nil {
			return fmt.Errorf("failed to read quote length: %v\n", err)
		}

		quote := make([]byte, quoteLen)
		_, err = io.ReadFull(reader, quote)
		if err != nil {
			return fmt.Errorf("failed to read quote: %v\n", err)
		}
		log.Println("quote:", string(quote))
	}
}

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		log.Fatal("ADDR is not set")
	}

	log.Println("client started.", "addr:", addr)

	for {
		err := quotes(addr)
		if err != nil {
			log.Println("connection failed:", err)
			time.Sleep(2 * time.Second)
		}
	}
}
