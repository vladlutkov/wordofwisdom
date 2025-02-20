package main

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"
)

func writeUInt8(conn net.Conn, v uint8) error {
	data := []byte{v}
	_, err := conn.Write(data)
	return err
}

func writeUInt32(conn net.Conn, v uint32) error {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, v)
	_, err := conn.Write(data)
	return err
}

func check(challenge []byte, nonce []byte, difficulty int) bool {
	solution := sha256.Sum256(append(challenge, nonce...))
	hashPrefix := binary.BigEndian.Uint32(solution[:4])
	if hashPrefix>>(31-difficulty) != 0 {
		return false
	}

	return true
}

func quote(random *rand.Rand, conn net.Conn, usr User) error {
	usr.Tick()

	difficulty := usr.Difficulty()
	err := writeUInt8(conn, uint8(difficulty))
	if err != nil {
		return fmt.Errorf("failed to send difficulty: %w", err)
	}

	challenge := make([]byte, 4)
	random.Read(challenge)
	_, err = conn.Write(challenge)
	if err != nil {
		return fmt.Errorf("failed to send challenge: %w", err)
	}

	nonce := make([]byte, 4)
	_, err = io.ReadFull(conn, nonce)
	if err != nil {
		return fmt.Errorf("failed to read nonce: %w", err)
	}

	if !check(challenge, nonce, difficulty) {
		// не надо наши загадки плохо решать!
		return errors.New("invalid proof of work")
	}

	q := QuotesGlob.Get()

	err = writeUInt32(conn, uint32(len(q)))
	if err != nil {
		return fmt.Errorf("failed to send quote length: %w", err)
	}

	_, err = conn.Write(q)
	if err != nil {
		return fmt.Errorf("failed to send quote: %w", err)
	}
	return nil
}

func handle(ctx context.Context, conn net.Conn, usr User) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			break
		default:
			err := quote(random, conn, usr)
			if err != nil {
				return
			}
		}
	}
}
