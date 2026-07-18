package main

// Named motions for harness scenarios.
// Naming follows Meta Lexical helpers (moveToLine*, deleteLine*) and Google
// Selection API units (lineboundary vs paragraphboundary). Harness Ctrl = Mac Cmd.
// Layout vs meaning units: Finseth, The Craft of Text Editing, ch. 4.

func keyMoveToVisualLineStart(shift bool) Key {
	return Key{Name: "ArrowLeft", Ctrl: true, Shift: shift}
}

func keyMoveToVisualLineEnd(shift bool) Key {
	return Key{Name: "ArrowRight", Ctrl: true, Shift: shift}
}

func keyMoveWordLeft(shift bool) Key {
	return Key{Name: "ArrowLeft", Alt: true, Shift: shift}
}

func keyMoveWordRight(shift bool) Key {
	return Key{Name: "ArrowRight", Alt: true, Shift: shift}
}

func keyMoveParagraphBackward(shift bool) Key {
	return Key{Name: "ArrowUp", Alt: true, Shift: shift}
}

func keyMoveParagraphForward(shift bool) Key {
	return Key{Name: "ArrowDown", Alt: true, Shift: shift}
}

func keyMoveDocStart(shift bool) Key {
	return Key{Name: "Home", Ctrl: true, Shift: shift}
}

func keyMoveDocEnd(shift bool) Key {
	return Key{Name: "End", Ctrl: true, Shift: shift}
}

func keyHomeVisualLine(shift bool) Key {
	return Key{Name: "Home", Shift: shift}
}

func keyEndVisualLine(shift bool) Key {
	return Key{Name: "End", Shift: shift}
}

func keyDeleteToVisualLineStart() Key {
	return Key{Name: "Backspace", Ctrl: true}
}

func keyDeleteWordBackward() Key {
	return Key{Name: "Backspace", Alt: true}
}

// expectVisualLineCaret collapses the caret at pos and keeps TextLen.
// Callers must pass a wrap-boundary index, not paragraph/doc end.
func expectVisualLineCaret(pos, textLen int) *StateExpect {
	return &StateExpect{
		Cursor:   intp(pos),
		SelStart: intp(pos),
		SelEnd:   intp(pos),
		TextLen:  intp(textLen),
	}
}

// isLineBoundaryScenario names that must assert visual-line ends (lineboundary),
// never paragraph/document end on a wrapped fixture.
func isLineBoundaryScenario(name string) bool {
	switch name {
	case "wrap-ctrl-left", "wrap-ctrl-right",
		"wrap-shift-ctrl-left", "wrap-shift-ctrl-right",
		"wrap-end-on-visual-line", "wrap-home-on-visual-line",
		"wrap-combo-ctrl-bs-line":
		return true
	default:
		return false
	}
}

// expectAllowsParagraphEnd reports whether an expect could accept caret at
// logical/paragraph end of a wrapped doc (the classic false green).
func expectAllowsParagraphEnd(exp *StateExpect, docLen int) bool {
	if exp == nil || docLen < 2 {
		return false
	}
	paraEnd := docLen - 1
	if exp.Cursor != nil && *exp.Cursor >= paraEnd {
		return true
	}
	if exp.CursorMax != nil && *exp.CursorMax >= paraEnd {
		return true
	}
	return false
}
