package main

// hwScenarios cover rM1 physical page-turn buttons (not keyboard arrows).
// Harness injects pageleft/pageright editor cmds — the same path the daemon
// should use after exclusive gpio grab on /dev/input/event1.
//
// Page step in QML is ~1500px (build-keywriter scrollUp/Down). Content must be
// tall enough that N=1/3/7 page-downs leave contentY well above the viewport.
func hwScenarios() []Scenario {
	tall := fixtureTall(250)

	return []Scenario{
		{
			Name:    "hw-page-right-scrolls-edit",
			Content: tall,
			Steps: []Step{
				{Expect: &StateExpect{Cursor: intp(0), ContentY: intp(0), Mode: intp(1)}},
				{Label: "page right x1 (N=1)", Cmd: "pageright", Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(0), ContentYMin: intp(1400), ContentYMax: intp(1600)}},
				{Label: "page right x2 more (N=3)", Cmd: "pageright", Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(0), ContentYMin: intp(4200), ContentYMax: intp(4800)}},
				{Label: "page right x4 more (N=7)", Cmd: "pageright", Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(0), ContentYMin: intp(10000), ContentYMax: intp(11000)}},
			},
		},
		{
			Name:    "hw-page-left-scrolls-edit",
			Content: tall,
			Steps: []Step{
				{Label: "seed page right x7", Cmd: "pageright", Repeat: 7},
				{Expect: &StateExpect{Cursor: intp(0), ContentYMin: intp(10000)}},
				{Label: "page left x1 (N=1)", Cmd: "pageleft", Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(0), ContentYMin: intp(8500), ContentYMax: intp(9500)}},
				{Label: "page left x2 more (N=3)", Cmd: "pageleft", Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(0), ContentYMin: intp(5500), ContentYMax: intp(6500)}},
				{Label: "page left x4 more (N=7)", Cmd: "pageleft", Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(0), ContentY: intp(0)}},
			},
		},
	}
}
