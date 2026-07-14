package main

import "strings"

// wrapParagraph is one logical line that wraps to multiple visual rows on e-ink.
var wrapParagraph = strings.TrimSpace(strings.Repeat("word ", 40))

// wrapScenarios test visual-line motion (not \n logical lines). Requires Scenario.Width
// and calibrated offsets in wrap_fixtures.go after device run.
func wrapScenarios() []Scenario {
	wp := wrapParagraph
	n := len(wp)
	return []Scenario{
		{
			Name:    "wrap-down-one-visual-line",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Label: "down one visual line", Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{TextLen: intp(n), CursorMin: intp(1), CursorMax: intp(n - 1)}},
			},
		},
		{
			Name:    "wrap-down-not-jump-paragraph",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{TextLen: intp(n), CursorMin: intp(1), CursorMax: intp(n - 1)}},
			},
		},
		{
			Name:    "wrap-up-from-visual-line-2",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Label: "after down", Expect: &StateExpect{CursorMin: intp(1), CursorMax: intp(n - 1)}},
				{Keys: []Key{{Name: "ArrowUp"}}},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "wrap-shift-down-one-visual",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Label: "shift+down", Keys: []Key{{Name: "ArrowDown", Shift: true}}},
				{Expect: &StateExpect{SelLenMin: intp(1), CursorMin: intp(1), CursorMax: intp(n)}},
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
				{Expect: &StateExpect{SelLenMin: intp(1)}},
			},
		},
		{
			Name:    "wrap-down-last-visual-line",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Label: "down on last visual line", Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(n), SelStart: intp(n), SelEnd: intp(n)}},
			},
		},
		{
			Name:    "wrap-shift-down-last-to-eof",
			Content: wp,
			Width:   harnessWrapWidth,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Label: "shift+down at wrap end", Keys: []Key{{Name: "ArrowDown", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(n), SelEnd: intp(n), SelLenMin: intp(1)}},
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
	}
}
