package parser_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func TestParse_PlainText(t *testing.T) {
	line := "this is a plain log line"
	e := parser.Parse(line)
	if e.Format != parser.FormatPlain {
		t.Fatalf("expected FormatPlain, got %d", e.Format)
	}
	if e.Message != line {
		t.Fatalf("expected message %q, got %q", line, e.Message)
	}
	if e.Raw != line {
		t.Fatalf("expected raw %q, got %q", line, e.Raw)
	}
}

func TestParse_JSON(t *testing.T) {
	line := `{"level":"info","msg":"server started","port":"8080"}`
	e := parser.Parse(line)
	if e.Format != parser.FormatJSON {
		t.Fatalf("expected FormatJSON, got %d", e.Format)
	}
	if e.Level != "INFO" {
		t.Fatalf("expected level INFO, got %q", e.Level)
	}
	if e.Message != "server started" {
		t.Fatalf("expected message 'server started', got %q", e.Message)
	}
	if e.Fields["port"] != "8080" {
		t.Fatalf("expected port 8080, got %q", e.Fields["port"])
	}
}

func TestParse_JSONSeverityAlias(t *testing.T) {
	line := `{"severity":"error","message":"disk full"}`
	e := parser.Parse(line)
	if e.Level != "ERROR" {
		t.Fatalf("expected ERROR, got %q", e.Level)
	}
	if e.Message != "disk full" {
		t.Fatalf("expected 'disk full', got %q", e.Message)
	}
}

func TestParse_Logfmt(t *testing.T) {
	line := `level=warn msg="connection refused" host=db01`
	e := parser.Parse(line)
	if e.Format != parser.FormatLogfmt {
		t.Fatalf("expected FormatLogfmt, got %d", e.Format)
	}
	if e.Level != "WARN" {
		t.Fatalf("expected WARN, got %q", e.Level)
	}
	if e.Message != "connection refused" {
		t.Fatalf("expected 'connection refused', got %q", e.Message)
	}
	if e.Fields["host"] != "db01" {
		t.Fatalf("expected host db01, got %q", e.Fields["host"])
	}
}

func TestParse_LogfmtNoKeyValue(t *testing.T) {
	line := "just some words without equals"
	e := parser.Parse(line)
	if e.Format != parser.FormatPlain {
		t.Fatalf("expected FormatPlain fallback, got %d", e.Format)
	}
}

func TestParse_EmptyLine(t *testing.T) {
	e := parser.Parse("")
	if e.Message != "" {
		t.Fatalf("expected empty message, got %q", e.Message)
	}
}

func TestParse_RawPreserved(t *testing.T) {
	raw := `  {"msg":"hello"}  `
	e := parser.Parse(raw)
	if e.Raw != raw {
		t.Fatalf("expected raw to be preserved, got %q", e.Raw)
	}
}
