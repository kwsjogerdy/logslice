package truncate

import (
	"strings"
	"testing"
)

func TestNew_DefaultSuffix(t *testing.T) {
	tr := New(10, "")
	if string(tr.suffix) != DefaultSuffix {
		t.Fatalf("expected default suffix %q, got %q", DefaultSuffix, tr.suffix)
	}
}

func TestEnabled(t *testing.T) {
	if New(0, "").Enabled() {
		t.Fatal("expected Enabled=false for maxBytes=0")
	}
	if !New(100, "").Enabled() {
		t.Fatal("expected Enabled=true for maxBytes=100")
	}
}

func TestApply_ShortLine_Unchanged(t *testing.T) {
	tr := New(100, "...")
	line := []byte("short line")
	got := tr.Apply(line)
	if string(got) != string(line) {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_ExactLimit_Unchanged(t *testing.T) {
	tr := New(5, "...")
	line := []byte("hello")
	got := tr.Apply(line)
	if string(got) != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestApply_LongLine_Truncated(t *testing.T) {
	tr := New(10, "...")
	line := []byte("this is a long line that should be clipped")
	got := tr.Apply(line)
	if len(got) != 10 {
		t.Fatalf("expected length 10, got %d: %q", len(got), got)
	}
	if !strings.HasSuffix(string(got), "...") {
		t.Fatalf("expected suffix '...', got %q", got)
	}
}

func TestApply_Disabled_NeverTruncates(t *testing.T) {
	tr := New(0, "...")
	long := []byte(strings.Repeat("x", 10000))
	got := tr.Apply(long)
	if len(got) != len(long) {
		t.Fatalf("expected no truncation, got length %d", len(got))
	}
}

func TestApply_UTF8_RuneBoundary(t *testing.T) {
	// "日本語" is 9 bytes (3 bytes per rune). Limit=8 should cut to 6 bytes (2 runes).
	tr := New(8, "..")
	line := []byte("日本語")
	got := tr.Apply(line)
	// cutAt = 8-2=6, which is a rune boundary; result = first 6 bytes + ".."
	if len(got) != 8 {
		t.Fatalf("expected 8 bytes, got %d: %q", len(got), got)
	}
	if !strings.HasSuffix(string(got), "..") {
		t.Fatalf("expected suffix '..', got %q", got)
	}
}

func TestApply_SuffixLongerThanMax(t *testing.T) {
	// When suffix >= maxBytes, cutAt collapses to 0; output is just the suffix.
	tr := New(2, "...")
	line := []byte("hello world")
	got := tr.Apply(line)
	if string(got) != "..." {
		t.Fatalf("expected %q, got %q", "...", got)
	}
}
