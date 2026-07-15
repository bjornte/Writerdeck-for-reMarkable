package main

// selectionScenarios cover shift+arrow selection edge cases.
func selectionScenarios() []Scenario {
	return []Scenario{
		{
			Name:    "shift-left-then-right-shrinks",
			Content: "abcdef",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 3},
				{Expect: &StateExpect{Cursor: intp(6), SelStart: intp(3), SelEnd: intp(6), SelLen: intp(3)}},
				{Keys: []Key{{Name: "ArrowRight", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(5), SelStart: intp(3), SelEnd: intp(5), SelLen: intp(2)}},
			},
		},
	}
}
