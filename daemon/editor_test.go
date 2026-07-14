package main

import (
	"strings"
	"testing"
)

func TestTranslateShiftArrow(t *testing.T) {
	cases := []struct {
		name string
		ev   wsMsg
		want string
	}{
		{
			name: "shift arrow down",
			ev:   wsMsg{Type: "key", Key: "ArrowDown", Shift: true},
			want: `{"t":"key","k":"ArrowDown","m":1}`,
		},
		{
			name: "delete key",
			ev:   wsMsg{Type: "key", Key: "Delete"},
			want: `{"t":"key","k":"Delete"}`,
		},
		{
			name: "shift end release",
			ev:   wsMsg{Type: "key", Key: "End", Shift: true, Action: "release"},
			want: `{"t":"key","k":"End","m":1,"u":1}`,
		},
		{
			name: "ctrl shift arrow left",
			ev:   wsMsg{Type: "key", Key: "ArrowLeft", Shift: true, Ctrl: true},
			want: `{"t":"key","k":"ArrowLeft","m":3}`,
		},
		{
			name: "meta shift arrow right",
			ev:   wsMsg{Type: "key", Key: "ArrowRight", Shift: true, Meta: true},
			want: `{"t":"key","k":"ArrowRight","m":9}`,
		},
		{
			name: "alt shift arrow up",
			ev:   wsMsg{Type: "key", Key: "ArrowUp", Shift: true, Alt: true},
			want: `{"t":"key","k":"ArrowUp","m":5}`,
		},
		{
			name: "ctrl meta shift arrow down",
			ev:   wsMsg{Type: "key", Key: "ArrowDown", Shift: true, Ctrl: true, Meta: true},
			want: `{"t":"key","k":"ArrowDown","m":11}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := string(translate(tc.ev))
			if got != tc.want {
				t.Fatalf("translate() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTranslateHarnessNamedKeys(t *testing.T) {
	for browser, feed := range namedKeys {
		ev := wsMsg{Type: "key", Key: browser}
		got := string(translate(ev))
		want := `{"t":"key","k":"` + feed + `"}`
		if got != want {
			t.Fatalf("translate(%q) = %q want %q", browser, got, want)
		}
	}
}

func TestTranslateShiftArrowNotDropped(t *testing.T) {
	ev := wsMsg{Type: "key", Key: "ArrowUp", Shift: true, Ctrl: true}
	line := translate(ev)
	if line == nil {
		t.Fatal("translate returned nil for Ctrl+Shift+ArrowUp")
	}
	if !strings.Contains(string(line), `"k":"ArrowUp"`) {
		t.Fatalf("unexpected payload: %s", line)
	}
	if !strings.Contains(string(line), `"m":3`) {
		t.Fatalf("modifier mask missing shift+ctrl: %s", line)
	}
}
