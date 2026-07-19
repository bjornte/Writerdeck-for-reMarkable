// notes-ui.js — typing view, paste modal, download-from-tablet offer.
import { state } from './state.js';
import { updateSyncBannerFromState, refreshSyncStatus } from './sync.js';
import { deps } from './deps.js';
import { grab, sendPaste, applyMode } from './connection.js';

var currentTypingFile = '';
var pendingDownloadName = '';

var IDLE_BODY_HTML =
  '<p>Open notes on the tablet Files tab.</p>' +
  '<p>Type here — keys go to the tablet.</p>';

var TYPING_BODY_HTML =
  '<p>Type here — your words appear on e-ink.</p>' +
  '<p>Press Home on the tablet when done.</p>';

var LOBBY_INPUT_LABELS = {
  'new': 'Type the new note name on the tablet keyboard.',
  'rename': 'Type the new name on the tablet keyboard.',
  'new-encrypted': 'Type the encrypted note name on the tablet keyboard.',
  'confirm-delete': 'Delete on tablet: tap Delete or Cancel (Enter / Esc also work).',
  'pin': 'Type the private PIN digits — they go to the tablet.',
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

function resetTypingBody() {
  document.getElementById('typing-body').innerHTML = TYPING_BODY_HTML;
  setPasteEnabled(true);
}

export function clearRemoteKeys() {
  state.remoteKeys = '';
  var banner = document.getElementById('remote-keys-banner');
  if (banner) banner.style.display = 'none';
  if (!state.typingMode || !currentTypingFile) {
    document.getElementById('typing-body').innerHTML = IDLE_BODY_HTML;
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
  document.getElementById('typing-note').textContent = 'Writerdeck';
  document.getElementById('typing-body').innerHTML = IDLE_BODY_HTML;
  setPasteEnabled(false);
  document.getElementById('typing').style.display = 'block';
  state.typingMode = true;
  applyMode();
}

// followTabletOpen: tablet opened a note (openedit) -- mirror it on the phone so
// keystrokes forward immediately. Skips when already typing that file, or when an
// overlay (PIN, settings, paste) has focus.
export function followTabletOpen(filename) {
  if (!filename) return;
  if (state.typingMode && currentTypingFile === filename) return;
  showTypingView(filename.replace(/\.md(\.enc)?$/, ''), filename);
}

export function showReadKeyView(filename) {
  if (!filename) return;
  state.remoteKeys = 'read';
  state.tabletOpenNote = filename;
  currentTypingFile = '';
  state.typingMode = true;
  document.getElementById('typing-note').textContent = displayName(filename);
  document.getElementById('typing-body').innerHTML =
    '<p>Reading on the tablet.</p>' +
    '<p>Press Esc on your keyboard to switch to edit mode.</p>';
  setPasteEnabled(false);
  document.getElementById('typing').style.display = 'block';
  applyMode();
}

export function showLobbyKeyView(mode) {
  if (!mode) {
    clearRemoteKeys();
    showIdleKeyboardView();
    return;
  }
  state.remoteKeys = 'lobby';
  state.typingMode = true;
  currentTypingFile = '';
  document.getElementById('typing-note').textContent = 'Writerdeck';
  document.getElementById('typing-body').innerHTML = IDLE_BODY_HTML;
  setPasteEnabled(false);
  document.getElementById('typing').style.display = 'block';
  var banner = document.getElementById('remote-keys-banner');
  banner.textContent = LOBBY_INPUT_LABELS[mode] || 'Keys go to the tablet.';
  banner.style.display = 'block';
  applyMode();
  if (deps.updateBannerOffset) deps.updateBannerOffset();
}

// showTypingView: show the "typing on the tablet" panel for an open note.
export function showTypingView(noteName, filename) {
  clearRemoteKeys();
  document.getElementById('typing-note').textContent = noteName;
  resetTypingBody();
  document.getElementById('typing').style.display = 'block';
  currentTypingFile = filename || '';
  state.tabletOpenNote = filename || '';
  state.typingMode = true;
  applyMode();
}

// hideTypingView: tablet left the note — stay on the keyboard shell.
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
var observeReady = false;

export function applyObserveEnabled(on) {
  state.observeEnabled = !!on;
  var btn = document.getElementById('typing-observe');
  if (!btn) return;
  if (!state.observeEnabled) {
    btn.style.display = 'none';
    setObserveBanner(false, 0, false);
    return;
  }
  btn.style.display = '';
  setObserveButton(observing);
  setObserveBanner(observing, 0, observeReady && !observing);
}

function setObserveBanner(active, steps, ready) {
  var banner = document.getElementById('observe-banner');
  if (!banner) return;
  if (!state.observeEnabled) {
    banner.style.display = 'none';
    banner.textContent = '';
  } else if (active) {
    banner.textContent = 'Observing' + (steps ? ' \u00b7 ' + steps + ' keys' : '') +
      ' \u2014 Stop when the bug shows.';
    banner.style.display = 'block';
  } else if (ready) {
    banner.textContent = 'Observation ready on the tablet. In Cursor chat, say you found a bug \u2014 it will be pulled automatically.';
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
  observeReady = !!(data.ready || (data.hasExport && !data.active));
  setObserveButton(!!data.active);
  setObserveBanner(!!data.active, data.steps || 0, observeReady);
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
