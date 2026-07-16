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
		n := editorLen(sc.Content)
		for i, step := range sc.Steps {
			if step.Expect == nil {
				continue
			}
			maxPos := n
			if step.Expect.TextLen != nil && *step.Expect.TextLen > maxPos {
				maxPos = *step.Expect.TextLen
			}
			check := func(field string, v *int) {
				if v == nil {
					return
				}
				if *v < 0 || *v > maxPos {
					t.Fatalf("%s step %d: %s=%d out of range for content len %d (max %d)", sc.Name, i, field, *v, n, maxPos)
				}
			}
			check("cursor", step.Expect.Cursor)
			check("cursorMin", step.Expect.CursorMin)
			check("cursorMax", step.Expect.CursorMax)
			check("selStart", step.Expect.SelStart)
			check("selEnd", step.Expect.SelEnd)
			if step.Expect.CursorMin != nil && step.Expect.CursorMax != nil && *step.Expect.CursorMin > *step.Expect.CursorMax {
				t.Fatalf("%s step %d: cursorMin > cursorMax", sc.Name, i)
			}
			if step.Expect.TextLen != nil && *step.Expect.TextLen < 0 {
				t.Fatalf("%s step %d: textLen=%d invalid", sc.Name, i, *step.Expect.TextLen)
			}
			if step.Expect.SelLen != nil && *step.Expect.SelLen < 0 {
				t.Fatalf("%s step %d: selLen=%d invalid", sc.Name, i, *step.Expect.SelLen)
			}
			if step.Expect.SelLenMin != nil && *step.Expect.SelLenMin < 0 {
				t.Fatalf("%s step %d: selLenMin=%d invalid", sc.Name, i, *step.Expect.SelLenMin)
			}
		}
	}
}

func TestScenarioCount(t *testing.T) {
	const want = 105
	if n := len(AllScenarios()); n != want {
		t.Fatalf("expected %d scenarios, got %d (update want after adding scenarios)", want, n)
	}
}

func TestWrapScenariosHaveWidth(t *testing.T) {
	for _, sc := range wrapScenarios() {
		if sc.Width <= 0 {
			t.Fatalf("wrap scenario %q missing Width", sc.Name)
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
