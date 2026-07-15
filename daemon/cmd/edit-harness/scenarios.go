package main

import "strings"

// AllScenarios returns keyboard/selection integration scenarios.
// Each step can send keys and/or assert editor state via /api/test/editor-state.
func AllScenarios() []Scenario {
	out := append([]Scenario{}, coreScenarios()...)
	out = append(out, regressionScenarios()...)
	out = append(out, cmScenarios()...)
	out = append(out, comboScenarios()...)
	out = append(out, bsScenarios()...)
	out = append(out, wrapScenarios()...)
	out = append(out, undoScenarios()...)
	out = append(out, gapScenarios()...)
	out = append(out, hwScenarios()...)
	out = append(out, readScenarios()...)
	out = append(out, touchScenarios()...)
	out = append(out, selectionScenarios()...)
	return attachScenarioTags(out)
}

// coreScenarios cover the basic caret/selection contract on fixtureProse.
// Directional cases prove N=1, N=3, and N=7 in both directions.
func coreScenarios() []Scenario {
	return []Scenario{
		{
			Name:    "load-cursor-at-start",
			Content: fixtureProse,
			Steps: []Step{
				{
					Label: "after open",
					Expect: &StateExpect{
						Cursor:   intp(0),
						SelStart: intp(0),
						SelEnd:   intp(0),
						TextLen:  intp(fixtureProseLen),
						Mode:     intp(1),
					},
				},
			},
		},
		{
			Name:    "home-clears-selection",
			Content: fixtureProse,
			Steps: []Step{
				{Label: "mid horizontal line", SetCursor: intp(proseHMid)},
				{Keys: []Key{{Name: "End"}}},
				{Label: "at line end", Expect: &StateExpect{Cursor: intp(proseHEditorEnd), SelStart: intp(proseHEditorEnd), SelEnd: intp(proseHEditorEnd)}},
				{Keys: []Key{{Name: "Home", Shift: true}}},
				{Expect: &StateExpect{SelStart: intp(proseHStart), SelEnd: intp(proseHEditorEnd), SelLen: intp(proseHLen)}},
				{Label: "Home clears selection", Keys: []Key{{Name: "Home"}}},
				{Expect: &StateExpect{Cursor: intp(proseHStart), SelStart: intp(proseHStart), SelEnd: intp(proseHStart)}},
			},
		},
		{
			// Reverse partner: shift-left-from-end.
			Name:    "shift-right-from-home",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHStart)},
				{Expect: &StateExpect{Cursor: intp(proseHStart), SelStart: intp(proseHStart), SelEnd: intp(proseHStart)}},
				{Label: "shift+right x1 (N=1)", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseHStart + 1), SelStart: intp(proseHStart), SelEnd: intp(proseHStart + 1)}},
				{Label: "shift+right x2 more (N=3)", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseHStart + 3), SelStart: intp(proseHStart), SelEnd: intp(proseHStart + 3)}},
				{Label: "shift+right x4 more (N=7)", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseHStart + 7), SelStart: intp(proseHStart), SelEnd: intp(proseHStart + 7)}},
			},
		},
		{
			// Reverse partner: shift-right-from-home. Mid-line start (not EOF).
			Name:    "shift-left-from-end",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHEditorEnd)},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd), SelStart: intp(proseHEditorEnd), SelEnd: intp(proseHEditorEnd)}},
				{Label: "shift+left x1 (N=1)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd), SelStart: intp(proseHEditorEnd - 1), SelEnd: intp(proseHEditorEnd)}},
				{Label: "shift+left x2 more (N=3)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd), SelStart: intp(proseHEditorEnd - 3), SelEnd: intp(proseHEditorEnd)}},
				{Label: "shift+left x4 more (N=7)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd), SelStart: intp(proseHEditorEnd - 7), SelEnd: intp(proseHEditorEnd)}},
			},
		},
		{
			Name:    "shift-right-after-home-no-stale-anchor",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHStart)},
				{Label: "N=1", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseHStart + 1), SelStart: intp(proseHStart), SelEnd: intp(proseHStart + 1)}},
				{Label: "N=2", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseHStart + 2), SelStart: intp(proseHStart), SelEnd: intp(proseHStart + 2)}},
				{Label: "N=7 (anchor still start)", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 5},
				{Expect: &StateExpect{Cursor: intp(proseHStart + 7), SelStart: intp(proseHStart), SelEnd: intp(proseHStart + 7)}},
			},
		},
		{
			Name:    "shift-left-after-end-no-stale-anchor",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHEditorEnd)},
				{Label: "N=1", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd), SelStart: intp(proseHEditorEnd - 1), SelEnd: intp(proseHEditorEnd)}},
				{Label: "N=2", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd), SelStart: intp(proseHEditorEnd - 2), SelEnd: intp(proseHEditorEnd)}},
				{Label: "N=7 (anchor still end)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 5},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd), SelStart: intp(proseHEditorEnd - 7), SelEnd: intp(proseHEditorEnd)}},
			},
		},
		{
			// Reverse partner: shift-up-after-arrow-down. Mid-doc vertical block.
			Name:    "shift-down-after-arrow-down",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseVLineStart(2))},
				{Label: "cursor on vertical line 2", Expect: &StateExpect{Cursor: intp(proseVLineStart(2)), SelStart: intp(proseVLineStart(2)), SelEnd: intp(proseVLineStart(2)), TextLen: intp(fixtureProseLen)}},
				{Label: "shift+down x1 (N=1)", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(3)), SelStart: intp(proseVLineStart(2)), SelEnd: intp(proseVLineStart(3))}},
				{Label: "shift+down x2 more (N=3)", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(5)), SelStart: intp(proseVLineStart(2)), SelEnd: intp(proseVLineStart(5))}},
				{Label: "shift+down x4 more (N=7)", Keys: []Key{{Name: "ArrowDown", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseVLineStart(9)), SelStart: intp(proseVLineStart(2)), SelEnd: intp(proseVLineStart(9))}},
			},
		},
		{
			// Reverse partner: shift-down-after-arrow-down.
			Name:    "shift-up-after-arrow-down",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseVLineEnd(9))},
				{Label: "cursor on vertical line 9 end", Expect: &StateExpect{Cursor: intp(proseVLineEnd(9)), SelStart: intp(proseVLineEnd(9)), SelEnd: intp(proseVLineEnd(9)), TextLen: intp(fixtureProseLen)}},
				{Label: "shift+up x1 (N=1)", Keys: []Key{{Name: "ArrowUp", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(9)), SelStart: intp(proseVLineEnd(8)), SelEnd: intp(proseVLineEnd(9))}},
				{Label: "shift+up x2 more (N=3)", Keys: []Key{{Name: "ArrowUp", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(9)), SelStart: intp(proseVLineEnd(6)), SelEnd: intp(proseVLineEnd(9))}},
				{Label: "shift+up x4 more (N=7)", Keys: []Key{{Name: "ArrowUp", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseVLineEnd(9)), SelStart: intp(proseVLineEnd(2)), SelEnd: intp(proseVLineEnd(9))}},
			},
		},
		{
			Name:    "ctrl-shift-left-select-line",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHEditorEnd)},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd)}},
				{Label: "shift+home selects line", Keys: []Key{{Name: "Home", Shift: true}}},
				{Expect: &StateExpect{SelStart: intp(proseHStart), SelEnd: intp(proseHEditorEnd), SelLen: intp(proseHLen)}},
			},
		},
	}
}

func scenarioNames() []string {
	var names []string
	for _, sc := range AllScenarios() {
		names = append(names, sc.Name)
	}
	return names
}

func findScenario(name string) (Scenario, bool) {
	for _, sc := range AllScenarios() {
		if sc.Name == name {
			return sc, true
		}
	}
	return Scenario{}, false
}

func findScenariosByPrefix(substr string) ([]Scenario, bool) {
	var out []Scenario
	for _, sc := range AllScenarios() {
		if strings.Contains(sc.Name, substr) {
			out = append(out, sc)
		}
	}
	return out, len(out) > 0
}
