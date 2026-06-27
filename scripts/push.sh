#!/usr/bin/env bash
# scripts/push.sh -- one short command to stage + commit + push.
# Replaces typing `git add ... && git commit -m ... && git push` by hand.
#
# Usage:
#   bash scripts/push.sh                 # uses a default message
#   bash scripts/push.sh "your message"  # custom message
#
# Tip: make it even shorter with an alias (run once):
#   echo 'alias rmpush="bash ~/dev/writerdeck-for-remarkable/scripts/push.sh"' >> ~/.zshrc
#   then just:  rmpush "message"

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
cd "$REPO"

MSG="${1:-device recon/probe result}"

git add -A
if git diff --cached --quiet; then
  echo "Nothing to commit (working tree clean). Pushing any unpushed commits..."
else
  git commit -m "$MSG"
fi
git push
git --no-pager log -1 --format="pushed %h - %s"
