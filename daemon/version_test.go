package main

import "testing"

func TestFormatVersionMessage(t *testing.T) {
	cases := []struct {
		local, latest string
		err           error
		want          string
	}{
		{"2026-07-18", "2026-07-18", nil, "Writerdeck version 2026-07-18 (latest)"},
		{"2026-07-01", "2026-07-18", nil, "Writerdeck version 2026-07-01. Latest on GitHub is 2026-07-18."},
		{"2026-07-01", "", fmtTestErr, "Writerdeck version 2026-07-01 (couldn't reach GitHub to check for updates)"},
		{"2026-07-01.2", "2026-07-01", nil, "Writerdeck version 2026-07-01.2. Latest on GitHub is 2026-07-01."},
	}
	for _, c := range cases {
		got := formatVersionMessage(c.local, c.latest, c.err)
		if got != c.want {
			t.Errorf("local=%q latest=%q err=%v\n got %q\nwant %q", c.local, c.latest, c.err, got, c.want)
		}
	}
}

var fmtTestErr = errString("network down")

type errString string

func (e errString) Error() string { return string(e) }
