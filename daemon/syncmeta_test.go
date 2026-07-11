package main

import "testing"

func TestStrHash(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", "5381"},
		{"hello", "261238937"},
		{"# Test\n", "3670303922"},
		{"æøå", "193635880"},
	}
	for _, c := range cases {
		got := strHash(c.in)
		if got != c.want {
			t.Errorf("strHash(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
