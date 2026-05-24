// Package timefilter provides time-range filtering for log lines.
// It parses timestamps from structured or plain log lines and discards
// entries that fall outside the configured [After, Before] window.
package timefilter

import (
	"time"
)

const (
	// DefaultLayout is the fallback timestamp layout tried for plain text.
	DefaultLayout = time.RFC3339
)

// Filter holds an optional time window. A zero Time value means unbounded.
type Filter struct {
	after  time.Time
	before time.Time
}

// New creates a Filter. Either bound may be zero to indicate "no limit".
func New(after, before time.Time) *Filter {
	return &Filter{after: after, before: before}
}

// Enabled reports whether at least one bound is set.
func (f *Filter) Enabled() bool {
	return !f.after.IsZero() || !f.before.IsZero()
}

// Allow returns true when t falls within the configured window.
// A zero After means "since the beginning of time"; a zero Before means
// "until the end of time".
func (f *Filter) Allow(t time.Time) bool {
	if t.IsZero() {
		// Cannot determine timestamp — let the line through.
		return true
	}
	if !f.after.IsZero() && !t.After(f.after) {
		return false
	}
	if !f.before.IsZero() && !t.Before(f.before) {
		return false
	}
	return true
}

// ParseTimestamp attempts to extract a time.Time from a set of candidate
// string values (e.g. the "time", "ts", or "timestamp" fields of a parsed
// log entry). It tries RFC3339Nano, RFC3339, and a handful of common layouts.
// Returns zero time if no value can be parsed.
func ParseTimestamp(candidates []string) time.Time {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"02/Jan/2006:15:04:05 -0700",
	}
	for _, v := range candidates {
		if v == "" {
			continue
		}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, v); err == nil {
				return t
			}
		}
	}
	return time.Time{}
}
