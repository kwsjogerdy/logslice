package config

import (
	"flag"
	"strings"
)

// multiStringFlag collects repeated flag values into a slice.
type multiStringFlag []string

func (m *multiStringFlag) String() string  { return strings.Join(*m, ",") }
func (m *multiStringFlag) Set(v string) error { *m = append(*m, v); return nil }

// RegisterRedactFlags registers --redact and --redact-placeholder flags on fs.
func RegisterRedactFlags(fs *flag.FlagSet, patterns *[]string, placeholder *string) {
	fs.Func("redact", "regex pattern to redact (repeatable)", func(v string) error {
		*patterns = append(*patterns, v)
		return nil
	})
	fs.StringVar(placeholder, "redact-placeholder", "",
		"replacement text for redacted values (default \"[REDACTED]\")",
	)
}

// ApplyRedact copies redaction settings from raw flag values into cfg.
func ApplyRedact(cfg *Config, patterns []string, placeholder string) {
	cfg.RedactPatterns = patterns
	cfg.RedactPlaceholder = placeholder
}
