# Two-way note sync with a private GitHub repo — does this exist, and is it wise?

A plain-language look at whether "type on the tablet, edit anywhere, notes sync both ways with a private GitHub repo" is a sane plan. Short answer: **yes — two-way sync is exactly what the common tools do, and it's safe as long as a clash never destroys anything.** Below is who's done it, what to copy, and the trade-offs.

**Effort:** the biggest item left — roughly 2–3× a typical Phase 9 polish task (font picker, PIN chooser). It's a new phone-side sync engine + GitHub calls + settings UI, but the tablet stays untouched, so no firmware/build risk.

## Recommendation (read this first)

Go two-way, with one golden rule: **never overwrite, never lose words.** Split it: the **tablet keeps the harmless settings** (sync on/off + repo name) so any connecting device knows the score; the **secret token lives only on each phone/laptop**. When sync is on but the device has no token, the tablet shows a gentle banner — "GitHub sync is on — log in on this device" — and notes still save locally meanwhile. Whichever device has a token is the sync engine: it reads notes, talks to the one private repo, writes changes back. On a clash, both versions are kept (`note.md` + `note (tablet copy).md`); all sorting-out happens in the browser, never on e-ink. Write in one place at a time and a slip costs a minute, never your text.

## Has anyone done this? Yes — it's common

| Tool | What it is | How close to us |
|---|---|---|
| **GitJournal** | Phone note-taking app whose only "save" is to a private GitHub/GitLab repo | Very close — proves the idea: notes-as-repo, no server, plain Markdown |
| **Obsidian Git** | Auto commit/pull/push for the popular Obsidian editor | Close — but warns mobile sync is "highly unstable"; the cause is our caution about conflicts |
| **GitSync** | Small phone app that just syncs a notes folder to git | Close — purpose-built, minimal, like an appliance |
| **Logseq / Dendron / Foam** | Note tools that keep plain Markdown you can drop in a repo | Looser — same "files in git" foundation |

So the pattern — **plain Markdown files living in a private git repo, no special server** — is proven. Our twist is just that the writer is an e-ink tablet, not a phone.

## What to copy
- Plain Markdown, one file per note, mirrored to the repo. No database, no lock-in.
- Repo name + on/off are plain settings (fine on the tablet); the token is the only secret and stays per-device.
- Sync = quietly back up while connected. Don't make the writer think about git.

