# Backing up notes to GitHub — does this idea exist, and is it wise?

A plain-language look at whether "type on the tablet, edit anywhere, notes sync to a private GitHub repo" is a sane plan. Short answer: **yes — two-way sync is exactly what the common tools do, and it's safe as long as a clash never destroys anything.** Below is who's done it, what to copy, and the trade-offs.

## Recommendation (read this first)

Go two-way, with one golden rule: **never overwrite, never lose words.** Because the Writerdeck only does anything while the companion phone is connected, **the phone is the sync engine and the only place the token lives** — the tablet stores no secret at all. The phone reads notes off the tablet, talks to one private repo (fine-grained token, that repo only) straight from the browser, and writes pulled changes back. If GitHub says "someone changed this first," it keeps both versions side by side (`note.md` and `note (tablet copy).md`) instead of clobbering. The tablet stays dumb; all "sort it out" work happens in the browser, never on e-ink. Write in one place at a time and a slip costs a minute, never your text.

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
- The token only sees one repo — and lives in the phone, never on the tablet.
- Sync = quietly back up while connected. Don't make the writer think about git.

## What everyone trips on (so we handle it)
- **Two devices, one note** → words get lost. Handle it: pull on open, push on save, and on a clash keep both copies — never overwrite.
- **Phone/tablet git is flaky** (Obsidian's own warning) → we sidestep it by using GitHub's plain web upload, not real git.
- **Conflicts need merge screens** → not e-ink friendly. So we don't merge on the tablet — we duplicate and let you sort it out in the companion browser.

## Who does what
- **Tablet (e-ink):** just stores `.md` files and serves them to the phone. Holds no token, makes no internet calls. Stays dumb.
- **Companion (phone browser):** the sync engine and control room. Keeps the token, talks to GitHub, copies-on-clash, and is where you pick "keep tablet / keep GitHub / keep both," rename, delete, and confirm. All decisions and all secrets live here.

## Every scenario, and how we don't get stuck

| Situation | What happens | Who resolves |
|---|---|---|
| Edit on tablet only | Pull, type, push. Done | Nobody |
| Edit in browser only | Push to repo; tablet pulls it next open | Nobody |
| Both edit same note | GitHub blocks the 2nd; phone saves `note (tablet copy).md` | Browser shows both, pick keep tablet / GitHub / both |
| Phone away / not connected | Tablet just saves locally; phone syncs next time it connects | Browser, when convenient |
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
- Needs the phone connected to sync; clashes leave two copies to tidy; token sits in the phone's browser storage.

**Opportunities**
- Free version history; pairs with the planned encryption; edit on any computer; even a public blog from the same notes.

**Threats**
- Overlapping edits leave stray copies (annoying, not data loss); a stolen token can wipe that repo; tokens expire (manual renew); GitHub outages.

## Bottom line
Common, sensible, low-risk **if** kept to one repo and clashes duplicate instead of overwrite. Two-way is fine — just write in one place at a time.
