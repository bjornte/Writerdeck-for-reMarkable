// Writerdeck-server -- reMarkable Wi-Fi typewriter daemon.
//
// Serves a WebSocket on 0.0.0.0:8000/ws and forwards received key events
// to the patched Writerdeck editor over a local Unix socket (/run/Writerdeck.sock).
//
// Architecture (two layers, two parsers):
//   Browser --WebSocket--> Writerdeck-server --Unix socket--> Writerdeck (patched keywriter)
//
// WebSocket message (JSON, from browser keydown):
//   {"type":"key","key":"<KeyboardEvent.key>"}
//
// Editor-feed wire format (NDJSON to keywriter's naive C++ parser):
//   {"t":"text","cp":<unicode-codepoint-int>}   -- single printable char
//   {"t":"key","k":"Escape|Return|Backspace|Tab|ArrowUp|ArrowDown|ArrowLeft|ArrowRight"}
//   {"t":"cmd","c":"home|open|notedeleted|noterenamed|..."}  -- editor commands
// Browser <- Writerdeck-server: {"type":"openedit","name":"<file>.md"} on tablet open/rename
//   {"type":"tabletcrud","op":"createnote|deletenote|renamenote","name":"...","oldName":"..."}
//     on tablet Lobby Files CRUD — server mirrors to GitHub when sync is configured.
//   {"type":"diskchanged","name":"<file>.md"} when disk is written under an open editor (slice 8).
//
// Integer codepoints are escaping-proof: JSON special chars in typed text
// can never corrupt the naive C++ substring parser (see socket-inject.patch).
//
// Usage on the device:
//   /home/root/Writerdeck-server               # serve on :8000
//   /home/root/Writerdeck-server -v            # verbose key logging
//   /home/root/Writerdeck-server --selftest    # one-shot hello world+Return
//   /home/root/Writerdeck-server --port 9000   # custom port
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

// Static assets embedded at compile time; all served with Cache-Control: no-store.
// app.css, app.js (+ connection/notes-ui/panels/deps modules), state.js, sync.js.
//
//go:embed index.html
var indexHTML []byte

//go:embed app.css
var appCSS []byte

//go:embed app.js
var appJS []byte

//go:embed deps.js
var depsJS []byte

//go:embed connection.js
var connectionJS []byte

//go:embed notes-ui.js
var notesUIJS []byte

//go:embed panels.js
var panelsJS []byte

//go:embed state.js
var stateJS []byte

//go:embed sync.js
var syncJS []byte

