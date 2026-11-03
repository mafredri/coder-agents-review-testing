package ringbuf

import (
	"slices"
	"testing"
)

func TestPushDropsOldest(t *testing.T) {
	b := New(3)
	for v := range 5 {
		b.Push(v)
	}
	if got, want := b.Values(), []int{2, 3, 4}; !slices.Equal(got, want) {
		t.Fatalf("Values() = %v, want %v", got, want)
	}
}

func TestClamp(t *testing.T) {
	b := New(3)
	for _, v := range []int{-5, 2, 9, 4} {
		b.Push(v)
	}
	b.Clamp(0, 5)
	if got, want := b.Values(), []int{2, 5, 4}; !slices.Equal(got, want) {
		t.Errorf("Values() after Clamp(0, 5) = %v, want %v", got, want)
	}
}
