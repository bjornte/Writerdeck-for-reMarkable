// connection.js — WebSocket, status bar, key capture.
import { state } from './state.js';
import {
  updateSyncBannerFromState, refreshSyncStatus, showSyncClash,
  recordEditorDiskBaseline, checkDiskDrift, notifyDiskChanged, respondToNeedToken
} from './sync.js';
import { deps } from './deps.js';

var RETRY_MS = 2000;
var ECHO_MAX = 300;

var dot  = document.getElementById('dot');
var msg  = document.getElementById('msg');
var echo = document.getElementById('echo');
var trap = document.getElementById('trap');
var buf = '';
var ws;
// state.typingMode: false=Browse (list, no capture), true=Type (capture + echo).
export function applyMode() {
  document.getElementById('foot').style.display = state.typingMode ? 'block' : 'none';
  document.body.classList.toggle('typing-dark', state.typingMode);
  if (state.typingMode) { grab(); }
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
        checkDiskDrift(data.openNote);
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

function send(key, shift, ctrl, alt, meta) {
  if (!ws || ws.readyState !== 1 /* OPEN */) { return; }
  ws.send(JSON.stringify({
    type: 'key', key: key,
    shift: !!shift, ctrl: !!ctrl, alt: !!alt, meta: !!meta
  }));
  appendEcho(key);
}

function onKey(e) {
  // Only capture keys in Type mode; in Browse mode let the browser handle them.
  if (!state.typingMode) { return; }
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
  return false;
}

function grab() {
  if (!state.typingMode) { return; }
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
    grab();
    deps.loadNotes();
    refreshSyncStatus();
  };

  ws.onclose = function () {
    wsReady = false;
    updateConnectionBar();
    setTimeout(connect, RETRY_MS);
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
      } else if (data.type === 'exitedit') {
        state.tabletOpenNote = '';
        state.editorDiskHash = '';
        if (state.typingMode) { deps.hideTypingView(); }
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
  if (ws && ws.readyState === WebSocket.OPEN) { ws.close(); }
}

export { connect, send, grab, startStatusPoll, stopStatusPoll, setStatus, overlayUp };
