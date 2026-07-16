package main

const harnessWrapWidth = 320

// Calibrated on device width=320 (2026-07-15, rM1 Writerdeck @ 1e62aff).
// N=3/N=7 offsets extrapolated from the one-/two-row calibrations; re-run
// `bash scripts/test-keyboard-harness.sh -m wrap -v` after font/width changes
// and update these if they drift.
const (
	wrapParagraphLen      = 199
	wrapDownOneCursor     = 20  // Ctrl+Home, Down×1 on wrapParagraph
	wrapDownTwoCursor     = 40  // Ctrl+Home, Down×2
	wrapDownThreeCursor   = 60  // Ctrl+Home, Down×3
	wrapDownFourCursor    = 80  // Ctrl+Home, Down×4 (linear ~20/row at W=320)
	wrapDownSixCursor     = 120 // Ctrl+Home, Down×6
	wrapDownSevenCursor   = 140 // Ctrl+Home, Down×7
	wrapGoalColDownCursor = 24  // "ab"+word×35, col 2, Down×1 at W=320
)
