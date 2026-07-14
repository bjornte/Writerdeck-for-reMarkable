package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	harnessNote   = "z-test-keyboard-harness.md"
	defaultKeyMs  = 120
	defaultStepMs = 200
	fastKeyMs     = 40
	fastStepMs    = 80
	httpTimeout   = 30 * time.Second
)

var (
	keyPause  = defaultKeyMs * time.Millisecond
	stepPause = defaultStepMs * time.Millisecond
)

var harnessHTTP = &http.Client{Timeout: httpTimeout}

type Key struct {
	Name  string `json:"name"`
	Shift bool   `json:"shift,omitempty"`
	Ctrl  bool   `json:"ctrl,omitempty"`
	Alt   bool   `json:"alt,omitempty"`
	Meta  bool   `json:"meta,omitempty"`
}

type StateExpect struct {
	Cursor    *int `json:"cursor,omitempty"`
	CursorMin *int `json:"cursorMin,omitempty"`
	CursorMax *int `json:"cursorMax,omitempty"`
	SelStart  *int `json:"selStart,omitempty"`
	SelEnd    *int `json:"selEnd,omitempty"`
	SelLen    *int `json:"selLen,omitempty"`
	SelLenMin *int `json:"selLenMin,omitempty"`
	TextLen   *int `json:"textLen,omitempty"`
	Mode      *int `json:"mode,omitempty"`
}

type Step struct {
	Label  string       `json:"label,omitempty"`
	Keys   []Key        `json:"keys,omitempty"`
	Repeat int          `json:"repeat,omitempty"`
	Expect *StateExpect `json:"expect,omitempty"`
}

type Scenario struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Width   int    `json:"width,omitempty"` // harness wrap width in pixels; 0 = default
	Steps   []Step `json:"steps"`
}

type EditorState struct {
	Cursor      int    `json:"cursor"`
	SelStart    int    `json:"selStart"`
	SelEnd      int    `json:"selEnd"`
	TextLen     int    `json:"textLen"`
	Mode        int    `json:"mode"`
	IsLobby     int    `json:"isLobby"`
	CurrentFile string `json:"currentFile"`
}

type Harness struct {
	base       string
	host       string
	port       int
	verbose    bool
	noPrepare  bool
	reloadPoll time.Duration
}

func main() {
	host := flag.String("host", "127.0.0.1", "tablet host")
	port := flag.Int("port", 8000, "server port")
	scenario := flag.String("scenario", "", "run one scenario by exact name")
	match := flag.String("match", "", "run scenarios whose name contains this substring")
	list := flag.Bool("list", false, "list scenario names")
	verbose := flag.Bool("v", false, "verbose step output")
	unit := flag.Bool("unit", false, "run translate unit tests only (no device)")
	fast := flag.Bool("fast", false, "shorter key/step pauses for dev iteration")
	noPrepare := flag.Bool("no-prepare", false, "skip sandbox prepare (reuse open buffer; same scenario only)")
	reportMD := flag.String("report-md", "", "write markdown results table to this path")
	flag.Parse()

	if *fast {
		keyPause = fastKeyMs * time.Millisecond
		stepPause = fastStepMs * time.Millisecond
	}

	if *unit {
		fmt.Println("Run: go test -run TestTranslate ./...")
		os.Exit(0)
	}

	h := &Harness{
		base:       fmt.Sprintf("http://%s:%d", *host, *port),
		host:       *host,
		port:       *port,
		verbose:    *verbose,
		noPrepare:  *noPrepare,
		reloadPoll: 200 * time.Millisecond,
	}
	if *fast {
		h.reloadPoll = 100 * time.Millisecond
	}

	names := scenarioNames()
	if *list {
		for _, n := range names {
			fmt.Println(n)
		}
		return
	}

	if *scenario != "" && *match != "" {
		fmt.Fprintln(os.Stderr, "use -scenario or -match, not both")
		os.Exit(2)
	}

	var run []Scenario
	switch {
	case *scenario != "":
		sc, ok := findScenario(*scenario)
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown scenario %q\n", *scenario)
			os.Exit(2)
		}
		run = []Scenario{sc}
	case *match != "":
		var ok bool
		run, ok = findScenariosByPrefix(*match)
		if !ok {
			fmt.Fprintf(os.Stderr, "no scenarios match %q\n", *match)
			os.Exit(2)
		}
	default:
		run = AllScenarios()
	}

	if h.verbose {
		fmt.Println("mode: sandbox-prepare (no editor restart)")
	}

	runStarted := time.Now()
	modeLabel := "sandbox-prepare (single session)"

	var results []scenarioResult
	for _, sc := range run {
		res := h.runScenarioTimed(sc)
		results = append(results, res)
		switch res.Kind {
		case outcomePass:
			fmt.Printf("PASS %s (%.1fs)\n", sc.Name, res.Duration.Seconds())
		case outcomePrepareFail:
			fmt.Fprintf(os.Stderr, "PREPARE_FAIL %s (%.1fs): %s\n", sc.Name, res.Duration.Seconds(), res.Err)
		default:
			fmt.Fprintf(os.Stderr, "FAIL %s (%.1fs): %s\n", sc.Name, res.Duration.Seconds(), res.Err)
		}
	}
	if report := formatContaminationReport(results); report != "" {
		fmt.Fprint(os.Stderr, report)
	}
	if *reportMD != "" {
		meta := runMeta{
			StartedAt:     runStarted,
			Target:        fmt.Sprintf("%s:%d", *host, *port),
			Mode:          modeLabel,
			Fast:          *fast,
			ScenarioCount: len(run),
		}
		md := formatMarkdownReport(meta, results)
		if err := os.WriteFile(*reportMD, []byte(md), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "report-md: %v\n", err)
		} else {
			fmt.Printf("report: %s\n", *reportMD)
		}
	}
	failed := 0
	for _, r := range results {
		if r.Kind != outcomePass {
			failed++
		}
	}
	if failed > 0 {
		os.Exit(1)
	}
}

