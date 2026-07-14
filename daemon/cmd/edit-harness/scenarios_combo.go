package main

// comboScenarios cover Alt/Ctrl and Shift+Alt/Ctrl arrow bindings (Mac-style, phone path).
func comboScenarios() []Scenario {
	const helloWorld = "hello world"
	const twoLines = "abc\ndef"
	const threeLines = "one\ntwo\nthree"
	const twoParas = "para1\n\npara2"

	return []Scenario{
		// Alt / Ctrl navigation (4 arrows).
		{
			Name:    "combo-alt-left",
			Content: helloWorld,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "ArrowLeft", Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(6), SelStart: intp(6), SelEnd: intp(6)}},
			},
		},
		{
			Name:    "combo-alt-right",
			Content: helloWorld,
			Steps: []Step{
				// prepare leaves cursor at 0; skip Ctrl+Home (press-only poisons next modified key)
				{Keys: []Key{{Name: "ArrowRight", Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(6), SelStart: intp(6), SelEnd: intp(6)}},
			},
		},
		{
			Name:    "combo-alt-up",
			Content: twoParas,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "ArrowUp", Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "combo-alt-down",
			Content: twoParas,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown", Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(7), SelStart: intp(7), SelEnd: intp(7)}},
			},
		},
		{
			Name:    "combo-ctrl-left",
			Content: helloWorld,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "ArrowLeft", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "combo-ctrl-right",
			Content: helloWorld,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowRight", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(11), SelStart: intp(11), SelEnd: intp(11)}},
			},
		},
		{
			Name:    "combo-ctrl-up",
			Content: threeLines,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "ArrowUp", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "combo-ctrl-down",
			Content: threeLines,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(13), SelStart: intp(13), SelEnd: intp(13)}},
			},
		},

		// Shift+Alt all arrows.
		{
			Name:    "combo-shift-alt-left",
			Content: helloWorld,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "ArrowLeft", Shift: true, Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(11), SelStart: intp(6), SelEnd: intp(11), SelLen: intp(5)}},
			},
		},
		{
			Name:    "combo-shift-alt-right",
			Content: helloWorld,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowRight", Shift: true, Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(6), SelStart: intp(0), SelEnd: intp(6), SelLen: intp(6)}},
			},
		},
		{
			Name:    "combo-shift-alt-up",
			Content: twoParas,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "ArrowUp", Shift: true, Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(12), SelStart: intp(0), SelEnd: intp(12), SelLen: intp(12)}},
			},
		},
		{
			Name:    "combo-shift-alt-down",
			Content: twoParas,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown", Shift: true, Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(7), SelStart: intp(0), SelEnd: intp(7), SelLen: intp(7)}},
			},
		},

		// Shift+Ctrl all arrows.
		{
			Name:    "combo-shift-ctrl-left",
			Content: helloWorld,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "ArrowLeft", Shift: true, Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(11), SelStart: intp(0), SelEnd: intp(11), SelLen: intp(11)}},
			},
		},
		{
			Name:    "combo-shift-ctrl-right",
			Content: helloWorld,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowRight", Shift: true, Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(11), SelStart: intp(0), SelEnd: intp(11), SelLen: intp(11)}},
			},
		},
		{
			Name:    "combo-shift-ctrl-up",
			Content: threeLines,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "ArrowUp", Shift: true, Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(13), SelStart: intp(0), SelEnd: intp(13), SelLen: intp(13)}},
			},
		},
		{
			Name:    "combo-shift-ctrl-down",
			Content: threeLines,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown", Shift: true, Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(13), SelStart: intp(0), SelEnd: intp(13), SelLen: intp(13)}},
			},
		},

		// Home / End combos.
		{
			Name:    "combo-shift-home-line",
			Content: twoLines,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "Home", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(7), SelStart: intp(0), SelEnd: intp(7), SelLen: intp(7)}},
			},
		},
		{
			Name:    "combo-shift-end-line",
			Content: twoLines,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "Home"}}},
				{Keys: []Key{{Name: "End", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(7), SelStart: intp(4), SelEnd: intp(7), SelLen: intp(3)}},
			},
		},
		{
			Name:    "combo-ctrl-home",
			Content: twoLines,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "combo-ctrl-end",
			Content: twoLines,
			Steps: []Step{
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(7), SelStart: intp(7), SelEnd: intp(7)}},
			},
		},
		{
			Name:    "combo-shift-ctrl-home",
			Content: twoLines,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "Home", Shift: true, Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(4), SelStart: intp(0), SelEnd: intp(4), SelLen: intp(4)}},
			},
		},
		{
			Name:    "combo-shift-ctrl-end",
			Content: twoLines,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "End", Shift: true, Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(7), SelStart: intp(4), SelEnd: intp(7), SelLen: intp(3)}},
			},
		},
	}
}
