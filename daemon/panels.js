// panels.js — PIN screen, Notes sync setup overlays.
import { state } from './state.js';
import {
  updateSyncBannerFromState, refreshSyncStatus, syncConfigured, waitForSyncIdle,
  ghToken, setSyncToken, clearSyncToken, pullTokenFromTablet, pushStoredTokenToTablet
} from './sync.js';
import { deps } from './deps.js';
import { connect, grab, setStatus, startStatusPoll, stopStatusPoll } from './connection.js';
import { loadNotes, showIdleKeyboardView } from './notes-ui.js';
import { t, tf, formatSyncAgo } from './i18n.js';

// state.syncOn and state.syncRepo (in state.js) mirror /api/settings

// Auto-advance poll: while the PIN screen is up, poll GET /api/notes every
// ~3 s. On 200 (owner switched to no-PIN, or PIN accepted from another
// client) auto-advance to the keyboard shell.
var pinPollTimer = null;
function startPinPoll() {
  if (pinPollTimer) return;
  pinPollTimer = setInterval(function() {
    fetch('/api/notes')
      .then(function(r) {
        if (r.ok) {
          stopPinPoll();
          hidePinScreen();
          connect();
          loadSyncConfig();
          showIdleKeyboardView();
          return r.json();
        }
      })
      .catch(function() {});
  }, 3000);
}
function stopPinPoll() {
  if (pinPollTimer) { clearInterval(pinPollTimer); pinPollTimer = null; }
}

function configurePinInput(pinDigits) {
  var pinInput = document.getElementById('pin-input');
  if (!pinInput) return;
  if (pinDigits === '4') {
    pinInput.maxLength = 4;
    pinInput.placeholder = '0000';
  } else {
    pinInput.maxLength = 6;
    pinInput.placeholder = '000000';
  }
}

export function updateBannerOffset() {
  var stack = document.getElementById('banner-stack');
  var h = 40;
  if (stack) {
    var visible = 0;
    ['sync-banner', 'clash-banner', 'drift-banner', 'remote-keys-banner', 'observe-banner'].forEach(function(id) {
      var el = document.getElementById(id);
      if (el && el.style.display !== 'none') visible += el.offsetHeight;
    });
    h += visible;
  }
  document.documentElement.style.setProperty('--below-bar', h + 'px');
}

export function showSync() {
  document.getElementById('sync-screen').style.display = 'flex';
  loadSyncPanel();
}

export function hideSync() {
  document.getElementById('sync-screen').style.display = 'none';
  grab();
}

export function wireOverlayDismiss(screenId, boxId, hideFn) {
  var screen = document.getElementById(screenId);
  var box = document.getElementById(boxId);
  screen.addEventListener('click', function(e) {
    if (e.target === screen) hideFn();
  });
  box.addEventListener('click', function(e) { e.stopPropagation(); });
}

function loadSyncPanel() {
  fetch('/api/settings')
    .then(function(r) { return r.json(); })
    .then(function(data) {
      state.syncOn = !!data.syncOn;
      state.syncRepo = data.syncRepo || '';
      renderSyncPanel();
      updateSyncBannerFromState();
    })
    .catch(function() {});
}

function updateSyncButtonStyles(saveBtn, syncBtn) {
  var hasToken = !!ghToken();
  if (hasToken) {
    syncBtn.className = 'sync-btn';
    saveBtn.className = 'sync-btn-secondary';
  } else {
    saveBtn.className = 'sync-btn';
    syncBtn.className = 'sync-btn-secondary';
  }
}

function renderSyncPanel() {
  var list = document.getElementById('sync-list');
  list.innerHTML = '';
  renderSyncControls(list);
}

