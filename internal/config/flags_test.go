package config

import (
	"flag"
	"testing"
)

func newFS() *flag.FlagSet {
	return flag.NewFlagSet("test", flag.ContinueOnError)
}

func TestFromFlags_Defaults(t *testing.T) {
	cfg, err := FromFlags(newFS(), []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Colorize {
		t.Error("expected Colorize=true")
	}
	if cfg.TimeField != "time" {
		t.Errorf("expected TimeField=\"time\", got %q", cfg.TimeField)
	}
}

func TestFromFlags_InputOutput(t *testing.T) {
	cfg, err := FromFlags(newFS(), []string{"-f", "app.log", "-o", "out.log"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.InputFile != "app.log" {
		t.Errorf("expected InputFile=\"app.log\", got %q", cfg.InputFile)
	}
	if cfg.OutputFile != "out.log" {
		t.Errorf("expected OutputFile=\"out.log\", got %q", cfg.OutputFile)
	}
}

func TestFromFlags_MultipleIncludeExclude(t *testing.T) {
	args := []string{"-include", "ERROR", "-include", "WARN", "-exclude", "healthcheck"}
	cfg, err := FromFlags(newFS(), args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Include) != 2 {
		t.Errorf("expected 2 include patterns, got %d", len(cfg.Include))
	}
	if len(cfg.Exclude) != 1 {
		t.Errorf("expected 1 exclude pattern, got %d", len(cfg.Exclude))
	}
}

func TestFromFlags_FollowAndQuiet(t *testing.T) {
	cfg, err := FromFlags(newFS(), []string{"-follow", "-quiet"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.FollowMode {
		t.Error("expected FollowMode=true")
	}
	if !cfg.Quiet {
		t.Error("expected Quiet=true")
	}
}
