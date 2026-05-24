package dedup

import (
	"fmt"
	"testing"
)

func TestNew_ClampsWindowSize(t *testing.T) {
	d := New(0)
	if d.size != 1 {
		t.Fatalf("expected size 1, got %d", d.size)
	}
}

func TestIsDuplicate_NewLineNotDuplicate(t *testing.T) {
	d := New(10)
	if d.IsDuplicate("hello world") {
		t.Fatal("first occurrence should not be a duplicate")
	}
}

func TestIsDuplicate_SameLineIsDuplicate(t *testing.T) {
	d := New(10)
	d.IsDuplicate("hello world")
	if !d.IsDuplicate("hello world") {
		t.Fatal("second occurrence should be a duplicate")
	}
}

func TestIsDuplicate_DifferentLinesNotDuplicate(t *testing.T) {
	d := New(10)
	d.IsDuplicate("line one")
	if d.IsDuplicate("line two") {
		t.Fatal("different line should not be a duplicate")
	}
}

func TestIsDuplicate_WindowEviction(t *testing.T) {
	// window size 2: after filling with unique lines the oldest is evicted
	d := New(2)
	d.IsDuplicate("a") // window: [a]
	d.IsDuplicate("b") // window: [a, b]
	d.IsDuplicate("c") // window: [b, c]  — "a" evicted
	if d.IsDuplicate("a") {
		t.Fatal("'a' should have been evicted from the window")
	}
}

func TestIsDuplicate_WindowStillHoldsRecent(t *testing.T) {
	d := New(3)
	d.IsDuplicate("x")
	d.IsDuplicate("y")
	if !d.IsDuplicate("x") {
		t.Fatal("'x' should still be in window")
	}
}

func TestSeen_CountsUpToWindowSize(t *testing.T) {
	d := New(4)
	for i := 0; i < 6; i++ {
		d.IsDuplicate(fmt.Sprintf("line%d", i))
	}
	if d.Seen() != 4 {
		t.Fatalf("expected Seen()=4, got %d", d.Seen())
	}
}

func TestIsDuplicate_EmptyStringHandled(t *testing.T) {
	d := New(5)
	d.IsDuplicate("")
	if !d.IsDuplicate("") {
		t.Fatal("empty string should be detected as duplicate")
	}
}
