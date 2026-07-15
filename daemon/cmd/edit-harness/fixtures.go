package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Shared fixtures for the keyboard harness. Most scenarios load fixtureProse
// via sandbox prepare into the harness-only note z-test-keyboard-harness.md
// (z-test- files are filtered from normal vault/user flows). Specialized
// Content remains for wrap geometry, empty-document, goal-column shapes, and
// tall page-scroll bodies.
//
// Editor positions (cursor, selStart, selEnd, textLen) match QML TextEdit /
// QString indexing: one unit per Unicode code point for the BMP characters
// used here (Norwegian æøå, accents). Use editorLen / rune helpers — never
// Go len() byte counts — when asserting against live editor state.

// fixtureProse is the shared dummy document: Norwegian prose, æøå and other
// specials, two bullet lists, and long enough uniform lines for N=1/3/7
// motion in both directions without hitting artificial ends early.
var fixtureProse = buildFixtureProse()

var (
	fixtureProseLen int

	proseHStart     int // first char of horizontal line
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

	prosePara2Start   int
	proseList1Start   int
	proseList2Item2   int
	proseMidDocCaret  int
	proseNearEOFCaret int
)

func buildFixtureProse() string {
	var b strings.Builder
	b.WriteString("Writerdeck harness dummy — ikke i vanlig notatliste\n")
	b.WriteString("Naïve café résumé: æøå på Færøyene, Zürich og São Paulo. Böcker, kjøttkaker, blåbær — «sitat» €£¥.\n")
	b.WriteString("\n")
	b.WriteString("Avsnitt to: den gamle typewriteren klapret mens skribenten noterte München, Köln og Ålesund før kvelden falt på.\n")
	b.WriteString("\n")
	b.WriteString("- første punkt: ta med sykkelen til kontoret før regnet\n")
	b.WriteString("- andre punkt: husk rødvin, ost og brød\n")
	b.WriteString("- tredje punkt: les Kapittel VII før middag\n")
	b.WriteString("alfa bravo charlie delta echo foxtrot golf hotel india juliett\n")
	b.WriteString("- eink: lesemodus skal ikke overscrolle\n")
	b.WriteString("- eink: page-up/down skal klemme ved bunnen\n")
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

	proseMidDocCaret = mustFind("café résumé")
	prosePara2Start = mustFind("Avsnitt to:")
	proseList1Start = mustFind("- første punkt:")
	proseList2Item2 = mustFind("- eink: page-up")

	proseWStart = mustFind("alfa bravo charlie")
	wByteStart := strings.Index(fixtureProse, "alfa bravo charlie")
	wByteEnd := strings.Index(fixtureProse, "\n- eink: lese")
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

	proseVStart = mustFind("Linje 00 æøå café!")
	proseVLen = proseVCount*proseVWidth + (proseVCount - 1) // newlines between
	proseNearEOFCaret = proseVLineStart(proseVCount-1) + proseVWidth/2
}

// proseVLineStart returns the editor position of logical vertical line i.
func proseVLineStart(i int) int {
	return proseVStart + i*(proseVWidth+1)
}

// proseVLineEnd returns the editor position one past the last character of
// vertical line i (Home/End within that line).
func proseVLineEnd(i int) int {
	return proseVLineStart(i) + proseVWidth
}

// editorLen is the TextEdit/QString length of s (BMP code points == runes).
func editorLen(s string) int {
	return utf8.RuneCountInString(s)
}

// editorIndexAtByte maps a UTF-8 byte offset in s to an editor position.
func editorIndexAtByte(s string, byteIndex int) int {
	if byteIndex <= 0 {
		return 0
	}
	if byteIndex >= len(s) {
		return editorLen(s)
	}
	return utf8.RuneCountInString(s[:byteIndex])
}

// trimLastRunes returns s without its last n runes (for Text expects).
func trimLastRunes(s string, n int) string {
	r := []rune(s)
	if n >= len(r) {
		return ""
	}
	return string(r[:len(r)-n])
}

// dropFirstRunes returns s without its first n runes.
func dropFirstRunes(s string, n int) string {
	r := []rune(s)
	if n >= len(r) {
		return ""
	}
	return string(r[n:])
}

// wordStarts / wordEnds are relative rune offsets within s.
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

// wrapParagraph is one logical line that wraps to multiple visual rows.
// Length must stay wrapParagraphLen (see wrap_fixtures.go).
var wrapParagraph = strings.TrimSpace(strings.Repeat("word ", 40))

// fixtureTall returns n filler lines for hardware / read page-scroll tests.
func fixtureTall(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "Page scroll filler æøå line %04d ", i)
		b.WriteString(strings.Repeat("x", 32))
		b.WriteByte('\n')
	}
	return b.String()
}

// fixtureTallRead is multi-screen but finite — reading-mode overscroll clamp.
func fixtureTallRead() string {
	return fixtureTall(80)
}

// Specialized tiny geometries (still real words, not abcdef toys).
const (
	fixtureGoalColDown = "tre\ni\nfemte"
	fixtureShorterDown = "tre\ni"
	fixtureShorterUp   = "i\ntre"
	fixtureTwoLines    = "ost\nost"
	fixtureThreeLines  = "en\nto\ntre"
	fixtureTwoParas    = "para1\n\npara2"
	fixtureHelloWorld  = "hello world"
)
