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

After phone UI deploy: hard-refresh the page — notes list loads; the bar shows Connected or Tablet offline, not stuck on connecting.

After sync or vault work: Sync setup still runs; a wrong vault PIN shows an error, not a blank editor.

You do not need to re-type the automated typing checks by hand while they stay passing. If real writing still feels wrong, note the keys and tell the agent.

## How to launch

From stock UI: USB Esc, both page buttons, tablet `~/wd`, or Mac `wd`.
