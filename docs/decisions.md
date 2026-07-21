# Architecture decisions

Why the project is built this way — document integrity and core editing first, device quirks later. How: [architecture.md](architecture.md). Open work: [../TODO.md](../TODO.md). Finished: [../DONE.md](../DONE.md). Gotchas: [lessons.md](lessons.md). Words: [terms.md](terms.md).

Numbers were renumbered 17 Jul 2026; prefer titles over old § IDs.

---

## Document integrity

Writerdeck is a typewriter. Your prose must survive as plain Markdown on disk. That rule comes before every feature. Full audit: [integrity-audit.md](integrity-audit.md).

Notes are UTF-8 Markdown, never Qt HTML. An open note must not be silently overwritten by sync or remote delete. Saves use defined paths, autosave about every 45 seconds, and save-before-stop on deploy; a hard kill can still lose recent typing. If the file on disk changes under you, reload or show a conflict — do not let a stale buffer win. GitHub sync backs up; it must not empty-push or delete against a live edit.

No change to the daemon, sync, build script, or note APIs ships without an integrity pass.

## Typing-test strategy is failing

We do **not** have a working automated typing-test strategy. We have an attempt — harness, dialect notes, critical tags, scenario counts — and it keeps failing the same way: suites go green while very basic editing bugs still reach a person typing. A full or critical green run is a score for those checks, not proof that basic editing works.

Do not remove this section (or the matching banners in [editor-testing/todo.md](editor-testing/todo.md) and [lessons.md](lessons.md)) until there is solid proof those misses have stopped recurring — meaning later human finds are ones the suite would already have failed. Working theory: [editor-testing/methodology-shortcomings.md](editor-testing/methodology-shortcomings.md). Claim×kill-test inventory (keep updated): [editor-testing/basic-claims.md](editor-testing/basic-claims.md). How we run the checks today: §13.

## Device verification

A successful deploy script is not enough. Rebuild when the QML changed, deploy, relaunch, read the tablet journal. Fail on QML parse errors or an editor that exits at once. After QML changes run the edit-session check. After caret or selection work run the automated typing tests (§13), remembering the strategy above is still failing. After Lobby or Home run the Lobby keyboard test (§15).

## 19. GitHub note-sync is a non-authoritative reconciler — delete/rename only from the browser
Status: built (Phase 9), device-verified. The optional two-way GitHub sync (off by default; the PAT lives only in the phone browser's `localStorage`, never on the tablet, which holds just non-secret `syncOn`/`syncRepo`) deliberately reconciles by *unioning* the tablet's and the repo's note lists and copying any note missing from one side to the other — it never deletes on its own. That is the intended safety property: a reconciler that cannot delete cannot lose a note, which sidesteps the documented real-git-on-mobile instability (isomorphic-git crashes / packfile corruption on low-RAM devices — the reason we sync via GitHub's plain Contents API, not real git). Accepted cost: destructive ops must go through the *phone browser*, which pairs them — the UI's Delete also calls `ghDelete` (Contents API DELETE + stored `sha`) and Rename deletes the old path then pushes the new. A delete or rename made *outside* the browser — VS Code, `git`, the GitHub web UI — reads only as "a note is missing from one side" and gets resurrected (or, for rename, duplicated) on the next sync. Not a bug; the price of never-delete-on-its-own. Cheap partial fix, now built (Opus-reviewed; browser/device-verify pending) in [../daemon/sync.js](../daemon/sync.js): the reconciler treats "on the tablet + carries a stored `sha` marker + pristine (unchanged since last sync) + now gone from GitHub" as a real delete rather than a new-file pull — reusing the per-note `sha`/hash already stored. It covers the common GitHub-side / committed-VS-Code delete, and (as a side effect) external *rename*, which reconciles as delete-old + pull-new instead of duplicating. Two safety invariants keep it from ever losing words: an unpushed local edit resurrects (push) rather than deletes, and the delete only fires after a fresh per-note `GET` returns 404 — the guard against `reconcileAll` mapping a transient remote-list failure to `[]` and mass-deleting the tablet (a false positive self-heals via the next pull). Still (correctly) leaves a purely-local unpushed delete alone: the tablet is a low-trust cache, so GitHub stays authoritative for deletes.

---

## 1. Plain-text edit mode

