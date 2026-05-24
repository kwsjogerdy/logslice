package highlight_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/highlight"
)

func TestScan_NoTerms(t *testing.T) {
	h := highlight.New(nil, false)
	r := h.Scan("hello world")
	if len(r.Matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(r.Matches))
	}
}

func TestScan_SingleMatch(t *testing.T) {
	h := highlight.New([]string{"error"}, false)
	r := h.Scan("an error occurred")
	if len(r.Matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(r.Matches))
	}
	if r.Matches[0].Start != 3 || r.Matches[0].End != 8 {
		t.Errorf("unexpected match bounds: %+v", r.Matches[0])
	}
}

func TestScan_MultipleOccurrences(t *testing.T) {
	h := highlight.New([]string{"log"}, false)
	r := h.Scan("log log log")
	if len(r.Matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(r.Matches))
	}
}

func TestScan_CaseInsensitive(t *testing.T) {
	h := highlight.New([]string{"ERROR"}, true)
	r := h.Scan("An error happened")
	if len(r.Matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(r.Matches))
	}
	// Bounds must reference the original line casing.
	got := r.Line[r.Matches[0].Start:r.Matches[0].End]
	if got != "error" {
		t.Errorf("expected 'error', got %q", got)
	}
}

func TestScan_MultipleTerms(t *testing.T) {
	h := highlight.New([]string{"foo", "bar"}, false)
	r := h.Scan("foo and bar")
	if len(r.Matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(r.Matches))
	}
}

func TestAnnotate_NoMatches(t *testing.T) {
	r := highlight.Result{Line: "plain line", Matches: nil}
	out := highlight.Annotate(r, "[", "]")
	if out != "plain line" {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestAnnotate_WithMatches(t *testing.T) {
	h := highlight.New([]string{"err"}, false)
	r := h.Scan("an err here")
	out := highlight.Annotate(r, "[", "]")
	expected := "an [err] here"
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestAnnotate_MultipleMatches(t *testing.T) {
	h := highlight.New([]string{"x"}, false)
	r := h.Scan("x y x")
	out := highlight.Annotate(r, "<", ">")
	expected := "<x> y <x>"
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}
