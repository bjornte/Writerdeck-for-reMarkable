// app.js — phone UI bootstrap: wire modules and event listeners.
// connection.js — WebSocket + key capture; notes-ui.js — typing / paste / download;
// panels.js — PIN screen, Notes sync setup overlays; sync.js — GitHub engine.
import { initSync } from './sync.js';
import { deps } from './deps.js';
import { initConnection, connect, applyMode } from './connection.js';
import {
  loadNotes, hideTypingView, followTabletOpen, showIdleKeyboardView,
  showPasteModal, hidePasteModal, submitPaste,
  showReadKeyView, showLobbyKeyView, clearRemoteKeys,
  showDownloadOffer, hideDownloadOffer, acceptDownloadOffer,
  toggleObserve, refreshObserveStatus, applyObserveStatus, applyObserveEnabled,
  applyEditorActive, launchWriterdeck
} from './notes-ui.js';
import {
  showSync, hideSync, showPinScreen,
  checkAuthAndInit, submitPIN, updateBannerOffset, wireOverlayDismiss
} from './panels.js';
import { loadI18n, t, applyPasteButtonLabel } from './i18n.js';

deps.loadNotes = loadNotes;
deps.hideTypingView = hideTypingView;
deps.followTabletOpen = followTabletOpen;
deps.showIdleKeyboardView = showIdleKeyboardView;
deps.showReadKeyView = showReadKeyView;
deps.showLobbyKeyView = showLobbyKeyView;
deps.clearRemoteKeys = clearRemoteKeys;
deps.applyEditorActive = applyEditorActive;
deps.showDownloadOffer = showDownloadOffer;
deps.showPinScreen = showPinScreen;
deps.connect = connect;
deps.applyObserveStatus = applyObserveStatus;
deps.applyObserveEnabled = applyObserveEnabled;
deps.refreshObserveStatus = refreshObserveStatus;
deps.updateBannerOffset = updateBannerOffset;

initConnection();

window.addEventListener('load', function () {
  loadI18n().then(function () {
    initSync({ onNotesChanged: loadNotes, onBannerChange: updateBannerOffset });
    updateBannerOffset();
    document.getElementById('typing-paste').addEventListener('click', function(e) {
      e.stopPropagation(); showPasteModal();
    });
    document.getElementById('typing-observe').addEventListener('click', toggleObserve);
    document.getElementById('typing-launch').addEventListener('click', launchWriterdeck);
    applyPasteButtonLabel();
    document.getElementById('paste-submit').addEventListener('click', function(e) {
      e.stopPropagation(); submitPaste();
    });
    document.getElementById('paste-cancel').addEventListener('click', function(e) {
      e.stopPropagation(); hidePasteModal();
    });
    document.getElementById('download-accept').addEventListener('click', function(e) {
      e.stopPropagation(); acceptDownloadOffer();
    });
    document.getElementById('download-cancel').addEventListener('click', function(e) {
      e.stopPropagation(); hideDownloadOffer();
    });
    document.getElementById('pin-btn').addEventListener('click', submitPIN);
    document.getElementById('pin-input').addEventListener('keydown', function(e) {
      if (e.key === 'Enter') { e.stopPropagation(); submitPIN(); }
    });
    document.getElementById('lobby-btn').addEventListener('click', function(e) {
      e.stopPropagation();
      var msgEl = document.getElementById('lobby-msg');
      msgEl.textContent = t('pin.asking');
      fetch('/api/lobby', { method: 'POST' })
        .then(function(r) {
          if (r.status === 429) {
            msgEl.textContent = t('pin.rateLimit');
            return;
          }
          if (!r.ok) {
            msgEl.textContent = t('pin.reachTablet');
            return;
          }
          msgEl.textContent = t('pin.appeared');
        })
        .catch(function() {
          msgEl.textContent = t('pin.reachServer');
        });
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
    wireOverlayDismiss('sync-screen', 'sync-box', hideSync);
    wireOverlayDismiss('download-modal', 'download-box', hideDownloadOffer);
    document.addEventListener('keydown', function(e) {
      if (e.key !== 'Escape') return;
      if (document.getElementById('sync-screen').style.display === 'flex') hideSync();
      if (document.getElementById('download-modal').style.display === 'flex') hideDownloadOffer();
    });
    applyMode();
    checkAuthAndInit();
    refreshObserveStatus();
  });
});
