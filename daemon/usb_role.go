package main

import (
	"fmt"
	"os"
	"strings"
)

// rM2 USB-C dual-role. Absent on rM1 (Micro-USB OTG). Stock default is gadget:
// the tablet charges and talks to a laptop (10.11.99.1). Host mode is only for
// a USB keyboard/hub and does not persist across reboot.
const usbCRolePath = "/sys/bus/platform/devices/ci_hdrc.0/role"
const usbGadgetUDCPath = "/sys/kernel/config/usb_gadget/g_ether/UDC"
const usbGadgetUDCName = "ci_hdrc.0"

func usbCRolePresent() bool {
	_, err := os.Stat(usbCRolePath)
	return err == nil
}

// restoreUSBGadget puts USB-C back to charge/laptop mode. Call before suspend,
// on SIGTERM, and never leave host mode across a power-off. No-op on rM1.
func restoreUSBGadget() {
	if !usbCRolePresent() {
		return
	}
	role, _ := os.ReadFile(usbCRolePath)
	needRole := strings.TrimSpace(string(role)) != "gadget"
	bind, bindErr := os.ReadFile(usbGadgetUDCPath)
	needUDC := bindErr == nil && strings.TrimSpace(string(bind)) == ""
	if !needRole && !needUDC {
		return
	}
	if needRole {
		if err := os.WriteFile(usbCRolePath, []byte("gadget\n"), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: USB-C gadget: %v\n", err)
			return
		}
	}
	if needUDC {
		if err := os.WriteFile(usbGadgetUDCPath, []byte(usbGadgetUDCName+"\n"), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: USB-C gadget UDC: %v\n", err)
		}
	}
	fmt.Fprintln(os.Stderr, "writerdeck-server: USB-C restored to gadget (charge / laptop)")
}
