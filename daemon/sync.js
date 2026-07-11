// sync.js -- GitHub two-way sync engine.
// Imported by app.js; never imports app.js (DAG: state.js <- sync.js <- app.js).
// Shared state (syncOn, syncRepo, tabletOpenNote, typingMode) lives in state.js
// so both modules can read/write the same values without circular imports.
import { state } from './state.js';

var SYNC_POLL_MS = 180000; // 3 min: safety-net reconcile for laptop-side edits

// initSync: wire back-references into app.js. Call once from the load
// handler: initSync({ onNotesChanged: loadNotes, onBannerChange: updateBannerOffset }).
var _onNotesChanged = function() {};
var _onBannerChange = function() {};

export function initSync(opts) {
  _onNotesChanged = opts.onNotesChanged || function() {};
  _onBannerChange = opts.onBannerChange || function() {};
}

// setSyncToken / clearSyncStorage: encapsulate the gh* localStorage key
// schema here so only sync.js knows the key names.
export function setSyncToken(token) {
  localStorage.setItem('ghToken', token);
}
export function clearSyncStorage() {
  localStorage.removeItem('ghToken');
  var keys = [];
  for (var i = 0; i < localStorage.length; i++) {
    var k = localStorage.key(i);
    if (k && (k.startsWith('ghSha_') || k.startsWith('ghPushFailed_') || k.startsWith('ghLocalHash_'))) keys.push(k);
  }
  keys.forEach(function(k) { localStorage.removeItem(k); });
}

export function ghToken() { return localStorage.getItem('ghToken') || ''; }
export function syncReady() { return state.syncOn && !!ghToken() && !!state.syncRepo; }

// UTF-8-safe base64 encode/decode (btoa/atob are ASCII-only without this wrapper).
function b64encode(str) { return btoa(unescape(encodeURIComponent(str))); }
function b64decode(str) { return decodeURIComponent(escape(atob(str.replace(/\s/g, '')))); }

function ghUrl(filename) {
  return 'https://api.github.com/repos/' + state.syncRepo + '/contents/' + encodeURIComponent(filename);
}
function ghHdrs() {
  return {
    'Authorization': 'Bearer ' + ghToken(),
    'Accept': 'application/vnd.github.v3+json',
    'Content-Type': 'application/json'
  };
}

export function updateSyncBannerFromState() {
  var el = document.getElementById('sync-banner');
  if (!el) return;
  if (state.syncOn && !ghToken()) {
    el.innerHTML = '\u26a0 GitHub sync is on \u2014 add your token in <strong>Sync</strong>.';
    el.style.display = 'block';
  } else {
    el.style.display = 'none';
  }
  _onBannerChange();
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

// handleClash: on a 409/422 push clash --
//   1. save current tablet content as "{name} (tablet copy).md"
//   2. fetch GitHub version, write it to "{name}.md" on the tablet
//   3. update stored SHA; show clash banner.
function handleClash(filename, tabletContent) {
  return fetch(ghUrl(filename), { headers: ghHdrs() })
    .then(function(r) { return r.ok ? r.json() : null; })
    .then(function(ghData) {
      if (!ghData) return;
      var ghContent = b64decode(ghData.content);
      // Not a real clash if both sides already hold identical text: adopt
      // GitHub's sha + fingerprint, no duplicate, no banner.
      if (ghContent === tabletContent) {
        localStorage.setItem('ghSha_' + filename, ghData.sha);
        localStorage.setItem('ghLocalHash_' + filename, strHash(tabletContent));
        localStorage.removeItem('ghPushFailed_' + filename);
        return;
      }
      // Accidental wipe: empty tablet vs non-empty GitHub is not a real clash --
      // restore from GitHub without creating a junk "(tablet copy)" duplicate.
      if (tabletContent === '' && ghContent !== '') {
        return fetch('/api/notes/' + encodeURIComponent(filename), {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ content: ghContent })
        }).then(function(r) {
          if (r && r.ok) {
            localStorage.setItem('ghSha_' + filename, ghData.sha);
            localStorage.setItem('ghLocalHash_' + filename, strHash(ghContent));
          }
          localStorage.removeItem('ghPushFailed_' + filename);
          _onNotesChanged();
        });
      }
      var copyBase = filename.replace(/\.md$/, '') + ' (tablet copy)';
      // Keep the tablet's version as a copy, then bring GitHub's into note.md.
      return fetch('/api/notes', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: copyBase, content: tabletContent })
      }).catch(function() {}).then(function() {
        return fetch('/api/notes/' + encodeURIComponent(filename), {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ content: ghContent })
        });
      }).then(function(r) {
        if (r && r.ok) {
          localStorage.setItem('ghSha_' + filename, ghData.sha);
          localStorage.setItem('ghLocalHash_' + filename, strHash(ghContent));
        }
        localStorage.removeItem('ghPushFailed_' + filename);
        showClashBanner(filename.replace(/\.md$/, ''), copyBase);
        _onNotesChanged();
      });
    })
    .catch(function() {});
}

