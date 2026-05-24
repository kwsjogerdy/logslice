// Package truncate provides line-length truncation for log output.
// Lines exceeding the configured maximum byte length are clipped and
// a configurable suffix (e.g. "...") is appended so the reader can
// tell that content was omitted.
package truncate

import "unicode/utf8"

const (
	// DefaultMaxBytes is the default maximum line length in bytes.
	DefaultMaxBytes = 2048
	// DefaultSuffix is appended to truncated lines.
	DefaultSuffix = "..."
)

// Truncator clips lines that exceed a byte threshold.
type Truncator struct {
	maxBytes int
	suffix   []byte
}

// New returns a Truncator. maxBytes <= 0 disables truncation.
// suffix is appended to any clipped line; pass "" to use DefaultSuffix.
func New(maxBytes int, suffix string) *Truncator {
	if suffix == "" {
		suffix = DefaultSuffix
	}
	return &Truncator{
		maxBytes: maxBytes,
		suffix:   []byte(suffix),
	}
}

// Apply returns line unchanged when truncation is disabled or the line
// fits within the limit. Otherwise it returns a clipped copy with the
// suffix appended. The cut point is adjusted to never split a UTF-8
// rune.
func (t *Truncator) Apply(line []byte) []byte {
	if t.maxBytes <= 0 || len(line) <= t.maxBytes {
		return line
	}

	cutAt := t.maxBytes - len(t.suffix)
	if cutAt < 0 {
		cutAt = 0
	}

	// Walk back to a valid rune boundary.
	for cutAt > 0 && !utf8.RuneStart(line[cutAt]) {
		cutAt--
	}

	out := make([]byte, cutAt+len(t.suffix))
	copy(out, line[:cutAt])
	copy(out[cutAt:], t.suffix)
	return out
}

// Enabled reports whether truncation is active.
func (t *Truncator) Enabled() bool { return t.maxBytes > 0 }
