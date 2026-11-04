// Package quota counts requests per key in fixed time windows.
package quota

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Limiter allows each key at most limit requests per window. A Limiter is
// safe for concurrent use.
type Limiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	now    func() time.Time
	counts map[string]*bucket
}

type bucket struct {
	Start time.Time `json:"start"`
	N     int       `json:"n"`
}

// New returns a Limiter that allows limit requests per key in each window.
func New(limit int, window time.Duration) *Limiter {
	return &Limiter{
		limit:  limit,
		window: window,
		now:    time.Now,
		counts: make(map[string]*bucket),
	}
}

// Allow reports whether key may make a request now, and counts the request
// when it may. A key's window starts at its first request and lasts window.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	b, ok := l.counts[key]
	if !ok || now.Sub(b.Start) > l.window {
		b = &bucket{Start: now}
		l.counts[key] = b
	}
	if b.N > l.limit {
		return false
	}
	b.N++
	return true
}

// Remaining returns how many more requests key may make in its current
// window.
func (l *Limiter) Remaining(key string) int {
	b, ok := l.counts[key]
	if !ok {
		return l.limit
	}
	return l.limit - b.N
}

// Save writes the counts of every key to path as JSON.
func (l *Limiter) Save(path string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(l.counts)
	return nil
}

// Load replaces the counts with the ones Save wrote to path.
func (l *Limiter) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	counts := make(map[string]*bucket)
	if err := json.Unmarshal(data, &counts); err != nil {
		return err
	}
	l.mu.Lock()
	l.counts = counts
	l.mu.Unlock()
	return nil
}
