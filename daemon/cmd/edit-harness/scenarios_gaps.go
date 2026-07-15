package main

// gapScenarios cover Mac-style editing cases missing from the core/combo/wrap blocks.
func gapScenarios() []Scenario {
	const resume = "test résumé"
	resumeLen := utf8Len(resume)
	resumeWordEnd := utf8ByteAtRune(resume, 4) // byte index after "test"

	return []Scenario{
		{
			Name:    "gap-up-at-doc-start",
			Content: "hello",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0)}},
				{Label: "up at start", Keys: []Key{{Name: "ArrowUp"}}},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "gap-plain-left-moves-caret",
			Content: "hello world",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(11)}},
				{Label: "plain left moves caret", Keys: []Key{{Name: "ArrowLeft"}}},
				{Expect: &StateExpect{Cursor: intp(10), SelStart: intp(10), SelEnd: intp(10)}},
			},
		},
		{
			Name:    "gap-plain-right-moves-caret",
			Content: "hello",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0)}},
				{Label: "plain right moves caret", Keys: []Key{{Name: "ArrowRight"}}},
				{Expect: &StateExpect{Cursor: intp(1), SelStart: intp(1), SelEnd: intp(1)}},
			},
		},
		{
			Name:    "gap-plain-left-at-doc-start",
			Content: "hello",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0)}},
				{Label: "left clamps at start", Keys: []Key{{Name: "ArrowLeft"}}},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "gap-plain-right-at-doc-end",
			Content: "hello",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(5)}},
				{Label: "right clamps at end", Keys: []Key{{Name: "ArrowRight"}}},
				{Expect: &StateExpect{Cursor: intp(5), SelStart: intp(5), SelEnd: intp(5)}},
			},
		},
		{
			Name:    "gap-collapse-selection-left",
			Content: "abcdef",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Label: "shift+left x3", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 3},
				{Expect: &StateExpect{Cursor: intp(6), SelStart: intp(3), SelEnd: intp(6), SelLen: intp(3)}},
				{Label: "left collapses selection", Keys: []Key{{Name: "ArrowLeft"}}},
				{Expect: &StateExpect{Cursor: intp(3), SelStart: intp(3), SelEnd: intp(3), SelLen: intp(0)}},
			},
		},
		{
			Name:    "gap-collapse-selection-right",
			Content: "abcdef",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Label: "shift+right x3", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 3},
				{Expect: &StateExpect{Cursor: intp(3), SelStart: intp(0), SelEnd: intp(3), SelLen: intp(3)}},
				{Label: "right collapses selection", Keys: []Key{{Name: "ArrowRight"}}},
				{Expect: &StateExpect{Cursor: intp(3), SelStart: intp(3), SelEnd: intp(3), SelLen: intp(0)}},
			},
		},
		{
			Name:    "gap-delete-forward",
			Content: "abcd",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowRight"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(2), TextLen: intp(4)}},
				{Label: "delete char at cursor", Keys: []Key{{Name: "Delete"}}},
				{Expect: &StateExpect{Cursor: intp(2), TextLen: intp(3), Text: strp("abd")}},
			},
		},
		{
			Name:    "gap-delete-with-selection",
			Content: "abcdef",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "Home", Shift: true}}},
				{Expect: &StateExpect{SelStart: intp(0), SelEnd: intp(6), SelLen: intp(6)}},
				{Label: "delete selection", Keys: []Key{{Name: "Delete"}}},
				{Expect: &StateExpect{Cursor: intp(0), TextLen: intp(0), SelStart: intp(0), SelEnd: intp(0), Text: strp("")}},
			},
		},
		{
			Name:    "gap-select-all",
			Content: "hello",
			Steps: []Step{
				{Label: "ctrl+a", Keys: []Key{{Name: "a", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(5), SelStart: intp(0), SelEnd: intp(5), SelLen: intp(5), TextLen: intp(5)}},
			},
		},
		{
			Name:    "gap-enter-new-line",
			Content: "ab",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(2), TextLen: intp(2)}},
				{Label: "enter inserts newline", Keys: []Key{{Name: "Enter"}}},
				{Expect: &StateExpect{Cursor: intp(3), TextLen: intp(3), Text: strp("ab\n")}},
			},
		},
		{
			Name:    "gap-type-replaces-selection",
			Content: "abcdef",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "Home", Shift: true}}},
				{Expect: &StateExpect{SelStart: intp(0), SelEnd: intp(6), SelLen: intp(6)}},
				{Label: "type replaces selection", Keys: []Key{{Name: "x"}}},
				{Expect: &StateExpect{Cursor: intp(1), TextLen: intp(1), SelStart: intp(1), SelEnd: intp(1), Text: strp("x")}},
			},
		},
		{
			Name:    "gap-redo-shift-ctrl-z",
			Content: "abc",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "a", Ctrl: true}}},
				{Keys: []Key{{Name: "Backspace"}}},
				{Expect: &StateExpect{TextLen: intp(0), Cursor: intp(0)}},
				{Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(3), Cursor: intp(3)}},
				{Label: "shift+ctrl+z redo", Keys: []Key{{Name: "z", Ctrl: true, Shift: true}}},
				{Expect: &StateExpect{TextLen: intp(0), Cursor: intp(0)}},
			},
		},
		{
			Name:    "gap-undo-chain",
			Content: "abc",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "Backspace"}}, Repeat: 2},
				{Expect: &StateExpect{TextLen: intp(1), Cursor: intp(1)}},
				{Label: "undo restore b", Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(2), Cursor: intp(2)}},
				{Label: "undo restore c", Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(3), Cursor: intp(3)}},
			},
		},
		{
			Name:    "gap-unicode-alt-backspace",
			Content: resume,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(resumeLen), TextLen: intp(resumeLen)}},
				{Label: "alt+backspace unicode word", Keys: []Key{{Name: "Backspace", Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(resumeWordEnd), TextLen: intp(resumeWordEnd), Text: strp("test")}},
			},
		},
		{
			Name:    "gap-empty-doc-backspace",
			Content: "",
			Steps: []Step{
				{Expect: &StateExpect{Cursor: intp(0), TextLen: intp(0), Text: strp("")}},
				{Label: "backspace on empty", Keys: []Key{{Name: "Backspace"}}},
				{Expect: &StateExpect{Cursor: intp(0), TextLen: intp(0), Text: strp("")}},
			},
		},
		{
			Name:    "gap-alt-bs-with-selection",
			Content: "hello world",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "ArrowLeft", Shift: true, Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(11), SelStart: intp(6), SelEnd: intp(11), SelLen: intp(5)}},
				{Label: "alt+backspace deletes selection", Keys: []Key{{Name: "Backspace", Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(5), TextLen: intp(5), SelStart: intp(5), SelEnd: intp(5), Text: strp("hello")}},
			},
		},
	}
}
