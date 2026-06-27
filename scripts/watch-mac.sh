#!/usr/bin/env bash
# scripts/watch-mac.sh -- auto-sync the git bridge FROM the Mac side.
#
# Pulls everything; but auto-commits/pushes ONLY new device outputs under
# docs/recon/ (WATCH_PATH). This keeps the Mac safe to use: stray or unfinished
# edits anywhere else are never swept into a commit -- they are reported, not
# pushed. A macOS GUI notification (Notification Center banner) pops on each
# sync, and when a pull brings in new commits from the PC.
#
# GUI banners fire on: ARM (start), each push, each pull-applied, when it sees
# unpushed edits OUTSIDE docs/recon/, and STOP. They are transient (vanishing)
# by design -- you will see the bridge arm/disarm and every sync, so you can
# tell at a glance whether it is running.
#
# Commits use THIS clone's configured git identity (same as push.sh). Make sure
# your Mac git identity is your PERSONAL one (bjornte@gmail.com), not a work
# email:  git config user.email
#
# Usage:
#   bash scripts/watch-mac.sh        # poll every 15s (default)
#   bash scripts/watch-mac.sh 30     # poll every 30s
#   WATCH_PATH=docs bash scripts/watch-mac.sh   # widen the auto-push scope
#   Ctrl-C to stop (you get a STOP banner).

set -uo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
cd "$REPO"

INTERVAL="${1:-15}"
# Only changes under this path are auto-committed/pushed. Everything else is
# reported but left alone. Override with the WATCH_PATH env var.
WATCH_PATH="${WATCH_PATH:-docs/recon}"
TITLE="rmwatch (Mac)"

note() {  # note "message" [sound-name]
  local msg="$1" snd="${2:-}"
  if [ -n "$snd" ]; then
    osascript -e "display notification \"$msg\" with title \"$TITLE\" sound name \"$snd\"" >/dev/null 2>&1 || true
  else
    osascript -e "display notification \"$msg\" with title \"$TITLE\"" >/dev/null 2>&1 || true
  fi
}

cleanup() {
  note "STOPPED - bridge is NO LONGER syncing" "Basso"
  echo
  echo "rmwatch (Mac): stopped."
  exit 0
}
trap cleanup INT TERM

note "ARMED - watching for changes" "Glass"
echo "rmwatch (Mac): polling every ${INTERVAL}s. Ctrl-C to stop."
echo "  auto-push scope: $WATCH_PATH/ (edits elsewhere are reported, not pushed)"
echo "  identity: $(git config user.name) <$(git config user.email)>"

last_other_sig=""
while true; do
  # 1) Stay in sync with the PC. autostash protects any in-flight edits.
  before="$(git rev-parse HEAD 2>/dev/null || true)"
  git pull --rebase --autostash >/dev/null 2>&1 || true
  after="$(git rev-parse HEAD 2>/dev/null || true)"
  if [ -n "$before" ] && [ "$before" != "$after" ]; then
    subj="$(git log -1 --format='%s' 2>/dev/null || true)"
    note "pulled: $subj" "Glass"
    echo "  pulled: $subj"
  fi

  # 2) Push ONLY new device outputs under WATCH_PATH (docs/recon by default).
  #    Scoped on purpose: never sweep up unrelated edits elsewhere.
  recon_dirty="$(git status --porcelain -- "$WATCH_PATH")"
  if [ -n "$recon_dirty" ]; then
    files="$(echo "$recon_dirty" | sed 's/^...//' | tr '\n' ' ')"
    git add -- "$WATCH_PATH"
    if git commit -m "auto(mac): sync $files" >/dev/null 2>&1 && git push >/dev/null 2>&1; then
      note "pushed: $files" "Glass"
      echo "  pushed: $files"
    else
      note "PUSH FAILED - see terminal" "Basso"
      echo "  push failed (run: git status)"
    fi
  fi

  # 2b) Report (do NOT commit) edits OUTSIDE the watched path. Re-warns only
  #     when the set changes, so it is informative without being noisy.
  other_sig="$(git status --porcelain | sed 's/^...//' | grep -v "^$WATCH_PATH/" | sort | tr '\n' ',' || true)"
  if [ -n "$other_sig" ] && [ "$other_sig" != "$last_other_sig" ]; then
    note "IGNORING edits outside $WATCH_PATH/ (not auto-pushed)" "Funk"
    echo "  note: leaving these alone (commit them yourself if intended):"
    git status --porcelain | sed 's/^...//' | grep -v "^$WATCH_PATH/" | sed 's/^/    /'
  fi
  last_other_sig="$other_sig"

  sleep "$INTERVAL"
done
