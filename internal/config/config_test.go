package config

import (
	"testing"
)

func TestDefault_Fields(t *testing.T) {
	cfg := Default()
	if !cfg.Colorize {
		t.Error("expected Colorize to be true by default")
	}
	if cfg.TimeField != "time" {
		t.Errorf("expected TimeField=\"time\", got %q", cfg.TimeField)
	}
	if cfg.Level != "level" {
		t.Errorf("expected Level=\"level\", got %q", cfg.Level)
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := &Config{
		Include: []string{"ERROR", "WARN"},
		Exclude: []string{"healthcheck"},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_BlankInclude(t *testing.T) {
	cfg := &Config{Include: []string{"  "}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for blank include pattern")
	}
}

func TestValidate_BlankExclude(t *testing.T) {
	cfg := &Config{Exclude: []string{""}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for blank exclude pattern")
	}
}

func TestValidate_EmptyConfig(t *testing.T) {
	cfg := &Config{}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("empty config should be valid, got: %v", err)
	}
}