func (h *Harness) runScenarioTimed(sc Scenario) scenarioResult {
	start := time.Now()
	res := h.runScenarioTracked(sc)
	res.Duration = time.Since(start)
	return res
}

func (h *Harness) runScenarioTracked(sc Scenario) scenarioResult {
	if h.noPrepare {
		if err := h.RunScenario(sc); err != nil {
			return scenarioResult{Name: sc.Name, Kind: outcomeFail, Err: err.Error()}
		}
		return scenarioResult{Name: sc.Name, Kind: outcomePass}
	}
	retries, err := h.prepareWithRetry(sc)
	if err != nil {
		return scenarioResult{
			Name:            sc.Name,
			Kind:            outcomePrepareFail,
			Err:             err.Error(),
			PrepareRecovered: retries > 0,
		}
	}
	if err := h.RunScenario(sc); err != nil {
		health := h.notePostScenarioHealth(sc.Name)
		return scenarioResult{Name: sc.Name, Kind: outcomeFail, Err: err.Error(), PrepareRecovered: retries > 0, HealthNotes: health}
	}
	return scenarioResult{Name: sc.Name, Kind: outcomePass, PrepareRecovered: retries > 0}
}

// prepareWithRetry sandbox-resets the live editor without quitting.
func (h *Harness) prepareWithRetry(sc Scenario) (retries int, err error) {
	const attempts = 5
	var last error
	for i := 0; i < attempts; i++ {
		if err := h.sandboxPrepare(sc); err == nil {
			return retries, nil
		} else {
			last = err
			retries = i
		}
		if i+1 < attempts {
			time.Sleep(time.Duration(i+1) * 150 * time.Millisecond)
		}
	}
	return retries, last
}

func (h *Harness) sandboxPrepare(sc Scenario) error {
	if err := h.writeNote(sc.Content); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if err := h.ensureHarnessEditor(); err != nil {
		return fmt.Errorf("ensure editor: %w", err)
	}
	st, err := h.queryState()
	if err != nil {
		return fmt.Errorf("state: %w", err)
	}
	if st.IsLobby == 1 || st.CurrentFile != harnessNote {
		if err := h.editorCmd("harnessopen", harnessNote, 0); err != nil {
			return fmt.Errorf("harnessopen: %w", err)
		}
	}
	if err := h.editorCmd("harnessprepare", "", sc.Width); err != nil {
		return fmt.Errorf("harnessprepare: %w", err)
	}
	if err := h.verifyPrepareState(sc.Content); err != nil {
		return err
	}
	return nil
}

func (h *Harness) notePostScenarioHealth(scenario string) []string {
	st, err := h.queryState()
	if err != nil {
		msg := fmt.Sprintf("editor unreachable after fail: %v", err)
		fmt.Fprintf(os.Stderr, "  HEALTH %s: %s\n", scenario, msg)
		return []string{msg}
	}
	var notes []string
	if st.IsLobby == 1 {
		msg := "editor in lobby after fail"
		fmt.Fprintf(os.Stderr, "  HEALTH %s: %s\n", scenario, msg)
		notes = append(notes, msg)
	}
	if st.Mode != 1 {
		msg := fmt.Sprintf("not in edit mode (mode=%d) after fail", st.Mode)
		fmt.Fprintf(os.Stderr, "  HEALTH %s: %s\n", scenario, msg)
		notes = append(notes, msg)
	}
	return notes
}

