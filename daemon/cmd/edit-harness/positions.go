package main

import "unicode/utf8"

// utf8ByteAtRune returns the byte offset at rune index i (0-based) in s.
func utf8ByteAtRune(s string, runeIndex int) int {
	if runeIndex <= 0 {
		return 0
	}
	i := 0
	for pos := range s {
		if i == runeIndex {
			return pos
		}
		i++
	}
	return len(s)
}

// utf8Len returns the UTF-8 byte length of s (same as len(s) but documents intent).
func utf8Len(s string) int {
	return len(s)
}

// utf8RuneCount returns the number of runes in s.
func utf8RuneCount(s string) int {
	return utf8.RuneCountInString(s)
}
