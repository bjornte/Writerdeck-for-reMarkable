#!/usr/bin/env bash
# Compare rM2 lobby PNG text cap heights to rM1 reference crop.
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REF="${1:-$DIR/../../../docs/screenshots/writerdeck-2026-07-22-lobby-5-choose-kbd-crop.png}"
RM2="${2:-$DIR/../../../docs/screenshots/writerdeck-rm2-lobby-keyboard-$(date +%Y-%m-%d).png}"
python3 - "$REF" "$RM2" <<'PY'
import sys
from PIL import Image

ref_path, rm2_path = sys.argv[1:3]
ref = Image.open(ref_path).convert("L")
rm2 = Image.open(rm2_path).convert("L")

def cap(im, box):
    x0, y0, x1, y1 = box
    sub = im.crop((x0, y0, x1, y1))
    px = sub.load()
    rows = [y for y in range(sub.height) if any(px[x, y] < 128 for x in range(sub.width))]
    return rows[-1] - rows[0] + 1 if rows else None

ref_regions = {
    "tab": (620, 8, 860, 58),
    "section": (35, 92, 520, 132),
    "help": (35, 132, 920, 175),
}
# rM2 keyboard tab (portrait 1404x1872) — approximate bands
rm2_regions = {
    "tab": (230, 8, 470, 68),
    "section": (35, 195, 700, 250),
    "help": (35, 250, 1300, 320),
    "url": (200, 420, 1200, 480),
}

print("reference", ref.size)
for k, b in ref_regions.items():
    print(f"  ref_{k}_cap", cap(ref, b))

print("rm2", rm2.size)
for k, b in rm2_regions.items():
    print(f"  rm2_{k}_cap", cap(rm2, b))

rt = cap(ref, ref_regions["tab"])
mt = cap(rm2, rm2_regions["tab"])
if rt and mt:
    print(f"tab scale factor to match ref: {rt / mt:.2f}")
PY