## What everyone trips on (so we handle it)
- **Two devices, one note** → words get lost. Handle it: pull on open, push on save, and on a clash keep both copies — never overwrite.
- **Phone/tablet git is flaky** (Obsidian's own warning) → we sidestep it by using GitHub's plain web upload, not real git.
- **Conflicts need merge screens** → not e-ink friendly. So we don't merge on the tablet — we duplicate and let you sort it out in the companion browser.

## Who does what
- **Tablet (e-ink):** stores `.md` files plus the non-secret config (sync on/off + repo). Holds no token. If sync's on and the device lacks one, it shows "GitHub sync is on — log in on this device."
- **Companion (phone or laptop browser):** whichever has the token is the sync engine and control room — talks to GitHub, copies-on-clash, picks keep tablet / GitHub / both, renames, deletes. Decisions and the token live here.
- **Settings ⚙:** a Sync section (off by default) — toggle + repo on the tablet; token per-device; status line (last sync / clashes / renew).

## Every scenario, and how we don't get stuck

| Situation | What happens | Who resolves |
|---|---|---|
| Edit on tablet only | Pull, type, push. Done | Nobody |
| Edit in browser only | Push to repo; tablet pulls it next open | Nobody |
| Both edit same note | GitHub blocks the 2nd; phone saves `note (tablet copy).md` | Browser shows both, pick keep tablet / GitHub / both |
| Phone away / not connected | Tablet just saves locally; phone syncs next time it connects | Browser, when convenient |
| Sync on, device has no token | Tablet shows "log in" banner; notes still save locally | Browser: paste token |
| Mid-save crash | Half-note never lands (one all-or-nothing upload); old version stays | Auto |
| New note, name clash | Repo wins; second becomes `note (2).md` | Browser rename |
| Edit a note deleted elsewhere | Phone re-adds it as a copy; nothing vanishes | Browser: keep or delete |
| Rename | Treated as delete + add; old copy lingers till browser tidies | Browser |
| Token expired | Sync pauses, banner in browser, notes still save on tablet | Browser: paste new token |
| GitHub down | Tablet keeps saving; phone retries later; nothing lost | Auto |

The spine throughout: **the tablet always saves locally; the phone syncs best-effort; clashes duplicate; the browser cleans up.** No screen ever waits on a merge, and no secret ever sits on the tablet.

## SWOT analysis

**Strengths**
- Proven idea; off-device backup; notes readable/editable on any computer; one fenced repo; token never on the open-LAN tablet; fits the existing daemon (phone makes the web calls, tablet unchanged).

**Weaknesses**
- Syncs only while the phone is connected; clashes leave two copies to tidy; token sits in the phone's browser storage; needs periodic token renewal.

**Opportunities**
- Free version history; pairs with the planned encryption; edit on any computer; even a public blog from the same notes.

**Threats**
- Overlapping edits leave stray copies (annoying, not data loss); a stolen token can wipe that repo; tokens expire (manual renew); GitHub outages.

## Bottom line
Common, sensible, low-risk **if** kept to one repo and clashes duplicate instead of overwrite. Two-way is fine — just write in one place at a time.

---

## Implementation (shipped — Sonnet built, Opus reviewed)

Built in [daemon/main.go](../daemon/main.go) + [daemon/index.html](../daemon/index.html). The token never touches the tablet; the sync engine is JS in the capture page. The Go side gained only a non-secret flag + repo name, reusing the existing `settingsData`/`/api/settings` pattern.

**Phase A — tablet config (Go).** ✅ `settingsData` gains `syncOn`/`syncRepo`; `settingsHandler` GET returns them, POST validates via `isValidGitHubRepo` (`owner/repo`, 400 otherwise). No token field — ever. (Note: the Lobby "log in" line is not wired; the phone shows the banner instead — see review.)

**Phase B — token + settings UI (phone).** ✅ Sync section in the ⚙ overlay: toggle + repo POST to `/api/settings`; token save/clear in `localStorage`; last-synced status line.

**Phase C — sync engine (phone JS).** ✅ GitHub Contents API per note: `GET …/contents/{name}` → `sha`; `PUT` with that `sha` updates, without creates; stale `sha` → 409/422 = clash. Open: `pullNoteAndUpdate` writes to the tablet via the new `PUT /api/notes/{name}` upsert. Exit typing: `pushNote` reads `/api/notes/{name}` and pushes. Per-note `sha` in `localStorage`.

**Phase D — copy-on-clash.** ✅ On 409/422, `handleClash` saves `note (tablet copy).md` and pulls GitHub's version into `note.md`; 30 s clash banner. Tablet never blocks.

**Phase E — robustness.** ✅ No token → yellow banner, local saves continue; 401/403 → "renew token" banner; offline/5xx swallowed. Deletes & renames do **not** propagate (see review).

Acceptance: one repo, no token on tablet, sync off = byte-identical to today. ✅

## Post-build review (Opus) — all fixed

Three gaps were found and closed; none ever risked tablet-local data.

1. **Push fired before the tablet saved (correctness).** ✅ Fixed. Push no longer runs on phone-back (`hideTypingView`); it runs on the tablet's post-save `exitedit`, tracked via a new `tabletOpenNote` that survives phone-back. Note-switch pushes the just-saved previous note. No more stale-push / poisoned-`sha` divergence.
2. **Delete & rename didn't propagate.** ✅ Fixed. New `ghDelete` (Contents API DELETE + stored `sha`) fires on delete so a note can't resurrect on next pull; rename deletes the old GitHub path then pushes the new name.
3. **Pull-then-open wasted for the already-open note.** ✅ Fixed. `openNote` skips the pre-open pull when `filename === tabletOpenNote`, so keywriter's save-on-load can't clobber it.

Every scenario-table row now holds. Ship behind the off-by-default toggle. Remaining nice-to-haves (not blocking): periodic background pull while idle, and a real diff view for clashes instead of duplicate-and-tidy.

**Caveat — the delete/rename fix is browser-scoped (by design).** `ghDelete` and rename-then-repush only fire when the *phone browser* performs the op, because the reconciler is deliberately non-authoritative: it unions the tablet + repo note lists and copies whatever's missing either way, never deleting on its own (a reconciler that can't delete can't lose a note — the whole reason we avoid real git on the low-RAM tablet). So a delete or rename made *outside* the browser — VS Code, `git`, the GitHub web UI — reads only as "a note is missing from one side" and resurrects (or, for rename, duplicates) on the next sync. Rule: do destructive ops in the phone browser. Durable record: [decisions.md](decisions.md) #19; cheap partial fix (marker-aware delete) tracked in [../TODO.md](../TODO.md).
