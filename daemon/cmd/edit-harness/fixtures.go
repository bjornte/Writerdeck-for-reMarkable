package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Shared fixtures for the keyboard harness. Most scenarios load fixtureProse
// via sandbox prepare into the harness-only note z-test-keyboard-harness.md
// (z-test- files are filtered from normal vault/user flows). Specialized
// Content remains for wrap geometry (Width=320), empty-document, goal-column
// shapes, and tall page-scroll bodies.
//
// Editor positions match QML TextEdit / QString indexing (BMP == runes here).
// Use editorLen — never Go len() byte counts — against live editor state.
//
// Motion/selection pattern (see pattern.go): uni 1, uni 5, bi 1+1, bi 3+5
// (overshoot), bi 7+7 — each direction pair covered.

var fixtureProse = buildFixtureProse()

var (
	fixtureProseLen int

	proseHStart     int
	proseHLen       int
	proseHMid       int
	proseH          string
	proseHEditorEnd int

	proseVCount = 12
	proseVWidth = 18
	proseVStart int
	proseVLen   int

	proseWStart     int
	proseW          string
	proseWordStarts []int
	proseWordEnds   []int
	proseWEditorEnd int

	prosePara1Start   int // long wrapping paragraph 1
	prosePara2Start   int // long wrapping paragraph 2
	prosePara3Start   int // third body paragraph
	proseList1Start   int
	proseList2Item2   int
	proseMidDocCaret    int // mid paragraph 1
	proseNearEOFCaret   int
	prosePara2Mid       int // caret mid paragraph 2 (different placement)
	prosePara1NearEnd   int // near end of para1 (Shift across blank into para2)
	prosePara2NearStart int // early in para2 (Shift+Up back into para1)
)

// Long prose lines intentionally exceed e-ink / wrap width so each paragraph
// paints across at least three visual lines on device (and many logical lines
// when hard-wrapped in the source for readability).

