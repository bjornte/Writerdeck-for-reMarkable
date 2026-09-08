# Spike: version-gated Writerdeck sidebar launch (rM2)

Goal: one native row in the stock left sidebar labeled Writerdeck that opens Lobby -- no phone, no keyboard, no AppLoad app grid.

Related: AppLoad issue 68 (hashtab lesson only), docs/decisions.md section 7 / 33.

## Approach

Learn AppLoad's sidebar QMLDiff hook, then own it:

| Piece | Role |
|-------|------|
| XOVI + qt-resource-rebuilder | Inject into xochitl (tethered; reboot = stock) |
| qt-command-executor | Let the sidebar row shell-out to wget /api/lobby |
| qml/writerdeck-sidebar.qmd | 3.27: SidebarFilterItem after Integrations |
| qml/writerdeck-sidebar-3.28.qmd | 3.28: ArkControls.SidebarItem after integrationsFoldout (candidate only) |
| tested-os.json | Refuse enable unless IMG_VERSION was smoke-tested |
| writerdeck-ensure-sidebar.sh | On-device: re-arm XOVI when Writerdeck starts / before xochitl |
| writerdeck-sidebar-launch.sh | Detach then POST /api/lobby (a blocking stop of xochitl from the click SEGVs 3.28 and reboots) |

AppLoad's .so is moved inactive on install. We do not ship or depend on its UI.

XOVI's systemd drop-in is tmpfs (reboot clears it). `install-native.sh` installs `/home/root/writerdeck-ensure-sidebar.sh`. `writerdeck.service` ExecStartPost and Writerdeck-server (before every `systemctl start xochitl`) call that script so the sidebar comes back with Writerdeck.

Sidebar icon: `resources/icon.png` packed as `qrc:/writerdeck/icon.png` (phone typewriter silhouette). Rebuild:

```bash
# after updating resources/icon.png
docker run --rm -v "$PWD/spike/rm2-sidebar-launch/resources:/work" -w /work \
  --platform linux/amd64 ghcr.io/toltec-dev/qt:v3.3 \
  rcc -binary -o writerdeck-icons.rcc writerdeck-icons.qrc
```

## Verdict on 3.27.3.0 AppLoad incompatible

Small. Missing hashtab after OTA/install. Rebuild hashtab, then start. Same as remagic issue 2.

## Demo

```bash
bash spike/rm2-sidebar-launch/scripts/install-native.sh
# optional manual arm (also happens on writerdeck start):
bash spike/rm2-sidebar-launch/scripts/enable.sh
bash scripts/install-service.sh --start   # picks up ExecStartPost ensure
```

On tablet: unlock, hamburger, Writerdeck. After reboot, start or restart `writerdeck` to re-arm.

```bash
bash spike/rm2-sidebar-launch/scripts/mark-tested.sh   # after smoke OK
bash spike/rm2-sidebar-launch/scripts/disable.sh       # or reboot
```

Legacy AppLoad-as-launcher path: install-demo.sh (kept for reference only).
