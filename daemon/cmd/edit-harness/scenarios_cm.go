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
			Steps: []Step{
				{SetCursor: intp(proseVLineStart(0))},
				{Label: "down x1 (N=1)", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(1)), SelStart: intp(proseVLineStart(1)), SelEnd: intp(proseVLineStart(1))}},
				{Label: "down x2 more (N=3)", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(3)), SelStart: intp(proseVLineStart(3)), SelEnd: intp(proseVLineStart(3))}},
				{Label: "down x4 more (N=7)", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(7)), SelStart: intp(proseVLineStart(7)), SelEnd: intp(proseVLineStart(7))}},
			},
		},
		{
			Name:    "cm-line-up-basic",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseVLineEnd(proseVCount - 1))},
				{Label: "up x1 (N=1)", Keys: []Key{{Name: "ArrowUp"}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(proseVCount - 2)), SelStart: intp(proseVLineEnd(proseVCount - 2)), SelEnd: intp(proseVLineEnd(proseVCount - 2))}},
				{Label: "up x2 more (N=3)", Keys: []Key{{Name: "ArrowUp"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(proseVCount - 4)), SelStart: intp(proseVLineEnd(proseVCount - 4)), SelEnd: intp(proseVLineEnd(proseVCount - 4))}},
				{Label: "up x4 more (N=7)", Keys: []Key{{Name: "ArrowUp"}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(proseVCount - 8)), SelStart: intp(proseVLineEnd(proseVCount - 8)), SelEnd: intp(proseVLineEnd(proseVCount - 8))}},
			},
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
				{Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(6), SelStart: intp(6), SelEnd: intp(6)}},
			},
		},
		{
			Name:    "cm-select-line-down",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseVLineStart(0))},
				{Label: "shift+down x1 (N=1)", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(1)), SelStart: intp(proseVLineStart(0)), SelEnd: intp(proseVLineStart(1)), SelLen: intp(proseVLineStart(1) - proseVLineStart(0))}},
				{Label: "shift+down x2 more (N=3)", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(3)), SelStart: intp(proseVLineStart(0)), SelEnd: intp(proseVLineStart(3)), SelLen: intp(proseVLineStart(3) - proseVLineStart(0))}},
				{Label: "shift+down x4 more (N=7)", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(7)), SelStart: intp(proseVLineStart(0)), SelEnd: intp(proseVLineStart(7)), SelLen: intp(proseVLineStart(7) - proseVLineStart(0))}},
			},
		},
		{
			Name:    "cm-select-line-down-mid",
			Content: three,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 1},
				{Label: "shift+down from mid", Keys: []Key{{Name: "ArrowDown", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(4), SelStart: intp(1), SelEnd: intp(4), SelLen: intp(3)}},
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
			Name:    "cm-select-up-basic",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseVLineEnd(proseVCount - 1))},
				{Label: "shift+up x1 (N=1)", Keys: []Key{{Name: "ArrowUp", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(proseVCount - 1)), SelStart: intp(proseVLineEnd(proseVCount - 2)), SelEnd: intp(proseVLineEnd(proseVCount - 1))}},
				{Label: "shift+up x2 more (N=3)", Keys: []Key{{Name: "ArrowUp", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(proseVCount - 1)), SelStart: intp(proseVLineEnd(proseVCount - 4)), SelEnd: intp(proseVLineEnd(proseVCount - 1))}},
				{Label: "shift+up x4 more (N=7)", Keys: []Key{{Name: "ArrowUp", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(proseVCount - 1)), SelStart: intp(proseVLineEnd(proseVCount - 8)), SelEnd: intp(proseVLineEnd(proseVCount - 1))}},
			},
		},
		{
			Name:    "cm-select-up-mid",
			Content: three,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Keys: []Key{{Name: "ArrowRight"}}},
				{Label: "shift+up from mid line3", Keys: []Key{{Name: "ArrowUp", Shift: true}}},
				{Expect: &StateExpect{SelLenMin: intp(1)}},
			},
		},
	}
}
