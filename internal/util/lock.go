package util

import (
	"sync"
)

type Locks struct {
	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

func (l *Locks) GetOrCreate(id string) *sync.Mutex {
	l.mu.Lock()
	defer l.mu.Unlock()
	mu, ok := l.locks[id]
	if !ok {
		mu = &sync.Mutex{}
		l.locks[id] = mu
	}
	return mu
}

func NewLocks() *Locks {
	return &Locks{
		locks: make(map[string]*sync.Mutex),
	}
}
