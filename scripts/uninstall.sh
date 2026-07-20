#!/usr/bin/env bash
# scripts/uninstall.sh -- Remove Writerdeck from the tablet (not a factory reset).
#
# Stops the service, deletes binaries / notes dir / settings / Qt / keymaps /
# systemd unit, starts stock xochitl. Does not touch Mac secrets or GitHub.
# Root password stays the same.
#
# Usage (run from repo root on the Mac):
#   bash scripts/uninstall.sh
#   bash scripts/uninstall.sh --reboot
#   bash scripts/uninstall.sh 192.168.1.8
#
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"

DO_REBOOT=0
TARGET="$RM_HOST"
for arg in "$@"; do
  case "$arg" in
    --reboot) DO_REBOOT=1 ;;
    -h|--help)
      sed -n '2,14p' "$0"
      exit 0
      ;;
    *) TARGET="$arg" ;;
  esac
done

echo "=== uninstall.sh  target=$TARGET ==="
echo "Removes Writerdeck from the tablet. Mac secrets and GitHub notes stay."
echo

if ! ping -c1 -W2 "$TARGET" >/dev/null 2>&1; then
  echo "ERROR: $TARGET unreachable (ping failed)." >&2
  exit 1
fi
if ! rm_test_key "$TARGET"; then
  echo "ERROR: key-based SSH to root@$TARGET failed." >&2
  exit 1
fi

echo "--- stop service + processes ---"
rm_ssh 'systemctl disable --now writerdeck 2>/dev/null || true
for p in $(pidof Writerdeck-server 2>/dev/null); do kill -TERM "$p" 2>/dev/null || true; done
for p in $(pidof Writerdeck 2>/dev/null); do kill -TERM "$p" 2>/dev/null || true; done
sleep 1
for p in $(pidof Writerdeck-server 2>/dev/null); do kill -KILL "$p" 2>/dev/null || true; done
for p in $(pidof Writerdeck 2>/dev/null); do kill -KILL "$p" 2>/dev/null || true; done
true' "$TARGET"
echo "  stopped."

echo "--- remove unit + files ---"
rm_ssh 'rm -f /etc/systemd/system/writerdeck.service /etc/systemd/system/rm1-writerdeck.service
systemctl daemon-reload 2>/dev/null || true
rm -f /home/root/Writerdeck /home/root/Writerdeck.new \
  /home/root/Writerdeck-server /home/root/Writerdeck-server.new \
  /home/root/Writerdeck-launcher.sh /home/root/wd \
  /home/root/rmkbd /home/root/keywriter /home/root/launch-keywriter.sh \
  /run/Writerdeck.sock /run/rmkbd.sock \
  /tmp/kw.log /tmp/wd-server.log /home/root/qt5.tar.gz
rm -rf /home/root/Writerdeck-user-documents /home/root/.Writerdeck \
  /home/root/keymaps /home/root/qt5 \
  /home/root/edit /home/root/.rmkbd
echo removed' "$TARGET"

echo "--- start stock UI (xochitl) ---"
rm_ssh 'systemctl start xochitl 2>/dev/null || true; sleep 1; systemctl is-active xochitl 2>&1 || true' "$TARGET"

echo "--- verify ---"
left="$(rm_ssh 'for p in /home/root/Writerdeck /home/root/Writerdeck-server /home/root/Writerdeck-launcher.sh /home/root/wd /home/root/Writerdeck-user-documents /home/root/.Writerdeck /home/root/keymaps /home/root/qt5 /etc/systemd/system/writerdeck.service; do [ -e "$p" ] && echo "$p"; done' "$TARGET" || true)"
if [ -n "${left:-}" ]; then
  echo "WARNING: still present:" >&2
  echo "$left" >&2
  exit 1
fi
echo "  clean."

if [ "$DO_REBOOT" -eq 1 ]; then
  echo "--- reboot ---"
  rm_ssh 'sync; /sbin/reboot' "$TARGET" || true
  echo "  reboot sent."
else
  echo
  echo "Done. Optional: bash scripts/uninstall.sh --reboot"
fi
echo "Reinstall: bash scripts/install.sh --start"
