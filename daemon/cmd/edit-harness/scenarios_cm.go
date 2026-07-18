package main

// cmScenarios port CodeMirror vertical motion cases. Prefer fixtureProse for
// N=1/3/7 walks; keep short specialized shapes for goal-column / shorter-line.
func cmScenarios() []Scenario {
	three := fixtureThreeLines
	threeLen := editorLen(three)
	return []Scenario{
		{
			Name:    "cm-line-down-basic",
			Content: fixtureProse,
			Steps:   verticalLinePattern(2, keyArrow("ArrowDown", false, false, false), keyArrow("ArrowUp", false, false, false), +1),
		},
		{
			Name:    "cm-line-up-basic",
			Content: fixtureProse,
			Steps:   verticalLinePattern(9, keyArrow("ArrowUp", false, false, false), keyArrow("ArrowDown", false, false, false), -1),
		},
		{
			Name:    "cm-line-down-shorter",
			Content: fixtureShorterDown,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(2)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(4), SelStart: intp(4), SelEnd: intp(4)}},
			},
		},
		{
			Name:    "cm-line-up-shorter",
			Content: fixtureShorterUp,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(2)}},
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(5)}},
				{Label: "up clamps onto shorter line above", Keys: []Key{{Name: "ArrowUp"}}},
				{Expect: &StateExpect{Cursor: intp(1), SelStart: intp(1), SelEnd: intp(1)}},
			},
		},
		{
			Name:    "cm-line-down-last-line",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseVLineStart(proseVCount-1) + 2)},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(proseVCount - 1)), SelStart: intp(proseVLineEnd(proseVCount - 1)), SelEnd: intp(proseVLineEnd(proseVCount - 1))}},
				{Label: "down stays clamped at doc end (N=3 total)", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(proseVCount - 1)), SelStart: intp(proseVLineEnd(proseVCount - 1)), SelEnd: intp(proseVLineEnd(proseVCount - 1))}},
				{Label: "still clamped (N=7 total)", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(proseVCount - 1)), SelStart: intp(proseVLineEnd(proseVCount - 1)), SelEnd: intp(proseVLineEnd(proseVCount - 1))}},
			},
		},
		{
			Name:    "cm-line-down-goal-col",
			Content: fixtureGoalColDown,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 2},
				// "tre\ni\nfemte": col 2 of "tre" sticky across short "i" → col 2 of "femte" (pos 8).
				// Allow 7–8 for e-ink glyph x rounding.
				{Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Expect: &StateExpect{CursorMin: intp(7), CursorMax: intp(8)}},
			},
		},
		{
			Name:    "cm-select-line-down",
			Content: fixtureProse,
			Steps:   shiftVerticalPattern(2, keyArrow("ArrowDown", true, false, false), keyArrow("ArrowUp", true, false, false), false),
		},
		{
			Name:    "cm-select-up-basic",
			Content: fixtureProse,
			Steps:   shiftVerticalPattern(9, keyArrow("ArrowUp", true, false, false), keyArrow("ArrowDown", true, false, false), true),
		},
		{
			// Mid-sentence wrapping prose (replaces toy en/to/tre mid-line).
			Name:    "cm-select-line-down-mid",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseMidDocCaret)},
				{Label: "shift+down from mid wrapping", Keys: []Key{{Name: "ArrowDown", Shift: true}}},
				{Expect: &StateExpect{
					SelStart:  intp(proseMidDocCaret),
					CursorMin: intp(proseMidDocCaret + 1),
					CursorMax: intp(proseMidDocCaret + 120),
					SelLenMin: intp(1),
					SelLenMax: intp(120),
				}},
			},
		},
		{
			Name:    "cm-select-down-up-doc-end",
			Content: three,
			Steps: []Step{
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(threeLen)}},
				{Label: "shift+down clamped at doc end (N=1)", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(threeLen), SelStart: intp(threeLen), SelEnd: intp(threeLen)}},
				{Label: "shift+down still clamped (N=3)", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(threeLen), SelStart: intp(threeLen), SelEnd: intp(threeLen)}},
				{Keys: []Key{{Name: "ArrowUp", Shift: true}}},
				{Expect: &StateExpect{SelLenMin: intp(1), Cursor: intp(threeLen)}},
			},
		},
		{
			Name:    "cm-select-up-mid",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(prosePara2Mid)},
				{Label: "shift+up from mid wrapping", Keys: []Key{{Name: "ArrowUp", Shift: true}}},
				{Expect: &StateExpect{
					SelEnd:    intp(prosePara2Mid),
					SelLenMin: intp(1),
					SelLenMax: intp(120),
				}},
			},
		},
	}
}
