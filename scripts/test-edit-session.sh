#!/usr/bin/env bash
# scripts/test-edit-session.sh -- regression: Edit from browser keeps Writerdeck up.
#
# Calls POST /api/open (same as phone Edit), then asserts:
#   - Writerdeck process stays running for the hold window
#   - xochitl stays stopped while editing
#   - /api/status reports editorActive=true
#
# Exit 0 = pass, 1 = fail. Logs to docs/recon/test-edit-session-<ts>.txt
#
# Usage (repo root):
#   bash scripts/test-edit-session.sh
#   bash scripts/test-edit-session.sh 192.168.1.8
#   bash scripts/test-edit-session.sh -s vetle.md   # note to open

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"
REPO="$(cd "$DIR/.." && pwd)"

NOTE="vetle.md"
TARGET="$RM_HOST"
for arg in "$@"; do
    case "$arg" in
        -s) shift; NOTE="${1:-vetle.md}"; shift || true ;;
        -*) echo "Unknown option: $arg" >&2; exit 2 ;;
        *)  TARGET="$arg" ;;
    esac
done

RECON_DIR="$REPO/docs/recon"
mkdir -p "$RECON_DIR"
TS="$(date +%Y-%m-%dT%H-%M-%S)"
LOG="$RECON_DIR/test-edit-session-$TS.txt"
exec > >(tee "$LOG") 2>&1

echo "=== test-edit-session  $TS  target=$TARGET  note=$NOTE ==="

_restore() {
    echo
    echo "--- Cleanup (restore stock UI if editor still up) ---"
    ssh -o StrictHostKeyChecking=accept-new -o ConnectTimeout=8 "root@$TARGET" \
        'wget -qO- --post-data="" http://127.0.0.1:8000/api/shutdown 2>/dev/null || true
         sleep 1
         systemctl start xochitl 2>/dev/null || true' || true
}
trap _restore EXIT

HOLD="${HOLD_SECS:-8}"
REMOTE_SCRIPT=$(cat <<EOF
set -e
NOTE='$NOTE'
HOLD=$HOLD
fail=0
log=/tmp/test-edit-session-\$\$.log
: > "\$log"

echo "=== before ===" | tee -a "\$log"
pidof Writerdeck 2>/dev/null || echo no-writerdeck | tee -a "\$log"
pidof xochitl 2>/dev/null || echo no-xochitl | tee -a "\$log"
wget -qO- http://127.0.0.1:8000/api/status | tee -a "\$log"
echo | tee -a "\$log"

echo "=== POST /api/open ===" | tee -a "\$log"
if ! wget -qO- --post-data="{\"name\":\"\$NOTE\"}" --header="Content-Type: application/json" \
    http://127.0.0.1:8000/api/open >> "\$log" 2>&1; then
  echo "FAIL: /api/open HTTP error" | tee -a "\$log"
  fail=1
fi

i=1
while [ "\$i" -le "\$HOLD" ]; do
  echo "--- t=\${i}s ---" | tee -a "\$log"
  kw=\$(pidof Writerdeck 2>/dev/null || true)
  xo=\$(pidof xochitl 2>/dev/null || true)
  st=\$(wget -qO- http://127.0.0.1:8000/api/status 2>/dev/null || echo '{}')
  echo "Writerdeck=\${kw:-none} xochitl=\${xo:-none} status=\$st" | tee -a "\$log"
  if [ -z "\$kw" ]; then
    echo "FAIL: Writerdeck exited at t=\${i}s" | tee -a "\$log"
    fail=1
    break
  fi
  if [ -n "\$xo" ]; then
    echo "FAIL: xochitl running while editor should be up (t=\${i}s)" | tee -a "\$log"
    fail=1
    break
  fi
  case "\$st" in
    *'"editorActive":true'*) ;;
    *) echo "FAIL: editorActive not true at t=\${i}s" | tee -a "\$log"; fail=1; break ;;
  esac
  i=\$((i + 1))
  sleep 1
done

if [ "\$fail" -eq 0 ]; then
  echo "PASS: Writerdeck stayed up for \${HOLD}s; xochitl down; editorActive=true" | tee -a "\$log"
fi
cat "\$log"
exit "\$fail"
EOF
)

if ! ssh -o StrictHostKeyChecking=accept-new -o ConnectTimeout=8 "root@$TARGET" "$REMOTE_SCRIPT"; then
    echo
    echo "=== verdict: FAIL ==="
    echo "Full log: $LOG"
    exit 1
fi

echo
echo "=== verdict: PASS ==="
echo "Full log: $LOG"
