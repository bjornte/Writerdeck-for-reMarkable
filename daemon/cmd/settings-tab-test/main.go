// settings-tab-test — tablet Settings tab socket ops + lobby page navigation.
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

const httpTimeout = 30 * time.Second

var client = &http.Client{Timeout: httpTimeout}

type editorState struct {
	IsLobby int `json:"isLobby"`
}

type settingsResp struct {
	ReadFont  string `json:"readFont"`
	PinDigits string `json:"pinDigits"`
}

func main() {
	host := flag.String("host", "127.0.0.1", "tablet host")
	port := flag.Int("port", 8000, "server port")
	flag.Parse()

	base := fmt.Sprintf("http://%s:%d", *host, *port)
	wsURL := fmt.Sprintf("ws://%s:%d/ws", *host, *port)

	fmt.Printf("=== settings-tab-test  host=%s ===\n", *host)
	if err := run(base, wsURL); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("PASS")
}

func run(base, wsURL string) error {
	before, err := getSettings(base)
	if err != nil {
		return fmt.Errorf("settings before: %w", err)
	}
	fmt.Printf("  before: readFont=%q pinDigits=%q\n", before.ReadFont, before.PinDigits)

	if code, err := post(base+"/api/lobby", nil); err != nil {
		return err
	} else if code != 200 {
		return fmt.Errorf("lobby: status %d", code)
	}
	time.Sleep(500 * time.Millisecond)

	st, err := queryState(base)
	if err != nil {
		return fmt.Errorf("lobby state: %w", err)
	}
	if st.IsLobby != 1 {
		return fmt.Errorf("after lobby: isLobby want 1 got %d", st.IsLobby)
	}
	fmt.Println("  lobby: isLobby=1")

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("websocket: %w", err)
	}
	defer ws.Close()
	if err := ws.WriteJSON(map[string]string{"type": "key", "key": "5"}); err != nil {
		return fmt.Errorf("key 5: %w", err)
	}
	time.Sleep(300 * time.Millisecond)
	fmt.Println("  sent Settings tab key (5)")

	targetFont := "Literata"
	if before.ReadFont == targetFont {
		targetFont = "EB Garamond"
	}
	if err := tabletReq(base, "setreadfont", targetFont); err != nil {
		return fmt.Errorf("setreadfont: %w", err)
	}
	afterFont, err := getSettings(base)
	if err != nil {
		return fmt.Errorf("settings after font: %w", err)
	}
	if afterFont.ReadFont != targetFont {
		return fmt.Errorf("readFont want %q got %q", targetFont, afterFont.ReadFont)
	}
	fmt.Printf("  setreadfont: %q ok\n", targetFont)

	targetPin := "4"
	if before.PinDigits == "4" {
		targetPin = "6"
	}
	if err := tabletReq(base, "setpindigits", targetPin); err != nil {
		return fmt.Errorf("setpindigits: %w", err)
	}
	afterPin, err := getSettings(base)
	if err != nil {
		return fmt.Errorf("settings after pin: %w", err)
	}
	if afterPin.PinDigits != targetPin {
		return fmt.Errorf("pinDigits want %q got %q", targetPin, afterPin.PinDigits)
	}
	fmt.Printf("  setpindigits: %q ok\n", targetPin)

	if err := tabletReq(base, "setreadfont", before.ReadFont); err != nil {
		return fmt.Errorf("restore font: %w", err)
	}
	if err := tabletReq(base, "setpindigits", before.PinDigits); err != nil {
		return fmt.Errorf("restore pin: %w", err)
	}
	restored, err := getSettings(base)
	if err != nil {
		return fmt.Errorf("settings restored: %w", err)
	}
	if restored.ReadFont != before.ReadFont || restored.PinDigits != before.PinDigits {
		return fmt.Errorf("restore mismatch: got %+v want %+v", restored, before)
	}
	fmt.Println("  restored original settings")
	return nil
}

func getSettings(base string) (settingsResp, error) {
	var out settingsResp
	resp, err := client.Get(base + "/api/settings")
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return out, fmt.Errorf("status %d: %s", resp.StatusCode, b)
	}
	err = json.NewDecoder(resp.Body).Decode(&out)
	return out, err
}

func tabletReq(base, op, name string) error {
	body, _ := json.Marshal(map[string]string{"op": op, "name": name})
	code, err := post(base+"/api/test/tablet-req", body)
	if err != nil {
		return err
	}
	if code != 200 {
		return fmt.Errorf("status %d", code)
	}
	time.Sleep(200 * time.Millisecond)
	return nil
}

func queryState(base string) (editorState, error) {
	var st editorState
	resp, err := client.Get(base + "/api/test/editor-state")
	if err != nil {
		return st, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return st, fmt.Errorf("status %d: %s", resp.StatusCode, b)
	}
	err = json.NewDecoder(resp.Body).Decode(&st)
	return st, err
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
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, nil
}
