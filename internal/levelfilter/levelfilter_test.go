package levelfilter_test

import (
	"testing"

	"github.com/your-org/logslice/internal/levelfilter"
)

func TestNew_BothEmpty_Disabled(t *testing.T) {
	f, err := levelfilter.New("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Enabled() {
		t.Fatal("expected filter to be disabled")
	}
}

func TestNew_InvalidMin_ReturnsError(t *testing.T) {
	_, err := levelfilter.New("nonsense", "")
	if err == nil {
		t.Fatal("expected error for invalid min level")
	}
}

func TestNew_InvalidMax_ReturnsError(t *testing.T) {
	_, err := levelfilter.New("", "bogus")
	if err == nil {
		t.Fatal("expected error for invalid max level")
	}
}

func TestAllow_DisabledFilter_AlwaysTrue(t *testing.T) {
	f, _ := levelfilter.New("", "")
	for _, lvl := range []string{"debug", "info", "warn", "error", "fatal"} {
		if !f.Allow(lvl) {
			t.Errorf("disabled filter should allow %q", lvl)
		}
	}
}

func TestAllow_MinWarn_DropsDebugInfo(t *testing.T) {
	f, err := levelfilter.New("warn", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Allow("debug") {
		t.Error("debug should be filtered out")
	}
	if f.Allow("info") {
		t.Error("info should be filtered out")
	}
	if !f.Allow("warn") {
		t.Error("warn should be allowed")
	}
	if !f.Allow("error") {
		t.Error("error should be allowed")
	}
}

func TestAllow_MaxWarn_DropsErrorFatal(t *testing.T) {
	f, err := levelfilter.New("", "warn")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow("debug") {
		t.Error("debug should be allowed")
	}
	if !f.Allow("warn") {
		t.Error("warn should be allowed")
	}
	if f.Allow("error") {
		t.Error("error should be filtered out")
	}
	if f.Allow("fatal") {
		t.Error("fatal should be filtered out")
	}
}

func TestAllow_UnknownLevel_Passes(t *testing.T) {
	f, _ := levelfilter.New("warn", "error")
	if !f.Allow("UNKNOWN") {
		t.Error("unknown level should pass through")
	}
}

func TestAllow_CaseInsensitive(t *testing.T) {
	f, err := levelfilter.New("warn", "error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow("WARN") {
		t.Error("uppercase WARN should be allowed")
	}
	if !f.Allow("Error") {
		t.Error("mixed-case Error should be allowed")
	}
}
