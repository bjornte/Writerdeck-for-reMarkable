# secrets/

Local-only credentials for rM1-Writerdeck. Nothing real in here is committed to git.

## How it works
- `remarkable.local.env` holds the real secrets. It is gitignored (see the root [.gitignore](../.gitignore)), so it lives on your disk only and cannot be pushed.
- `remarkable.local.env.example` is the committed template — placeholders only, safe to share. Copy it to create your real file:
  ```powershell
  Copy-Item secrets\remarkable.local.env.example secrets\remarkable.local.env
  ```
- The bash device scripts (`bootstrap.sh`, `recon.sh`, `deploy.sh`, `deploy-keywriter.sh`) read `remarkable.local.env` via `_env.sh`.

## What's stored
| Key | Meaning |
|---|---|
| `RM_ROOT_PASSWORD` | reMarkable root SSH password. On device: Settings → Help → Copyrights and licenses → General information (scroll down). Regenerates after every firmware update. |
| `RM_HOST_USB` | Device IP over USB — always `10.11.99.1`. |
| `RM_HOST_WIFI` | Device IP on Wi-Fi. |

## Why this is "good enough" here
The threat model is low: this password is shown on the tablet's own settings screen to anyone holding it, and the device lives on your home LAN. The real danger is accidentally committing it to git and pushing it public — which the gitignore prevents. Plaintext in a gitignored file is a reasonable, pragmatic choice for this case. (Day-to-day access is by SSH key, installed once by `bootstrap.sh`; the password is kept only because firmware updates reset it.)
