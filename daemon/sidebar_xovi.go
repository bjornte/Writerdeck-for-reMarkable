package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// ensureSidebarXOVI best-effort arms the rM2 XOVI sidebar hook before xochitl
// starts. No-op when the on-device script is missing (rM1 / not installed).
func ensureSidebarXOVI() {
	const path = "/home/root/writerdeck-ensure-sidebar.sh"
	if _, err := os.Stat(path); err != nil {
		return
	}
	cmd := exec.Command("/bin/sh", path)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()
	select {
	case err := <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: ensure-sidebar: %v\n", err)
		}
	case <-time.After(45 * time.Second):
		fmt.Fprintln(os.Stderr, "writerdeck-server: ensure-sidebar: timed out")
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}
}

// startXochitl arms the sidebar hook (if installed) then starts stock UI.
func startXochitl() {
	ensureSidebarXOVI()
	if err := exec.Command("systemctl", "start", "xochitl").Run(); err != nil {
		fmt.Fprintf(os.Stderr, "writerdeck-server: warning: start xochitl: %v\n", err)
	}
}
