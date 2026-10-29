// Package ringbuf implements a fixed-size ring buffer of ints.
package ringbuf

// Buffer holds the most recent values pushed to it, up to its capacity.
type Buffer struct {
	vals  []int
	start int
	n     int
}

// New returns an empty Buffer with the given capacity. It panics if
// capacity is not positive.
func New(capacity int) *Buffer {
	if capacity <= 0 {
		panic("ringbuf: capacity must be positive")
	}
	return &Buffer{vals: make([]int, capacity)}
}

// Push adds v, dropping the oldest value when the buffer is full.
func (b *Buffer) Push(v int) {
	if b.n < len(b.vals) {
		b.vals[(b.start+b.n)%len(b.vals)] = v
		b.n++
		return
	}
	b.vals[b.start] = v
	b.start = (b.start + 1) % len(b.vals)
}

// Len returns the number of values in the buffer.
func (b *Buffer) Len() int { return b.n }

// Values returns the values from oldest to newest.
func (b *Buffer) Values() []int {
	out := make([]int, b.n)
	for i := range b.n {
		out[i] = b.vals[(b.start+i)%len(b.vals)]
	}
	return out
}