// pushNote: read note from tablet and push to GitHub; handle clash and auth errors.
// MUST return its promise: reconcileAll sequences pushes through a reduce chain,
// and GitHub creates one commit per PUT parented on the current branch HEAD --
// if pushes fire concurrently, only one commit wins per round and the rest 409,
// which is exactly the "one file synced per attempt" failure this return fixes.
export function pushNote(filename) {
  if (!syncReady()) return Promise.resolve();
  return fetch('/api/notes/' + encodeURIComponent(filename))
    .then(function(r) { return r.ok ? r.text() : null; })
    .then(function(content) {
      if (content === null) return;
      var storedHash = localStorage.getItem('ghLocalHash_' + filename);
      var emptyHash = strHash('');
      // Safety net: never push an empty tablet file over a previously-synced note.
      // (Lobby Home used to wipe files this way; see lessons.md.)
      if (content === '' && storedHash && storedHash !== emptyHash) {
        localStorage.setItem('ghPushFailed_' + filename, '1');
        return;
      }
      var sha = localStorage.getItem('ghSha_' + filename) || null;
      var body = { message: 'Writerdeck: ' + filename, content: b64encode(content) };
      if (sha) body.sha = sha;
      return fetch(ghUrl(filename), {
        method: 'PUT', headers: ghHdrs(), body: JSON.stringify(body)
      }).then(function(r) {
        if (r.ok) {
          return r.json().then(function(d) {
            localStorage.setItem('ghSha_' + filename, d.content.sha);
            localStorage.setItem('ghLocalHash_' + filename, strHash(content));
            localStorage.removeItem('ghPushFailed_' + filename);
            var ts = new Date().toLocaleString();
            localStorage.setItem('ghLastSync', ts);
            var els = document.querySelectorAll('.sync-status-line');
            for (var i = 0; i < els.length; i++) els[i].textContent = 'Last synced: ' + ts;
          });
        }
        if (r.status === 409 || r.status === 422) { return handleClash(filename, content); }
        if (r.status === 401 || r.status === 403) {
          localStorage.setItem('ghPushFailed_' + filename, '1');
          var banner = document.getElementById('sync-banner');
          if (banner) {
            banner.innerHTML = '\u26a0 GitHub token rejected \u2014 renew it in <strong>Sync</strong>.';
            banner.style.display = 'block';
            _onBannerChange();
          }
        } else {
          localStorage.setItem('ghPushFailed_' + filename, '1');
        }
      });
    }).catch(function() { localStorage.setItem('ghPushFailed_' + filename, '1'); });
}

// pullNoteAndUpdate: check GitHub for a newer version and write it to the tablet.
// If SHA matches stored SHA, nothing is fetched. Returns a promise.
export function pullNoteAndUpdate(filename) {
  if (!syncReady()) return Promise.resolve();
  return fetch(ghUrl(filename), { headers: ghHdrs() })
    .then(function(r) {
      if (r.status === 404) return null; // not on GitHub yet
      if (r.status === 401 || r.status === 403) {
        var banner = document.getElementById('sync-banner');
        if (banner) {
          banner.innerHTML = '\u26a0 GitHub token rejected \u2014 renew it in <strong>Sync</strong>.';
          banner.style.display = 'block';
        }
        return null;
      }
      return r.ok ? r.json() : null;
    })
    .then(function(data) {
      if (!data) return;
      if (localStorage.getItem('ghSha_' + filename) === data.sha) return; // already up to date
      var ghContent = b64decode(data.content);
      return fetch('/api/notes/' + encodeURIComponent(filename), {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content: ghContent })
      }).then(function(r) {
        if (r.ok) {
          localStorage.setItem('ghSha_' + filename, data.sha);
          localStorage.setItem('ghLocalHash_' + filename, strHash(ghContent));
        }
      });
    })
    .catch(function() {});
}

// ghDelete: remove a note from GitHub via the Contents API (needs the file's
// current sha, tracked per note). No-op if the note was never synced from here.
export function ghDelete(filename) {
  if (!syncReady()) return Promise.resolve();
  var sha = localStorage.getItem('ghSha_' + filename);
  if (!sha) return Promise.resolve();
  return fetch(ghUrl(filename), {
    method: 'DELETE', headers: ghHdrs(),
    body: JSON.stringify({ message: 'Writerdeck: delete ' + filename, sha: sha })
  }).then(function() {
    localStorage.removeItem('ghSha_' + filename);
    localStorage.removeItem('ghLocalHash_' + filename);
    localStorage.removeItem('ghPushFailed_' + filename);
  }).catch(function() {});
}