While you type, the editor stays plain text — raw Markdown on screen and in memory. Esc preview may look fancy. Headings and bold belong there, not in edit mode.

Rich text in edit was tried. Pulling formatted text back into the note produced empty files, HTML on disk, and broken previews. Real hide-the-markers editing stays out of scope unless someone writes a new decision and re-proves integrity.

## 2. Display sync after Lobby clears the buffer

Going Home clears the on-screen box. Reload and mode toggle must push the real note text (or its preview HTML) back onto the screen — never read fancy HTML back into the note. Without that, Home can save a zero-length file.

## 3. Socket input, not uinput

This kernel cannot load a uinput (fake keyboard device). Do not retry that path. Keystrokes arrive on a local socket and become Qt key events. The phone resolves the keymap; the tablet gets characters and Mac-style Ctrl/Alt shortcuts (Ctrl stands for ⌘). USB keyboards use that same set — no second map. Binding authority for tests: Apple Cocoa prose first; CodeMirror for Home/End caret-to-visual-line (stock Apple Home only scrolls); Qt’s Ctrl+arrow = word is not the target. Details: [editor-testing/scenario-catalog.md](editor-testing/scenario-catalog.md).

## 4. Owned keywriter fork

The editor is Singleton’s keywriter, rebuilt as Writerdeck from our fork. The old binary does not load on current firmware. We build it in CI (GitHub Actions) with a Qt sysroot and draw to the real framebuffer (`/dev/fb0`).

Fork: [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter). CI clones it; the build script only checks and compiles. New editor behavior goes in the fork, assembled with `./assemble-qml.sh` into committed `main.qml`. Math, undo, shortcuts, and wrapped-line motion live in C++ `EditHelper`; QML draws and applies. Migrations: [editor-migration-1-to-QML](editor-migration-1-to-QML/todo-handoff-keywriter-fork.md), [editor-migration-2-to-cpp](editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md). Keep §5–§6.

Pull from Singleton’s original on purpose, not every session. Histories are linked (`5946cae`); ordinary merges work. After a merge: rebuild, deploy, edit-session, then the 38 basic typing checks.

## 5. Keep wrap gaps and custom EditHelper undo

Qt’s text box does not give clean “visual line” indexes. Wrapped Up/Down therefore uses small hand-tuned pixel gaps and a sticky horizontal goal. That is not elegant — but it matches Mac/Linux motion and the wrap tests. Change the numbers only if those tests start failing.

Undo is ours in `EditHelper`, not Qt’s built-in stack. Socket typing and shortcut deletes must share one history. Redesign only if integrity or undo tests force it. Do not replace Qt’s text box for purity — §6.

## 6. Do not fork Qt’s text box

Moving our helpers into the fork kept the same typing behavior under passing automated tests — evenings, not months. Replacing Qt’s text box is a different job.

That box is open source, but the useful slice is tens of thousands of lines inside Qt. Forking it means maintaining a patched Qt forever. We already pass integrity and editing on device with §5. Stop here unless stock TextEdit cannot fix a real product bug, with tests ready to catch regressions.

## 7. No Toltec

Toltec pins firmware and can leave the tablet unable to update normally on unsupported versions. That fights OTA (over-the-air updates). Skip it unless the owner accepts the lock.

## 8. Static Go binary

The server is one static ARM binary. No Python or extra runtimes on the tablet. It builds on the Mac in about a second and survives firmware updates.

## 9. Always-on server, on-demand editor

The server keeps the phone reachable under the stock UI. Boot stays on xochitl; Writerdeck opens on demand (page buttons, USB Esc, `/api/lobby` / `wd`). Only the editor session flips xochitl. Keep-awake covers the editor child, not the whole device. Units that start from `/home/root` must wait for that mount or cold boot fails with `203/EXEC`.

## 10. Companion model

No phone app — the tablet is the server. Day-to-day files and settings live in the Lobby; the phone is for typing, paste-at-cursor, accepting a tablet Download offer, the sync token, and launching the Lobby when closed (“Show PIN on tablet”). Home twice: edit to Lobby, Lobby to quit.

## 11. GitHub sync copies missing notes both ways

