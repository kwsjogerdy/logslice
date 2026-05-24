// Command logslice is a fast log filtering tool that supports structured
// and unstructured log formats.
//
// Usage:
//
//	logslice [flags] [file ...]
//
// If no file is given, logslice reads from standard input.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/config"
	"github.com/yourorg/logslice/internal/pipeline"
)

const version = "0.1.0"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: logslice [flags] [file ...]\n\n")
		fmt.Fprintf(os.Stderr, "A fast log filtering tool for structured and unstructured log formats.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fs.PrintDefaults()
	}

	var showVersion bool
	fs.BoolVar(&showVersion, "version", false, "print version and exit")

	cfg, err := config.FromFlags(fs, args)
	if err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	if showVersion {
		fmt.Printf("logslice %s\n", version)
		return nil
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	p, err := pipeline.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialise pipeline: %w", err)
	}

	stats, err := p.Run()
	if err != nil {
		return fmt.Errorf("pipeline error: %w", err)
	}

	if !cfg.Quiet {
		fmt.Fprintf(os.Stderr, "%s\n", stats)
	}

	return nil
}
