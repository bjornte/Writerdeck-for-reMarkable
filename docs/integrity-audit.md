# Document integrity

Last updated 17 Jul 2026. Contract: [architecture.md](architecture.md), [decisions.md](decisions.md).

Writerdeck is a typewriter. Integrity means your words land as real Markdown, stay readable, and are not silently emptied or replaced.

For normal solo use — one note at a time, Home or autosave — things are in good shape. That is not bank-grade durability or live collaboration, and we are not aiming for either.

## Still open

Pull power or hard-kill the process and you can lose up to about forty-five seconds since the last autosave. Sync skips the open note; other notes still sync. If disk changes under an open session, the phone can warn, but the tablet will not auto-reload — keep typing and your buffer can win. Clashes often leave a “(tablet copy)” file to sort by hand. PIN off on a shared Wi-Fi means anyone there can read or change notes.

## What already shipped

Edit lease for the open note. Plain Markdown saves with an HTML guard. Atomic writes. Autosave. Save before deploy and stop. Empty-push guard. Clash copies instead of silent overwrite. Vault disable refuses while encrypted user notes exist; a failed decrypt shows an error, not a blank editor.

## Before you ship note or sync changes

Ask: can this lose text, write wrong bytes, or overwrite without the user knowing? If yes, fix it or record an explicit acceptance here.
