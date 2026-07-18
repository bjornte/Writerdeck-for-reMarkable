package main

// comboScenarios cover Alt/Ctrl and Shift+Alt/Ctrl arrow bindings (Mac-style).
// Word jumps use the word line inside fixtureProse at N=1/3/7 both directions.
// Doc-bound jumps prove idempotent clamp at N=1/3/7.
func comboScenarios() []Scenario {
	twoLines := fixtureTwoLines
	twoLen := editorLen(twoLines)
	three := fixtureThreeLines
	threeLen := editorLen(three)
	twoParas := fixtureTwoParas
	twoParasLen := editorLen(twoParas)
	hello := fixtureHelloWorld
	helloLen := editorLen(hello)

	return []Scenario{
		{
			Name:    "combo-alt-left",
			Content: fixtureProse,
			// Seed word index 7; Left needs room for uni5/bi7 and Right overshoot.
			Steps: wordCaretPattern(7, keyArrow("ArrowLeft", false, true, false), keyArrow("ArrowRight", false, true, false), -1),
		},
		{
			Name:    "combo-alt-right",
			Content: fixtureProse,
			Steps:   wordCaretPattern(4, keyArrow("ArrowRight", false, true, false), keyArrow("ArrowLeft", false, true, false), +1),
		},
		{
			Name:    "combo-alt-up",
			Content: twoParas,
			// Mac Option+Up: end of para2 → start of para2 → doc start.
			// Ctrl+End: plain End is line-end only.
			Steps: []Step{
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Label: "alt+up to para2 start (N=1)", Keys: []Key{{Name: "ArrowUp", Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(7), SelStart: intp(7), SelEnd: intp(7)}},
				{Label: "alt+up to doc start (N=3)", Keys: []Key{{Name: "ArrowUp", Alt: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "alt+up clamp (N=7)", Keys: []Key{{Name: "ArrowUp", Alt: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "combo-alt-down",
			Content: twoParas,
			// Mac Option+Down: start → end of para1 → EOF.
			Steps: []Step{
				{Label: "alt+down to para1 end (N=1)", Keys: []Key{{Name: "ArrowDown", Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(5), SelStart: intp(5), SelEnd: intp(5)}},
				{Label: "alt+down further / clamp (N=3)", Keys: []Key{{Name: "ArrowDown", Alt: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(twoParasLen), SelStart: intp(twoParasLen), SelEnd: intp(twoParasLen)}},
				{Label: "alt+down clamp (N=7)", Keys: []Key{{Name: "ArrowDown", Alt: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(twoParasLen), SelStart: intp(twoParasLen), SelEnd: intp(twoParasLen)}},
			},
		},
		{
			// Two consecutive blank lines between paras. Apple paragraphs include
			// empty \n segments — Alt+Up must not skip the blank band to doc start.
			Name:    "combo-alt-up-double-blank",
			Content: fixtureDoubleBlankParas,
			Steps: []Step{
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Label: "alt+up to para2 start", Keys: []Key{{Name: "ArrowUp", Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(8), SelStart: intp(8), SelEnd: intp(8)}},
				{Label: "alt+up into blank paragraphs (not doc start)", Keys: []Key{{Name: "ArrowUp", Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{CursorMin: intp(5), CursorMax: intp(7)}},
				{Label: "alt+up further reaches doc start", Keys: []Key{{Name: "ArrowUp", Alt: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "combo-alt-down-double-blank",
			Content: fixtureDoubleBlankParas,
			Steps: []Step{
				{Label: "alt+down to para1 end", Keys: []Key{{Name: "ArrowDown", Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(5), SelStart: intp(5), SelEnd: intp(5)}},
				{Label: "alt+down into blank paragraphs (not EOF)", Keys: []Key{{Name: "ArrowDown", Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{CursorMin: intp(6), CursorMax: intp(8)}},
				{Label: "alt+down further reaches EOF", Keys: []Key{{Name: "ArrowDown", Alt: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(13), SelStart: intp(13), SelEnd: intp(13)}},
			},
		},
		{
			// Alt+Up across the prose fixture's trailing double-blank section.
			Name:    "combo-alt-up-prose-double-blank",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseDoubleBlank + 10)},
				{Label: "alt+up to double-blank section start", Keys: []Key{{Name: "ArrowUp", Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseDoubleBlank), SelStart: intp(proseDoubleBlank), SelEnd: intp(proseDoubleBlank)}},
				{Label: "alt+up stays in or just above blank band", Keys: []Key{{Name: "ArrowUp", Alt: true}}, Repeat: 1},
				// Must not jump from section start all the way past the blanks in one press.
				{Expect: &StateExpect{CursorMin: intp(proseDoubleBlank - 3), CursorMax: intp(proseDoubleBlank)}},
			},
		},
		{
			// Hard-\n only: visual line == logical. Wrap proof is wrap-ctrl-*.
			Name:    "combo-ctrl-left",
			Content: three,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(5)}}, // end of "to"
				{Label: "ctrl+left to line start (N=1)", Keys: []Key{{Name: "ArrowLeft", Ctrl: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(3), SelStart: intp(3), SelEnd: intp(3)}},
				{Label: "ctrl+left stays (N=3)", Keys: []Key{{Name: "ArrowLeft", Ctrl: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(3), SelStart: intp(3), SelEnd: intp(3)}},
				{Label: "ctrl+left stays (N=7)", Keys: []Key{{Name: "ArrowLeft", Ctrl: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(3), SelStart: intp(3), SelEnd: intp(3)}},
			},
		},
		{
			// Hard-\n only: visual line == logical. Wrap proof is wrap-ctrl-*.
			Name:    "combo-ctrl-right",
			Content: three,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "Home"}}},
				{Expect: &StateExpect{Cursor: intp(3)}}, // start of "to"
				{Label: "ctrl+right to line end (N=1)", Keys: []Key{{Name: "ArrowRight", Ctrl: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(5), SelStart: intp(5), SelEnd: intp(5)}},
				{Label: "ctrl+right stays (N=3)", Keys: []Key{{Name: "ArrowRight", Ctrl: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(5), SelStart: intp(5), SelEnd: intp(5)}},
				{Label: "ctrl+right stays (N=7)", Keys: []Key{{Name: "ArrowRight", Ctrl: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(5), SelStart: intp(5), SelEnd: intp(5)}},
			},
		},
		{
			Name:    "combo-ctrl-up",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseNearEOFCaret)},
				{Label: "ctrl+up to doc start (N=1)", Keys: []Key{{Name: "ArrowUp", Ctrl: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "ctrl+up stays (N=3)", Keys: []Key{{Name: "ArrowUp", Ctrl: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "ctrl+up stays (N=7)", Keys: []Key{{Name: "ArrowUp", Ctrl: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "combo-ctrl-down",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(prosePara2Start)},
				{Label: "ctrl+down to doc end (N=1)", Keys: []Key{{Name: "ArrowDown", Ctrl: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(fixtureProseLen), SelStart: intp(fixtureProseLen), SelEnd: intp(fixtureProseLen)}},
				{Label: "ctrl+down stays (N=3)", Keys: []Key{{Name: "ArrowDown", Ctrl: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(fixtureProseLen), SelStart: intp(fixtureProseLen), SelEnd: intp(fixtureProseLen)}},
				{Label: "ctrl+down stays (N=7)", Keys: []Key{{Name: "ArrowDown", Ctrl: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(fixtureProseLen), SelStart: intp(fixtureProseLen), SelEnd: intp(fixtureProseLen)}},
			},
		},
		{
			Name:    "combo-shift-alt-left",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseWEditorEnd)},
				// 12 words (alfa…lima); from end: N=1→lima[11], N=3→juliett[9], N=7→foxtrot[5].
				{Label: "N=1", Keys: []Key{{Name: "ArrowLeft", Shift: true, Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseWEditorEnd), SelStart: intp(proseWordStarts[11]), SelEnd: intp(proseWEditorEnd), SelLen: intp(proseWEditorEnd - proseWordStarts[11])}},
				{Label: "N=3", Keys: []Key{{Name: "ArrowLeft", Shift: true, Alt: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseWEditorEnd), SelStart: intp(proseWordStarts[9]), SelEnd: intp(proseWEditorEnd)}},
				{Label: "N=7", Keys: []Key{{Name: "ArrowLeft", Shift: true, Alt: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseWEditorEnd), SelStart: intp(proseWordStarts[5]), SelEnd: intp(proseWEditorEnd)}},
			},
		},
		{
			// Alias kept for catalog continuity; same coverage as combo-shift-alt-left.
			Name:    "combo-shift-alt-left-repeat",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseWEditorEnd)},
				{Label: "N=1", Keys: []Key{{Name: "ArrowLeft", Shift: true, Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseWEditorEnd), SelStart: intp(proseWordStarts[11]), SelEnd: intp(proseWEditorEnd)}},
				{Label: "N=3", Keys: []Key{{Name: "ArrowLeft", Shift: true, Alt: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseWEditorEnd), SelStart: intp(proseWordStarts[9]), SelEnd: intp(proseWEditorEnd)}},
				{Label: "N=7", Keys: []Key{{Name: "ArrowLeft", Shift: true, Alt: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseWEditorEnd), SelStart: intp(proseWordStarts[5]), SelEnd: intp(proseWEditorEnd)}},
			},
		},
		{
			// After Shift-select + typing, Shift+Alt+Left must re-anchor at the
			// post-type caret (not a leftover shiftHead from before the edit).
			Name:    "combo-shift-alt-left-after-type",
			Content: "hello world test",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Label: "shift+left seeds stale head", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{SelLen: intp(4), TextLen: intp(16)}},
				{Label: "type replaces selection", Keys: []Key{{Name: "x"}}},
				{Expect: &StateExpect{Text: strp("hello world x"), TextLen: intp(13), Cursor: intp(13), SelLen: intp(0)}},
				{Label: "shift+alt+left from post-type caret", Keys: []Key{{Name: "ArrowLeft", Shift: true, Alt: true}}},
				{Expect: &StateExpect{Cursor: intp(13), SelStart: intp(12), SelEnd: intp(13), SelLen: intp(1)}},
			},
		},
		{
			// Same stale-head risk on plain Shift+Left (different code path than
			// Shift+Alt / selectionExtendFrom). One type-then-nav guard is enough
			// for this axis; not every combo needs an after-type twin.
			Name:    "combo-shift-left-after-type",
			Content: "hello world test",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Label: "shift+left seeds stale head", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{SelLen: intp(4), TextLen: intp(16)}},
				{Label: "type replaces selection", Keys: []Key{{Name: "x"}}},
				{Expect: &StateExpect{Text: strp("hello world x"), TextLen: intp(13), Cursor: intp(13), SelLen: intp(0)}},
				{Label: "shift+left from post-type caret", Keys: []Key{{Name: "ArrowLeft", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(13), SelStart: intp(12), SelEnd: intp(13), SelLen: intp(1)}},
			},
		},
		{
			Name:    "combo-shift-alt-right",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseWStart)},
				{Label: "N=1", Keys: []Key{{Name: "ArrowRight", Shift: true, Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseWordStarts[1]), SelStart: intp(proseWStart), SelEnd: intp(proseWordStarts[1])}},
				{Label: "N=3", Keys: []Key{{Name: "ArrowRight", Shift: true, Alt: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseWordStarts[3]), SelStart: intp(proseWStart), SelEnd: intp(proseWordStarts[3])}},
				{Label: "N=7", Keys: []Key{{Name: "ArrowRight", Shift: true, Alt: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseWordStarts[7]), SelStart: intp(proseWStart), SelEnd: intp(proseWordStarts[7])}},
			},
		},
		{
			Name:    "combo-shift-alt-right-repeat",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseWStart)},
				{Label: "N=1", Keys: []Key{{Name: "ArrowRight", Shift: true, Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseWordStarts[1]), SelStart: intp(proseWStart), SelEnd: intp(proseWordStarts[1])}},
				{Label: "N=3", Keys: []Key{{Name: "ArrowRight", Shift: true, Alt: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseWordStarts[3]), SelStart: intp(proseWStart), SelEnd: intp(proseWordStarts[3])}},
				{Label: "N=7", Keys: []Key{{Name: "ArrowRight", Shift: true, Alt: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseWordStarts[7]), SelStart: intp(proseWStart), SelEnd: intp(proseWordStarts[7])}},
			},
		},
		{
			Name:    "combo-shift-alt-up",
			Content: twoParas,
			Steps: []Step{
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowUp", Shift: true, Alt: true}}},
				// From EOF: select back to start of current paragraph (para2).
				{Expect: &StateExpect{Cursor: intp(twoParasLen), SelStart: intp(7), SelEnd: intp(twoParasLen), SelLen: intp(twoParasLen - 7)}},
			},
		},
		{
			Name:    "combo-shift-alt-down",
			Content: twoParas,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown", Shift: true, Alt: true}}},
				// From doc start: select to end of current paragraph (para1).
				{Expect: &StateExpect{Cursor: intp(5), SelStart: intp(0), SelEnd: intp(5), SelLen: intp(5)}},
			},
		},
		{
			Name:    "combo-shift-ctrl-left",
			Content: hello,
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Label: "N=1", Keys: []Key{{Name: "ArrowLeft", Shift: true, Ctrl: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(helloLen), SelStart: intp(0), SelEnd: intp(helloLen), SelLen: intp(helloLen)}},
				{Label: "N=3", Keys: []Key{{Name: "ArrowLeft", Shift: true, Ctrl: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(helloLen), SelStart: intp(0), SelEnd: intp(helloLen), SelLen: intp(helloLen)}},
				{Label: "N=7", Keys: []Key{{Name: "ArrowLeft", Shift: true, Ctrl: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(helloLen), SelStart: intp(0), SelEnd: intp(helloLen), SelLen: intp(helloLen)}},
			},
		},
		{
			Name:    "combo-shift-ctrl-left-multiline",
			Content: three,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "ArrowLeft", Shift: true, Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(5), SelStart: intp(3), SelEnd: intp(5), SelLen: intp(2)}},
			},
		},
		{
			Name:    "combo-shift-ctrl-right",
			Content: hello,
			Steps: []Step{
				{Label: "N=1", Keys: []Key{{Name: "ArrowRight", Shift: true, Ctrl: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(helloLen), SelStart: intp(0), SelEnd: intp(helloLen), SelLen: intp(helloLen)}},
				{Label: "N=3", Keys: []Key{{Name: "ArrowRight", Shift: true, Ctrl: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(helloLen), SelStart: intp(0), SelEnd: intp(helloLen), SelLen: intp(helloLen)}},
				{Label: "N=7", Keys: []Key{{Name: "ArrowRight", Shift: true, Ctrl: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(helloLen), SelStart: intp(0), SelEnd: intp(helloLen), SelLen: intp(helloLen)}},
			},
		},
		{
			Name:    "combo-shift-ctrl-up",
			Content: three,
			Steps: []Step{
				{Keys: []Key{{Name: "End", Ctrl: true}}},
				{Label: "N=1", Keys: []Key{{Name: "ArrowUp", Shift: true, Ctrl: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(threeLen), SelStart: intp(0), SelEnd: intp(threeLen), SelLen: intp(threeLen)}},
				{Label: "N=3", Keys: []Key{{Name: "ArrowUp", Shift: true, Ctrl: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(threeLen), SelStart: intp(0), SelEnd: intp(threeLen), SelLen: intp(threeLen)}},
				{Label: "N=7", Keys: []Key{{Name: "ArrowUp", Shift: true, Ctrl: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(threeLen), SelStart: intp(0), SelEnd: intp(threeLen), SelLen: intp(threeLen)}},
			},
		},
		{
			Name:    "combo-shift-ctrl-down",
			Content: three,
			Steps: []Step{
				{Label: "N=1", Keys: []Key{{Name: "ArrowDown", Shift: true, Ctrl: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(threeLen), SelStart: intp(0), SelEnd: intp(threeLen), SelLen: intp(threeLen)}},
				{Label: "N=3", Keys: []Key{{Name: "ArrowDown", Shift: true, Ctrl: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(threeLen), SelStart: intp(0), SelEnd: intp(threeLen), SelLen: intp(threeLen)}},
				{Label: "N=7", Keys: []Key{{Name: "ArrowDown", Shift: true, Ctrl: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(threeLen), SelStart: intp(0), SelEnd: intp(threeLen), SelLen: intp(threeLen)}},
			},
		},
		{
			Name:    "combo-shift-home-line",
			Content: twoLines,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "End"}}},
				{Keys: []Key{{Name: "Home", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(twoLen), SelStart: intp(4), SelEnd: intp(twoLen), SelLen: intp(3)}},
			},
		},
		{
			Name:    "combo-shift-end-line",
			Content: twoLines,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Keys: []Key{{Name: "Home"}}},
				{Keys: []Key{{Name: "End", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(twoLen), SelStart: intp(4), SelEnd: intp(twoLen), SelLen: intp(3)}},
			},
		},
		{
			Name:    "combo-ctrl-home",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseVLineStart(5))},
				{Label: "N=1", Keys: []Key{{Name: "Home", Ctrl: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "N=3", Keys: []Key{{Name: "Home", Ctrl: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "N=7", Keys: []Key{{Name: "Home", Ctrl: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "combo-ctrl-end",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseList2Item2)},
				{Label: "N=1", Keys: []Key{{Name: "End", Ctrl: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(fixtureProseLen), SelStart: intp(fixtureProseLen), SelEnd: intp(fixtureProseLen)}},
				{Label: "N=3", Keys: []Key{{Name: "End", Ctrl: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(fixtureProseLen), SelStart: intp(fixtureProseLen), SelEnd: intp(fixtureProseLen)}},
				{Label: "N=7", Keys: []Key{{Name: "End", Ctrl: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(fixtureProseLen), SelStart: intp(fixtureProseLen), SelEnd: intp(fixtureProseLen)}},
			},
		},
		{
			Name:    "combo-shift-ctrl-home",
			Content: twoLines,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(4)}},
				{Keys: []Key{{Name: "Home", Shift: true, Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(4), SelStart: intp(0), SelEnd: intp(4), SelLen: intp(4)}},
			},
		},
		{
			Name:    "combo-shift-ctrl-end",
			Content: twoLines,
			Steps: []Step{
				{Keys: []Key{{Name: "ArrowDown"}}},
				{Expect: &StateExpect{Cursor: intp(4)}},
				{Keys: []Key{{Name: "End", Shift: true, Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(twoLen), SelStart: intp(4), SelEnd: intp(twoLen), SelLen: intp(3)}},
			},
		},
	}
}
