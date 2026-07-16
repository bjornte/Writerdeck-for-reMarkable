package main

import "strconv"

// Motion/selection coverage pattern (docs/editor-testing/todo.md):
//
//	uni 1       — one press one way
//	uni 5       — five presses one way
//	bi 1+1      — grow 1, reverse 1 (must shrink / return)
//	bi 3+5      — grow 3, reverse 5 (intentional overshoot past the anchor)
//	bi 7+7      — grow 7, reverse 7
//
// Each axis has a forward and reverse partner scenario (Left↔Right, Up↔Down,
// page-right↔page-left). Between pattern blocks the caret is re-seeded with
// SetCursor so a pure 1+1 / 3+5 / 7+7 is observed — never "grow to 7 first,
// then peel".

func keyArrow(name string, shift, alt, ctrl bool) Key {
	return Key{Name: name, Shift: shift, Alt: alt, Ctrl: ctrl}
}

// caretAxisPattern: plain caret motion along one axis.
// delta is +1 for Right/Down and -1 for Left/Up.
func caretAxisPattern(seed int, forwardKey, backKey Key, delta int) []Step {
	pos := func(n int) int { return seed + delta*n }
	collapse := func(p int) *StateExpect {
		return &StateExpect{Cursor: intp(p), SelStart: intp(p), SelEnd: intp(p)}
	}
	var out []Step
	out = append(out, Step{Label: "reset uni1", SetCursor: intp(seed)})
	out = append(out, Step{Label: "uni 1", Keys: []Key{forwardKey}, Repeat: 1})
	out = append(out, Step{Expect: collapse(pos(1))})
	out = append(out, Step{Label: "reset uni5", SetCursor: intp(seed)})
	out = append(out, Step{Label: "uni 5", Keys: []Key{forwardKey}, Repeat: 5})
	out = append(out, Step{Expect: collapse(pos(5))})
	out = append(out, Step{Label: "reset bi1+1", SetCursor: intp(seed)})
	out = append(out, Step{Label: "bi 1+1 forward", Keys: []Key{forwardKey}, Repeat: 1})
	out = append(out, Step{Expect: collapse(pos(1))})
	out = append(out, Step{Label: "bi 1+1 reverse", Keys: []Key{backKey}, Repeat: 1})
	out = append(out, Step{Expect: collapse(seed)})
	out = append(out, Step{Label: "reset bi3+5", SetCursor: intp(seed)})
	out = append(out, Step{Label: "bi 3+5 forward x3", Keys: []Key{forwardKey}, Repeat: 3})
	out = append(out, Step{Expect: collapse(pos(3))})
	out = append(out, Step{Label: "bi 3+5 reverse x5 (overshoot)", Keys: []Key{backKey}, Repeat: 5})
	out = append(out, Step{Expect: collapse(seed+(-delta)*2)})
	out = append(out, Step{Label: "reset bi7+7", SetCursor: intp(seed)})
	out = append(out, Step{Label: "bi 7+7 forward x7", Keys: []Key{forwardKey}, Repeat: 7})
	out = append(out, Step{Expect: collapse(pos(7))})
	out = append(out, Step{Label: "bi 7+7 reverse x7", Keys: []Key{backKey}, Repeat: 7})
	out = append(out, Step{Expect: collapse(seed)})
	return out
}

