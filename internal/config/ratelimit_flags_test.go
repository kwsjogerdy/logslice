package config

import (
	"flag"
	"testing"
)

func TestRegisterRateLimitFlags_Default(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	rate := RegisterRateLimitFlags(fs)
	if err := fs.Parse([]string{}); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if *rate != 0 {
		t.Fatalf("expected default rate 0, got %d", *rate)
	}
}

func TestRegisterRateLimitFlags_Explicit(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	rate := RegisterRateLimitFlags(fs)
	if err := fs.Parse([]string{"-rate", "500"}); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if *rate != 500 {
		t.Fatalf("expected rate 500, got %d", *rate)
	}
}

func TestApplyRateLimit_SetsField(t *testing.T) {
	cfg := Default()
	ApplyRateLimit(cfg, 250)
	if cfg.RateLimit != 250 {
		t.Fatalf("expected RateLimit 250, got %d", cfg.RateLimit)
	}
}

func TestApplyRateLimit_Zero_Unlimited(t *testing.T) {
	cfg := Default()
	ApplyRateLimit(cfg, 0)
	if cfg.RateLimit != 0 {
		t.Fatalf("expected RateLimit 0 (unlimited), got %d", cfg.RateLimit)
	}
}
