# USB keyboard layouts

Writerdeck loads USB keyboards through Qt’s `.qmap` files, not the usual Linux console maps.

us.qmap — US QWERTY (default). no.qmap — Norwegian. Pick in Lobby → Keyboard; it applies on the next editor launch.

Regenerate with `bash keymaps/generate.sh` (Docker). Deploy copies `*.qmap` onto the tablet. The Alt+Left/Right include stops Linux console-switch keys from arriving as fake Escape.
