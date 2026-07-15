package main

const harnessWrapWidth = 320

// Calibrated on device width=320 (2026-07-15, rM1 Writerdeck @ 1e62aff). Re-run
// `bash scripts/test-keyboard-harness.sh -m wrap -v` after font or width changes.
const (
	wrapParagraphLen    = 199
	wrapDownOneCursor     = 20 // Ctrl+Home, Down×1 on wrapParagraph
	wrapDownTwoCursor     = 40 // Ctrl+Home, Down×2 on wrapParagraph
	wrapGoalColDownCursor = 24 // "ab"+word×35, col 2, Down×1 at W=320
)