// applyRemoteDelete: a previously-synced, locally-unchanged note has vanished from
// GitHub -> treat as a real upstream delete. Confirms with a fresh per-note GET
// (guards against a stale/empty bulk list or a transient network error mapping
// failure -> [] in reconcileAll) before removing it from the tablet. A false
// positive self-heals on the next sync (re-pulled via the !hasLocal && hasRemote
// branch). Never touches the currently-open note.
function applyRemoteDelete(name) {
  if (!syncReady() || name === state.tabletOpenNote) return Promise.resolve();
  return fetch(ghUrl(name), { headers: ghHdrs() })
    .then(function(r) {
      if (r.status !== 404) return;                    // still there / uncertain -> do nothing (safe)
      return fetch('/api/notes/' + encodeURIComponent(name), { method: 'DELETE' })
        .then(function(dr) {
          if (dr && dr.ok) {
            localStorage.removeItem('ghSha_' + name);
            localStorage.removeItem('ghLocalHash_' + name);
            localStorage.removeItem('ghPushFailed_' + name);
          }
        });
    })
    .catch(function() {});                             // network error -> no delete
}

// strHash: cheap deterministic fingerprint (djb2) of a note's text, used to
// tell whether the tablet copy changed since the last sync -- the missing
// signal that lets reconcile distinguish a local-only edit from a real clash.
function strHash(s) {
  var h = 5381;
  for (var i = 0; i < s.length; i++) { h = ((h << 5) + h + s.charCodeAt(i)) | 0; }
  return String(h >>> 0);
}

// reconcileAll: full two-way sync of every note -- the trigger the event-only
// model was missing. Runs on first connect/verify, on reconnect, on a periodic
// poll, and from "Sync now". Per note it delegates to the same safe primitives
// used elsewhere: remote-only -> pull; local-only -> push; both -> compare
// fingerprints and push / pull / keep-both. The actively-open note is skipped
// (its own open-pull + Home-push own it). Resolves to the note count.
var syncing = false;
export function reconcileAll(reason, opts) {
  opts = opts || {};
  if (!syncReady()) return Promise.resolve(0);
  if (syncing) {
    if (opts.wait) {
      return new Promise(function(resolve) {
        (function poll() {
          if (!syncing) resolve(reconcileAll(reason, opts));
          else setTimeout(poll, 250);
        })();
      });
    }
    return Promise.resolve(0);
  }
  syncing = true;
  var statusEls = document.querySelectorAll('.sync-status-line');
  for (var s = 0; s < statusEls.length; s++) statusEls[s].textContent = 'Syncing\u2026';
  var remoteList = fetch('https://api.github.com/repos/' + state.syncRepo + '/contents/', { headers: ghHdrs() })
    .then(function(r) {
      if (r.status === 404) return []; // empty repo
      if (r.status === 401 || r.status === 403) {
        var b = document.getElementById('sync-banner');
        if (b) { b.innerHTML = '\u26a0 GitHub token rejected \u2014 renew it in <strong>Sync</strong>.'; b.style.display = 'block'; _onBannerChange(); }
        return null; // auth failure sentinel
      }
      return r.ok ? r.json() : [];
    }).catch(function() { return []; });
  var localList = fetch('/api/notes').then(function(r) { return r.ok ? r.json() : []; }).catch(function() { return []; });
  return Promise.all([remoteList, localList]).then(function(res) {
    var remote = res[0], local = res[1];
    if (remote === null) throw new Error('auth'); // banner already shown; skip success line
    var remoteMap = {};
    remote.forEach(function(e) {
      if (e && e.type === 'file' && /\.md$/.test(e.name)) remoteMap[e.name] = e.sha;
    });
    var names = {};
    Object.keys(remoteMap).forEach(function(n) { names[n] = true; });
    (local || []).forEach(function(e) { if (e && e.name) names[e.name] = true; });
    var list = Object.keys(names).filter(function(n) { return n !== state.tabletOpenNote; });
    // Sequential: gentle on the rate limit, no concurrent tablet writes.
    return list.reduce(function(chain, name) {
      return chain.then(function() { return reconcileOne(name, remoteMap[name]); });
    }, Promise.resolve()).then(function() { return list.length; });
  }).then(function(count) {
    var ts = new Date().toLocaleString();
    localStorage.setItem('ghLastSync', ts);
    var els = document.querySelectorAll('.sync-status-line');
    for (var i = 0; i < els.length; i++) els[i].textContent = 'Last synced: ' + ts;
    _onNotesChanged();
    // Tell the daemon so Lobby can show "Last sync was …" and power-sleep can proceed.
    fetch('/api/sync/ack', { method: 'POST', credentials: 'same-origin' }).catch(function() {});
    syncing = false;
    return count;
  }).catch(function() {
    // Failed or auth-rejected: don't claim a sync happened. Restore the line.
    var ls = localStorage.getItem('ghLastSync');
    var els = document.querySelectorAll('.sync-status-line');
    for (var j = 0; j < els.length; j++) {
      els[j].textContent = ls ? 'Last synced: ' + ls : 'Never synced on this device';
    }
    syncing = false;
    return 0;
  });
}