function renderSyncControls(list) {
  var syncToggle = document.createElement('div');
  syncToggle.className = 'font-row' + (state.syncOn ? ' active' : '');
  syncToggle.innerHTML = '<span>' + t('sync.toggle') + '</span>' +
    (state.syncOn ? '<span class="font-check">&#10003;</span>' : '');
  syncToggle.addEventListener('click', function(e) {
    e.stopPropagation();
    var newVal = !state.syncOn;
    fetch('/api/settings', {
      method: 'POST', headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({ syncOn: newVal })
    }).then(function(r) {
      if (r.ok) {
        state.syncOn = newVal;
        renderSyncPanel();
        refreshSyncStatus().then(function(s) { updateSyncBannerFromState(s); });
      }
    }).catch(function(){});
  });
  list.appendChild(syncToggle);

  if (state.syncOn) {
    var repoLabel = document.createElement('div');
    repoLabel.style.cssText = 'color:#888;font-size:12px;margin-top:6px;padding:0 2px;';
    repoLabel.textContent = t('sync.repoLabel');
    list.appendChild(repoLabel);

    var repoInput = document.createElement('input');
    repoInput.type = 'text'; repoInput.className = 'token-input';
    repoInput.style.width = '100%';
    repoInput.placeholder = t('sync.repoPlaceholder'); repoInput.value = state.syncRepo;
    list.appendChild(repoInput);

    if (state.syncRepo) {
      var repoLink = document.createElement('div');
      repoLink.className = 'repo-link';
      repoLink.innerHTML = '<a href="https://github.com/' + state.syncRepo +
        '" target="_blank" rel="noopener noreferrer">github.com/' + state.syncRepo + '</a>';
      list.appendChild(repoLink);
    }

    var tokLabel = document.createElement('div');
    tokLabel.style.cssText = 'color:#888;font-size:12px;margin-top:8px;padding:0 2px;';
    tokLabel.textContent = t('sync.tokenLabel');
    list.appendChild(tokLabel);

    var tokInput = document.createElement('input');
    tokInput.type = 'password'; tokInput.className = 'token-input';
    tokInput.style.width = '100%';
    tokInput.placeholder = t('sync.tokenPlaceholder'); tokInput.autocomplete = 'off';
    if (ghToken() || syncConfigured) tokInput.value = '\u2022'.repeat(16);
    var tokTouched = false;
    tokInput.addEventListener('focus', function() {
      if (!tokTouched) { tokInput.value = ''; tokTouched = true; }
    });
    list.appendChild(tokInput);

    var verifyLine = document.createElement('div');
    verifyLine.className = 'sync-verify-line';
    verifyLine.style.cssText = 'font-size:12px;padding:6px 2px;min-height:16px;color:#888;';

    var actionRow = document.createElement('div');
    actionRow.style.cssText = 'display:flex;gap:8px;align-items:center;margin-top:8px;';

    var saveBtn = document.createElement('button');
    saveBtn.textContent = t('sync.save');

    var syncBtn = document.createElement('button');
    syncBtn.textContent = t('sync.sync');
    updateSyncButtonStyles(saveBtn, syncBtn);

    saveBtn.addEventListener('click', function(e) {
      e.stopPropagation();
      var repoVal = repoInput.value.trim();
      fetch('/api/settings', {
        method: 'POST', headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ syncRepo: repoVal })
      }).then(function(r) {
        if (!r.ok) {
          verifyLine.style.color = '#e57373';
          verifyLine.textContent = t('sync.invalidRepo');
          return null;
        }
        state.syncRepo = repoVal;
        var link = list.querySelector('.repo-link');
        if (repoVal) {
          if (!link) {
            link = document.createElement('div');
            link.className = 'repo-link';
            repoInput.parentNode.insertBefore(link, repoInput.nextSibling);
          }
          link.innerHTML = '<a href="https://github.com/' + repoVal +
            '" target="_blank" rel="noopener noreferrer">github.com/' + repoVal + '</a>';
        } else if (link) { link.remove(); }
        if (tokTouched && tokInput.value.trim()) {
          setSyncToken(tokInput.value.trim());
          tokInput.value = '\u2022'.repeat(16); tokTouched = false;
          updateSyncButtonStyles(saveBtn, syncBtn);
        }
        var token = ghToken();
        if (!token) {
          if (!repoVal) {
            verifyLine.style.color = '#888';
            verifyLine.textContent = t('sync.savedRepo');
            return refreshSyncStatus();
          }
          verifyLine.style.color = '#888';
          verifyLine.textContent = t('sync.savedEnterToken');
          return refreshSyncStatus();
        }
        if (!repoVal) {
          verifyLine.style.color = '#888';
          verifyLine.textContent = t('sync.savedEnterRepo');
          return refreshSyncStatus();
        }
        return fetch('/api/sync/token', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          credentials: 'same-origin',
          body: JSON.stringify({ token: token })
        }).then(function(tr) {
          tokInput.value = '\u2022'.repeat(16); tokTouched = false;
          updateSyncButtonStyles(saveBtn, syncBtn);
          if (tr.status === 401) {
            clearSyncToken();
            tokInput.value = ''; tokTouched = true;
            updateSyncButtonStyles(saveBtn, syncBtn);
            verifyLine.style.color = '#e57373';
            verifyLine.textContent = t('sync.tokenRejected');
            return refreshSyncStatus();
          }
          if (tr.status === 404) {
            verifyLine.style.color = '#e57373';
            verifyLine.textContent = t('sync.repoNotFound');
            return refreshSyncStatus();
          }
          if (!tr.ok) {
            verifyLine.style.color = '#e57373';
            verifyLine.textContent = t('sync.verifyFailed');
            return refreshSyncStatus();
          }
          verifyLine.style.color = '#4caf50';
          verifyLine.textContent = t('sync.tokenSaved');
          return refreshSyncStatus();
        });
      }).catch(function() {
        verifyLine.style.color = '#e57373';
        verifyLine.textContent = t('sync.reachTablet');
      });
    });

    syncBtn.addEventListener('click', function(e) {
      e.stopPropagation();
      verifyLine.style.color = '#888';
      verifyLine.textContent = t('sync.syncing');
      var runSync = function() {
        return fetch('/api/sync/run', {
          method: 'POST', credentials: 'same-origin'
        }).then(function(r) {
          if (r.status === 400) {
            verifyLine.style.color = '#e57373';
            verifyLine.textContent = t('sync.notConfigured');
            return null;
          }
          if (!r.ok) {
            verifyLine.style.color = '#e57373';
            verifyLine.textContent = t('sync.failed');
            return null;
          }
          return refreshSyncStatus().then(function(before) {
            var baseline = (before && before.lastSyncAt) || 0;
            return waitForSyncIdle({ baselineLastSync: baseline });
          });
        });
      };
      var prep = ghToken()
        ? refreshSyncStatus().then(function(s) {
            if (s && s.configured) return;
            return pushStoredTokenToTablet();
          })
        : pullTokenFromTablet().then(function(pulled) {
            if (pulled) {
              tokInput.value = '\u2022'.repeat(16); tokTouched = false;
              updateSyncButtonStyles(saveBtn, syncBtn);
              return refreshSyncStatus();
            }
            if (!ghToken()) {
              verifyLine.style.color = '#e57373';
              verifyLine.textContent = t('sync.noToken');
              return null;
            }
          });
      prep.then(function() {
        if (verifyLine.textContent.indexOf('\u2717') === 0) return;
        return runSync();
      }).then(function(s) {
        if (!s) {
          if (verifyLine.textContent === t('sync.syncing')) {
            verifyLine.style.color = '#e57373';
            verifyLine.textContent = t('sync.reachTablet');
          }
          return;
        }
        loadNotes();
        if (s.lastError) {
          verifyLine.style.color = '#e57373';
          verifyLine.textContent = '\u2717 ' + s.lastError;
          return;
        }
        verifyLine.style.color = '#4caf50';
        var when = formatSyncAgo(s.lastSyncAt) || t('sync.agoJustNow');
        return fetch('/api/status', { credentials: 'same-origin' })
          .then(function(r) { return r.ok ? r.json() : null; })
          .then(function(st) {
            if (st && st.openNote) {
              verifyLine.textContent = tf('sync.syncedSkippedOpen',
                st.openNote.replace(/\.md$/, ''));
              return;
            }
            verifyLine.textContent = tf('sync.synced', when);
          });
      }).catch(function() {
        verifyLine.style.color = '#e57373';
        verifyLine.textContent = t('sync.reachTablet');
      });
    });

    var clearBtn = document.createElement('button');
    clearBtn.className = 'sync-btn-secondary'; clearBtn.textContent = t('sync.clearToken');
    clearBtn.addEventListener('click', function(e) {
      e.stopPropagation();
      if (!confirm(t('sync.clearTokenConfirm'))) return;
      clearSyncToken();
      updateSyncButtonStyles(saveBtn, syncBtn);
      fetch('/api/sync/token', {
        method: 'POST', headers: {'Content-Type': 'application/json'},
        credentials: 'same-origin',
        body: JSON.stringify({ token: '' })
      }).then(function() {
        tokInput.value = ''; tokTouched = true;
        verifyLine.style.color = '#888';
        verifyLine.textContent = t('sync.tokenCleared');
        return refreshSyncStatus();
      }).then(function(s) { updateSyncBannerFromState(s); });
    });

    actionRow.appendChild(saveBtn); actionRow.appendChild(syncBtn); actionRow.appendChild(clearBtn);
    list.appendChild(actionRow);
    list.appendChild(verifyLine);

    var hintLine = document.createElement('div');
    hintLine.style.cssText = 'font-size:11px;color:#888;padding:4px 2px;';
    hintLine.textContent = t('sync.hint');
    list.appendChild(hintLine);

    var statusLine = document.createElement('div');
    statusLine.className = 'sync-status-line';
    statusLine.style.cssText = 'font-size:12px;color:#888;padding:4px 2px;';
    statusLine.textContent = t('sync.loadingStatus');
    list.appendChild(statusLine);
    var hadLocalToken = !!ghToken();
    refreshSyncStatus().then(function() {
      updateSyncButtonStyles(saveBtn, syncBtn);
      if (!hadLocalToken && ghToken()) {
        verifyLine.style.color = '#4caf50';
        verifyLine.textContent = t('sync.tokenRestored');
        tokInput.value = '\u2022'.repeat(16);
        tokTouched = false;
      } else if (ghToken() || syncConfigured) {
        tokInput.value = '\u2022'.repeat(16);
        tokTouched = false;
      }
    });
  }
}

