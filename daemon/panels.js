// panels.js — PIN screen, Preferences, GitHub Setup overlays.
import { state } from './state.js';
import {
  updateSyncBannerFromState, refreshSyncStatus, syncConfigured, waitForSyncIdle,
  ghToken, setSyncToken, clearSyncToken
} from './sync.js';
import { deps } from './deps.js';
import { connect, grab, setStatus, startStatusPoll, stopStatusPoll, closeWebSocket } from './connection.js';
import { renderNotes, loadNotes } from './notes-ui.js';

// ---- Auth ----

// --- Settings panel (font + PIN picker) ---
var currentFont = '';
var settingsFonts = [];
var currentPinDigits = '6';
// state.syncOn and state.syncRepo (in state.js) mirror /api/settings

// Auto-advance poll: while the PIN screen is up, poll GET /api/notes every
// ~3 s. On 200 (owner switched to no-PIN, or PIN accepted from another
// client) auto-advance to the notes view.
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
          return r.json();
        }
      })
      .then(function(notes) { if (notes) renderNotes(notes); })
      .catch(function() {});
  }, 3000);
}
function stopPinPoll() {
  if (pinPollTimer) { clearInterval(pinPollTimer); pinPollTimer = null; }
}

function applyPinDigits(d) {
  currentPinDigits = d;
  renderSettingsList(); // immediate UI feedback
  fetch('/api/settings', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({pinDigits: d})
  }).catch(function() {});
}

export function updateBannerOffset() {
  var stack = document.getElementById('banner-stack');
  var h = 40;
  if (stack) {
    var visible = 0;
    ['sync-banner', 'clash-banner', 'drift-banner'].forEach(function(id) {
      var el = document.getElementById(id);
      if (el && el.style.display !== 'none') visible += el.offsetHeight;
    });
    h += visible;
  }
  document.documentElement.style.setProperty('--below-bar', h + 'px');
}

export function showSettings() {
  document.getElementById('settings-screen').style.display = 'flex';
  loadSettingsPanel();
}