func (h *Harness) RunScenario(sc Scenario) error {
	ws, err := h.dialWS()
	if err != nil {
		return fmt.Errorf("websocket: %w", err)
	}
	defer ws.Close()

	for i, step := range sc.Steps {
		label := step.Label
		if label == "" {
			label = fmt.Sprintf("step %d", i+1)
		}
		repeat := step.Repeat
		if repeat <= 0 {
			repeat = 1
		}
		for r := 0; r < repeat; r++ {
			for _, k := range step.Keys {
				if err := h.sendKey(ws, k); err != nil {
					return fmt.Errorf("%s: send %s: %w", label, k.Name, err)
				}
			}
		}
		if len(step.Keys) > 0 {
			time.Sleep(stepPause)
		}
		if step.Expect != nil {
			st, err := h.queryState()
			if err != nil {
				return fmt.Errorf("%s: state: %w", label, err)
			}
			if h.verbose {
				b, _ := json.Marshal(st)
				fmt.Printf("  %s: got %s\n", label, b)
			}
			if err := matchExpect(st, *step.Expect); err != nil {
				return fmt.Errorf("%s: %w", label, err)
			}
		}
	}
	return nil
}

func (h *Harness) verifyPrepareState(content string) error {
	st, err := h.queryState()
	if err != nil {
		return fmt.Errorf("post-prepare state: %w", err)
	}
	if st.TextLen != len(content) {
		return fmt.Errorf("textLen want %d got %d", len(content), st.TextLen)
	}
	if st.Cursor != 0 || st.SelStart != 0 || st.SelEnd != 0 {
		return fmt.Errorf("cursor/selection not clean: %v", st)
	}
	if st.Mode != 1 {
		return fmt.Errorf("want edit mode 1 got %d", st.Mode)
	}
	if st.IsLobby == 1 {
		return fmt.Errorf("editor in lobby")
	}
	if st.CurrentFile != harnessNote {
		return fmt.Errorf("currentFile want %q got %q", harnessNote, st.CurrentFile)
	}
	return nil
}

func (h *Harness) ensureHarnessEditor() error {
	st, err := h.queryState()
	if err == nil && st.Mode == 1 && st.IsLobby == 0 {
		return nil
	}
	if err := h.openNote(); err != nil {
		return err
	}
	deadline := time.Now().Add(12 * time.Second)
	for time.Now().Before(deadline) {
		st, err := h.queryState()
		if err == nil && st.Mode == 1 && st.IsLobby == 0 {
			return nil
		}
		time.Sleep(h.reloadPoll)
	}
	return fmt.Errorf("editor not ready after open")
}

func (h *Harness) editorCmd(c, name string, width int) error {
	body := map[string]interface{}{"c": c}
	if name != "" {
		body["name"] = name
	}
	if width > 0 {
		body["w"] = width
	}
	raw, _ := json.Marshal(body)
	code, err := h.post("/api/test/editor-cmd", raw)
	if err != nil {
		return err
	}
	if code != 200 {
		return fmt.Errorf("editor-cmd %s HTTP %d", c, code)
	}
	time.Sleep(h.reloadPoll)
	return nil
}

// writeNote upserts harness content without deleting the file or quitting the editor.
func (h *Harness) writeNote(content string) error {
	get, err := harnessHTTP.Get(h.base + "/api/notes/" + harnessNote)
	if err != nil {
		return err
	}
	io.Copy(io.Discard, get.Body) //nolint:errcheck
	status := get.StatusCode
	get.Body.Close()
	if status == 200 {
		return h.putNoteContent(content)
	}
	if status != 404 {
		return fmt.Errorf("read note HTTP %d", status)
	}
	body, _ := json.Marshal(map[string]string{
		"name":    strings.TrimSuffix(harnessNote, ".md"),
		"content": content,
	})
	code, err := h.post("/api/notes", body)
	if err != nil {
		return err
	}
	if code == 200 || code == 201 {
		return nil
	}
	if code == 409 {
		return h.putNoteContent(content)
	}
	return fmt.Errorf("create note HTTP %d", code)
}

