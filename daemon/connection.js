// connection.js — WebSocket, status bar, key capture.
import { state } from './state.js';
import {
  updateSyncBannerFromState, refreshSyncStatus, showSyncClash,
  recordEditorDiskBaseline, notifyDiskChanged, respondToNeedToken
} from './sync.js';
import { deps } from './deps.js';

var RETRY_MS = 2000;
var ECHO_MAX = 300;

// Cursor (and Electron shells) load this page for agent checks; they must not
// count as a keyboard path for the Lobby tip. Real Safari/Chrome/Firefox do not match.
function isIdeBrowser() {
  try {
    if (typeof window !== 'undefined' && window.cursorBrowser) return true;
    var ua = (typeof navigator !== 'undefined' && navigator.userAgent) ? navigator.userAgent : '';
    if (ua.indexOf('Cursor/') !== -1 || ua.indexOf('Electron/') !== -1) return true;
  } catch (e) {}
  return false;
}

var dot  = document.getElementById('dot');
var msg  = document.getElementById('msg');
var echo = document.getElementById('echo');
var trap = document.getElementById('trap');
var buf = '';
var ws;
var allowReconnect = true;
// state.typingMode: true = keyboard shell (capture + echo). Default after auth.
// state.remoteKeys: '' | 'read' | 'lobby' -- still set for banners / Esc into edit.
function keysActive() {
  return state.typingMode || state.remoteKeys !== '';
}

export function applyMode() {
  document.getElementById('foot').style.display = keysActive() ? 'block' : 'none';
  document.body.classList.toggle('typing-dark', state.typingMode);
  if (keysActive()) { grab(); }
}

function setStatus(cls, text) {
  dot.className = cls;
  msg.textContent = text;
}

var statusTimer = null;
var STATUS_MS = 5000;
var STATUS_TIMEOUT_MS = 4000;
var statusPolling = false;
var tabletReachable = false;
var tabletInfo = null;
var wsReady = false;

function tabletExtras(data) {
  var parts = [];
  if (data.battery >= 0) {
    parts.push(data.battery + '%' + (data.charging ? ' +' : ''));
  }
  if (!data.wifi) { parts.push('no Wi-Fi'); }
  return parts.length ? ' \u00b7 ' + parts.join(' \u00b7 ') : '';
}

// One bar state: HTTP /api/status is the fast truth for reachability; the
// WebSocket only gates "ready to type" once the tablet is actually there.
function updateConnectionBar() {
  if (!statusPolling) {
    dot.className = '';
    msg.textContent = '';
    return;
  }
  if (!tabletReachable) {
    setStatus('off', 'Tablet offline');
    return;
  }
  if (wsReady) {
    setStatus('on', 'Connected' + tabletExtras(tabletInfo || {}));
    return;
  }
  setStatus('', 'Connecting\u2026');
}

function markTabletOffline() {
  tabletReachable = false;
  tabletInfo = null;
  updateConnectionBar();
  if (ws && ws.readyState === WebSocket.OPEN) { ws.close(); }
}

function refreshTabletStatus() {
  var ctrl = new AbortController();
  var timer = setTimeout(function() { ctrl.abort(); }, STATUS_TIMEOUT_MS);
  fetch('/api/status', { signal: ctrl.signal, credentials: 'same-origin' })
    .then(function(r) {
      clearTimeout(timer);
      return r.ok ? r.json() : null;
    })
    .then(function(data) {
      if (!data) { markTabletOffline(); return; }
      tabletReachable = true;
      tabletInfo = data;
      if (data.editorActive && data.openNote) {
        if (state.tabletOpenNote !== data.openNote) {
          state.tabletOpenNote = data.openNote;
          recordEditorDiskBaseline(data.openNote);
        } else if (!state.editorDiskHash) {
          recordEditorDiskBaseline(data.openNote);
        }
        // Do not poll-check drift here: tablet autosave changes disk under the open note and
        // would false-alarm "Disk changed". Real external writes arrive via WS diskchanged.
      }
      updateConnectionBar();
    })
    .catch(function() {
      clearTimeout(timer);
      markTabletOffline();
    });
}

