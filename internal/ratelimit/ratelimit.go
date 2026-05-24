package ratelimit

import (
	"sync"
	"time"
)

// Limiter enforces a maximum number of lines emitted per second.
// A rate of zero means unlimited.
type Limiter struct {
	mu       sync.Mutex
	rate     int           // lines per second
	bucket   int           // current token count
	lastFill time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows at most ratePerSec lines per second.
// If ratePerSec <= 0 the limiter is disabled (all lines pass through).
func New(ratePerSec int) *Limiter {
	if ratePerSec < 0 {
		ratePerSec = 0
	}
	return &Limiter{
		rate:     ratePerSec,
		bucket:   ratePerSec,
		lastFill: time.Now(),
		clock:    time.Now,
	}
}

// Allow returns true when the line should be emitted, false when it should
// be dropped to stay within the configured rate limit.
func (l *Limiter) Allow() bool {
	if l.rate == 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastFill)
	if elapsed >= time.Second {
		periods := int(elapsed / time.Second)
		l.bucket += periods * l.rate
		if l.bucket > l.rate {
			l.bucket = l.rate
		}
		l.lastFill = l.lastFill.Add(time.Duration(periods) * time.Second)
	}

	if l.bucket <= 0 {
		return false
	}
	l.bucket--
	return true
}

// Rate returns the configured lines-per-second limit (0 = unlimited).
func (l *Limiter) Rate() int {
	return l.rate
}
