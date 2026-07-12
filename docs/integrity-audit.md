# Document integrity — audit

**Last updated:** 2026-07-12

## Summary (plain language)

Writerdeck is a typewriter. **Document integrity** means your words end up as real Markdown files on the tablet, stay readable, and are not silently replaced, emptied, or forked without you noticing.

**Where we stand today:** For normal solo use — one person, one tablet, typing one note at a time, saving with Home or waiting for autosave — your text is in good shape. Saves go to disk as plain text; GitHub sync will not delete notes on its own; if sync disagrees with you, it tends to keep both copies rather than pick a loser.

**What is still not “solved”:**

- Sudden death — pull the battery or kill the app hard, and you can lose up to ~45 seconds of typing since the last autosave. Home and normal exit save; a crash does not.
- While you are editing — the note on screen is protected from sync overwriting it. Other notes can still sync. If GitHub (or sync) changes a file while you have it open, you get a warning on the phone; the tablet does not auto-refresh. If you keep typing and save, your buffer can win over the newer file on disk — you have to choose reload or keep going.
- Sync mess, not sync loss — clashes and outside edits often leave two files (`note.md` and `note (tablet copy).md`) or resurrect old versions. You rarely lose words entirely; you may need to merge or delete duplicates by hand.
- Fast switching — jumping between notes very quickly (Ctrl-K, Lobby Files) can occasionally interleave saves and loads in the wrong order. Uncommon in normal writing; possible under stress.
- Trust the Wi‑Fi — with PIN turned off, anyone on your network can read or change notes. That is a privacy and tampering risk, not just a sync issue.
- Unknown unknowns — firmware updates, rare Qt bugs, and paths we have not hammered in testing could still surprise us.

Bottom line: Integrity here means “don't lose my words silently.” On that bar, we are strong for everyday writing, weaker under crash, clash, and multi-place editing, and honest that sync tidiness is your job. We are not at “bank-grade” or “real-time collaborative editing” integrity — and we are not trying to be.

Technical detail, shipped mitigations, and open items follow. Contract: [architecture.md](architecture.md) § Document integrity, [decisions.md](decisions.md) § Document integrity.

---



## What document integrity means (for this product)

This is not abstract security auditing. For Writerdeck it is a small set of promises to the person writing:


| Promise                                          | In practice                                                                                                                                                                             |
| ------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Your files are really Markdown**               | What lands on disk is UTF-8 text you can open anywhere — not hidden editor formats or HTML masquerading as notes.                                                                       |
| **Saving means saved**                           | When you leave a note normally (Home, switch note, exit), or when autosave fires, text is written through defined paths to disk — not left only in memory.                              |
| **The open note is sacred**                      | While you are editing a file on the tablet, sync and remote changes must not silently replace that file or push an empty version over good content.                                     |
| **If disk and screen disagree, you should know** | When something else changes a file under you, the system should warn you — not pretend nothing happened.                                                                                |
| **Sync helps, it does not rule**                 | GitHub backup is optional. It may copy, clash, or duplicate; it must not be the reason words vanish. Deletes on the repo side only propagate when the tablet copy is clearly untouched. |
| **Writes should not tear mid-crash**             | A save should either leave the old file intact or the full new file — not a half-written fragment.                                                                                      |


What it **does not** mean:

- Zero data loss under every physics failure (power cut mid-sentence).
- Automatic merge of two people editing the same paragraph.
- A single canonical version when you edited the same note on GitHub and the tablet without looking.
- Protection from a malicious guest on your Wi‑Fi if you disabled the PIN.

The **feature gate** for developers: before shipping anything that touches notes, saves, opens, sync, or editor lifecycle — *can this lose text, write wrong bytes, or overwrite without the user knowing?* If yes, fix it, or the owner explicitly accepts the risk and it is logged here.

---



## Critical assessment (2026-07-12)



### Strengths — believe these

