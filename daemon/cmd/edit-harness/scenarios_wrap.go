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
			Name:    "wrap-shift-down-then-up-shrinks",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Label: "extend shift+down x3", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 3},
				{Expect: &StateExpect{SelLenMin: intp(1), TextLen: intp(n)}},
				{Label: "shrink shift+up x1 (N=1)", Keys: []Key{{Name: "ArrowUp", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{SelLenMin: intp(1), TextLen: intp(n)}},
				// Full reverse of visual Shift+Down×3 collapses (old expect kept
				// SelLenMin because buggy EOF jump made 3 ups leave a leftover).
				{Label: "shrink shift+up x2 more (N=3 full reverse)", Keys: []Key{{Name: "ArrowUp", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(wrapDownTwoCursor), SelStart: intp(wrapDownTwoCursor), SelEnd: intp(wrapDownTwoCursor), TextLen: intp(n)}},
			},
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
			Name:    "wrap-down-goal-column",
			Content: goalColContent,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(2)}},
				{Label: "down keeps visual x", Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(wrapGoalColDownCursor), TextLen: intp(goalColLen)}},
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
			Name:    "wrap-combo-ctrl-bs-line",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(n), TextLen: intp(n)}},
				{Label: "ctrl+backspace wrapped line", Keys: []Key{{Name: "Backspace", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(0), Cursor: intp(0), Text: strp("")}},
			},
		},
		{
			Name:    "wrap-shift-left-across-wrap",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
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
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor), TextLen: intp(n)}},
				{Label: "home on visual row start", Keys: []Key{{Name: "Home"}}},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor), TextLen: intp(n)}},
			},
		},
		{
			Name:    "wrap-end-on-visual-line",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor), TextLen: intp(n)}},
				{Label: "end on visual row end", Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{CursorMin: intp(wrapDownOneCursor), CursorMax: intp(n - 1), TextLen: intp(n)}},
			},
		},
	}
}
