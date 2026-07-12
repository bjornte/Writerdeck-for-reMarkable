# Improvements

Open wish-list and design notes — not tracked work. Shipped features → [../DONE.md](../DONE.md). Actionable items → [../TODO.md](../TODO.md). Capability matrix → [browser-vs-tablet.md](browser-vs-tablet.md). Integrity audit → [integrity-audit.md](integrity-audit.md).

---

## Encrypted / password-protected note subset

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

---

## USB keyboard locales (Norwegian and others)

**Two input paths behave differently** ([decisions.md](decisions.md) #2):

| Path | Layout resolution | Norwegian æøå today |
|---|---|---|
| **Phone / browser → WebSocket** | `KeyboardEvent.key` — OS layout applied | **Works** ([DONE.md](../DONE.md)) |
| **USB keyboard → Qt QPA evdev** | Default **US QWERTY** in Writerdeck | **Broken** — wrong letters for NO, DE, etc. |

Qt on reMarkable does **not** use `loadkeys`, `localectl`, or `setxkbmap`. Keyboard goes through `QT_QPA_EVDEV_KEYBOARD_PARAMETERS` and a **`.qmap`** file ([remarkable-keywriter#1](https://github.com/dps/remarkable-keywriter/issues/1)).

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
- Ship **multiple qmaps** (e.g. `us`, `no`, `de`) and persist choice in `settings.json` → `keyboardLayout` → launcher picks file. Lobby **Keyboard** tab is the picker (tablet-only; phone duplicate removed).
- **Build pipeline:** add `keymaps/` in repo; CI or maintainer script regenerates from `ckbcomp` layouts; deploy beside Qt sysroot (~tens of KB per layout).
- **AltGr / compose:** qmap includes dead keys for standard layouts; test æøå on device after deploy.

**Norwegian variants:** `no` (standard), `no(nodeadkeys)`, `no(winkeys)` — pick one default (likely `no` or `no(winkeys)` for ANSI/ISO USB keyboards).

---

## Edit view

- More VS Code-like shortcuts — e.g. newline in indented list continues indentation.
- Optional cursor block to the left of the active line.
- Headline navigation interface.
- Status bar: title, terse confirmations, zoom, time, battery (battery already on phone status bar via `/api/status`).

---

## Browser

- Bulk select / multi-delete in note list.
- Search across note titles and bodies.
- HTTPS / secure context for native Share sheet ([decisions.md](decisions.md) #9).
- Encrypted-folder unlock UI (paired with tablet encryption above).
