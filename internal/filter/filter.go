// Package filter provides log line filtering based on include/exclude patterns.
package filter

import (
	"regexp"
	"strings"
)

// Options holds configuration for a Filter.
type Options struct {
	// Include keeps only lines matching any of these patterns.
	Include []string
	// Exclude drops lines matching any of these patterns.
	Exclude []string
	// CaseSensitive controls whether pattern matching is case-sensitive.
	CaseSensitive bool
	// UseRegex treats patterns as regular expressions instead of plain substrings.
	UseRegex bool
}

// Filter evaluates log lines against include/exclude rules.
type Filter struct {
	opts     Options
	includes []*regexp.Regexp
	excludes []*regexp.Regexp
}

// New creates a Filter from the provided Options.
func New(opts Options) (*Filter, error) {
	f := &Filter{opts: opts}

	compile := func(patterns []string) ([]*regexp.Regexp, error) {
		var compiled []*regexp.Regexp
		for _, p := range patterns {
			if !opts.UseRegex {
				p = regexp.QuoteMeta(p)
			}
			if !opts.CaseSensitive {
				p = "(?i)" + p
			}
			re, err := regexp.Compile(p)
			if err != nil {
				return nil, err
			}
			compiled = append(compiled, re)
		}
		return compiled, nil
	}

	var err error
	if f.includes, err = compile(opts.Include); err != nil {
		return nil, err
	}
	if f.excludes, err = compile(opts.Exclude); err != nil {
		return nil, err
	}
	return f, nil
}

// Match reports whether the given log line should be kept.
func (f *Filter) Match(line string) bool {
	if !f.opts.CaseSensitive {
		_ = strings.ToLower(line) // normalisation handled by (?i) flag
	}

	for _, re := range f.excludes {
		if re.MatchString(line) {
			return false
		}
	}

	if len(f.includes) == 0 {
		return true
	}
	for _, re := range f.includes {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}
