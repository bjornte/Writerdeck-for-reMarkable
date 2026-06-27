#!/usr/bin/env bash
# scripts/install-alias.sh -- one-time setup for short shell shortcuts.
# Adds aliases to your shell rc file, using this repo's actual location.
# Idempotent: re-running only adds aliases that are missing.
#
#   rmpush ["message"]   -> scripts/push.sh           (commit + push)
#   rmkw                 -> deploy-keywriter.sh -b      (binary-only fast deploy, ~1s)
#
# Usage (once):   bash scripts/install-alias.sh
# Then forever:   rmpush     /     rmkw

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"

# Pick the right rc file for the current shell (zsh is macOS default).
case "${SHELL:-}" in
  *zsh) RC="$HOME/.zshrc" ;;
  *bash) RC="$HOME/.bashrc" ;;
  *) RC="$HOME/.zshrc" ;;
esac
touch "$RC"

# name|alias-body pairs.
add_alias() {
  local name="$1" body="$2"
  if grep -qF "alias ${name}=" "$RC"; then
    echo "${name} alias already present in $RC -- skipping."
  else
    printf '\n# rM1-Writerdeck shortcut\nalias %s=%s\n' "$name" "'$body'" >> "$RC"
    echo "Added ${name} alias to $RC"
  fi
}

add_alias rmpush "bash \"$DIR/push.sh\""
add_alias rmkw   "bash \"$DIR/deploy-keywriter.sh\" -b"

echo
echo "Activate now with:   source $RC"
echo "Then just run:        rmpush   /   rmkw"
