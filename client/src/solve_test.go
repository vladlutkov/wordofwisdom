package main

import (
	"crypto/sha256"
	"encoding/binary"
	"testing"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		name       string
		difficulty int
		challenge  []byte
	}{
		{"easyChallenge", 1, []byte("hello")},
		{"moderateChallenge", 3, []byte("challenge")},
		{"hardChallenge", 5, []byte("GoSolve")},
		{"emptyChallenge", 2, []byte("")},
		{"highDifficulty", 10, []byte("high difficulty")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nonce := solve(tt.difficulty, tt.challenge)

			if nonce == nil {
				t.Fatalf("Expected a nonce, got nil")
			}

			hash := sha256.Sum256(append(tt.challenge, nonce...))
			hashPrefix := binary.BigEndian.Uint32(hash[:4])
			if hashPrefix>>(31-tt.difficulty) != 0 {
				t.Errorf("Expected at least %d leading zero bits, got %032b...", tt.difficulty, hashPrefix)
			}
		})
	}
}