// shiftSelectAxisPattern: Shift+arrow grow then reverse-shrink/overshoot.
// towardLow: true when growKey moves the active end toward lower indices.
func shiftSelectAxisPattern(anchor int, growKey, reverseKey Key, towardLow bool) []Step {
	growExpect := func(n int) *StateExpect {
		if towardLow {
			return &StateExpect{
				Cursor: intp(anchor), SelStart: intp(anchor - n), SelEnd: intp(anchor), SelLen: intp(n),
			}
		}
		return &StateExpect{
			Cursor: intp(anchor + n), SelStart: intp(anchor), SelEnd: intp(anchor + n), SelLen: intp(n),
		}
	}
	afterReverse := func(growN, revK int) *StateExpect {
		if towardLow {
			if revK < growN {
				return &StateExpect{
					Cursor: intp(anchor), SelStart: intp(anchor - growN + revK), SelEnd: intp(anchor), SelLen: intp(growN - revK),
				}
			}
			if revK == growN {
				return &StateExpect{Cursor: intp(anchor), SelStart: intp(anchor), SelEnd: intp(anchor), SelLen: intp(0)}
			}
			over := revK - growN
			return &StateExpect{
				Cursor: intp(anchor + over), SelStart: intp(anchor), SelEnd: intp(anchor + over), SelLen: intp(over),
			}
		}
		if revK < growN {
			return &StateExpect{
				Cursor: intp(anchor + growN - revK), SelStart: intp(anchor), SelEnd: intp(anchor + growN - revK), SelLen: intp(growN - revK),
			}
		}
		if revK == growN {
			return &StateExpect{Cursor: intp(anchor), SelStart: intp(anchor), SelEnd: intp(anchor), SelLen: intp(0)}
		}
		over := revK - growN
		return &StateExpect{
			Cursor: intp(anchor), SelStart: intp(anchor - over), SelEnd: intp(anchor), SelLen: intp(over),
		}
	}

	var out []Step
	block := func(label string, growN, revN int, bi bool) {
		out = append(out, Step{Label: "reset " + label, SetCursor: intp(anchor)})
		out = append(out, Step{Label: label + " grow x" + strconv.Itoa(growN), Keys: []Key{growKey}, Repeat: growN})
		out = append(out, Step{Expect: growExpect(growN)})
		if !bi {
			return
		}
		out = append(out, Step{Label: label + " reverse x" + strconv.Itoa(revN), Keys: []Key{reverseKey}, Repeat: revN})
		out = append(out, Step{Expect: afterReverse(growN, revN)})
	}

	block("uni1", 1, 0, false)
	block("uni5", 5, 0, false)
	block("bi1+1", 1, 1, true)
	block("bi3+5", 3, 5, true)
	block("bi7+7", 7, 7, true)
	return out
}

// verticalLinePattern: Down/Up by logical lines on equal-width proseV lines.
func verticalLinePattern(seedLine int, forwardKey, backKey Key, deltaLines int) []Step {
	seed := proseVLineStart(seedLine)
	linePos := func(n int) int { return proseVLineStart(seedLine + deltaLines*n) }
	collapse := func(p int) *StateExpect {
		return &StateExpect{Cursor: intp(p), SelStart: intp(p), SelEnd: intp(p)}
	}
	var out []Step
	out = append(out, Step{Label: "reset uni1", SetCursor: intp(seed)})
	out = append(out, Step{Label: "uni 1", Keys: []Key{forwardKey}, Repeat: 1})
	out = append(out, Step{Expect: collapse(linePos(1))})
	out = append(out, Step{Label: "reset uni5", SetCursor: intp(seed)})
	out = append(out, Step{Label: "uni 5", Keys: []Key{forwardKey}, Repeat: 5})
	out = append(out, Step{Expect: collapse(linePos(5))})
	out = append(out, Step{Label: "reset bi1+1", SetCursor: intp(seed)})
	out = append(out, Step{Label: "bi 1+1 forward", Keys: []Key{forwardKey}, Repeat: 1})
	out = append(out, Step{Expect: collapse(linePos(1))})
	out = append(out, Step{Label: "bi 1+1 reverse", Keys: []Key{backKey}, Repeat: 1})
	out = append(out, Step{Expect: collapse(seed)})
	out = append(out, Step{Label: "reset bi3+5", SetCursor: intp(seed)})
	out = append(out, Step{Label: "bi 3+5 forward x3", Keys: []Key{forwardKey}, Repeat: 3})
	out = append(out, Step{Expect: collapse(linePos(3))})
	out = append(out, Step{Label: "bi 3+5 reverse x5 (overshoot)", Keys: []Key{backKey}, Repeat: 5})
	out = append(out, Step{Expect: collapse(proseVLineStart(seedLine + deltaLines*(-2)))})
	out = append(out, Step{Label: "reset bi7+7", SetCursor: intp(seed)})
	out = append(out, Step{Label: "bi 7+7 forward x7", Keys: []Key{forwardKey}, Repeat: 7})
	out = append(out, Step{Expect: collapse(linePos(7))})
	out = append(out, Step{Label: "bi 7+7 reverse x7", Keys: []Key{backKey}, Repeat: 7})
	out = append(out, Step{Expect: collapse(seed)})
	return out
}

