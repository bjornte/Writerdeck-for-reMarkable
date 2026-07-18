#!/usr/bin/env bash
# scripts/test-install-reuse.sh -- Regression for installer secret reuse +
# GitHub vault recovery (secret/vault). Does not print secrets.
#
# Prerequisites:
#   - secrets/remarkable.local.env has password, Wi-Fi, PIN_DIGITS, and
#     preferably SYNC_REPO + GH_TOKEN (needed for the vault half)
#   - tablet awake on RM_HOST_WIFI; writerdeck running
#
# Usage (repo root):
#   bash scripts/test-install-reuse.sh           # reuse check + vault recovery
#   bash scripts/test-install-reuse.sh --reuse   # only ensure-secrets reuse
#   bash scripts/test-install-reuse.sh --vault   # only vault wipe/restore
#
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"

DO_REUSE=1
DO_VAULT=1
for arg in "$@"; do
  case "$arg" in
    --reuse) DO_VAULT=0 ;;
    --vault) DO_REUSE=0 ;;
    -h|--help) sed -n '2,16p' "$0"; exit 0 ;;
  esac
done

if [ -t 2 ]; then R=$'\033[1;31m'; Z=$'\033[0m'; else R=''; Z=''; fi
err() { printf '%sERROR:%s %s\n' "$R" "$Z" "$*" >&2; }
ok()  { printf '  OK  %s\n' "$*"; }

_get() {
  sed -n -E "s/^[[:space:]]*$1[[:space:]]*=[[:space:]]*([^#]*).*/\1/p" \
    "$REPO/secrets/remarkable.local.env" | head -n1 | sed -E 's/[[:space:]]+$//'
}

echo "=== test-install-reuse  target=$RM_HOST ==="
echo

PASS="$(_get RM_ROOT_PASSWORD)"
WIFI="$(_get RM_HOST_WIFI)"
PIN="$(_get PIN_DIGITS)"
REPO_SLUG="$(_get SYNC_REPO)"
TOK="$(_get GH_TOKEN)"

if [ -z "$PASS" ] || [ -z "$WIFI" ]; then
  err "need RM_ROOT_PASSWORD + RM_HOST_WIFI in secrets"
  exit 1
fi

if [ "$DO_REUSE" -eq 1 ]; then
  echo "--- ensure-secrets reuse (non-interactive) ---"
  OUT="$(bash "$DIR/ensure-secrets.sh" </dev/null 2>&1 || true)"
  printf '%s\n' "$OUT"
  if ! printf '%s\n' "$OUT" | grep -q 'reusing saved values\|Nothing new to ask'; then
    # PIN/sync may still be optional-skip on empty PIN; password+wifi must not prompt.
    if printf '%s\n' "$OUT" | grep -q 'Password:\|Wi-Fi IP:'; then
      err "ensure-secrets asked for password/Wi-Fi despite saved values"
      exit 1
    fi
  fi
  ok "secret reuse path did not re-ask password/Wi-Fi"
  echo
fi

if [ "$DO_VAULT" -eq 1 ]; then
  if [ -z "$REPO_SLUG" ] || [ -z "$TOK" ]; then
    err "vault test needs SYNC_REPO + GH_TOKEN in secrets"
    exit 1
  fi
  if ! ping -c1 -W2 "$RM_HOST" >/dev/null 2>&1; then
    err "$RM_HOST unreachable"
    exit 1
  fi

  echo "--- GitHub secret/vault present ---"
  GH_TOKEN="$TOK" SYNC_REPO="$REPO_SLUG" python3 - <<'PY'
import json, urllib.request, base64, sys, os
tok = os.environ["GH_TOKEN"]
repo = os.environ["SYNC_REPO"]
req = urllib.request.Request(
    f"https://api.github.com/repos/{repo}/contents/secret/vault",
    headers={
        "Authorization": f"Bearer {tok}",
        "Accept": "application/vnd.github+json",
        "User-Agent": "writerdeck-test-install-reuse",
    },
)
try:
    with urllib.request.urlopen(req) as r:
        data = json.load(r)
except Exception as e:
    print("FAIL: cannot read secret/vault:", e, file=sys.stderr)
    sys.exit(1)
raw = base64.b64decode(data["content"])
obj = json.loads(raw)
for k in ("salt", "wrappedDataKey", "verifier"):
    if not obj.get(k):
        print("FAIL: incomplete secret/vault", file=sys.stderr)
        sys.exit(1)
print("github secret/vault OK")
open("/tmp/wd-vault-expected.json", "w").write(json.dumps(obj))
PY
  ok "GitHub has usable secret/vault"
  echo

  echo "--- wipe tablet vault fields, restore via sync ---"
  # shellcheck disable=SC2086
  ssh $RM_SSH_OPTS -o ConnectTimeout=8 "root@$RM_HOST" \
    "cat /home/root/.Writerdeck/settings.json" > /tmp/wd-settings-before.json
  python3 - <<'PY'
import json
from pathlib import Path
d = json.loads(Path("/tmp/wd-settings-before.json").read_text())
if not d.get("wrappedDataKey"):
    raise SystemExit("tablet has no vault material to test")
d["encryptionEnabled"] = False
d["vaultSalt"] = ""
d["vaultVerifier"] = ""
d["wrappedDataKey"] = ""
meta = d.get("syncMeta") or {}
meta.pop("secret/pin", None)
meta.pop("secret/vault", None)
d["syncMeta"] = meta
Path("/tmp/wd-settings-wiped.json").write_text(json.dumps(d, separators=(",", ":")))
PY
  # shellcheck disable=SC2086
  scp $RM_SSH_OPTS /tmp/wd-settings-wiped.json \
    "root@$RM_HOST:/home/root/.Writerdeck/settings.json" >/dev/null
  rm_ssh "systemctl restart writerdeck" >/dev/null
  sleep 2
  bash "$DIR/configure-sync.sh" >/dev/null
  rm_ssh 'wget -qO- --post-data="{}" --header="Content-Type: application/json" http://127.0.0.1:8000/api/sync/run >/dev/null' || true
  sleep 6
  # shellcheck disable=SC2086
  ssh $RM_SSH_OPTS -o ConnectTimeout=8 "root@$RM_HOST" \
    "cat /home/root/.Writerdeck/settings.json" > /tmp/wd-settings-after.json
  RM_HOST="$RM_HOST" python3 - <<'PY'
import json, sys, subprocess, os
from pathlib import Path
exp = json.loads(Path("/tmp/wd-vault-expected.json").read_text())
after = json.loads(Path("/tmp/wd-settings-after.json").read_text())
ok = after.get("encryptionEnabled") is True
for src, dst in (("salt", "vaultSalt"), ("wrappedDataKey", "wrappedDataKey"), ("verifier", "vaultVerifier")):
    if after.get(dst) != exp.get(src):
        print(f"FAIL: {dst} mismatch after sync", file=sys.stderr)
        ok = False
if not ok:
    host = os.environ["RM_HOST"]
    subprocess.check_call(["scp", "/tmp/wd-settings-before.json", f"root@{host}:/home/root/.Writerdeck/settings.json"])
    sys.exit(1)
print("vault fields restored from GitHub")
PY
  ok "vault recovery via secret/vault"
  rm -f /tmp/wd-vault-expected.json /tmp/wd-settings-before.json \
    /tmp/wd-settings-wiped.json /tmp/wd-settings-after.json
  echo
fi

echo "=== verdict: PASS ==="
