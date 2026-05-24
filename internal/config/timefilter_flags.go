package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/yourorg/logslice/internal/timefilter"
)

// RegisterTimeFilterFlags adds --after and --before flags to fs.
// Values must be RFC3339 strings, e.g. "2024-01-01T00:00:00Z".
func RegisterTimeFilterFlags(fs *flag.FlagSet) {
	fs.String("after", "", "only show log lines with a timestamp after this value (RFC3339)")
	fs.String("before", "", "only show log lines with a timestamp before this value (RFC3339)")
}

// ApplyTimeFilter reads --after and --before from fs and, if set, constructs
// a *timefilter.Filter and assigns it to cfg.TimeFilter.
// Returns an error if either value cannot be parsed.
func ApplyTimeFilter(fs *flag.FlagSet, cfg *Config) error {
	afterStr := flagString(fs, "after")
	beforeStr := flagString(fs, "before")

	var after, before time.Time

	if afterStr != "" {
		t, err := time.Parse(time.RFC3339, afterStr)
		if err != nil {
			return fmt.Errorf("--after: %w", err)
		}
		after = t
	}
	if beforeStr != "" {
		t, err := time.Parse(time.RFC3339, beforeStr)
		if err != nil {
			return fmt.Errorf("--before: %w", err)
		}
		before = t
	}

	cfg.TimeFilter = timefilter.New(after, before)
	return nil
}

// flagString is a helper that safely retrieves a string flag value.
func flagString(fs *flag.FlagSet, name string) string {
	f := fs.Lookup(name)
	if f == nil {
		return ""
	}
	return f.Value.String()
}
