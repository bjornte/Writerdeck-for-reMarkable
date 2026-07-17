# secrets/

Local credentials — never committed. Copy `remarkable.local.env.example` to `remarkable.local.env`. Scripts read it through `_env.sh`.

RM_ROOT_PASSWORD — SSH password from the tablet settings screen; regenerates after OTA.

RM_HOST_WIFI — tablet address on your network. Prefer a DHCP reservation.

RM_HOST_USB — unused on the Mac; Wi-Fi is the path.

The password is already visible on the device. The real risk is committing it. Day-to-day access uses the SSH key from `bootstrap.sh`.
