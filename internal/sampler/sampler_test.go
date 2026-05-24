package sampler

import (
	"testing"
)

func TestNew_ClampsRate(t *testing.T) {
	s := New(-0.5, false)
	if s.Rate() != 0 {
		t.Fatalf("expected 0, got %f", s.Rate())
	}
	s2 := New(1.5, false)
	if s2.Rate() != 1.0 {
		t.Fatalf("expected 1.0, got %f", s2.Rate())
	}
}

func TestKeep_RateOne_AlwaysKeeps(t *testing.T) {
	s := New(1.0, true)
	for i := 0; i < 1000; i++ {
		if !s.Keep() {
			t.Fatal("expected Keep to return true for rate=1.0")
		}
	}
}

func TestKeep_RateZero_NeverKeeps(t *testing.T) {
	s := New(0.0, true)
	for i := 0; i < 1000; i++ {
		if s.Keep() {
			t.Fatal("expected Keep to return false for rate=0.0")
		}
	}
}

func TestKeep_RateHalf_Approximate(t *testing.T) {
	s := New(0.5, true)
	kept := 0
	const n = 10000
	for i := 0; i < n; i++ {
		if s.Keep() {
			kept++
		}
	}
	ratio := float64(kept) / float64(n)
	if ratio < 0.40 || ratio > 0.60 {
		t.Fatalf("expected ratio near 0.5, got %f", ratio)
	}
}

func TestSeen_CountsCallsToKeep(t *testing.T) {
	s := New(0.5, true)
	for i := 0; i < 20; i++ {
		s.Keep()
	}
	if s.Seen() != 20 {
		t.Fatalf("expected 20, got %d", s.Seen())
	}
}

func TestSeen_RateOne_DoesNotCountWhenShortCircuited(t *testing.T) {
	// rate=1.0 short-circuits before incrementing counter
	s := New(1.0, true)
	for i := 0; i < 5; i++ {
		s.Keep()
	}
	if s.Seen() != 0 {
		t.Fatalf("expected 0, got %d", s.Seen())
	}
}
