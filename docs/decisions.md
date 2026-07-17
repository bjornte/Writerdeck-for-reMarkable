# Architecture decisions

Why the project is built this way — document integrity and core editing first, device quirks later. How: [architecture.md](architecture.md). Open work: [../TODO.md](../TODO.md). Finished: [../DONE.md](../DONE.md). Gotchas: [lessons.md](lessons.md). Words: [terms.md](terms.md).

Numbers were renumbered 17 Jul 2026; prefer titles over old § IDs.

---

## Document integrity

Writerdeck is a typewriter. Your prose must survive as plain Markdown on disk. That rule comes before every feature. Full audit: [integrity-audit.md](integrity-audit.md).

Notes are UTF-8 Markdown, never Qt HTML. An open note must not be silently overwritten by sync or remote delete. Saves use defined paths, autosave about every 45 seconds, and save-before-stop on deploy; a hard kill can still lose recent typing. If the file on disk changes under you, reload or show a conflict — do not let a stale buffer win. GitHub sync backs up; it must not empty-push or delete against a live edit.

No change to the daemon, sync, build script, or note APIs ships without an integrity pass.

## Device verification

A successful deploy script is not enough. Rebuild when the QML changed, deploy, relaunch, read the tablet journal. Fail on QML parse errors or an editor that exits at once. After QML changes run the edit-session check. After caret or selection work run the automated typing tests (§13). After Lobby or Home run the Lobby keyboard test (§15).

---

## 1. Plain-text edit mode

While you type, the editor stays plain text — raw Markdown on screen and in memory. Esc preview may look fancy. Headings and bold belong there, not in edit mode.

Rich text in edit was tried. Pulling formatted text back into the note produced empty files, HTML on disk, and broken previews. Real hide-the-markers editing stays out of scope unless someone writes a new decision and re-proves integrity.

## 2. Display sync after Lobby clears the buffer

Going Home clears the on-screen box. Reload and mode toggle must push the real note text (or its preview HTML) back onto the screen — never read fancy HTML back into the note. Without that, Home can save a zero-length file.

## 3. Socket input, not uinput

This kernel cannot load a uinput (fake keyboard device). Do not retry that path. Keystrokes arrive on a local socket and become Qt key events. The phone resolves the keymap; the tablet gets characters and Mac/Linux-style Ctrl/Alt shortcuts. USB Linux keyboards already use those shortcuts — no second shortcut set.

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

The server keeps the phone reachable under the stock UI. Only the editor session flips xochitl. Keep-awake covers the editor child, not the whole device. Units that start from `/home/root` must wait for that mount or cold boot fails with `203/EXEC`.

## 10. Companion model

No phone app — the tablet is the server. Day-to-day files and settings live in the Lobby; the phone is for typing, upload/download, paste-at-cursor, and the sync token — [browser-vs-tablet.md](browser-vs-tablet.md). A PIN appears on the e-ink each boot. Home twice: edit to Lobby, Lobby to quit.

## 11. GitHub sync copies missing notes both ways

Sync runs on the tablet. The token sits in the browser and in tablet RAM — never on disk. After restart the phone can repost it. Sync copies missing notes both ways; it does not delete on its own. Empty push over a known-good note is refused. Details: [server-sync-implementation.md](server-sync-implementation.md).

## 12. Optional at-rest encryption (private notes)

Pairing PIN and vault PIN are separate. Encrypted notes are `.md.enc` beside plain `.md`. Each file is sealed; the vault PIN unlocks a key held in RAM only while that note (or a one-shot encrypt/decrypt) is active. PIN every open, including note switches. Tablet-only entry. Per-note encrypt/decrypt — no bulk lock. Sync copies encrypted files without reading them and mirrors recovery material under `secret/`. Markdown integrity applies to `.md` only. See [integrity-audit.md](integrity-audit.md).

## 13. Automated typing tests

We prove caret and selection on the real tablet, over the same path the phone uses — not by reading saved files. About 110 checks; 38 are the “basic editing works” set. Pass/fail log: [editor-testing/](editor-testing/). Test notes use the `z-test-` prefix (§32). USB layout quirks still need a human check after qmap changes.

## 14. Edit-session check

Opening a note from outside must keep the editor up for several seconds. Instant exit usually means a broken QML, not a server bug. Run after Writerdeck or QML deploy.

## 15. Lobby keyboard check

After Home from edit, Lobby keys must still work. The Lobby keyboard script checks that path. Run it with edit-session after Lobby or Home QML changes.

## 16. Physical Home: exclusive gpio grab

While Writerdeck is open, the server exclusively grabs the tablet’s Home/Power/page buttons so Qt never sees a second Home. Release on exit so the stock UI works again. USB Home is unchanged. Idle page-button launch (§23) needs no grab. Handoff: [todo-handoff-physical-home-input.md](todo-handoff-physical-home-input.md).

## 17. PIN and per-IP lockout

Six digits, four, or none. None means anyone on the Wi-Fi can hit the notes API — the UI warns. Five wrong guesses lock that IP for a minute, not the whole device. Store length as `"6"`, `"4"`, or `"none"`.

## 18. Show PIN on tablet

A second phone mid-edit sees the note, not the PIN. Ask the tablet to return to Lobby so the PIN is readable on e-ink. That call is pre-auth and rate-limited; it reveals nothing over the wire.

## 19. Tablet file create/rename/delete via trusted socket

Lobby file ops use the same local socket as keystrokes. The server does the disk work and can nudge the phone. Launch Lobby with `wd` on the Mac or `~/wd` on the tablet.

## 20. Bluetooth remote key capture on the phone

Bluetooth pairs to the phone. Capture only when the tablet asks (edit, read, or a Lobby text field). Leave Browse alone so normal browser shortcuts still work.

## 21. Lobby Files: Edit and Read

Two opens: Edit to type, Read to preview. The phone learns which mode so a Bluetooth keyboard can still Esc into edit. Enter edits; `v` reads.

## 22. Lobby tab order (Files first)

Boot and Home land on the note list, not the welcome screen. Coming back from edit reselects the note you left.

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

## 30. Upload reuses the safe create route

Phone upload goes through the same create API as a new note — path checks and “already exists” included. The server enforces size; the client check is courtesy.

## 31. Display rotation persists in settings

Rotation is saved on the tablet and pushed when the editor connects. Change it from Lobby Settings or Ctrl-R / Ctrl-arrows. Phone Preferences no longer rotate.

## 32. Device test note names

Automated tests use filenames starting with `z-test-` so they sort last and stay obvious in the Files list.

---

## Open risks

Over-the-air update may wipe the boot service and reset the SSH password — redeploy and re-enable. Rootfs is nearly full; everything we ship lives under `/home/root`. Do not resize rootfs. uinput is closed (§3). The editor lives in the fork (§4); residual risk is a clash when merging Singleton’s original, not a patch-script pile. Calling typing work done still means all 110 automated typing checks (§13). Integrity leftovers: [integrity-audit.md](integrity-audit.md).
