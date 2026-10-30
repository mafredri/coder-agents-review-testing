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

func TestLast(t *testing.T) {
	b := New(3)
	for v := range 5 {
		b.Push(v)
	}
	tests := []struct {
		n    int
		want []int
	}{
		{0, []int{}},
		{1, []int{4}},
		{3, []int{2, 3, 4}},
		{10, []int{2, 3, 4}},
	}
	for _, tt := range tests {
		if got := b.Last(tt.n); !slices.Equal(got, tt.want) {
			t.Errorf("Last(%d) = %v, want %v", tt.n, got, tt.want)
		}
	}
}

func TestSum(t *testing.T) {
	b := New(4)
	b.Push(2)
	b.Push(5)
	if got := b.Sum(); got != 7 {
		t.Errorf("Sum() = %d, want 7", got)
	}
}
