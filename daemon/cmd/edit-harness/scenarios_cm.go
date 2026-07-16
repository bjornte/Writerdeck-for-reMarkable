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
				{Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(6), SelStart: intp(6), SelEnd: intp(6)}},
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
