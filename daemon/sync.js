// sync.js -- disk drift UX + sync status display.
// GitHub reconcile runs on Writerdeck-server (Go); phone supplies token via POST /api/sync/token.
import { state } from './state.js';
import { t, tf, formatSyncAgo } from './i18n.js';

var _onNotesChanged = function() {};
var _onBannerChange = function() {};

export var syncConfigured = false;
export var syncOffline = false;

var _tokenPushPromise = null;
var _tokenPullPromise = null;

export function ghToken() {
  return localStorage.getItem('ghToken') || '';
}

export function setSyncToken(token) {
  localStorage.setItem('ghToken', token);
}

export function clearSyncToken() {
  localStorage.removeItem('ghToken');
}

export function pullTokenFromTablet() {
  if (ghToken()) return Promise.resolve(false);
  if (_tokenPullPromise) return _tokenPullPromise;
  _tokenPullPromise = fetch('/api/sync/token', { credentials: 'same-origin' })
    .then(function(r) { return r.ok ? r.json() : null; })
    .then(function(data) {
      _tokenPullPromise = null;
      if (data && data.configured && data.token) {
        setSyncToken(data.token);
        return true;
      }
      return false;
    })
    .catch(function() {
      _tokenPullPromise = null;
      return false;
    });
  return _tokenPullPromise;
}

export function pushStoredTokenToTablet() {
  var token = ghToken();
  if (!token || !state.syncOn || !state.syncRepo) return Promise.resolve(false);
  if (_tokenPushPromise) return _tokenPushPromise;
  _tokenPushPromise = fetch('/api/sync/token', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    credentials: 'same-origin',
    body: JSON.stringify({ token: token })
  }).then(function(r) {
    _tokenPushPromise = null;
    if (r.status === 401) clearSyncToken();
    return r.ok;
  }).catch(function() {
    _tokenPushPromise = null;
    return false;
  });
  return _tokenPushPromise;
}

export function respondToNeedToken() {
  if (!ghToken()) return Promise.resolve(false);
  return fetchSyncStatus().then(function(data) {
    if (!data || !data.syncOn || !data.syncRepo) return false;
    state.syncOn = true;
    state.syncRepo = data.syncRepo;
    if (data.configured) return true;
    return pushStoredTokenToTablet().then(function(ok) {
      if (!ok) return false;
      return refreshSyncStatus().then(function() { return true; });
    });
  });
}

export function initSync(opts) {
  _onNotesChanged = opts.onNotesChanged || function() {};
  _onBannerChange = opts.onBannerChange || function() {};
}

function fetchSyncStatus() {
  return fetch('/api/sync/status', { credentials: 'same-origin' })
    .then(function(r) { return r.ok ? r.json() : null; });
}

function reportSyncOffline() {
  syncOffline = true;
  syncConfigured = false;
  updateSyncBannerFromState(null);
  var els = document.querySelectorAll('.sync-status-line');
  for (var i = 0; i < els.length; i++) {
    els[i].textContent = t('sync.statusOffline');
    els[i].style.color = '#e57373';
  }
}

export function refreshSyncStatus() {
  return fetchSyncStatus()
    .then(function(data) {
      if (!data) {
        reportSyncOffline();
        return null;
      }
      syncOffline = false;
      if (data.syncOn && !ghToken() && data.configured) {
        return pullTokenFromTablet().then(function(pulled) {
          if (!pulled) return data;
          return fetchSyncStatus();
        });
      }
      if (data.syncOn && !data.configured && ghToken()) {
        return pushStoredTokenToTablet().then(function(ok) {
          if (!ok) return data;
          return fetchSyncStatus();
        });
      }
      return data;
    })
    .then(function(data) {
      if (!data) return null;
      syncConfigured = !!data.configured;
      updateSyncBannerFromState(data);
      updateSyncStatusLines(data);
      return data;
    })
    .catch(function() {
      reportSyncOffline();
      return null;
    });
}

export function updateSyncBannerFromState(status) {
  var el = document.getElementById('sync-banner');
  if (!el) return;
  if (syncOffline) {
    el.innerHTML = t('sync.bannerOffline');
    el.style.display = 'block';
  } else if (state.syncOn && status && !status.configured) {
    el.innerHTML = t('sync.bannerNeedToken');
    el.style.display = 'block';
  } else if (status && status.lastError) {
    el.innerHTML = tf('sync.bannerErrorRenew', status.lastError);
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
    var missing = ghToken()
      ? t('sync.statusRestoring')
      : t('sync.statusTokenMissing');
    for (var i = 0; i < els.length; i++) {
      els[i].textContent = missing;
      els[i].style.color = '#b45309';
    }
    return;
  }
  if (data.lastError) {
    for (var e = 0; e < els.length; e++) {
      els[e].textContent = tf('sync.statusFailed', data.lastError);
      els[e].style.color = '#e57373';
    }
    return;
  }
  if (data.syncing) {
    for (var s = 0; s < els.length; s++) {
      els[s].textContent = t('sync.syncing');
      els[s].style.color = '#888';
    }
    return;
  }
  var text = data.lastSyncAt
    ? tf('sync.lastSynced', formatSyncAgo(data.lastSyncAt))
    : t('sync.never');
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
  var idleTicks = 0;
  var baseline = opts.baselineLastSync || 0;
  return new Promise(function(resolve) {
    function tick() {
      refreshSyncStatus().then(function(s) {
        if (!s) { resolve(null); return; }
        if (s.syncing) {
          sawSyncing = true;
          idleTicks = 0;
        } else {
          idleTicks++;
        }
        if (s.lastError) { resolve(s); return; }
        if (!s.syncing && (sawSyncing || s.lastSyncAt > baseline || idleTicks >= 2)) {
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
  el.innerHTML = '<strong>' + t('sync.clashTitle') + '</strong> ' +
    tf('sync.clashBody', noteName, copyName);
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
    .then(function(body) {
      state.editorDiskHash = body !== null ? strHash(body) : '';
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
  el.innerHTML = '<strong>' + t('sync.driftTitle') + '</strong> ' +
    tf('sync.driftBody', label) +
    '<button type="button" id="drift-reload-btn">' + t('sync.driftReload') + '</button>' +
    t('sync.driftOrKeep');
  el.style.display = 'block';
  var btn = document.getElementById('drift-reload-btn');
  if (btn) {
    btn.onclick = function(e) {
      e.stopPropagation();
      fetch('/api/reload', { method: 'POST', credentials: 'same-origin' })
        .then(function(r) {
          if (!r.ok) {
            alert(t('sync.reloadFailed'));
            return;
          }
          hideDriftBanner();
          return recordEditorDiskBaseline(filename);
        })
        .catch(function() { alert(t('generic.reachServer')); });
    };
  }
  _onBannerChange();
}

export function checkDiskDrift(filename) {
  if (!filename || !state.editorDiskHash) return Promise.resolve();
  return fetch('/api/notes/' + encodeURIComponent(filename), { credentials: 'same-origin' })
    .then(function(r) { return r.ok ? r.text() : null; })
    .then(function(body) {
      if (body === null) return;
      if (strHash(body) !== state.editorDiskHash) showDriftBanner(filename);
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
