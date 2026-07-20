package main

import "testing"

func TestLineBoundaryScenariosHaveWidth(t *testing.T) {
	for _, sc := range AllScenarios() {
		if !isLineBoundaryScenario(sc.Name) {
			continue
		}
		if sc.Width <= 0 {
			t.Fatalf("%s claims lineboundary but Width is unset", sc.Name)
		}
	}
}

func TestLineBoundaryExpectsRejectParagraphEnd(t *testing.T) {
	for _, sc := range AllScenarios() {
		if !isLineBoundaryScenario(sc.Name) {
			continue
		}
		n := editorLen(sc.Content)
		for i, step := range sc.Steps {
			if !expectAllowsParagraphEnd(step.Expect, n) {
				continue
			}
			t.Fatalf("%s step %d (%q): expect allows paragraph/doc end (docLen=%d); use a wrap-point cursor",
				sc.Name, i, step.Label, n)
		}
	}
}

func TestNamedMotionKeys(t *testing.T) {
	if k := keyMoveToVisualLineEnd(false); !k.Ctrl || k.Name != "ArrowRight" || k.Alt {
		t.Fatalf("keyMoveToVisualLineEnd = %+v", k)
	}
	if k := keyMoveWordLeft(true); !k.Alt || !k.Shift || k.Name != "ArrowLeft" {
		t.Fatalf("keyMoveWordLeft(shift) = %+v", k)
	}
	if k := keyDeleteToVisualLineStart(); !k.Ctrl || k.Name != "Backspace" {
		t.Fatalf("keyDeleteToVisualLineStart = %+v", k)
	}
}
