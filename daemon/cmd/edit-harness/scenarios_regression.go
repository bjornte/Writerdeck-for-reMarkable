package main

// regressionScenarios cover reported editing bugs on fixtureProse.
// Note: the harness's "cursor" field always mirrors the larger of selStart/
// selEnd (Qt selectionEnd semantics), not necessarily the direction just
// moved — see docs/editor-testing/scenario-cookbook.md.
func regressionScenarios() []Scenario {
	return []Scenario{
		{
			// Reverse partner: up-one-logical-line. Proven at N=1, N=3, N=7.
			Name:    "down-one-logical-line",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseVLineStart(0))},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(0)), TextLen: intp(fixtureProseLen)}},
				{Label: "down x1 (N=1)", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(1)), SelStart: intp(proseVLineStart(1)), SelEnd: intp(proseVLineStart(1))}},
				{Label: "down x2 more (N=3)", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(3)), SelStart: intp(proseVLineStart(3)), SelEnd: intp(proseVLineStart(3))}},
				{Label: "down x4 more (N=7)", Keys: []Key{{Name: "ArrowDown"}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(7)), SelStart: intp(proseVLineStart(7)), SelEnd: intp(proseVLineStart(7))}},
			},
		},
		{
			// Reverse partner: down-one-logical-line.
			Name:    "up-one-logical-line",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseVLineStart(proseVCount - 1))},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(proseVCount - 1)), TextLen: intp(fixtureProseLen)}},
				{Label: "up x1 (N=1)", Keys: []Key{{Name: "ArrowUp"}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(proseVCount - 2)), SelStart: intp(proseVLineStart(proseVCount - 2)), SelEnd: intp(proseVLineStart(proseVCount - 2))}},
				{Label: "up x2 more (N=3)", Keys: []Key{{Name: "ArrowUp"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(proseVCount - 4)), SelStart: intp(proseVLineStart(proseVCount - 4)), SelEnd: intp(proseVLineStart(proseVCount - 4))}},
				{Label: "up x4 more (N=7)", Keys: []Key{{Name: "ArrowUp"}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(proseVCount - 8)), SelStart: intp(proseVLineStart(proseVCount - 8)), SelEnd: intp(proseVLineStart(proseVCount - 8))}},
			},
		},
		{
			// Extend-then-shrink at N=1, N=3, N=7.
			Name:    "shift-down-then-up-shrinks",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseVLineStart(2))},
				{Label: "cursor line2", Expect: &StateExpect{Cursor: intp(proseVLineStart(2))}},
				{Label: "extend down x8", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 8},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(10)), SelStart: intp(proseVLineStart(2)), SelEnd: intp(proseVLineStart(10)), SelLen: intp(proseVLineStart(10) - proseVLineStart(2))}},
				{Label: "shift+up shrinks x1 (N=1)", Keys: []Key{{Name: "ArrowUp", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(9)), SelStart: intp(proseVLineStart(2)), SelEnd: intp(proseVLineStart(9)), SelLen: intp(proseVLineStart(9) - proseVLineStart(2))}},
				{Label: "shift+up shrinks x2 more (N=3)", Keys: []Key{{Name: "ArrowUp", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(7)), SelStart: intp(proseVLineStart(2)), SelEnd: intp(proseVLineStart(7)), SelLen: intp(proseVLineStart(7) - proseVLineStart(2))}},
				{Label: "shift+up shrinks x4 more (N=7)", Keys: []Key{{Name: "ArrowUp", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(3)), SelStart: intp(proseVLineStart(2)), SelEnd: intp(proseVLineStart(3)), SelLen: intp(proseVLineStart(3) - proseVLineStart(2))}},
			},
		},
		{
			// Mid-horizontal-line variant (not only doc/line end).
			Name:    "shift-left-repeat-from-end",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHEditorEnd)},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd)}},
				{Label: "shift+left x1 (N=1)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd), SelStart: intp(proseHEditorEnd - 1), SelEnd: intp(proseHEditorEnd), SelLen: intp(1)}},
				{Label: "shift+left x2 more (N=3)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd), SelStart: intp(proseHEditorEnd - 3), SelEnd: intp(proseHEditorEnd), SelLen: intp(3)}},
				{Label: "shift+left x4 more (N=7)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd), SelStart: intp(proseHEditorEnd - 7), SelEnd: intp(proseHEditorEnd), SelLen: intp(7)}},
			},
		},
		{
			// Backward word deletion from end of word line, N=1/3/7.
			Name:    "alt-backspace-deletes-word",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseWEditorEnd)},
				{Expect: &StateExpect{Cursor: intp(proseWEditorEnd), TextLen: intp(fixtureProseLen)}},
				{Label: "alt+backspace x1 word (N=1)", Keys: []Key{{Name: "Backspace", Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseWordEnds[8]), TextLen: intp(fixtureProseLen - (proseWordEnds[9] - proseWordEnds[8]))}},
				{Label: "alt+backspace x2 more words (N=3)", Keys: []Key{{Name: "Backspace", Alt: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseWordEnds[6])}},
				{Label: "alt+backspace x4 more words (N=7)", Keys: []Key{{Name: "Backspace", Alt: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseWordEnds[2])}},
			},
		},
		{
			// Ctrl+Backspace from last vertical line, N=1/3/7.
			Name:    "ctrl-backspace-deletes-line",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseVLineEnd(proseVCount - 1))},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(proseVCount - 1)), TextLen: intp(fixtureProseLen)}},
				{Label: "ctrl+backspace x1 (N=1)", Keys: []Key{{Name: "Backspace", Ctrl: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(proseVCount - 1)), TextLen: intp(fixtureProseLen - proseVWidth)}},
				{Label: "ctrl+backspace x2 more (N=3)", Keys: []Key{{Name: "Backspace", Ctrl: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(proseVCount - 2))}},
				{Label: "ctrl+backspace x4 more (N=7)", Keys: []Key{{Name: "Backspace", Ctrl: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(proseVCount - 4))}},
			},
		},
		{
			// Mid-document horizontal selection on a vertical line end.
			Name:    "shift-left-repeat-mid-doc",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseVLineEnd(5))},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(5))}},
				{Label: "shift+left x1 (N=1)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(5)), SelStart: intp(proseVLineEnd(5) - 1), SelEnd: intp(proseVLineEnd(5)), SelLen: intp(1)}},
				{Label: "shift+left x2 more (N=3)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(5)), SelStart: intp(proseVLineEnd(5) - 3), SelEnd: intp(proseVLineEnd(5)), SelLen: intp(3)}},
				{Label: "shift+left x4 more (N=7)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(5)), SelStart: intp(proseVLineEnd(5) - 7), SelEnd: intp(proseVLineEnd(5)), SelLen: intp(7)}},
			},
		},
	}
}
