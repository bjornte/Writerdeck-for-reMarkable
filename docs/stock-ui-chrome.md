# Stock UI chrome (button and headline template)

Measured from reMarkable’s General settings on an rM2 (software 2026-06-12, screen 1404×1872). The pair is Restart and Turn off. Use these numbers for Writerdeck Lobby buttons and section headlines so we match the stock UI.

Source grab: [screenshots/rm2-general-settings-2026-08-22.png](screenshots/rm2-general-settings-2026-08-22.png). Lobby look still lives in `lobby-ui.json` — this file is the target to aim that JSON (and QML) at.

## Headline (“General settings”)

Font: reMarkable Serif VF, Regular (variable face; webui ships `/usr/share/remarkable/webui/reMarkableSerif.woff2`). Converted TTF for Qt is in `config/fonts/reMarkableSerif.ttf` (family name `reMarkable Serif VF`). Closest stock system fallback: EB Garamond under `/usr/share/fonts/ttf/ebgaramond/`.

Size: 50 px. Color: black (`#000000`). Weight: Regular.

Ink box for that title was about 338×47 px.

## Body and list text

Row labels and help copy in General settings use the same face as button labels: Noto Sans UI Regular. Measured ink height is about 20 px; treating point size as 27 matches the Restart / Turn off labels on this panel DPI. Writerdeck Lobby maps `sectionPointSize`, `rowPointSize`, `labelPointSize`, and `bannerPointSize` to 27, with `helpPointSize` at 20 for secondary lines.

Noto Sans UI does not include ↵ (U+21B5). Shortcuts text and Enter-key badges use `font.families: ["Noto Sans UI", "Nimbus Sans"]` with `config/fonts/NimbusSans-Regular.otf` installed on the tablet (URW base35, symbol fallback only).

## Buttons (Restart / Turn off)

The two buttons are identical. They sit in one row, share the content column width, and flush to the column’s right edge.

| | |
| --- | ---: |
| Outer width | 379 px each |
| Outer height | 80 px |
| Gap between | 24 px |
| Corner radius | 0 (square) |
| Border | 2 px solid black (`#000000`) |
| Fill | white (`#FFFFFF`) |
| Label color | black (`#000000`) |

Content is an icon then a label, centered as a group.

| | |
| --- | ---: |
| Vertical padding | ~20 px (centers the row in 80 px) |
| Horizontal padding | ~108–111 px each side |
| Icon size | ~36–38 × 35–36 px |
| Icon-to-label gap | ~23–24 px |
| Label font | Noto Sans UI (on device: `/usr/share/fonts/ttf/noto/`) |
| Label size | ~27 px Regular |

Icons: Restart uses a circular arrow; Turn off uses a power symbol. Stroke weight is close to the 2 px border.

## Copy-paste defaults

```
Headline:
  fontFamily: "reMarkable Serif VF"   // device TTF: config/fonts/reMarkableSerif.ttf
  fontSize: 50
  color: "#000000"
  fontWeight: Normal

Body / list / button label:
  fontFamily: "Noto Sans UI"
  fontSize: 27                       // help lines: 20
  color: "#000000"

Button:
  height: 80
  width: 379                       // or (rowWidth - 24) / 2
  radius: 0
  border.width: 2
  border.color: "#000000"
  color: "#ffffff"
  spacing (pair): 24
  icon: 36×36, 24 px before label
  label.fontFamily: "Noto Sans UI"
  label.fontSize: 27
  label.color: "#000000"
  padding: center content (~20 vertical)
```

Lobby maps these into `config/lobby-ui.json` (`titlePointSize` 50, label/section/row/banner 27, `helpPointSize` 20, `btnRadius` 0, white fills, selected tab fill `#000000`). Font-sample rows keep `font.family: modelData.id` for the face samples.

## How this was measured

Framebuffer dump from the running stock UI (`xochitl`) via the XOVI framebuffer spy config string, then pixel bounds on the PNG. Font family and size were matched by rendering candidate faces against the cropped glyphs.
