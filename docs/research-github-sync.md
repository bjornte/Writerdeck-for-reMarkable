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

## Implementation plan (Opus plans, Sonnet builds)

Concrete enough to hand to a build session. The token never touches the tablet; the sync engine is JS in the capture page. The Go side gains only a non-secret flag + repo name, reusing the existing `settingsData`/`/api/settings` pattern in [daemon/main.go](../daemon/main.go).

**Phase A — tablet config (small, Go).** Extend `settingsData` with `syncOn bool` + `syncRepo string`; surface in `settingsHandler` GET/POST with `notesSafe`-style validation (`owner/repo`, 400 otherwise). No token field — ever. Lobby gains one conditional line: sync on + no token on this device → "GitHub sync is on — log in on this device." Mirrors decisions #16–18.

**Phase B — token + settings UI (phone).** Add a Sync section to the ⚙ overlay: toggle, repo, token field, status line (last sync / clashes / renew). Token persists to `localStorage` only; on/off + repo POST to `/api/settings`.

**Phase C — sync engine (phone JS).** GitHub Contents API per note: `GET …/contents/{name}` → `sha`; `PUT` with that `sha` updates, without creates; a stale `sha` returns 409 = clash. Open: pull, write to tablet via `POST /api/notes`. Save: read `/api/notes/{name}`, push. Track each note's last-synced `sha` in `localStorage`.

**Phase D — copy-on-clash.** On 409, save GitHub's copy as `note (tablet copy).md` (existing `notesSafe` `(2)` rule), surface "keep tablet / GitHub / both" in browser. Tablet never blocks.

**Phase E — robustness.** No token / 401 → banner, local saves continue; offline/5xx → retry; deletes & renames log not auto-propagate; verify each on-device per [TODO.md](../TODO.md) cadence.

Acceptance: one repo, two devices, every scenario-table row holds; no token on tablet; sync off = byte-identical to today.
