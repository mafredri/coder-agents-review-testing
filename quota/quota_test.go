package quota

import (
	"path/filepath"
	"testing"
	"time"
)

func newTestLimiter(limit int, window time.Duration) (*Limiter, *time.Time) {
	l := New(limit, window)
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	l.now = func() time.Time { return now }
	return l, &now
}

func TestAllowFirstRequests(t *testing.T) {
	l, _ := newTestLimiter(3, time.Minute)
	for i := range 3 {
		if !l.Allow("a") {
			t.Fatalf("request %d denied", i+1)
		}
	}
}

func TestAllowNewWindow(t *testing.T) {
	l, now := newTestLimiter(1, time.Minute)
	l.Allow("a")
	*now = now.Add(2 * time.Minute)
	if !l.Allow("a") {
		t.Fatal("request in a new window denied")
	}
}

func TestRemaining(t *testing.T) {
	l, _ := newTestLimiter(5, time.Minute)
	if got := l.Remaining("a"); got != 5 {
		t.Fatalf("Remaining before any request = %d, want 5", got)
	}
	l.Allow("a")
	l.Allow("a")
	if got := l.Remaining("a"); got != 3 {
		t.Fatalf("Remaining after 2 requests = %d, want 3", got)
	}
}

func TestSaveLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "counts.json")
	l, _ := newTestLimiter(5, time.Minute)
	l.Allow("a")
	if err := l.Save(path); err != nil {
		t.Fatal(err)
	}
	l2, _ := newTestLimiter(5, time.Minute)
	if err := l2.Load(path); err != nil {
		t.Fatal(err)
	}
	if got := l2.Remaining("a"); got != 4 {
		t.Fatalf("Remaining after Load = %d, want 4", got)
	}
}
