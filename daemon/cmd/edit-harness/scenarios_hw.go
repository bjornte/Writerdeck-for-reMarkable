package main

import "strings"

// hwScenarios cover rM1 physical page-turn buttons (not keyboard arrows).
// Harness injects pageleft/pageright editor cmds — the same path the daemon
// should use after exclusive gpio grab on /dev/input/event1.
//
// Page step in QML is ~1500px (build-keywriter scrollUp/Down). Content must be
// tall enough that two page-downs leave contentY well above the viewport.
func hwScenarios() []Scenario {
	// ~250 lines: with ~40–60px line height this is several screens + headroom
	// for two 1500px page steps (viewport ~1872px).
	var b strings.Builder
	for i := 0; i < 250; i++ {
		b.WriteString("Page scroll filler line ")
		b.WriteString(strings.Repeat("x", 40))
		b.WriteByte('\n')
	}
	tall := b.String()

	// Expect near-full 1500px steps; allow small clamp/rounding slack.
	return []Scenario{
		{
			Name:    "hw-page-right-scrolls-edit",
			Content: tall,
			Steps: []Step{
				{Expect: &StateExpect{Cursor: intp(0), ContentY: intp(0)}},
				{Label: "page right 1", Cmd: "pageright"},
				{Expect: &StateExpect{Cursor: intp(0), ContentYMin: intp(1400), ContentYMax: intp(1600)}},
				{Label: "page right 2", Cmd: "pageright"},
				{Expect: &StateExpect{Cursor: intp(0), ContentYMin: intp(2800), ContentYMax: intp(3200)}},
			},
		},
		{
			Name:    "hw-page-left-scrolls-edit",
			Content: tall,
			Steps: []Step{
				{Label: "seed page right 1", Cmd: "pageright"},
				{Label: "seed page right 2", Cmd: "pageright"},
				{Expect: &StateExpect{Cursor: intp(0), ContentYMin: intp(2800)}},
				{Label: "page left 1", Cmd: "pageleft"},
				{Expect: &StateExpect{Cursor: intp(0), ContentYMin: intp(1400), ContentYMax: intp(1600)}},
				{Label: "page left 2", Cmd: "pageleft"},
				{Expect: &StateExpect{Cursor: intp(0), ContentY: intp(0)}},
			},
		},
	}
}
