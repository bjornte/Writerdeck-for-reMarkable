package main

import "testing"

func TestProseAnchorsFitN7(t *testing.T) {
	if proseHMid-7 < proseHStart {
		t.Fatalf("horizontal mid too close to start: mid=%d start=%d", proseHMid, proseHStart)
	}
	if proseHMid+7 > proseHEditorEnd {
		t.Fatalf("horizontal mid too close to end: mid=%d end=%d", proseHMid, proseHEditorEnd)
	}
	if len(proseWordStarts) < 10 {
		t.Fatalf("words=%d", len(proseWordStarts))
	}
	if proseVCount < 12 {
		t.Fatalf("vcount=%d", proseVCount)
	}
	if fixtureProseLen < 200 {
		t.Fatalf("prose too short: %d", fixtureProseLen)
	}
}

func TestEditorLenMatchesRuneCount(t *testing.T) {
	s := "æøå café résumé"
	if got, want := editorLen(s), 15; got != want {
		t.Fatalf("editorLen(%q)=%d want %d", s, got, want)
	}
}
