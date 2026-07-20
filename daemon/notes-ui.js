// notes-ui.js — typing view, paste modal, download-from-tablet offer.
import { state } from './state.js';
import { updateSyncBannerFromState, refreshSyncStatus } from './sync.js';
import { deps } from './deps.js';
import { grab, sendPaste, applyMode } from './connection.js';

var currentTypingFile = '';
var pendingDownloadName = '';

var IDLE_BODY_HTML = '';

// Writerdeck quit — tablet is back on stock xochitl; Launch button only.
var CLOSED_BODY_HTML = '';

var TYPING_BODY_HTML = '';

var editorActive = null; // null until first /api/status (avoid Launch flash)

// Only modes that still need a phone banner. Name/PIN prompts stay silent.
var LOBBY_INPUT_LABELS = {
  'confirm-delete': 'Delete on tablet: tap Delete or Cancel (Enter / Esc also work).',
  'no-keyboard': 'Scan the QR on the tablet (or open the URL) to connect this phone as a keyboard.'
};

// Auth + sync only — the phone no longer lists notes.
export function loadNotes() {
  fetch('/api/notes')
    .then(function(r) {
      if (r.status === 401) { deps.showPinScreen(); return null; }
      if (!r.ok) throw new Error('server error');
      return r.json();
    })
    .then(function(notes) {
      if (notes !== null) {
        refreshSyncStatus().then(function(s) { updateSyncBannerFromState(s); });
      }
    })
    .catch(function(err) { console.error('loadNotes failed:', err); });
}

function displayName(filename) {
  return (filename || '').replace(/\.md(\.enc)?$/, '');
}

// Paste inserts at the tablet cursor — only while a note is open for edit, not Files/Lobby/read.
function pasteAllowed() {
  return !!currentTypingFile && state.remoteKeys !== 'lobby' && state.remoteKeys !== 'read';
}

function setPasteEnabled(on) {
  var btn = document.getElementById('typing-paste');
  if (!btn) return;
  btn.disabled = !on;
  btn.setAttribute('aria-disabled', on ? 'false' : 'true');
  btn.style.display = '';
}

function setTypingGuide(html) {
  var guide = document.getElementById('typing-guide');
  if (guide) guide.innerHTML = html;
}

function isIdleShell() {
  return !currentTypingFile && state.remoteKeys !== 'read' && state.remoteKeys !== 'lobby';
}

function setLaunchVisible(on) {
  var btn = document.getElementById('typing-launch');
  if (!btn) return;
  btn.style.display = on ? 'inline-block' : 'none';
  if (!on) {
    btn.disabled = false;
    btn.textContent = 'Launch Writerdeck';
  }
}

function updateIdleChrome() {
  if (!isIdleShell()) {
    setLaunchVisible(false);
    return;
  }
  // Only show Launch when status has confirmed the editor is down.
  if (editorActive === false) {
    setTypingGuide(CLOSED_BODY_HTML);
    setLaunchVisible(true);
  } else {
    setTypingGuide(IDLE_BODY_HTML);
    setLaunchVisible(false);
  }
  setPasteEnabled(false);
}

// Called from /api/status polls and WebSocket open/exit so the shell
// matches Files-tab vs stock UI without waiting for a note open.
export function applyEditorActive(active) {
  editorActive = !!active;
  if (isIdleShell()) updateIdleChrome();
  else setLaunchVisible(false);
}

export function launchWriterdeck(e) {
  if (e) e.stopPropagation();
  var btn = document.getElementById('typing-launch');
  if (btn) {
    btn.disabled = true;
    btn.textContent = 'Launching\u2026';
  }
  fetch('/api/launch', { method: 'POST', credentials: 'same-origin' })
    .then(function(r) {
      if (r.status === 401) {
        deps.showPinScreen();
        return null;
      }
      // 409 = already running — treat as success and refresh chrome.
      if (r.ok || r.status === 409) {
        applyEditorActive(true);
        return null;
      }
      return r.text().then(function(t) {
        throw new Error(t || ('HTTP ' + r.status));
      });
    })
    .catch(function(err) {
      if (btn) {
        btn.disabled = false;
        btn.textContent = 'Launch Writerdeck';
      }
      alert('Could not launch Writerdeck: ' +
        (err && err.message ? err.message : 'error'));
    });
}

function resetTypingBody() {
  setTypingGuide(TYPING_BODY_HTML);
  setPasteEnabled(true);
  setLaunchVisible(false);
}

export function clearRemoteKeys() {
  state.remoteKeys = '';
  var banner = document.getElementById('remote-keys-banner');
  if (banner) banner.style.display = 'none';
  if (!state.typingMode || !currentTypingFile) {
    updateIdleChrome();
  } else {
    resetTypingBody();
  }
  applyMode();
}

