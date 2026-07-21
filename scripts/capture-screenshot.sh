#!/usr/bin/env bash
# scripts/capture-screenshot.sh -- Grab the reMarkable 1 screen into docs/screenshots/.
#
# Usage:
#   bash scripts/capture-screenshot.sh              # label from what's on screen
#   bash scripts/capture-screenshot.sh lobby-files  # optional override
#   bash scripts/capture-screenshot.sh --rotate 0
#
# Writes: docs/screenshots/writerdeck-YYYY-MM-DD-<label>.png
#
# Auto label (via localhost /api/test/editor-state on the tablet):
#   lobby | edit-<note> | read-<note> | stock
# Lobby tab name is not exposed yet, so Lobby shots are just "lobby".
#
# Streams /dev/fb0 over SSH and encodes PNG on the Mac (stdlib Python; no
# ffmpeg). Rotation follows Writerdeck settings.rotation (Lobby Settings),
# unless --rotate overrides. Mac alias (after install-alias.sh): rmshot

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"

ROTATE="" # empty = from tablet settings

usage() {
  echo "Usage: bash scripts/capture-screenshot.sh [--rotate 0|90|180|270] [label]" >&2
  echo "  label    optional slug; default from screen (lobby / edit-note / read-note / stock)" >&2
  echo "  --rotate force PNG clockwise degrees (default: from settings.rotation)" >&2
  echo "  -> docs/screenshots/writerdeck-YYYY-MM-DD-<label>.png" >&2
  exit 1
}

while [ $# -gt 0 ]; do
  case "$1" in
    --rotate)
      [ $# -ge 2 ] || usage
      ROTATE="$2"
      shift 2
      ;;
    -h|--help) usage ;;
    -*)
      err "unknown option: $1"
      usage
      ;;
    *) break ;;
  esac
done

LABEL="${1:-}"
if [ -n "$LABEL" ]; then
  shift || true
fi
[ $# -eq 0 ] || usage

slugify() {
  # note title / path -> filename-safe slug
  printf '%s' "$1" | python3 -c '
import re, sys
s = sys.stdin.read().strip()
s = re.sub(r"\.md(\.enc)?$", "", s, flags=re.I)
s = s.lower()
s = re.sub(r"[^a-z0-9]+", "-", s)
s = s.strip("-")[:48] or "note"
print(s)
'
}

# Infer label from Writerdeck editor state (localhost APIs; PIN not needed).
auto_label() {
  local st http status kw xo
  # wget prints body; capture HTTP code via a tiny wrapper.
  st="$(rm_ssh 'wget -qO- http://127.0.0.1:8000/api/test/editor-state 2>/dev/null' || true)"
  if [ -n "$st" ]; then
    printf '%s' "$st" | python3 -c '
import json, re, sys
st = json.load(sys.stdin)
def slug(s):
    s = re.sub(r"\.md(\.enc)?$", "", s or "", flags=re.I)
    s = re.sub(r"[^a-z0-9]+", "-", s.lower()).strip("-")[:48]
    return s or "note"
if st.get("isLobby") == 1:
    print("lobby")
elif st.get("currentFile"):
    kind = "edit" if st.get("mode") == 1 else "read"
    print(kind + "-" + slug(st["currentFile"]))
else:
    print("lobby")
'
    return 0
  fi
  status="$(rm_ssh 'wget -qO- http://127.0.0.1:8000/api/status 2>/dev/null' || true)"
  if [ -n "$status" ]; then
    note="$(printf '%s' "$status" | python3 -c 'import sys,json; print(json.load(sys.stdin).get("openNote") or "")' 2>/dev/null || true)"
    active="$(printf '%s' "$status" | python3 -c 'import sys,json; print("1" if json.load(sys.stdin).get("editorActive") else "0")' 2>/dev/null || true)"
    if [ "$active" = "1" ] && [ -n "$note" ]; then
      echo "edit-$(slugify "$note")"
      return 0
    fi
    if [ "$active" = "1" ]; then
      echo "lobby"
      return 0
    fi
  fi
  kw="$(rm_ssh 'pidof Writerdeck 2>/dev/null' || true)"
  xo="$(rm_ssh 'pidof xochitl 2>/dev/null' || true)"
  if [ -z "$kw" ] && [ -n "$xo" ]; then
    echo "stock"
    return 0
  fi
  echo "screen"
}

if [ -z "$LABEL" ]; then
  LABEL="$(auto_label)"
  echo "label=$LABEL (from screen)"
fi

if ! printf '%s' "$LABEL" | grep -Eq '^[A-Za-z0-9][A-Za-z0-9_-]*$'; then
  err "label must be letters/digits/-/_ (got: $LABEL)"
  exit 1
fi

# Read display rotation from the running server, else settings.json on disk.
# settings.rotation is how Writerdeck draws (0 / 90 / 180 / 270). The dump is
# always in panel portrait; undoing that transform makes an upright PNG:
#   png_cw = (360 - settings.rotation) % 360
read_device_rotation() {
  # Prefer settings.json over SSH: /api/settings needs the phone PIN.
  local deg="" json=""
  json="$(rm_ssh "cat '$DEVICE_SETTINGS_FILE'" 2>/dev/null || true)"
  if [ -n "$json" ]; then
    deg="$(printf '%s' "$json" | python3 -c 'import sys,json; print(json.load(sys.stdin)["rotation"])' 2>/dev/null || true)"
  fi
  printf '%s' "$deg"
}

if [ -z "$ROTATE" ]; then
  DEV_ROT="$(read_device_rotation)"
  if [ -z "$DEV_ROT" ]; then
    err "could not read settings.rotation from tablet"
    exit 1
  fi
  case "$DEV_ROT" in
    0|90|180|270) ;;
    *)
      # normalizeRotation may leave other multiples; fold to cardinal.
      DEV_ROT=$(( (DEV_ROT % 360 + 360) % 360 ))
      case "$DEV_ROT" in
        0|90|180|270) ;;
        *)
          err "unexpected settings.rotation: $DEV_ROT"
          exit 1
          ;;
      esac
      ;;
  esac
  ROTATE=$(( (360 - DEV_ROT) % 360 ))
  echo "settings.rotation=${DEV_ROT} -> PNG rotate ${ROTATE} CW"
