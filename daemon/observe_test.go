package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestExpectFromState(t *testing.T) {
	st := EditorState{Cursor: 10, SelStart: 3, SelEnd: 10, TextLen: 40, Text: "hello"}
	exp := expectFromState(st, false)
	if exp.Cursor == nil || *exp.Cursor != 10 {
		t.Fatalf("cursor: %+v", exp.Cursor)
	}
	if exp.Text != nil {
		t.Fatal("text should be omitted without withText")
	}
	exp2 := expectFromState(st, true)
	if exp2.Text == nil || *exp2.Text != "hello" {
		t.Fatalf("text: %+v", exp2.Text)
	}
}

func TestObserveBuildExport(t *testing.T) {
	o := &observeRecorder{
		started: time.Date(2026, 7, 18, 15, 4, 5, 0, time.UTC),
		note:    "demo.md",
		content: "line one\nline two\n",
		steps: []observeStep{
			{Keys: []observeKey{{Name: "ArrowDown", Shift: true}}},
			{
				Keys: []observeKey{{Name: "ArrowDown", Shift: true}},
				Expect: &observeExpect{
					Cursor:   intPtr(20),
					SelStart: intPtr(0),
					SelEnd:   intPtr(20),
					TextLen:  intPtr(18),
				},
			},
		},
	}
	raw, err := o.buildExportLocked()
	if err != nil {
		t.Fatal(err)
	}
	var sc observeScenario
	if err := json.Unmarshal(raw, &sc); err != nil {
		t.Fatal(err)
	}
	if sc.Name != "observed-2026-07-18-150405" {
		t.Fatalf("name: %q", sc.Name)
	}
	if sc.Content != "line one\nline two\n" {
		t.Fatalf("content: %q", sc.Content)
	}
	if len(sc.Tags) != 1 || sc.Tags[0] != "observed" {
		t.Fatalf("tags: %+v", sc.Tags)
	}
	if len(sc.Steps) != 2 {
		t.Fatalf("steps: %d", len(sc.Steps))
	}
	if !sc.Steps[0].Keys[0].Shift || sc.Steps[0].Keys[0].Name != "ArrowDown" {
		t.Fatalf("key: %+v", sc.Steps[0].Keys[0])
	}
	if sc.Steps[1].Expect == nil || *sc.Steps[1].Expect.SelEnd != 20 {
		t.Fatalf("expect: %+v", sc.Steps[1].Expect)
	}
}

func TestObserveRecordKeyInactive(t *testing.T) {
	o := &observeRecorder{}
	o.recordKey(wsMsg{Type: "key", Key: "a"})
	if len(o.steps) != 0 {
		t.Fatal("inactive recorder should ignore keys")
	}
}

func intPtr(n int) *int { return &n }
