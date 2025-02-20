package main

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"time"
)

func try(r *rand.Rand, difficulty int, challenge []byte) []byte {
	nonce := make([]byte, 4)
	r.Read(nonce)

	solution := append(challenge, nonce...)
	hash := sha256.Sum256(solution)

	hashPrefix := binary.BigEndian.Uint32(hash[:4])
	if hashPrefix>>(31-difficulty) != 0 {
		return nil
	}

	return nonce
}

func solve(difficulty int, challenge []byte) []byte {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	for {
		nonce := try(random, difficulty, challenge)
		if nonce == nil {
			continue
		}
		return nonce
	}
}
