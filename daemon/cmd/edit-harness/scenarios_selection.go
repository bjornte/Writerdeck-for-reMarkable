package main

// selectionScenarios cover shift+arrow selection edge cases — grow then shrink
// both ways at N=1/3/7 on fixtureProse. Cursor follows the moving (far) end.
func selectionScenarios() []Scenario {
	return []Scenario{
		{
			Name:    "shift-left-then-right-shrinks",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHMid)},
				{Label: "extend left x1 (N=1)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseHMid), SelStart: intp(proseHMid - 1), SelEnd: intp(proseHMid), SelLen: intp(1)}},
				{Label: "extend left x2 more (N=3)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseHMid), SelStart: intp(proseHMid - 3), SelEnd: intp(proseHMid), SelLen: intp(3)}},
				{Label: "extend left x4 more (N=7)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseHMid), SelStart: intp(proseHMid - 7), SelEnd: intp(proseHMid), SelLen: intp(7)}},
				{Label: "shrink right x1 (N=1)", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseHMid - 1), SelStart: intp(proseHMid - 7), SelEnd: intp(proseHMid - 1), SelLen: intp(6)}},
				{Label: "shrink right x2 more (N=3)", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseHMid - 3), SelStart: intp(proseHMid - 7), SelEnd: intp(proseHMid - 3), SelLen: intp(4)}},
				{Label: "shrink right x4 more (N=7)", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseHMid - 7), SelStart: intp(proseHMid - 7), SelEnd: intp(proseHMid - 7), SelLen: intp(0)}},
			},
		},
		{
			Name:    "shift-right-then-left-shrinks",
			Content: fixtureProse,
			Steps: []Step{
				{SetCursor: intp(proseHMid)},
				{Label: "extend right x1 (N=1)", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseHMid + 1), SelStart: intp(proseHMid), SelEnd: intp(proseHMid + 1), SelLen: intp(1)}},
				{Label: "extend right x2 more (N=3)", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseHMid + 3), SelStart: intp(proseHMid), SelEnd: intp(proseHMid + 3), SelLen: intp(3)}},
				{Label: "extend right x4 more (N=7)", Keys: []Key{{Name: "ArrowRight", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseHMid + 7), SelStart: intp(proseHMid), SelEnd: intp(proseHMid + 7), SelLen: intp(7)}},
				{Label: "shrink left x1 (N=1)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 1},
				{Expect: &StateExpect{Cursor: intp(proseHMid + 6), SelStart: intp(proseHMid), SelEnd: intp(proseHMid + 6), SelLen: intp(6)}},
				{Label: "shrink left x2 more (N=3)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 2},
				{Expect: &StateExpect{Cursor: intp(proseHMid + 4), SelStart: intp(proseHMid), SelEnd: intp(proseHMid + 4), SelLen: intp(4)}},
				{Label: "shrink left x4 more (N=7)", Keys: []Key{{Name: "ArrowLeft", Shift: true}}, Repeat: 4},
				{Expect: &StateExpect{Cursor: intp(proseHMid), SelStart: intp(proseHMid), SelEnd: intp(proseHMid), SelLen: intp(0)}},
			},
		},
	}
}