export function showPinScreen() {
  stopStatusPoll();
  document.getElementById('pin-screen').style.display = 'flex';
  startPinPoll();
  var pinInput = document.getElementById('pin-input');
  if (pinInput) {
    pinInput.focus();
    // Soft keyboards / USB on phone: ensure the field is ready without a tap.
    try { pinInput.select(); } catch (e) {}
  }
}

export function hidePinScreen() {
  stopPinPoll();
  document.getElementById('pin-screen').style.display = 'none';
  document.getElementById('sync-btn').style.display = 'block';
  startStatusPoll();
}

export function loadSyncConfig() {
  return fetch('/api/settings')
    .then(function(r) { return r.ok ? r.json() : null; })
    .then(function(data) {
      if (!data) return;
      state.syncOn = !!data.syncOn;
      state.syncRepo = data.syncRepo || '';
      if (deps.applyObserveEnabled) deps.applyObserveEnabled(!!data.observe);
      if (data.observe && deps.refreshObserveStatus) deps.refreshObserveStatus();
      return refreshSyncStatus().then(function(s) { updateSyncBannerFromState(s); });
    })
    .catch(function() {});
}

export function checkAuthAndInit() {
  fetch('/api/settings')
    .then(function(r) { return r.ok ? r.json() : null; })
    .then(function(settings) {
      if (settings) configurePinInput(settings.pinDigits || '6');
      if (settings && deps.applyObserveEnabled) {
        deps.applyObserveEnabled(!!settings.observe);
      }
      if (settings && settings.pinDigits === 'none') {
        document.getElementById('pin-screen').style.display = 'none';
      }
      return fetch('/api/notes');
    })
    .then(function(r) {
      if (r.status === 401) { showPinScreen(); return null; }
      if (!r.ok) throw new Error('server error');
      return r.json();
    })
    .then(function(notes) {
      if (notes === null) return;
      hidePinScreen();
      connect();
      showIdleKeyboardView();
      loadSyncConfig();
    })
    .catch(function() {
      setStatus('off', t('conn.retrying'));
      setTimeout(checkAuthAndInit, 2000);
    });
}

export function submitPIN(e) {
  if (e) e.stopPropagation();
  var pin = document.getElementById('pin-input').value.trim();
  var errEl = document.getElementById('pin-err');
  errEl.textContent = '';
  fetch('/api/pin', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ pin: pin })
  }).then(function(r) {
    if (r.status === 429) {
      errEl.textContent = t('pin.tooMany');
      document.getElementById('pin-input').value = '';
      return;
    }
    if (r.status === 401) {
      errEl.textContent = t('pin.wrong');
      document.getElementById('pin-input').value = '';
      document.getElementById('pin-input').focus();
      return;
    }
    if (!r.ok) { errEl.textContent = t('pin.serverError'); return; }
    hidePinScreen();
    connect();
    showIdleKeyboardView();
    loadNotes();
    loadSyncConfig();
  }).catch(function() { errEl.textContent = t('generic.reachServer'); });
}
