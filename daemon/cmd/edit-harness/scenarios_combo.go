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
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Label: "alt+up x1 (N=1)", Keys: []Key{{Name: "ArrowUp", Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "alt+up clamp (N=3)", Keys: []Key{{Name: "ArrowUp", Alt: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "alt+up clamp (N=7)", Keys: []Key{{Name: "ArrowUp", Alt: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "combo-alt-down",
			Content: twoParas,
			Steps: []Step{
				{Label: "alt+down x1 (N=1)", Keys: []Key{{Name: "ArrowDown", Alt: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(7), SelStart: intp(7), SelEnd: intp(7)}},
				{Label: "alt+down further / clamp (N=3)", Keys: []Key{{Name: "ArrowDown", Alt: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(twoParasLen), SelStart: intp(twoParasLen), SelEnd: intp(twoParasLen)}},
				{Label: "alt+down clamp (N=7)", Keys: []Key{{Name: "ArrowDown", Alt: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(twoParasLen), SelStart: intp(twoParasLen), SelEnd: intp(twoParasLen)}},
			},
		},
		{
			Name:    "combo-ctrl-left",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseMidDocCaret)},
				{Label: "ctrl+left to doc start (N=1)", Keys: []Key{{Name: "ArrowLeft", Ctrl: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "ctrl+left stays (N=3)", Keys: []Key{{Name: "ArrowLeft", Ctrl: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "ctrl+left stays (N=7)", Keys: []Key{{Name: "ArrowLeft", Ctrl: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "combo-ctrl-right",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseList1Start)},
				{Label: "ctrl+right to doc end (N=1)", Keys: []Key{{Name: "ArrowRight", Ctrl: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(fixtureProseLen), SelStart: intp(fixtureProseLen), SelEnd: intp(fixtureProseLen)}},
				{Label: "ctrl+right stays (N=3)", Keys: []Key{{Name: "ArrowRight", Ctrl: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(fixtureProseLen), SelStart: intp(fixtureProseLen), SelEnd: intp(fixtureProseLen)}},
				{Label: "ctrl+right stays (N=7)", Keys: []Key{{Name: "ArrowRight", Ctrl: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(fixtureProseLen), SelStart: intp(fixtureProseLen), SelEnd: intp(fixtureProseLen)}},
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
				{Expect: &StateExpect{Cursor: intp(twoParasLen), SelStart: intp(7), SelEnd: intp(twoParasLen), SelLen: intp(twoParasLen - 7)}},
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
