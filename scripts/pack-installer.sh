#!/usr/bin/env bash
# scripts/pack-installer.sh -- Build a slim end-user installer ZIP.
#
# Includes only scripts + keymaps + secrets template needed for:
#   bash scripts/install.sh --start
# Binaries are downloaded from GitHub Releases at install time.
#
# Output: dist/Writerdeck-installer.zip
#
# Usage (from repo root):
#   bash scripts/pack-installer.sh
#
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
OUT_DIR="$REPO/dist"
STAGING="$OUT_DIR/Writerdeck-installer"
ZIP="$OUT_DIR/Writerdeck-installer.zip"

cd "$REPO"
rm -rf "$STAGING"
mkdir -p "$STAGING/scripts" "$STAGING/secrets" "$STAGING/keymaps" \
  "$STAGING/third_party/keywriter/dist"

SCRIPTS=(
  install.sh
  uninstall.sh
  ensure-secrets.sh
  preflight.sh
  bootstrap.sh
  fix-hostkey.sh
  fetch-keywriter-dist.sh
  fetch-server-dist.sh
  deploy-keywriter.sh
  deploy-rmkbd.sh
  install-service.sh
  configure-sync.sh
  _env.sh
  paths.sh
  migrate-device-layout.sh
  Writerdeck-launcher.sh
  writerdeck.service
  wd
)

for f in "${SCRIPTS[@]}"; do
  cp "$REPO/scripts/$f" "$STAGING/scripts/$f"
done
chmod +x "$STAGING"/scripts/*.sh "$STAGING/scripts/wd" 2>/dev/null || true

cp "$REPO/secrets/remarkable.local.env.example" \
  "$STAGING/secrets/remarkable.local.env.example"
cp "$REPO/keymaps/us.qmap" "$REPO/keymaps/no.qmap" "$STAGING/keymaps/"

cat > "$STAGING/third_party/keywriter/dist/README.md" <<'EOF'
Editor binaries land here at install time (`bash scripts/fetch-keywriter-dist.sh`).
Do not put secrets here. Rolling Release tag: keywriter.
EOF

cat > "$STAGING/README.md" <<'EOF'
# Writerdeck installer

Install Writerdeck on a reMarkable **1** from a Mac or Linux computer on the same Wi-Fi.

1. Keep the tablet awake.
2. In a terminal in this folder:

```bash
bash scripts/install.sh --start
```

It asks for the tablet Wi-Fi IP, root password, and optional GitHub notes sync. Binaries download from GitHub Releases (no Go required).

Remove later (keeps your notes on GitHub):

```bash
bash scripts/uninstall.sh
```

Full project and docs: https://github.com/bjornte/Writerdeck-for-reMarkable
EOF

# Prefer zip; fall back to tar.gz named .zip only if zip missing -- use real zip.
rm -f "$ZIP"
(
  cd "$OUT_DIR"
  if command -v zip >/dev/null 2>&1; then
    zip -rq "Writerdeck-installer.zip" "Writerdeck-installer"
  else
    echo "ERROR: zip command not found (install zip)." >&2
    exit 1
  fi
)

rm -rf "$STAGING"
ls -lh "$ZIP"
echo "Packed: $ZIP"
