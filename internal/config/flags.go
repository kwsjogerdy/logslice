package config

import (
	"flag"
	"strings"
)

// multiFlag allows a flag to be specified multiple times.
type multiFlag []string

func (m *multiFlag) String() string  { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error { *m = append(*m, v); return nil }

// FromFlags parses os.Args using the provided FlagSet and returns a populated
// Config. Call fs.Parse(args) before calling FromFlags when testing.
func FromFlags(fs *flag.FlagSet, args []string) (*Config, error) {
	cfg := Default()

	var include, exclude multiFlag

	fs.StringVar(&cfg.InputFile, "f", "", "input log file (default: stdin)")
	fs.StringVar(&cfg.OutputFile, "o", "", "output file (default: stdout)")
	fs.BoolVar(&cfg.Colorize, "color", true, "colorize output")
	fs.BoolVar(&cfg.FollowMode, "follow", false, "follow file (tail -f mode)")
	fs.BoolVar(&cfg.Quiet, "quiet", false, "suppress stats on exit")
	fs.StringVar(&cfg.MinLevel, "level", "", "minimum log level to include")
	fs.StringVar(&cfg.TimeField, "time-field", "time", "JSON/logfmt field name for timestamp")
	fs.Var(&include, "include", "include lines matching pattern (repeatable)")
	fs.Var(&exclude, "exclude", "exclude lines matching pattern (repeatable)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg.Include = []string(include)
	cfg.Exclude = []string(exclude)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}
