package formatter_test

import (
	"testing"

	"github.com/user/logslice/internal/formatter"
)

func TestParseLevel(t *testing.T) {
	cases := []struct {
		input string
		want  formatter.Level
	}{
		{"debug", formatter.LevelDebug},
		{"DBG", formatter.LevelDebug},
		{"trace", formatter.LevelDebug},
		{"info", formatter.LevelInfo},
		{"INFORMATION", formatter.LevelInfo},
		{"warn", formatter.LevelWarn},
		{"WARNING", formatter.LevelWarn},
		{"error", formatter.LevelError},
		{"ERR", formatter.LevelError},
		{"fatal", formatter.LevelFatal},
		{"CRITICAL", formatter.LevelFatal},
		{"panic", formatter.LevelFatal},
		{"unknown", formatter.LevelUnknown},
		{"", formatter.LevelUnknown},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := formatter.ParseLevel(tc.input)
			if got != tc.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestFormatter_NoColor(t *testing.T) {
	f := formatter.New(false)
	line := "2024-01-01 INFO hello world"
	got := f.Format(line, "info")
	if got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestFormatter_ColorizeReturnsNonEmpty(t *testing.T) {
	f := formatter.New(true)
	levels := []string{"debug", "info", "warn", "error", "fatal", ""}
	line := "some log line"
	for _, lvl := range levels {
		got := f.Format(line, lvl)
		if got == "" {
			t.Errorf("Format(%q, %q) returned empty string", line, lvl)
		}
	}
}

func TestFormatter_UnknownLevelNoColor(t *testing.T) {
	f := formatter.New(true)
	line := "plain log line"
	got := f.Format(line, "")
	if got != line {
		t.Errorf("expected unchanged line for unknown level, got %q", got)
	}
}