function startStatusPoll() {
  statusPolling = true;
  refreshTabletStatus();
  if (statusTimer) clearInterval(statusTimer);
  statusTimer = setInterval(refreshTabletStatus, STATUS_MS);
}

function stopStatusPoll() {
  statusPolling = false;
  if (statusTimer) { clearInterval(statusTimer); statusTimer = null; }
  tabletReachable = false;
  tabletInfo = null;
  wsReady = false;
  updateConnectionBar();
}

function appendEcho(key) {
  if (key === 'Backspace') {
    buf = buf.slice(0, -1);
  } else if (key === 'Enter') {
    buf += '\n';
  } else if (key === 'Tab') {
    buf += '\t';
  } else if (key === 'Escape') {
    buf += '[Esc]';
  } else if (key.length === 1) {
    buf += key;
  }
  if (buf.length > ECHO_MAX) { buf = buf.slice(-ECHO_MAX); }
  echo.textContent = buf;
  // Keep the latest keystrokes in view -- the box clips its top, not its
  // bottom -- so the echo always reflects what you just typed.
  echo.scrollTop = echo.scrollHeight;
}

function send(key, shift, ctrl, alt, meta, action) {
  if (!ws || ws.readyState !== 1 /* OPEN */) { return; }
  var msg = {
    type: 'key', key: key,
    shift: !!shift, ctrl: !!ctrl, alt: !!alt, meta: !!meta
  };
  if (action) msg.action = action;
  ws.send(JSON.stringify(msg));
  if (!action) appendEcho(key);
}

// Bulk paste for the typing view. Server accepts only while a note is open (not Lobby Files).
export function sendPaste(text) {
  if (!ws || ws.readyState !== 1 /* OPEN */) { return; }
  if (!text) return;
  ws.send(JSON.stringify({ type: 'paste', text: text }));
  appendEcho('[paste]');
}

function onKey(e) {
  // Capture keys in Type mode and when the tablet asked for remote input (read/lobby).
  if (!keysActive()) { return; }
  // Pass browser-navigation shortcuts through so the page stays manageable
  // (Cmd/Ctrl + R=reload, T=new tab, W=close, N=new window, L=address bar).
  // Everything else -- including Ctrl/Cmd+C/V/X/A/Z/K and modifier+arrows --
  // is captured and forwarded to the tablet.
  var k = e.key.toLowerCase();
  if ((e.ctrlKey || e.metaKey) && (k==='r'||k==='t'||k==='w'||k==='n'||k==='l')) {
    return;
  }
  e.preventDefault();
  send(e.key, e.shiftKey, e.ctrlKey, e.altKey, e.metaKey);
  // Esc toggles edit/preview on key-up in Writerdeck. Socket inject no longer
  // auto-releases Escape (harness sends an explicit release), so the phone must too.
  if (e.key === 'Escape') {
    send(e.key, e.shiftKey, e.ctrlKey, e.altKey, e.metaKey, 'release');
  }
}

// overlayUp: true when the PIN screen or the paste modal is showing, so the
// tablet-capture logic (grab/onKey) must stand down and let the overlay's own
// inputs receive focus + keystrokes. Without this, opening the paste modal in
// the typing view would steal focus to the trap and forward keys to the tablet.
function overlayUp() {
  var ps = document.getElementById('pin-screen');
  if (ps && ps.style.display !== 'none') return true;
  var pm = document.getElementById('paste-modal');
  if (pm && pm.style.display === 'flex') return true;
  var dm = document.getElementById('download-modal');
  if (dm && dm.style.display === 'flex') return true;
  var ss = document.getElementById('sync-screen');
  if (ss && ss.style.display === 'flex') return true;
  return false;
}

function grab() {
  if (!keysActive()) { return; }
  if (overlayUp()) return;
  trap.focus();
}

// Primary listener on the trap textarea.
trap.addEventListener('keydown', onKey);

// Fallback: catch keydown on the whole document in case focus escapes.
// Skip entirely while the PIN screen is up so the PIN input can receive typed digits.
document.addEventListener('keydown', function (e) {
  if (overlayUp()) return;
  if (document.activeElement !== trap) { onKey(e); }
});

// Any click/tap on the page re-focuses the trap -- but not while the PIN screen is up.
document.addEventListener('click', function(e) {
  if (overlayUp()) return;
  grab();
});

