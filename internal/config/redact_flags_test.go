package config_test

import (
	"flag"
	"testing"

	"github.com/yourorg/logslice/internal/config"
)

func TestRegisterRedactFlags_Default(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var patterns []string
	var placeholder string
	config.RegisterRedactFlags(fs, &patterns, &placeholder)

	if err := fs.Parse([]string{}); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(patterns) != 0 {
		t.Fatalf("expected no patterns, got %v", patterns)
	}
	if placeholder != "" {
		t.Fatalf("expected empty placeholder, got %q", placeholder)
	}
}

func TestRegisterRedactFlags_SinglePattern(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var patterns []string
	var placeholder string
	config.RegisterRedactFlags(fs, &patterns, &placeholder)

	if err := fs.Parse([]string{"--redact", `\d+`}); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(patterns) != 1 || patterns[0] != `\d+` {
		t.Fatalf("expected [\\d+], got %v", patterns)
	}
}

func TestRegisterRedactFlags_MultiplePatterns(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var patterns []string
	var placeholder string
	config.RegisterRedactFlags(fs, &patterns, &placeholder)

	err := fs.Parse([]string{"--redact", `password=\S+`, "--redact", `token=\S+`})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(patterns) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(patterns))
	}
}

func TestApplyRedact_SetsFields(t *testing.T) {
	cfg := config.Default()
	config.ApplyRedact(cfg, []string{`secret`}, "***")
	if len(cfg.RedactPatterns) != 1 || cfg.RedactPatterns[0] != `secret` {
		t.Fatalf("unexpected RedactPatterns: %v", cfg.RedactPatterns)
	}
	if cfg.RedactPlaceholder != "***" {
		t.Fatalf("unexpected RedactPlaceholder: %q", cfg.RedactPlaceholder)
	}
}