else
  case "$ROTATE" in
    0|90|180|270) ;;
    *)
      err "rotate must be 0, 90, 180, or 270 (got: $ROTATE)"
      exit 1
      ;;
  esac
  echo "PNG rotate ${ROTATE} CW (--rotate override)"
fi

OUT_DIR="$DIR/../docs/screenshots"
mkdir -p "$OUT_DIR"
DATE="$(date +%Y-%m-%d)"
OUT="$OUT_DIR/writerdeck-${DATE}-${LABEL}.png"
if [ -e "$OUT" ]; then
  n=1
  while [ -e "$OUT_DIR/writerdeck-${DATE}-${LABEL}-${n}.png" ]; do
    n=$((n + 1))
  done
  OUT="$OUT_DIR/writerdeck-${DATE}-${LABEL}-${n}.png"
fi

# rM1 linuxfb: 16 bpp, virtual 1408x3840 (double buffer). Visible frame is the
# first 1408x1872. Qt uses size=1404x1872; keep full 1408 stride so the dump
# lines up with /dev/fb0.
FB_W=1408
FB_H=1872
FB_BYTES=$((FB_W * FB_H * 2))

TMP="$(mktemp "${TMPDIR:-/tmp}/writerdeck-fb.XXXXXX")"
trap 'rm -f "$TMP"' EXIT

echo "Capturing framebuffer from $RM_HOST..."
# Exact byte count via dd; stderr (records in/out) stays on the tablet.
rm_ssh "dd if=/dev/fb0 bs=$FB_BYTES count=1 2>/dev/null" >"$TMP"

GOT="$(wc -c <"$TMP" | tr -d '[:space:]')"
if [ "$GOT" -ne "$FB_BYTES" ]; then
  err "framebuffer dump size mismatch (got $GOT, expected $FB_BYTES)"
  exit 1
fi

python3 - "$TMP" "$OUT" "$FB_W" "$FB_H" "$ROTATE" <<'PY'
import struct, sys, zlib
from pathlib import Path

raw_path, out_path, w_s, h_s, rot_s = sys.argv[1:6]
w, h, rotate = int(w_s), int(h_s), int(rot_s)
fb = Path(raw_path).read_bytes()
need = w * h * 2
if len(fb) < need:
    sys.exit(f"short framebuffer: {len(fb)} < {need}")
fb = fb[:need]

# RGB565 little-endian -> 8-bit grayscale (e-ink looks monochrome).
pixels = bytearray(w * h)
for y in range(h):
    off = y * w * 2
    row = y * w
    for x in range(w):
        p = fb[off + 2 * x] | (fb[off + 2 * x + 1] << 8)
        r = ((p >> 11) & 31) * 255 // 31
        g = ((p >> 5) & 63) * 255 // 63
        b = (p & 31) * 255 // 31
        pixels[row + x] = (r * 30 + g * 59 + b * 11) // 100

def rot90(src: bytearray, sw: int, sh: int):
    dw, dh = sh, sw
    dst = bytearray(dw * dh)
    for y in range(sh):
        for x in range(sw):
            # clockwise: (x, y) -> (sh - 1 - y, x)
            dst[x * dw + (sh - 1 - y)] = src[y * sw + x]
    return dst, dw, dh

if rotate == 90:
    pixels, w, h = rot90(pixels, w, h)
elif rotate == 180:
    pixels, w, h = rot90(*rot90(pixels, w, h))
elif rotate == 270:
    pixels, w, h = rot90(*rot90(*rot90(pixels, w, h)))

rows = []
for y in range(h):
    row = bytearray(1 + w)
    row[1:] = pixels[y * w : (y + 1) * w]
    rows.append(row)

def chunk(tag: bytes, data: bytes) -> bytes:
    return (
        struct.pack(">I", len(data))
        + tag
        + data
        + struct.pack(">I", zlib.crc32(tag + data) & 0xFFFFFFFF)
    )

ihdr = struct.pack(">IIBBBBB", w, h, 8, 0, 0, 0, 0)  # 8-bit grayscale
png = (
    b"\x89PNG\r\n\x1a\n"
    + chunk(b"IHDR", ihdr)
    + chunk(b"IDAT", zlib.compress(b"".join(rows), 9))
    + chunk(b"IEND", b"")
)
Path(out_path).write_bytes(png)
print(out_path)
PY

echo "Wrote $OUT"
