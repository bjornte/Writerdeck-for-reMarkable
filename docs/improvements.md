# Improvements

Owner wish-list and design notes — not tracked work. Shipped features live in [../DONE.md](../DONE.md); actionable items derived here land in [../TODO.md](../TODO.md). How it works: [architecture.md](architecture.md). Why: [decisions.md](decisions.md).

**Document integrity:** the non-negotiable product contract lives in [architecture.md](architecture.md) § Document integrity and [decisions.md](decisions.md) § Document integrity. This file's § Document integrity below is the **audit** — what is fixed, what is known-open, and what is unknown — not a wish-list.

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

## Document integrity — audit (2026-07-11)

Living audit for [architecture.md](architecture.md) § Document integrity. **Slices 1–11 are shipped** (`b1ce2bc`…`f72282d`, device-verified). This section is not a backlog — open items stay open until closed or accepted.

### Fixed (slices 1–11)

| Slice | What |
|---|---|
| 1 | Edit lease — `notifyOpen` + `openedit` WS; reconcile skips `tabletOpenNote` |
| 2 | Content fidelity — plain-markdown save contract, load sanitizer, `toggleMode` fix, server HTML guard |
| 3 | `notedeleted` / `noterenamed` — editor notified on phone rename/delete of open file |
| 4 | Reconcile policy — `openNote` in `/api/status`; `reconcileAll` gated on edit lease |
| 5 | OCC — GET `ETag`; PUT overwrite requires `If-Match`; sync sends revision |
| 6 | Atomic server writes — `writeNoteFile` temp+rename |
| 7 | Tablet CRUD → GitHub — `tabletcrud` WS + `pendingSync` queue |
| 8 | Disk↔buffer drift — `diskchanged` WS, phone drift banner, `POST /api/reload` |
| 9 | Autosave — 45 s `autosaveTimer` while editing |
| 10 | Tablet atomic saves — loopback `PUT /api/notes` → `writeNoteFile` |
| 11 | Save before deploy/stop — `POST /api/flush-save`, deploy graceful wait, SIGTERM flush |

Also under contract: empty-push guard (#19), Lobby Home wipe fix (#24). Details: [DONE.md](../DONE.md) § Document integrity.

### Known open (residual risks)

Mitigated or bounded — not eliminated. Revisit before calling integrity “closed.”

**Durability**
- SIGKILL, editor segfault, or power loss between autosave ticks (up to ~45 s of typing).
- `POST /api/flush-save` / SIGTERM path fails if editor socket down, `autosavenow` ack times out, or mismatched Writerdeck binary (no loopback save).

**Save / load timing**
- `doLoad` is async — rapid note switch (Ctrl-K, Lobby Files) can interleave saves and loads.
- `/api/open` continues on save-ack timeout — previous note may not be flushed before switch.

**Concurrency / sync**
- Real edit clashes create `(tablet copy)` duplicates; clash overwrites disk without auto-reloading the editor (drift banner is manual).
- Stale `tabletOpenNote` after phone-back can skip the wrong file in reconcile until status poll refreshes.
- External GitHub edits (VS Code, web, git) — marker-aware delete helps; duplicates/resurrections still possible (#19).
- `localStorage` loss (new browser, cleared site data) — surprise push/pull/clash on next reconcile.
- Multiple browser tabs share one sync state — serialized but confusing.
- Power sleep: up to 45 s wait for phone reconcile ack; GitHub may lag if phone offline.

**Auth / ops**
- PIN `none` — anyone on LAN can mutate notes (integrity + confidentiality on untrusted Wi-Fi).
- `restore-wiped-notes.sh` — Git history only; no partial-truncation or duplicate cleanup beyond `(tablet copy)` names.

### Unknown (unbounded)

No claim that the threat surface is complete.

- Bugs in paths not stress-tested (rapid switching, multi-device, sleep/wake edge cases, first boot with sync on).
- Firmware OTA side effects on deployed binaries or unit file.
- Qt / QML regressions not caught by `test-edit-session.sh` (save paths, autosave, loopback PUT).
- Novel failure modes (network partition mid-save, partial HTTP write despite atomic rename, future feature regressions).

**Process:** before shipping note/save/sync/lifecycle changes, ask: *can this lose text, write wrong bytes, or overwrite without the user knowing?* If yes, it does not ship until mitigated or **explicitly accepted** by the owner and logged here.

### Reference patterns

| Pattern | Status | Pointer |
|---|---|---|
| Edit lease | Shipped | [Ink & Switch local-first](https://www.inkandswitch.com/essay/local-first/) |
| Plain-text save contract | Shipped | — |
| Optimistic concurrency | Shipped (HTTP) | [RFC 7232](https://httpwg.org/specs/rfc7232.html) |
| Conflict copies | Partial (`(tablet copy)`) | [Syncthing conflicts](https://docs.syncthing.net/users/syncing.html) |
| Atomic durable write | Shipped | [google/renameio](https://github.com/google/renameio) |
| CRDT / OT | Not planned | [Automerge](https://automerge.org/) |

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

**Sync pairing:** tablet socket CRUD queues `pendingSync` + `tabletcrud` WS (slice 7); phone pairs `ghDelete`/`pushNote` immediately or on connect. Phone HTTP CRUD still pairs inline. If sync pushes an accidental empty file, `pushNote` refuses when `ghLocalHash` was non-empty (#24).

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
