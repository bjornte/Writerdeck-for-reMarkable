// app.js — phone UI bootstrap: wire modules and event listeners.
// connection.js — WebSocket + key capture; notes-ui.js — file manager;
// panels.js — PIN screen, Preferences, Notes sync setup; sync.js — GitHub engine.
import { state } from './state.js';
import { initSync } from './sync.js';
import { deps } from './deps.js';
import { initConnection, connect, grab, applyMode } from './connection.js';
import {
  loadNotes, showList, hideTypingView, createNote, uploadFile,
  showPasteModal, hidePasteModal, submitPaste
} from './notes-ui.js';
import {
  showSettings, hideSettings, showSync, hideSync, showPinScreen,
  checkAuthAndInit, submitPIN, updateBannerOffset, wireOverlayDismiss
} from './panels.js';

deps.loadNotes = loadNotes;
deps.hideTypingView = hideTypingView;
deps.showPinScreen = showPinScreen;
deps.connect = connect;

initConnection();

window.addEventListener('load', function () {
  initSync({ onNotesChanged: loadNotes, onBannerChange: updateBannerOffset });
  updateBannerOffset();
  document.getElementById('new-btn').addEventListener('click', createNote);
  document.getElementById('upload-btn').addEventListener('click', function(e) {
    e.stopPropagation();
    var fi = document.getElementById('file-input');
    fi.value = '';
    fi.click();
  });
  document.getElementById('file-input').addEventListener('change', function() {
    uploadFile(this.files[0]);
    this.value = '';
  });
  document.getElementById('typing-paste').addEventListener('click', function(e) {
    e.stopPropagation(); showPasteModal();
  });
  // Clipboard-source label: "from here" is always correct (whatever device you
  // hold). Upgrade to "from phone" only on a high-confidence phone UA (iPhone/
  // iPod, or Android + "Mobile"). iPad is excluded on purpose: since iPadOS 13 it
  // reports as desktop Safari, so it stays the safe "from here".
  if (/iPhone|iPod/.test(navigator.userAgent) ||
      (/Android/.test(navigator.userAgent) && /Mobile/.test(navigator.userAgent))) {
    document.getElementById('typing-paste').textContent = 'Paste from phone';
  }
  document.getElementById('paste-submit').addEventListener('click', function(e) {
    e.stopPropagation(); submitPaste();
  });
  document.getElementById('paste-cancel').addEventListener('click', function(e) {
    e.stopPropagation(); hidePasteModal();
  });
  document.getElementById('preview-back').addEventListener('click', showList);
  document.getElementById('typing-back').addEventListener('click', hideTypingView);
  document.getElementById('pin-btn').addEventListener('click', submitPIN);
  document.getElementById('pin-input').addEventListener('keydown', function(e) {
    if (e.key === 'Enter') { e.stopPropagation(); submitPIN(); }
  });
  document.getElementById('lobby-btn').addEventListener('click', function(e) {
    e.stopPropagation();
    var msgEl = document.getElementById('lobby-msg');
    msgEl.textContent = 'Asking the tablet to show the PIN\u2026';
    fetch('/api/lobby', { method: 'POST' })
      .then(function(r) {
        if (r.status === 429) {
          msgEl.textContent = 'Just a moment, try again shortly.';
          return;
        }
        if (!r.ok) {
          msgEl.textContent = 'Could not reach the tablet.';
          return;
        }
        msgEl.textContent = 'PIN should now appear on the tablet.';
      })
      .catch(function() {
        msgEl.textContent = 'Could not reach the server.';
      });
  });
  document.getElementById('settings-btn').addEventListener('click', function(e) {
    e.stopPropagation(); showSettings();
  });
  document.getElementById('settings-done').addEventListener('click', function(e) {
    e.stopPropagation(); hideSettings();
  });
  document.getElementById('settings-close').addEventListener('click', function(e) {
    e.stopPropagation(); hideSettings();
  });
  document.getElementById('sync-btn').addEventListener('click', function(e) {
    e.stopPropagation(); showSync();
  });
  document.getElementById('sync-done').addEventListener('click', function(e) {
    e.stopPropagation(); hideSync();
  });
  document.getElementById('sync-close').addEventListener('click', function(e) {
    e.stopPropagation(); hideSync();
  });
  wireOverlayDismiss('settings-screen', 'settings-box', hideSettings);
  wireOverlayDismiss('sync-screen', 'sync-box', hideSync);
  document.addEventListener('keydown', function(e) {
    if (e.key !== 'Escape') return;
    if (document.getElementById('settings-screen').style.display === 'flex') hideSettings();
    else if (document.getElementById('sync-screen').style.display === 'flex') hideSync();
  });
  applyMode();
  checkAuthAndInit();
});
