package main

const harnessWrapWidth = 320

// Calibrated on device width=320 (2026-07-14, rM1 Writerdeck). Re-run
// `bash scripts/test-keyboard-harness.sh -m wrap -v` after font or width changes.
const (
	wrapParagraphLen    = 199
	wrapDownOneCursor     = 10 // Ctrl+Home, Down×1 on wrapParagraph
	wrapDownTwoCursor     = 20 // Ctrl+Home, Down×2 on wrapParagraph (approx; tighten after calibrate)
	wrapGoalColDownCursor = 10 // "ab"+word×35, col 2, Down×1 at W=320
)
