package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	tmp := filepath.Join(t.TempDir(), "input.log")
	if err := os.WriteFile(tmp, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return tmp
}

func TestPipeline_BasicRun(t *testing.T) {
	input := writeTempFile(t, "hello world\nfoo bar\n")
	output := filepath.Join(t.TempDir(), "out.log")

	cfg := Config{
		InputFile:  input,
		OutputFile: output,
	}
	p, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	stats, err := p.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if stats.Read != 2 {
		t.Errorf("expected Read=2, got %d", stats.Read)
	}
	if stats.Written != 2 {
		t.Errorf("expected Written=2, got %d", stats.Written)
	}
}

func TestPipeline_FilterExcludes(t *testing.T) {
	input := writeTempFile(t, "keep this line\ndrop secret line\nkeep another\n")
	output := filepath.Join(t.TempDir(), "out.log")

	cfg := Config{
		Exclude:    []string{"secret"},
		InputFile:  input,
		OutputFile: output,
	}
	p, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	stats, err := p.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if stats.Written != 2 {
		t.Errorf("expected Written=2, got %d", stats.Written)
	}
	if stats.Filtered != 1 {
		t.Errorf("expected Filtered=1, got %d", stats.Filtered)
	}
}

// TestPipeline_OutputContents verifies that the lines written to the output
// file match the lines that were not excluded by the filter.
func TestPipeline_OutputContents(t *testing.T) {
	input := writeTempFile(t, "keep this line\ndrop secret line\nkeep another\n")
	output := filepath.Join(t.TempDir(), "out.log")

	cfg := Config{
		Exclude:    []string{"secret"},
		InputFile:  input,
		OutputFile: output,
	}
	p, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := p.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}

	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	got := string(data)
	if strings.Contains(got, "secret") {
		t.Errorf("output should not contain excluded line, got: %q", got)
	}
	if !strings.Contains(got, "keep this line") || !strings.Contains(got, "keep another") {
		t.Errorf("output missing expected lines, got: %q", got)
	}
}

func TestStats_String(t *testing.T) {
	s := Stats{Read: 10, Filtered: 3, Written: 7}
	got := s.String()
	if !strings.Contains(got, "read=10") || !strings.Contains(got, "written=7") {
		t.Errorf("unexpected Stats.String(): %q", got)
	}
}

func TestStats_FilterRatio(t *testing.T) {
	s := Stats{Read: 4, Filtered: 1, Written: 3}
	if got := s.FilterRatio(); got != 0.25 {
		t.Errorf("expected 0.25, got %f", got)
	}
	empty := Stats{}
	if got := empty.FilterRatio(); got != 0 {
		t.Errorf("expected 0 for empty stats, got %f", got)
	}
}
