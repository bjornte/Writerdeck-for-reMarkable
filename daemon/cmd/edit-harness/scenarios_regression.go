package main

// regressionScenarios cover reported editing bugs. Some may fail until handleKey is fixed.
func regressionScenarios() []Scenario {
	return []Scenario{
		{
			// Tests Down across explicit \n — not wrapped visual lines (see wrap-* scenarios).
			Name:    "down-one-logical-line",
			Content: "aa\nbb",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0), TextLen: intp(5)}},
				{Label: "down one line", Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(3), SelStart: intp(3), SelEnd: intp(3)}},
			},
		},
		{
			Name:    "shift-down-then-up-shrinks",
			Content: "line1\nline2\nline3\nline4",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Label: "cursor line3", Expect: &StateExpect{Cursor: intp(12)}},
				{Label: "shift+down", Keys: []Key{{Name: "ArrowDown", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(18), SelStart: intp(12), SelEnd: intp(18), SelLen: intp(6)}},
				{Label: "shift+up shrinks downward selection", Keys: []Key{{Name: "ArrowUp", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(17), SelStart: intp(12), SelEnd: intp(17), SelLen: intp(5)}},
			},
		},
		{
			Name:    "shift-left-repeat-from-end",
			Content: "one two three",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(13)}},
				{Label: "shift+left x3", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 3},
				{Expect: &StateExpect{Cursor: intp(13), SelStart: intp(10), SelEnd: intp(13), SelLen: intp(3)}},
			},
		},
		{
			Name:    "alt-backspace-deletes-word",
			Content: "hello world",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(11), TextLen: intp(11)}},
				{Label: "alt+backspace", Keys: []Key{{Name: "Backspace", Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(5), TextLen: intp(5)}},
			},
		},
		{
			Name:    "ctrl-backspace-deletes-line",
			Content: "line1\nline2",
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(11), TextLen: intp(11)}},
				{Label: "ctrl+backspace", Keys: []Key{{Name: "Backspace", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(6)}},
			},
		},
		{
			Name:    "shift-left-repeat-mid-doc",
			Content: "abc\ndef ghij",
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(12)}},
				{Label: "shift+left x3", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 3},
				{Expect: &StateExpect{Cursor: intp(12), SelStart: intp(9), SelEnd: intp(12), SelLen: intp(3)}},
			},
		},
	}
}
