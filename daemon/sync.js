// sync.js -- disk drift UX + sync status display.
// GitHub reconcile runs on Writerdeck-server (Go); phone supplies token via POST /api/sync/token.
import { state } from './state.js';

var _onNotesChanged = function() {};
var _onBannerChange = function() {};

export var syncConfigured = false;

export function initSync(opts) {
  _onNotesChanged = opts.onNotesChanged || function() {};
  _onBannerChange = opts.onBannerChange || function() {};
}

export function refreshSyncStatus() {
  return fetch('/api/sync/status', { credentials: 'same-origin' })
    .then(function(r) { return r.ok ? r.json() : null; })
    .then(function(data) {
      if (!data) return null;
      syncConfigured = !!data.configured;
      updateSyncBannerFromState(data);
      updateSyncStatusLines(data);
      return data;
    })
    .catch(function() { return null; });
}

export function updateSyncBannerFromState(status) {
  var el = document.getElementById('sync-banner');
  if (!el) return;
  if (state.syncOn && status && !status.configured) {
    el.innerHTML = '\u26a0 GitHub sync is on \u2014 add your token in <strong>Setup</strong>.';
    el.style.display = 'block';
  } else if (status && status.lastError) {
    el.innerHTML = '\u26a0 ' + status.lastError +
      ' \u2014 renew token in <strong>Setup</strong>.';
    el.style.display = 'block';
  } else {
    el.style.display = 'none';
  }
  _onBannerChange();
}

function updateSyncStatusLines(data) {
  if (!data) return;
  var els = document.querySelectorAll('.sync-status-line');
  if (data.syncOn && !data.configured) {
    for (var i = 0; i < els.length; i++) {
      els[i].textContent = 'Token not on tablet \u2014 use Save & verify below';
      els[i].style.color = '#b45309';
    }
    return;
  }
  var text = data.lastSyncAgo ? 'Last synced: ' + data.lastSyncAgo : 'Never synced';
  for (var j = 0; j < els.length; j++) {
    els[j].textContent = text;
    els[j].style.color = '#888';
  }
}

// waitForSyncIdle: poll until background reconcile finishes (token verify runs async on tablet).
export function waitForSyncIdle(opts) {
  opts = opts || {};
  var deadline = Date.now() + (opts.timeoutMs || 90000);
  var sawSyncing = false;
  var baseline = opts.baselineLastSync || 0;
  return new Promise(function(resolve) {
    function tick() {
      refreshSyncStatus().then(function(s) {
        if (!s) { resolve(null); return; }
        if (s.syncing) sawSyncing = true;
        if (s.lastError) { resolve(s); return; }
        if (!s.syncing && (sawSyncing || s.lastSyncAt > baseline)) {
          resolve(s); return;
        }
        if (Date.now() > deadline) { resolve(s); return; }
        setTimeout(tick, 500);
      });
    }
    setTimeout(tick, 300);
  });
}

function showClashBanner(noteName, copyName) {
  var el = document.getElementById('clash-banner');
  if (!el) return;
  el.innerHTML = '<strong>Sync clash:</strong> \u201c' + noteName + '\u201d was also edited on GitHub. ' +
    'Your tablet version is now in \u201c' + copyName + '\u201d; ' +
    '\u201c' + noteName + '\u201d now holds the GitHub version. Review both, delete the one you don\u2019t want.';
  el.style.display = 'block';
  setTimeout(function() { el.style.display = 'none'; _onBannerChange(); }, 30000);
  _onBannerChange();
}

export function showSyncClash(noteName, copyName) {
  showClashBanner(noteName.replace(/\.md$/, ''), copyName.replace(/\.md$/, ''));
  _onNotesChanged();
}

export function recordEditorDiskBaseline(filename) {
  if (!filename) {
    state.editorDiskHash = '';
    return Promise.resolve();
  }
  return fetch('/api/notes/' + encodeURIComponent(filename), { credentials: 'same-origin' })
    .then(function(r) { return r.ok ? r.text() : null; })
    .then(function(t) {
      state.editorDiskHash = t !== null ? strHash(t) : '';
    })
    .catch(function() {});
}

function hideDriftBanner() {
  var el = document.getElementById('drift-banner');
  if (el) el.style.display = 'none';
  _onBannerChange();
}

function showDriftBanner(filename) {
  var el = document.getElementById('drift-banner');
  if (!el) return;
  var label = filename.replace(/\.md$/, '');
  el.innerHTML = '<strong>Disk changed:</strong> \u201c' + label +
    '\u201d was updated on disk while open on the tablet. ' +
    '<button type="button" id="drift-reload-btn">Reload on tablet</button> ' +
    'or keep editing (unsaved buffer wins on save).';
  el.style.display = 'block';
  var btn = document.getElementById('drift-reload-btn');
  if (btn) {
    btn.onclick = function(e) {
      e.stopPropagation();
      fetch('/api/reload', { method: 'POST', credentials: 'same-origin' })
        .then(function(r) {
          if (!r.ok) {
            alert('Could not reload \u2014 is the note still open on the tablet?');
            return;
          }
          hideDriftBanner();
          return recordEditorDiskBaseline(filename);
        })
        .catch(function() { alert('Could not reach server.'); });
    };
  }
  _onBannerChange();
}

export function checkDiskDrift(filename) {
  if (!filename || !state.editorDiskHash) return Promise.resolve();
  return fetch('/api/notes/' + encodeURIComponent(filename), { credentials: 'same-origin' })
    .then(function(r) { return r.ok ? r.text() : null; })
    .then(function(t) {
      if (t === null) return;
      if (strHash(t) !== state.editorDiskHash) showDriftBanner(filename);
    })
    .catch(function() {});
}

export function notifyDiskChanged(filename) {
  if (!filename) return;
  if (state.tabletOpenNote && filename !== state.tabletOpenNote) return;
  checkDiskDrift(filename);
}

function strHash(s) {
  var h = 5381;
  for (var i = 0; i < s.length; i++) { h = ((h << 5) + h + s.charCodeAt(i)) | 0; }
  return String(h >>> 0);
}
