package main

import "strings"

// wrapScenarios test visual-line motion (not \n logical lines). Requires Scenario.Width
// and calibrated offsets in wrap_fixtures.go. Wrap geometry stays specialized
// (word×40 at W=320); motion/selection cases prove N=1/3/7 both directions.
func wrapScenarios() []Scenario {
	wp := wrapParagraph
	n := wrapParagraphLen
	if len(wp) != n {
		panic("wrapParagraphLen mismatch: update wrap_fixtures.go")
	}
	goalColContent := "ab" + strings.Repeat("word ", 35)
	goalColLen := len(goalColContent)
	return []Scenario{
		{
			Name:    "wrap-down-one-visual-line",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Label: "down x1 (N=1)", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 1},
				{Expect: &StateExpect{TextLen: intp(n), Cursor: intp(wrapDownOneCursor)}},
				{Label: "down x2 more (N=3)", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Expect: &StateExpect{TextLen: intp(n), Cursor: intp(wrapDownThreeCursor)}},
				{Label: "down x4 more (N=7)", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 4},
				{Expect: &StateExpect{TextLen: intp(n), Cursor: intp(wrapDownSevenCursor)}},
			},
		},
		{
			Name:    "wrap-down-not-jump-paragraph",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{TextLen: intp(n), Cursor: intp(wrapDownOneCursor)}},
				{Label: "still inside wrapped paragraph after N=3", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Expect: &StateExpect{TextLen: intp(n), CursorMin: intp(wrapDownOneCursor), CursorMax: intp(n - 1)}},
			},
		},
		{
			Name:    "wrap-up-from-visual-line-2",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}, Repeat: 7},
				{Expect: &StateExpect{Cursor: intp(wrapDownSevenCursor), TextLen: intp(n)}},
				{Label: "up x1 (N=1)", Keys: []Key{{Name: "ArrowUp"}}, Repeat: 1},
				{Expect: &StateExpect{CursorMin: intp(wrapDownSixCursor - 5), CursorMax: intp(wrapDownSevenCursor), TextLen: intp(n)}},
				{Label: "up x2 more (N=3)", Keys: []Key{{Name: "ArrowUp"}}, Repeat: 2},
				{Expect: &StateExpect{CursorMin: intp(wrapDownFourCursor - 5), CursorMax: intp(wrapDownFourCursor + 5), TextLen: intp(n)}},
				{Label: "up x4 more (N=7)", Keys: []Key{{Name: "ArrowUp"}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0), TextLen: intp(n)}},
			},
		},
		{
			Name:    "wrap-shift-down-one-visual",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Label: "shift+down x1 (N=1)", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor), SelEnd: intp(wrapDownOneCursor), SelLenMin: intp(1), TextLen: intp(n)}},
				{Label: "shift+down x2 more (N=3)", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(wrapDownThreeCursor), SelEnd: intp(wrapDownThreeCursor), SelLenMin: intp(wrapDownOneCursor), TextLen: intp(n)}},
				{Label: "shift+down x4 more (N=7)", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(wrapDownSevenCursor), SelEnd: intp(wrapDownSevenCursor), SelLenMin: intp(wrapDownThreeCursor), TextLen: intp(n)}},
			},
		},
		{
			// Mid-sentence on the wrapped block (not Ctrl+Home / visual col 0).
			Name:    "wrap-shift-down-then-up-shrinks",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: shiftVisualProsePatternStep(wrapDownTwoCursor+5,
				keyArrow("ArrowDown", true, false, false),
				keyArrow("ArrowUp", true, false, false), false, 20, wrapParagraphLen),
		},
		{
			Name:    "wrap-down-last-visual-line",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Label: "down clamp x1 (N=1)", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(n), SelStart: intp(n), SelEnd: intp(n), TextLen: intp(n)}},
				{Label: "down clamp x3 total", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(n), SelStart: intp(n), SelEnd: intp(n), TextLen: intp(n)}},
				{Label: "down clamp x7 total", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(n), SelStart: intp(n), SelEnd: intp(n), TextLen: intp(n)}},
			},
		},
		{
			Name:    "wrap-shift-down-last-to-eof",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Label: "shift+down at wrap end", Keys: []Key{{Name: "ArrowDown", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(n), SelEnd: intp(n), SelLenMin: intp(1), TextLen: intp(n)}},
			},
		},
		{
			Name:    "wrap-mixed-newline-and-wrap",
			Content: "kort\n" + strings.Repeat("langord ", 12),
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Label: "into wrapped tail", Expect: &StateExpect{CursorMin: intp(5)}},
			},
		},
		{
			// Kill-test: Down from col 2 must not snap to col 0 of the next
			// visual row. Up must restore col 2 — index-only Down can false-green
			// if the landing index is merely "nearby."
			Name:    "wrap-down-goal-column",
			Content: goalColContent,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(2)}},
				{Label: "down keeps visual x", Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(wrapGoalColDownCursor), TextLen: intp(goalColLen)}},
				{Label: "up restores goal column", Keys: []Key{{Name: "ArrowUp"}}},
				{Expect: &StateExpect{Cursor: intp(2), TextLen: intp(goalColLen)}},
			},
		},
		{
			Name:    "wrap-combo-alt-left-word",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(n), TextLen: intp(n)}},
				{Label: "alt+left x1 (N=1)", Keys: []Key{{Name: "ArrowLeft", Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{CursorMin: intp(1), CursorMax: intp(n - 1), TextLen: intp(n)}},
				{Label: "alt+left x2 more (N=3)", Keys: []Key{{Name: "ArrowLeft", Alt: true}}, Repeat: 2},
				{Expect: &StateExpect{CursorMin: intp(1), CursorMax: intp(n - 1), TextLen: intp(n)}},
				{Label: "alt+left x4 more (N=7)", Keys: []Key{{Name: "ArrowLeft", Alt: true}}, Repeat: 4},
				{Expect: &StateExpect{CursorMin: intp(0), CursorMax: intp(n - 1), TextLen: intp(n)}},
			},
		},
		{
			Name:    "wrap-combo-alt-right-word",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Label: "alt+right x1 (N=1)", Keys: []Key{{Name: "ArrowRight", Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{CursorMin: intp(1), CursorMax: intp(n - 1), TextLen: intp(n)}},
				{Label: "alt+right x2 more (N=3)", Keys: []Key{{Name: "ArrowRight", Alt: true}}, Repeat: 2},
				{Expect: &StateExpect{CursorMin: intp(1), CursorMax: intp(n - 1), TextLen: intp(n)}},
				{Label: "alt+right x4 more (N=7)", Keys: []Key{{Name: "ArrowRight", Alt: true}}, Repeat: 4},
				{Expect: &StateExpect{CursorMin: intp(1), CursorMax: intp(n), TextLen: intp(n)}},
			},
		},
		{
			// Apple ⌘⌫ / CM deleteLineBoundaryBackward: delete to visual-line start only.
			Name:    "wrap-combo-ctrl-bs-line",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{keyMoveDocStart(false)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 5},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor + 5), TextLen: intp(n)}},
				{Label: "deleteToVisualLineStart", Keys: []Key{keyDeleteToVisualLineStart()}},
				{Expect: &StateExpect{
					Cursor:  intp(wrapDownOneCursor),
					TextLen: intp(n - 5),
				}},
			},
		},
		{
			Name:    "wrap-shift-left-across-wrap",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{keyMoveDocStart(false)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor)}},
				{Label: "shift+left x1 (N=1)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{SelLenMin: intp(1), TextLen: intp(n)}},
				{Label: "shift+left x2 more (N=3)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{SelLenMin: intp(1), TextLen: intp(n)}},
				{Label: "shift+left x4 more (N=7)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{SelLenMin: intp(1), CursorMin: intp(0), CursorMax: intp(wrapDownOneCursor), TextLen: intp(n)}},
			},
		},
		{
			Name:    "wrap-home-on-visual-line",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{keyMoveDocStart(false)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor), TextLen: intp(n)}},
				{Label: "homeVisualLine", Keys: []Key{keyHomeVisualLine(false)}},
				{Expect: expectVisualLineCaret(wrapDownOneCursor, n)},
			},
		},
		{
			Name:    "wrap-end-on-visual-line",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{keyMoveDocStart(false)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor), TextLen: intp(n)}},
				{Label: "capture caretY on visual row 1", CaptureCaretY: true},
				{Label: "endVisualLine (lineboundary)", Keys: []Key{keyEndVisualLine(false)}},
				{Expect: &StateExpect{
					Cursor:   intp(wrapEndVisualRow1),
					SelStart: intp(wrapEndVisualRow1),
					SelEnd:   intp(wrapEndVisualRow1),
					TextLen:  intp(n),
					Assoc:    intp(-1),
				}},
				{Label: "caretY stays on visual row 1", ExpectCaretYEqCaptured: true},
			},
		},
		{
			// Apple ⌘←: visual-line start; further presses stay (not paragraph/doc start).
			Name:    "wrap-ctrl-left",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{keyMoveDocStart(false)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 5},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor + 5), TextLen: intp(n)}},
				{Label: "moveToVisualLineStart (N=1)", Keys: []Key{keyMoveToVisualLineStart(false)}, Repeat: 1},
				{Expect: expectVisualLineCaret(wrapDownOneCursor, n)},
				{Label: "moveToVisualLineStart stays (N=3)", Keys: []Key{keyMoveToVisualLineStart(false)}, Repeat: 2},
				{Expect: expectVisualLineCaret(wrapDownOneCursor, n)},
				{Label: "moveToVisualLineStart stays (N=7)", Keys: []Key{keyMoveToVisualLineStart(false)}, Repeat: 4},
				{Expect: expectVisualLineCaret(wrapDownOneCursor, n)},
			},
		},
		{
			// Apple ⌘→: visual-line end; further presses stay (not paragraph end).
			Name:    "wrap-ctrl-right",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{keyMoveDocStart(false)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 5},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor + 5), TextLen: intp(n)}},
				{Label: "capture caretY on visual row 1", CaptureCaretY: true},
				{Label: "moveToVisualLineEnd (N=1)", Keys: []Key{keyMoveToVisualLineEnd(false)}, Repeat: 1},
				{Expect: &StateExpect{
					Cursor:  intp(wrapEndVisualRow1),
					SelStart: intp(wrapEndVisualRow1),
					SelEnd:   intp(wrapEndVisualRow1),
					TextLen:  intp(n),
					Assoc:    intp(-1),
				}},
				// Same painted row as before Ctrl+Right — not the next visual line.
				{Label: "caretY stays on visual row 1", ExpectCaretYEqCaptured: true},
				{Label: "moveToVisualLineEnd stays (N=3)", Keys: []Key{keyMoveToVisualLineEnd(false)}, Repeat: 2},
				{Expect: expectVisualLineCaret(wrapEndVisualRow1, n)},
				{Label: "moveToVisualLineEnd stays (N=7)", Keys: []Key{keyMoveToVisualLineEnd(false)}, Repeat: 4},
				{Expect: expectVisualLineCaret(wrapEndVisualRow1, n)},
			},
		},
		{
			// False-green catch: shared wrap index looks like next-row start unless
			// affinity sticks. Ctrl+Left must return to *this* visual row's start.
			Name:    "wrap-ctrl-right-then-left",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{keyMoveDocStart(false)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 5},
				{Keys: []Key{keyMoveToVisualLineEnd(false)}},
				{Expect: &StateExpect{Cursor: intp(wrapEndVisualRow1), Assoc: intp(-1), TextLen: intp(n)}},
				{Label: "moveToVisualLineStart of same row", Keys: []Key{keyMoveToVisualLineStart(false)}},
				{Expect: expectVisualLineCaret(wrapDownOneCursor, n)},
			},
		},
		{
			// After End on wrap, Up must climb to the previous visual row — not
			// treat the wrap point as already being on the next row.
			Name:    "wrap-end-then-up",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{keyMoveDocStart(false)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor), TextLen: intp(n)}},
				{Label: "endVisualLine", Keys: []Key{keyEndVisualLine(false)}},
				{Expect: &StateExpect{Cursor: intp(wrapEndVisualRow1), Assoc: intp(-1), TextLen: intp(n)}},
				{Label: "up to previous visual row", Keys: []Key{{Name: "ArrowUp"}}},
				{Expect: &StateExpect{
					CursorMin: intp(0),
					CursorMax: intp(wrapDownOneCursor - 1),
					TextLen:   intp(n),
				}},
			},
		},
		{
			// Kill-test: End at wrap point (shared index) then Down must advance
			// one visual row to that row's end — not no-op (col 0 of "next") or
			// skip a row (treating wrap point as already on the next row).
			Name:    "wrap-end-then-down",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{keyMoveDocStart(false)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor), TextLen: intp(n)}},
				{Label: "capture caretY on visual row 1", CaptureCaretY: true},
				{Label: "endVisualLine", Keys: []Key{keyEndVisualLine(false)}},
				{Expect: &StateExpect{Cursor: intp(wrapEndVisualRow1), Assoc: intp(-1), TextLen: intp(n)}},
				{Label: "caretY stays on visual row 1", ExpectCaretYEqCaptured: true},
				{Label: "down to next visual row end", Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{
					Cursor:   intp(wrapEndVisualRow2),
					SelStart: intp(wrapEndVisualRow2),
					SelEnd:   intp(wrapEndVisualRow2),
					TextLen:  intp(n),
				}},
			},
		},
		{
			Name:    "wrap-shift-ctrl-left",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{keyMoveDocStart(false)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 5},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor + 5), TextLen: intp(n)}},
				{Label: "selectToVisualLineStart", Keys: []Key{keyMoveToVisualLineStart(true)}},
				{Expect: &StateExpect{
					// Qt select() parks cursorPosition at the high end; shiftHead is low.
					Cursor:   intp(wrapDownOneCursor + 5),
					SelStart: intp(wrapDownOneCursor),
					SelEnd:   intp(wrapDownOneCursor + 5),
					SelLen:   intp(5),
					TextLen:  intp(n),
				}},
			},
		},
		{
			Name:    "wrap-shift-ctrl-right",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{keyMoveDocStart(false)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 5},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor + 5), TextLen: intp(n)}},
				{Label: "selectToVisualLineEnd", Keys: []Key{keyMoveToVisualLineEnd(true)}},
				{Expect: &StateExpect{
					Cursor:   intp(wrapEndVisualRow1),
					SelStart: intp(wrapDownOneCursor + 5),
					SelEnd:   intp(wrapEndVisualRow1),
					SelLen:   intp(wrapEndVisualRow1 - (wrapDownOneCursor + 5)),
					TextLen:  intp(n),
				}},
			},
		},
	}
}
