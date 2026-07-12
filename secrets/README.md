# secrets/

Local credentials — gitignored. Nothing real is committed.

Copy `remarkable.local.env.example` → `remarkable.local.env`. Bash scripts read it via `_env.sh`.

| Key | Meaning |
|---|---|
| `RM_ROOT_PASSWORD` | Root SSH password (device settings screen; regenerates on OTA) |
| `RM_HOST_USB` | USB IP — `10.11.99.1` (dead on Mac; Wi-Fi is the path) |
| `RM_HOST_WIFI` | Wi-Fi IP (DHCP — reserve tablet MAC on router) |

Threat model: password is visible on the device; home LAN. Real risk is git leakage — gitignore prevents that. Day-to-day access is SSH key via `bootstrap.sh`.
