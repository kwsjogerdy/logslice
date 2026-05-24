package fieldfilter_test

import (
	"testing"

	"github.com/your-org/logslice/internal/fieldfilter"
)

func TestEnabled_NoFields(t *testing.T) {
	ff := fieldfilter.New(nil, nil)
	if ff.Enabled() {
		t.Fatal("expected Enabled() == false when no fields configured")
	}
}

func TestEnabled_WithInclude(t *testing.T) {
	ff := fieldfilter.New([]string{"level"}, nil)
	if !ff.Enabled() {
		t.Fatal("expected Enabled() == true with include fields")
	}
}

func TestApply_NoRules_ReturnsSameMap(t *testing.T) {
	ff := fieldfilter.New(nil, nil)
	input := map[string]string{"level": "info", "msg": "hello"}
	out := ff.Apply(input)
	if len(out) != len(input) {
		t.Fatalf("expected %d fields, got %d", len(input), len(out))
	}
}

func TestApply_IncludeFields(t *testing.T) {
	ff := fieldfilter.New([]string{"level", "msg"}, nil)
	input := map[string]string{"level": "info", "msg": "hello", "ts": "123", "caller": "main.go"}
	out := ff.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out))
	}
	if out["level"] != "info" || out["msg"] != "hello" {
		t.Fatal("unexpected field values after include filter")
	}
}

func TestApply_ExcludeFields(t *testing.T) {
	ff := fieldfilter.New(nil, []string{"ts", "caller"})
	input := map[string]string{"level": "info", "msg": "hello", "ts": "123", "caller": "main.go"}
	out := ff.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out))
	}
	if _, ok := out["ts"]; ok {
		t.Fatal("ts should have been excluded")
	}
}

func TestApply_IncludeAndExclude(t *testing.T) {
	// include wins first, then exclude removes from that set
	ff := fieldfilter.New([]string{"level", "msg", "ts"}, []string{"ts"})
	input := map[string]string{"level": "warn", "msg": "hi", "ts": "999", "extra": "drop"}
	out := ff.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out))
	}
	if _, ok := out["ts"]; ok {
		t.Fatal("ts should have been excluded")
	}
}

func TestApply_CaseInsensitiveKeys(t *testing.T) {
	ff := fieldfilter.New([]string{"Level"}, nil)
	input := map[string]string{"level": "debug", "msg": "test"}
	out := ff.Apply(input)
	if _, ok := out["level"]; !ok {
		t.Fatal("expected 'level' key to be retained via case-insensitive include")
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 field, got %d", len(out))
	}
}

func TestApply_OriginalMapUnmodified(t *testing.T) {
	ff := fieldfilter.New(nil, []string{"secret"})
	input := map[string]string{"msg": "hi", "secret": "token"}
	_ = ff.Apply(input)
	if _, ok := input["secret"]; !ok {
		t.Fatal("original map should not be modified")
	}
}
