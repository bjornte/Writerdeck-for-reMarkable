package main

// gapScenarios cover Mac-style editing cases missing from core/combo/wrap.
// Motion/selection gaps use fixtureProse at N=1/3/7 both directions.
func gapScenarios() []Scenario {
	const resume = "test résumé æøå"
	resumeLen := editorLen(resume)
	resumeWordEnd := editorLen("test résumé")

	return []Scenario{
		{
			Name:    "gap-up-at-doc-start",
			Content: fixtureProse,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0)}},
				{Label: "up at start x1", Keys: []Key{{Name: "ArrowUp"}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "up at start x3 total", Keys: []Key{{Name: "ArrowUp"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "up at start x7 total", Keys: []Key{{Name: "ArrowUp"}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "gap-plain-left-moves-caret",
			Content: fixtureProse,
			Steps:   caretAxisPattern(proseHMid, keyArrow("ArrowLeft", false, false, false), keyArrow("ArrowRight", false, false, false), -1),
		},
		{
			Name:    "gap-plain-right-moves-caret",
			Content: fixtureProse,
			Steps:   caretAxisPattern(proseHMid, keyArrow("ArrowRight", false, false, false), keyArrow("ArrowLeft", false, false, false), +1),
		},
		{
			// Mid paragraph for variety (not only horizontal line).
			Name:    "gap-plain-left-in-paragraph",
			Content: fixtureProse,
			Steps:   caretAxisPattern(proseMidDocCaret, keyArrow("ArrowLeft", false, false, false), keyArrow("ArrowRight", false, false, false), -1),
		},
		{
			Name:    "gap-plain-right-in-paragraph",
			Content: fixtureProse,
			Steps:   caretAxisPattern(prosePara2Mid, keyArrow("ArrowRight", false, false, false), keyArrow("ArrowLeft", false, false, false), +1),
		},
		{
			Name:    "gap-plain-left-at-doc-start",
			Content: fixtureProse,
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0)}},
				{Label: "left clamps at start (N=1)", Keys: []Key{{Name: "ArrowLeft"}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "left clamps (N=3)", Keys: []Key{{Name: "ArrowLeft"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "left clamps (N=7)", Keys: []Key{{Name: "ArrowLeft"}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "gap-plain-right-at-doc-end",
			Content: fixtureProse,
			Steps: []Step{
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(fixtureProseLen)}},
				{Label: "right clamps at end (N=1)", Keys: []Key{{Name: "ArrowRight"}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(fixtureProseLen), SelStart: intp(fixtureProseLen), SelEnd: intp(fixtureProseLen)}},
				{Label: "right clamps (N=3)", Keys: []Key{{Name: "ArrowRight"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(fixtureProseLen), SelStart: intp(fixtureProseLen), SelEnd: intp(fixtureProseLen)}},
				{Label: "right clamps (N=7)", Keys: []Key{{Name: "ArrowRight"}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(fixtureProseLen), SelStart: intp(fixtureProseLen), SelEnd: intp(fixtureProseLen)}},
			},
		},
		{
			Name:    "gap-collapse-selection-left",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHMid)},
				{Label: "shift+left x3", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 3},
				{Expect: &StateExpect{Cursor: intp(proseHMid), SelStart: intp(proseHMid - 3), SelEnd: intp(proseHMid), SelLen: intp(3)}},
				{Label: "left collapses selection", Keys: []Key{{Name: "ArrowLeft"}}},
				{Expect: &StateExpect{Cursor: intp(proseHMid - 3), SelStart: intp(proseHMid - 3), SelEnd: intp(proseHMid - 3), SelLen: intp(0)}},
			},
		},
		{
			Name:    "gap-collapse-selection-right",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHMid)},
				{Label: "shift+right x3", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 3},
				{Expect: &StateExpect{Cursor: intp(proseHMid + 3), SelStart: intp(proseHMid), SelEnd: intp(proseHMid + 3), SelLen: intp(3)}},
				{Label: "right collapses selection", Keys: []Key{{Name: "ArrowRight"}}},
				{Expect: &StateExpect{Cursor: intp(proseHMid + 3), SelStart: intp(proseHMid + 3), SelEnd: intp(proseHMid + 3), SelLen: intp(0)}},
			},
		},
		{
			Name:    "gap-delete-forward",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHMid)},
				{Expect: &StateExpect{Cursor: intp(proseHMid), TextLen: intp(fixtureProseLen)}},
				{Label: "delete x1 (N=1)", Keys: []Key{{Name: "Delete"}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseHMid), TextLen: intp(fixtureProseLen - 1)}},
				{Label: "delete x2 more (N=3)", Keys: []Key{{Name: "Delete"}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseHMid), TextLen: intp(fixtureProseLen - 3)}},
				{Label: "delete x4 more (N=7)", Keys: []Key{{Name: "Delete"}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseHMid), TextLen: intp(fixtureProseLen - 7)}},
			},
		},
		{
			Name:    "gap-delete-with-selection",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHStart + 2)},
				{Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{SelStart: intp(proseHStart + 2), SelEnd: intp(proseHStart + 6), SelLen: intp(4)}},
				{Label: "delete selection", Keys: []Key{{Name: "Delete"}}},
				{Expect: &StateExpect{Cursor: intp(proseHStart + 2), TextLen: intp(fixtureProseLen - 4), SelStart: intp(proseHStart + 2), SelEnd: intp(proseHStart + 2)}},
			},
		},
		{
			Name:    "gap-select-all",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseMidDocCaret)},
				{Label: "ctrl+a", Keys: []Key{{Name: "a", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(fixtureProseLen), SelStart: intp(0), SelEnd: intp(fixtureProseLen), SelLen: intp(fixtureProseLen), TextLen: intp(fixtureProseLen)}},
			},
		},
		{
			// In-editor clipboard over the phone/WebSocket path (Bluetooth).
			Name:    "gap-copy-paste",
			Content: "hello world",
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 5},
				{Expect: &StateExpect{SelStart: intp(0), SelEnd: intp(5), SelLen: intp(5), TextLen: intp(11)}},
				{Label: "ctrl+c copy", Keys: []Key{{Name: "c", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(11), Text: strp("hello world")}},
				{Keys: []Key{{Name: "End"}}},
				{Label: "ctrl+v paste", Keys: []Key{{Name: "v", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(16), Text: strp("hello worldhello"), Cursor: intp(16), SelStart: intp(16), SelEnd: intp(16)}},
			},
		},
		{
			Name:    "gap-cut-paste",
			Content: "hello world",
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 5},
				{Expect: &StateExpect{SelStart: intp(0), SelEnd: intp(5), SelLen: intp(5), TextLen: intp(11)}},
				{Label: "ctrl+x cut", Keys: []Key{{Name: "x", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(6), Text: strp(" world"), Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Keys: []Key{{Name: "End"}}},
				{Label: "ctrl+v paste", Keys: []Key{{Name: "v", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(11), Text: strp(" worldhello"), Cursor: intp(11), SelStart: intp(11), SelEnd: intp(11)}},
			},
		},
		{
			Name:    "gap-enter-new-line",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHEditorEnd)},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd), TextLen: intp(fixtureProseLen)}},
				{Label: "enter inserts newline", Keys: []Key{{Name: "Enter"}}},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd + 1), TextLen: intp(fixtureProseLen + 1)}},
			},
		},
		{
			Name:    "gap-type-replaces-selection",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHStart + 4)},
				{Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 3},
				{Expect: &StateExpect{SelStart: intp(proseHStart + 4), SelEnd: intp(proseHStart + 7), SelLen: intp(3)}},
				{Label: "type replaces selection", Keys: []Key{{Name: "x"}}},
				{Expect: &StateExpect{Cursor: intp(proseHStart + 5), TextLen: intp(fixtureProseLen - 2), SelStart: intp(proseHStart + 5), SelEnd: intp(proseHStart + 5)}},
			},
		},
		{
			Name:    "gap-redo-shift-ctrl-z",
			Content: "abc æøå",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "a", Ctrl: true}}},
				{Keys: []Key{{Name: "Backspace"}}},
				{Expect: &StateExpect{TextLen: intp(0), Cursor: intp(0)}},
				{Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(editorLen("abc æøå")), Cursor: intp(editorLen("abc æøå"))}},
				{Label: "shift+ctrl+z redo", Keys: []Key{{Name: "z", Ctrl: true, Shift: true}}},
				{Expect: &StateExpect{TextLen: intp(0), Cursor: intp(0)}},
			},
		},
		{
			Name:    "gap-undo-chain",
			Content: "abc æøå",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "Backspace"}}, Repeat: 2},
				{Expect: &StateExpect{TextLen: intp(editorLen("abc æøå") - 2), Cursor: intp(editorLen("abc æøå") - 2)}},
				{Label: "undo restore one", Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(editorLen("abc æøå") - 1)}},
				{Label: "undo restore two", Keys: []Key{{Name: "z", Ctrl: true}}},
				{Expect: &StateExpect{TextLen: intp(editorLen("abc æøå"))}},
			},
		},
		{
			Name:    "gap-unicode-alt-backspace",
			Content: resume,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(resumeLen), TextLen: intp(resumeLen)}},
				{Label: "alt+backspace unicode word", Keys: []Key{{Name: "Backspace", Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(resumeWordEnd), TextLen: intp(resumeWordEnd), Text: strp("test résumé")}},
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
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseWEditorEnd)},
				{Keys: []Key{{Name: "ArrowLeft", Shift: true, Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(proseWEditorEnd), SelStart: intp(proseWordStarts[11]), SelEnd: intp(proseWEditorEnd), SelLen: intp(proseWEditorEnd - proseWordStarts[11])}},
				{Label: "alt+backspace deletes selection", Keys: []Key{{Name: "Backspace", Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(proseWordStarts[11]), TextLen: intp(fixtureProseLen - (proseWEditorEnd - proseWordStarts[11])), SelStart: intp(proseWordStarts[11]), SelEnd: intp(proseWordStarts[11])}},
			},
		},
		{
			// Mid-sentence Shift+Down/Up across wrapping paragraphs (owner report).
			// Pattern: uni1/uni5 + bi1+1/bi3+5/bi7+7 both directions.
			Name:    "gap-shift-down-mid-wrapping-paras",
			Content: fixtureProse,
			Steps: shiftVisualProsePattern(proseMidDocCaret,
				keyArrow("ArrowDown", true, false, false),
				keyArrow("ArrowUp", true, false, false), false),
		},
		{
			Name:    "gap-shift-up-mid-wrapping-paras",
			Content: fixtureProse,
			Steps: shiftVisualProsePattern(prosePara2Mid,
				keyArrow("ArrowUp", true, false, false),
				keyArrow("ArrowDown", true, false, false), true),
		},
		{
			// Near end of wrapping para1: Shift+Down must enter para2 (not jump oddly).
			Name:    "gap-shift-down-across-para-break",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(prosePara1NearEnd)},
				{Label: "shift+down across blank into para2", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 5},
				{Expect: &StateExpect{
					SelStart:  intp(prosePara1NearEnd),
					CursorMin: intp(prosePara2Start),
					SelLenMin: intp(prosePara2Start - prosePara1NearEnd),
				}},
				{Label: "shrink back with shift+up x5", Keys: []Key{{Name: "ArrowUp", Shift: true}}, Repeat: 5},
				{Expect: &StateExpect{SelLenMax: intp(90), CursorMin: intp(prosePara1NearEnd - 90), CursorMax: intp(prosePara1NearEnd + 90)}},
			},
		},
		{
			// Early in wrapping para2: Shift+Up must reach back into para1.
			Name:    "gap-shift-up-across-para-break",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(prosePara2NearStart)},
				{Label: "shift+up across blank into para1", Keys: []Key{{Name: "ArrowUp", Shift: true}}, Repeat: 5},
				{Expect: &StateExpect{
					SelEnd:    intp(prosePara2NearStart),
					SelLenMin: intp(prosePara2NearStart - prosePara2Start + 2),
				}},
				{Label: "shrink back with shift+down x5", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 5},
				{Expect: &StateExpect{SelLenMax: intp(90), CursorMin: intp(prosePara2NearStart - 90), CursorMax: intp(prosePara2NearStart + 90)}},
			},
		},
		{
			// Mid-column Shift+Down/Up on short hard-newline lines (col 0 already critical).
			Name:    "gap-shift-down-mid-short-lines",
			Content: fixtureProse,
			Steps: shiftVerticalMidPattern(4, proseVWidth/2,
				keyArrow("ArrowDown", true, false, false),
				keyArrow("ArrowUp", true, false, false), false),
		},
	}
}