func buildFixtureProse() string {
	para1 := strings.Join([]string{
		"Første avsnitt — Naïve café résumé med æøå på Færøyene, Zürich og São Paulo.",
		"Böcker, kjøttkaker og blåbær fyller hyllen mens typewriteren klaprer gjennom natten,",
		"og skribenten noterer München, Köln, Ålesund og «sitat» med €£¥ før regnet treffer glasset.",
		"Denne prosaen skal speile ekte lesing: nok setninger til at avsnittet bryter visuelt over minst tre linjer.",
	}, " ")

	para2 := strings.Join([]string{
		"Andre avsnitt fortsetter historien om eink-skrivingen: page-up og page-down skal klemme ved bunnen,",
		"aldri scrolle inn i tomrom under dokumentet. Kapittel VII venter før middag, og sykkelen står klar",
		"til kontoret hvis skyene letter. Rødvin, ost og brød ligger i vesken — små bevis på at testteksten",
		"er ekte prosa, ikke abcdef-leker, med plass til både venstre/høyre-og opp/ned-vandring.",
	}, " ")

	para3 := strings.Join([]string{
		"Tredje avsnitt er kortere men fortsatt flersorglinjet: foxtrot over hotelindia,",
		"juliett over kilo, lima over november — ordrekke for Alt-hopp uten å forlate avsnittet.",
	}, " ")

	var b strings.Builder
	b.WriteString("Writerdeck harness dummy — ikke i vanlig notatliste\n")
	b.WriteString("\n")
	b.WriteString(para1)
	b.WriteString("\n\n")
	b.WriteString(para2)
	b.WriteString("\n\n")
	b.WriteString(para3)
	b.WriteString("\n\n")
	b.WriteString("- første punkt: ta med sykkelen til kontoret før regnet\n")
	b.WriteString("- andre punkt: husk rødvin, ost og brød\n")
	b.WriteString("- tredje punkt: les Kapittel VII før middag\n")
	b.WriteString("\n")
	b.WriteString("alfa bravo charlie delta echo foxtrot golf hotel india juliett kilo lima\n")
	b.WriteString("\n")
	b.WriteString("- eink: lesemodus skal ikke overscrolle\n")
	b.WriteString("- eink: page-up/down skal klemme ved bunnen\n")
	b.WriteString("\n")
	b.WriteString("Horisontal æøå: abcdefghijklmnopqrstuvwxyzæøåcafé!!\n")
	for i := 0; i < proseVCount; i++ {
		line := fmt.Sprintf("Linje %02d æøå café!", i)
		if editorLen(line) != proseVWidth {
			panic(fmt.Sprintf("proseVWidth mismatch line %d: got %d want %d (%q)", i, editorLen(line), proseVWidth, line))
		}
		b.WriteString(line)
		if i < proseVCount-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func init() {
	fixtureProseLen = editorLen(fixtureProse)
	mustFind := func(sub string) int {
		i := strings.Index(fixtureProse, sub)
		if i < 0 {
			panic("fixtureProse missing " + sub)
		}
		return editorIndexAtByte(fixtureProse, i)
	}

	prosePara1Start = mustFind("Første avsnitt")
	prosePara2Start = mustFind("Andre avsnitt")
	prosePara3Start = mustFind("Tredje avsnitt")
	proseMidDocCaret = mustFind("typewriteren klaprer")
	prosePara2Mid = mustFind("page-up og page-down")
	prosePara1NearEnd = mustFind("bryter visuelt")
	prosePara2NearStart = mustFind("historien om eink")
	proseList1Start = mustFind("- første punkt:")
	proseList2Item2 = mustFind("- eink: page-up")

	proseWStart = mustFind("alfa bravo charlie")
	wByteStart := strings.Index(fixtureProse, "alfa bravo charlie")
	wByteEnd := strings.Index(fixtureProse, "\n\n- eink: lese")
	if wByteEnd < 0 {
		wByteEnd = strings.Index(fixtureProse, "\n- eink: lese")
	}
	if wByteEnd < 0 {
		panic("fixtureProse word line end missing")
	}
	proseW = fixtureProse[wByteStart:wByteEnd]
	proseWEditorEnd = proseWStart + editorLen(proseW)
	rels := wordStarts(proseW)
	rele := wordEnds(proseW)
	if len(rels) < 10 || len(rele) < 10 {
		panic(fmt.Sprintf("need >=10 words, got starts=%d ends=%d", len(rels), len(rele)))
	}
	proseWordStarts = make([]int, len(rels))
	proseWordEnds = make([]int, len(rele))
	for j, r := range rels {
		proseWordStarts[j] = proseWStart + r
	}
	for j, r := range rele {
		proseWordEnds[j] = proseWStart + r
	}

	proseHStart = mustFind("Horisontal æøå:")
	hByte := strings.Index(fixtureProse, "Horisontal æøå:")
	hEndRel := strings.Index(fixtureProse[hByte:], "\n")
	if hEndRel < 0 {
		panic("horizontal line missing newline")
	}
	proseH = fixtureProse[hByte : hByte+hEndRel]
	proseHLen = editorLen(proseH)
	proseHEditorEnd = proseHStart + proseHLen
	proseHMid = proseHStart + proseHLen/2
	if proseHMid-7 < proseHStart || proseHMid+7 > proseHEditorEnd {
		panic(fmt.Sprintf("horizontal mid %d cannot host ±7 (start=%d end=%d)", proseHMid, proseHStart, proseHEditorEnd))
	}

	proseVStart = mustFind("Linje 00 æøå café!")
	proseVLen = proseVCount*proseVWidth + (proseVCount - 1)
	proseNearEOFCaret = proseVLineStart(proseVCount-1) + proseVWidth/2
}

func proseVLineStart(i int) int {
	return proseVStart + i*(proseVWidth+1)
}

func proseVLineEnd(i int) int {
	return proseVLineStart(i) + proseVWidth
}

func editorLen(s string) int {
	return utf8.RuneCountInString(s)
}

func editorIndexAtByte(s string, byteIndex int) int {
	if byteIndex <= 0 {
		return 0
	}
	if byteIndex >= len(s) {
		return editorLen(s)
	}
	return utf8.RuneCountInString(s[:byteIndex])
}

func trimLastRunes(s string, n int) string {
	r := []rune(s)
	if n >= len(r) {
		return ""
	}
	return string(r[:len(r)-n])
}

func dropFirstRunes(s string, n int) string {
	r := []rune(s)
	if n >= len(r) {
		return ""
	}
	return string(r[n:])
}

func wordStarts(s string) []int {
	var out []int
	inWord := false
	i := 0
	for _, r := range s {
		if r != ' ' && r != '\n' {
			if !inWord {
				out = append(out, i)
				inWord = true
			}
		} else {
			inWord = false
		}
		i++
	}
	return out
}

func wordEnds(s string) []int {
	var out []int
	inWord := false
	i := 0
	for _, r := range s {
		if r != ' ' && r != '\n' {
			inWord = true
		} else if inWord {
			out = append(out, i)
			inWord = false
		}
		i++
	}
	if inWord {
		out = append(out, i)
	}
	return out
}

var wrapParagraph = strings.TrimSpace(strings.Repeat("word ", 40))

func fixtureTall(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "Page scroll filler æøå line %04d ", i)
		b.WriteString(strings.Repeat("x", 32))
		b.WriteByte('\n')
	}
	return b.String()
}

func fixtureTallRead() string {
	return fixtureTall(80)
}

const (
	fixtureGoalColDown = "tre\ni\nfemte"
	fixtureShorterDown = "tre\ni"
	fixtureShorterUp   = "i\ntre"
	fixtureTwoLines    = "ost\nost"
	fixtureThreeLines  = "en\nto\ntre"
	fixtureTwoParas    = "para1\n\npara2"
	fixtureHelloWorld  = "hello world"
)
