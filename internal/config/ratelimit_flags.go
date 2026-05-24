package config

import "flag"

// RegisterRateLimitFlags registers the --rate flag on the supplied FlagSet.
// It returns a pointer to the integer value that will be populated after
// flag.Parse is called.
func RegisterRateLimitFlags(fs *flag.FlagSet) *int {
	return fs.Int("rate", 0, "maximum lines per second to emit (0 = unlimited)")
}

// ApplyRateLimit copies the parsed rate-limit value into cfg.
func ApplyRateLimit(cfg *Config, rate int) {
	cfg.RateLimit = rate
}
