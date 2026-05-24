// Package levelfilter provides log-level based filtering.
// Lines whose parsed severity falls outside the configured
// minimum/maximum level range are dropped from the pipeline.
package levelfilter

import (
	"strings"

	"github.com/your-org/logslice/internal/formatter"
)

// Filter drops log lines whose level is outside [Min, Max].
type Filter struct {
	min formatter.Level
	max formatter.Level
	enabled bool
}

// New creates a Filter from human-readable level strings (e.g. "warn", "error").
// An empty string for either bound means "no bound on that side".
// If both strings are empty the filter is disabled.
func New(minLevel, maxLevel string) (*Filter, error) {
	if minLevel == "" && maxLevel == "" {
		return &Filter{enabled: false}, nil
	}

	min := formatter.LevelDebug
	max := formatter.LevelFatal

	if minLevel != "" {
		parsed, err := formatter.ParseLevel(strings.ToLower(minLevel))
		if err != nil {
			return nil, err
		}
		min = parsed
	}

	if maxLevel != "" {
		parsed, err := formatter.ParseLevel(strings.ToLower(maxLevel))
		if err != nil {
			return nil, err
		}
		max = parsed
	}

	return &Filter{
		min:     min,
		max:     max,
		enabled: true,
	}, nil
}

// Enabled reports whether the filter has any active bounds.
func (f *Filter) Enabled() bool { return f.enabled }

// Allow returns true when the given level string falls within [min, max].
// Lines with an unrecognised level string are always allowed through.
func (f *Filter) Allow(levelStr string) bool {
	if !f.enabled {
		return true
	}
	lvl, err := formatter.ParseLevel(strings.ToLower(levelStr))
	if err != nil {
		// Unknown level — let it pass so we don't silently drop unstructured lines.
		return true
	}
	return lvl >= f.min && lvl <= f.max
}
