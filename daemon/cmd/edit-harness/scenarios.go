package main

import "strings"

// AllScenarios returns keyboard/selection integration scenarios.
// Each step can send keys and/or assert editor state via /api/test/editor-state.
func AllScenarios() []Scenario {
	return append(coreScenarios(), regressionScenarios()...)
}

func coreScenarios() []Scenario {
	return []Scenario{
		{
			Name:    "load-cursor-at-start",
			Content: "abcdef",
			Steps: []Step{
				{
					Label: "after open",
					Expect: &StateExpect{
						Cursor:   intp(0),
						SelStart: intp(0),
						SelEnd:   intp(0),
						TextLen:  intp(6),
						Mode:     intp(1),
					},
				},
			},
		},
		{
			Name:    "home-clears-selection",
			Content: "abcdef",
			Steps: []Step{
				{Label: "start", Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "End", Keys: []Key{{Name: "End"}}},
				{Label: "at end", Expect: &StateExpect{Cursor: intp(6), SelStart: intp(6), SelEnd: intp(6)}},
				{Label: "Home", Keys: []Key{{Name: "Home"}}},
				{Label: "cursor line start", Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
			},
		},
		{
			Name:    "shift-right-from-home",
			Content: "abcdef",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Label: "shift+right x3", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 3},
				{Expect: &StateExpect{Cursor: intp(3), SelStart: intp(0), SelEnd: intp(3)}},
			},
		},
		{
			Name:    "shift-left-from-end",
			Content: "abcdef",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(6), SelStart: intp(6), SelEnd: intp(6)}},
				{Label: "shift+left x3", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 3},
				{Expect: &StateExpect{Cursor: intp(6), SelStart: intp(3), SelEnd: intp(6)}},
			},
		},
		{
			Name:    "shift-right-after-home-no-stale-anchor",
			Content: "abcdef",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Expect: &StateExpect{Cursor: intp(0), SelStart: intp(0), SelEnd: intp(0)}},
				{Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(1), SelStart: intp(0), SelEnd: intp(1)}},
				{Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(2), SelStart: intp(0), SelEnd: intp(2)}},
			},
		},
		{
			Name:    "shift-down-after-arrow-down",
			Content: "line1\nline2\nline3\nline4",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Label: "cursor on line3", Expect: &StateExpect{Cursor: intp(12), SelStart: intp(12), SelEnd: intp(12), TextLen: intp(23)}},
				{Label: "shift+down extends downward", Keys: []Key{{Name: "ArrowDown", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(18), SelStart: intp(12), SelEnd: intp(18)}},
			},
		},
		{
			Name:    "shift-up-after-arrow-down",
			Content: "line1\nline2\nline3\nline4",
			Steps: []Step{
				{Keys: []Key{{Name: "Home", Ctrl: true}}},
				{Keys: []Key{{Name: "ArrowDown"}}, Repeat: 2},
				{Label: "cursor on line3", Expect: &StateExpect{Cursor: intp(12), SelStart: intp(12), SelEnd: intp(12)}},
				{Label: "shift+up extends upward", Keys: []Key{{Name: "ArrowUp", Shift: true}}},
				{Expect: &StateExpect{Cursor: intp(12), SelStart: intp(6), SelEnd: intp(12)}},
			},
		},
		{
			Name:    "ctrl-shift-left-select-line",
			Content: "abcdef",
			Steps: []Step{
				{Keys: []Key{{Name: "End"}}},
				{Expect: &StateExpect{Cursor: intp(6)}},
				{Label: "shift+home selects line", Keys: []Key{{Name: "Home", Shift: true}}},
				{Expect: &StateExpect{SelStart: intp(0), SelEnd: intp(6)}},
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
