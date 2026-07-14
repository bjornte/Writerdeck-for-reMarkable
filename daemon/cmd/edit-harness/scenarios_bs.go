package main

// bsScenarios extend backspace/delete coverage beyond basic regression cases.
func bsScenarios() []Scenario {
	return []Scenario{
		{
			Name:    "bs-alt-word-mid",
			Content: "hello world",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 8},
				{Expect: &StateExpect{Cursor: intp(8), TextLen: intp(11)}},
				{Label: "alt+backspace mid word", Keys: []Key{{Name: "Backspace", Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(5), TextLen: intp(5)}},
			},
		},
		{
			Name:    "bs-ctrl-line-start",
			Content: "line1\nline2",
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(6), TextLen: intp(11)}},
				{Label: "ctrl+backspace at line2 start", Keys: []Key{{Name: "Backspace", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(6), TextLen: intp(6)}},
			},
		},
		{
			Name:    "bs-shift-with-selection",
			Content: "abcd",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "Home", Shift: true}}},
				{Expect: &StateExpect{SelStart: intp(0), SelEnd: intp(4), SelLen: intp(4)}},
				{Keys: []Key{{Name: "Backspace", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(0), TextLen: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "bs-plain",
			Content: "abcd",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "Backspace"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(2), TextLen: intp(2)}},
			},
		},
	}
}
