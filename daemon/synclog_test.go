package main

import (
	"strings"
	"testing"
)

func TestFormatReasonCounts(t *testing.T) {
	got := formatReasonCounts(map[string]int{"token": 5, "boot": 1})
	if got != "boot, token×5" {
		t.Fatalf("got %q", got)
	}
}

func TestSyncSkipCoalesce(t *testing.T) {
	b := &syncSkipBucket{reasons: map[string]int{}}
	b.note("token")
	if b.count != 1 || !b.firstLogged {
		t.Fatalf("first note: count=%d first=%v", b.count, b.firstLogged)
	}
	b.note("token")
	b.note("boot")
	if b.count != 3 {
		t.Fatalf("count=%d want 3", b.count)
	}
	b.flush()
	if b.count != 0 || b.firstLogged {
		t.Fatalf("after flush: count=%d first=%v", b.count, b.firstLogged)
	}
}

func TestPluralAndJoin(t *testing.T) {
	if plural(1) != "" || plural(2) != "s" {
		t.Fatal("plural")
	}
	details := []string{"pushed a.md", "pulled b.md"}
	if strings.Join(details, "; ") != "pushed a.md; pulled b.md" {
		t.Fatal("join")
	}
}
