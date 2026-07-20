package main

// hwScenarios: pageleft/pageright with uni1/uni5/bi1+1/bi3+5/bi7+7.
// Page step is ~85% of Flickable height (see flick.scrollDown), not a fixed
// 1500px — landscape uses the short side and must step less than portrait.
func hwScenarios() []Scenario {
	tall := fixtureTall(600)
	// Force portrait (rotation 0): body height = screen.height ≈ 1872,
	// step ≈ round(1872*0.85) = 1591. (Saved tablet rotation may be 90.)
	const step = 1590
	const slack = 150

	expectY := func(pages int) *StateExpect {
		y := pages * step
		return &StateExpect{
			Cursor: intp(0), Mode: intp(1),
			ContentYMin: intp(y - slack), ContentYMax: intp(y + slack),
		}
	}
	at0 := &StateExpect{Cursor: intp(0), ContentY: intp(0), Mode: intp(1)}
	forcePortrait := Step{Label: "force portrait", Cmd: "setrotation", Degrees: intp(0), PauseMs: 500}

	pageRightPattern := func() []Step {
		var out []Step
		out = append(out, forcePortrait)
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
		out = append(out, forcePortrait)
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
		{
			// Portrait (rot 0, tall) must step more than landscape (rot 90, short).
			// Fixed 1500px treated both as the same height and overshot landscape.
			Name:    "hw-page-step-shrinks-in-landscape",
			Content: tall,
			Steps: []Step{
				{Label: "portrait", Cmd: "setrotation", Degrees: intp(0), PauseMs: 500},
				{Expect: at0},
				{Label: "portrait page", Cmd: "pageright", CaptureContentY: true,
					Expect: expectY(1)},
				{Label: "reset portrait", Cmd: "pageleft"},
				{Expect: at0},
				{Label: "landscape", Cmd: "setrotation", Degrees: intp(90), PauseMs: 500},
				{Label: "landscape page", Cmd: "pageright", ExpectContentYLtCaptured: true,
					Expect: &StateExpect{
						Cursor: intp(0), Mode: intp(1),
						// ~0.85 * 1404 ≈ 1193
						ContentYMin: intp(1000), ContentYMax: intp(1350),
					}},
				{Label: "restore portrait", Cmd: "setrotation", Degrees: intp(0), PauseMs: 300},
			},
		},
	}
}
