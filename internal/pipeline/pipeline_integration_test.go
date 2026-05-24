package pipeline_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/config"
	"github.com/yourorg/logslice/internal/pipeline"
)

// writeLines writes a slice of log lines to a temp file and returns the path.
func writeLines(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logs-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	_, err = f.WriteString(strings.Join(lines, "\n") + "\n")
	if err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return f.Name()
}

// TestIntegration_PlainTextInclude verifies that plain-text include filtering
// passes only lines matching the include pattern.
func TestIntegration_PlainTextInclude(t *testing.T) {
	input := []string{
		"2024-01-01 INFO  starting server",
		"2024-01-01 ERROR database connection failed",
		"2024-01-01 INFO  request received",
		"2024-01-01 WARN  slow query detected",
	}
	inFile := writeLines(t, input)
	outFile := filepath.Join(t.TempDir(), "out.log")

	cfg := config.Default()
	cfg.InputFile = inFile
	cfg.OutputFile = outFile
	cfg.Include = []string{"ERROR"}
	cfg.NoColor = true

	p, err := pipeline.New(cfg)
	if err != nil {
		t.Fatalf("pipeline.New: %v", err)
	}
	if err := p.Run(); err != nil {
		t.Fatalf("pipeline.Run: %v", err)
	}

	got, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(got)), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 output line, got %d: %v", len(lines), lines)
	}
	if !strings.Contains(lines[0], "ERROR") {
		t.Errorf("expected output to contain ERROR, got: %s", lines[0])
	}
}

// TestIntegration_JSONExclude verifies that JSON-formatted logs are parsed and
// lines matching the exclude pattern are dropped.
func TestIntegration_JSONExclude(t *testing.T) {
	input := []string{
		`{"level":"info","msg":"server started","port":8080}`,
		`{"level":"debug","msg":"health check ok"}`,
		`{"level":"error","msg":"unhandled panic","err":"nil pointer"}`,
		`{"level":"debug","msg":"cache miss","key":"user:42"}`,
	}
	inFile := writeLines(t, input)
	outFile := filepath.Join(t.TempDir(), "out.log")

	cfg := config.Default()
	cfg.InputFile = inFile
	cfg.OutputFile = outFile
	cfg.Exclude = []string{"debug"}
	cfg.NoColor = true

	p, err := pipeline.New(cfg)
	if err != nil {
		t.Fatalf("pipeline.New: %v", err)
	}
	if err := p.Run(); err != nil {
		t.Fatalf("pipeline.Run: %v", err)
	}

	got, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	output := string(got)
	if strings.Contains(strings.ToLower(output), "debug") {
		t.Errorf("expected no debug lines in output, got:\n%s", output)
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 output lines (info + error), got %d:\n%s", len(lines), output)
	}
}

// TestIntegration_StatsAccuracy checks that pipeline stats reflect the correct
// total and filtered line counts after a run.
func TestIntegration_StatsAccuracy(t *testing.T) {
	input := []string{
		"line one alpha",
		"line two beta",
		"line three alpha",
		"line four gamma",
		"line five alpha",
	}
	inFile := writeLines(t, input)
	outFile := filepath.Join(t.TempDir(), "out.log")

	cfg := config.Default()
	cfg.InputFile = inFile
	cfg.OutputFile = outFile
	cfg.Include = []string{"alpha"}
	cfg.NoColor = true

	p, err := pipeline.New(cfg)
	if err != nil {
		t.Fatalf("pipeline.New: %v", err)
	}
	if err := p.Run(); err != nil {
		t.Fatalf("pipeline.Run: %v", err)
	}

	stats := p.Stats()
	if stats.Total != 5 {
		t.Errorf("expected Total=5, got %d", stats.Total)
	}
	if stats.Matched != 3 {
		t.Errorf("expected Matched=3, got %d", stats.Matched)
	}
}
