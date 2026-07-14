package main

import (
	"strings"
	"testing"
)

func TestScenarioNamesUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, sc := range AllScenarios() {
		if sc.Name == "" {
			t.Fatal("scenario with empty name")
		}
		if seen[sc.Name] {
			t.Fatalf("duplicate scenario name %q", sc.Name)
		}
		seen[sc.Name] = true
	}
}

func TestScenarioContentValid(t *testing.T) {
	for _, sc := range AllScenarios() {
		n := len(sc.Content)
		for i, step := range sc.Steps {
			if step.Expect == nil {
				continue
			}
			check := func(field string, v *int) {
				if v == nil {
					return
				}
				if *v < 0 || *v > n {
					t.Fatalf("%s step %d: %s=%d out of range for content len %d", sc.Name, i, field, *v, n)
				}
			}
			check("cursor", step.Expect.Cursor)
			check("selStart", step.Expect.SelStart)
			check("selEnd", step.Expect.SelEnd)
			if step.Expect.TextLen != nil && *step.Expect.TextLen < 0 {
				t.Fatalf("%s step %d: textLen=%d invalid", sc.Name, i, *step.Expect.TextLen)
			}
			if step.Expect.SelLen != nil && *step.Expect.SelLen < 0 {
				t.Fatalf("%s step %d: selLen=%d invalid", sc.Name, i, *step.Expect.SelLen)
			}
		}
	}
}

func TestFindScenarioByPrefix(t *testing.T) {
	got, ok := findScenariosByPrefix("shift-left")
	if !ok {
		t.Fatal("expected prefix match")
	}
	for _, sc := range got {
		if !strings.Contains(sc.Name, "shift-left") {
			t.Fatalf("unexpected scenario %q", sc.Name)
		}
	}
}
