package main

// touchScenarios verify visual goal-x is preserved after tap (harnessSetCursor simulates touch).
func touchScenarios() []Scenario {
	return []Scenario{
		{
			Name:    "touch-down-goal-column",
			Content: "one\ntwo",
			Steps: []Step{
				{Label: "tap mid line 1", SetCursor: intp(2)},
				{Expect: &StateExpect{Cursor: intp(2), SelStart: intp(2), SelEnd: intp(2)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(6), SelStart: intp(6), SelEnd: intp(6)}},
			},
		},
		{
			Name:    "touch-up-goal-column",
			Content: "one\ntwo\nthree",
			Steps: []Step{
				{Label: "tap mid line 2", SetCursor: intp(6)},
				{Expect: &StateExpect{Cursor: intp(6)}},
				{Keys: []Key{{Name: "ArrowUp"}}},
				{Expect: &StateExpect{Cursor: intp(2), SelStart: intp(2), SelEnd: intp(2)}},
			},
		},
		{
			Name:    "touch-down-shorter-line",
			Content: "one\nt",
			Steps: []Step{
				{Label: "tap mid line 1", SetCursor: intp(2)},
				{Expect: &StateExpect{Cursor: intp(2)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(4), SelStart: intp(4), SelEnd: intp(4)}},
			},
		},
	}
}
