package main

// regressionScenarios: logical-line caret motion and word/line delete.
// Caret Down/Up use uni1/uni5/bi1+1/bi3+5/bi7+7. Shift shrink lives in selection.
func regressionScenarios() []Scenario {
	down := keyArrow("ArrowDown", false, false, false)
	up := keyArrow("ArrowUp", false, false, false)

	return []Scenario{
		{
			Name:    "down-one-logical-line",
			Content: fixtureProse,
			Steps:   verticalLinePattern(2, down, up, +1),
		},
		{
			Name:    "up-one-logical-line",
			Content: fixtureProse,
			Steps:   verticalLinePattern(9, up, down, -1),
		},
		{
			// Covered fully by selectionScenarios shift-down-then-up-shrinks;
			// keep name as critical alias with the same pattern.
			Name:    "shift-down-then-up-shrinks",
			Content: fixtureProse,
			Steps:   shiftVerticalPattern(2, keyArrow("ArrowDown", true, false, false), keyArrow("ArrowUp", true, false, false), false),
		},
		{
			Name:    "shift-left-repeat-from-end",
			Content: fixtureProse,
			Steps:   shiftSelectAxisPattern(proseHEditorEnd, keyArrow("ArrowLeft", true, false, false), keyArrow("ArrowRight", true, false, false), true),
		},
		{
			Name:    "alt-backspace-deletes-word",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseWEditorEnd)},
				{Label: "uni 1", Keys: []Key{{Name: "Backspace", Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseWordEnds[len(proseWordEnds)-2])}},
				{Label: "reset uni5", SetCursor: intp(proseWEditorEnd)},
				{Label: "uni 5", Keys: []Key{{Name: "Backspace", Alt: true}}, Repeat: 5},
				{Expect: &StateExpect{Cursor: intp(proseWordEnds[len(proseWordEnds)-6])}},
			},
		},
		{
			Name:    "ctrl-backspace-deletes-line",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseVLineEnd(proseVCount - 1))},
				{Label: "uni 1", Keys: []Key{{Name: "Backspace", Ctrl: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(proseVCount - 1))}},
				{Label: "reset", SetCursor: intp(proseVLineEnd(proseVCount - 1))},
				{Label: "uni 5", Keys: []Key{{Name: "Backspace", Ctrl: true}}, Repeat: 5},
				{Expect: &StateExpect{CursorMin: intp(proseVLineStart(proseVCount - 5)), CursorMax: intp(proseVLineStart(proseVCount - 1))}},
			},
		},
		{
			Name:    "shift-left-repeat-mid-doc",
			Content: fixtureProse,
			Steps:   shiftSelectAxisPattern(proseVLineEnd(5), keyArrow("ArrowLeft", true, false, false), keyArrow("ArrowRight", true, false, false), true),
		},
	}
}
