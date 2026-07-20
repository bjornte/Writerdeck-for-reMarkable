package main

import (
	"strings"
	"testing"
)

func TestProseAnchorsFitN7(t *testing.T) {
	if proseHMid-7 < proseHStart || proseHMid+7 > proseHEditorEnd {
		t.Fatalf("horizontal mid too tight: mid=%d start=%d end=%d", proseHMid, proseHStart, proseHEditorEnd)
	}
	if proseMidDocCaret < 7 || proseMidDocCaret+7 > fixtureProseLen {
		t.Fatalf("para1 mid too tight: %d", proseMidDocCaret)
	}
	if prosePara2Mid < 7 || prosePara2Mid+7 > fixtureProseLen {
		t.Fatalf("para2 mid too tight: %d", prosePara2Mid)
	}
	if len(proseWordStarts) < 12 {
		t.Fatalf("words=%d want >=12", len(proseWordStarts))
	}
	if proseVCount < 12 {
		t.Fatalf("vcount=%d", proseVCount)
	}
	if !strings.Contains(fixtureProse, "Første avsnitt") || !strings.Contains(fixtureProse, "Andre avsnitt") {
		t.Fatal("missing wrapping paragraphs")
	}
	if !strings.Contains(fixtureProse, "Dobbeltblank seksjon") {
		t.Fatal("missing double-blank section")
	}
	if proseDoubleBlank <= 0 || proseDoubleBlank >= prosePara1Start {
		t.Fatalf("double-blank anchor %d should sit before para1 %d", proseDoubleBlank, prosePara1Start)
	}
	between := fixtureProse[strings.Index(fixtureProse, "notatliste"):strings.Index(fixtureProse, "Dobbeltblank")]
	if !strings.Contains(between, "\n\n\n") {
		t.Fatal("expected two consecutive blank lines before double-blank section")
	}
	if fixtureProseLen < 400 {
		t.Fatalf("prose too short: %d", fixtureProseLen)
	}
}

func TestEditorLenMatchesRuneCount(t *testing.T) {
	s := "æøå café résumé"
	if got, want := editorLen(s), 15; got != want {
		t.Fatalf("editorLen(%q)=%d want %d", s, got, want)
	}
}
