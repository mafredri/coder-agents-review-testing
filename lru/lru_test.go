package lru

import (
	"strings"
	"testing"
)

func TestPutGet(t *testing.T) {
	c := New(2)
	c.Put("a", "1")
	if v, ok := c.Get("a"); !ok || v != "1" {
		t.Fatalf("Get(a) = %q, %v, want 1, true", v, ok)
	}
	if _, ok := c.Get("b"); ok {
		t.Fatal("Get(b) found a key that was never put")
	}
}

func TestPutReplaces(t *testing.T) {
	c := New(2)
	c.Put("a", "1")
	c.Put("a", "2")
	if v, _ := c.Get("a"); v != "2" {
		t.Fatalf("Get(a) = %q, want 2", v)
	}
	if got := c.Len(); got != 1 {
		t.Fatalf("Len = %d, want 1", got)
	}
}

func TestNewestSurvives(t *testing.T) {
	c := New(2)
	c.Put("a", "1")
	c.Put("b", "2")
	c.Put("c", "3")
	if _, ok := c.Get("c"); !ok {
		t.Fatal("newest key was evicted")
	}
}

func TestHitRatio(t *testing.T) {
	c := New(2)
	c.Put("a", "1")
	c.Get("a")
	c.Get("b")
	if got := c.HitRatio(); got != 0.5 {
		t.Fatalf("HitRatio = %v, want 0.5", got)
	}
}

func TestDump(t *testing.T) {
	c := New(2)
	c.Put("a", "1")
	c.Put("b", "2")
	var sb strings.Builder
	if err := c.Dump(&sb); err != nil {
		t.Fatal(err)
	}
	if got, want := sb.String(), "b=2\na=1\n"; got != want {
		t.Fatalf("Dump = %q, want %q", got, want)
	}
}

func TestNewPanicsOnNonPositive(t *testing.T) {
	for _, capacity := range []int{0, -1} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("New(%d) did not panic", capacity)
				}
			}()
			New(capacity)
		}()
	}
}
