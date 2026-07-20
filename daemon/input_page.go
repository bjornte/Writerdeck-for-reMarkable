package main

// physicalPageCmd maps bezel left/right buttons to editor pageleft (scroll up)
// / pageright (scroll down), accounting for display rotation.
//
// Portrait (0) and 90 deg CCW (270): left = up, right = down.
// Upside-down (180) and 90 deg CW (90): left/right are flipped.
func physicalPageCmd(isLeft bool, rot int) string {
	rot = normalizeRotation(rot)
	flip := rot == 90 || rot == 180
	up := isLeft != flip
	if up {
		return "pageleft"
	}
	return "pageright"
}
