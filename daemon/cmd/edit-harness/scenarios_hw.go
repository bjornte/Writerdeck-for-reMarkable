package main

// hwScenarios: pageleft/pageright with uni1/uni5/bi1+1/bi3+5/bi7+7.
// ~1500px per page. Tall content so 7 pages stay below max scroll.
func hwScenarios() []Scenario {
	tall := fixtureTall(400)
	const step = 1500
	const slack = 150

	expectY := func(pages int) *StateExpect {
		y := pages * step
		return &StateExpect{
			Cursor: intp(0), Mode: intp(1),
			ContentYMin: intp(y - slack), ContentYMax: intp(y + slack),
		}
	}
	at0 := &StateExpect{Cursor: intp(0), ContentY: intp(0), Mode: intp(1)}

	pageRightPattern := func() []Step {
		var out []Step
		out = append(out, Step{Expect: at0})
		// uni 1
		out = append(out, Step{Label: "uni 1", Cmd: "pageright", Repeat: 1})
		out = append(out, Step{Expect: expectY(1)})
		// reset via pageleft
		out = append(out, Step{Label: "reset after uni1", Cmd: "pageleft", Repeat: 1})
		out = append(out, Step{Expect: at0})
		// uni 5
		out = append(out, Step{Label: "uni 5", Cmd: "pageright", Repeat: 5})
		out = append(out, Step{Expect: expectY(5)})
		out = append(out, Step{Label: "reset after uni5", Cmd: "pageleft", Repeat: 5})
		out = append(out, Step{Expect: at0})
		// bi 1+1
		out = append(out, Step{Label: "bi 1+1 forward", Cmd: "pageright", Repeat: 1})
		out = append(out, Step{Expect: expectY(1)})
		out = append(out, Step{Label: "bi 1+1 reverse", Cmd: "pageleft", Repeat: 1})
		out = append(out, Step{Expect: at0})
		// bi 3+5 overshoot: after reverse 3 back at 0; 2 more try below 0 → clamp 0
		out = append(out, Step{Label: "bi 3+5 forward x3", Cmd: "pageright", Repeat: 3})
		out = append(out, Step{Expect: expectY(3)})
		out = append(out, Step{Label: "bi 3+5 reverse x5 (overshoot clamp)", Cmd: "pageleft", Repeat: 5})
		out = append(out, Step{Expect: at0})
		// bi 7+7
		out = append(out, Step{Label: "bi 7+7 forward x7", Cmd: "pageright", Repeat: 7})
		out = append(out, Step{Expect: expectY(7)})
		out = append(out, Step{Label: "bi 7+7 reverse x7", Cmd: "pageleft", Repeat: 7})
		out = append(out, Step{Expect: at0})
		return out
	}

	pageLeftPattern := func() []Step {
		var out []Step
		// Seed to page 7 so Left has room; uni1 from there.
		out = append(out, Step{Label: "seed to y7", Cmd: "pageright", Repeat: 7})
		out = append(out, Step{Expect: expectY(7)})
		out = append(out, Step{Label: "uni 1 left", Cmd: "pageleft", Repeat: 1})
		out = append(out, Step{Expect: expectY(6)})
		out = append(out, Step{Label: "reset seed", Cmd: "pageright", Repeat: 1})
		out = append(out, Step{Expect: expectY(7)})
		out = append(out, Step{Label: "uni 5 left", Cmd: "pageleft", Repeat: 5})
		out = append(out, Step{Expect: expectY(2)})
		out = append(out, Step{Label: "reseed y7", Cmd: "pageright", Repeat: 5})
		out = append(out, Step{Expect: expectY(7)})
		// bi 1+1
		out = append(out, Step{Label: "bi 1+1 left", Cmd: "pageleft", Repeat: 1})
		out = append(out, Step{Expect: expectY(6)})
		out = append(out, Step{Label: "bi 1+1 reverse right", Cmd: "pageright", Repeat: 1})
		out = append(out, Step{Expect: expectY(7)})
		// bi 3+5: left 3 → y4; right 5 → y9
		out = append(out, Step{Label: "bi 3+5 left x3", Cmd: "pageleft", Repeat: 3})
		out = append(out, Step{Expect: expectY(4)})
		out = append(out, Step{Label: "bi 3+5 right x5 (overshoot)", Cmd: "pageright", Repeat: 5})
		out = append(out, Step{Expect: expectY(9)})
		// return to y7 then bi 7+7
		out = append(out, Step{Label: "back to y7", Cmd: "pageleft", Repeat: 2})
		out = append(out, Step{Expect: expectY(7)})
		out = append(out, Step{Label: "bi 7+7 left x7", Cmd: "pageleft", Repeat: 7})
		out = append(out, Step{Expect: at0})
		out = append(out, Step{Label: "bi 7+7 right x7", Cmd: "pageright", Repeat: 7})
		out = append(out, Step{Expect: expectY(7)})
		return out
	}

	return []Scenario{
		{Name: "hw-page-right-scrolls-edit", Content: tall, Steps: pageRightPattern()},
		{Name: "hw-page-left-scrolls-edit", Content: tall, Steps: pageLeftPattern()},
	}
}
