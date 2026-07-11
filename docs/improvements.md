# Improvements

Owner wish-list and design notes — not tracked work. Shipped features live in [../DONE.md](../DONE.md); actionable items derived here land in [../TODO.md](../TODO.md). How it works: [architecture.md](architecture.md). Why: [decisions.md](decisions.md).

**Document integrity:** the non-negotiable product contract lives in [architecture.md](architecture.md) § Document integrity and [decisions.md](decisions.md) § Document integrity. This file's § Document integrity below is the **risk matrix and implementation backlog** for that contract — not a downgrade to "wish-list".

---

## Browser vs tablet — what exists today

The phone/Mac companion (`daemon/index.html` + `app.js` + `sync.js`) is the full control surface. The e-ink Lobby (six-tab pager) and Ctrl-K picker cover most day-to-day file ops on tablet.

| Capability | Browser (authed) | Tablet (USB / Lobby) |
|---|---|---|
| List notes | Yes — scrollable list with size/date | **Files** tab — `ListView` from socket `noteslist`; note count on Home |
| Create note | **New**, **Upload**, paste-on-create | **Files → New** or **Ctrl-K** → type name → Enter |
| Open / edit | **Edit** → Type mode + `POST /api/open` | **Files → Open** / Enter / double-tap; **Ctrl-K** picker |
| Read / preview | Read view (`textContent`, safe) | Esc toggles edit/preview in Writerdeck |
| Rename | Read view → **Rename** (`PATCH /api/notes/{name}`) | **Files → Rename** or `r` (socket `renamenote`) |
| Delete | Read view → **Delete** (`DELETE`); pairs with GitHub if sync on | **Files → Delete** or `d` + confirm (socket `deletenote`) |
| Download / copy | **Download**, **Copy** (http fallback) | No |
| Paste at cursor | **Paste from here** (Type mode) | No |
| Font (read view) | Preferences → pick Inter / Literata / EB Garamond / DejaVu | Phone pushes `setfont`; Settings tab shows current font |
| PIN length | Preferences → 6 / 4 / none | Display only (Lobby shows PIN; no change on device) |
| Display rotation | Preferences → **Rotate tablet 90°** | **Ctrl-R** / Ctrl+←/→; Settings tab **Rotate 90** button |
| Sync config | **Sync** panel — on/off, repo, token (`localStorage`) | Read-only **Sync** tab |
| Sync run | Auto on connect/poll/Home/Power; **Sync now** | Triggers only if a phone browser is connected |
| Connection status | Top bar — offline / connecting / connected + battery | Not shown on e-ink |
| Show PIN on tablet | **Show PIN on tablet** (`POST /api/lobby`, pre-auth) | N/A — you are looking at the PIN |
| Exit Writerdeck | Preferences → **Exit Writerdeck** (`POST /api/shutdown`) | **Ctrl-Q** or Home from Lobby |
| Launch Lobby | **Edit** without note (starts session) | **Esc** / L+R page buttons; Mac `wd`; tablet `~/wd` |
| Keystrokes | WebSocket — layout resolved by phone OS | USB — Qt evdev, **US QWERTY default**; BT — same as browser path |

