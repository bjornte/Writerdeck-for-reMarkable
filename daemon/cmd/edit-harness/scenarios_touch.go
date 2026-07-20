package main

// touchScenarios verify visual goal-x is preserved after tap (harnessSetCursor).
func touchScenarios() []Scenario {
	return []Scenario{
		{
			Name:    "touch-down-goal-column",
			Content: fixtureThreeLines,
			Steps: []Step{
				{Label: "tap mid line 1", SetCursor: intp(1)},
				{Expect: &StateExpect{Cursor: intp(1), SelStart: intp(1), SelEnd: intp(1)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(4), SelStart: intp(4), SelEnd: intp(4)}},
			},
		},
		{
			Name:    "touch-up-goal-column",
			Content: fixtureThreeLines,
			Steps: []Step{
				{Label: "tap mid line 2", SetCursor: intp(4)},
				{Expect: &StateExpect{Cursor: intp(4)}},
				{Keys: []Key{{Name: "ArrowUp"}}},
				{Expect: &StateExpect{Cursor: intp(1), SelStart: intp(1), SelEnd: intp(1)}},
			},
		},
		{
			Name:    "touch-down-shorter-line",
			Content: fixtureShorterDown,
			Steps: []Step{
				{Label: "tap mid line 1", SetCursor: intp(2)},
				{Expect: &StateExpect{Cursor: intp(2)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(4), SelStart: intp(4), SelEnd: intp(4)}},
			},
		},
	}
}
