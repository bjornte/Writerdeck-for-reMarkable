package main

// bsScenarios extend backspace/delete coverage. Plain backspace/delete use the
// horizontal line of fixtureProse at N=1/3/7 in both directions.
func bsScenarios() []Scenario {
	return []Scenario{
		{
			Name:    "bs-alt-word-mid",
			Content: fixtureProse,
			Steps: []Step{
				// Mid "foxtrot" on the word line.
				{SetCursor: intp(proseWordStarts[5] + 3)},
				{Expect: &StateExpect{Cursor: intp(proseWordStarts[5] + 3), TextLen: intp(fixtureProseLen)}},
				{Label: "alt+backspace mid word", Keys: []Key{{Name: "Backspace", Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(proseWordStarts[5]), TextLen: intp(fixtureProseLen - 3)}},
			},
		},
		{
			Name:    "bs-ctrl-line-start",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseVLineStart(1))},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(1)), TextLen: intp(fixtureProseLen)}},
				{Label: "ctrl+backspace at line start", Keys: []Key{{Name: "Backspace", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(0)), TextLen: intp(fixtureProseLen - proseVWidth - 1)}},
			},
		},
		{
			Name:    "bs-shift-with-selection",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHEditorEnd)},
				{Keys: []Key{{Name: "Home", Shift: true}}},
				{Expect: &StateExpect{SelStart: intp(proseHStart), SelEnd: intp(proseHEditorEnd), SelLen: intp(proseHLen)}},
				{Keys: []Key{{Name: "Backspace", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(proseHStart), TextLen: intp(fixtureProseLen - proseHLen), SelStart: intp(proseHStart), SelEnd: intp(proseHStart)}},
			},
		},
		{
			Name:    "bs-plain",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHEditorEnd)},
				{Label: "backspace x1 (N=1)", Keys: []Key{{Name: "Backspace"}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd - 1), TextLen: intp(fixtureProseLen - 1)}},
				{Label: "backspace x2 more (N=3)", Keys: []Key{{Name: "Backspace"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd - 3), TextLen: intp(fixtureProseLen - 3)}},
				{Label: "backspace x4 more (N=7)", Keys: []Key{{Name: "Backspace"}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd - 7), TextLen: intp(fixtureProseLen - 7)}},
			},
		},
		{
			Name:    "delete-repeat-forward",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHStart)},
				{Label: "delete x1 (N=1)", Keys: []Key{{Name: "Delete"}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseHStart), TextLen: intp(fixtureProseLen - 1)}},
				{Label: "delete x2 more (N=3)", Keys: []Key{{Name: "Delete"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseHStart), TextLen: intp(fixtureProseLen - 3)}},
				{Label: "delete x4 more (N=7)", Keys: []Key{{Name: "Delete"}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseHStart), TextLen: intp(fixtureProseLen - 7)}},
			},
		},
	}
}