**Takeaway:** upload, download, copy, paste, GitHub token entry, and the sync engine remain browser-only. Tablet now has **Files** CRUD via trusted socket — shipped and device-verified 2026-07-11 ([decisions.md](decisions.md) #23).

---

## Document integrity — risk matrix & fix backlog (2026-07-11)

Implementation backlog for the product contract in [architecture.md](architecture.md) § Document integrity and [decisions.md](decisions.md) § Document integrity. Second-pass review after the Lobby Home wipe fix (#24).

**As built** ([architecture.md](architecture.md)): tablet-primary `.md` on disk via Go `/api/notes`; phone `sync.js` reconcile + `ghSha_*`/`ghLocalHash_*`; open tracking via editor `notifyOpen` + server `openedit` (slice 1, 2026-07-11); poll gated on `tabletOpenNote`; tablet Files CRUD over socket `req`; WebSocket `exitedit` on Home/delete; `settings.json` already write-temp-rename.

### Content fidelity — save format (separate from lifecycle)

Lifecycle hardening guards *when* writes happen; it does **not** guarantee *what* bytes land on disk. Upstream keywriter `saveFile()` writes `doc` verbatim — no plain-markdown contract.

**Observed (2026-07-11):** `likelønnsdrodling.md` saved as Qt `qrichtext` HTML (`<!DOCTYPE…`, single `<p>`), newlines collapsed, then synced to GitHub. Good markdown still in Git history (`909ff1a`). Root: `toggleMode()` (Esc) does `doc = query.text` + `saveFile()` across PlainText/RichText boundary; presentation state treated as storage format.

**Improvement:**
- **Save contract** — `saveFile()` writes plain UTF-8 from `query.text` in edit mode only; reject `qrichtext` / `<!DOCTYPE HTML` on write.
- **Load sanitizer** — `doLoad` detects HTML wrapper, strip to plain text (or refuse + log).
- **Fix `toggleMode`** — no blind `saveFile()` on every Esc; no `doc = query.text` while RichText active.
- **Server guard** — `PUT /api/notes` rejects non-markdown/HTML payloads (optional fingerprint).

### Open-file visibility

Today only **phone** `POST /api/open` sets server `currentNote` and browser `tabletOpenNote`. **Tablet** opens (Lobby Files, Ctrl-K omni) update neither. Consequences:

- **`typingMode` is phone-only, not “tablet is editing”.** Poll reconcile is skipped when the phone is in Type mode (`typingMode` true), not when the tablet has an active editor session. USB-only editing with the phone on the note list still runs the **3-minute poll** — and the same overwrite risk applies to **connect**, **reconnect**, **startup**, **Sync now**, and **sync toggle** reconciles, which also ignore tablet editor state.
- **Sync reconcile** skips only `tabletOpenNote` on the phone — not the file actually being edited on tablet. A phone on the note list (`typingMode` false) can `PUT` GitHub content over the live file on disk while the editor holds a different buffer; the next save writes the buffer back (fork or undo the pull).
- **Stale `tabletOpenNote`:** phone-back from Type mode deliberately keeps `tabletOpenNote` (so push waits for `exitedit`). If the user then opens a *different* note on the tablet, reconcile keeps skipping the *old* name and may overwrite the *current* one.
- **Phone delete/rename while the tablet is editing:** `notedeleted` / `currentNote` guards only match **phone-initiated** opens (`currentNote` set by `/api/open`). Delete from phone preview removes the file on disk; the next tablet Home/save **recreates** the old path from the in-memory buffer. **Rename has no editor notification at all** — server `currentNote` and phone `tabletOpenNote` update, but editor `currentFile` stays on the old basename, so the next save can **resurrect the old filename** even for a phone-opened note.
- **Power wake** reopens `currentNote` captured before suspend — empty after tablet-only open, so wake may land in Lobby with no note restored.

**Mitigation already in place (tablet CRUD):** Lobby **Files** delete/rename requires being in the Lobby first; `handleHome` / `showLobby` save and clear `currentFile`, so tablet-initiated destructive ops from Files do not hit the “stale buffer resurrects path” trap. The gap is **concurrent phone ops** (or reconcile pulls) while the tablet editor still holds a buffer.

**Improvement:** ~~Editor should report open file…~~ **Shipped slice 1 (2026-07-11).** ~~`notedeleted` + `noterenamed`~~ **Shipped slice 3 (2026-07-11).** Remaining: gate all reconcile triggers on `tabletOpenNote` (not only poll).

### Save & load edge cases

- **`doLoad` is async;** Home or `saveAndLoad` before XHR completes can save the wrong file or skip save when `currentFile` is still empty (softer now with edit 2b, but timing remains). Rapid note switching (Ctrl-K / Lobby Files) can interleave saves and loads.
- **`/api/open` continues on save-ack timeout** — switch may proceed before the previous note is flushed to disk.
- **Disk↔buffer drift** — reconcile `pullNoteAndUpdate`, clash resolution, and phone-side `PUT /api/notes` write disk only; the editor is never told to reload. The buffer stays stale until the user opens another note or saves (then the buffer wins and may undo the pull).
- **No atomic write** for notes — QML synchronous `PUT` and Go `os.WriteFile` write in place; power loss mid-save could truncate (unobserved but possible). `settings.json` on tablet already uses write-temp-then-rename ([architecture.md](architecture.md) Preferences bullet); notes should match server-side.
- **Editor crash / deploy kill / SIGTERM** — unsaved `query.text` buffer is lost; disk has last explicit save only.

### Sync & GitHub

- **Tablet delete/rename** updates disk only; GitHub learns on next phone reconcile — delete can **resurrect** from remote; rename can leave **two paths** until paired.
- **Real edit clashes** (both sides changed) still create `(tablet copy)` duplicates — only empty-tablet vs non-empty-remote is special-cased. Clash resolution overwrites disk without reloading the editor (same drift class as pull).
- **External GitHub edits** (VS Code, web, git) — marker-aware delete/rename helps but duplicates/resurrections remain (#19).
- **Empty-push guard** blocks accidental wipe to GitHub when phone remembers prior hash — intentional full clear may need manual re-push or hash reset.
- **`localStorage` loss** (new phone browser, cleared site data) — loses per-note `sha`/fingerprints; next reconcile may surprise-push, surprise-pull, or clash-copy.
- **Multiple browser tabs** on same origin share one sync state — concurrent reconciles are serialized but UX can confuse.
- **Power sleep** waits up to 45 s for phone reconcile ack (`/api/sync/ack`); if the phone is disconnected or reconcile fails, sleep proceeds with disk saved but GitHub may lag.

### Auth & LAN

- **PIN `none`** — anyone on Wi-Fi can `PUT`/`DELETE` notes via HTTP; integrity and confidentiality risk on untrusted LANs.

### Recovery & ops

- **`scripts/restore-wiped-notes.sh`** — recovers zero-byte files from Git history; does not cover partial truncation or duplicate cleanup beyond `(tablet copy)` names.
- **No editor↔server version/fingerprint** — cannot detect silent disk↔buffer drift without full text compare at save time.

### Suggested implementation order

Stack: **content contract** + **edit lease** + **OCC** + **atomic writes** + **conflict copies**.

1. ~~**Edit lease**~~ — shipped 2026-07-11.
2. ~~**Content fidelity**~~ — shipped 2026-07-11 (`39cbdd3`: save contract, load sanitizer, `toggleMode` fix, server HTML guard).
3. ~~**`notedeleted` + `noterenamed`**~~ — shipped 2026-07-11 (editor notified on phone rename/delete of open file; `noteDeleted` clears buffer).
4. **Reconcile policy** — skip pull/overwrite for server-known open file on all triggers.
5. **OCC on disk** — `PUT /api/notes` requires base revision / `If-Match`.
6. **Atomic note writes** on server (same pattern as `settings.json`).
7. Tablet CRUD → queued GitHub ops or Lobby “sync pending”.
8. Optional: reload or conflict banner when disk hash ≠ editor buffer.

### Models & further reading (terse)

| Pattern | Fits Writerdeck | Pointer |
|---|---|---|
| Edit lease / single writer | Shipped slice 1 | [Ink & Switch local-first](https://www.inkandswitch.com/essay/local-first/) |
| Plain-text save contract | Blocks qrichtext/HTML saves | — |
| Optimistic concurrency (`If-Match`, 412) | `ghSha` on push; extend to tablet `PUT` | [RFC 7232](https://httpwg.org/specs/rfc7232.html), [CouchDB `_rev`](https://docs.couchdb.org/en/stable/replication/conflicts.html) |
| Conflict copies + versioning | `(tablet copy).md` today; optional `.stversions` archive | [Syncthing conflicts](https://docs.syncthing.net/users/syncing.html) |
| Atomic durable write | Notes `PUT` — settings already do this | [google/renameio](https://github.com/google/renameio) |
| Remote CAS + serialize | Sequential `reconcileAll`; 409 → clash handler | [GitHub Contents API](https://docs.github.com/en/rest/repos/contents) |
| CRDT / OT | Overkill for current whole-file markdown model | [Automerge](https://automerge.org/) |

---

## On reMarkable — bring browser features to the tablet

### Lobby subpages — shipped (2026-07-11)

Six-tab pager: **Home · Files · Keyboard · Sync · Settings · Shortcuts**. Touch tabs or keyboard Tab / arrows / digits 1–6. Implementation: `lobbyPage` + `lobby_subpages.qml.inc`; Files uses `ListView` + socket `req` API via `lobby_bridge.cpp`. Device-verified; open/Home wipe bug fixed (#24 in [decisions.md](decisions.md)).

| Page | Content | Inputs |
|---|---|---|
| **Home** | Title, tagline, note count, hints | Tab / 1 |
| **Files** | Scrollable note list; New/Open/Rename/Delete | ↑/↓, Enter open, `n`/`r`/`d`, touch buttons |
| **Keyboard** | USB OTG, phone URL + PIN | Tab / 2 |
| **Sync** | `syncOn`, repo URL, last sync | Tab / 3 |
| **Settings** | Read font, rotation, **Rotate 90** button | Tab / 4 |
| **Shortcuts** | Esc, L+R, Ctrl-K/R/Q, Home | Tab / 5 |

### File CRUD on tablet — shipped (socket path)

Operations map to Writerdeck-server via trusted socket — no cookie, no LAN HTTP from the editor:

| Op | Channel today | Tablet UX |
|---|---|---|
| List | `{"t":"req","op":"noteslist"}` | Files subpage `ListView` |
| Create | `{"t":"req","op":"createnote","name":…}` | Files **New** or `n` |
| Open | `saveAndLoad(name)` (local) | Files **Open** / Enter / double-tap |
| Rename | `{"t":"req","op":"renamenote",…}` | Files **Rename** or `r` |
| Delete | `{"t":"req","op":"deletenote",…}` | Files **Delete** or `d` + confirm |
| Read-only preview | Esc in Writerdeck | Unchanged |

**Auth:** browser uses `writerdeck_token` cookie; tablet uses Unix socket `req` ops (editor is trusted — #23).

**Sync pairing:** tablet delete/rename updates disk only; GitHub catches up on next phone reconcile (#19). If sync pushes an accidental empty file, `pushNote` refuses when `ghLocalHash` was non-empty (#24).

**Not worth porting to tablet:** Upload, Download, Copy, GitHub token entry — need a file picker or clipboard; phone is the right surface.

### Password-protected / encrypted note subset

**Current model:** one global PIN gates the whole notes API and WebSocket ([decisions.md](decisions.md) #17). All `.md` files in `Writerdeck-user-documents/` are equally accessible once authed.

**Goal:** protect a *subset* of files (e.g. journal in `private/` or tagged notes) with a separate passphrase — distinct from the LAN PIN.

| Approach | Pros | Cons |
|---|---|---|
| **Encrypted subfolder** (`private/*.md.enc` or opaque blobs) | Clear separation; list dir can show locked names only | Need encrypt/decrypt in Go; editor sees plaintext only while unlocked; sync must treat ciphertext as opaque or exclude from GitHub |
| **Per-note passphrase** | Fine-grained | UX heavy on e-ink; many keys to remember |
| **Second PIN tier** | Simple mental model | Still one bucket — doesn’t split “shopping list” vs “diary” |

**Recommended direction:** encrypted subfolder, owner opt-in.

1. Notes under `Writerdeck-user-documents/.private/` (or `private/`) stored as encrypted at rest (e.g. AES-GCM via Go `crypto/*`, key derived from passphrase with scrypt/argon2id).
2. Passphrase **never stored on tablet** — entered on phone to decrypt for a session, or on tablet via a Lobby/Files unlock overlay (digits/phrase on USB keyboard).
3. **List API:** `GET /api/notes` returns `locked: true` entries without body; `GET` content requires `POST /api/notes/{name}/unlock` with passphrase (rate-limited like PIN).
4. **Browser path:** decrypt in memory for preview/edit; WebSocket forwards plaintext codepoints as today.
5. **USB path:** unchanged once note is open in Writerdeck.
6. **GitHub sync:** default **exclude** `.private/` from reconcile, or sync ciphertext only (readable on GitHub but useless without key). Document clearly.
7. **Threat model:** honest LAN guest with PIN cannot read private notes; physical access + disk image still needs passphrase. Not a HSM — proportionate for home device.

Open design choices: single folder passphrase vs per-file; timeout to re-lock; whether Lobby shows private note *titles* or only “N locked notes”.

### USB keyboard locales (Norwegian and others)

**Two input paths behave differently** ([decisions.md](decisions.md) #2):

| Path | Layout resolution | Norwegian æøå today |
|---|---|---|
| **Phone / browser → WebSocket** | `KeyboardEvent.key` — OS layout applied | **Works** (documented in [DONE.md](../DONE.md)) |
| **USB keyboard → Qt QPA evdev** | Default **US QWERTY** in Writerdeck | **Broken** — wrong letters for NO, DE, etc. |

Qt on reMarkable does **not** use `loadkeys`, `localectl`, or `setxkbmap`. Keyboard goes through `QT_QPA_EVDEV_KEYBOARD_PARAMETERS` and a **`.qmap`** file ([remarkable-keywriter#1](https://github.com/dps/remarkable-keywriter/issues/1) — community-verified on rM1/rM2 for keywriter).

**Proven workaround (upstream):**

```bash
# On a Linux host with xkb:
ckbcomp -layout no > no.kmap
# Build or obtain kmap2qmap from Qt 5.15 (Eeems published a binary — see issue #1)
./kmap2qmap no.kmap no.qmap
# Deploy to tablet, then in Writerdeck-launcher.sh:
export QT_QPA_EVDEV_KEYBOARD_PARAMETERS="/dev/input/eventN:grab=1:keymap=/home/root/keymaps/no.qmap"
```

**Caveats:**

- **`grab=1`** dedicates the keyboard to Writerdeck until exit — acceptable here (xochitl is stopped during sessions).
- **Event node varies** (`event2`, `event3`, …) — may need hotplug match or the same rescan logic Writerdeck-server uses for Escape launch (`findKeyboardInputDevices`).
- **reMarkable 1** uses standard Qt linuxfb evdev, not rM2 folio’s custom epaper-qpa keymap — qmap path applies to Writerdeck, not stock xochitl.
- Ship **multiple qmaps** (e.g. `us`, `no`, `de`) and persist choice in `settings.json` → `keyboardLayout` → launcher picks file. Lobby **Keyboard** subpage shows current layout + switch shortcut.
- **Build pipeline:** add `keymaps/` in repo; CI or maintainer script regenerates from `ckbcomp` layouts; deploy beside Qt sysroot (~tens of KB per layout).
- **AltGr / compose:** qmap includes dead keys for standard layouts; test æøå on device after deploy.

**Norwegian variants:** `no` (standard), `no(nodeadkeys)`, `no(winkeys)` — pick one default (likely `no` or `no(winkeys)` for ANSI/ISO USB keyboards).

### Edit view (unchanged wish-list)

- More VS Code-like shortcuts — e.g. newline in indented list continues indentation.
- Optional cursor block to the left of the active line.
- Headline navigation interface.
- Status bar: title, terse confirmations, zoom, time, battery (battery already on phone status bar via `/api/status`).

---

## In browser

*(No open browser items — shipped polish tracked in [DONE.md](../DONE.md).)*

Possible future browser items (lower priority than tablet parity):

- Bulk select / multi-delete in note list.
- Search across note titles and bodies.
- HTTPS / secure context for native Share sheet ([decisions.md](decisions.md) #9).
- Encrypted-folder unlock UI (paired with tablet encryption above).
