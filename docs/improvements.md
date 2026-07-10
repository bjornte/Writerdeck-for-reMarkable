# Improvements

## On reMarkable
* General
   * ~~Power button should be functional~~ (implemented — save, sleep screen, suspend; wake restores editor; device test pending)
   * ~~Some sleep logic should be included~~ (same — power-button sleep path in `rmkbd`)
   * Files could be encrypted
* For USB-connected keyboards
   * Activate the program using a keyboard shortcut
   * Support additional keyboard locales, incl. Norwegian
* In Edit view
   * Page navigation and corresponding keyboard shortcuts generally work, but leave a lot to be desired
   * More VS Code-like shortcuts & behaviour. E.g., when in indented list, newline could add spaces to match previous indentation
   * ~~On the Das Keyboard, some Command-arrow combinations are confusing~~ (done, device-verified — Mac-style nav: Alt=word/paragraph, Cmd=line/doc, Shift=select, Home/End)
   * Maybe also a cursor block to the left of the in-focus line
   * Interface for navigating between headlines
   * Status bar with affordances for: Title, terse confirmations, zooming, time, battery
   * ~~The rotate button (in the browser) seemingly does not work.~~ Done — in Settings → Display; affects all tablet screens.
* In lobby:
   * Could include an affordace (button) to open the file picker (`ctrl-K`) — blocked until USB keyboard
   * ~~When the notes sync to GitHub is enabled, the repo URL should be listed in the Lobby~~ (done, device-verified — green `Sync: github.com/owner/repo` line when sync on)
* In settings:
   * ~~possible to exit service entirely (e.g. for battery performance reasons)~~ (done — Settings → Exit Writerdeck on phone)

## In browser

* On main screen
   * ~~Settings button with label `Settings`, not just cog wheel~~ (done 2026-07-10)
   * ~~Dedicated `Sync` button, and separate sync UI from settings.~~ (done 2026-07-10)
   * ~~Battery/Wi-Fi status in the top bar~~ (done — `GET /api/status`)
* On the settings screen
   * ~~Possible to exit service entirely on reMarkable (e.g. for battery performance reasons)~~ (done — Settings → Exit Writerdeck)
   * ~~When the notes sync to GitHub is enabled, there should be a link to this repo from the browser.~~ (done 2026-07-10 — in Sync panel)
   * ~~In settings, when in a laptop browser, clicking outside the pop-up or ESC closes it. Also, small close button in top right.~~ (done 2026-07-10 — settings + sync overlays)
