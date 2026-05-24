package timefilter_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/timefilter"
)

func mustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestEnabled_NoBounds(t *testing.T) {
	f := timefilter.New(time.Time{}, time.Time{})
	if f.Enabled() {
		t.Fatal("expected Enabled() == false when no bounds set")
	}
}

func TestEnabled_WithAfter(t *testing.T) {
	f := timefilter.New(mustParse("2024-01-01T00:00:00Z"), time.Time{})
	if !f.Enabled() {
		t.Fatal("expected Enabled() == true")
	}
}

func TestAllow_ZeroTime_AlwaysPasses(t *testing.T) {
	f := timefilter.New(mustParse("2024-01-01T00:00:00Z"), mustParse("2024-12-31T23:59:59Z"))
	if !f.Allow(time.Time{}) {
		t.Fatal("zero time should always be allowed")
	}
}

func TestAllow_WithinWindow(t *testing.T) {
	f := timefilter.New(mustParse("2024-01-01T00:00:00Z"), mustParse("2024-12-31T23:59:59Z"))
	ts := mustParse("2024-06-15T12:00:00Z")
	if !f.Allow(ts) {
		t.Fatal("timestamp inside window should be allowed")
	}
}

func TestAllow_BeforeAfterBound(t *testing.T) {
	f := timefilter.New(mustParse("2024-06-01T00:00:00Z"), time.Time{})
	ts := mustParse("2024-01-01T00:00:00Z")
	if f.Allow(ts) {
		t.Fatal("timestamp before 'after' bound should be rejected")
	}
}

func TestAllow_AfterBeforeBound(t *testing.T) {
	f := timefilter.New(time.Time{}, mustParse("2024-06-01T00:00:00Z"))
	ts := mustParse("2024-12-01T00:00:00Z")
	if f.Allow(ts) {
		t.Fatal("timestamp after 'before' bound should be rejected")
	}
}

func TestAllow_ExactBoundaryExcluded(t *testing.T) {
	boundary := mustParse("2024-06-01T00:00:00Z")
	f := timefilter.New(time.Time{}, boundary)
	// boundary itself is NOT before boundary, so should be rejected
	if f.Allow(boundary) {
		t.Fatal("timestamp equal to 'before' bound should be rejected")
	}
}

func TestParseTimestamp_RFC3339(t *testing.T) {
	ts := timefilter.ParseTimestamp([]string{"2024-06-15T12:00:00Z"})
	if ts.IsZero() {
		t.Fatal("expected non-zero timestamp")
	}
}

func TestParseTimestamp_SpaceLayout(t *testing.T) {
	ts := timefilter.ParseTimestamp([]string{"2024-06-15 12:00:00"})
	if ts.IsZero() {
		t.Fatal("expected non-zero timestamp for space-separated layout")
	}
}

func TestParseTimestamp_Empty(t *testing.T) {
	ts := timefilter.ParseTimestamp([]string{"", ""})
	if !ts.IsZero() {
		t.Fatal("expected zero time for empty candidates")
	}
}

func TestParseTimestamp_FirstValidWins(t *testing.T) {
	ts := timefilter.ParseTimestamp([]string{"not-a-date", "2024-03-01T00:00:00Z"})
	if ts.IsZero() {
		t.Fatal("should parse second candidate when first is invalid")
	}
}
