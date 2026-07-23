# Writerdeck for reMarkable 1

A distraction-free word processor for the first generation reMarkable. Supports Bluetooth and USB keyboards. Optionally syncs your documents to a private GitHub repository of your choice. Optionally encrypts files. Saves files as Markdown.

Natively, the reMarkable 1 supports the "draw", "write by hand" and "read" use cases. With this app, "use as typewriter" is also supported.

Bluetooth keyboards pair to your phone and bridge over Wi-Fi. USB keyboards use an [OTG cable](https://en.wikipedia.org/wiki/USB_On-The-Go#OTG_micro_cables).

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="/img/Writerdeck-for-reMarkable-two-photos-dark-bg.jpg">
  <source media="(prefers-color-scheme: light)" srcset="/img/Writerdeck-for-reMarkable-two-photos-light-bg.jpg">
  <img alt="Two photos of Writerdeck for reMarkable 1" src="/img/Writerdeck-for-reMarkable-two-photos.jpg">
</picture>

The reMarkable 1 has a large e-ink screen and a quiet OS, but no word processor and no keyboard support. This fills the gap.

![Keyboard support](docs/screenshots/writerdeck-2026-07-22-lobby-5-choose-kbd-crop.png)

The project is heavily LLM-assisted and partly human-reviewed. Primary sources: Singleton’s [keywriter](https://github.com/dps/remarkable-keywriter) (forked as [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter)) and ideas from [crazy-cow](https://github.com/machinelevel/sp425-crazy-cow).

## Status

A usable appliance. As with any software there's improvement potential, but a whole lot is working as intended. Latest: five languages supported, both in UI and as USB keyboard layouts.

## Install

Need: a reMarkable **1** (not 2), a Mac or Linux computer on the **same Wi-Fi**, and the tablet **awake**.

1. [Download the installer](https://github.com/bjornte/Writerdeck-for-reMarkable/releases/download/installer/Writerdeck-installer.zip). Unzip it and open a terminal in that folder.

2. Run:

   ```bash
   bash scripts/install.sh --start
   ```

If Writerdeck is stuck or something looks wrong after a bad install:

```bash
systemctl disable --now writerdeck && systemctl start xochitl
```

### Optional GitHub sync

As suggested during install: Create a private, personal GitHub repo for your notes. Set a fine-grained access token with Contents read/write on that repo only. See `img/create-token.png` for details. When the install is finished, use the phone interface to submit the token. For security reasons, is never stored on the reMarkable.

## Starting and quitting

Launch by mashing both page buttons simultaneously. Quit by pressing the home/ middle button (from editor to documents, and from documents to main reMarkable UI). Esc toggles edit and preview inside Writerdeck, or launches to Lobby from the stock UI with a USB keyboard.

## For developers

Feel free to contribute in every way! As mentioned above, my fork of keywriter has [its own repo](https://github.com/bjornte/Writerdeck-keywriter). Maybe the improved modifier key / chord handling part can be a usable element for other keywriter-derived projects, too.

## Pieces

* Writerdeck — full-screen editor from our keywriter fork
* Writerdeck-server — e.g. serves the phone page
* Phone page — captures Bluetooth keyboard keystrokes

## License

[MIT](LICENSE) © 2026 Bjørn Tennøe. Keywriter is third-party with its own license.

## Elsewhere

[Post on Discord](https://discord.com/channels/385916768696139794/1529598069346406511/1529598069346406511), [Closed post on r/writerdeck](https://www.reddit.com/r/writerDeck/comments/1uhtsnu/)
