// lobby-keyboard-test — lobby keys + Home from read must not quit Writerdeck.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	testNote    = "z-test-keyboard-harness.md"
	httpTimeout = 30 * time.Second
)

var client = &http.Client{Timeout: httpTimeout}

type editorState struct {
	TextLen int `json:"textLen"`
	Mode    int `json:"mode"`
	IsLobby int `json:"isLobby"`
}

func main() {
	host := flag.String("host", "127.0.0.1", "tablet host")
	port := flag.Int("port", 8000, "server port")
	note := flag.String("note", testNote, "note to open")
	flag.Parse()

	base := fmt.Sprintf("http://%s:%d", *host, *port)
	wsURL := fmt.Sprintf("ws://%s:%d/ws", *host, *port)

	fmt.Printf("=== lobby-keyboard-test  host=%s  note=%s ===\n", *host, *note)

	if err := run(base, wsURL, *note); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("PASS")
}

func run(base, wsURL, note string) error {
	st0, err := openNote(base, note)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	if st0.TextLen == 0 {
		return fmt.Errorf("after open: textLen=0")
	}
	fmt.Printf("  edit: textLen=%d\n", st0.TextLen)

	wsPre, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("websocket: %w", err)
	}
	if err := wsPre.WriteJSON(map[string]string{"type": "key", "key": "Escape"}); err != nil {
		wsPre.Close()
		return fmt.Errorf("escape: %w", err)
	}
	wsPre.Close()
	time.Sleep(400 * time.Millisecond)

	stRead, err := queryState(base)
	if err != nil {
		return fmt.Errorf("read mode: %w", err)
	}
	if stRead.Mode != 0 {
		return fmt.Errorf("after Esc: mode want 0 got %d", stRead.Mode)
	}
	fmt.Printf("  read: mode=%d\n", stRead.Mode)

	if code, err := post(base+"/api/test/home", nil); err != nil {
		return err
	} else if code != 200 {
		return fmt.Errorf("home from read HTTP %d", code)
	}
	time.Sleep(600 * time.Millisecond)

	stLobby, err := queryState(base)
	if err != nil {
		return fmt.Errorf("post-home-read: %w", err)
	}
	if stLobby.IsLobby != 1 {
		return fmt.Errorf("home from read: isLobby want 1 got %d", stLobby.IsLobby)
	}
	active, err := editorActive(base)
	if err != nil {
		return err
	}
	if !active {
		return fmt.Errorf("home from read quit Writerdeck")
	}
	fmt.Println("  home from read: still in session")

	body, _ := json.Marshal(map[string]string{"name": note})
	if code, err := post(base+"/api/open", body); err != nil {
		return err
	} else if code != 200 {
		return fmt.Errorf("re-open HTTP %d", code)
	}
	time.Sleep(400 * time.Millisecond)

	if code, err := post(base+"/api/test/home", nil); err != nil {
		return err
	} else if code != 200 {
		return fmt.Errorf("home HTTP %d", code)
	}
	time.Sleep(500 * time.Millisecond)

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("websocket: %w", err)
	}
	defer ws.Close()

	for _, k := range []string{"Tab", "Enter"} {
		if err := ws.WriteJSON(map[string]string{"type": "key", "key": k}); err != nil {
			return fmt.Errorf("key %s: %w", k, err)
		}
		time.Sleep(300 * time.Millisecond)
	}

	wantLen := st0.TextLen
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		st1, err := queryState(base)
		if err != nil {
			return fmt.Errorf("post-enter: %w", err)
		}
		if st1.TextLen >= wantLen {
			fmt.Printf("  after lobby Enter: textLen=%d\n", st1.TextLen)
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	st1, err := queryState(base)
	if err != nil {
		return err
	}
	return fmt.Errorf("lobby keyboard dead: textLen=%d want >=%d", st1.TextLen, wantLen)
}

func openNote(base, note string) (editorState, error) {
	body, _ := json.Marshal(map[string]string{"name": note})
	code, err := post(base+"/api/open", body)
	if err != nil {
		return editorState{}, err
	}
	if code != 200 {
		return editorState{}, fmt.Errorf("open HTTP %d", code)
	}
	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		st, err := queryState(base)
		if err == nil && st.TextLen > 0 {
			return st, nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return queryState(base)
}

func queryState(base string) (editorState, error) {
	resp, err := client.Get(base + "/api/test/editor-state")
	if err != nil {
		return editorState{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return editorState{}, fmt.Errorf("editor-state HTTP %d: %s", resp.StatusCode, bytes.TrimSpace(b))
	}
	var st editorState
	if err := json.NewDecoder(resp.Body).Decode(&st); err != nil {
		return editorState{}, err
	}
	return st, nil
}

func editorActive(base string) (bool, error) {
	resp, err := client.Get(base + "/api/status")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var st struct {
		EditorActive bool `json:"editorActive"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&st); err != nil {
		return false, err
	}
	return st.EditorActive, nil
}

func post(url string, body []byte) (int, error) {
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	req, err := http.NewRequest(http.MethodPost, url, r)
	if err != nil {
		return 0, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	io.Copy(io.Discard, resp.Body) //nolint:errcheck
	resp.Body.Close()
	return resp.StatusCode, nil
}
