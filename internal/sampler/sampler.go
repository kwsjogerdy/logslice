package sampler

import (
	"math/rand"
	"sync/atomic"
)

// Sampler decides whether a log line should be kept based on a sampling rate.
// A rate of 1.0 keeps all lines; 0.1 keeps approximately 10% of lines.
type Sampler struct {
	rate    float64
	counter atomic.Uint64
	rng     *rand.Rand
}

// New returns a Sampler with the given rate clamped to [0.0, 1.0].
// If deterministic is true a fixed seed is used (useful for tests).
func New(rate float64, deterministic bool) *Sampler {
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	var src rand.Source
	if deterministic {
		src = rand.NewSource(42)
	} else {
		src = rand.NewSource(rand.Int63())
	}
	return &Sampler{
		rate: rate,
		rng:  rand.New(src), //nolint:gosec
	}
}

// Keep reports whether the current line should be kept.
func (s *Sampler) Keep() bool {
	if s.rate >= 1.0 {
		return true
	}
	if s.rate <= 0.0 {
		return false
	}
	s.counter.Add(1)
	return s.rng.Float64() < s.rate
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() float64 {
	return s.rate
}

// Seen returns the total number of times Keep has been called.
func (s *Sampler) Seen() uint64 {
	return s.counter.Load()
}
