// Пакет управляет пользователями, их запросами и динамической сложностью
package main

import (
	"context"
	"sync"
	"time"
)

type address string

type User interface {
	Difficulty() int
	Tick()
}

type user struct {
	difficulty int

	requestCount       int
	connectionAttempts int
}

func (u *user) Difficulty() int {
	return u.difficulty
}

func (u *user) Tick() {
	u.requestCount++
}

type Users interface {
	Register(addr string) User
}

type users struct {
	m                 sync.RWMutex
	u                 map[address]*user
	requestsThreshold int
}

func NewUsers(ctx context.Context, requestsCount int, inTime time.Duration) Users {
	usrs := &users{
		u:                 make(map[address]*user, 4),
		requestsThreshold: requestsCount,
	}

	go func() {
		ticker := time.NewTicker(inTime)
		for {
			select {
			case <-ticker.C:
				usrs.recalculate()
			case <-ctx.Done():
				return
			}
		}
	}()

	return usrs
}

// Корректирует сложность и сбрасывает счётчики.
func (u *users) recalculate() {
	u.m.Lock()
	defer u.m.Unlock()

	for _, usr := range u.u {
		if usr.requestCount > u.requestsThreshold {
			usr.difficulty++
		} else {
			usr.difficulty--
		}

		if usr.difficulty < 1 {
			usr.difficulty = 1
		}
		if usr.difficulty > 32 {
			usr.difficulty = 32
		}

		usr.connectionAttempts = 0
		usr.requestCount = 0
	}
}

func (u *users) Register(addr string) User {
	u.m.Lock()
	defer u.m.Unlock()

	usr, ok := u.u[address(addr)]
	if !ok {
		usr = &user{
			difficulty:         10,
			connectionAttempts: 1,
		}
		u.u[address(addr)] = usr
	} else {
		usr.connectionAttempts++
		holdOn := time.Duration(usr.connectionAttempts) * time.Second
		time.Sleep(holdOn)
	}
	return usr
}
