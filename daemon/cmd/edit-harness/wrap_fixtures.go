package main

const harnessWrapWidth = 320

// Calibrated on device width=320 (2026-07-16, visual-line next-y walk @ fee53bc/6dfade8).
// Prior 20-char/row values matched the old height-based stepper, which skipped
// every other visual row on this font/width. Re-run wrap tag after font changes.
//
// Down from Ctrl+Home lands on the start of each visual row. End / Ctrl+Right from
// the start of a row lands on the same index as Down to the *next* row (wrap point).
const (
	wrapParagraphLen      = 199
	wrapDownOneCursor     = 10  // Ctrl+Home, Downx1 on wrapParagraph
	wrapDownTwoCursor     = 20  // Ctrl+Home, Downx2
	wrapDownThreeCursor   = 30  // Ctrl+Home, Downx3
	wrapDownFourCursor    = 40  // Ctrl+Home, Downx4
	wrapDownSixCursor     = 60  // Ctrl+Home, Downx6
	wrapDownSevenCursor   = 70  // Ctrl+Home, Downx7
	wrapGoalColDownCursor = 14  // "ab"+wordx35, col 2, Downx1 at W=320

	// Visual-line ends (== next row start). Must stay << wrapParagraphLen.
	wrapEndVisualRow0 = wrapDownOneCursor // End / Ctrl+Right from 0
	wrapEndVisualRow1 = wrapDownTwoCursor // End / Ctrl+Right from wrapDownOneCursor
)
