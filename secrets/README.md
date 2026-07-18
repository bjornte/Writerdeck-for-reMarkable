# secrets/

Local credentials — never committed. `bash scripts/install.sh` (via `ensure-secrets.sh`) creates `remarkable.local.env` from the example and asks only for missing fields. Saved values are reused on the next install.

RM_ROOT_PASSWORD — SSH password from the tablet settings screen; regenerates after OTA.

RM_HOST_WIFI — tablet address on your network. Prefer a DHCP reservation.

RM_HOST_USB — unused on the Mac; Wi-Fi is the path.

PIN_DIGITS — phone connection PIN length: `6`, `4`, or `none`. Applied to the tablet after install.

SYNC_REPO — optional `owner/repo` for private GitHub notes sync.

GH_TOKEN — optional fine-grained token (Contents read/write). Stored on the computer only; pushed into tablet RAM by `configure-sync.sh`. Never written to the tablet’s disk.

SYNC_SKIP=1 — you declined sync; installer will not ask again until you clear this.

The password is already visible on the device. The real risk is committing secrets. Day-to-day access uses the SSH key from `bootstrap.sh`.
