// notes-ui.js — note list, preview, typing view, paste modal.
import { state } from './state.js';
import { updateSyncBannerFromState, refreshSyncStatus } from './sync.js';
import { deps } from './deps.js';
import { grab, send, applyMode } from './connection.js';

var currentTypingFile = ''; // phone-view-only; clears on phone-back
var previewFilename = '';

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

export function showList(e) {
  if (e) e.stopPropagation();
  document.getElementById('preview').style.display = 'none';
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
export function hideTypingView(e) {
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
export function openNote(filename, displayName) {
  fetch('/api/open', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name: filename })
  }).then(function(r) {
    if (r.status === 501) { alert('Not in supervisor mode.'); return; }
    if (!r.ok) { alert('Could not open note (' + r.status + ').'); return; }
    showTypingView(displayName || filename.replace(/\.md$/, ''), filename);
  }).catch(function() { alert('Could not reach server.'); });
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

export function showPreview(filename, displayName) {
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

export function createNote(e) {
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
  });
}

function deleteNote(filename, displayName) {
  if (!confirm('Delete "' + displayName + '"? This cannot be undone.')) return;
  fetch('/api/notes/' + encodeURIComponent(filename), { method: 'DELETE' })
    .then(function(r) {
      if (!r.ok) { alert('Could not delete note.'); return; }
      if (state.tabletOpenNote === filename) { state.tabletOpenNote = ''; }
      showList();
      loadNotes();
    });
}

// ---- Paste-note modal ----

// The paste modal inserts clipboard text into the note open in the editor, at
// the cursor, via the existing keystroke pipeline (typeText) -- no new note,
// no new protocol.
export function showPasteModal() {
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

export function hidePasteModal() {
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

export function submitPaste() {
  var content = document.getElementById('paste-content').value;
  if (!content) { alert('Nothing to insert \u2014 paste some text first.'); return; }
  hidePasteModal();
  typeText(content);
  grab(); // restore trap focus so manual typing continues
}
