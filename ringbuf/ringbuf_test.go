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

func TestLastPartial(t *testing.T) {
	b := New(4)
	if got := b.Last(2); len(got) != 0 {
		t.Errorf("empty Last(2) = %v, want []", got)
	}
	b.Push(7)
	b.Push(8)
	if got, want := b.Last(1), []int{8}; !slices.Equal(got, want) {
		t.Errorf("Last(1) = %v, want %v", got, want)
	}
	if got, want := b.Last(5), []int{7, 8}; !slices.Equal(got, want) {
		t.Errorf("Last(5) = %v, want %v", got, want)
	}
}

func TestLastIsCopy(t *testing.T) {
	b := New(2)
	b.Push(1)
	b.Push(2)
	got := b.Last(2)
	b.Push(3)
	if want := []int{1, 2}; !slices.Equal(got, want) {
		t.Errorf("Last(2) after Push = %v, want %v", got, want)
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		name   string
		pushes []int
		want   int
	}{
		{"empty", nil, 0},
		{"partial", []int{2, 5}, 7},
		{"wrapped", []int{1, 2, 3, 4, 5}, 12},
	}
	for _, tt := range tests {
		b := New(3)
		for _, v := range tt.pushes {
			b.Push(v)
		}
		if got := b.Sum(); got != tt.want {
			t.Errorf("%s: Sum() = %d, want %d", tt.name, got, tt.want)
		}
	}
}
