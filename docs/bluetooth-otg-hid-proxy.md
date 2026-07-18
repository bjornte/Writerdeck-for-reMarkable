# Bluetooth keyboard over USB OTG

The reMarkable 1 has no Bluetooth. Writerdeck already accepts wired USB keyboards through a Micro-USB OTG adapter. A wireless Bluetooth keyboard needs a USB stick that pairs to the keyboard itself and presents as a normal wired USB keyboard (HID). Ordinary PC Bluetooth dongles (TP-Link, UGREEN, and the like) will not work — they need an operating-system Bluetooth stack the tablet does not have.

## Combined plug

No single Micro-USB stick was found that does OTG and Bluetooth HID proxy in one piece. Plan on two parts.

## OTG adapter

A rigid Micro-USB male to USB-A female OTG block is easy to buy (Delock 65549, Amazon “Micro USB OTG” plugs, Elkjøp’s compact OTG adapter). Prefer a solid block with no cable.

## HID proxy hardware

Ready-made CSR8510 “HID proxy” dongles that behave like a Logitech-style receiver for Bluetooth are nearly gone from retail. Specialist sellers of that exact part have been out of stock.

What you can buy today are programmable USB boards that become a HID proxy after you load bridge firmware once on a computer:

- Nordic nRF52840 Dongle (USB-A stick; plugs into the OTG adapter). Firmware example: [zephyr-ble-hid-usb-translator](https://github.com/charlesmst/zephyr-ble-hid-usb-translator).
- Seeed XIAO nRF52840 (smaller board, USB-C). Same chip family; needs a USB-C link to the tablet.
- Raspberry Pi Pico W with [HidProxy](https://github.com/MatanBright/HidProxy) (Bluetooth Classic only).

Power is not a concern if the tablet already runs a Das Keyboard Ultimate over OTG. Active draw for an nRF52840 bridge is a small fraction of a mechanical keyboard.

## XIAO size vs plug size

The XIAO board is the smallest of the three (~21×17.5 mm). The Nordic dongle is larger as a board but often smaller in practice at the tablet, because it plugs straight into USB-A OTG. The XIAO needs Micro-USB OTG plus a USB-C connection.

A rigid Micro-USB male to USB-C male block that would join tablet and XIAO directly is rare. Catalog listings that claim it (e.g. Networx AD-USB-MMCF) are weak: shared variant pages, conflicting lengths, little photo proof, and no confirmed OTG ID-pin wiring for host mode. Do not rely on those links. Prefer a known Micro-USB OTG to USB-A block, or a short cable sold explicitly as Micro-USB OTG host to USB-C male.

## Keyboard type

BLE keyboards suit the Nordic/XIAO bridge path. Older Bluetooth Classic-only keyboards need a different bridge (e.g. Pico W HidProxy). Many “wireless” keyboards include a proprietary 2.4 GHz USB receiver instead of Bluetooth; that receiver is already a driverless HID device and needs no proxy — OTG adapter plus receiver is enough.

## Product path today

Writerdeck’s supported wireless path remains the phone: Bluetooth pairs to the phone, which forwards keys to the tablet. Direct tablet Bluetooth via HID proxy is possible with the hardware above, but it is a DIY firmware step, not a finished consumer plug.