// reconcileOne: reconcile a single note given its remote sha (undefined if not
// on GitHub). Classifies via stored sha (remote change) + stored fingerprint
// (local change) into push / pull / keep-both.
function reconcileOne(name, remoteSha) {
  var hasRemote = !!remoteSha;
  return fetch('/api/notes/' + encodeURIComponent(name))
    .then(function(r) { return r.ok ? r.text() : null; })
    .then(function(localContent) {
      var hasLocal = localContent !== null;
      if (hasLocal && !hasRemote) {
        if (!localStorage.getItem('ghSha_' + name)) { return pushNote(name); }              // never synced -> new note -> push
        if (localStorage.getItem('ghLocalHash_' + name) !== strHash(localContent)) {
          return pushNote(name);                                                             // edited since last sync -> keep words, resurrect
        }
        return applyRemoteDelete(name);                                                     // synced + pristine + gone -> confirm & delete
      }
      if (!hasLocal && hasRemote) { return pullNoteAndUpdate(name); }
      if (!hasLocal && !hasRemote) { return; }
      var storedSha = localStorage.getItem('ghSha_' + name);
      var storedHash = localStorage.getItem('ghLocalHash_' + name);
      var remoteChanged = storedSha !== remoteSha;             // includes no-stored-sha
      var localChanged = storedHash !== strHash(localContent); // includes no-stored-hash
      if (remoteChanged && localChanged) { return handleClash(name, localContent); }
      if (localChanged) {
        var emptyHash = strHash('');
        if (localContent === '' && storedHash && storedHash !== emptyHash && hasRemote) {
          return pullNoteAndUpdate(name);
        }
        return pushNote(name);
      }
      if (remoteChanged) { return pullNoteAndUpdate(name); }
      return; // both unchanged
    })
    .catch(function() {});
}

// startSyncPoll: periodic safety-net reconcile so notes edited on the laptop
// land without any user action. Skipped while the tablet editor holds a file
// open (server-known edit lease via openedit, not phone typingMode).
var syncPollTimer = null;
export function startSyncPoll() {
  if (syncPollTimer) return;
  syncPollTimer = setInterval(function() {
    if (syncReady() && !state.tabletOpenNote) { reconcileAll('poll'); }
  }, SYNC_POLL_MS);
}

// verifyGitHubRepo: probe GET /repos/{owner}/{repo} with the token and report
// a plain-language result into statusEl. A 200 confirms both the token works
// and it can see the repo -- the exact thing that was silently unconfirmed before.
export function verifyGitHubRepo(repo, token, statusEl) {
  statusEl.style.color = '#888';
  statusEl.textContent = 'Verifying with GitHub\u2026';
  fetch('https://api.github.com/repos/' + repo, {
    headers: { 'Authorization': 'Bearer ' + token, 'Accept': 'application/vnd.github.v3+json' }
  }).then(function(r) {
    if (r.status === 200) {
      statusEl.style.color = '#4caf50';
      statusEl.textContent = '\u2713 Connected \u2014 syncing your notes\u2026';
      // First-connect reconcile: turn a just-verified token into an actual sync.
      reconcileAll('after verify').then(function(n) {
        statusEl.style.color = '#4caf50';
        statusEl.textContent = '\u2713 Connected \u2014 synced ' + n + ' note' + (n === 1 ? '' : 's') + ' with ' + repo + '.';
      });
    } else if (r.status === 401 || r.status === 403) {
      statusEl.style.color = '#e57373';
      statusEl.textContent = '\u2717 Token rejected \u2014 check it is correct and not expired.';
    } else if (r.status === 404) {
      statusEl.style.color = '#e57373';
      statusEl.textContent = '\u2717 Repo not found \u2014 check owner/repo, and that the token grants access to it.';
    } else {
      statusEl.style.color = '#e57373';
      statusEl.textContent = '\u2717 GitHub error (' + r.status + ').';
    }
  }).catch(function() {
    statusEl.style.color = '#e57373';
    statusEl.textContent = '\u2717 Could not reach GitHub (offline, or the phone has no internet).';
  });
}
