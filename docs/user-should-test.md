# Things you should test by hand

Scripts cover a lot. This list is only what still needs your fingers on the tablet.

## Due now — Physical Home

While Writerdeck is open, the server owns the middle button ([decisions.md](decisions.md) §16). Scripts cannot press it.

Start Writerdeck, then:

- Edit → Home once. You should land in Lobby Files with your typing saved — not quit to the stock UI.
- Read → Home once. Lobby again — not quit.
- Lobby → Home once. Writerdeck quits; stock UI returns. The phone can still reach port 8000.
- From stock UI, left + right page buttons together should launch to Lobby.
- Power while editing: sleep, wake, keep writing.

With a USB keyboard: Home from edit goes to Lobby; Home from Lobby quits.

When these pass, check them off or tell the agent. Detail: [todo-handoff-physical-home-input.md](todo-handoff-physical-home-input.md).

## Spot-checks after certain changes

After USB keymap work: Norwegian characters and the Lobby layout picker.

After phone UI deploy: hard-refresh the page — keyboard shell with the gray logo (no notes list); the bar shows Connected or Tablet offline, not stuck on connecting. After Exit (stock UI) Launch Writerdeck appears under the logo; on the Files tab that button is hidden. Paste and Sync setup still open.

After Lobby Download: with the phone page open, Download (or **g**) on Files should prompt “Download here?” on the phone; without a phone page, the tablet should say to open one first.

After Files pagination chrome: with more notes than fit one screen, Prev / Page N/M / Next and a line should sit above New/Edit; with few notes, that strip stays hidden.

After New / Rename dialogs: try a name that already exists — the dialog should stay open with a short message under the name, not only a box above the list. Also try the same letters in different case (`Doc` vs `doc`), and a plain name that already exists as encrypted (or the reverse).

After editing `/home/root/.Writerdeck/lobby-ui.json` over SSH: a color, help line, or Ctrl-letter (for example change Read from `v` to `g`) should update on Lobby within a couple of seconds without redeploying the binary. Bad JSON should leave the previous look in place. Defaults live in repo `config/lobby-ui.json`. Avoid browser-reserved letters (R T W N L) and do not bind Lobby chrome to bare letters or digits ([decisions.md](decisions.md) §37).

After sync or vault work: Sync setup still runs; a wrong vault PIN shows an error, not a blank editor. With change-driven sync: idle tablet should stay quiet in the journal (one `sync: nothing to do` at most per quiet streak, then a ×N summary when something finally changes). Edit a note, Home — expect `sync: pushed … (home)`.

You do not need to re-type the automated typing checks by hand while they stay passing. If real writing still feels wrong, note the keys and tell the agent.

## How to launch

From stock UI: USB Esc, both page buttons, phone **Show PIN on tablet**, tablet `~/wd`, or Mac `wd`. Cold boot stays on the stock UI.
