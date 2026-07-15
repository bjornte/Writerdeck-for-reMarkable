package main

// cmScenarios port CodeMirror vertical motion cases (explicit \n lines).
// See docs/editor-testing/scenario-cookbook.md.
func cmScenarios() []Scenario {
	return []Scenario{
		{
			Name:    "cm-line-down-basic",
			Content: "one\ntwo",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(4), SelStart: intp(4), SelEnd: intp(4)}},
			},
		},
		{
			Name:    "cm-line-down-shorter",
			Content: "one\nt",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(2)}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(4), SelStart: intp(4), SelEnd: intp(4)}},
			},
		},
		{
			Name:    "cm-line-down-last-line",
			Content: "one",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 2},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(3), SelStart: intp(3), SelEnd: intp(3)}},
			},
		},
		{
			Name:    "cm-line-down-goal-col",
			Content: "one\nt\nthree",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 2},
				{Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(6), SelStart: intp(6), SelEnd: intp(6)}},
			},
		},
		{
			Name:    "cm-select-line-down",
			Content: "one\ntwo\nthree",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Label: "shift+down", Keys: []Key{{Name: "ArrowDown", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(4), SelStart: intp(0), SelEnd: intp(4), SelLen: intp(4)}},
			},
		},
		{
			Name:    "cm-select-line-down-mid",
			Content: "one\ntwo\nthree",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 2},
				{Label: "shift+down from mid", Keys: []Key{{Name: "ArrowDown", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(7), SelStart: intp(2), SelEnd: intp(7), SelLen: intp(5)}},
			},
		},
		{
			Name:    "cm-select-down-up-doc-end",
			Content: "one\ntwo\nthree",
			Steps: []Step{
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(13)}},
				{Keys: []Key{{Name: "ArrowDown", Shift: true}}},
				{Keys: []Key{{Name: "ArrowUp", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(12), SelStart: intp(8), SelEnd: intp(12), SelLen: intp(4)}},
			},
		},
		{
			Name:    "cm-select-up-basic",
			Content: "one\ntwo\nthree",
			Steps: []Step{
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Label: "shift+up", Keys: []Key{{Name: "ArrowUp", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(13), SelStart: intp(4), SelEnd: intp(13), SelLen: intp(9)}},
			},
		},
		{
			Name:    "cm-select-up-mid",
			Content: "one\ntwo\nthree",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Keys: []Key{{Name: "ArrowRight"}}},
				{Label: "shift+up from mid line3", Keys: []Key{{Name: "ArrowUp", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(9), SelStart: intp(4), SelEnd: intp(9), SelLen: intp(5)}},
			},
		},
	}
}
