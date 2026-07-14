# Improvements

Wish-list and design notes — not tracked work. Shipped: [../DONE.md](../DONE.md). Actionable verify items: [../TODO.md](../TODO.md). Capability split: [browser-vs-tablet.md](browser-vs-tablet.md).

## USB keyboard locales (future)

`us` and `no` qmaps ship today with a Lobby Keyboard tab picker ([DONE.md](../DONE.md)). Additional layouts: add kmap sources under `keymaps/src/`, run `bash keymaps/generate.sh`, verify on hardware (phone path does not exercise evdev — see [lessons.md](lessons.md) § Keyboard and selection).

## Edit view (future)

More VS Code-like list indentation; optional cursor block beside the active line; headline navigation; status bar with title, confirmations, zoom, time — battery is already on the phone status bar.

WYSIWYG Markdown in edit mode (large headings, bold, italic without visible `**` / `#`) is not planned — see [decisions.md](decisions.md) §26. Esc/read mode already renders sundown output.

## Phone Markdown mirror (future)

While editing on e-ink, show the open note rendered as Markdown on the phone — read-only context, not a second editor. Tapping a word in that view could move the tablet cursor there and scroll that section into view on e-ink. Depends on cursor position over the socket (already published for the keyboard harness) and a rendered Markdown pane on the phone.

## Browser (future)

Bulk select and multi-delete; search across titles and bodies; HTTPS for native Share sheet; phone-side vault PIN UI (tablet-only PIN entry is shipped — [decisions.md](decisions.md) §31).

## Install and onboarding (future)

First-time setup still assumes SSH, systemd, and optional `gh`. Checklist: [todo-install-onboarding.md](todo-install-onboarding.md).
