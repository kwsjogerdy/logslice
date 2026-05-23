package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/output"
)

func TestWriter_WriteLinePlain(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, false)

	if err := w.WriteLine("hello world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}

	got := buf.String()
	if got != "hello world\n" {
		t.Errorf("expected %q, got %q", "hello world\n", got)
	}
}

func TestWriter_WriteLineColorized(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, true)

	if err := w.WriteLine("colored line"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "colored line") {
		t.Errorf("output missing original text, got %q", got)
	}
	// Expect ANSI escape codes to be present.
	if !strings.Contains(got, "\033[") {
		t.Errorf("expected ANSI codes in colorized output, got %q", got)
	}
}

func TestWriter_MultipleLines(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, false)

	lines := []string{"line one", "line two", "line three"}
	for _, l := range lines {
		if err := w.WriteLine(l); err != nil {
			t.Fatalf("unexpected error writing %q: %v", l, err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}

	got := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(got) != len(lines) {
		t.Fatalf("expected %d lines, got %d", len(lines), len(got))
	}
	for i, want := range lines {
		if got[i] != want {
			t.Errorf("line %d: expected %q, got %q", i, want, got[i])
		}
	}
}
