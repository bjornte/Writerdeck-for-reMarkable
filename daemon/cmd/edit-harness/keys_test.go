package main

import (
	"strings"
	"testing"
)

func TestWrapParagraphLenConstant(t *testing.T) {
	if len(wrapParagraph) != wrapParagraphLen {
		t.Fatalf("wrapParagraphLen=%d but len(wrapParagraph)=%d; recalibrate wrap_fixtures.go",
			wrapParagraphLen, len(wrapParagraph))
	}
}

func TestUtf8ByteAtRune(t *testing.T) {
	const s = "test résumé"
	if got := utf8ByteAtRune(s, 4); got != 4 {
		t.Fatalf("byteAtRune(4) = %d want 4", got)
	}
	if got := utf8Len(s); got != 13 {
		t.Fatalf("utf8Len = %d want 13", got)
	}
}

func TestScenarioKeysValid(t *testing.T) {
	for _, sc := range AllScenarios() {
		for i, step := range sc.Steps {
			for j, k := range step.Keys {
				if msg := validateScenarioKey(k); msg != "" {
					t.Fatalf("%s step %d key %d: %s", sc.Name, i, j, msg)
				}
			}
		}
	}
}

func TestScenarioTagsAssigned(t *testing.T) {
	for _, sc := range AllScenarios() {
		if len(sc.Tags) == 0 {
			t.Fatalf("scenario %q has no tags", sc.Name)
		}
	}
}

func TestFindScenariosByTag(t *testing.T) {
	got, ok := findScenariosByTag("wrap")
	if !ok {
		t.Fatal("expected wrap tag")
	}
	for _, sc := range got {
		found := false
		for _, tag := range sc.Tags {
			if tag == "wrap" {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("scenario %q missing wrap tag", sc.Name)
		}
	}
}

func TestStepNeedsModifiedPrime(t *testing.T) {
	if stepNeedsModifiedPrime(Step{Keys: []Key{{Name: "Home", Ctrl: true}}}) {
		t.Fatal("ctrl+home should not prime")
	}
	if !stepNeedsModifiedPrime(Step{Keys: []Key{{Name: "End", Ctrl: true}}}) {
		t.Fatal("ctrl+end from 0 should prime")
	}
	if !stepNeedsModifiedPrime(Step{Keys: []Key{{Name: "ArrowRight", Ctrl: true}}}) {
		t.Fatal("ctrl+right from 0 should prime")
	}
	if stepNeedsModifiedPrime(Step{Keys: []Key{{Name: "ArrowRight", Shift: true}}}) {
		t.Fatal("shift-only arrow should not prime")
	}
}

func TestKeyModifiedNav(t *testing.T) {
	kAlt := Key{Name: "ArrowRight", Alt: true}
	if !kAlt.isModifiedNav() {
		t.Fatal("alt+arrow should be modified nav")
	}
	kPlain := Key{Name: "ArrowRight"}
	if kPlain.isModifiedNav() {
		t.Fatal("plain arrow should not be modified nav")
	}
	kRel := Key{Name: "ArrowLeft", Shift: true, Ctrl: true}
	if !kRel.needsExplicitRelease() {
		t.Fatal("ctrl+shift+arrow needs explicit release")
	}
}

func TestInferScenarioTags(t *testing.T) {
	if tags := inferScenarioTags("wrap-down-one-visual-line"); len(tags) != 1 || tags[0] != "wrap" {
		t.Fatalf("wrap tags = %v", tags)
	}
	if tags := inferScenarioTags("down-one-logical-line"); len(tags) != 1 || tags[0] != "regression" {
		t.Fatalf("regression tags = %v", tags)
	}
}

func TestScenarioContentValidTextExpect(t *testing.T) {
	for _, sc := range AllScenarios() {
		for i, step := range sc.Steps {
			if step.Expect == nil || step.Expect.Text == nil {
				continue
			}
			if step.Expect.TextLen != nil && *step.Expect.TextLen != len(*step.Expect.Text) {
				t.Fatalf("%s step %d: textLen %d != len(text) %d", sc.Name, i, *step.Expect.TextLen, len(*step.Expect.Text))
			}
		}
	}
}

func TestSortedTagNames(t *testing.T) {
	got := sortedTagNames(map[string]bool{"wrap": true, "core": true, "gap": true})
	if strings.Join(got, ",") != "core,gap,wrap" {
		t.Fatalf("sorted tags = %v", got)
	}
}
