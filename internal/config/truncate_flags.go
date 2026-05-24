package config

import (
	"flag"

	"github.com/your-org/logslice/internal/truncate"
)

// RegisterTruncateFlags adds --max-line-bytes and --truncate-suffix to fs.
func RegisterTruncateFlags(fs *flag.FlagSet) {
	fs.Int("max-line-bytes", truncate.DefaultMaxBytes,
		"maximum bytes per output line (0 = unlimited)")
	fs.String("truncate-suffix", truncate.DefaultSuffix,
		"string appended to truncated lines")
}

// ApplyTruncate reads truncation flags from fs and stores them in cfg.
func ApplyTruncate(fs *flag.FlagSet, cfg *Config) {
	if f := fs.Lookup("max-line-bytes"); f != nil {
		if v, ok := f.Value.(flag.Getter); ok {
			cfg.MaxLineBytes = v.Get().(int)
		}
	}
	if f := fs.Lookup("truncate-suffix"); f != nil {
		if v, ok := f.Value.(flag.Getter); ok {
			cfg.TruncateSuffix = v.Get().(string)
		}
	}
}
