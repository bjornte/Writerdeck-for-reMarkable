// Writerdeck-server — see main.go for overview.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// LAN use; no auth in Phase 3.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// logEvery controls how often the terse (non-verbose) log prints a running
// key count. Per-key translation detail is gated behind -v (for keymap
// debugging); by default the log stays quiet -- a periodic count plus a
// per-session total is enough to confirm keys are flowing, without flooding
// the device log with one line per keystroke.
const logEvery = 25

// --- WS broadcast hub ---
//
// Every connected browser is registered as a wsClient. A dedicated writer
// goroutine per client owns all conn.Write calls (gorilla WS forbids
// concurrent writers). broadcast() fans out a server-push message to all
// clients; sends are non-blocking so a slow/dead client cannot stall the
// caller.
type wsClient struct {
	conn *websocket.Conn
	send chan []byte
}

var (
	wsClientsMu sync.Mutex
	wsClients   = make(map[*wsClient]bool)

	syncAckMu sync.Mutex
	syncAckCh chan struct{} // set while power-sleep waits for browser reconcile
)

// signalSyncAck unblocks sleepForPower after the phone browser finishes GitHub sync.
func signalSyncAck() {
	syncAckMu.Lock()
	ch := syncAckCh
	syncAckCh = nil
	syncAckMu.Unlock()
	if ch != nil {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func beginSyncWait() {
	syncAckMu.Lock()
	syncAckCh = make(chan struct{}, 1)
	syncAckMu.Unlock()
}

func waitSyncAck(timeout time.Duration) {
	syncAckMu.Lock()
	ch := syncAckCh
	syncAckMu.Unlock()
	if ch == nil {
		return
	}
	select {
	case <-ch:
		fmt.Fprintln(os.Stderr, "writerdeck-server: sync ack received")
	case <-time.After(timeout):
		fmt.Fprintln(os.Stderr, "writerdeck-server: sync ack timeout -- proceeding")
	}
	syncAckMu.Lock()
	if syncAckCh == ch {
		syncAckCh = nil
	}
	syncAckMu.Unlock()
}

// broadcastOpenEdit tells phone clients which note the tablet editor holds open.
func broadcastOpenEdit(name string) {
	if name == "" {
		return
	}
	msg, _ := json.Marshal(struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}{"openedit", name})
	broadcast(msg)
}

// broadcast pushes msg to every registered browser client.
func broadcast(msg []byte) {
	wsClientsMu.Lock()
	defer wsClientsMu.Unlock()
	for c := range wsClients {
		select {
		case c.send <- msg:
		default:
		}
	}
}

var (
	needTokenMu     sync.Mutex
	needTokenLast   time.Time
	needTokenMinGap = 8 * time.Second
)

func sendNeedToken(client *wsClient) {
	if !syncEng.needsBrowserToken() {
		return
	}
	msg, _ := json.Marshal(struct {
		Type string `json:"type"`
	}{"needtoken"})
	select {
	case client.send <- msg:
	default:
	}
}

// maybeBroadcastNeedToken asks connected browsers to POST a saved GitHub token
// when sync is on but tablet RAM has none (e.g. after server restart).
func maybeBroadcastNeedToken() {
	if !syncEng.needsBrowserToken() {
		return
	}
	wsClientsMu.Lock()
	n := len(wsClients)
	wsClientsMu.Unlock()
	if n == 0 {
		return
	}
	needTokenMu.Lock()
	if time.Since(needTokenLast) < needTokenMinGap {
		needTokenMu.Unlock()
		return
	}
	needTokenLast = time.Now()
	needTokenMu.Unlock()
	msg, _ := json.Marshal(struct {
		Type string `json:"type"`
	}{"needtoken"})
	broadcast(msg)
}

// notifyTabletCrud queues the op and tells connected browsers to refresh; server syncs to GitHub.
func notifyTabletCrud(op, name, oldName string) {
	if name != "" {
		if p := notesSafe(name); p != "" {
			name = filepath.Base(p)
		}
	}
	if oldName != "" {
		if p := notesSafe(oldName); p != "" {
			oldName = filepath.Base(p)
		}
	}
	enqueuePendingSync(op, name, oldName)
	msg, _ := json.Marshal(struct {
		Type    string `json:"type"`
		Op      string `json:"op"`
		Name    string `json:"name"`
		OldName string `json:"oldName,omitempty"`
	}{"tabletcrud", op, name, oldName})
	broadcast(msg)
	pushLobbyInfo()
	syncEng.trySyncAfterCrud(op, name, oldName)
}

// maybeBroadcastDiskChanged notifies phone browsers when disk was written for the open note.
func maybeBroadcastDiskChanged(name string) {
	base := filepath.Base(name)
	currentNoteMu.Lock()
	open := currentNote != "" && currentNote == base
	currentNoteMu.Unlock()
	if !open {
		return
	}
	msg, _ := json.Marshal(struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}{"diskchanged", base})
	broadcast(msg)
}

// currentNote is the basename (.md) of the note the editor currently has open.
// Protected by currentNoteMu. Set by openHandler and editor {"t":"open"} reports;
// cleared by watchHomeButton, session.end(), and the DELETE handler on a match.
var (
	currentNoteMu sync.Mutex
	currentNote   string
)

func wsHandler(ec *editorConn, verbose bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authMu.Lock()
		required := pinRequired
		tok := authToken
		authMu.Unlock()
		if required {
			cookie, err := r.Cookie("writerdeck_token")
			if err != nil || cookie.Value != tok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: WS upgrade error: %v\n", err)
			return
		}
		// Register in the broadcast hub; writer goroutine owns all conn writes.
		client := &wsClient{conn: conn, send: make(chan []byte, 8)}
		wsClientsMu.Lock()
		wsClients[client] = true
		wsClientsMu.Unlock()
		defer func() {
			wsClientsMu.Lock()
			delete(wsClients, client)
			wsClientsMu.Unlock()
			close(client.send) // signals writer goroutine to drain and exit
			conn.Close()
		}()
		go func() {
			for msg := range client.send {
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					conn.Close() // read loop will see error and exit
					return
				}
			}
		}()
		// Tell new clients which note the tablet editor holds open (edit lease).
		currentNoteMu.Lock()
		openNote := currentNote
		currentNoteMu.Unlock()
		if openNote != "" {
			msg, _ := json.Marshal(struct {
				Type string `json:"type"`
				Name string `json:"name"`
			}{"openedit", openNote})
			select {
			case client.send <- msg:
			default:
			}
		}
		sendNeedToken(client)
		remote := r.RemoteAddr
		fmt.Fprintf(os.Stderr, "writerdeck-server: client connected %s\n", remote)
		var keys int
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Fprintf(os.Stderr, "writerdeck-server: client disconnected %s: %v (%d keys forwarded)\n", remote, err, keys)
				return
			}
			var ev wsMsg
			if err := json.Unmarshal(msg, &ev); err != nil || ev.Type != "key" {
				continue
			}
			line := translate(ev)
			if line == nil {
				continue
			}
			keys++
			if verbose {
				fmt.Fprintf(os.Stderr, "writerdeck-server: key=%q -> %s\n", ev.Key, line)
			} else if keys%logEvery == 0 {
				fmt.Fprintf(os.Stderr, "writerdeck-server: forwarded %d keys\n", keys)
			}
			ec.write(line)
		}
	}
}