// shiftVerticalPattern: Shift+Down/Up grow then reverse on proseV lines (line starts).
func shiftVerticalPattern(seedLine int, growKey, reverseKey Key, towardLow bool) []Step {
	anchor := proseVLineStart(seedLine)
	// Vertical selection uses line-start landing; sel length spans whole lines between.
	growExpect := func(n int) *StateExpect {
		if towardLow {
			end := proseVLineStart(seedLine - n)
			return &StateExpect{Cursor: intp(anchor), SelStart: intp(end), SelEnd: intp(anchor), SelLen: intp(anchor - end)}
		}
		end := proseVLineStart(seedLine + n)
		return &StateExpect{Cursor: intp(end), SelStart: intp(anchor), SelEnd: intp(end), SelLen: intp(end - anchor)}
	}
	afterReverse := func(growN, revK int) *StateExpect {
		if towardLow {
			if revK < growN {
				end := proseVLineStart(seedLine - growN + revK)
				return &StateExpect{Cursor: intp(anchor), SelStart: intp(end), SelEnd: intp(anchor), SelLen: intp(anchor - end)}
			}
			if revK == growN {
				return &StateExpect{Cursor: intp(anchor), SelStart: intp(anchor), SelEnd: intp(anchor), SelLen: intp(0)}
			}
			over := revK - growN
			end := proseVLineStart(seedLine + over)
			return &StateExpect{Cursor: intp(end), SelStart: intp(anchor), SelEnd: intp(end), SelLen: intp(end - anchor)}
		}
		if revK < growN {
			end := proseVLineStart(seedLine + growN - revK)
			return &StateExpect{Cursor: intp(end), SelStart: intp(anchor), SelEnd: intp(end), SelLen: intp(end - anchor)}
		}
		if revK == growN {
			return &StateExpect{Cursor: intp(anchor), SelStart: intp(anchor), SelEnd: intp(anchor), SelLen: intp(0)}
		}
		over := revK - growN
		end := proseVLineStart(seedLine - over)
		return &StateExpect{Cursor: intp(anchor), SelStart: intp(end), SelEnd: intp(anchor), SelLen: intp(anchor - end)}
	}
	var out []Step
	block := func(label string, growN, revN int, bi bool) {
		out = append(out, Step{Label: "reset " + label, SetCursor: intp(anchor)})
		out = append(out, Step{Label: label + " grow x" + strconv.Itoa(growN), Keys: []Key{growKey}, Repeat: growN})
		out = append(out, Step{Expect: growExpect(growN)})
		if !bi {
			return
		}
		out = append(out, Step{Label: label + " reverse x" + strconv.Itoa(revN), Keys: []Key{reverseKey}, Repeat: revN})
		out = append(out, Step{Expect: afterReverse(growN, revN)})
	}
	block("uni1", 1, 0, false)
	block("uni5", 5, 0, false)
	block("bi1+1", 1, 1, true)
	block("bi3+5", 3, 5, true)
	block("bi7+7", 7, 7, true)
	return out
}

// wordCaretPattern: Alt+Left/Right by word index on the prose word line.
func wordCaretPattern(seedWordIdx int, forwardKey, backKey Key, deltaWord int) []Step {
	seed := proseWordStarts[seedWordIdx]
	wordAt := func(n int) int {
		return proseWordStarts[seedWordIdx+deltaWord*n]
	}
	collapse := func(p int) *StateExpect {
		return &StateExpect{Cursor: intp(p), SelStart: intp(p), SelEnd: intp(p)}
	}
	var out []Step
	out = append(out, Step{Label: "reset uni1", SetCursor: intp(seed)})
	out = append(out, Step{Label: "uni 1", Keys: []Key{forwardKey}, Repeat: 1})
	out = append(out, Step{Expect: collapse(wordAt(1))})
	out = append(out, Step{Label: "reset uni5", SetCursor: intp(seed)})
	out = append(out, Step{Label: "uni 5", Keys: []Key{forwardKey}, Repeat: 5})
	out = append(out, Step{Expect: collapse(wordAt(5))})
	out = append(out, Step{Label: "reset bi1+1", SetCursor: intp(seed)})
	out = append(out, Step{Label: "bi 1+1 forward", Keys: []Key{forwardKey}, Repeat: 1})
	out = append(out, Step{Expect: collapse(wordAt(1))})
	out = append(out, Step{Label: "bi 1+1 reverse", Keys: []Key{backKey}, Repeat: 1})
	out = append(out, Step{Expect: collapse(seed)})
	out = append(out, Step{Label: "reset bi3+5", SetCursor: intp(seed)})
	out = append(out, Step{Label: "bi 3+5 forward x3", Keys: []Key{forwardKey}, Repeat: 3})
	out = append(out, Step{Expect: collapse(wordAt(3))})
	out = append(out, Step{Label: "bi 3+5 reverse x5 (overshoot)", Keys: []Key{backKey}, Repeat: 5})
	out = append(out, Step{Expect: collapse(wordAt(-2))})
	out = append(out, Step{Label: "reset bi7+7", SetCursor: intp(seed)})
	out = append(out, Step{Label: "bi 7+7 forward x7", Keys: []Key{forwardKey}, Repeat: 7})
	out = append(out, Step{Expect: collapse(wordAt(7))})
	out = append(out, Step{Label: "bi 7+7 reverse x7", Keys: []Key{backKey}, Repeat: 7})
	out = append(out, Step{Expect: collapse(seed)})
	return out
}
