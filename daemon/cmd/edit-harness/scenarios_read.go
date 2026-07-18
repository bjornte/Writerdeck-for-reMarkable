package main

// readScenarios cover preview/read-mode scrolling. Esc toggles edit↔preview;
// hardware page cmds still inject via pageleft/pageright.
func readScenarios() []Scenario {
	tall := fixtureTallRead()
	return []Scenario{
		{
			// After reaching document end in read mode, further page-downs must
			// clamp (no endless empty scroll). Ten extra downs leave contentY
			// unchanged; one page-up returns into the document.
			Name:    "read-overscroll-clamps",
			Content: tall,
			Steps: []Step{
				{Expect: &StateExpect{Mode: intp(1), ContentY: intp(0)}},
				{Label: "Esc to preview", Keys: []Key{{Name: "Escape"}}},
				{Expect: &StateExpect{Mode: intp(0), ContentY: intp(0)}},
				{Label: "page toward end (x20)", Cmd: "pageright", Repeat: 20},
				{Label: "capture end contentY", CaptureContentY: true, Expect: &StateExpect{Mode: intp(0)}},
				{Label: "ten extra page-downs must clamp", Cmd: "pageright", Repeat: 10},
				{Label: "still at end", ExpectContentYEqCaptured: true, Expect: &StateExpect{Mode: intp(0)}},
				{Label: "one page-up returns into doc", Cmd: "pageleft", Repeat: 1},
				{Label: "contentY decreased", ExpectContentYLtCaptured: true, Expect: &StateExpect{Mode: intp(0)}},
				{Label: "Esc back to edit", Keys: []Key{{Name: "Escape"}}},
				{Expect: &StateExpect{Mode: intp(1)}},
			},
		},
	}
}
