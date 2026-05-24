// Package dedup provides line-level deduplication for log streams.
// It tracks recently seen lines using a fixed-size ring buffer of hashes
// and suppresses consecutive or near-consecutive duplicate entries.
package dedup

import (
	"hash/fnv"
)

// Deduplicator tracks seen log lines and reports whether a line is a duplicate.
type Deduplicator struct {
	window []uint64
	size   int
	head   int
	count  int
}

// New creates a Deduplicator that remembers the last windowSize unique lines.
// windowSize must be at least 1; values less than 1 are clamped to 1.
func New(windowSize int) *Deduplicator {
	if windowSize < 1 {
		windowSize = 1
	}
	return &Deduplicator{
		window: make([]uint64, windowSize),
		size:   windowSize,
	}
}

// IsDuplicate returns true if line was seen within the current window.
// If the line is new it is recorded and false is returned.
func (d *Deduplicator) IsDuplicate(line string) bool {
	h := hash(line)
	for i := 0; i < d.count; i++ {
		idx := (d.head - 1 - i + d.size) % d.size
		if d.window[idx] == h {
			return true
		}
	}
	d.record(h)
	return false
}

// Seen returns the number of entries currently held in the window.
func (d *Deduplicator) Seen() int {
	return d.count
}

func (d *Deduplicator) record(h uint64) {
	d.window[d.head] = h
	d.head = (d.head + 1) % d.size
	if d.count < d.size {
		d.count++
	}
}

func hash(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}
