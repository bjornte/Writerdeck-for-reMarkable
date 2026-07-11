// app.js -- phone UI shell: WebSocket capture, notes file-manager, settings, auth.
// Sync engine lives in sync.js; shared mutable state in state.js.
import { state } from './state.js';
import {
  reconcileAll, pushNote, pullNoteAndUpdate, ghDelete,
  updateSyncBannerFromState, verifyGitHubRepo,
  syncReady, ghToken, startSyncPoll, initSync,
  setSyncToken, clearSyncStorage,
  applyTabletCrud, clearPendingSync,
  recordEditorDiskBaseline, checkDiskDrift, notifyDiskChanged
} from './sync.js';

(function () {
  'use strict';

  var RETRY_MS = 2000;
  var ECHO_MAX = 300;

  var dot  = document.getElementById('dot');
  var msg  = document.getElementById('msg');
  var echo = document.getElementById('echo');
  var trap = document.getElementById('trap');
  var buf = '';
  var ws;
  // state.typingMode: false=Browse (list/read, no capture), true=Type (capture + echo).
  // state.tabletOpenNote: server-known open file (.md); from openedit WS or phone
  //   /api/open; clears on exitedit (post-save), so GitHub push waits for save.
  //   right file even after a phone-back. Both live in state.js so sync.js can read them.
  var currentTypingFile = ''; // phone-view-only; clears on phone-back
  var previewFilename = '';

  function applyMode() {
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
    var ss = document.getElementById('settings-screen');
    if (ss && ss.style.display !== 'none') return true;
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
      loadNotes();
      // Connect-reconcile: the trigger the event-only model lacked. Runs on
      // first connect and every reconnect (the natural "back online" hook).
      if (syncReady()) { reconcileAll('connect'); }
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
        } else if (data.type === 'exitedit') {
          var saved = state.tabletOpenNote;
          state.tabletOpenNote = '';
          state.editorDiskHash = '';
          if (state.typingMode) { hideTypingView(); }
          var awaitAck = !!data.awaitSync;
          function ackSync() {
            if (awaitAck) {
              fetch('/api/sync/ack', { method: 'POST', credentials: 'same-origin' });
            }
          }
          if (data.source === 'home' || data.source === 'power') {
            if (syncReady()) {
              reconcileAll('tablet-' + data.source, { wait: true }).finally(ackSync);
            } else {
              ackSync();
            }
          } else if (syncReady()) {
            reconcileAll('tablet-exit');
          } else if (saved) {
            pushNote(saved);
          }
        } else if (data.type === 'tabletcrud') {
          if (data.op === 'deletenote' && state.tabletOpenNote === data.name) {
            state.tabletOpenNote = '';
          } else if (data.op === 'renamenote' && data.oldName && state.tabletOpenNote === data.oldName) {
            state.tabletOpenNote = data.name || '';
          }
          loadNotes();
          if (syncReady()) {
            applyTabletCrud(data.op, data.name, data.oldName).then(clearPendingSync);
          }
        } else if (data.type === 'diskchanged') {
          notifyDiskChanged(data.name || '');
        }
        // Unknown types are silently ignored -- forward-compatible.
      } catch (e) {}
    };
  }

  // ---- Notes API ----

  function loadNotes() {
    fetch('/api/notes')
      .then(function(r) {
        if (r.status === 401) { showPinScreen(); return null; }
        if (!r.ok) throw new Error('server error');
        return r.json();
      })
      .then(function(notes) {
        if (notes !== null) { renderNotes(notes); updateSyncBannerFromState(); }
      })
      .catch(function(err) { console.error('loadNotes failed:', err); });
  }

  function renderNotes(notes) {
    var el = document.getElementById('notes-items');
    if (!notes || notes.length === 0) {
      el.innerHTML = '<div id="notes-empty">No notes yet \u2014 tap &quot;New&quot; to create one.</div>';
      return;
    }
    el.innerHTML = '';
    notes.forEach(function(note) {
      var displayName = note.name.replace(/\.md$/, '');
      var row = document.createElement('div');
      row.className = 'note-row';

      var nameEl = document.createElement('span');
      nameEl.className = 'note-name';
      nameEl.textContent = displayName;
      nameEl.title = note.name;
      nameEl.addEventListener('click', function(e) {
        e.stopPropagation();
        showPreview(note.name, displayName);
      });

      row.appendChild(nameEl);
      el.appendChild(row);
    });
  }

  function showList(e) {
    if (e) e.stopPropagation();
    document.getElementById('preview').style.display = 'none';
    document.getElementById('typing').style.display = 'none';
    document.getElementById('notes').style.display = '';
  }

  // showTypingView: hide everything else, show the "typing on the tablet" panel.
  // noteName is the display name (without .md) shown in the header.
  function showTypingView(noteName, filename) {
    document.getElementById('notes').style.display = 'none';
    document.getElementById('preview').style.display = 'none';
    document.getElementById('typing-note').textContent = noteName;
    document.getElementById('typing').style.display = 'block';
    currentTypingFile = filename || '';
    state.tabletOpenNote = filename || '';
    state.typingMode = true;
    applyMode();
  }

  // hideTypingView: leave the editor session running on the tablet, just
  // return the phone to Browse (note list). The session continues; the user
  // can re-enter it by tapping Edit on any note.
  function hideTypingView(e) {
    if (e) e.stopPropagation();
    // Phone-back only: the tablet keeps editing and has NOT saved, so we must not
    // push here (that would ship a stale file and poison the stored SHA). The push
    // happens on the tablet's post-save `exitedit` instead.
    currentTypingFile = '';
    state.typingMode = false;
    applyMode();
    document.getElementById('typing').style.display = 'none';
    document.getElementById('notes').style.display = '';
    loadNotes();
  }

  // openNote: POST /api/open to tell the tablet to save the current note and
  // open the chosen one; then switch the phone to the typing view.
  function openNote(filename, displayName) {
    // The tablet's /api/open does saveAndLoad: it SAVES the currently-open note
    // before loading the new one. Capture it so we can push that saved note.
    var prev = state.tabletOpenNote;
    function doOpen() {
      fetch('/api/open', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: filename })
      }).then(function(r) {
        if (r.status === 501) { alert('Not in supervisor mode.'); return; }
        if (!r.ok) { alert('Could not open note (' + r.status + ').'); return; }
        showTypingView(displayName || filename.replace(/\.md$/, ''), filename);
        // The note-switch just saved `prev` on the tablet -- push it to GitHub.
        if (prev && prev !== filename) { pushNote(prev); }
      }).catch(function() { alert('Could not reach server.'); });
    }
    // Skip the pre-open pull when the note is already open on the tablet:
    // keywriter's save-on-load would immediately overwrite the pulled file.
    if (!syncReady() || filename === state.tabletOpenNote) { doOpen(); return; }
    // Pull newest GitHub version to tablet first (best-effort; open regardless of outcome).
    pullNoteAndUpdate(filename).then(doOpen, doOpen);
  }

  function rewirePreviewButtons(filename, displayName) {
    function rewire(id, fn) {
      var old = document.getElementById(id);
      var fresh = old.cloneNode(true);
      old.parentNode.replaceChild(fresh, old);
      fresh.addEventListener('click', function(e) { e.stopPropagation(); fn(); });
    }
    rewire('preview-edit',   function() { openNote(filename, displayName); });
    rewire('preview-rename', function() { renameNote(filename, displayName); });
    rewire('preview-delete', function() { deleteNote(filename, displayName); });
    rewire('preview-download', function() { downloadNote(filename); });
    rewire('preview-copy',   function() {
      var label = document.getElementById('preview-title').textContent;
      doCopy(document.getElementById('preview-body').textContent, label);
    });
  }

  function showPreview(filename, displayName) {
    previewFilename = filename;
    fetch('/api/notes/' + encodeURIComponent(filename))
      .then(function(r) {
        if (!r.ok) throw new Error('not found');
        return r.text();
      })
      .then(function(text) {
        document.getElementById('preview-title').textContent = displayName;
        document.getElementById('preview-body').textContent = text;
        document.getElementById('notes').style.display = 'none';
        document.getElementById('preview').style.display = 'block';
        rewirePreviewButtons(filename, displayName);
      })
      .catch(function() { alert('Could not load note.'); });
  }

  // downloadNote: navigate the browser to the Content-Disposition:attachment endpoint.
  // The server gates it with the session cookie, which the browser sends automatically
  // on same-origin requests -- no extra auth work needed here.
  function downloadNote(filename) {
    var a = document.createElement('a');
    a.href = '/api/notes/' + encodeURIComponent(filename) + '/download';
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
  }

  // doCopy: write text to the clipboard, falling back to execCommand on plain-http
  // (navigator.clipboard requires a secure context; the Companion runs on http LAN).
  function doCopy(text, label) {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(text)
        .then(function() { alert('Copied \u201c' + label + '\u201d to clipboard.'); })
        .catch(function() { fallbackCopy(text, label); });
    } else {
      fallbackCopy(text, label);
    }
  }

  function fallbackCopy(text, label) {
    var ta = document.createElement('textarea');
    ta.value = text;
    ta.style.position = 'fixed';
    ta.style.left = '-9999px';
    ta.style.top = '0';
    document.body.appendChild(ta);
    ta.focus();
    ta.select();
    var ok = false;
    try { ok = document.execCommand('copy'); } catch (e) {}
    document.body.removeChild(ta);
    alert(ok
      ? ('Copied \u201c' + label + '\u201d to clipboard.')
      : 'Copy failed \u2014 try long-pressing the text above.');
  }

  function uploadFile(file) {
    if (!file) return;
    if (!/\.(md|markdown|txt)$/i.test(file.name)) {
      alert('Only .md, .markdown, or .txt files can be uploaded.');
      return;
    }
    if (file.size > 1024 * 1024) {
      alert('File is too large (max 1 MB).');
      return;
    }
    var reader = new FileReader();
    reader.onload = function () {
      var base = file.name.replace(/\.(md|markdown|txt)$/i, '');
      fetch('/api/notes', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: base, content: reader.result })
      }).then(function(r) {
        if (r.status === 409) {
          var alt = prompt('"' + base + '" already exists. Save as:', base + '-2');
          if (!alt || !alt.trim()) return;
          var altBase = alt.trim().replace(/\.(md|markdown|txt)$/i, '');
          fetch('/api/notes', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: altBase, content: reader.result })
          }).then(function(r2) {
            if (!r2.ok) { alert('Could not upload note.'); return; }
            var fn = altBase.endsWith('.md') ? altBase : altBase + '.md';
            openNote(fn, altBase.replace(/\.md$/, ''));
          });
          return;
        }
        if (!r.ok) { alert('Could not upload note.'); return; }
        var fn = base.endsWith('.md') ? base : base + '.md';
        openNote(fn, base.replace(/\.md$/, ''));
      }).catch(function() { alert('Could not reach server.'); });
    };
    reader.readAsText(file);
  }

  function createNote(e) {
    e.stopPropagation();
    var name = prompt('New note name:');
    if (!name || !name.trim()) return;
    fetch('/api/notes', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: name.trim() })
    }).then(function(r) {
      if (r.status === 409) { alert('"' + name.trim() + '" already exists.'); return; }
      if (!r.ok) { alert('Could not create note.'); return; }
      // Open the new note immediately so the user lands in editing, not
      // back in the (possibly empty) list. This also replaces the old
      // Write-button entry path for an empty library.
      var filename = name.trim().endsWith('.md') ? name.trim() : name.trim() + '.md';
      openNote(filename, name.trim().replace(/\.md$/, ''));
    });
  }

  function renameNote(filename, displayName) {
    var newName = prompt('Rename "' + displayName + '" to:', displayName);
    if (!newName || !newName.trim() || newName.trim() === displayName) return;
    fetch('/api/notes/' + encodeURIComponent(filename), {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: newName.trim() })
    }).then(function(r) {
      if (r.status === 409) { alert('"' + newName.trim() + '" already exists.'); return; }
      if (!r.ok) { alert('Could not rename note.'); return; }
      // Update the preview in place: new title + rebind all buttons to the new filename.
      var newFilename = newName.trim().endsWith('.md') ? newName.trim() : newName.trim() + '.md';
      var newDisplayName = newName.trim().replace(/\.md$/, '');
      previewFilename = newFilename;
      document.getElementById('preview-title').textContent = newDisplayName;
      rewirePreviewButtons(newFilename, newDisplayName);
      loadNotes();
      if (state.tabletOpenNote === filename) { state.tabletOpenNote = newFilename; }
      // Propagate the rename to GitHub: delete the old path, then push the new one.
      if (syncReady()) { ghDelete(filename).then(function() { pushNote(newFilename); }); }
    });
  }

  function deleteNote(filename, displayName) {
    if (!confirm('Delete "' + displayName + '"? This cannot be undone.')) return;
    fetch('/api/notes/' + encodeURIComponent(filename), { method: 'DELETE' })
      .then(function(r) {
        if (!r.ok) { alert('Could not delete note.'); return; }
        if (state.tabletOpenNote === filename) { state.tabletOpenNote = ''; }
        ghDelete(filename); // propagate to GitHub so it doesn't resurrect on next pull
        showList();
        loadNotes();
      });
  }

  // ---- Paste-note modal ----

  // The paste modal inserts clipboard text into the note open in the editor, at
  // the cursor, via the existing keystroke pipeline (typeText) -- no new note,
  // no new protocol.
  function showPasteModal() {
    var contentEl = document.getElementById('paste-content');
    contentEl.value = '';
    document.getElementById('paste-modal').style.display = 'flex';
    // Try to pre-fill from clipboard; may prompt the user on iOS Safari.
    if (navigator.clipboard && navigator.clipboard.readText) {
      navigator.clipboard.readText()
        .then(function(text) {
          if (text) { contentEl.value = text; }
        })
        .catch(function() {}); // ignore -- user can long-press Paste manually
    }
    contentEl.focus();
  }

  function hidePasteModal() {
    document.getElementById('paste-modal').style.display = 'none';
  }

  // typeText: replay a string into the editor as key events, so it lands at the
  // current cursor in keywriter -- no new protocol, just the existing keystroke
  // path. Iterates by code point (Array.from) so emoji/astral chars stay intact.
  function typeText(text) {
    text = text.replace(/\r\n/g, '\n').replace(/\r/g, '\n');
    Array.from(text).forEach(function(ch) {
      if (ch === '\n') { send('Enter'); }
      else if (ch === '\t') { send('Tab'); }
      else { send(ch); }
    });
  }

  function submitPaste() {
    var content = document.getElementById('paste-content').value;
    if (!content) { alert('Nothing to insert \u2014 paste some text first.'); return; }
    hidePasteModal();
    typeText(content);
    grab(); // restore trap focus so manual typing continues
  }

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

  function updateBannerOffset() {
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

  function showSettings() {
    document.getElementById('settings-screen').style.display = 'flex';
    loadSettingsPanel();
  }

  function hideSettings() {
    document.getElementById('settings-screen').style.display = 'none';
    grab();
  }

  function showSync() {
    document.getElementById('sync-screen').style.display = 'flex';
    loadSyncPanel();
  }

  function hideSync() {
    document.getElementById('sync-screen').style.display = 'none';
    grab();
  }

  function wireOverlayDismiss(screenId, boxId, hideFn) {
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
          if (ws) { ws.close(); }
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
          updateSyncBannerFromState();
          if (newVal && syncReady()) { reconcileAll('toggle'); }
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
      tokLabel.textContent = 'GitHub token (stays on this device \u2014 never sent to the tablet)';
      list.appendChild(tokLabel);

      var tokInput = document.createElement('input');
      tokInput.type = 'password'; tokInput.className = 'token-input';
      tokInput.style.width = '100%';
      tokInput.placeholder = 'github_pat_\u2026 or ghp_\u2026'; tokInput.autocomplete = 'off';
      if (ghToken()) tokInput.value = '\u2022'.repeat(16);
      var tokTouched = false;
      tokInput.addEventListener('focus', function() {
        if (!tokTouched) { tokInput.value = ''; tokTouched = true; }
      });
      list.appendChild(tokInput);

      // One primary action for the whole section: save repo (to tablet) + token
      // (to this browser), then verify both against GitHub so you get a clear
      // yes/no instead of a silent green flash.
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
        // Persist the token only if the field was actually edited (an untouched
        // masked field means "keep the existing token"); never wipe on blank here.
        if (tokTouched && tokInput.value.trim()) {
          setSyncToken(tokInput.value.trim());
          tokInput.value = '\u2022'.repeat(16); tokTouched = false;
        }
        var token = ghToken();
        // Save the repo to the tablet first (non-secret), then verify.
        fetch('/api/settings', {
          method: 'POST', headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({ syncRepo: repoVal })
        }).then(function(r) {
          if (!r.ok) {
            verifyLine.style.color = '#e57373';
            verifyLine.textContent = '\u2717 Invalid repo \u2014 use owner/repo format.';
            return;
          }
          state.syncRepo = repoVal;
          updateSyncBannerFromState();
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
          if (!repoVal || !token) {
            verifyLine.style.color = '#888';
            verifyLine.textContent = 'Saved. Enter both repo and token to verify.';
            return;
          }
          verifyGitHubRepo(repoVal, token, verifyLine);
        }).catch(function() {
          verifyLine.style.color = '#e57373';
          verifyLine.textContent = '\u2717 Could not reach the tablet to save.';
        });
      });

      var clearBtn = document.createElement('button');
      clearBtn.className = 'sync-btn-secondary'; clearBtn.textContent = 'Clear token';
      clearBtn.addEventListener('click', function(e) {
        e.stopPropagation();
        if (!confirm('Remove GitHub token from this device?')) return;
        clearSyncStorage();
        tokInput.value = ''; tokTouched = true;
        verifyLine.style.color = '#888';
        verifyLine.textContent = 'Token removed from this device.';
        updateSyncBannerFromState();
      });

      var syncNowBtn = document.createElement('button');
      syncNowBtn.className = 'sync-btn-secondary'; syncNowBtn.textContent = 'Sync now';
      syncNowBtn.addEventListener('click', function(e) {
        e.stopPropagation();
        if (!syncReady()) {
          verifyLine.style.color = '#e57373';
          verifyLine.textContent = '\u2717 Save a repo and token first.';
          return;
        }
        verifyLine.style.color = '#888';
        verifyLine.textContent = 'Syncing all notes\u2026';
        reconcileAll('manual').then(function(n) {
          verifyLine.style.color = '#4caf50';
          verifyLine.textContent = '\u2713 Synced ' + n + ' note' + (n === 1 ? '' : 's') + '.';
        });
      });

      actionRow.appendChild(saveBtn); actionRow.appendChild(syncNowBtn); actionRow.appendChild(clearBtn);
      list.appendChild(actionRow);
      list.appendChild(verifyLine);

      var statusLine = document.createElement('div');
      statusLine.className = 'sync-status-line';
      statusLine.style.cssText = 'font-size:12px;color:#888;padding:4px 2px;';
      var ls = localStorage.getItem('ghLastSync');
      statusLine.textContent = ls ? 'Last synced: ' + ls : 'Never synced on this device';
      list.appendChild(statusLine);
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

  function showPinScreen() {
    stopStatusPoll();
    document.getElementById('pin-screen').style.display = 'flex';
    startPinPoll();
  }

  function hidePinScreen() {
    stopPinPoll();
    document.getElementById('pin-screen').style.display = 'none';
    document.getElementById('settings-btn').style.display = 'block';
    document.getElementById('sync-btn').style.display = 'block';
    startStatusPoll();
  }

  // loadSyncConfig: pull the non-secret syncOn/syncRepo flags at startup, not
  // just when the Settings panel opens. Without this, connect-reconcile and the
  // poll both saw the default `false` at page load and silently skipped -- the
  // note would connect fine yet read "never synced". Runs post-auth; on success
  // it kicks an immediate reconcile so first connect actually moves files.
  function loadSyncConfig() {
    return fetch('/api/settings')
      .then(function(r) { return r.ok ? r.json() : null; })
      .then(function(data) {
        if (!data) return;
        state.syncOn = !!data.syncOn;
        state.syncRepo = data.syncRepo || '';
        updateSyncBannerFromState();
        if (syncReady()) { reconcileAll('startup'); }
      })
      .catch(function() {});
  }

  function checkAuthAndInit() {
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

  function submitPIN(e) {
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
    startSyncPoll();
  });
}());
