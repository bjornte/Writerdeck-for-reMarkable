package main

import "testing"

func TestPhysicalPageCmd(t *testing.T) {
	cases := []struct {
		left bool
		rot  int
		want string
	}{
		{true, 0, "pageleft"},
		{false, 0, "pageright"},
		{true, 270, "pageleft"},  // 90° CCW: left scrolls up
		{false, 270, "pageright"},
		{true, 180, "pageright"}, // upside down: flip
		{false, 180, "pageleft"},
		{true, 90, "pageright"},
		{false, 90, "pageleft"},
		{true, -90, "pageleft"}, // normalize to 270
	}
	for _, c := range cases {
		got := physicalPageCmd(c.left, c.rot)
		if got != c.want {
			t.Errorf("physicalPageCmd(left=%v, rot=%d)=%q want %q", c.left, c.rot, got, c.want)
		}
	}
}
