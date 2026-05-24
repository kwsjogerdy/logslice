package redact_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/redact"
)

func TestNew_InvalidPattern(t *testing.T) {
	_, err := redact.New([]string{"[invalid"}, "")
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestEnabled_NoPatterns(t *testing.T) {
	r, _ := redact.New(nil, "")
	if r.Enabled() {
		t.Fatal("expected Enabled() == false with no patterns")
	}
}

func TestEnabled_WithPatterns(t *testing.T) {
	r, _ := redact.New([]string{`\d+`}, "")
	if !r.Enabled() {
		t.Fatal("expected Enabled() == true with patterns")
	}
}

func TestApply_NoPatterns_Unchanged(t *testing.T) {
	r, _ := redact.New(nil, "")
	got := r.Apply("hello world")
	if got != "hello world" {
		t.Fatalf("expected unchanged, got %q", got)
	}
}

func TestApply_RedactsMatch(t *testing.T) {
	r, _ := redact.New([]string{`\b\d{4}-\d{4}-\d{4}-\d{4}\b`}, "")
	input := "card: 1234-5678-9012-3456 processed"
	got := r.Apply(input)
	want := "card: [REDACTED] processed"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	r, _ := redact.New([]string{`password=\S+`}, "***")
	got := r.Apply("password=secret123 user=alice")
	want := "*** user=alice"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApply_MultiplePatterns(t *testing.T) {
	r, _ := redact.New([]string{`\b\d{3}-\d{2}-\d{4}\b`, `token=\S+`}, "")
	got := r.Apply("ssn=123-45-6789 token=abc123")
	want := "ssn=[REDACTED] token=[REDACTED]"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApplyBytes_RedactsMatch(t *testing.T) {
	r, _ := redact.New([]string{`secret`}, "")
	got := r.ApplyBytes([]byte("my secret value"))
	want := "my [REDACTED] value"
	if string(got) != want {
		t.Fatalf("expected %q, got %q", want, string(got))
	}
}