func serveJS(w http.ResponseWriter, data []byte) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.Write(data) //nolint:errcheck
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(indexHTML) //nolint:errcheck
}
func main() {
	port := flag.Int("port", defaultPort, "WebSocket server port")
	doSelftest := flag.Bool("selftest", false, "send hardcoded hello world+Return and exit")
	verbose := flag.Bool("v", false, "verbose: log every translated key (default logs only a periodic count)")
	editorPath := flag.String("editor", "", "path to editor launch script; rmkbd spawns it as a child and owns its lifecycle (supervisor mode, used by systemd unit)")
	flag.StringVar(&notesDirPath, "notes-dir", notesDirPath, "directory for .md notes (default: /home/root/Writerdeck-user-documents; override for local dev)")
	flag.StringVar(&settingsFilePath, "settings-file", settingsFilePath, "path to settings JSON (default /home/root/.Writerdeck/settings.json; override for local dev)")
	flag.Parse()

	loadSettings()
	startSyncBackground()

	// Determine PIN length from persisted settings (loaded above).
	var bootPinLen int
	switch curSettings.PinDigits {
	case "4":
		bootPinLen = 4
	case "none":
		bootPinLen = 0
	default:
		bootPinLen = 6 // "6" or any unrecognised value
	}
	authMu.Lock()
	authPIN = generatePIN(bootPinLen)
	authToken = generateToken()
	pinRequired = bootPinLen > 0
	authMu.Unlock()
	if authPIN != "" {
		fmt.Fprintf(os.Stderr, "writerdeck-server: PIN is %s (will be shown on tablet Lobby; for now read from this log)\n", authPIN)
	} else {
		fmt.Fprintln(os.Stderr, "writerdeck-server: no PIN required (pinDigits=none)")
	}

	ec := &editorConn{}
	globalEC = ec
	go dialLoop(ec)
	go watchLobbyIP()

	if *doSelftest {
		selftest(ec)
		return
	}

	addr := fmt.Sprintf(":%d", *port)
	fmt.Fprintf(os.Stderr, "writerdeck-server: serving capture page on http://<device-ip>%s/\n", addr)
	fmt.Fprintf(os.Stderr, "writerdeck-server: serving WebSocket on %s/ws\n", addr)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/app.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(appCSS) //nolint:errcheck
	})
	http.HandleFunc("/app.js", func(w http.ResponseWriter, r *http.Request) { serveJS(w, appJS) })
	http.HandleFunc("/deps.js", func(w http.ResponseWriter, r *http.Request) { serveJS(w, depsJS) })
	http.HandleFunc("/connection.js", func(w http.ResponseWriter, r *http.Request) { serveJS(w, connectionJS) })
	http.HandleFunc("/notes-ui.js", func(w http.ResponseWriter, r *http.Request) { serveJS(w, notesUIJS) })
	http.HandleFunc("/panels.js", func(w http.ResponseWriter, r *http.Request) { serveJS(w, panelsJS) })
	http.HandleFunc("/state.js", func(w http.ResponseWriter, r *http.Request) { serveJS(w, stateJS) })
	http.HandleFunc("/sync.js", func(w http.ResponseWriter, r *http.Request) { serveJS(w, syncJS) })
	http.HandleFunc("/ws", wsHandler(ec, *verbose))
	http.HandleFunc("/api/vault/status", vaultStatusHandler)
	http.HandleFunc("/api/pin", pinHandler)
	http.HandleFunc("/api/launch", launchHandler)
	http.HandleFunc("/api/open", openHandler)
	http.HandleFunc("/api/notes", notesListHandler)
	http.HandleFunc("/api/notes/", notesItemHandler)
	http.HandleFunc("/api/settings", settingsHandler)
	http.HandleFunc("/api/lobby", lobbyHandler) // pre-auth: reveals PIN on e-ink only
	http.HandleFunc("/api/status", statusHandler)
	http.HandleFunc("/api/sync/ack", syncAckHandler)
	http.HandleFunc("/api/sync/pending", pendingSyncHandler)
	http.HandleFunc("/api/sync/pending/clear", pendingClearHandler)
	http.HandleFunc("/api/sync/token", syncTokenHandler)
	http.HandleFunc("/api/sync/status", syncStatusHandler)
	http.HandleFunc("/api/sync/run", syncRunHandler)
	http.HandleFunc("/api/reload", reloadHandler)
	http.HandleFunc("/api/flush-save", flushSaveHandler)
	http.HandleFunc("/api/test/reset", testResetHandler)
	http.HandleFunc("/api/test/editor-state", testEditorStateHandler)
	http.HandleFunc("/api/test/home", testHomeHandler)
	http.HandleFunc("/api/test/tablet-req", testTabletReqHandler)
	http.HandleFunc("/api/test/editor-cmd", testEditorCmdHandler)

	if *editorPath != "" {
		// Supervisor mode: rmkbd is always-on; editor sessions are on-demand.
		// xochitl stop/start happens per session in Go (session.start/end).
		// ExecStopPost=start xochitl in the unit stays as a safety net.
		activeSess = &session{editorPath: *editorPath, ec: ec}

		// Always-on Home watcher: loops for rmkbd's lifetime.
		go watchHomeButton(activeSess, ec)
		go watchUSBKeyboardForLaunch(activeSess)
		go watchUSBKeyboardHotplug(activeSess)

		// HTTP server always-on in the background.
		go func() {
			if err := http.ListenAndServe(addr, nil); err != nil {
				fmt.Fprintf(os.Stderr, "writerdeck-server: HTTP server: %v\n", err)
			}
		}()

		// Reconcile: kill any stray keywriter from a previous crash so two
		// editors don't fight for the framebuffer on startup.
		fmt.Fprintln(os.Stderr, "writerdeck-server: reconcile: killing any stray Writerdeck editor")
		exec.Command("sh", "-c", "for p in $(pidof Writerdeck 2>/dev/null); do kill $p 2>/dev/null; done").Run() //nolint:errcheck
		time.Sleep(500 * time.Millisecond)

		// Auto-launch first session: power-on = typewriter (unchanged behaviour).
		fmt.Fprintln(os.Stderr, "writerdeck-server: auto-launching editor session on boot")
		if err := activeSess.start(); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: auto-launch failed: %v\n", err)
		}

		// Block until SIGTERM/SIGINT; gracefully end any active session first.
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		sig := <-sigCh
		fmt.Fprintf(os.Stderr, "writerdeck-server: signal %v received\n", sig)
		flushEditorSave()
		if activeSess.isActive() {
			fmt.Fprintln(os.Stderr, "writerdeck-server: ending active session before exit")
			activeSess.quit()
		}
		fmt.Fprintln(os.Stderr, "writerdeck-server: exiting (ExecStopPost safety net restarts xochitl if needed)")
		os.Exit(0)
	}

	// Stand-alone mode: dev/test scripts manage keywriter separately.
	// Still watch the home button: sends a single quit to ec (one-shot).
	go watchHomeButton(nil, ec)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintf(os.Stderr, "writerdeck-server: server error: %v\n", err)
		os.Exit(1)
	}
}
