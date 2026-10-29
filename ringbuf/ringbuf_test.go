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