export function hideSettings() {
  document.getElementById('settings-screen').style.display = 'none';
  grab();
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

function loadSettingsPanel() {
  fetch('/api/settings')
    .then(function(r) { return r.json(); })
    .then(function(data) {
      currentFont = data.readFont;
      settingsFonts = data.fonts;
      currentPinDigits = data.pinDigits || '6';
      state.syncOn = !!data.syncOn;
      state.syncRepo = data.syncRepo || '';
      renderSettingsList();
    })
    .catch(function() {});
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

function renderSyncPanel() {
  var list = document.getElementById('sync-list');
  list.innerHTML = '';
  renderSyncControls(list);
}

function requestRotate(onMsg) {
  return fetch('/api/rotate', { method: 'POST' })
    .then(function(r) {
      if (r.ok) {
        if (onMsg) { onMsg.style.color = '#4caf50'; onMsg.textContent = 'Rotated 90\u00b0.'; }
        return;
      }
      var err = 'Rotate failed (' + r.status + ').';
      if (r.status === 409) err = 'No editor session \u2014 open a note on the tablet.';
      else if (r.status === 401) err = 'Not authorized \u2014 reconnect and enter PIN.';
      if (onMsg) { onMsg.style.color = '#e57373'; onMsg.textContent = err; }
      else { alert(err); }
    })
    .catch(function() {
      var err = 'Could not reach the tablet.';
      if (onMsg) { onMsg.style.color = '#e57373'; onMsg.textContent = err; }
      else { alert(err); }
    });
}

function renderSettingsList() {
  var list = document.getElementById('settings-list');
  list.innerHTML = '';

  // ---- Font section ----
  var fh = document.createElement('div');
  fh.className = 'settings-section';
  fh.textContent = 'Reading font';
  list.appendChild(fh);
  settingsFonts.forEach(function(f) {
    var row = document.createElement('div');
    row.className = 'font-row' + (f.id === currentFont ? ' active' : '');
    var check = f.id === currentFont ? '<span class="font-check">&#10003;</span>' : '';
    row.innerHTML = '<span>' + f.label + '</span>' + check;
    row.addEventListener('click', function(e) {
      e.stopPropagation();
      applyFont(f.id);
    });
    list.appendChild(row);
  });

  // ---- PIN section ----
  var ph = document.createElement('div');
  ph.className = 'settings-section';
  ph.textContent = 'PIN length';
  list.appendChild(ph);
  var pinOpts = [
    {id: '6', label: '6 digits'},
    {id: '4', label: '4 digits'},
    {id: 'none', label: 'No PIN', warn: 'Anyone on your Wi-Fi can read & edit your notes'}
  ];
  // Sync pin-input maxlength/placeholder with the active PIN length.
  var pinInput = document.getElementById('pin-input');
  if (pinInput) {
    if (currentPinDigits === '4') {
      pinInput.maxLength = 4; pinInput.placeholder = '0000';
    } else {
      pinInput.maxLength = 6; pinInput.placeholder = '000000';
    }
  }
  pinOpts.forEach(function(opt) {
    var row = document.createElement('div');
    row.className = 'font-row' + (opt.id === currentPinDigits ? ' active' : '');
    var inner = '<div><div>' + opt.label + '</div>';
    if (opt.warn) {
      inner += '<div class="row-warn">' + opt.warn + '</div>';
    }
    inner += '</div>';
    if (opt.id === currentPinDigits) { inner += '<span class="font-check">&#10003;</span>'; }
    row.innerHTML = inner;
    row.addEventListener('click', function(e) {
      e.stopPropagation();
      applyPinDigits(opt.id);
    });
    list.appendChild(row);
  });

  // ---- Display section ----
  var dh = document.createElement('div');
  dh.className = 'settings-section';
  dh.textContent = 'Display';
  list.appendChild(dh);

  var rotateMsg = document.createElement('div');
  rotateMsg.style.cssText = 'font-size:12px;padding:4px 2px;min-height:16px;color:#888;';

  var rotateBtn = document.createElement('button');
  rotateBtn.className = 'sync-btn';
  rotateBtn.style.width = '100%';
  rotateBtn.textContent = 'Rotate tablet 90\u00b0';
  rotateBtn.addEventListener('click', function(e) {
    e.stopPropagation();
    rotateMsg.style.color = '#888';
    rotateMsg.textContent = 'Rotating\u2026';
    requestRotate(rotateMsg);
  });
  list.appendChild(rotateBtn);
  list.appendChild(rotateMsg);

  // ---- Service section ----
  var svh = document.createElement('div');
  svh.className = 'settings-section';
  svh.textContent = 'Service';
  list.appendChild(svh);

  var exitWarn = document.createElement('div');
  exitWarn.className = 'row-warn';
  exitWarn.style.cssText = 'padding:4px 2px 8px;line-height:1.4;';
  exitWarn.textContent = 'Stop Writerdeck and return the tablet to the stock reMarkable UI. Reconnect later via SSH or reboot.';
  list.appendChild(exitWarn);

  var exitMsg = document.createElement('div');
  exitMsg.style.cssText = 'font-size:12px;padding:4px 2px;min-height:16px;color:#888;';

  var exitBtn = document.createElement('button');
  exitBtn.className = 'sync-btn-danger';
  exitBtn.textContent = 'Exit Writerdeck';
  exitBtn.addEventListener('click', function(e) {
    e.stopPropagation();
    if (!confirm('Exit Writerdeck on the tablet? This page will disconnect.')) return;
    exitMsg.style.color = '#888';
    exitMsg.textContent = 'Stopping\u2026';
    fetch('/api/shutdown', { method: 'POST' })
      .then(function(r) {
        if (!r.ok) throw new Error('status ' + r.status);
        exitMsg.style.color = '#4caf50';
        exitMsg.textContent = 'Stopped. Stock UI should be back on the tablet.';
          closeWebSocket();
        setStatus('off', 'Writerdeck stopped on tablet');
        stopStatusPoll();
      })
      .catch(function() {
        exitMsg.style.color = '#e57373';
        exitMsg.textContent = 'Could not stop -- try again or use SSH.';
      });
  });
  list.appendChild(exitBtn);
  list.appendChild(exitMsg);
}

function renderSyncControls(list) {
  var syncToggle = document.createElement('div');
  syncToggle.className = 'font-row' + (state.syncOn ? ' active' : '');
  syncToggle.innerHTML = '<span>Sync notes to GitHub</span>' +
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
    repoLabel.textContent = 'Private repo (owner/repo)';
    list.appendChild(repoLabel);

    var repoInput = document.createElement('input');
    repoInput.type = 'text'; repoInput.className = 'token-input';
    repoInput.style.width = '100%';
    repoInput.placeholder = 'e.g. alice/my-notes'; repoInput.value = state.syncRepo;
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
    tokLabel.textContent = 'GitHub token (saved in this browser; copied to tablet RAM when verified)';
    list.appendChild(tokLabel);

    var tokInput = document.createElement('input');
    tokInput.type = 'password'; tokInput.className = 'token-input';
    tokInput.style.width = '100%';
    tokInput.placeholder = 'github_pat_\u2026 or ghp_\u2026'; tokInput.autocomplete = 'off';
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
    saveBtn.className = 'sync-btn'; saveBtn.textContent = 'Save & verify';
    saveBtn.addEventListener('click', function(e) {
      e.stopPropagation();
      var repoVal = repoInput.value.trim();
      fetch('/api/settings', {
        method: 'POST', headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ syncRepo: repoVal })
      }).then(function(r) {
        if (!r.ok) {
          verifyLine.style.color = '#e57373';
          verifyLine.textContent = '\u2717 Invalid repo \u2014 use owner/repo format.';
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
        }
        var token = ghToken();
        if (!token) {
          if (!repoVal) {
            verifyLine.style.color = '#888';
            verifyLine.textContent = 'Saved repo.';
            return refreshSyncStatus();
          }
          verifyLine.style.color = '#888';
          verifyLine.textContent = 'Saved. Enter token to verify.';
          return refreshSyncStatus();
        }
        if (!repoVal) {
          verifyLine.style.color = '#888';
          verifyLine.textContent = 'Saved token. Enter repo to verify.';
          return refreshSyncStatus();
        }
        return fetch('/api/sync/token', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          credentials: 'same-origin',
          body: JSON.stringify({ token: token })
        }).then(function(tr) {
          tokInput.value = '\u2022'.repeat(16); tokTouched = false;
          if (tr.status === 401) {
            clearSyncToken();
            tokInput.value = ''; tokTouched = true;
            verifyLine.style.color = '#e57373';
            verifyLine.textContent = '\u2717 Token rejected.';
            return refreshSyncStatus();
          }
          if (tr.status === 404) {
            verifyLine.style.color = '#e57373';
            verifyLine.textContent = '\u2717 Repo not found.';
            return refreshSyncStatus();
          }
          if (!tr.ok) {
            verifyLine.style.color = '#e57373';
            verifyLine.textContent = '\u2717 Could not verify.';
            return refreshSyncStatus();
          }
          verifyLine.style.color = '#888';
          verifyLine.textContent = 'Syncing on tablet\u2026';
          return refreshSyncStatus().then(function(before) {
            var baseline = (before && before.lastSyncAt) || 0;
            return waitForSyncIdle({ baselineLastSync: baseline });
          }).then(function(s) {
            loadNotes();
            if (s && s.lastError) {
              verifyLine.style.color = '#e57373';
              verifyLine.textContent = '\u2717 ' + s.lastError;
              return;
            }
            verifyLine.style.color = '#4caf50';
            var when = (s && s.lastSyncAgo) ? s.lastSyncAgo : 'just now';
            return fetch('/api/status', { credentials: 'same-origin' })
              .then(function(r) { return r.ok ? r.json() : null; })
              .then(function(st) {
                if (st && st.openNote) {
                  verifyLine.textContent = '\u2713 Token saved \u2014 synced other notes; \u201c' +
                    st.openNote.replace(/\.md$/, '') + '\u201d skipped while open.';
                  return;
                }
                verifyLine.textContent = '\u2713 Connected \u2014 last synced ' + when + '.';
              });
          });
        });
      }).catch(function() {
        verifyLine.style.color = '#e57373';
        verifyLine.textContent = '\u2717 Could not reach the tablet.';
      });
    });

    var clearBtn = document.createElement('button');
    clearBtn.className = 'sync-btn-secondary'; clearBtn.textContent = 'Clear token';
    clearBtn.addEventListener('click', function(e) {
      e.stopPropagation();
      if (!confirm('Remove GitHub token from this browser and the tablet?')) return;
      clearSyncToken();
      fetch('/api/sync/token', {
        method: 'POST', headers: {'Content-Type': 'application/json'},
        credentials: 'same-origin',
        body: JSON.stringify({ token: '' })
      }).then(function() {
        tokInput.value = ''; tokTouched = true;
        verifyLine.style.color = '#888';
        verifyLine.textContent = 'Token cleared.';
        return refreshSyncStatus();
      }).then(function(s) { updateSyncBannerFromState(s); });
    });

    actionRow.appendChild(saveBtn); actionRow.appendChild(clearBtn);
    list.appendChild(actionRow);
    list.appendChild(verifyLine);

    var hintLine = document.createElement('div');
    hintLine.style.cssText = 'font-size:11px;color:#888;padding:4px 2px;';
    hintLine.textContent = 'Run sync from the tablet Lobby \u2192 Sync tab.';
    list.appendChild(hintLine);

    var statusLine = document.createElement('div');
    statusLine.className = 'sync-status-line';
    statusLine.style.cssText = 'font-size:12px;color:#888;padding:4px 2px;';
    statusLine.textContent = 'Loading sync status\u2026';
    list.appendChild(statusLine);
    refreshSyncStatus();
  }
}

