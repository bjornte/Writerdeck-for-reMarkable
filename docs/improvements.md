# Improvements

Owner wish-list and design notes — not tracked work. Shipped features live in [../DONE.md](../DONE.md); actionable items derived here land in [../TODO.md](../TODO.md). How it works: [architecture.md](architecture.md). Why: [decisions.md](decisions.md).

---

## Browser vs tablet — what exists today

The phone/Mac companion (`daemon/index.html` + `app.js` + `sync.js`) is the full control surface. The e-ink Lobby and Ctrl-K picker cover a subset.

| Capability | Browser (authed) | Tablet (USB / Lobby) |
|---|---|---|
| List notes | Yes — scrollable list with size/date | Partial — note count in Lobby; full list only via **Ctrl-K** omni picker |
| Create note | **New**, **Upload**, paste-on-create | **Ctrl-K** → type new name → Enter |
| Open / edit | **Edit** → Type mode + `POST /api/open` | **Ctrl-K** → pick note; boot opens last session |
| Read / preview | Read view (`textContent`, safe) | Esc toggles edit/preview in Writerdeck |
| Rename | Read view → **Rename** (`PATCH /api/notes/{name}`) | No |
| Delete | Read view → **Delete** (`DELETE`); pairs with GitHub if sync on | No (phone delete pushes `notedeleted` if that file was open) |
| Download / copy | **Download**, **Copy** (http fallback) | No |
| Paste at cursor | **Paste from here** (Type mode) | No |
| Font (read view) | Preferences → pick Inter / Literata / EB Garamond / DejaVu | Phone pushes `setfont`; no on-tablet picker |
| PIN length | Preferences → 6 / 4 / none | Display only (Lobby shows PIN; no change on device) |
| Display rotation | Preferences → **Rotate tablet 90°** | **Ctrl-R** / Ctrl+←/→ in Lobby and preview |
| Sync config | **Sync** panel — on/off, repo, token (`localStorage`) | Read-only status in Lobby Syncing section |
| Sync run | Auto on connect/poll/Home/Power; **Sync now** | Triggers only if a phone browser is connected |
| Connection status | Top bar — offline / connecting / connected + battery | Not shown on e-ink |
| Show PIN on tablet | **Show PIN on tablet** (`POST /api/lobby`, pre-auth) | N/A — you are looking at the PIN |
| Exit Writerdeck | Preferences → **Exit Writerdeck** (`POST /api/shutdown`) | **Ctrl-Q** or Home from Lobby |
| Keystrokes | WebSocket — layout resolved by phone OS | USB — Qt evdev, **US QWERTY default**; BT — same as browser path |

**Takeaway:** file CRUD beyond open/create, export, paste, sync engine, and most settings are browser-only by design ([decisions.md](decisions.md) #8 — phone owns files, tablet owns editing). The Lobby today is informational plus **Ctrl-K**; it does not expose rename, delete, upload, or settings changes.

---

## On reMarkable — bring browser features to the tablet

### Lobby subpages (navigation model)

The Lobby is one long `Flickable` column (Notes · Syncing · Keyboard · Shortcuts). As features grow, split into **subpages** with a simple tab or ←/→ pager — same pattern as the omni note picker (`isOmni` overlay), but for Lobby sections.

Proposed pages:

| Page | Content | Inputs |
|---|---|---|
| **Home** | Title, tagline, note count, “press → for menu” | → or digit `1` |
| **Files** | Scrollable note list; highlight selection | ↑/↓, Enter open, `n` new, `d` delete, `r` rename |
| **Keyboard** | USB OTG hint, phone URL + PIN, **layout picker** (when locales ship) | → |
| **Sync** | `syncOn`, repo URL, last sync (already in `pushLobbyInfo`) | → ; optional “sync needs phone” line |
| **Settings** | Rotation (reuse Ctrl-R), read font cycle, PIN mode summary | → |
| **Shortcuts** | Esc / Ctrl-K / Ctrl-R / Ctrl-Q / Home (today’s footer block) | → |

Implementation sketch: `property int lobbyPage: 0` (or string enum), left/right or number keys when `isLobby && !isOmni`, page indicator at bottom. Reuse `ListView` from omni for Files. Server already exposes everything needed over HTTP — Writerdeck would need a thin local client (socket command → Go performs API op, or QML `XMLHttpRequest` to `http://127.0.0.1:8000` with cookie-less local trust — prefer socket commands to avoid baking auth into the editor).

### File CRUD on tablet (feasibility)

Operations map cleanly to existing API — no new server primitives required for basic CRUD:

| Op | API today | Tablet UX option |
|---|---|---|
| List | `GET /api/notes` | Files subpage or enhanced Ctrl-K list with modified date |
| Create | `POST /api/notes` `{name}` | `n` in Files page, or Ctrl-K (already works) |
| Open | `POST /api/open` + `saveAndLoad` | Enter on Files row (already via Ctrl-K → `doLoad`) |
| Rename | `PATCH /api/notes/{old}` `{name}` | `r` → inline rename prompt (QML `TextInput` overlay) |
| Delete | `DELETE /api/notes/{name}` | `d` with confirm; server already notifies editor via `notedeleted` |
| Read-only preview | `GET /api/notes/{name}` | Optional — Esc already toggles preview for open note |

**Auth gap:** browser calls use the `writerdeck_token` cookie. A tablet-local path should use **Unix-socket or localhost commands** handled by Writerdeck-server without cookie (editor is already trusted — it holds the open file). Do not expose unauthenticated HTTP CRUD on LAN.

**Sync pairing:** delete/rename on tablet should mirror browser behaviour — if sync is on, either (a) require phone for GitHub pairing, or (b) queue ops and let the next browser reconcile ([decisions.md](decisions.md) #19). Simpler v1: tablet delete/rename updates disk only; GitHub catches up on next phone sync (same as pre-marker-aware era) with a Lobby warning.

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
