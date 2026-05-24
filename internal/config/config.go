// Package config defines the runtime configuration for logslice and
// provides helpers for building it from CLI flags or sensible defaults.
package config

import "errors"

// Config holds all settings that control a single logslice run.
type Config struct {
	// Input is the path to the log file to read. Empty means stdin.
	Input string
	// Output is the path to write filtered lines. Empty means stdout.
	Output string
	// Include holds substrings; a line must match at least one to pass
	// (unless the slice is empty, in which case all lines pass).
	Include []string
	// Exclude holds substrings; a line is dropped if it matches any.
	Exclude []string
	// CaseInsensitive makes both Include and Exclude matching case-blind.
	CaseInsensitive bool
	// Color enables ANSI colour in the output.
	Color bool
	// Highlight lists terms that should be visually emphasised in output.
	Highlight []string
	// Quiet suppresses the stats summary printed after processing.
	Quiet bool
	// Follow keeps the input file open and streams new lines (tail -f).
	Follow bool
}

// Default returns a Config populated with sensible out-of-the-box values.
func Default() Config {
	return Config{
		CaseInsensitive: true,
		Color:           true,
	}
}

// Validate checks that the Config is internally consistent and returns a
// descriptive error when it is not.
func (c Config) Validate() error {
	for _, inc := range c.Include {
		if inc == "" {
			return errors.New("config: include pattern must not be blank")
		}
	}
	for _, exc := range c.Exclude {
		if exc == "" {
			return errors.New("config: exclude pattern must not be blank")
		}
	}
	for _, hl := range c.Highlight {
		if hl == "" {
			return errors.New("config: highlight term must not be blank")
		}
	}
	if len(c.Include) == 0 && len(c.Exclude) == 0 && len(c.Highlight) == 0 {
		// A completely unconfigured run is valid — it acts as a passthrough.
		return nil
	}
	return nil
}