Sync runs on the tablet. The token sits in the browser and in tablet RAM — never on disk. After restart the phone can repost it. Sync is change-driven (boot, app open, document open, wake, Wi-Fi up, Home, power sleep, CRUD, Sync now, token verify) — not a timer. Boot, app open, wake, Wi-Fi up, and Sync now always list GitHub so remote-only edits pull without tapping Sync; Home, power sleep, and token verify skip the network when nothing local is dirty. Opening a note from the Lobby pulls that file first. Unchanged notes and vault secrets are not pushed, so clean runs do not create empty GitHub commits. The journal tells a short story (`sync: pushed …`, coalesced `nothing to do` streaks) rather than one line per quiet skip. Sync copies missing notes both ways; it does not delete on its own. Empty push over a known-good note is refused. Details: [server-sync-implementation.md](server-sync-implementation.md).

## 12. Optional at-rest encryption (private notes)

Pairing PIN and vault PIN are separate. Encrypted notes are `.md.enc` beside plain `.md`. Each file is sealed; the vault PIN unlocks a key held in RAM only while that note (or a one-shot encrypt/decrypt) is active. PIN every open, including note switches. Tablet-only entry. Per-note encrypt/decrypt — no bulk lock. Sync copies encrypted files without reading them and mirrors recovery material under `secret/`. Markdown integrity applies to `.md` only. See [integrity-audit.md](integrity-audit.md).

## 13. Automated typing tests (how we run them)

Harness and scenario mechanics only — strategy status is above under **Typing-test strategy is failing**. We run caret and selection checks on the real tablet, over the same path the phone uses — not by reading saved files. Counts and pass/fail log: [editor-testing/](editor-testing/). Test notes use the `z-test-` prefix (§32). USB layout quirks still need a human check after qmap changes.

## 14. Edit-session check

Opening a note from outside must keep the editor up for several seconds. Instant exit usually means a broken QML, not a server bug. Run after Writerdeck or QML deploy.

## 15. Lobby keyboard check

After Home from edit, Lobby keys must still work. Focus stays on Lobby after touch. The Lobby keyboard script checks the Home path. Run it with edit-session after Lobby or Home QML changes.

## 16. Physical Home: exclusive gpio grab

While Writerdeck is open, the server exclusively grabs the tablet’s Home/Power/page buttons so Qt never sees a second Home. Release on exit so the stock UI works again. Note → Lobby and Lobby → stock UI are two bindings in `lobby-ui.json` (`global.toLobby`, `global.quit`); both default to `hardware_home` (the physical middle button). The USB / phone keyboard Home key is not that button — it keeps caret Home/End and must not drive those two paths. Idle page-button launch (§23) needs no grab. Handoff: [todo-handoff-physical-home-input.md](todo-handoff-physical-home-input.md). Chord defaults: [todo-lobby-ui-shortcuts.md](todo-lobby-ui-shortcuts.md), `config/lobby-ui.json`.

## 17. PIN and per-IP lockout

Six digits, four, or none. None means anyone on the Wi-Fi can hit the notes API — the UI warns. Five wrong guesses lock that IP for a minute, not the whole device. Store length as `"6"`, `"4"`, or `"none"`.

## 18. Show PIN on tablet

A second phone mid-edit sees the note, not the PIN. Ask the tablet to return to Lobby so the PIN is readable on e-ink. That call is pre-auth and rate-limited; it reveals nothing over the wire.

## 19. Tablet file create/rename/delete via trusted socket

Lobby file ops use the same local socket as keystrokes. The server does the disk work and can nudge the phone. Launch Lobby with `wd` on the Mac or `~/wd` on the tablet. A create or rename that hits an existing name fails on the server and shows the message inside the New / Rename / New encrypted dialog (same pattern as a wrong vault PIN staying on the pad). Names are unique by case-insensitive title: `Doc` and `doc` collide, and a plain note may not share that title with an encrypted twin (`doc.md` vs `doc.md.enc`). The Files list is sorted the same way. Other Files feedback (vault open errors, “open the phone page” for Download) still uses the dismissable box above the list.

## 20. Bluetooth remote key capture on the phone

Bluetooth pairs to the phone. After connect, the phone stays on the keyboard shell and captures keys by default (Lobby navigation included). Cmd/Ctrl+R/T/W/N/L still pass through so the browser stays manageable — phones also reserve Cmd/Ctrl+1–9 for browser tabs. Overlays (PIN, paste, download offer, sync setup) pause capture. Lobby Ctrl-letter chords avoid those reserved combos (§37).

