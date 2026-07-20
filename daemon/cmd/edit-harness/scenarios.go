package main

import "strings"

// coreScenarios: load/home plus Shift extend with uni1/uni5/bi1+1/bi3+5/bi7+7.
func coreScenarios() []Scenario {
	left := keyArrow("ArrowLeft", true, false, false)
	right := keyArrow("ArrowRight", true, false, false)
	up := keyArrow("ArrowUp", true, false, false)
	down := keyArrow("ArrowDown", true, false, false)

	return []Scenario{
		{
			Name:    "load-cursor-at-start",
			Content: fixtureProse,
			Steps: []Step{
				{
					Label: "after open",
					Expect: &StateExpect{
						Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0),
						TextLen: intp(fixtureProseLen), Mode: intp(1),
					},
				},
			},
		},
		{
			Name:    "home-clears-selection",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHMid)},
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(proseHEditorEnd), SelStart: intp(proseHEditorEnd), SelEnd: intp(proseHEditorEnd)}},
				{Keys: []Key{{Name: "Home", Shift: true}}},
				{Expect: &StateExpect{SelStart: intp(proseHStart), SelEnd: intp(proseHEditorEnd), SelLen: intp(proseHLen)}},
				{Keys: []Key{{Name: "Home"}}},
				{Expect: &StateExpect{Cursor: intp(proseHStart), SelStart: intp(proseHStart), SelEnd: intp(proseHStart)}},
			},
		},
		{
			// Mid-line (not column 0) so bi 3+5 overshoot Left has room.
			Name:    "shift-right-from-home",
			Content: fixtureProse,
			Steps:   shiftSelectAxisPattern(proseHStart+10, right, left, false),
		},
		{
			// Near line end so bi 3+5 overshoot Right has room past the grow.
			Name:    "shift-left-from-end",
			Content: fixtureProse,
			Steps:   shiftSelectAxisPattern(proseHEditorEnd-10, left, right, true),
		},
		{
			Name:    "shift-right-after-home-no-stale-anchor",
			Content: fixtureProse,
			Steps:   shiftSelectAxisPattern(proseHStart+10, right, left, false),
		},
		{
			Name:    "shift-left-after-end-no-stale-anchor",
			Content: fixtureProse,
			Steps:   shiftSelectAxisPattern(proseHEditorEnd-10, left, right, true),
		},
		{
			Name:    "shift-down-after-arrow-down",
			Content: fixtureProse,
			Steps:   shiftVerticalPattern(2, down, up, false),
		},
		{
			Name:    "shift-up-after-arrow-down",
			Content: fixtureProse,
			Steps:   shiftVerticalPattern(9, up, down, true),
		},
		{
			Name:    "ctrl-shift-left-select-line",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHEditorEnd)},
				{Keys: []Key{{Name: "Home", Shift: true}}},
				{Expect: &StateExpect{SelStart: intp(proseHStart), SelEnd: intp(proseHEditorEnd), SelLen: intp(proseHLen)}},
			},
		},
	}
}

// AllScenarios returns keyboard/selection integration scenarios.
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
