#!/usr/bin/env bash
# Check a spike screenshot is the right size and not blank. Usage: verify-png.sh <png> [min_dark]
set -euo pipefail
PNG="${1:?png path}"
MIN_DARK="${2:-400}"
python3 - "$PNG" "$MIN_DARK" <<'PY'
import sys
from pathlib import Path

try:
    from PIL import Image
except ImportError:
    sys.stderr.write("ERROR: python3 Pillow required (pip install Pillow)\n")
    sys.exit(2)

path = Path(sys.argv[1])
min_dark = int(sys.argv[2])
img = Image.open(path).convert("L")
w, h = img.size
if (w, h) != (1404, 1872):
    sys.stderr.write(f"ERROR: expected 1404x1872, got {w}x{h}\n")
    sys.exit(1)
data = img.getdata()
dark = sum(1 for p in data if p < 200)
mean = sum(data) / len(data)
if mean > 252 and dark < min_dark:
    sys.stderr.write(f"ERROR: mostly blank (mean={mean:.1f}, dark={dark})\n")
    sys.exit(1)
if dark < min_dark:
    sys.stderr.write(f"ERROR: too few dark pixels ({dark} < {min_dark})\n")
    sys.exit(1)
print(f"OK {w}x{h} mean={mean:.1f} dark={dark}")
PY
