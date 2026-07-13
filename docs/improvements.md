# Improvements

Wish-list and design notes — not tracked work. Shipped: [../DONE.md](../DONE.md). Actionable verify items: [../TODO.md](../TODO.md). Capability split: [browser-vs-tablet.md](browser-vs-tablet.md).

## USB keyboard locales

Browser to WebSocket resolves layout in the phone OS — Norwegian works today. USB uses Qt evdev with shipped `us` and `no` qmaps and a Lobby Keyboard tab picker; Norwegian national characters and AltGr verified on hardware ([DONE.md](../DONE.md)).

Alt+Left/Right on standard Linux kmaps become fake Escape in Qt — see [lessons.md](lessons.md). `writerdeck-alt-arrows.inc` keeps them as arrows in `no` and `us` qmaps.

Qt ignores `loadkeys` and `setxkbmap`. Use `QT_QPA_EVDEV_KEYBOARD_PARAMETERS` and a qmap file — [remarkable-keywriter#1](https://github.com/dps/remarkable-keywriter/issues/1).

Regenerate with `bash keymaps/generate.sh`. Phone-path selection regressions use `scripts/test-keyboard-harness.sh`; that path does not exercise evdev — re-check qmaps on hardware after qmap edits.

## Edit view (future)

More VS Code-like list indentation; optional cursor block beside the active line; headline navigation; status bar with title, confirmations, zoom, time — battery is already on the phone status bar.

WYSIWYG Markdown in edit mode (large headings, bold, italic without visible `**` / `#`) is not planned — see [decisions.md](decisions.md) §26. Esc/read mode already renders sundown output.

## Phone Markdown mirror (future)

While editing on e-ink, show the open note rendered as Markdown on the phone — read-only context, not a second editor. Tapping a word in that view could move the tablet cursor there and scroll that section into view on e-ink. Depends on cursor position over the socket (already published for the keyboard harness) and a rendered Markdown pane on the phone.

## Browser (future)

Bulk select and multi-delete; search across titles and bodies; HTTPS for native Share sheet; phone-side vault PIN UI (tablet-only PIN entry is shipped — [decisions.md](decisions.md) §31).
