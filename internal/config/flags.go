package config

import (
	"flag"
	"strings"
)

type multiFlag []string

func (m *multiFlag) String() string  { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error { *m = append(*m, v); return nil }

// RegisterFlags registers all logslice flags on fs and returns pointers that
// FromFlags will read.
func RegisterFlags(fs *flag.FlagSet) {
	fs.String("input", "", "input log file (default: stdin)")
	fs.String("output", "", "output file (default: stdout)")
	fs.Bool("color", true, "colorize output")
	fs.Bool("quiet", false, "suppress stats summary")
	fs.Bool("follow", false, "follow file for new lines (like tail -f)")
	fs.Float64("sample-rate", 1.0, "fraction of lines to keep (0,1]")
}

// FromFlags builds a Config from a parsed FlagSet.
// It expects flags registered by RegisterFlags plus any -include/-exclude
// flags added externally.
func FromFlags(fs *flag.FlagSet) Config {
	c := Default()

	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "input":
			c.Input = f.Value.String()
		case "output":
			c.Output = f.Value.String()
		case "color":
			c.Color = f.Value.String() == "true"
		case "quiet":
			c.Quiet = f.Value.String() == "true"
		case "follow":
			c.Follow = f.Value.String() == "true"
		case "include":
			c.Include = append(c.Include, f.Value.String())
		case "exclude":
			c.Exclude = append(c.Exclude, f.Value.String())
		case "sample-rate":
			// parsed as float64 by flag package
			if g, ok := f.Value.(flag.Getter); ok {
				if v, ok2 := g.Get().(float64); ok2 {
					c.SampleRate = v
				}
			}
		}
	})
	return c
}