function applyFont(id) {
  currentFont = id;
  renderSettingsList(); // immediate UI feedback
  fetch('/api/settings', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({readFont: id})
  }).catch(function() {}); // Go pushes setfont to the editor automatically
}

export function showPinScreen() {
  stopStatusPoll();
  document.getElementById('pin-screen').style.display = 'flex';
  startPinPoll();
}

export function hidePinScreen() {
  stopPinPoll();
  document.getElementById('pin-screen').style.display = 'none';
  document.getElementById('settings-btn').style.display = 'block';
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
      return refreshSyncStatus().then(function(s) { updateSyncBannerFromState(s); });
    })
    .catch(function() {});
}

export function checkAuthAndInit() {
  fetch('/api/notes')
    .then(function(r) {
      if (r.status === 401) { showPinScreen(); return null; }
      if (!r.ok) throw new Error('server error');
      return r.json();
    })
    .then(function(notes) {
      if (notes === null) return;
      hidePinScreen();
      connect();
      grab();
      renderNotes(notes);
      loadSyncConfig(); // populate syncOn/syncRepo so auto-sync can fire
    })
    .catch(function() { setTimeout(checkAuthAndInit, 2000); });
}

export function submitPIN(e) {
  if (e) e.stopPropagation();
  var pin = document.getElementById('pin-input').value.trim();
  var errEl = document.getElementById('pin-err');
  errEl.textContent = '';
  if (!pin) return;
  fetch('/api/pin', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ pin: pin })
  }).then(function(r) {
    if (r.status === 429) {
      errEl.textContent = 'Too many attempts. Wait a minute, then try again.';
      document.getElementById('pin-input').value = '';
      return;
    }
    if (r.status === 401) {
      errEl.textContent = 'Wrong PIN, try again.';
      document.getElementById('pin-input').value = '';
      document.getElementById('pin-input').focus();
      return;
    }
    if (!r.ok) { errEl.textContent = 'Server error.'; return; }
    hidePinScreen();
    connect();
    grab();
    loadNotes();
    loadSyncConfig();
  }).catch(function() { errEl.textContent = 'Could not reach server.'; });
}
