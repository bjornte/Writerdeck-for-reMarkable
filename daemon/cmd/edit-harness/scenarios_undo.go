package main

// undoScenarios port Qt/CodeMirror undo cases.
// Single-char keys use the WebSocket text path (translate → {"t":"text","cp":…}).
func undoScenarios() []Scenario {
	const body = "abc æøå"
	bodyLen := editorLen(body)
	const multi = "fem\nlinjer\ni\ndette\nnotatet"
	multiLen := editorLen(multi)

	return []Scenario{
		{
			Name:    "undo-redo-len",
			Content: body,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(bodyLen), TextLen: intp(bodyLen)}},
				{Label: "select all", Keys: []Key{{Name: "a", Ctrl: true}}},
				{Expect: &StateExpect{SelLen: intp(bodyLen), Cursor: intp(bodyLen)}},
				{Label: "delete all", Keys: []Key{{Name: "Backspace"}}},
				{Expect: &StateExpect{Cursor: intp(0), TextLen: intp(0)}},
				{Label: "ctrl+z restore", Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(bodyLen), TextLen: intp(bodyLen)}},
				{Label: "ctrl+y redo delete", Keys: []Key{{Name: "y", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0), TextLen: intp(0)}},
			},
		},
		{
			Name:    "undo-cursor-reposition",
			Content: multi,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0), TextLen: intp(multiLen)}},
				{Label: "insert Blah", Keys: []Key{{Name: "B"}, {Name: "l"}, {Name: "a"}, {Name: "h"}}},
				{Expect: &StateExpect{Cursor: intp(4), TextLen: intp(multiLen + 4)}},
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(multiLen + 4), TextLen: intp(multiLen + 4)}},
				{Label: "undo from eof", Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0), TextLen: intp(multiLen)}},
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(multiLen), TextLen: intp(multiLen)}},
				{Label: "redo restores insert cursor", Keys: []Key{{Name: "y", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(4), TextLen: intp(multiLen + 4)}},
			},
		},
		{
			Name:    "undo-mid-line-delete",
			Content: "ost\nrømme",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(4), TextLen: intp(editorLen("ost\nrømme"))}},
				{Label: "select line2", Keys: []Key{{Name: "End", Shift: true}}},
				{Expect: &StateExpect{SelStart: intp(4), SelEnd: intp(editorLen("ost\nrømme")), SelLen: intp(editorLen("rømme"))}},
				{Label: "delete line2", Keys: []Key{{Name: "Backspace"}}},
				{Expect: &StateExpect{Cursor: intp(4), TextLen: intp(4)}},
				{Label: "undo", Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(editorLen("ost\nrømme")), TextLen: intp(editorLen("ost\nrømme"))}},
			},
		},
		{
			Name:    "redo-cleared-by-new-edit",
			Content: body,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "a", Ctrl: true}}},
				{Keys: []Key{{Name: "Backspace"}}},
				{Expect: &StateExpect{TextLen: intp(0)}},
				{Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(bodyLen), Cursor: intp(bodyLen)}},
				{Keys: []Key{{Name: "Backspace"}}},
				{Expect: &StateExpect{TextLen: intp(bodyLen - 1), Cursor: intp(bodyLen - 1)}},
				{Label: "redo dead after new edit", Keys: []Key{{Name: "y", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(bodyLen - 1), Cursor: intp(bodyLen - 1)}},
			},
		},
		{
			Name:    "undo-after-select-delete",
			Content: body,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Label: "shift+home select line", Keys: []Key{{Name: "Home", Shift: true}}},
				{Expect: &StateExpect{SelStart: intp(0), SelEnd: intp(bodyLen), Cursor: intp(bodyLen), SelLen: intp(bodyLen)}},
				{Keys: []Key{{Name: "Backspace"}}},
				{Expect: &StateExpect{TextLen: intp(0), Cursor: intp(0)}},
				{Label: "undo", Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(bodyLen), Cursor: intp(bodyLen), SelStart: intp(bodyLen), SelEnd: intp(bodyLen)}},
			},
		},
	}
}
