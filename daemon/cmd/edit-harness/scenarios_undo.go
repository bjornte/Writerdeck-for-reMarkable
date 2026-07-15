package main

// undoScenarios port Qt/CodeMirror undo cases (see docs/editor-testing/scenario-cookbook.md).
// Single-char keys use the WebSocket text path (translate → {"t":"text","cp":…}).
func undoScenarios() []Scenario {
	return []Scenario{
		{
			Name:    "undo-redo-len",
			Content: "abc d",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(5), TextLen: intp(5)}},
				{Label: "select all", Keys: []Key{{Name: "a", Ctrl: true}}},
				{Expect: &StateExpect{SelLen: intp(5), Cursor: intp(5)}},
				{Label: "delete all", Keys: []Key{{Name: "Backspace"}}},
				{Expect: &StateExpect{Cursor: intp(0), TextLen: intp(0)}},
				{Label: "ctrl+z restore", Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(5), TextLen: intp(5)}},
				{Label: "ctrl+y redo delete", Keys: []Key{{Name: "y", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0), TextLen: intp(0)}},
			},
		},
		{
			Name:    "undo-cursor-reposition",
			Content: "five\nlines\nin\nthis\ntextedit",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0), TextLen: intp(27)}},
				{Label: "insert Blah", Keys: []Key{{Name: "B"}, {Name: "l"}, {Name: "a"}, {Name: "h"}}},
				{Expect: &StateExpect{Cursor: intp(4), TextLen: intp(31)}},
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(31), TextLen: intp(31)}},
				{Label: "undo from eof", Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0), TextLen: intp(27)}},
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(27), TextLen: intp(27)}},
				{Label: "redo restores insert cursor", Keys: []Key{{Name: "y", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(4), TextLen: intp(31)}},
			},
		},
		{
			Name:    "undo-mid-line-delete",
			Content: "abc\ndef",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(4), TextLen: intp(7)}},
				{Label: "select line2", Keys: []Key{{Name: "End", Shift: true}}},
				{Expect: &StateExpect{SelStart: intp(4), SelEnd: intp(7), SelLen: intp(3)}},
				{Label: "delete line2", Keys: []Key{{Name: "Backspace"}}},
				{Expect: &StateExpect{Cursor: intp(4), TextLen: intp(4)}},
				{Label: "undo", Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(4), TextLen: intp(7)}},
			},
		},
		{
			Name:    "redo-cleared-by-new-edit",
			Content: "abc",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "a", Ctrl: true}}},
				{Keys: []Key{{Name: "Backspace"}}},
				{Expect: &StateExpect{TextLen: intp(0)}},
				{Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(3), Cursor: intp(3)}},
				{Keys: []Key{{Name: "Backspace"}}},
				{Expect: &StateExpect{TextLen: intp(2), Cursor: intp(2)}},
				{Label: "redo dead after new edit", Keys: []Key{{Name: "y", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(2), Cursor: intp(2)}},
			},
		},
		{
			Name:    "undo-after-select-delete",
			Content: "abcdef",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Label: "shift+home select line", Keys: []Key{{Name: "Home", Shift: true}}},
				{Expect: &StateExpect{SelStart: intp(0), SelEnd: intp(6), Cursor: intp(6), SelLen: intp(6)}},
				{Keys: []Key{{Name: "Backspace"}}},
				{Expect: &StateExpect{TextLen: intp(0), Cursor: intp(0)}},
				{Label: "undo", Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(6), Cursor: intp(6), SelStart: intp(6), SelEnd: intp(6)}},
			},
		},
	}
}
