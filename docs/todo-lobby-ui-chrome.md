# Lobby chrome and string i18n

Labels, help copy, visual knobs, and language packs for the Lobby. Goal: change them over SSH without a rebuild ([decisions.md](decisions.md) §36). Parent: [../TODO.md](../TODO.md). Shortcuts: [todo-lobby-ui-shortcuts.md](todo-lobby-ui-shortcuts.md). Wishlist: [improvements.md](improvements.md).

## On disk

Main file: `/home/root/.Writerdeck/lobby-ui.json` — `language`, `visual`, optional `strings` overrides, `shortcuts`. Repo: [../config/lobby-ui.json](../config/lobby-ui.json).

Language packs: `/home/root/.Writerdeck/lobby-ui-i18n/<lang>.json` (flat string maps). Repo: [../config/lobby-ui-i18n/](../config/lobby-ui-i18n/). Supported: `en`, `no`, `es`, `de`, `fr`. Set `"language"` in `lobby-ui.json`. Pack load order: embedded English defaults → disk pack for that language → overlay `strings` in `lobby-ui.json`.

Korean, Chinese, and Japanese need CJK fonts on the tablet first — not string-only. See [improvements.md](improvements.md). Language picker in Settings is wishlist only; change language in JSON for now.

Keycap badge borders use `visual.badgeBorderColor` (separate from outer `borderColor`). Tab titles are `tabs.files` … `tabs.about` in the pack.

## Checklist — copy and labels

- [x] Tab titles (`tabs.files` … `tabs.about`).
- [x] Files: Prev / Next / page template; action labels; list marker and `[private]` suffix.
- [x] About: brand, tagline, document-count templates, version lines, open-source line (`home.tip` in pack).
- [x] Keyboard: Bluetooth / USB headlines and connected suffixes, bodies, Layout, PIN prefix, layout button labels.
- [x] Sync: title, status templates, banners, button states, footnote.
- [x] Settings: page title; section headers; Enable / Change PIN / Exit; PIN option labels and Wi-Fi warn; rotation labels.
- [x] Vault pad: all prompts; Cancel via `dialog.cancel`.
- [x] Errors: name-already-exists, Operation failed, wrong PIN, related fallbacks.

## Checklist — visual

- [x] Button / tab fills and selected fills.
- [x] Page / dialog backgrounds, dialog scrim, vault wash.
- [x] Corner radii (button, dialog, keycap, banner).
- [x] Layout leftovers: page-strip height, list row inset, tab-row extra height, dialog width fraction and padding.
- [x] Type sizes beyond label/badge (title, section, row, dialog title, banner, help).
- [x] `badgeBorderColor` for shortcut keycaps.

## Verify

- [x] Edit tablet JSON (language or a color), confirm `lobby-ui: loaded …` in the journal, spot-check. Prefer one Writerdeck deploy per wiring batch. After Lobby QML: `bash scripts/test-lobby-keyboard.sh`.

## Follow-up prompt

> Project: Writerdeck for reMarkable 1. Branch: `lobby-ui-chrome` (this repo) and matching Writerdeck-keywriter branch.
> Finish Lobby chrome + Latin i18n per [docs/todo-lobby-ui-chrome.md](todo-lobby-ui-chrome.md): check off any remaining checklist rows, refresh long `shortcuts.body` translations if English drifted, deploy fork binary, seed `lobby-ui-i18n` if missing, SSH-test `"language": "no"` and `badgeBorderColor`, run `bash scripts/test-lobby-keyboard.sh`. Do not ship ko/zh/ja until CJK fonts exist ([improvements.md](improvements.md)). Settings language picker stays wishlist.
> Read: architecture, decisions §36, terms (lobby-ui.json), lessons (lobby-ui reload). Device: `secrets/remarkable.local.env`.