// Default phone surface: keyboard ready, no file list.
export function showIdleKeyboardView() {
  clearRemoteKeys();
  currentTypingFile = '';
  state.tabletOpenNote = state.tabletOpenNote || '';
  document.getElementById('typing-note').textContent = '';
  updateIdleChrome();
  setPasteEnabled(false);
  document.getElementById('typing').style.display = 'flex';
  state.typingMode = true;
  applyMode();
}

// followTabletOpen: tablet opened a note (openedit) -- mirror it on the phone so
// keystrokes forward immediately. Skips when already typing that file, or when an
// overlay (PIN, settings, paste) has focus.
export function followTabletOpen(filename) {
  if (!filename) return;
  applyEditorActive(true);
  if (state.typingMode && currentTypingFile === filename) return;
  showTypingView(filename.replace(/\.md(\.enc)?$/, ''), filename);
}

export function showReadKeyView(filename) {
  if (!filename) return;
  applyEditorActive(true);
  state.remoteKeys = 'read';
  state.tabletOpenNote = filename;
  currentTypingFile = '';
  state.typingMode = true;
  document.getElementById('typing-note').textContent = '';
  setTypingGuide('');
  setPasteEnabled(false);
  setLaunchVisible(false);
  document.getElementById('typing').style.display = 'flex';
  applyMode();
}

export function showLobbyKeyView(mode) {
  applyEditorActive(true);
  if (!mode) {
    clearRemoteKeys();
    showIdleKeyboardView();
    return;
  }
  state.remoteKeys = 'lobby';
  state.typingMode = true;
  currentTypingFile = '';
  document.getElementById('typing-note').textContent = '';
  setTypingGuide(IDLE_BODY_HTML);
  setPasteEnabled(false);
  setLaunchVisible(false);
  document.getElementById('typing').style.display = 'flex';
  var banner = document.getElementById('remote-keys-banner');
  var label = LOBBY_INPUT_LABELS[mode];
  if (label) {
    banner.textContent = label;
    banner.style.display = 'block';
  } else {
    banner.textContent = '';
    banner.style.display = 'none';
  }
  applyMode();
  if (deps.updateBannerOffset) deps.updateBannerOffset();
}

// showTypingView: show the "typing on the tablet" panel for an open note.
export function showTypingView(noteName, filename) {
  applyEditorActive(true);
  clearRemoteKeys();
  document.getElementById('typing-note').textContent = '';
  resetTypingBody();
  document.getElementById('typing').style.display = 'flex';
  currentTypingFile = filename || '';
  state.tabletOpenNote = filename || '';
  state.typingMode = true;
  applyMode();
}

// hideTypingView: tablet left the note — stay on the keyboard shell.
// Caller sets applyEditorActive when the whole app quit (exitedit).
export function hideTypingView(e) {
  if (e) e.stopPropagation();
  currentTypingFile = '';
  showIdleKeyboardView();
}

export function showDownloadOffer(filename) {
  if (!filename) return;
  pendingDownloadName = filename;
  document.getElementById('download-name').textContent = displayName(filename);
  document.getElementById('download-modal').style.display = 'flex';
}

export function hideDownloadOffer() {
  pendingDownloadName = '';
  document.getElementById('download-modal').style.display = 'none';
  grab();
}

export function acceptDownloadOffer() {
  var name = pendingDownloadName;
  hideDownloadOffer();
  if (name) downloadNote(name);
}

export function downloadNote(filename) {
  var url = '/api/notes/' + encodeURIComponent(filename) + '/download';
  fetch(url, { credentials: 'same-origin' })
    .then(function(r) {
      if (r.status === 401) { deps.showPinScreen(); return null; }
      if (r.status === 423) {
        alert('Enter private PIN on tablet');
        return waitForVaultPIN().then(function(ok) {
          if (!ok) { alert('Timed out waiting for tablet PIN.'); return null; }
          return fetch(url, { credentials: 'same-origin' });
        });
      }
      if (!r.ok) { alert('Download failed.'); return null; }
      return r.blob();
    })
    .then(function(blob) {
      if (!blob) return;
      var a = document.createElement('a');
      var dlName = filename.replace(/\.enc$/, '');
      a.href = URL.createObjectURL(blob);
      a.download = dlName;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(a.href);
    })
    .catch(function() { alert('Could not reach server.'); });
}