## 21. Lobby Files: Edit and Read

Two opens: Edit to type, Read to preview. The phone learns which mode so a Bluetooth keyboard can still Esc into edit. Enter edits; Ctrl-V reads.

## 22. Lobby tab order (Files first)

Boot and Home land on the note list, not the welcome screen. Coming back from edit reselects the note you left. Switch Lobby pages with Tab, Shift-Tab, and plain Left/Right. Optional per-tab Ctrl-letters live in `lobby-ui.json` (`tabs.files` … `tabs.about`); ship empty. Do not hardwire digits 1–6 for tab jump — bare digits stay free for a later Finder-style note jump (§37).

## 23. Idle launch from stock UI

From the stock UI, USB Escape or both page buttons together open Writerdeck to the Lobby. Power still owns sleep/wake. While Writerdeck is already up, Escape is left to the editor.

## 24. Mac builds the server; GitHub builds Writerdeck

Deploy starts on the Mac, which can reach the tablet. The Go server builds locally. Writerdeck needs the Qt container, so CI (GitHub Actions) produces that binary.

## 25. Wi-Fi is the dev path

USB ethernet to the tablet is unused. Set `RM_HOST_WIFI` (or hotspot) in secrets and deploy over Wi-Fi SSH.

## 26. Deploy transport is gzip-over-ssh

`scp` stalls on this link. We stream gzip over SSH and check the size afterward.

## 27. Secrets in a gitignored env file

The tablet password is already on the device screen. Keeping it in a local ignored env file is fine; committing it is not.

## 28. LF and ASCII on device files

Wrong line endings or stray fancy characters break shell and systemd on the tablet. Markdown may use Unicode; scripts and device files may not.

## 29. On-device Writerdeck naming

On the tablet the names are Writerdeck, Writerdeck-server, Writerdeck-user-documents, and `/run/Writerdeck.sock`. Some repo script names are older history.

## 30. Note create stays on the safe create route

New notes (Lobby New, and any future import UI) go through the same create API — path checks and “already exists” included. The server enforces size. The old phone Upload control is gone with the phone notes list; bringing import back should reuse this route, not invent a second path.

## 31. Display rotation persists in settings

Rotation is saved on the tablet and pushed when the editor connects. Change it from Lobby Settings (chord in `lobby-ui.json`). Phone Preferences no longer rotate.

## 32. Device test note names

Automated tests use filenames starting with `z-test-` so they sort last and stay obvious in the Files list.

## 33. reMarkable 1 first; rM2 only if the community wants it

Writerdeck targets the reMarkable 1. Install and docs say so on purpose.

The editor draws through linuxfb on `/dev/fb0` with the epaper scene graph — that path is rM1. On rM2 the panel is driven differently, so the same binary does not light the screen. Community shims such as rm2fb usually mean Toltec; that conflicts with keeping over-the-air updates (§ Constraints in [architecture.md](architecture.md)). A Toltec-free path means a real rM2 display backend while keeping the Qt editor — roughly the same weight of work as making the typewriter trustworthy (EditHelper, wrap, undo, harness). Not insurmountable; not a weekend config change either. Replacing Qt wholesale would be larger, because typing behavior would have to be rebuilt.

Launch without page buttons and Home is already partly covered: phone Show PIN / `/api/lobby`, `wd`, USB Esc. Power-button patterns or touch could fill a tablet-only gap later.

Do not start rM2 work unless there is clear community demand. Wishlist: [improvements.md](improvements.md).

## 34. Lobby tip: real browsers only, not Cursor

Touch Edit / New / Rename (and similar) show a connect tip unless a USB keyboard is present or a phone/laptop page has sent WebSocket `hello`. A real Safari/Chrome/Firefox session must count. Cursor’s embedded browser must not — it loads the same page for agent checks and was skipping the tip with no intentional keyboard. The page skips `hello` when it detects Cursor/Electron; the server also ignores `hello` from User-Agents containing `Cursor/` or `Electron/`. How: [architecture.md](architecture.md). Gotcha: [lessons.md](lessons.md).

## 35. Lobby Files paginates on e-ink

