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
	"os/exec"
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
	if err := run(base, wsURL, *host); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("PASS")
}

func run(base, wsURL, host string) error {
	before, err := getSettings(base)
	if err != nil {
		return fmt.Errorf("settings before: %w", err)
	}
	defer func() {
		if rerr := sshRestoreSettings(host, before); rerr != nil {
			fmt.Fprintf(os.Stderr, "  cleanup warning: %v\n", rerr)
		}
	}()
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
	// Files -> Keyboard -> Sync -> Settings (Tab three times).
	for i := 0; i < 3; i++ {
		if err := ws.WriteJSON(map[string]string{"type": "key", "key": "Tab"}); err != nil {
			return fmt.Errorf("Tab to Settings: %w", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(200 * time.Millisecond)
	fmt.Println("  sent Tab×3 to Settings")

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
	afterPin, err := readPinDigitsSSH(host)
	if err != nil {
		return fmt.Errorf("settings after pin (ssh): %w", err)
	}
	if afterPin != targetPin {
		return fmt.Errorf("pinDigits want %q got %q", targetPin, afterPin)
	}
	fmt.Printf("  setpindigits: %q ok\n", targetPin)

	restored, err := readSettingsSSH(host)
	if err != nil {
		return fmt.Errorf("final settings: %w", err)
	}
	if restored.ReadFont != targetFont || restored.PinDigits != targetPin {
		return fmt.Errorf("final mismatch: got %+v want font=%q pin=%q", restored, targetFont, targetPin)
	}
	fmt.Println("  font and pin persisted on disk")
	return nil
}

func readSettingsSSH(host string) (settingsResp, error) {
	cmd := exec.Command("ssh", "-o", "BatchMode=no", "-o", "ConnectTimeout=10",
		"root@"+host, "cat /home/root/.Writerdeck/settings.json")
	out, err := cmd.Output()
	if err != nil {
		return settingsResp{}, err
	}
	var disk settingsResp
	if json.Unmarshal(out, &disk) != nil {
		return settingsResp{}, fmt.Errorf("parse settings.json")
	}
	return disk, nil
}

func readPinDigitsSSH(host string) (string, error) {
	s, err := readSettingsSSH(host)
	return s.PinDigits, err
}

func sshRestoreSettings(host string, want settingsResp) error {
	script := fmt.Sprintf(
		`set -e
f=%q
p=%q
sed -i "s/\"readFont\":\"[^\"]*\"/\"readFont\":\"$f\"/" /home/root/.Writerdeck/settings.json
sed -i "s/\"pinDigits\":\"[^\"]*\"/\"pinDigits\":\"$p\"/" /home/root/.Writerdeck/settings.json
systemctl restart writerdeck
`, want.ReadFont, want.PinDigits)
	cmd := exec.Command("ssh", "-o", "BatchMode=no", "-o", "ConnectTimeout=10",
		"root@"+host, "sh", "-c", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(out))
	}
	time.Sleep(4 * time.Second)
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
