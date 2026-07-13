# Improvements

Wish-list and design notes — not tracked work. Shipped: [../DONE.md](../DONE.md). Actionable: [../TODO.md](../TODO.md). Capability split: [browser-vs-tablet.md](browser-vs-tablet.md).

## Encrypted note subset

Today one global PIN gates the whole notes API. Shipped: optional per-note encryption with a separate 6-digit vault PIN (tablet only), `.md.enc` files, and GitHub `secret/` recovery — see [decisions.md](decisions.md) §31.

## USB keyboard locales

Browser to WebSocket resolves layout in the phone OS — Norwegian works today. USB goes through Qt evdev with a `.qmap` file; national layouts need shipped qmaps and the Lobby picker, which is done, pending device test for AltGr and æøå.

Alt+Left/Right on standard Linux kmaps become fake Escape in Qt — see [lessons.md](lessons.md). `writerdeck-alt-arrows.inc` keeps them as arrows in `no` and `us` qmaps.

Qt ignores `loadkeys` and `setxkbmap`. Use `QT_QPA_EVDEV_KEYBOARD_PARAMETERS` and a qmap file — [remarkable-keywriter#1](https://github.com/dps/remarkable-keywriter/issues/1).

Regenerate with `bash keymaps/generate.sh`. Deploy via `deploy-keywriter.sh`. Lobby Keyboard tab picks layout; applies on next editor launch. `grab=1` dedicates the keyboard to Writerdeck during a session; the event node varies by device — test AltGr and æøå on hardware. Phone-path selection regressions use `scripts/test-keyboard-harness.sh`; that path does not exercise evdev.

## Edit view (future)

More VS Code-like list indentation; optional cursor block beside the active line; headline navigation; status bar with title, confirmations, zoom, time — battery is already on the phone status bar.

WYSIWYG Markdown in edit mode (large headings, bold, italic without visible `**` / `#`) is not planned — see [decisions.md](decisions.md) §26. Esc/read mode already renders sundown output.

Scroll/cursor niceties in `build-keywriter.sh` (visual last-line detection via `positionToRectangle`, edge-threshold `ensureVisible`) are patched but need CI rebuild and device verify after deploy.

## Phone Markdown mirror (future)

While editing on e-ink, show the open note rendered as Markdown on the phone — read-only context, not a second editor. Tapping a word in that view could move the tablet cursor there and scroll that section into view on e-ink. Depends on cursor position over the socket (already published for the keyboard harness) and a rendered Markdown pane on the phone.

## Browser (future)

Bulk select and multi-delete; search across titles and bodies; HTTPS for native Share sheet; phone-side vault unlock UI (tablet unlock flow is shipped — [decisions.md](decisions.md) §31).
