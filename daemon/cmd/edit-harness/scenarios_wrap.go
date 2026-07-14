package main

import "strings"

// wrapParagraph is one logical line that wraps to multiple visual rows on e-ink.
var wrapParagraph = strings.TrimSpace(strings.Repeat("word ", 40))

// wrapScenarios test visual-line motion (not \n logical lines). Requires Scenario.Width
// and calibrated offsets in wrap_fixtures.go after device run.
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
				{Label: "down one visual line", Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{TextLen: intp(n), Cursor: intp(wrapDownOneCursor)}},
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
			},
		},
		{
			Name:    "wrap-up-from-visual-line-2",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Label: "after down", Expect: &StateExpect{Cursor: intp(wrapDownOneCursor), TextLen: intp(n)}},
				{Keys: []Key{{Name: "ArrowUp"}}},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0), TextLen: intp(n)}},
			},
		},
		{
			Name:    "wrap-shift-down-one-visual",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Label: "shift+down", Keys: []Key{{Name: "ArrowDown", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(wrapDownOneCursor), SelEnd: intp(wrapDownOneCursor), SelLenMin: intp(1), TextLen: intp(n)}},
			},
		},
		{
			Name:    "wrap-shift-down-then-up-shrinks",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Label: "shift+down", Keys: []Key{{Name: "ArrowDown", Shift: true}}},
				{Label: "shift+up shrinks", Keys: []Key{{Name: "ArrowUp", Shift: true}}},
				{Expect: &StateExpect{SelLenMin: intp(1), TextLen: intp(n)}},
			},
		},
		{
			Name:    "wrap-down-last-visual-line",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Label: "down on last visual line", Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(n), SelStart: intp(n), SelEnd: intp(n), TextLen: intp(n)}},
			},
		},
		{
			Name:    "wrap-shift-down-last-to-eof",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Label: "shift+down at wrap end", Keys: []Key{{Name: "ArrowDown", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(n), SelEnd: intp(n), SelLenMin: intp(1), TextLen: intp(n)}},
			},
		},
		{
			Name:    "wrap-mixed-newline-and-wrap",
			Content: "short\n" + strings.Repeat("longword ", 12),
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Label: "into wrapped tail", Expect: &StateExpect{CursorMin: intp(6)}},
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
				{Label: "down keeps goal column", Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(wrapGoalColDownCursor), TextLen: intp(goalColLen)}},
			},
		},
		{
			Name:    "wrap-combo-alt-left-word",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(n), TextLen: intp(n)}},
				{Label: "alt+left word on wrap", Keys: []Key{{Name: "ArrowLeft", Alt: true}}},
				{Expect: &StateExpect{CursorMin: intp(1), CursorMax: intp(n - 1), TextLen: intp(n)}},
			},
		},
		{
			Name:    "wrap-combo-ctrl-bs-line",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
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
				{Label: "shift+left across wrap", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 3},
				{Expect: &StateExpect{SelLenMin: intp(1), CursorMin: intp(wrapDownOneCursor - 3), CursorMax: intp(wrapDownOneCursor), TextLen: intp(n)}},
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
