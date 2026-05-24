// Package config handles loading and validation of logslice configuration
// from CLI flags, environment variables, and optional config files.
package config

import (
	"errors"
	"strings"
)

// Config holds all runtime configuration for a logslice run.
type Config struct {
	// Input
	InputFile string // empty means stdin

	// Output
	OutputFile string // empty means stdout
	Colorize   bool

	// Filtering
	Include []string
	Exclude []string

	// Formatting
	Level     string
	MinLevel  string
	TimeField string

	// Behaviour
	FollowMode bool // tail -f style
	Quiet      bool // suppress stats output
}

// Validate checks that the Config fields are consistent and returns an
// error describing the first problem found.
func (c *Config) Validate() error {
	for _, pat := range c.Include {
		if strings.TrimSpace(pat) == "" {
			return errors.New("config: include pattern must not be blank")
		}
	}
	for _, pat := range c.Exclude {
		if strings.TrimSpace(pat) == "" {
			return errors.New("config: exclude pattern must not be blank")
		}
	}
	return nil
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Colorize:  true,
		TimeField: "time",
		Level:     "level",
	}
}
