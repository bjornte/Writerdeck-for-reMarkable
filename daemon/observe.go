// Writerdeck-server -- see main.go for overview.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Observation mode records phone/WebSocket keys plus caret/selection snapshots
// so a user can demonstrate a typing bug to an LLM. USB-on-tablet keys are not
// on this path. Export shape matches edit-harness Scenario JSON.

const (
	observeMaxSteps   = 2000
	observeStateQueue = 64
)

type observeKey struct {
	Name   string `json:"name"`
	Shift  bool   `json:"shift,omitempty"`
	Ctrl   bool   `json:"ctrl,omitempty"`
	Alt    bool   `json:"alt,omitempty"`
	Meta   bool   `json:"meta,omitempty"`
	Action string `json:"action,omitempty"`
}

type observeExpect struct {
	Cursor   *int    `json:"cursor,omitempty"`
	SelStart *int    `json:"selStart,omitempty"`
	SelEnd   *int    `json:"selEnd,omitempty"`
	TextLen  *int    `json:"textLen,omitempty"`
	Text     *string `json:"text,omitempty"`
}

type observeStep struct {
	Label  string         `json:"label,omitempty"`
	Keys   []observeKey   `json:"keys,omitempty"`
	Expect *observeExpect `json:"expect,omitempty"`
}

type observeScenario struct {
	Name    string        `json:"name"`
	Content string        `json:"content"`
	Note    string        `json:"note,omitempty"`
	Tags    []string      `json:"tags"`
	Steps   []observeStep `json:"steps"`
}

type observeStatus struct {
	Enabled   bool   `json:"enabled"` // settings.observe — phone button visible when true
	Active    bool   `json:"active"`
	Steps     int    `json:"steps"`
	Note      string `json:"note,omitempty"`
	Started   string `json:"started,omitempty"`
	HasExport bool   `json:"hasExport"`
	Ready     bool   `json:"ready"` // stopped with an export waiting for the laptop/agent
}

type observeJob struct {
	stepIdx int
	final   bool // attach full text when true
}

type observeRecorder struct {
	mu      sync.Mutex
	active  bool
	started time.Time
	note    string
	content string
	steps   []observeStep

	jobs   chan observeJob
	stopCh chan struct{}
	wg     sync.WaitGroup

	lastExport []byte
}

var observe = &observeRecorder{}

func observeEnabled() bool {
	settingsMu.Lock()
	on := curSettings.Observe
	settingsMu.Unlock()
	return on
}

func (o *observeRecorder) status() observeStatus {
	enabled := observeEnabled()
	// Warm memory from disk so Ready stays true after a server restart.
	_ = o.loadExport()
	o.mu.Lock()
	defer o.mu.Unlock()
	has := len(o.lastExport) > 0
	st := observeStatus{
		Enabled:   enabled,
		Active:    o.active,
		Steps:     len(o.steps),
		Note:      o.note,
		HasExport: has,
		Ready:     has && !o.active,
	}
	if o.active && !o.started.IsZero() {
		st.Started = o.started.UTC().Format(time.RFC3339)
	}
	return st
}

func observeExportPath() string {
	return filepath.Join(filepath.Dir(settingsFilePath), "last-observation.json")
}

func (o *observeRecorder) persistExport(raw []byte) {
	path := observeExportPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "writerdeck-server: observe persist mkdir: %v\n", err)
		return
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "writerdeck-server: observe persist write: %v\n", err)
		return
	}
	if err := os.Rename(tmp, path); err != nil {
		fmt.Fprintf(os.Stderr, "writerdeck-server: observe persist rename: %v\n", err)
	}
}

// loadExport returns the last stop payload from memory, or from disk after restart.
func (o *observeRecorder) loadExport() []byte {
	o.mu.Lock()
	raw := o.lastExport
	o.mu.Unlock()
	if len(raw) > 0 {
		return raw
	}
	data, err := os.ReadFile(observeExportPath())
	if err != nil {
		return nil
	}
	o.mu.Lock()
	if len(o.lastExport) == 0 {
		o.lastExport = data
	}
	raw = o.lastExport
	o.mu.Unlock()
	return raw
}

func (o *observeRecorder) start() error {
	if !observeEnabled() {
		return fmt.Errorf("observation disabled")
	}
	o.mu.Lock()
	if o.active {
		o.mu.Unlock()
		return fmt.Errorf("already observing")
	}
	if activeSess == nil || !activeSess.isActive() {
		o.mu.Unlock()
		return fmt.Errorf("no active editor session")
	}
	o.mu.Unlock()

	st, err := queryEditorState()
	if err != nil {
		return err
	}
	content := st.Text
	if content == "" {
		currentNoteMu.Lock()
		name := currentNote
		currentNoteMu.Unlock()
		if name != "" {
			if body, rerr := readNoteBody(name); rerr == nil {
				content = body
			}
		}
	}

	currentNoteMu.Lock()
	note := currentNote
	currentNoteMu.Unlock()
	if note == "" {
		note = st.CurrentFile
	}

	o.mu.Lock()
	o.active = true
	o.started = time.Now().UTC()
	o.note = note
	o.content = content
	o.steps = nil
	o.jobs = make(chan observeJob, observeStateQueue)
	o.stopCh = make(chan struct{})
	o.wg.Add(1)
	go o.stateWorker()
	o.mu.Unlock()

	fmt.Fprintf(os.Stderr, "writerdeck-server: observe started note=%q textLen=%d\n", note, len(content))
	broadcastObserveStatus()
	return nil
}

