package main

// selectionScenarios: horizontal Shift reverse + vertical Shift reverse partner
// (shift-down-then-up-shrinks lives in regression for critical tag continuity).
func selectionScenarios() []Scenario {
	left := keyArrow("ArrowLeft", true, false, false)
	right := keyArrow("ArrowRight", true, false, false)
	up := keyArrow("ArrowUp", true, false, false)
	down := keyArrow("ArrowDown", true, false, false)

	return []Scenario{
		{
			Name:    "shift-left-then-right-shrinks",
			Content: fixtureProse,
			Steps:   shiftSelectAxisPattern(proseHMid, left, right, true),
		},
		{
			Name:    "shift-right-then-left-shrinks",
			Content: fixtureProse,
			Steps:   shiftSelectAxisPattern(proseHMid, right, left, false),
		},
		{
			Name:    "shift-up-then-down-shrinks",
			Content: fixtureProse,
			Steps:   shiftVerticalPattern(9, up, down, true),
		},
	}
}
