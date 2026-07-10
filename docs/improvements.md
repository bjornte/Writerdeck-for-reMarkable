# Improvements

## On reMarkable
* General
   * Power button should be functional
   * Some sleep logic should be included
   * Files could be encrypted
   * When WiFi and thus URL changes, syncing fails. Setting page on new URL misses the token. Manual fix is to paste it in again. Do the browsers' security architecture make it impossible to share the token across URLs?
* For USB-connected keyboards
   * Activate the program using a keyboard shortcut
   * Support additional keyboard locales, incl. Norwegian
* In Edit view
   * Page navigation and corresponding keyboard shortcuts generally work, but  leave a lot to be desired
   * More VS Code-like shortcuts & behaviour. E.g., when in indented list, newline could add spaces to match previous indentation
   * On the Das Keyboard, some Command-arrow combinations are confusing
   * Maybe also a cursor block to the left of the in-focus line
   * Interface for navigating between headlines
   * Status bar with affordances for: Title, terse confirmations, zooming, time, battery
   * ~~The rotate button (in the browser) seemingly does not work.~~ Done — in Settings → Display; affects all tablet screens.
* In reading view:
   * Paragraph distances could be greater. Postponed.
   * ~~should not automatically scroll to the bottom~~ (done, device-verified — `ensureVisible` gated to edit mode)
* In lobby:
   * Could include an affordace (button) to open the file picker (`ctrl-K`)
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
   * For later: Choose more fonts (for reMarkable edit view, and views in the browser. Maybe one font across all?) and later maybe other things
   * ~~In settings, when in a laptop browser, clicking outside the pop-up or ESC closes it. Also, small close button in top right.~~ (done 2026-07-10 — settings + sync overlays)
