//go:build !linux

package main

// Stubs for host-side unit tests. Device builds use the linux ioctl path.

func evdevGrab(fd uintptr) error { return nil }

func evdevUngrab(fd uintptr) error { return nil }