function waitForVaultPIN() {
  return new Promise(function(resolve) {
    var deadline = Date.now() + 120000;
    function poll() {
      fetch('/api/vault/status', { credentials: 'same-origin' })
        .then(function(r) { return r.ok ? r.json() : null; })
        .then(function(st) {
          if (st && st.enabled && !st.locked) { resolve(true); return; }
          if (Date.now() > deadline) { resolve(false); return; }
          setTimeout(poll, 1500);
        })
        .catch(function() {
          if (Date.now() > deadline) resolve(false);
          else setTimeout(poll, 1500);
        });
    }
    poll();
  });
}

// ---- Paste-note modal (edit only — not Lobby Files) ----

export function showPasteModal() {
  if (!pasteAllowed()) return;
  var contentEl = document.getElementById('paste-content');
  contentEl.value = '';
  document.getElementById('paste-modal').style.display = 'flex';
  if (navigator.clipboard && navigator.clipboard.readText) {
    navigator.clipboard.readText()
      .then(function(text) {
        if (text) { contentEl.value = text; }
      })
      .catch(function() {});
  }
  contentEl.focus();
}

export function hidePasteModal() {
  document.getElementById('paste-modal').style.display = 'none';
}

export function submitPaste() {
  if (!pasteAllowed()) {
    hidePasteModal();
    return;
  }
  var content = document.getElementById('paste-content').value;
  if (!content) { alert('Nothing to insert \u2014 paste some text first.'); return; }
  hidePasteModal();
  // Dedicated paste message — server ignores it when no note is open (Lobby Files).
  sendPaste(content);
  grab();
}

// ---- Observation mode (bug demos for LLM) ----

var observing = false;

export function applyObserveEnabled(on) {
  state.observeEnabled = !!on;
  var btn = document.getElementById('typing-observe');
  if (!btn) return;
  if (!state.observeEnabled) {
    btn.style.display = 'none';
    setObserveBanner(false, 0);
    return;
  }
  btn.style.display = '';
  setObserveButton(observing);
  setObserveBanner(observing, 0);
}

function setObserveBanner(active, steps) {
  var banner = document.getElementById('observe-banner');
  if (!banner) return;
  if (!state.observeEnabled) {
    banner.style.display = 'none';
    banner.textContent = '';
  } else if (active) {
    banner.textContent = 'Observing' + (steps ? ' \u00b7 ' + steps + ' keys' : '') +
      ' \u2014 Stop when the bug shows.';
    banner.style.display = 'block';
  } else {
    banner.style.display = 'none';
    banner.textContent = '';
  }
  if (deps.updateBannerOffset) deps.updateBannerOffset();
}

function setObserveButton(active) {
  var btn = document.getElementById('typing-observe');
  if (!btn) return;
  observing = !!active;
  if (active) {
    btn.textContent = 'Stop observe';
    btn.classList.add('observe-on');
    btn.classList.remove('danger');
  } else {
    btn.textContent = 'Observe';
    btn.classList.remove('observe-on', 'danger');
  }
}

export function applyObserveStatus(data) {
  if (!data) return;
  if (typeof data.enabled === 'boolean') {
    applyObserveEnabled(data.enabled);
  }
  if (!state.observeEnabled) return;
  setObserveButton(!!data.active);
  setObserveBanner(!!data.active, data.steps || 0);
}

export function refreshObserveStatus() {
  if (!state.observeEnabled) {
    applyObserveEnabled(false);
    return Promise.resolve(null);
  }
  return fetch('/api/observe/status', { credentials: 'same-origin' })
    .then(function(r) { return r.ok ? r.json() : null; })
    .then(function(st) {
      if (st) applyObserveStatus(st);
      return st;
    })
    .catch(function() { return null; });
}

export function toggleObserve(e) {
  if (e) e.stopPropagation();
  if (!state.observeEnabled) return;
  if (observing) {
    fetch('/api/observe/stop', { method: 'POST', credentials: 'same-origin' })
      .then(function(r) {
        if (!r.ok) {
          return r.text().then(function(t) { throw new Error(t || 'stop failed'); });
        }
        return r.text();
      })
      .then(function() {
        applyObserveStatus({ active: false, ready: true, hasExport: true, steps: 0, enabled: true });
      })
      .catch(function(err) {
        alert('Could not stop observation: ' + (err && err.message ? err.message : 'error'));
        refreshObserveStatus();
      });
    return;
  }
  fetch('/api/observe/start', { method: 'POST', credentials: 'same-origin' })
    .then(function(r) {
      if (!r.ok) {
        return r.text().then(function(t) { throw new Error(t || 'start failed'); });
      }
      return r.json();
    })
    .then(function(st) {
      applyObserveStatus(st);
    })
    .catch(function(err) {
      alert('Could not start observation: ' + (err && err.message ? err.message : 'open a note on the tablet first'));
    });
}
