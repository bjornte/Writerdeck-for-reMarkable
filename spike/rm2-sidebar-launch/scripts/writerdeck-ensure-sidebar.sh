#!/bin/sh
# On-device: keep XOVI sidebar hook armed while Writerdeck sidebar files exist.
# No-op if not installed (rM1). Safe before every systemctl start xochitl and on boot.
#
# OS gate: IMG_VERSION in tested-os.json "tested", or in "candidates" when
# accept-untested flag file is present.
set -eu

XOVI_ROOT=/home/root/xovi
QRR_DIR=$XOVI_ROOT/exthome/qt-resource-rebuilder
QMD=$QRR_DIR/writerdeck-sidebar.qmd
HASHTAB=$QRR_DIR/hashtab
TESTED_JSON=$QRR_DIR/tested-os.json
ACCEPT_FLAG=$QRR_DIR/accept-untested
DROP_DIR=/etc/systemd/system/xochitl.service.d
SRC_DIR=$XOVI_ROOT/services/xochitl.service

log() { echo "writerdeck-ensure-sidebar: $*" >&2; }

if [ ! -x "$XOVI_ROOT/start" ] || [ ! -f "$QMD" ] || [ ! -s "$HASHTAB" ]; then
  exit 0
fi
if [ ! -f "$XOVI_ROOT/extensions.d/qt-resource-rebuilder.so" ]; then
  exit 0
fi
if [ ! -f "$XOVI_ROOT/extensions.d/qt-command-executor.so" ]; then
  log "skip: qt-command-executor.so missing"
  exit 0
fi
if [ -f "$XOVI_ROOT/extensions.d/appload.so" ]; then
  log "skip: appload.so still active (run install-native.sh)"
  exit 0
fi
if [ ! -f "$TESTED_JSON" ]; then
  log "skip: missing $TESTED_JSON"
  exit 0
fi

OS_VER=
if [ -f /etc/os-release ]; then
  OS_VER=$(sed -n 's/^IMG_VERSION=//p' /etc/os-release | head -n 1 | tr -d '"')
fi
if [ -z "$OS_VER" ] && [ -f /etc/version ]; then
  OS_VER=$(grep -Eo '[0-9]+\.[0-9]+\.[0-9]+(\.[0-9]+)?' /etc/version | head -n 1 || true)
fi
if [ -z "$OS_VER" ]; then
  log "skip: could not read OS version"
  exit 0
fi

gate_ok=0
if command -v python3 >/dev/null 2>&1; then
  PY=python3
elif command -v python >/dev/null 2>&1; then
  PY=python
else
  PY=
fi
if [ -n "$PY" ]; then
  if $PY - "$TESTED_JSON" "$OS_VER" "$ACCEPT_FLAG" <<'PY'
import json, os, sys
path, ver, flag = sys.argv[1], sys.argv[2], sys.argv[3]
data = json.load(open(path))
tested = set(data.get("tested") or [])
cands = set(data.get("candidates") or [])
if ver in tested:
    raise SystemExit(0)
if os.path.isfile(flag) and ver in cands:
    raise SystemExit(0)
raise SystemExit(1)
PY
  then
    gate_ok=1
  fi
else
  # No python: allow only if version string appears under a "tested" section snippet.
  if grep -A50 '"tested"' "$TESTED_JSON" | grep -F "\"$OS_VER\"" >/dev/null 2>&1; then
    gate_ok=1
  elif [ -f "$ACCEPT_FLAG" ] && grep -A50 '"candidates"' "$TESTED_JSON" | grep -F "\"$OS_VER\"" >/dev/null 2>&1; then
    gate_ok=1
  fi
fi

if [ "$gate_ok" != 1 ]; then
  log "skip: OS $OS_VER not allowed (update tested-os.json or accept-untested)"
  exit 0
fi

# Already armed?
if [ -f "$DROP_DIR/00-xovi.conf" ] && grep -q 'LD_PRELOAD=/home/root/xovi/xovi.so' "$DROP_DIR/00-xovi.conf" 2>/dev/null; then
  if mount | grep -q " on $DROP_DIR "; then
    exit 0
  fi
fi

log "arming XOVI sidebar hook (os=$OS_VER)"

mkdir -p "$DROP_DIR"
umount -q "$DROP_DIR" 2>/dev/null || true
mount -t tmpfs tmpfs "$DROP_DIR"
if [ -d "$SRC_DIR" ]; then
  cp -ra "$SRC_DIR/." "$DROP_DIR/" 2>/dev/null || true
fi
cat > "$DROP_DIR/00-xovi.conf" <<EOF
[Service]
Environment="LD_PRELOAD=/home/root/xovi/xovi.so"
Environment="XOVI_ROOT=$SRC_DIR/"
EOF
systemctl daemon-reload

# Do not start xochitl over a live editor session.
if pidof Writerdeck >/dev/null 2>&1; then
  log "drop-in mounted; editor running -- xochitl not restarted"
  exit 0
fi

if systemctl is-active --quiet xochitl 2>/dev/null; then
  log "restarting xochitl with XOVI"
  systemctl restart xochitl
else
  log "drop-in mounted; xochitl inactive -- will pick up XOVI on next start"
fi

exit 0
