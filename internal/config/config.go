// Package config holds the runtime configuration for logslice.
package config

import (
	"errors"
	"strings"
)

// Config holds all runtime options for a logslice run.
type Config struct {
	// Input / output
	Input  string
	Output string

	// Filtering
	Include []string
	Exclude []string

	// Display
	Color  bool
	Quiet  bool
	Follow bool

	// Sampling: 0 < SampleRate <= 1.0; 1.0 means keep all lines.
	SampleRate float64
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{
		Color:      true,
		SampleRate: 1.0,
	}
}

// Validate returns an error if the Config contains invalid values.
func (c *Config) Validate() error {
	for _, p := range c.Include {
		if strings.TrimSpace(p) == "" {
			return errors.New("include pattern must not be blank")
		}
	}
	for _, p := range c.Exclude {
		if strings.TrimSpace(p) == "" {
			return errors.New("exclude pattern must not be blank")
		}
	}
	if len(c.Include) == 0 && len(c.Exclude) == 0 && !c.Follow {
		if c.Input == "" && c.Output == "" {
			return errors.New("no filters or options specified")
		}
	}
	if c.SampleRate <= 0 || c.SampleRate > 1.0 {
		return errors.New("sample-rate must be in range (0, 1]")
	}
	return nil
}