func (h *Harness) putNoteContent(content string) error {
	get, err := harnessHTTP.Get(h.base + "/api/notes/" + harnessNote)
	if err != nil {
		return err
	}
	io.Copy(io.Discard, get.Body) //nolint:errcheck
	etag := get.Header.Get("ETag")
	status := get.StatusCode
	get.Body.Close()
	if status != 200 {
		return fmt.Errorf("read note HTTP %d", status)
	}
	body, _ := json.Marshal(map[string]string{"content": content})
	req, err := http.NewRequest(http.MethodPut, h.base+"/api/notes/"+harnessNote, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if etag != "" {
		req.Header.Set("If-Match", etag)
	}
	resp, err := harnessHTTP.Do(req)
	if err != nil {
		return err
	}
	io.Copy(io.Discard, resp.Body) //nolint:errcheck
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("put note HTTP %d", resp.StatusCode)
	}
	return nil
}

func (h *Harness) openNote() error {
	body, _ := json.Marshal(map[string]string{"name": harnessNote})
	for attempt := 0; attempt < 3; attempt++ {
		code, err := h.post("/api/open", body)
		if err == nil && code == 200 {
			return nil
		}
		time.Sleep(time.Second)
		if err != nil && attempt < 2 {
			continue
		}
		if err != nil {
			return err
		}
		if code != 200 {
			return fmt.Errorf("open HTTP %d", code)
		}
	}
	return fmt.Errorf("open failed after retries")
}

func (h *Harness) queryState() (EditorState, error) {
	resp, err := harnessHTTP.Get(h.base + "/api/test/editor-state")
	if err != nil {
		return EditorState{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return EditorState{}, fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	var st EditorState
	if err := json.NewDecoder(resp.Body).Decode(&st); err != nil {
		return EditorState{}, err
	}
	return st, nil
}

func (h *Harness) post(path string, body []byte) (int, error) {
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	req, err := http.NewRequest(http.MethodPost, h.base+path, r)
	if err != nil {
		return 0, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := harnessHTTP.Do(req)
	if err != nil {
		return 0, err
	}
	io.Copy(io.Discard, resp.Body) //nolint:errcheck
	resp.Body.Close()
	return resp.StatusCode, nil
}

func (h *Harness) dialWS() (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", h.host, h.port), Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	return conn, err
}

func (h *Harness) sendKey(ws *websocket.Conn, k Key) error {
	ev := map[string]interface{}{
		"type": "key",
		"key":  k.Name,
	}
	if k.Shift {
		ev["shift"] = true
	}
	if k.Ctrl {
		ev["ctrl"] = true
	}
	if k.Alt {
		ev["alt"] = true
	}
	if k.Meta {
		ev["meta"] = true
	}
	if err := ws.WriteJSON(ev); err != nil {
		return err
	}
	time.Sleep(keyPause)
	return nil
}

func matchExpect(got EditorState, exp StateExpect) error {
	var errs []string
	check := func(name string, want *int, have int) {
		if want == nil {
			return
		}
		if *want != have {
			errs = append(errs, fmt.Sprintf("%s want %d got %d", name, *want, have))
		}
	}
	check("cursor", exp.Cursor, got.Cursor)
	if exp.CursorMin != nil && got.Cursor < *exp.CursorMin {
		errs = append(errs, fmt.Sprintf("cursorMin want >= %d got %d", *exp.CursorMin, got.Cursor))
	}
	if exp.CursorMax != nil && got.Cursor > *exp.CursorMax {
		errs = append(errs, fmt.Sprintf("cursorMax want <= %d got %d", *exp.CursorMax, got.Cursor))
	}
	check("selStart", exp.SelStart, got.SelStart)
	check("selEnd", exp.SelEnd, got.SelEnd)
	check("textLen", exp.TextLen, got.TextLen)
	check("mode", exp.Mode, got.Mode)
	if exp.SelLen != nil {
		have := got.selLen()
		if *exp.SelLen != have {
			errs = append(errs, fmt.Sprintf("selLen want %d got %d", *exp.SelLen, have))
		}
	}
	if exp.SelLenMin != nil {
		have := got.selLen()
		if have < *exp.SelLenMin {
			errs = append(errs, fmt.Sprintf("selLenMin want >= %d got %d", *exp.SelLenMin, have))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s; state=%v", strings.Join(errs, "; "), got)
	}
	return nil
}

func intp(v int) *int { return &v }

func (s EditorState) selLen() int {
	if s.SelStart == s.SelEnd {
		return 0
	}
	if s.SelStart < s.SelEnd {
		return s.SelEnd - s.SelStart
	}
	return s.SelStart - s.SelEnd
}
