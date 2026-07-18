# Install and onboarding

First-time install for visitors: download the [installer ZIP](https://github.com/bjornte/Writerdeck-for-reMarkable/releases/download/installer/Writerdeck-installer.zip) → `bash scripts/install.sh --start`. Reuses saved secrets; asks only for missing password / Wi-Fi / phone PIN length / optional GitHub notes sync (can open a prefilled token page). Downloads editor + server from Releases, deploys, health-checks, enables autostart, configures PIN + sync when saved.

Wipe + restore notes: [reinstall-cheatsheet.md](reinstall-cheatsheet.md). Remove Writerdeck only: `bash scripts/uninstall.sh`. Regression: `bash scripts/test-install-reuse.sh`. Rebuild slim ZIP: `bash scripts/pack-installer.sh`.

Checklist and remaining follow-up (boot bricking): [todo-install-onboarding.md](todo-install-onboarding.md). Wishlist: [../improvements.md](../improvements.md).
