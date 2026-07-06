# Improvements

## On reMarkable
* General
   * For USB-connected keyboards, it should be possible to activate the program using a keyboard shortcut
   * For USB-connected keyboards, support additional keyboard locales, e.g. Norwegian
   * Power button should be functional
   * Some sleep logic should be included
   * Files could be encrypted
   * When WiFi and thus URL changes, syncing fails. Setting page on new URL misses the token. Manual fix is to paste it in again. Do the browsers' security architecture make it impossible to share the token across URLs?
* In Edit view
   * Page navigation and corresponding keyboard shortcuts generally work, but  leave a lot to be desired
   * More VS Code-like shortcuts & behaviour. E.g., when in indented list, newline could add spaces to match previous indentation
   * On the Das Keyboard, some Command-arrow combinations are confusing
   * Maybe also a cursor block to the left of the in-focus line
   * Interface for navigating between headlines
   * Status bar with affordances for: Title, terse confirmations, zooming, time, battery
   * The rotate button (in the browser) seemingly does not work. The direct USB keyboard shortcut works, however.
* In reading view:
   * Paragraph distances could be greater. Postponed.
   * should not automatically scroll to the bottom
* In lobby:
   * UX bug: On boot, URL shows "?". Fix: Refresh at intervals until on Wi-Fi, then update URL
   *  Could include an affordace (button) to open the file picker (`ctrl-K`)
   * When the notes sync to GitHub is enabled, the repo URL should be listed in the Lobby 
* In settings:
   * possible to exit service entirely (e.g. for battery performance reasons)

## In browser

* On main screen
   * Settings button with label `Settings`, not just cog wheel
   * Dedicated `Sync` button, and separate sync UI from settings.
* On the settings screen
   * Possible to exit service entirely on reMarkable (e.g. for battery performance reasons)
   * When the notes sync to GitHub is enabled, there should be a link to this repo from the browser.
   * For later: Choose more fonts (for reMarkable edit view, and views in the browser. Maybe one font across all?) and later maybe other things
   * In settings, when in a laptop browser, clicking outside the pop-up or ESC closes it. Also, small close button in top right. Or use a whole page for settings, rather than a pop-up. Pop-ups are not great on mobile.
