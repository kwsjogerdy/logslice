package redact

import (
	"regexp"
	"strings"
)

// Redactor replaces sensitive patterns in log lines with a placeholder.
type Redactor struct {
	patterns    []*regexp.Regexp
	placeholder string
}

// New creates a Redactor that replaces matches of the given regex patterns
// with placeholder. If placeholder is empty, "[REDACTED]" is used.
func New(patterns []string, placeholder string) (*Redactor, error) {
	if placeholder == "" {
		placeholder = "[REDACTED]"
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return &Redactor{patterns: compiled, placeholder: placeholder}, nil
}

// Enabled reports whether the Redactor has any active patterns.
func (r *Redactor) Enabled() bool {
	return len(r.patterns) > 0
}

// Apply replaces all pattern matches in line with the placeholder.
func (r *Redactor) Apply(line string) string {
	if !r.Enabled() {
		return line
	}
	for _, re := range r.patterns {
		line = re.ReplaceAllString(line, r.placeholder)
	}
	return line
}

// ApplyBytes is a convenience wrapper for []byte input.
func (r *Redactor) ApplyBytes(line []byte) []byte {
	if !r.Enabled() {
		return line
	}
	return []byte(r.Apply(strings.Clone(string(line))))
}
