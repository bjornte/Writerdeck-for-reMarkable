# Things you should test by hand

Automated scripts cover a lot. This list is only what still needs **your fingers on the tablet** (and sometimes a USB keyboard or phone). Check items off when done; add new ones when a change ships that scripts cannot cover.

Agents: after a deploy that needs owner fingers, add a short unchecked section here and mention it in chat.

---

## Due now — Physical Home (July 2026)

Server now owns the middle button while Writerdeck is open ([decisions.md](decisions.md) §28). Scripts cannot press the real button.

Start Writerdeck (boot, Esc, or L+R page buttons). Then:

- [ ] **Edit → Home** — open a note, type a few characters, press the middle button once. You should land in the Lobby (Files). The note should still be there with your typing. Writerdeck must **not** quit straight to the stock UI.
- [ ] **Read → Home** — open a note, Esc into read/preview, press Home once. Lobby again — not quit, not a crash.
- [ ] **Lobby → Home** — from the Lobby, press Home once. Writerdeck quits; stock reMarkable UI returns. Phone can still reach `http://<tablet-ip>:8000/`.
- [ ] **Lobby Home once only** — from Lobby, one firm press. Should not feel like a double-fire (instant quit-then-weirdness). A second press after you are already in stock UI is a different world (xochitl).
- [ ] **Idle page buttons still work** — quit to stock UI, open a notebook if you like, flip pages with the side buttons. Then hold **left + right** together: Writerdeck should launch to the Lobby.
- [ ] **Power while editing** — in a note, press power: sleep screen, then wake with power again and confirm you can keep writing (or return sensibly).

If you have a USB keyboard plugged in:

- [ ] **USB Home from edit** — should go to Lobby (same as middle button), not scramble the line.
- [ ] **USB Home from Lobby** — should quit Writerdeck.

When all of the above pass, tell the agent (or check them off and commit). Detail: [todo-handoff-physical-home-input.md](todo-handoff-physical-home-input.md).

---

## Recurring spot-checks (not every session)

Do these after the kind of change named — not after every tiny docs edit.

### After USB keymap / launcher changes

- [ ] Norwegian USB: æ ø å, AltGr, `@`, `{` `}`
- [ ] Layout picker in Lobby → Keyboard still switches `us` / `no`

### After phone UI (`daemon/*.js`) deploy

- [ ] Hard-refresh `http://<tablet-ip>:8000/` — notes list loads, bar shows Connected (or Tablet offline), not stuck on `connecting...`

### After sync / vault work

- [ ] Sync setup on the phone still saves and runs Sync
- [ ] Encrypted note: open with vault PIN on the tablet; wrong PIN shows an error, not a blank editor

### Keyboard editing (product sign-off)

Full harness is green (**107/107**), including mid-sentence Shift+Up/Down in long wrapping paragraphs. You do **not** need to re-type harness cases by hand. If something still feels wrong while writing for real, note the keys you pressed and tell the agent.

---

## How to launch Writerdeck for testing

- From stock UI: USB **Esc**, or **left + right** page buttons together, or tablet SSH `~/wd`, or Mac `wd` / `bash scripts/lobby.sh`
- Phone: open the tablet IP on port 8000 after Writerdeck is up