1. **Plain-text contract is enforced.** Saves reject rich-text garbage; loads sanitize; the on-disk format matches the product story.
2. **Atomic writes on the server path.** Temp file + rename reduces “corrupt half file” risk on power loss during a save.
3. **Edit lease is real and tested.** Reconcile skips the note the tablet is editing; other notes still sync (fixed 2026-07-12). Empty-push and Home-wipe regressions are guarded after real incidents (#19, #24).
4. **Sync engine on the tablet.** Backup no longer depends on the phone staying open — a major integrity win for “will my notes reach GitHub?”
5. **Non-destructive reconciler.** Sync union-copies; it does not mass-delete the tablet when GitHub hiccups. Marker-aware delete is conservative (404 per note before acting).
6. **Operational save-before-deploy.** Deploy and graceful shutdown try to flush the open buffer first — appropriate for a device you SSH into.



### Weaknesses — take seriously

1. **Coherence is incomplete (biggest gap).** We detect disk drift and offer a manual reload on the phone. We do **not** block save when disk changed; we do **not** auto-reload the tablet after a clash pull. A attentive user is still part of the safety system. That is a integrity **process** burden, not full **system** guarantee.
2. **The 45-second autosave window.** Acceptable for a typewriter metaphor, but dishonest to call “durable” without caveat. SIGKILL, segfault, or forced kill between ticks = lost words. No CRDT, no continuous journal.
3. **Clash UX creates debt.** `(tablet copy)` preserves text but violates mental model of “one note, one name.” Users must reconcile duplicates manually; stale `syncMeta` or external edits can still resurrect or fork.
4. **Async note switching.** `doLoad` and save-ack timeouts under rapid open/switch are a known race class. Mitigated by typical pacing; not proven absent under fuzzing.
5. **Two binaries must match.** Server expects loopback PUT and `autosavenow` from a current Writerdeck. Deploy server-only or stale editor = save paths silently degraded.
6. **Testing depth.** `test-edit-session.sh` guards “editor stays up on Edit” — not “save under load,” “clash while typing,” or “power sleep with open note.” Integrity confidence is partly **design review + spot checks**, not exhaustive automation.
7. **PIN** `none` **on LAN.** Integrity includes “only I can change my notes.” Open LAN mode fails that for confidentiality and tampering.



### Severity matrix


| Scenario                                   | Likely outcome today                                       | Severity              |
| ------------------------------------------ | ---------------------------------------------------------- | --------------------- |
| Normal typing, Home to Lobby, sync on      | Saved; sync runs for other notes                           | OK                    |
| Token restore, Save & verify while editing | Token OK; other notes sync; open note skipped              | OK                    |
| Power loss mid-paragraph                   | Up to ~45 s lost                                           | Medium                |
| GitHub edited same note while tablet open  | Drift banner; manual reload; save may overwrite if ignored | Medium–high           |
| Sync clash                                 | Both copies kept; manual cleanup                           | Low loss, medium mess |
| Rapid Ctrl-K / Files switching             | Possible wrong interleave (uncommon)                       | Medium (rare)         |
| Guest on LAN, PIN off                      | Can read/change/delete notes                               | High (ops choice)     |
| Firmware OTA                               | Unknown; binaries may need redeploy                        | Unknown               |




### Verdict

**Shipped integrity (slices 1–11 + server sync + 2026-07 fixes) is adequate for a personal typewriter with optional GitHub backup** — provided the owner understands autosave timing, reads clash/drift warnings, and does not expect Google Docs–grade merge.

**Not closed.** Residual risks in the next section remain open until eliminated or explicitly accepted. Calling integrity “done” would overstate what we have built.

---



## Fixed (slices 1–11 + follow-ups)


| Slice | What                                                                                                 |
| ----- | ---------------------------------------------------------------------------------------------------- |
| 1     | Edit lease — `notifyOpen` + `openedit` WS; reconcile skips open note                                 |
| 2     | Content fidelity — plain-markdown save contract, load sanitizer, `toggleMode` fix, server HTML guard |
| 3     | `notedeleted` / `noterenamed` — editor notified on phone rename/delete of open file                  |
| 4     | Reconcile policy — `openNote` in `/api/status`; `reconcileAll` skips **only** the open note          |
| 5     | OCC — GET `ETag`; PUT overwrite requires `If-Match`; sync sends revision                             |
| 6     | Atomic server writes — `writeNoteFile` temp+rename                                                   |
| 7     | Tablet CRUD → GitHub — server-side sync (was `pendingSync` + phone)                                  |
| 8     | Disk↔buffer drift — `diskchanged` WS, phone drift banner, `POST /api/reload`                         |
| 9     | Autosave — 45 s `autosaveTimer` while editing                                                        |
| 10    | Tablet atomic saves — loopback `PUT /api/notes` → `writeNoteFile`                                    |
| 11    | Save before deploy/stop — `POST /api/flush-save`, deploy graceful wait, SIGTERM flush                |


Also under contract: empty-push guard (#19), Lobby Home wipe fix (#24), server-side sync engine (2026-07-11/12). Shipped summary: [../DONE.md](../DONE.md) § Document integrity.

## Known open (residual risks)

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
- `localStorage` ~~loss~~ **Closed (2026-07-11):** `syncMeta` on tablet; engine server-side.
- ~~Multiple browser tabs~~ **Closed:** phone no longer runs reconcile.
- ~~Power sleep waited on phone~~ **Closed:** server reconcile before suspend.
- ~~Phone required for sync~~ **Closed:** server engine + tablet **Sync now**.
- ~~Save & verify stall while editing~~ **Closed (2026-07-12):** partial reconcile + UI fix.

**Auth / ops**

- PIN `none` — anyone on LAN can mutate notes (integrity + confidentiality on untrusted Wi-Fi).
- `restore-wiped-notes.sh` — Git history only; no partial-truncation or duplicate cleanup beyond `(tablet copy)` names.



## Unknown (unbounded)

No claim that the threat surface is complete.

- Bugs in paths not stress-tested (rapid switching, multi-device, sleep/wake edge cases, first boot with sync on).
- Firmware OTA side effects on deployed binaries or unit file.
- Qt / QML regressions not caught by `test-edit-session.sh` (save paths, autosave, loopback PUT).
- Novel failure modes (network partition mid-save, partial HTTP write despite atomic rename, future feature regressions).

**Process:** before shipping note/save/sync/lifecycle changes, ask: *can this lose text, write wrong bytes, or overwrite without the user knowing?* If yes, it does not ship until mitigated or **explicitly accepted** by the owner and logged here.

## Reference patterns


| Pattern                  | Status                    | Pointer                                                                     |
| ------------------------ | ------------------------- | --------------------------------------------------------------------------- |
| Edit lease               | Shipped                   | [Ink & Switch local-first](https://www.inkandswitch.com/essay/local-first/) |
| Plain-text save contract | Shipped                   | —                                                                           |
| Optimistic concurrency   | Shipped (HTTP)            | [RFC 7232](https://httpwg.org/specs/rfc7232.html)                           |
| Conflict copies          | Partial (`(tablet copy)`) | [Syncthing conflicts](https://docs.syncthing.net/users/syncing.html)        |
| Atomic durable write     | Shipped                   | [google/renameio](https://github.com/google/renameio)                       |
| CRDT / OT                | Not planned               | [Automerge](https://automerge.org/)                                         |


