#!/bin/sh
# Put rM2 USB-C back to gadget so the port can charge and talk to a laptop.
# No-op on rM1 (no ci_hdrc.0 role file). Host mode must not survive sleep or halt.
ROLE=/sys/bus/platform/devices/ci_hdrc.0/role
UDC=/sys/kernel/config/usb_gadget/g_ether/UDC
[ -f "$ROLE" ] || exit 0
cur=$(cat "$ROLE" 2>/dev/null | tr -d ' \n\r')
if [ "$cur" != gadget ]; then
  printf 'gadget\n' > "$ROLE" 2>/dev/null || true
fi
if [ -f "$UDC" ]; then
  bind=$(cat "$UDC" 2>/dev/null | tr -d ' \n\r')
  if [ -z "$bind" ]; then
    printf 'ci_hdrc.0\n' > "$UDC" 2>/dev/null || true
  fi
fi
exit 0
