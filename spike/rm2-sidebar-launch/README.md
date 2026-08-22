# Spike: version-gated Writerdeck sidebar launch (rM2)

Goal: one native row in the stock left sidebar labeled Writerdeck that opens Lobby -- no phone, no keyboard, no AppLoad app grid.

Related: AppLoad issue 68 (hashtab lesson only), docs/decisions.md section 7 / 33.

## Approach

Learn AppLoad's sidebar QMLDiff hook, then own it:

| Piece | Role |
|-------|------|
| XOVI + qt-resource-rebuilder | Inject into xochitl (tethered; reboot = stock) |
| qt-command-executor | Let the sidebar row shell-out to wget /api/lobby |
| qml/writerdeck-sidebar.qmd | Inserts SidebarFilterItem titled Writerdeck after Integrations |
| tested-os.json | Refuse enable unless IMG_VERSION was smoke-tested |

AppLoad's .so is moved inactive on install. We do not ship or depend on its UI.

## Verdict on 3.27.3.0 AppLoad incompatible

Small. Missing hashtab after OTA/install. Rebuild hashtab, then start. Same as remagic issue 2.

## Demo

```bash
bash spike/rm2-sidebar-launch/scripts/install-native.sh
ACCEPT_UNTESTED=1 bash spike/rm2-sidebar-launch/scripts/enable.sh
```

On tablet: unlock, hamburger, Writerdeck.

```bash
bash spike/rm2-sidebar-launch/scripts/mark-tested.sh   # after smoke OK
bash spike/rm2-sidebar-launch/scripts/disable.sh       # or reboot
```

Legacy AppLoad-as-launcher path: install-demo.sh (kept for reference only).
