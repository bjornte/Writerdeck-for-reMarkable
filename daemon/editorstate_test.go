package main

import "testing"

func TestParseEditorState(t *testing.T) {
	line := []byte(`{"t":"state","cursor":3,"selStart":0,"selEnd":3,"textLen":6,"mode":1}`)
	st, ok := parseEditorState(line)
	if !ok {
		t.Fatal("parse failed")
	}
	if st.Cursor != 3 || st.SelStart != 0 || st.SelEnd != 3 || st.TextLen != 6 || st.Mode != 1 {
		t.Fatalf("unexpected %+v", st)
	}
}
