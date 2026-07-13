// notes-ui.js — note list (upload/download), typing view, paste modal.
import { state } from './state.js';
import { updateSyncBannerFromState, refreshSyncStatus } from './sync.js';
import { deps } from './deps.js';
import { grab, send, applyMode } from './connection.js';

var currentTypingFile = ''; // phone-view-only; clears on phone-back

// ---- Notes API ----

export function loadNotes() {
  fetch('/api/notes')
    .then(function(r) {
      if (r.status === 401) { deps.showPinScreen(); return null; }
      if (!r.ok) throw new Error('server error');
      return r.json();
    })
    .then(function(notes) {
      if (notes !== null) {
        renderNotes(notes);
        refreshSyncStatus().then(function(s) { updateSyncBannerFromState(s); });
      }
    })
    .catch(function(err) { console.error('loadNotes failed:', err); });
}

export function renderNotes(notes) {
  var el = document.getElementById('notes-items');
  if (!notes || notes.length === 0) {
    el.innerHTML = '<div id="notes-empty">No notes yet \u2014 create one on the tablet Files tab, or Upload here.</div>';
    return;
  }
  el.innerHTML = '';
  notes.forEach(function(note) {
    var displayName = note.name.replace(/\.md\.enc$/, '').replace(/\.md$/, '');
    var row = document.createElement('div');
    row.className = 'note-row';

    var nameEl = document.createElement('span');
    nameEl.className = 'note-name';
    if (note.encrypted) {
      nameEl.textContent = (note.locked ? '[locked] ' : '[private] ') + displayName;
    } else {
      nameEl.textContent = displayName;
    }
    nameEl.title = note.name;

    var dlBtn = document.createElement('button');
    dlBtn.className = 'note-btn';
    dlBtn.textContent = 'Download';
    dlBtn.disabled = !!(note.encrypted && note.locked);
    dlBtn.addEventListener('click', function(e) {
      e.stopPropagation();
      downloadNote(note.name);
    });

    row.appendChild(nameEl);
    row.appendChild(dlBtn);
    el.appendChild(row);
  });
}

export function showList(e) {
  if (e) e.stopPropagation();
  document.getElementById('typing').style.display = 'none';
  document.getElementById('notes').style.display = '';
}

// followTabletOpen: tablet opened a note (openedit) -- mirror it on the phone so
// keystrokes forward immediately. Skips when already typing that file, or when an
// overlay (PIN, settings, paste) has focus.
export function followTabletOpen(filename) {
  if (!filename) return;
  if (state.typingMode && currentTypingFile === filename) return;
  showTypingView(filename.replace(/\.md$/, ''), filename);
}

// showTypingView: hide everything else, show the "typing on the tablet" panel.
// noteName is the display name (without .md) shown in the header.
export function showTypingView(noteName, filename) {
  document.getElementById('notes').style.display = 'none';
  document.getElementById('typing-note').textContent = noteName;
  document.getElementById('typing').style.display = 'block';
  currentTypingFile = filename || '';
  state.tabletOpenNote = filename || '';
  state.typingMode = true;
  applyMode();
}

// hideTypingView: leave the editor session running on the tablet, just
// return the phone to the note list. The session continues on e-ink.
export function hideTypingView(e) {
  if (e) e.stopPropagation();
  currentTypingFile = '';
  state.typingMode = false;
  applyMode();
  document.getElementById('typing').style.display = 'none';
  document.getElementById('notes').style.display = '';
  loadNotes();
}

function downloadNote(filename) {
  var url = '/api/notes/' + encodeURIComponent(filename) + '/download';
  fetch(url, { credentials: 'same-origin' })
    .then(function(r) {
      if (r.status === 401) { deps.showPinScreen(); return null; }
      if (r.status === 423) {
        alert('Enter private PIN on tablet');
        return waitForVaultUnlock().then(function(ok) {
          if (!ok) { alert('Timed out waiting for tablet unlock.'); return null; }
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

function waitForVaultUnlock() {
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

export function uploadFile(file) {
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
          loadNotes();
        });
        return;
      }
      if (!r.ok) { alert('Could not upload note.'); return; }
      loadNotes();
    }).catch(function() { alert('Could not reach server.'); });
  };
  reader.readAsText(file);
}

// ---- Paste-note modal ----

export function showPasteModal() {
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

function typeText(text) {
  text = text.replace(/\r\n/g, '\n').replace(/\r/g, '\n');
  Array.from(text).forEach(function(ch) {
    if (ch === '\n') { send('Enter'); }
    else if (ch === '\t') { send('Tab'); }
    else { send(ch); }
  });
}

export function submitPaste() {
  var content = document.getElementById('paste-content').value;
  if (!content) { alert('Nothing to insert \u2014 paste some text first.'); return; }
  hidePasteModal();
  typeText(content);
  grab();
}
