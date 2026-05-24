package ratelimit

import (
	"testing"
	"time"
)

func TestNew_ZeroRate_Unlimited(t *testing.T) {
	l := New(0)
	for i := 0; i < 10000; i++ {
		if !l.Allow() {
			t.Fatal("expected unlimited limiter to always allow")
		}
	}
}

func TestNew_NegativeRate_Unlimited(t *testing.T) {
	l := New(-5)
	if l.Rate() != 0 {
		t.Fatalf("expected rate 0, got %d", l.Rate())
	}
	if !l.Allow() {
		t.Fatal("expected unlimited limiter to allow")
	}
}

func TestAllow_ExhaustsBucket(t *testing.T) {
	l := New(3)
	// First 3 calls should succeed.
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("call %d should be allowed", i+1)
		}
	}
	// 4th call should be dropped.
	if l.Allow() {
		t.Fatal("4th call should be dropped")
	}
}

func TestAllow_RefillsAfterOneSec(t *testing.T) {
	base := time.Now()
	l := New(2)
	l.clock = func() time.Time { return base }

	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("should be dropped before refill")
	}

	// Advance clock by 1 second.
	l.clock = func() time.Time { return base.Add(time.Second) }
	if !l.Allow() {
		t.Fatal("should be allowed after refill")
	}
}

func TestAllow_MultiPeriodRefill(t *testing.T) {
	base := time.Now()
	l := New(5)
	l.clock = func() time.Time { return base }

	// Drain bucket.
	for i := 0; i < 5; i++ {
		l.Allow()
	}

	// Advance 3 seconds — bucket should cap at rate (5), not 15.
	l.clock = func() time.Time { return base.Add(3 * time.Second) }
	allowed := 0
	for i := 0; i < 10; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 5 {
		t.Fatalf("expected 5 allowed after 3s refill, got %d", allowed)
	}
}

func TestRate_ReturnsConfiguredRate(t *testing.T) {
	l := New(42)
	if l.Rate() != 42 {
		t.Fatalf("expected rate 42, got %d", l.Rate())
	}
}