function connect() {
  wsReady = false;
  updateConnectionBar();

  ws = new WebSocket('ws://' + window.location.host + '/ws');

  ws.onopen = function () {
    wsReady = true;
    updateConnectionBar();
    // Marks this tab as a real phone/laptop UI for the Lobby keyboard tip.
    // Cursor's embedded browser must not count (skips the tip otherwise).
    if (!isIdeBrowser()) {
      try { ws.send(JSON.stringify({ type: 'hello' })); } catch (e) {}
    }
    if (deps.showIdleKeyboardView) deps.showIdleKeyboardView();
    else grab();
    deps.loadNotes();
    refreshSyncStatus();
  };

  ws.onclose = function () {
    wsReady = false;
    updateConnectionBar();
    // Do not reconnect after pagehide/unload — leftover tabs were keeping
    // phoneConnected true and skipping the Lobby keyboard tip.
    if (allowReconnect) setTimeout(connect, RETRY_MS);
  };

  ws.onerror = function () {
    // onclose fires after onerror; nothing extra needed.
  };

  ws.onmessage = function (event) {
    try {
      var data = JSON.parse(event.data);
      if (data.type === 'openedit') {
        state.tabletOpenNote = data.name || '';
        recordEditorDiskBaseline(state.tabletOpenNote);
        if (state.tabletOpenNote && !overlayUp()) {
          deps.followTabletOpen(state.tabletOpenNote);
        }
      } else if (data.type === 'openread') {
        state.tabletOpenNote = data.name || '';
        if (state.tabletOpenNote && !overlayUp()) {
          deps.showReadKeyView(state.tabletOpenNote);
        }
      } else if (data.type === 'lobbyinput') {
        if (!overlayUp()) {
          deps.showLobbyKeyView(data.mode || '');
        }
      } else if (data.type === 'exitedit') {
        state.tabletOpenNote = '';
        state.editorDiskHash = '';
        if (deps.hideTypingView) deps.hideTypingView();
        else if (deps.clearRemoteKeys) deps.clearRemoteKeys();
      } else if (data.type === 'tabletcrud') {
        if (data.op === 'deletenote' && state.tabletOpenNote === data.name) {
          state.tabletOpenNote = '';
        } else if (data.op === 'renamenote' && data.oldName && state.tabletOpenNote === data.oldName) {
          state.tabletOpenNote = data.name || '';
        }
        deps.loadNotes();
      } else if (data.type === 'syncclash') {
        showSyncClash(data.note || '', data.copyName || '');
      } else if (data.type === 'diskchanged') {
        notifyDiskChanged(data.name || '');
      } else if (data.type === 'needtoken') {
        respondToNeedToken();
      } else if (data.type === 'vaultpingranted') {
        deps.loadNotes();
      } else if (data.type === 'downloadoffer') {
        if (data.name && deps.showDownloadOffer && !overlayUp()) {
          deps.showDownloadOffer(data.name);
        }
      } else if (data.type === 'observe') {
        if (deps.applyObserveStatus) {
          deps.applyObserveStatus({
            active: !!data.active,
            steps: data.steps || 0,
            ready: !!data.ready,
            hasExport: !!data.hasExport
          });
        }
      }
      // Unknown types are silently ignored -- forward-compatible.
    } catch (e) {}
  };
}

export function initConnection() {
  trap.addEventListener('keydown', onKey);
  document.addEventListener('keydown', function (e) {
    if (overlayUp()) return;
    if (document.activeElement !== trap) { onKey(e); }
  });
  document.addEventListener('click', function(e) {
    if (overlayUp()) return;
    grab();
  });
}
export function closeWebSocket() {
  allowReconnect = false;
  if (ws && ws.readyState === WebSocket.OPEN) { ws.close(); }
  ws = null;
}

window.addEventListener('pagehide', function () {
  allowReconnect = false;
  if (ws) {
    try { ws.close(); } catch (e) {}
    ws = null;
  }
});

window.addEventListener('pageshow', function () {
  allowReconnect = true;
});

export { connect, send, sendPaste, grab, startStatusPoll, stopStatusPoll, setStatus, overlayUp };
