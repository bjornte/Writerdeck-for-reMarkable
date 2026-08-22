#!/usr/bin/env bash
# Capture rM2 /dev/fb0 to PNG (spike only). Uses pan offset + stride from sysfs.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
HOST="${RM2_HOST:-${RM2_HOST_WIFI:-10.11.99.1}}"
OUT_DIR="$ROOT/spike/rm2-official-sdk/screenshots"
mkdir -p "$OUT_DIR"
LABEL="${1:-hello}"
DATE="$(date +%Y-%m-%d)"
OUT="$OUT_DIR/rm2-spike-${DATE}-${LABEL}.png"
n=0
while [ -e "$OUT" ]; do
  n=$((n + 1))
  OUT="$OUT_DIR/rm2-spike-${DATE}-${LABEL}-${n}.png"
done

read -r FB_W FB_H FB_BPP FB_STRIDE FB_PAN_Y <<EOF
$(ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" 'bash -s' <<'REMOTE'
set -e
vs=$(cat /sys/class/graphics/fb0/virtual_size)
w=${vs%%,*}
h=${vs#*,}
bpp=$(cat /sys/class/graphics/fb0/bits_per_pixel)
stride=$(cat /sys/class/graphics/fb0/stride)
pan=$(cat /sys/class/graphics/fb0/pan)
pan_y=${pan#*,}
# Visible panel slice used by xochitl/hello (portrait 1404x1872 stored as 260-wide columns).
vis_w=260
vis_h=1408
echo "$vis_w $vis_h $bpp $stride $pan_y"
REMOTE
)
EOF

FB_BYTES=$((FB_STRIDE * FB_H))
VIS_BYTES=$((FB_STRIDE * FB_H))
# Start at pan y row (lines), each row stride bytes.
SKIP=$((FB_PAN_Y * FB_STRIDE))
CAP=$((FB_STRIDE * FB_H))

TMP="$(mktemp "${TMPDIR:-/tmp}/rm2-fb.XXXXXX")"
trap 'rm -f "$TMP"' EXIT

echo "Capturing fb0 from $HOST (skip=$SKIP cap=$CAP stride=$FB_STRIDE)..."
ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
  "dd if=/dev/fb0 bs=$FB_STRIDE skip=$((FB_PAN_Y + FB_H - 1408)) count=1408 2>/dev/null" >"$TMP" || \
ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
  "dd if=/dev/fb0 bs=$FB_STRIDE skip=$FB_PAN_Y count=1408 2>/dev/null" >"$TMP"

GOT="$(wc -c <"$TMP" | tr -d '[:space:]')"
EXP=$((FB_STRIDE * 1408))
if [ "$GOT" -ne "$EXP" ]; then
  echo "WARN: got $GOT bytes, expected $EXP; trying full pan slice" >&2
  ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
    "dd if=/dev/fb0 bs=$FB_STRIDE skip=$FB_PAN_Y count=$FB_H 2>/dev/null" >"$TMP"
  GOT="$(wc -c <"$TMP" | tr -d '[:space:]')"
fi

python3 - "$TMP" "$OUT" "$FB_STRIDE" "1408" <<'PY'
import struct, sys, zlib
from pathlib import Path

raw_path, out_path, stride_s, rows_s = sys.argv[1:5]
stride, rows = int(stride_s), int(rows_s)
raw = Path(raw_path).read_bytes()
need = stride * rows
raw = raw[:need]
# 32 bpp BGRA little-endian -> grayscale
w = stride // 4
h = min(rows, len(raw) // stride)
pixels = bytearray(w * h)
for y in range(h):
    off = y * stride
    for x in range(w):
        b = raw[off + 4 * x]
        g = raw[off + 4 * x + 1]
        r = raw[off + 4 * x + 2]
        pixels[y * w + x] = (r * 30 + g * 59 + b * 11) // 100

def rot90(src, sw, sh):
    dw, dh = sh, sw
    dst = bytearray(dw * dh)
    for y in range(sh):
        for x in range(sw):
            dst[x * dw + (sh - 1 - y)] = src[y * sw + x]
    return dst, dw, dh

# Portrait panel: rotate 90 CW for readable PNG (1408 x 260 -> 260 x 1408 display coords)
pixels, w, h = rot90(pixels, w, h)

rows_out = []
for y in range(h):
    row = bytearray(1 + w)
    row[1:] = pixels[y * w : (y + 1) * w]
    rows_out.append(row)

def chunk(tag, data):
    return struct.pack(">I", len(data)) + tag + data + struct.pack(">I", zlib.crc32(tag + data) & 0xFFFFFFFF)

ihdr = struct.pack(">IIBBBBB", w, h, 8, 0, 0, 0, 0)
png = (
    b"\x89PNG\r\n\x1a\n"
    + chunk(b"IHDR", ihdr)
    + chunk(b"IDAT", zlib.compress(b"".join(rows_out), 9))
    + chunk(b"IEND", b"")
)
Path(out_path).write_bytes(png)
print(out_path)
PY

echo "Wrote $OUT"
