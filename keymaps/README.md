# USB keyboard layouts

Writerdeck loads USB keyboards through Qt’s `.qmap` files, not the usual Linux console maps.

Layouts (Lobby → Keyboard; applies on the next editor launch):

| Id | Layout |
|----|--------|
| us | US QWERTY (default) |
| no | Norwegian |
| es | Spanish |
| de | German |
| fr | French |

Regenerate with `bash keymaps/generate.sh` (Docker). Deploy copies `*.qmap` onto the tablet. The Alt+Left/Right include stops Linux console-switch keys from arriving as fake Escape.
