# USB keyboard layouts (Qt evdev `.qmap`)

Writerdeck reads USB keyboards through Qt evdev, not `loadkeys` / `setxkbmap`.
Each layout is a binary `.qmap` file loaded via `QT_QPA_EVDEV_KEYBOARD_PARAMETERS`
in `Writerdeck-launcher.sh`.

| File | Layout |
|------|--------|
| `us.qmap` | US QWERTY (default) |
| `no.qmap` | Norwegian (`no-latin1` — æ ø å, AltGr symbols) |

Owner picks the layout in the Lobby **Keyboard** tab (`settings.json` → `keyboardLayout`).
The choice applies on the next editor launch (Lobby or note).

## Regenerate

```bash
bash keymaps/generate.sh          # us + no (needs Docker / Colima)
bash keymaps/generate.sh de       # extra layout if src/i386/qwerty/de.map exists
```

Sources live under `keymaps/src/i386/` (from device `/usr/share/keymaps/i386/`).
Build uses vendored `kmap2qmap-main.cpp` (Qt 5.15) + `qevdevkeyboardhandler_p.h`.

Deploy: `bash scripts/deploy-keywriter.sh` copies `*.qmap` to `/home/root/keymaps/`.

`writerdeck-alt-arrows.inc` overrides Linux Alt+Left/Right (console switch) so Qt delivers arrow keys for word navigation instead of fake Escape.
