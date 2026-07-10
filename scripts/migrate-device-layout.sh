#!/usr/bin/env bash
# scripts/migrate-device-layout.sh -- One-time rename of on-device paths.
# SOURCE ONLY: defines migrate_device_layout() for deploy/install scripts.

migrate_device_layout() {
  local target="${1:-$RM_HOST}"
  rm_ssh "$(cat <<'EOS'
set -e
# Notes: edit/ -> Writerdeck-user-documents/
if [ -d /home/root/edit ] && [ ! -e /home/root/Writerdeck-user-documents ]; then
  mv /home/root/edit /home/root/Writerdeck-user-documents
  echo "  migrated notes: edit -> Writerdeck-user-documents"
elif [ -d /home/root/edit ] && [ -d /home/root/Writerdeck-user-documents ]; then
  cp -a /home/root/edit/. /home/root/Writerdeck-user-documents/
  rm -rf /home/root/edit
  echo "  merged notes: edit -> Writerdeck-user-documents"
fi
mkdir -p /home/root/Writerdeck-user-documents

# Settings: .rmkbd/ -> .Writerdeck/
if [ -d /home/root/.rmkbd ] && [ ! -e /home/root/.Writerdeck ]; then
  mv /home/root/.rmkbd /home/root/.Writerdeck
  echo "  migrated settings: .rmkbd -> .Writerdeck"
elif [ -f /home/root/.rmkbd/settings.json ] && [ ! -f /home/root/.Writerdeck/settings.json ]; then
  mkdir -p /home/root/.Writerdeck
  mv /home/root/.rmkbd/settings.json /home/root/.Writerdeck/
  rmdir /home/root/.rmkbd 2>/dev/null || true
  echo "  migrated settings.json"
fi
mkdir -p /home/root/.Writerdeck

# Stop legacy processes (BusyBox: kill by pidof).
for p in $(pidof rmkbd 2>/dev/null); do kill "$p" 2>/dev/null || true; done
for p in $(pidof keywriter 2>/dev/null); do kill "$p" 2>/dev/null || true; done
pkill -f /home/root/rmkbd 2>/dev/null || true
pkill -f /home/root/keywriter 2>/dev/null || true
pkill -f /home/root/Writerdeck-server 2>/dev/null || true
for p in $(pidof Writerdeck 2>/dev/null); do kill "$p" 2>/dev/null || true; done
sleep 0.5

# systemd: rm1-writerdeck -> writerdeck
if [ -f /etc/systemd/system/rm1-writerdeck.service ]; then
  systemctl disable --now rm1-writerdeck 2>/dev/null || true
  rm -f /etc/systemd/system/rm1-writerdeck.service
  echo "  removed legacy unit rm1-writerdeck.service"
fi
if [ -f /etc/systemd/system/rmnetwriter.service ]; then
  systemctl disable --now rmnetwriter 2>/dev/null || true
  rm -f /etc/systemd/system/rmnetwriter.service
  echo "  removed legacy unit rmnetwriter.service"
fi
systemctl daemon-reload 2>/dev/null || true

# Stale socket from old editor binary name.
rm -f /run/rmkbd.sock /run/Writerdeck.sock 2>/dev/null || true
EOS
)" "$target"
}
