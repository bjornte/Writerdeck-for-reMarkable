package main

import "testing"

func TestCompareProductVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"2026-07-18", "2026-07-18", 0},
		{"2026-07-01", "2026-07-18", -1},
		{"2026-07-18", "2026-07-01", 1},
		{"2026-07-18", "2026-07-18.2", -1},
		{"2026-07-18.2", "2026-07-18", 1},
		{"2026-07-18.2", "2026-07-18.3", -1},
		{"unknown", "2026-07-18", -1},
		{"2026-07-18", "unknown", 1},
	}
	for _, c := range cases {
		got := compareProductVersions(c.a, c.b)
		if got != c.want {
			t.Errorf("compare(%q,%q)=%d want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestCombineProductVersion(t *testing.T) {
	cases := []struct {
		server, editor, want string
		mismatch             bool
	}{
		{"2026-07-20", "2026-07-20", "2026-07-20", false},
		{"2026-07-20", "2026-07-18", "2026-07-18", true},
		{"2026-07-18", "2026-07-20", "2026-07-18", true},
		{"2026-07-20", "unknown", "unknown", true},
		{"unknown", "2026-07-20", "unknown", true},
	}
	for _, c := range cases {
		got, mis := combineProductVersion(c.server, c.editor)
		if got != c.want || mis != c.mismatch {
			t.Errorf("combine(%q,%q)=(%q,%v) want (%q,%v)",
				c.server, c.editor, got, mis, c.want, c.mismatch)
		}
	}
}

func TestFormatVersionMessage(t *testing.T) {
	cases := []struct {
		product, latest string
		err             error
		mismatched      bool
		want            string
	}{
		{"2026-07-18", "2026-07-18", nil, false, "Writerdeck version 2026-07-18 (latest)"},
		{"2026-07-01", "2026-07-18", nil, false, "Writerdeck version 2026-07-01. Latest on GitHub is 2026-07-18."},
		{"2026-07-01", "", fmtTestErr, false, "Writerdeck version 2026-07-01 (couldn't reach GitHub to check for updates)"},
		{"2026-07-18.2", "2026-07-18", nil, false, "Writerdeck version 2026-07-18.2 (newer than GitHub 2026-07-18)"},
		{"2026-07-18", "2026-07-20", nil, true, "Writerdeck version 2026-07-18 (server and editor differ — update both). Latest on GitHub is 2026-07-20."},
		{"2026-07-18", "2026-07-18", nil, true, "Writerdeck version 2026-07-18 (server and editor differ — update both)"},
	}
	for _, c := range cases {
		got := formatVersionMessage(c.product, c.latest, c.err, c.mismatched)
		if got != c.want {
			t.Errorf("product=%q latest=%q mis=%v err=%v\n got %q\nwant %q",
				c.product, c.latest, c.mismatched, c.err, got, c.want)
		}
	}
}

var fmtTestErr = errString("network down")

type errString string

func (e errString) Error() string { return string(e) }
