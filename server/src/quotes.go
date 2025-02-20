package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

type Quotes interface {
	Get() []byte
}

type quotes struct {
	r *rand.Rand
	l [][]byte
}

func NewQuotes() Quotes {
	s := rand.NewSource(time.Now().UnixNano())
	q := &quotes{
		l: make([][]byte, 0, 1000),
		r: rand.New(s),
	}

	file, err := os.Open("quotes.txt")
	if err != nil {
		log.Fatalf("failed to open quoutes file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatalf("error reading file: %v", err)
		}

		if len(line) > 0 {
			bytes.TrimSpace(line)
			q.l = append(q.l, line)
		}

		if err == io.EOF {
			break
		}
	}

	return q
}

func (q *quotes) Get() []byte {
	return q.l[q.r.Intn(len(q.l))]
}

var QuotesGlob = NewQuotes()