func (o *observeRecorder) stop() ([]byte, error) {
	o.mu.Lock()
	if !o.active {
		exp := o.lastExport
		o.mu.Unlock()
		if len(exp) == 0 {
			return nil, fmt.Errorf("not observing")
		}
		return exp, nil
	}
	o.active = false
	stopCh := o.stopCh
	o.mu.Unlock()

	close(stopCh)
	o.wg.Wait()

	// Final caret/selection + live text after the last key settles.
	var finalExpect *observeExpect
	if st, err := queryEditorState(); err == nil {
		finalExpect = expectFromState(st, true)
	}

	o.mu.Lock()
	if finalExpect != nil {
		if n := len(o.steps); n > 0 {
			o.steps[n-1].Expect = finalExpect
		} else {
			o.steps = append(o.steps, observeStep{
				Label:  "initial",
				Expect: finalExpect,
			})
		}
	}
	raw, err := o.buildExportLocked()
	if err != nil {
		o.mu.Unlock()
		return nil, err
	}
	o.lastExport = raw
	steps := len(o.steps)
	o.mu.Unlock()

	o.persistExport(raw)
	fmt.Fprintf(os.Stderr, "writerdeck-server: observe stopped steps=%d bytes=%d path=%s\n", steps, len(raw), observeExportPath())
	broadcastObserveStatus()
	return raw, nil
}

func (o *observeRecorder) buildExportLocked() ([]byte, error) {
	name := fmt.Sprintf("observed-%s", o.started.Format("2006-01-02-150405"))
	sc := observeScenario{
		Name:    name,
		Content: o.content,
		Note:    o.note,
		Tags:    []string{"observed"},
		Steps:   append([]observeStep(nil), o.steps...),
	}
	return json.MarshalIndent(sc, "", "  ")
}

func (o *observeRecorder) recordKey(ev wsMsg) {
	o.mu.Lock()
	if !o.active {
		o.mu.Unlock()
		return
	}
	if len(o.steps) >= observeMaxSteps {
		o.mu.Unlock()
		return
	}
	k := observeKey{
		Name:   ev.Key,
		Shift:  ev.Shift,
		Ctrl:   ev.Ctrl,
		Alt:    ev.Alt,
		Meta:   ev.Meta,
		Action: ev.Action,
	}
	idx := len(o.steps)
	o.steps = append(o.steps, observeStep{Keys: []observeKey{k}})
	jobs := o.jobs
	o.mu.Unlock()

	select {
	case jobs <- observeJob{stepIdx: idx}:
	default:
		// Queue full: keep the key; skip this snapshot rather than stall typing.
	}
}

func (o *observeRecorder) stateWorker() {
	defer o.wg.Done()
	for {
		select {
		case <-o.stopCh:
			// Drain remaining jobs without new queries after stop begins;
			// stop() takes the authoritative final snapshot.
			for {
				select {
				case <-o.jobs:
				default:
					return
				}
			}
		case job := <-o.jobs:
			st, err := queryEditorState()
			if err != nil {
				continue
			}
			exp := expectFromState(st, job.final)
			o.mu.Lock()
			if job.stepIdx >= 0 && job.stepIdx < len(o.steps) {
				o.steps[job.stepIdx].Expect = exp
			}
			o.mu.Unlock()
		}
	}
}

func expectFromState(st EditorState, withText bool) *observeExpect {
	c, ss, se, tl := st.Cursor, st.SelStart, st.SelEnd, st.TextLen
	exp := &observeExpect{
		Cursor:   &c,
		SelStart: &ss,
		SelEnd:   &se,
		TextLen:  &tl,
	}
	if withText && st.Text != "" {
		t := st.Text
		exp.Text = &t
	}
	return exp
}

// readNoteBody loads a note from disk for observe start when live text is empty.
func readNoteBody(name string) (string, error) {
	path := notesSafe(name)
	if path == "" {
		return "", fmt.Errorf("bad note name")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func broadcastObserveStatus() {
	st := observe.status()
	msg, _ := json.Marshal(struct {
		Type      string `json:"type"`
		Active    bool   `json:"active"`
		Steps     int    `json:"steps"`
		Ready     bool   `json:"ready"`
		HasExport bool   `json:"hasExport"`
	}{"observe", st.Active, st.Steps, st.Ready, st.HasExport})
	broadcast(msg)
}

func observeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if !checkAuth(w, r) {
		return
	}

	path := r.URL.Path
	switch {
	case path == "/api/observe/status" && r.Method == http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(observe.status()) //nolint:errcheck
		return

	case path == "/api/observe/start" && r.Method == http.MethodPost:
		if err := observe.start(); err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(observe.status()) //nolint:errcheck
		return

	case path == "/api/observe/stop" && r.Method == http.MethodPost:
		raw, err := observe.stop()
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(raw) //nolint:errcheck
		return

	case path == "/api/observe/export" && r.Method == http.MethodGet:
		raw := observe.loadExport()
		if len(raw) == 0 {
			http.Error(w, "no observation export", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(raw) //nolint:errcheck
		return
	}

	http.NotFound(w, r)
}
