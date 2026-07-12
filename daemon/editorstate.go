package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// EditorState is a cursor/selection snapshot from the tablet TextEdit.
type EditorState struct {
	Cursor   int `json:"cursor"`
	SelStart int `json:"selStart"`
	SelEnd   int `json:"selEnd"`
	TextLen  int `json:"textLen"`
	Mode     int `json:"mode"` // 0=preview, 1=edit
}

const stateQueryTimeout = 3 * time.Second

type stateWait struct {
	ch chan EditorState
}

func (e *editorConn) registerStateWait() chan EditorState {
	ch := make(chan EditorState, 1)
	e.stateMu.Lock()
	e.stateWait = ch
	e.stateMu.Unlock()
	return ch
}

func (e *editorConn) cancelStateWait(ch chan EditorState) {
	e.stateMu.Lock()
	if e.stateWait == ch {
		e.stateWait = nil
	}
	e.stateMu.Unlock()
}

func (e *editorConn) deliverState(st EditorState) {
	e.stateMu.Lock()
	ch := e.stateWait
	e.stateWait = nil
	e.stateMu.Unlock()
	if ch == nil {
		return
	}
	select {
	case ch <- st:
	default:
	}
}

// queryEditorState asks Writerdeck to publish cursor/selection and waits for the reply.
func queryEditorState() (EditorState, error) {
	if globalEC == nil || !globalEC.ready() {
		return EditorState{}, fmt.Errorf("editor socket not connected")
	}
	ch := globalEC.registerStateWait()
	globalEC.write([]byte(`{"t":"cmd","c":"editorstate"}`))
	select {
	case st := <-ch:
		return st, nil
	case <-time.After(stateQueryTimeout):
		globalEC.cancelStateWait(ch)
		return EditorState{}, fmt.Errorf("editor state timeout")
	}
}

func parseEditorState(line []byte) (EditorState, bool) {
	var raw struct {
		T        string `json:"t"`
		Cursor   int    `json:"cursor"`
		SelStart int    `json:"selStart"`
		SelEnd   int    `json:"selEnd"`
		TextLen  int    `json:"textLen"`
		Mode     int    `json:"mode"`
	}
	if err := json.Unmarshal(line, &raw); err != nil || raw.T != "state" {
		return EditorState{}, false
	}
	return EditorState{
		Cursor:   raw.Cursor,
		SelStart: raw.SelStart,
		SelEnd:   raw.SelEnd,
		TextLen:  raw.TextLen,
		Mode:     raw.Mode,
	}, true
}

func (s EditorState) selEmpty() bool {
	return s.SelStart == s.SelEnd
}

func (s EditorState) selLen() int {
	if s.selEmpty() {
		return 0
	}
	if s.SelStart < s.SelEnd {
		return s.SelEnd - s.SelStart
	}
	return s.SelStart - s.SelEnd
}

func (s EditorState) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}