The Files list shows a fixed page of rows that fit the screen. Up/Down move the selection one row within the page; crossing the edge turns to the next or previous page. PgUp/PgDn jump a full page. Do not flick-scroll and do not slide the window one row at a time — both fight e-ink (extra redraws, muddy motion). How: [browser-vs-tablet.md](browser-vs-tablet.md), [lessons.md](lessons.md).

## 36. Lobby look and chords on disk

Lobby look, wording, language, and shortcut chords live on disk — not baked into every binary change. Edit over SSH; Writerdeck reloads without a rebuild. Main file: `/home/root/.Writerdeck/lobby-ui.json` (`language`, `visual`, optional `strings` overrides, `shortcuts`). Label packs: `lobby-ui-i18n/<lang>.json` for `en` / `no` / `es` / `de` / `fr`. Repo defaults under `config/`. `deploy-keywriter.sh` seeds missing files only so local edits survive. That JSON is the source of truth for which key each Lobby action uses — do not mirror chord tables here. Letter values mean Ctrl (or Cmd) plus that letter. Two special values need no Ctrl: `enter` (badge glyph always ↩) and `hardware_home` (tablet physical Home, not keyboard Home). Note → Lobby and Lobby → stock UI are separate lines (`global.toLobby`, `global.quit`), both defaulting to `hardware_home`. Tab shortcut lines (`tabs.files` … `tabs.about`) ship as `""`; tab *titles* are strings in the language pack. Shortcut keycap borders use `visual.badgeBorderColor`. No hardwired Ctrl-K note picker; K and Q are ordinary letters. Missing or corrupt JSON falls back to embedded defaults (or keeps the last good load after a successful start). Optional `strings` in `lobby-ui.json` override the language pack — leave that object empty unless you mean to pin a few lines; leftover English keys there will beat Norwegian (and other) packs. How: [architecture.md](architecture.md). Gotchas: [lessons.md](lessons.md). Chords: [todo-lobby-ui-shortcuts.md](todo-lobby-ui-shortcuts.md). Chrome / i18n: [todo-lobby-ui-chrome.md](todo-lobby-ui-chrome.md).

## 37. Phone-safe Lobby chords; Finder-style jump later

Phone browsers (Safari/Chrome) keep Cmd/Ctrl+R, T, W, N, L and Cmd/Ctrl+1–9 for their own chrome. Those often never reach Writerdeck from a phone keyboard, so Lobby *action* Ctrl-letters avoid R/T/W/N/L. Which letter each action uses lives only in `lobby-ui.json` (see §36). Copy/cut/paste and editor Mac chords stay Ctrl-C/X/V/A/Z.

Do not hardwire digits 1–6 for Lobby tab jump. Plain letters and digits without Ctrl stay free for a later MacOS Finder-style document jump — do not spend bare digits on Lobby tabs without deciding that trade-off.

How: [browser-vs-tablet.md](browser-vs-tablet.md). Defaults: `config/lobby-ui.json`.

## 38. One product version, auto date, no stale lies

Writerdeck is one product on the tablet even though it is two binaries. About shows a single date stamp, not separate server/editor numbers.

The stamp is an auto calendar date (`YYYY-MM-DD`). Same-day server and editor builds share that date. A second intentional ship the same day uses `scripts/product-version.sh --bump` (`.2`, `.3`, …). Do not hand-edit `VERSION` for ordinary work — CI and `deploy-rmkbd.sh` call `--write`.

Both binaries must be stamped. About combines them: if they disagree, show the older stamp and never claim “(latest)” until they match each other and GitHub `VERSION` on `main`. Editor builds write the stamp into `product_version.h` (a `-D` date string is parsed as arithmetic). How: [architecture.md](architecture.md). Gotchas: [lessons.md](lessons.md).

---

## Open risks

Over-the-air update may wipe the boot service and reset the SSH password — redeploy and re-enable. Rootfs is nearly full; everything we ship lives under `/home/root`. Do not resize rootfs. uinput is closed (§3). The editor lives in the fork (§4); residual risk is a clash when merging Singleton’s original, not a patch-script pile. Typing-test strategy is still failing (section above §1) — green is not “basic editing works.” Integrity leftovers: [integrity-audit.md](integrity-audit.md).
