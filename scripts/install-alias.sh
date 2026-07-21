#!/usr/bin/env bash
# scripts/install-alias.sh -- one-time setup for short shell shortcuts.
# Adds aliases to your shell rc file, using this repo's actual location.
# Idempotent: re-running only adds aliases that are missing; migrates old names.
#
#   rmpush ["message"]   -> scripts/push.sh              (commit + push)
#   rmkw                 -> deploy-keywriter.sh -b         (binary-only fast deploy)
#   rmlobby              -> lobby.sh                       (show Lobby on tablet)
#   rmshot [label]       -> capture-screenshot.sh [label]  (screen -> docs/screenshots/)
#
# Usage (once):   bash scripts/install-alias.sh
# Then forever:   rmpush  /  rmkw  /  rmlobby  /  rmshot
# (New aliases are not visible until you source the rc file or open a new terminal.)
#
# Names use the rm* prefix (same as rmpush / rmkw). Avoids wd / wds / wdl, which
# collide with Oh My Zsh warp-directory and other "working directory" tools.

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
ADDED=0
MIGRATED=0

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
  if grep -qE "^[[:space:]]*alias ${name}=" "$RC"; then
    echo "${name} alias already present in $RC -- skipping."
  else
    printf '\n# Writerdeck shortcut\nalias %s=%s\n' "$name" "'$body'" >> "$RC"
    echo "Added ${name} alias to $RC"
    ADDED=$((ADDED + 1))
  fi
}

# Comment out a prior alias line if present (keeps history in the rc file).
retire_alias() {
  local name="$1"
  if grep -qE "^[[:space:]]*alias ${name}=" "$RC"; then
    # macOS sed needs '' after -i; GNU sed accepts -i'' too with this form.
    if sed --version >/dev/null 2>&1; then
      sed -i -E "s/^([[:space:]]*alias ${name}=)/# migrated: \\1/" "$RC"
    else
      sed -i '' -E "s/^([[:space:]]*alias ${name}=)/# migrated: \\1/" "$RC"
    fi
    echo "Retired old ${name} alias in $RC (commented out)."
    MIGRATED=$((MIGRATED + 1))
  fi
}

# Old Mac shortcut was wd -> lobby.sh; that name fights Oh My Zsh wd (warp dir).
retire_alias wd

add_alias rmpush  "bash \"$DIR/push.sh\""
add_alias rmkw    "bash \"$DIR/deploy-keywriter.sh\" -b"
add_alias rmlobby "bash \"$DIR/lobby.sh\""
add_alias rmshot  "bash \"$DIR/capture-screenshot.sh\""

echo
if [ "$ADDED" -gt 0 ] || [ "$MIGRATED" -gt 0 ]; then
  echo ">>> Run this in the SAME terminal before using the new alias(es):"
  echo "    source $RC"
  echo
fi
echo "Then:  rmpush   /   rmkw   /   rmlobby   /   rmshot"
echo "(Until you source, use e.g.  bash scripts/lobby.sh  instead of  rmlobby)"
